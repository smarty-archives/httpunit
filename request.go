package httpunit

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

type RequestBuilder struct {
	Context context.Context
	Method  string
	URL     string
	Headers http.Header
	Body    string
	JSON    interface{}
}

func NewRequestBuilder() *RequestBuilder {
	return &RequestBuilder{
		Method:  http.MethodGet,
		URL:     "/",
		Headers: make(http.Header),
	}
}

func (this *RequestBuilder) Build() *http.Request {
	request := httptest.NewRequest(this.Method, this.URL, this.body())
	if this.Context != nil {
		request = request.WithContext(this.Context)
	}
	for key, values := range this.Headers {
		request.Header[key] = values
	}
	return request
}
func (this *RequestBuilder) body() (body io.Reader) {
	if len(this.Body) > 0 {
		return strings.NewReader(this.Body)
	}
	if this.JSON == nil {
		return nil
	}
	raw, err := json.Marshal(this.JSON)
	if err != nil {
		panic(err)
	}
	this.Headers.Set("Content-Type", "application/json; charset=utf-8")
	return bytes.NewReader(raw)
}
