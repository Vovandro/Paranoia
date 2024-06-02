package interfaces

import (
	"time"
)

type ICache interface {
	Init(app IService) error
	Stop() error
	String() string

	/*
	   Has checks if the given key exists.
	   It returns true if the key exists, false otherwise.

	   Parameters:
	   - key (string): The key to check.

	   Returns:
	   - bool: true if the key exists, false otherwise.
	*/
	Has(key string) bool

	/*
	   Set sets a key with the provided data and expiration time.

	   Parameters:
	   - key: the key to set
	   - args: the data to be stored, can be of type []byte, string, or any other type that can be converted to string
	   - timeout: the duration after which the key will expire

	   Returns:
	   - error: an error if the key could not be set
	*/
	Set(key string, args any, timeout time.Duration) error

	/*
	   SetIn updates the value of a specific key within a map stored. If the key does not exist, a new map is created.
	   It takes the key, the subkey, the value to set, and the timeout duration as parameters.

	   Parameters:
	   - key: the key of the map
	   - key2: the subkey within the map to update
	   - args: the value to set
	   - timeout: the duration for which the data should be stored

	   Returns:
	   - error: an error if the operation fails, nil otherwise
	*/
	SetIn(key string, key2 string, args any, timeout time.Duration) error

	/*
	   SetMap sets a key-value pair by marshaling the input args into JSON format before calling the Set method.

	   Parameters:
	   - key: the key to set
	   - args: the value to set, will be marshaled into JSON format
	   - timeout: the duration after which the key-value pair will expire

	   Returns:
	   - error: an error if the operation encounters any issues, nil otherwise
	*/
	SetMap(key string, args any, timeout time.Duration) error

	/*
	   Get retrieves the value associated with the given key.

	   Parameters:
	   - key (string): The key to look up.

	   Returns:
	   - any: The value associated with the key.
	   - error: An error, if any occurred during the retrieval process. Returns ErrKeyNotFound if the key is not found in the cache.
	*/
	Get(key string) (any, error)

	/*
	   GetIn retrieves the value associated with the specified key and subkey.

	   Parameters:
	   - key: The key to look up.
	   - key2: The subkey to look up within the key's associated value.

	   Returns:
	   - any: The value associated with the key and subkey.
	   - error: An error if the key, subkey, or associated value is not found.

	   Notes:
	   - If the key or subkey is not found, an error with ErrKeyNotFound will be returned.
	*/
	GetIn(key string, key2 string) (any, error)

	/*
	   GetMap retrieves a value using the specified key and returns it as a map[string]interface{}.

	   Parameters:
	   - key (string): The key used to retrieve the value.

	   Returns:
	   - map[string]interface{}: The value retrieved stored as a map.
	   - error: An error if any occurred during the retrieval or unmarshalling process.
	*/
	GetMap(key string) (any, error)

	/*
	   Increment increments the value of the given key by the specified amount.
	   If the key does not exist, it sets the key with the provided value and expiration time.

	   Parameters:
	   - key: The key to increment.
	   - val: The amount by which to increment the value of the key.
	   - timeout: The duration after which the key will expire.

	   Returns:
	   - int64: New value
	   - error: An error if any occurred during the increment operation, or setting the key, or touching the key with the new expiration time.
	*/
	Increment(key string, val int64, timeout time.Duration) (int64, error)

	/*
	   IncrementIn increments the value associated with key2 in the map stored at key by the specified value val.
	   If the key does not exist, a new map is created. If an error occurs during the retrieval of the map, it is returned.
	   If the value at key2 is not an integer, an error is returned.

	   Parameters:
	   - key: The key of the map to be incremented.
	   - key2: The key within the map whose value will be incremented.
	   - val: The value by which to increment the existing value at key2.
	   - timeout: The duration after which the operation times out.

	   Returns:
	   - int64: New value
	   - error: An error if any occurred during the operation.
	*/
	IncrementIn(key string, key2 string, val int64, timeout time.Duration) (int64, error)

	/*
	   Decrement decrements the value associated with the given key by the specified amount.
	   If the key is not found in the cache, it returns ErrKeyNotFound.
	   If any other error occurs during the decrement operation, that error is returned.
	   After decrementing the value, it updates the expiration time of the key with the provided timeout.

	   Parameters:
	   - key: The key for which the value needs to be decremented.
	   - val: The amount by which the value should be decremented.
	   - timeout: The duration after which the key should expire if not accessed.

	   Returns:
	   - int64: New value
	   - error: An error if the decrement operation or updating the expiration time fails.
	*/
	Decrement(key string, val int64, timeout time.Duration) (int64, error)

	/*
	   DecrementIn decrements the value associated with key2 in the map stored at key by the specified val.
	   If the key does not exist, a new map will be created.
	   If key2 does not exist in the map, key2 will be added with the negative value of val.
	   The updated map will then be stored with the specified timeout.

	   Parameters:
	   - key: The key under which the map is stored.
	   - key2: The key within the map whose value needs to be decremented.
	   - val: The value by which key2 should be decremented.
	   - timeout: The duration after which the updated map will expire.

	   Returns:
	   - int64: New value
	   - error: An error if the operation encounters any issues, nil otherwise.
	*/
	DecrementIn(key string, key2 string, val int64, timeout time.Duration) (int64, error)

	/*
	   Delete deletes the value for a key.

	   If the key is not found in the cache, it returns ErrKeyNotFound.
	   If an error occurs during the deletion operation, that error is returned.

	   Parameters:
	   - key (string): The key for which the value needs to be deleted.

	   Returns:
	   - error: An error if the deletion operation encounters any issues.
	*/
	Delete(key string) error
}
