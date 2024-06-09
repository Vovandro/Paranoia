package decoder

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

func Decode(from interface{}, to interface{}, tag string, weakly bool) error {
	return decode(from, reflect.ValueOf(to).Elem(), tag, weakly)
}

func decode(input interface{}, outVal reflect.Value, tag string, weakly bool) error {
	var inputVal reflect.Value

	if input != nil {
		inputVal = reflect.ValueOf(input)

		// We need to check here if input is a typed nil. Typed nils won't
		// match the "input == nil" below so we check that here.
		if inputVal.Kind() == reflect.Ptr && inputVal.IsNil() {
			input = nil
		}
	}

	if input == nil {
		// If the data is nil, then we don't set anything, unless ZeroFields is set
		// to true.
		outVal.Set(reflect.Zero(outVal.Type()))
		return nil
	}

	if !inputVal.IsValid() {
		// If the input value is invalid, then we just set the value
		// to be the zero value.
		outVal.Set(reflect.Zero(outVal.Type()))
		return nil
	}

	var err error
	outputKind := getKind(outVal)

	switch outputKind {
	case reflect.Bool:
		err = decodeBool(input, outVal, tag, weakly)
	case reflect.Interface:
		err = decodeBasic(input, outVal, tag, weakly)
	case reflect.String:
		err = decodeString(input, outVal, tag, weakly)
	case reflect.Int:
		if outVal.MethodByName("Nanoseconds").IsValid() {
			err = decodeDuration(input, outVal, tag, weakly)
		} else {
			err = decodeInt(input, outVal, tag, weakly)
		}
	case reflect.Uint:
		err = decodeUint(input, outVal, tag, weakly)
	case reflect.Float32:
		err = decodeFloat(input, outVal, tag, weakly)
	case reflect.Struct:
		err = decodeStruct(input, outVal, tag, weakly)
	case reflect.Map:
		err = decodeMap(input, outVal, tag, weakly)
	case reflect.Ptr:
		err = decodePtr(input, outVal, tag, weakly)
	case reflect.Slice:
		err = decodeSlice(input, outVal, tag, weakly)
	case reflect.Array:
		err = decodeArray(input, outVal, tag, weakly)
	case reflect.Func:
		err = decodeFunc(input, outVal, tag, weakly)
	default:
		// If we reached this point then we weren't able to decode it
		return fmt.Errorf("unsupported type: %s", outputKind)
	}

	return err
}

// This decodes a basic type (bool, int, string, etc.) and sets the
// value to "data" of that type.
func decodeBasic(data interface{}, val reflect.Value, tag string, weakly bool) error {
	if val.IsValid() && val.Elem().IsValid() {
		elem := val.Elem()

		// If we can't address this element, then its not writable. Instead,
		// we make a copy of the value (which is a pointer and therefore
		// writable), decode into that, and replace the whole value.
		copied := false
		if !elem.CanAddr() {
			copied = true

			// Make *T
			copy := reflect.New(elem.Type())

			// *T = elem
			copy.Elem().Set(elem)

			// Set elem so we decode into it
			elem = copy
		}

		// Decode. If we have an error then return. We also return right
		// away if we're not a copy because that means we decoded directly.
		if err := decode(data, elem, tag, weakly); err != nil || !copied {
			return err
		}

		// If we're a copy, we need to set te final result
		val.Set(elem.Elem())
		return nil
	}

	dataVal := reflect.ValueOf(data)

	// If the input data is a pointer, and the assigned type is the dereference
	// of that exact pointer, then indirect it so that we can assign it.
	// Example: *string to string
	if dataVal.Kind() == reflect.Ptr && dataVal.Type().Elem() == val.Type() {
		dataVal = reflect.Indirect(dataVal)
	}

	if !dataVal.IsValid() {
		dataVal = reflect.Zero(val.Type())
	}

	dataValType := dataVal.Type()
	if !dataValType.AssignableTo(val.Type()) {
		return fmt.Errorf(
			"expected type '%s', got '%s'",
			val.Type(), dataValType)
	}

	val.Set(dataVal)
	return nil
}

func decodeString(data interface{}, val reflect.Value, tag string, weakly bool) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)

	converted := true
	switch {
	case dataKind == reflect.String:
		val.SetString(dataVal.String())
	case dataKind == reflect.Bool && weakly:
		if dataVal.Bool() {
			val.SetString("1")
		} else {
			val.SetString("0")
		}
	case dataKind == reflect.Int && weakly:
		val.SetString(strconv.FormatInt(dataVal.Int(), 10))
	case dataKind == reflect.Uint && weakly:
		val.SetString(strconv.FormatUint(dataVal.Uint(), 10))
	case dataKind == reflect.Float32 && weakly:
		val.SetString(strconv.FormatFloat(dataVal.Float(), 'f', -1, 64))
	case dataKind == reflect.Slice && weakly,
		dataKind == reflect.Array && weakly:
		dataType := dataVal.Type()
		elemKind := dataType.Elem().Kind()
		switch elemKind {
		case reflect.Uint8:
			var uints []uint8
			if dataKind == reflect.Array {
				uints = make([]uint8, dataVal.Len(), dataVal.Len())
				for i := range uints {
					uints[i] = dataVal.Index(i).Interface().(uint8)
				}
			} else {
				uints = dataVal.Interface().([]uint8)
			}
			val.SetString(string(uints))
		default:
			converted = false
		}
	default:
		converted = false
	}

	if !converted {
		return fmt.Errorf(
			"expected type '%s', got unconvertible type '%s', value: '%v'",
			val.Type(), dataVal.Type(), data)
	}

	return nil
}

func decodeInt(data interface{}, val reflect.Value, tag string, weakly bool) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)
	dataType := dataVal.Type()

	switch {
	case dataKind == reflect.Int:
		val.SetInt(dataVal.Int())
	case dataKind == reflect.Uint:
		val.SetInt(int64(dataVal.Uint()))
	case dataKind == reflect.Float32:
		val.SetInt(int64(dataVal.Float()))
	case dataKind == reflect.Bool && weakly:
		if dataVal.Bool() {
			val.SetInt(1)
		} else {
			val.SetInt(0)
		}
	case dataKind == reflect.String && weakly:
		str := dataVal.String()
		if str == "" {
			str = "0"
		}

		i, err := strconv.ParseInt(str, 0, val.Type().Bits())
		if err == nil {
			val.SetInt(i)
		} else {
			return fmt.Errorf("cannot parse as int: %s", err)
		}
	case dataType.PkgPath() == "encoding/json" && dataType.Name() == "Number":
		jn := data.(json.Number)
		i, err := jn.Int64()
		if err != nil {
			return fmt.Errorf(
				"error decoding json.Number into: %s", err)
		}
		val.SetInt(i)
	default:
		return fmt.Errorf(
			"expected type '%s', got unconvertible type '%s', value: '%v'",
			val.Type(), dataVal.Type(), data)
	}

	return nil
}

func decodeDuration(data interface{}, val reflect.Value, tag string, weakly bool) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)
	dataType := dataVal.Type()

	switch {
	case dataKind == reflect.Int:
		val.SetInt(dataVal.Int())
	case dataKind == reflect.Uint:
		val.SetInt(int64(dataVal.Uint()))
	case dataKind == reflect.Float32:
		val.SetInt(int64(dataVal.Float()))
	case dataKind == reflect.Bool && weakly:
		if dataVal.Bool() {
			val.SetInt(1)
		} else {
			val.SetInt(0)
		}
	case dataKind == reflect.String && weakly:
		str := dataVal.String()
		if str == "" {
			str = "0"
		}

		d, err := time.ParseDuration(str)

		if err == nil {
			val.SetInt(int64(d))
		} else {
			return fmt.Errorf("cannot parse as duration: %s", err)
		}
	case dataType.PkgPath() == "encoding/json" && dataType.Name() == "Number":
		jn := data.(json.Number)
		i, err := jn.Int64()
		if err != nil {
			return fmt.Errorf(
				"error decoding json.Number into: %s", err)
		}
		val.SetInt(i)
	default:
		return fmt.Errorf(
			"expected type '%s', got unconvertible type '%s', value: '%v'",
			val.Type(), dataVal.Type(), data)
	}

	return nil
}

func decodeUint(data interface{}, val reflect.Value, tag string, weakly bool) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)
	dataType := dataVal.Type()

	switch {
	case dataKind == reflect.Int:
		i := dataVal.Int()
		if i < 0 && !weakly {
			return fmt.Errorf("cannot parse, %d overflows uint", i)
		}
		val.SetUint(uint64(i))
	case dataKind == reflect.Uint:
		val.SetUint(dataVal.Uint())
	case dataKind == reflect.Float32:
		f := dataVal.Float()
		if f < 0 && !weakly {
			return fmt.Errorf("cannot parse, %f overflows uint", f)
		}
		val.SetUint(uint64(f))
	case dataKind == reflect.Bool && weakly:
		if dataVal.Bool() {
			val.SetUint(1)
		} else {
			val.SetUint(0)
		}
	case dataKind == reflect.String && weakly:
		str := dataVal.String()
		if str == "" {
			str = "0"
		}

		i, err := strconv.ParseUint(str, 0, val.Type().Bits())
		if err == nil {
			val.SetUint(i)
		} else {
			return fmt.Errorf("cannot parse as uint: %s", err)
		}
	case dataType.PkgPath() == "encoding/json" && dataType.Name() == "Number":
		jn := data.(json.Number)
		i, err := strconv.ParseUint(string(jn), 0, 64)
		if err != nil {
			return fmt.Errorf(
				"error decoding json.Number into: %s", err)
		}
		val.SetUint(i)
	default:
		return fmt.Errorf(
			"expected type '%s', got unconvertible type '%s', value: '%v'",
			val.Type(), dataVal.Type(), data)
	}

	return nil
}

func decodeBool(data interface{}, val reflect.Value, tag string, weakly bool) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)

	switch {
	case dataKind == reflect.Bool:
		val.SetBool(dataVal.Bool())
	case dataKind == reflect.Int && weakly:
		val.SetBool(dataVal.Int() != 0)
	case dataKind == reflect.Uint && weakly:
		val.SetBool(dataVal.Uint() != 0)
	case dataKind == reflect.Float32 && weakly:
		val.SetBool(dataVal.Float() != 0)
	case dataKind == reflect.String && weakly:
		b, err := strconv.ParseBool(dataVal.String())
		if err == nil {
			val.SetBool(b)
		} else if dataVal.String() == "" {
			val.SetBool(false)
		} else {
			return fmt.Errorf("cannot parse as bool: %s", err)
		}
	default:
		return fmt.Errorf(
			"expected type '%s', got unconvertible type '%s', value: '%v'",
			val.Type(), dataVal.Type(), data)
	}

	return nil
}

func decodeFloat(data interface{}, val reflect.Value, tag string, weakly bool) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)
	dataType := dataVal.Type()

	switch {
	case dataKind == reflect.Int:
		val.SetFloat(float64(dataVal.Int()))
	case dataKind == reflect.Uint:
		val.SetFloat(float64(dataVal.Uint()))
	case dataKind == reflect.Float32:
		val.SetFloat(dataVal.Float())
	case dataKind == reflect.Bool && weakly:
		if dataVal.Bool() {
			val.SetFloat(1)
		} else {
			val.SetFloat(0)
		}
	case dataKind == reflect.String && weakly:
		str := dataVal.String()
		if str == "" {
			str = "0"
		}

		f, err := strconv.ParseFloat(str, val.Type().Bits())
		if err == nil {
			val.SetFloat(f)
		} else {
			return fmt.Errorf("cannot parse as float: %s", err)
		}
	case dataType.PkgPath() == "encoding/json" && dataType.Name() == "Number":
		jn := data.(json.Number)
		i, err := jn.Float64()
		if err != nil {
			return fmt.Errorf(
				"error decoding json.Number into: %s", err)
		}
		val.SetFloat(i)
	default:
		return fmt.Errorf(
			"expected type '%s', got unconvertible type '%s', value: '%v'",
			val.Type(), dataVal.Type(), data)
	}

	return nil
}

func decodeMap(data interface{}, val reflect.Value, tag string, weakly bool) error {
	valType := val.Type()
	valKeyType := valType.Key()
	valElemType := valType.Elem()

	// By default we overwrite keys in the current map
	valMap := val

	// If the map is nil or we're purposely zeroing fields, make a new map
	if valMap.IsNil() {
		// Make a new map to hold our result
		mapType := reflect.MapOf(valKeyType, valElemType)
		valMap = reflect.MakeMap(mapType)
	}

	// Check input type and based on the input type jump to the proper func
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	switch dataVal.Kind() {
	case reflect.Map:
		return decodeMapFromMap(dataVal, val, valMap, tag, weakly)

	case reflect.Struct:
		return decodeMapFromStruct(dataVal, val, valMap, tag, weakly)

	case reflect.Array, reflect.Slice:
		if weakly {
			return decodeMapFromSlice(dataVal, val, valMap, tag, weakly)
		}

		fallthrough

	default:
		return fmt.Errorf("expected a map, got '%s'", dataVal.Kind())
	}
}

func decodeMapFromSlice(dataVal reflect.Value, val reflect.Value, valMap reflect.Value, tag string, weakly bool) error {
	// Special case for BC reasons (covered by tests)
	if dataVal.Len() == 0 {
		val.Set(valMap)
		return nil
	}

	for i := 0; i < dataVal.Len(); i++ {
		err := decode(dataVal.Index(i).Interface(), val, tag, weakly)
		if err != nil {
			return err
		}
	}

	return nil
}

func decodeMapFromMap(dataVal reflect.Value, val reflect.Value, valMap reflect.Value, tag string, weakly bool) error {
	valType := val.Type()
	valKeyType := valType.Key()
	valElemType := valType.Elem()

	// Accumulate errors
	errors := make([]string, 0)

	// If the input data is empty, then we just match what the input data is.
	if dataVal.Len() == 0 {
		if dataVal.IsNil() {
			if !val.IsNil() {
				val.Set(dataVal)
			}
		} else {
			// Set to empty allocated value
			val.Set(valMap)
		}

		return nil
	}

	for _, k := range dataVal.MapKeys() {
		// First decode the key into the proper type
		currentKey := reflect.Indirect(reflect.New(valKeyType))
		if err := decode(k.Interface(), currentKey, tag, weakly); err != nil {
			errors = append(errors, err.Error())
			continue
		}

		// Next decode the data into the proper type
		v := dataVal.MapIndex(k).Interface()
		currentVal := reflect.Indirect(reflect.New(valElemType))
		if err := decode(v, currentVal, tag, weakly); err != nil {
			errors = append(errors, err.Error())
			continue
		}

		valMap.SetMapIndex(currentKey, currentVal)
	}

	// Set the built up map to the value
	val.Set(valMap)

	// If we had errors, return those
	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "\n"))
	}

	return nil
}

func decodeMapFromStruct(dataVal reflect.Value, val reflect.Value, valMap reflect.Value, tag string, weakly bool) error {
	typ := dataVal.Type()
	for i := 0; i < typ.NumField(); i++ {
		// Get the StructField first since this is a cheap operation. If the
		// field is unexported, then ignore it.
		f := typ.Field(i)
		if f.PkgPath != "" {
			continue
		}

		// Next get the actual value of this field and verify it is assignable
		// to the map value.
		v := dataVal.Field(i)
		if !v.Type().AssignableTo(valMap.Type().Elem()) {
			return fmt.Errorf("cannot assign type '%s' to map value field of type '%s'", v.Type(), valMap.Type().Elem())
		}

		tagValue := f.Tag.Get(tag)
		keyName := f.Name

		if tagValue == "" {
			continue
		}

		// If Squash is set in the config, we squash the field down.
		squash := v.Kind() == reflect.Struct && f.Anonymous

		v = dereferencePtrToStructIfNeeded(v, tag)

		// Determine the name of the key in the map
		if index := strings.Index(tagValue, ","); index != -1 {
			if tagValue[:index] == "-" {
				continue
			}
			// If "omitempty" is specified in the tag, it ignores empty values.
			if strings.Index(tagValue[index+1:], "omitempty") != -1 && isEmptyValue(v) {
				continue
			}

			// If "squash" is specified in the tag, we squash the field down.
			squash = squash || strings.Index(tagValue[index+1:], "squash") != -1
			if squash {
				// When squashing, the embedded type can be a pointer to a struct.
				if v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct {
					v = v.Elem()
				}

				// The final type must be a struct
				if v.Kind() != reflect.Struct {
					return fmt.Errorf("cannot squash non-struct type '%s'", v.Type())
				}
			}
			if keyNameTagValue := tagValue[:index]; keyNameTagValue != "" {
				keyName = keyNameTagValue
			}
		} else if len(tagValue) > 0 {
			if tagValue == "-" {
				continue
			}
			keyName = tagValue
		}

		switch v.Kind() {
		// this is an embedded struct, so handle it differently
		case reflect.Struct:
			x := reflect.New(v.Type())
			x.Elem().Set(v)

			vType := valMap.Type()
			vKeyType := vType.Key()
			vElemType := vType.Elem()
			mType := reflect.MapOf(vKeyType, vElemType)
			vMap := reflect.MakeMap(mType)

			// Creating a pointer to a map so that other methods can completely
			// overwrite the map if need be (looking at you decodeMapFromMap). The
			// indirection allows the underlying map to be settable (CanSet() == true)
			// where as reflect.MakeMap returns an unsettable map.
			addrVal := reflect.New(vMap.Type())
			reflect.Indirect(addrVal).Set(vMap)

			err := decode(x.Interface(), reflect.Indirect(addrVal), tag, weakly)
			if err != nil {
				return err
			}

			// the underlying map may have been completely overwritten so pull
			// it indirectly out of the enclosing value.
			vMap = reflect.Indirect(addrVal)

			if squash {
				for _, k := range vMap.MapKeys() {
					valMap.SetMapIndex(k, vMap.MapIndex(k))
				}
			} else {
				valMap.SetMapIndex(reflect.ValueOf(keyName), vMap)
			}

		default:
			valMap.SetMapIndex(reflect.ValueOf(keyName), v)
		}
	}

	if val.CanAddr() {
		val.Set(valMap)
	}

	return nil
}

func decodePtr(data interface{}, val reflect.Value, tag string, weakly bool) error {
	// If the input data is nil, then we want to just set the output
	// pointer to be nil as well.
	isNil := data == nil
	if !isNil {
		switch v := reflect.Indirect(reflect.ValueOf(data)); v.Kind() {
		case reflect.Chan,
			reflect.Func,
			reflect.Interface,
			reflect.Map,
			reflect.Ptr,
			reflect.Slice:
			isNil = v.IsNil()
		}
	}
	if isNil {
		if !val.IsNil() && val.CanSet() {
			nilValue := reflect.New(val.Type()).Elem()
			val.Set(nilValue)
		}

		return nil
	}

	// Create an element of the concrete (non pointer) type and decode
	// into that. Then set the value of the pointer to this type.
	valType := val.Type()
	valElemType := valType.Elem()
	if val.CanSet() {
		realVal := val
		if realVal.IsNil() {
			realVal = reflect.New(valElemType)
		}

		if err := decode(data, reflect.Indirect(realVal), tag, weakly); err != nil {
			return err
		}

		val.Set(realVal)
	} else {
		if err := decode(data, reflect.Indirect(val), tag, weakly); err != nil {
			return err
		}
	}
	return nil
}

func decodeFunc(data interface{}, val reflect.Value, tag string, weakly bool) error {
	// Create an element of the concrete (non pointer) type and decode
	// into that. Then set the value of the pointer to this type.
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	if val.Type() != dataVal.Type() {
		return fmt.Errorf(
			"expected type '%s', got unconvertible type '%s', value: '%v'", val.Type(), dataVal.Type(), data)
	}
	val.Set(dataVal)
	return nil
}

func decodeSlice(data interface{}, val reflect.Value, tag string, weakly bool) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataValKind := dataVal.Kind()
	valType := val.Type()
	valElemType := valType.Elem()
	sliceType := reflect.SliceOf(valElemType)

	// If we have a non array/slice type then we first attempt to convert.
	if dataValKind != reflect.Array && dataValKind != reflect.Slice {
		if weakly {
			switch {
			// Slice and array we use the normal logic
			case dataValKind == reflect.Slice, dataValKind == reflect.Array:
				break

			// Empty maps turn into empty slices
			case dataValKind == reflect.Map:
				if dataVal.Len() == 0 {
					val.Set(reflect.MakeSlice(sliceType, 0, 0))
					return nil
				}
				// Create slice of maps of other sizes
				return decodeSlice([]interface{}{data}, val, tag, weakly)

			case dataValKind == reflect.String && valElemType.Kind() == reflect.Uint8:
				return decodeSlice([]byte(dataVal.String()), val, tag, weakly)

			// All other types we try to convert to the slice type
			// and "lift" it into it. i.e. a string becomes a string slice.
			default:
				// Just re-try this function with data as a slice.
				return decodeSlice([]interface{}{data}, val, tag, weakly)
			}
		}

		return fmt.Errorf(
			"source data must be an array or slice, got %s", dataValKind)
	}

	// If the input value is nil, then don't allocate since empty != nil
	if dataValKind != reflect.Array && dataVal.IsNil() {
		return nil
	}

	valSlice := val
	if valSlice.IsNil() {
		// Make a new slice to hold our result, same size as the original data.
		valSlice = reflect.MakeSlice(sliceType, dataVal.Len(), dataVal.Len())
	} else if valSlice.Len() > dataVal.Len() {
		valSlice = valSlice.Slice(0, dataVal.Len())
	}

	// Accumulate any errors
	errors := make([]string, 0)

	for i := 0; i < dataVal.Len(); i++ {
		currentData := dataVal.Index(i).Interface()
		for valSlice.Len() <= i {
			valSlice = reflect.Append(valSlice, reflect.Zero(valElemType))
		}
		currentField := valSlice.Index(i)

		if err := decode(currentData, currentField, tag, weakly); err != nil {
			errors = append(errors, err.Error())
		}
	}

	// Finally, set the value to the slice we built up
	val.Set(valSlice)

	// If there were errors, we return those
	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "\n"))
	}

	return nil
}

func decodeArray(data interface{}, val reflect.Value, tag string, weakly bool) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataValKind := dataVal.Kind()
	valType := val.Type()
	valElemType := valType.Elem()
	arrayType := reflect.ArrayOf(valType.Len(), valElemType)

	valArray := val

	if valArray.Interface() == reflect.Zero(valArray.Type()).Interface() {
		// Check input type
		if dataValKind != reflect.Array && dataValKind != reflect.Slice {
			if weakly {
				switch {
				// Empty maps turn into empty arrays
				case dataValKind == reflect.Map:
					if dataVal.Len() == 0 {
						val.Set(reflect.Zero(arrayType))
						return nil
					}

				// All other types we try to convert to the array type
				// and "lift" it into it. i.e. a string becomes a string array.
				default:
					// Just re-try this function with data as a slice.
					return decodeArray([]interface{}{data}, val, tag, weakly)
				}
			}

			return fmt.Errorf(
				"source data must be an array or slice, got %s", dataValKind)

		}
		if dataVal.Len() > arrayType.Len() {
			return fmt.Errorf(
				"expected source data to have length less or equal to %d, got %d", arrayType.Len(), dataVal.Len())

		}

		// Make a new array to hold our result, same size as the original data.
		valArray = reflect.New(arrayType).Elem()
	}

	// Accumulate any errors
	errors := make([]string, 0)

	for i := 0; i < dataVal.Len(); i++ {
		currentData := dataVal.Index(i).Interface()
		currentField := valArray.Index(i)

		if err := decode(currentData, currentField, tag, weakly); err != nil {
			errors = append(errors, err.Error())
		}
	}

	// Finally, set the value to the array we built up
	val.Set(valArray)

	// If there were errors, we return those
	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "\n"))
	}

	return nil
}

func decodeStruct(data interface{}, val reflect.Value, tag string, weakly bool) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))

	// If the type of the value to write to and the data match directly,
	// then we just set it directly instead of recursing into the structure.
	if dataVal.Type() == val.Type() {
		val.Set(dataVal)
		return nil
	}

	dataValKind := dataVal.Kind()
	switch dataValKind {
	case reflect.Map:
		return decodeStructFromMap(dataVal, val, tag, weakly)

	case reflect.Struct:
		// Not the most efficient way to do this but we can optimize later if
		// we want to. To convert from struct to struct we go to map first
		// as an intermediary.

		// Make a new map to hold our result
		mapType := reflect.TypeOf((map[string]interface{})(nil))
		mval := reflect.MakeMap(mapType)

		// Creating a pointer to a map so that other methods can completely
		// overwrite the map if need be (looking at you decodeMapFromMap). The
		// indirection allows the underlying map to be settable (CanSet() == true)
		// where as reflect.MakeMap returns an unsettable map.
		addrVal := reflect.New(mval.Type())

		reflect.Indirect(addrVal).Set(mval)
		if err := decodeMapFromStruct(dataVal, reflect.Indirect(addrVal), mval, tag, weakly); err != nil {
			return err
		}

		result := decodeStructFromMap(reflect.Indirect(addrVal), val, tag, weakly)
		return result

	default:
		return fmt.Errorf("expected a map, got '%s'", dataVal.Kind())
	}
}

func decodeStructFromMap(dataVal, val reflect.Value, tag string, weakly bool) error {
	dataValType := dataVal.Type()
	if kind := dataValType.Key().Kind(); kind != reflect.String && kind != reflect.Interface {
		return fmt.Errorf(
			"needs a map with string keys, has '%s' keys", dataValType.Key().Kind())
	}

	dataValKeys := make(map[reflect.Value]struct{})
	dataValKeysUnused := make(map[interface{}]struct{})
	for _, dataValKey := range dataVal.MapKeys() {
		dataValKeys[dataValKey] = struct{}{}
		dataValKeysUnused[dataValKey.Interface()] = struct{}{}
	}

	targetValKeysUnused := make(map[interface{}]struct{})
	errors := make([]string, 0)

	// This slice will keep track of all the structs we'll be decoding.
	// There can be more than one struct if there are embedded structs
	// that are squashed.
	structs := make([]reflect.Value, 1, 5)
	structs[0] = val

	// Compile the list of all the fields that we're going to be decoding
	// from all the structs.
	type field struct {
		field reflect.StructField
		val   reflect.Value
	}

	// remainField is set to a valid field set with the "remain" tag if
	// we are keeping track of remaining values.
	var remainField *field

	fields := []field{}
	for len(structs) > 0 {
		structVal := structs[0]
		structs = structs[1:]

		structType := structVal.Type()

		for i := 0; i < structType.NumField(); i++ {
			fieldType := structType.Field(i)
			fieldVal := structVal.Field(i)
			if fieldVal.Kind() == reflect.Ptr && fieldVal.Elem().Kind() == reflect.Struct {
				// Handle embedded struct pointers as embedded structs.
				fieldVal = fieldVal.Elem()
			}

			// If "squash" is specified in the tag, we squash the field down.
			squash := fieldVal.Kind() == reflect.Struct && fieldType.Anonymous
			remain := false

			// We always parse the tags cause we're looking for other tags too
			tagParts := strings.Split(fieldType.Tag.Get(tag), ",")
			for _, tag := range tagParts[1:] {
				if tag == "squash" {
					squash = true
					break
				}

				if tag == "remain" {
					remain = true
					break
				}
			}

			if squash {
				if fieldVal.Kind() != reflect.Struct {
					errors = append(errors,
						fmt.Sprintf("%s: unsupported type for squash: %s", fieldType.Name, fieldVal.Kind()))
				} else {
					structs = append(structs, fieldVal)
				}
				continue
			}

			// Build our field
			if remain {
				remainField = &field{fieldType, fieldVal}
			} else {
				// Normal struct field, store it away
				fields = append(fields, field{fieldType, fieldVal})
			}
		}
	}

	// for fieldType, field := range fields {
	for _, f := range fields {
		field, fieldValue := f.field, f.val
		fieldName := field.Name

		tagValue := field.Tag.Get(tag)
		tagValue = strings.SplitN(tagValue, ",", 2)[0]
		if tagValue != "" {
			fieldName = tagValue
		}

		rawMapKey := reflect.ValueOf(fieldName)
		rawMapVal := dataVal.MapIndex(rawMapKey)
		if !rawMapVal.IsValid() {
			// Do a slower search by iterating over each key and
			// doing case-insensitive search.
			for dataValKey := range dataValKeys {
				mK, ok := dataValKey.Interface().(string)
				if !ok {
					// Not a string key
					continue
				}

				if mK == fieldName {
					rawMapKey = dataValKey
					rawMapVal = dataVal.MapIndex(dataValKey)
					break
				}
			}

			if !rawMapVal.IsValid() {
				// There was no matching key in the map for the value in
				// the struct. Remember it for potential errors and metadata.
				targetValKeysUnused[fieldName] = struct{}{}
				continue
			}
		}

		if !fieldValue.IsValid() {
			// This should never happen
			panic("field is not valid")
		}

		// If we can't set the field, then it is unexported or something,
		// and we just continue onwards.
		if !fieldValue.CanSet() {
			continue
		}

		// Delete the key we're using from the unused map so we stop tracking
		delete(dataValKeysUnused, rawMapKey.Interface())

		if err := decode(rawMapVal.Interface(), fieldValue, tag, weakly); err != nil {
			errors = append(errors, err.Error())
		}
	}

	// If we have a "remain"-tagged field and we have unused keys then
	// we put the unused keys directly into the remain field.
	if remainField != nil && len(dataValKeysUnused) > 0 {
		// Build a map of only the unused values
		remain := map[interface{}]interface{}{}
		for key := range dataValKeysUnused {
			remain[key] = dataVal.MapIndex(reflect.ValueOf(key)).Interface()
		}

		// Decode it as-if we were just decoding this map onto our map.
		if err := decodeMap(remain, remainField.val, tag, weakly); err != nil {
			errors = append(errors, err.Error())
		}

		// Set the map to nil so we have none so that the next check will
		// not error (ErrorUnused)
		dataValKeysUnused = nil
	}

	if len(dataValKeysUnused) > 0 {
		keys := make([]string, 0, len(dataValKeysUnused))
		for rawKey := range dataValKeysUnused {
			keys = append(keys, rawKey.(string))
		}
		sort.Strings(keys)

		err := fmt.Errorf("has invalid keys: %s", strings.Join(keys, ", "))
		errors = append(errors, err.Error())
	}

	if len(targetValKeysUnused) > 0 {
		keys := make([]string, 0, len(targetValKeysUnused))
		for rawKey := range targetValKeysUnused {
			keys = append(keys, rawKey.(string))
		}
		sort.Strings(keys)

		err := fmt.Errorf("has unset fields: %s", strings.Join(keys, ", "))
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "\n"))
	}

	return nil
}

func isEmptyValue(v reflect.Value) bool {
	switch getKind(v) {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func getKind(val reflect.Value) reflect.Kind {
	kind := val.Kind()

	switch {
	case kind >= reflect.Int && kind <= reflect.Int64:
		return reflect.Int
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return reflect.Uint
	case kind >= reflect.Float32 && kind <= reflect.Float64:
		return reflect.Float32
	default:
		return kind
	}
}

func isStructTypeConvertibleToMap(typ reflect.Type, checkMapstructureTags bool, tagName string) bool {
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.PkgPath == "" && !checkMapstructureTags { // check for unexported fields
			return true
		}
		if checkMapstructureTags && f.Tag.Get(tagName) != "" { // check for mapstructure tags inside
			return true
		}
	}
	return false
}

func dereferencePtrToStructIfNeeded(v reflect.Value, tagName string) reflect.Value {
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return v
	}
	deref := v.Elem()
	derefT := deref.Type()
	if isStructTypeConvertibleToMap(derefT, true, tagName) {
		return deref
	}
	return v
}
