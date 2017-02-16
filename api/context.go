package api

import (
	"fmt"
)

func (ctx *RequestContext) Error(status int, err interface{}, msg ...interface{}) {
	ret := map[string]string{
		"error": fmt.Sprint(err),
	}
	if len(msg) > 0 {
		ret["message"] = fmt.Sprint(msg...)
	}
	ctx.rnd.JSON(status, ret)
}

func (ctx *RequestContext) Message(status int, msg ...interface{}) {
	switch len(msg) {
	case 0:
		ctx.res.WriteHeader(status)
	case 1:
		ctx.rnd.JSON(status, map[string]interface{}{"message": msg[0]})
	default:
		ret := map[string]interface{}{"message": msg[0]}
		for i := 1; i < len(msg); i++ {
			ret[fmt.Sprint("message", i)] = msg[i]
		}
		ctx.rnd.JSON(status, ret)
	}
}

// func (ctx *RequestContext) EventHandle() {
// 	ctx.res.Header().Set("Content-Type", "application/json")
// 	ctx.res.WriteHeader(200)

// 	for {
// 		select {
// 		case evt := <-ctx.eventChan:
// 			_, err := ctx.res.Write([]byte(time.Now().String() + ":" + evt + "\n\n"))
// 			if err != nil {
// 				return
// 			}
// 		case <-time.Tick(5 * time.Second):
// 			_, err := ctx.res.Write([]byte("this is a ping."))
// 			if err != nil {
// 				return
// 			}
// 		}

// 		if flu, ok := ctx.res.(http.Flusher); ok {
// 			flu.Flush()
// 		}
// 	}
// 	for ; ; time.Sleep(1 * time.Second) {
// 		fmt.Println(time.Now().String())
// 		_, err := ctx.res.Write([]byte(time.Now().String() + "\n\n"))
// 		if err != nil {
// 			return
// 		}
// 		if flu, ok := ctx.res.(http.Flusher); ok {
// 			flu.Flush()
// 		}
// 	}
// }
