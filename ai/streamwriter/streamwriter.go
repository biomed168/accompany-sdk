package streamwriter

import (
	"accompany-sdk/pkg/misc"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"accompany-sdk/pkg/ternary"
	"github.com/gorilla/websocket"
	"github.com/openimsdk/tools/log"
)

type StreamWriter struct {
	ws         *websocket.Conn
	r          *http.Request
	w          http.ResponseWriter
	enableCors bool

	once      sync.Once
	sseInited bool
	debug     bool

	onClosedSync sync.Once
	onClosed     func()
}

var corsHeaders = http.Header{
	"Access-Control-Allow-Origin":  []string{"*"},
	"Access-Control-Allow-Headers": []string{"*"},
	"Access-Control-Allow-Methods": []string{"GET,POST,OPTIONS,HEAD,PUT,PATCH,DELETE"},
}

type WSError struct {
	Code  int    `json:"code,omitempty"`
	Error string `json:"error,omitempty"`
}

func (e WSError) JSON() []byte {
	data, _ := json.Marshal(e)
	return data
}

func (sw *StreamWriter) Close() {
	if sw.debug {
		log.ZDebug(context.Background(), "close stream writer")
	}

	sw.handleClosed()
}

type InitRequest[T any] interface {
	Init() T
}

func (sw *StreamWriter) SetOnClosed(cb func()) {
	sw.onClosed = cb
}

func (sw *StreamWriter) handleClosed() {
	sw.onClosedSync.Do(func() {
		if sw.ws != nil {
			_ = sw.ws.Close()
		} else {
			if sw.sseInited {
				// 写入结束标志
				_, _ = sw.w.Write([]byte("data: [DONE]\n\n"))
				if f, ok := sw.w.(http.Flusher); ok {
					f.Flush()
				}
			}
		}

		if sw.onClosed != nil {
			sw.onClosed()
		}
	})
}

func New[T InitRequest[T]](enableWs bool, enableCors bool, r *http.Request, w http.ResponseWriter) (*StreamWriter, *T, error) {
	sw := &StreamWriter{
		r:          r,
		w:          w,
		enableCors: enableCors,
	}

	var req T
	if enableWs {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		if wsConn, err := upgrader.Upgrade(w, r, ternary.If(enableCors, corsHeaders, http.Header{})); err != nil {
			sw.writeJSON(NewErrorResponse(fmt.Errorf("upgrade websocket failed: %v", err)), http.StatusInternalServerError)
			return nil, nil, err
		} else {
			sw.ws = wsConn

			if sw.debug {
				log.ZDebug(context.Background(), "websocket connected: %s", wsConn.RemoteAddr())
			}

			// 读取第一条消息，用于获取用户输入
			_, msg, err := wsConn.ReadMessage()
			if err != nil {
				misc.NoError(context.Background(), sw.WriteStream(NewErrorResponse(fmt.Errorf("read websocket message failed: %v", err))))
				misc.NoError(context.Background(), wsConn.Close())
				return nil, nil, err
			}

			if err := json.Unmarshal(msg, &req); err != nil {
				misc.NoError(context.Background(), sw.WriteStream(NewErrorWithCodeResposne(fmt.Errorf("invalid request: %v", err), http.StatusBadRequest)))
				misc.NoError(context.Background(), wsConn.Close())
				return nil, nil, err
			}

			go func() {
				defer func() {
					sw.handleClosed()
				}()
				for {
					typ, msg, err := wsConn.ReadMessage()
					if err != nil {
						return
					}

					log.ZWarn(context.Background(), "receive message from websocket: (%d) %s", err, typ, string(msg))
				}
			}()
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sw.writeJSON(NewErrorWithCodeResposne(fmt.Errorf("invalid request: %v", err), http.StatusBadRequest), http.StatusBadRequest)
			return nil, nil, err
		}
	}

	req = req.Init()
	return sw, &req, nil
}

func (sw *StreamWriter) initSSE() {
	if sw.ws != nil {
		return
	}

	sw.once.Do(func() {
		sw.sseInited = true

		if sw.debug {
			log.ZDebug(context.Background(), "init sse")
		}

		sw.wrapRawResponse(sw.w, func() {
			sw.w.Header().Set("Content-Type", "text/event-stream")
			sw.w.Header().Set("Cache-Control", "no-cache")
			sw.w.Header().Set("Connection", "keep-alive")
		})
	})
}

func (sw *StreamWriter) WriteErrorStream(err error, statusCode int) error {
	return sw.WriteStream(NewErrorWithCodeResposne(err, statusCode))
}

func (sw *StreamWriter) WriteStream(payload any) error {
	var data []byte

	if str, ok := payload.(string); ok {
		data = []byte(str)
	} else {
		data, _ = json.Marshal(payload)
	}

	if sw.debug {
		log.ZDebug(context.Background(), "write stream: %s", string(data))
	}

	if sw.ws != nil {
		return sw.ws.WriteMessage(websocket.TextMessage, data)
	}

	sw.initSSE()

	if _, err := sw.w.Write([]byte("data: " + string(data) + "\n\n")); err != nil {
		return err
	}

	if f, ok := sw.w.(http.Flusher); ok {
		f.Flush()
	}

	return nil
}

func (sw *StreamWriter) wrapRawResponse(w http.ResponseWriter, cb func()) {
	// 允许跨域
	if sw.enableCors {
		for k, v := range corsHeaders {
			for _, v1 := range v {
				w.Header().Set(k, v1)
			}
		}
	}

	cb()
}

func (sw *StreamWriter) writeJSON(payload any, statusCode int) {
	sw.wrapRawResponse(sw.w, func() {
		data, err := json.Marshal(payload)
		if err != nil {
			sw.w.WriteHeader(http.StatusInternalServerError)
			sw.w.Write(ErrorResponse{Error: err.Error()}.ToJSON())
			return
		}

		sw.w.Header().Set("Content-Type", "application/json; charset=utf-8")
		sw.w.WriteHeader(statusCode)
		d, err := sw.w.Write(data)
		misc.NoError2(context.Background(), d, err)
	})
}

type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{Error: err.Error(), Code: http.StatusInternalServerError}
}

func NewErrorWithCodeResposne(err error, code int) ErrorResponse {
	return ErrorResponse{Error: err.Error(), Code: code}
}

func (resp ErrorResponse) ToJSON() []byte {
	data, _ := json.Marshal(resp)
	return data
}
