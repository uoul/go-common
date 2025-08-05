package servemux

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
	"github.com/uoul/go-common/collections"
	"github.com/uoul/go-common/serialization"
)

// -----------------------------------------------------------------------------------------------------------
// Types
// -----------------------------------------------------------------------------------------------------------

type HttpCtx[T any] struct {
	req        *http.Request
	respWriter http.ResponseWriter
	serializer serialization.ISerializer
	maxMemSize int64
	decoder    *schema.Decoder

	statusCode   int
	responseBody any
	errors       []error
}

type ServerSentEvent struct {
	Event string // Event type (e.g. "message")
	Data  any    // Data (will be serialized using given serializer)
}

func (s *ServerSentEvent) marshal(serializer serialization.ISerializer) ([]byte, error) {
	data, err := serializer.Marshal(s.Data)
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("event: %s\ndata: %s\n\n", s.Event, data)), nil
}

// -----------------------------------------------------------------------------------------------------------
// Public
// -----------------------------------------------------------------------------------------------------------

func (h *HttpCtx[T]) GetBody() (T, error) {
	contentType := h.GetHeader("Content-Type")
	if collections.ContainsSlice(contentType, func(t string) bool { return t == "multipart/form-data" }) {
		return h.parseFormData(h.req.Form)
	} else if collections.ContainsSlice(contentType, func(t string) bool { return t == "application/x-www-form-urlencoded" }) {
		return h.parseFormData(h.req.MultipartForm.Value)
	} else {
		return h.parseBody()
	}
}

func (h *HttpCtx[T]) Context() context.Context {
	return h.req.Context()
}

func (h *HttpCtx[T]) GetFiles() map[string][]*multipart.FileHeader {
	return h.req.MultipartForm.File
}

func (h *HttpCtx[T]) GetHeader(name string) []string {
	return h.req.Header[name]
}

func (h *HttpCtx[T]) GetQueryParam(name string) string {
	return h.req.URL.Query().Get(name)
}

func (h *HttpCtx[T]) GetUrlParam(name string) string {
	return h.req.PathValue(name)
}

func (h *HttpCtx[T]) GetStatusCode() int {
	return h.statusCode
}

func (h *HttpCtx[T]) SetStatusCode(code int) {
	h.statusCode = code
}

func (h *HttpCtx[T]) Errors() []error {
	return h.errors
}

func (h *HttpCtx[T]) Error(err error) {
	h.errors = append(h.errors, err)
}

func (h *HttpCtx[T]) GetResponseBody() any {
	return h.responseBody
}

func (h *HttpCtx[T]) SetResponseBody(value any) {
	h.responseBody = value
}

func (h *HttpCtx[T]) Stream(step func() (ServerSentEvent, bool)) (bool, error) {
	// Set Headers for streaming
	h.respWriter.Header().Set("Content-Type", "text/event-stream")
	h.respWriter.Header().Set("Cache-Control", "no-cache")
	h.respWriter.Header().Set("Connection", "keep-alive")
	h.respWriter.Header().Set("Transfer-Encoding", "chunked")
	// Create http flusher
	flusher, ok := h.respWriter.(http.Flusher)
	if !ok {
		return false, fmt.Errorf("failed to create http flusher")
	}
	// Run Stream
	for {
		select {
		case <-h.Context().Done():
			return true, nil
		default:
			msg, proceed := step()
			if !proceed {
				return false, nil
			}
			sse, err := msg.marshal(h.serializer)
			if err != nil {
				return false, err
			}
			if _, err := h.respWriter.Write([]byte(sse)); err != nil {
				return false, err
			}
			flusher.Flush()
		}
	}
}

// -----------------------------------------------------------------------------------------------------------
// Private
// -----------------------------------------------------------------------------------------------------------

func (h *HttpCtx[T]) parseFormData(data url.Values) (T, error) {
	body := *new(T)
	err := h.decoder.Decode(&body, data)
	return body, err
}

func (h *HttpCtx[T]) parseBody() (T, error) {
	reader, err := h.req.GetBody()
	if err != nil {
		return *new(T), err
	}
	raw, err := io.ReadAll(reader)
	if err != nil {
		return *new(T), err
	}
	body := *new(T)
	if err := h.serializer.Unmarshal(raw, &body); err != nil {
		return *new(T), err
	}
	return body, nil
}

// -----------------------------------------------------------------------------------------------------------
// Constructor
// -----------------------------------------------------------------------------------------------------------

func NewHttpCtx[T any](req *http.Request, responseWriter http.ResponseWriter, serializer serialization.ISerializer, maxMemSize int64) *HttpCtx[T] {
	// Create new HttpCtx
	h := &HttpCtx[T]{
		req:        req,
		respWriter: responseWriter,
		serializer: serializer,
		maxMemSize: maxMemSize,

		statusCode:   http.StatusOK,
		decoder:      schema.NewDecoder(),
		responseBody: nil,
		errors:       []error{},
	}
	// Parse Form data if form data encoding
	contentType := h.GetHeader("Content-Type")
	if collections.ContainsSlice(contentType, func(t string) bool { return t == "multipart/form-data" }) {
		h.req.ParseMultipartForm(h.maxMemSize)
	} else if collections.ContainsSlice(contentType, func(t string) bool { return t == "application/x-www-form-urlencoded" }) {
		h.req.ParseForm()
	}
	// Return context
	return h
}
