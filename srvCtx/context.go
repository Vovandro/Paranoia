package srvCtx

//func FromKafka(msg *kafka.Message) *Ctx {
//	ctx := ContextPool.Get().(*Ctx)
//	ctx.Request.Body = msg.Value
//	ctx.Request.Headers = make(http.Header, len(msg.Headers))
//
//	for _, v := range msg.Headers {
//		if val, ok := ctx.Request.Headers[v.Key]; ok {
//			ctx.Request.Headers[v.Key] = append(val, string(v.Value))
//		} else {
//			ctx.Request.Headers[v.Key] = []string{string(v.Value)}
//		}
//	}
//
//	ctx.Request.Ip = ""
//	ctx.Request.URI = *msg.TopicPartition.Topic
//	ctx.Request.Method = "KAFKA"
//	ctx.Request.Host = ""
//	ctx.Request.PostForm = nil
//
//	ctx.Values = map[string]interface{}{}
//
//	ctx.Response.Body = ctx.Response.Body[:0]
//	ctx.Response.StatusCode = 200
//
//	return ctx
//}
