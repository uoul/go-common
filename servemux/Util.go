package servemux

import (
	"net/http"
)

// -----------------------------------------------------------------------------------------------------------
// Type
// -----------------------------------------------------------------------------------------------------------

type HandlerFunc[T any] func(ctx *HttpCtx[T])

// -----------------------------------------------------------------------------------------------------------
// Public Functions
// -----------------------------------------------------------------------------------------------------------

func Handle[T any](config *HandlerConfig, handlers ...HandlerFunc[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create HttpCtx
		httpCtx := NewHttpCtx[T](r, w, config.serializer, config.maxMemSize)
		// Call handlers
		for _, handler := range handlers {
			handler(httpCtx)
			if httpCtx.IsAborted() {
				break
			}
		}
		// Parse Response body
		respBody, _ := config.serializer.Marshal(
			httpCtx.GetResponseBody(),
		)
		// Write ResultHeader
		w.WriteHeader(httpCtx.GetStatusCode())
		// Write Result
		w.Write(respBody)
	}
}
