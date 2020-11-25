package httpunit

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
)

type HTTPFixture struct {
	t T

	DumpHandler      *DumpHandler
	RequestBuilder   *RequestBuilder
	ResponseRecorder *httptest.ResponseRecorder
}

func NewFixture(t T, inner http.Handler) *HTTPFixture {
	return &HTTPFixture{
		t: t,

		RequestBuilder:   NewRequestBuilder(),
		DumpHandler:      NewDumpHandler(t, inner),
		ResponseRecorder: httptest.NewRecorder(),
	}
}
func (this *HTTPFixture) Teardown() {
	this.DumpHandler.Teardown()
}
func (this *HTTPFixture) Serve() {
	this.DumpHandler.ServeHTTP(this.ResponseRecorder, this.RequestBuilder.Build())
}

func (this *HTTPFixture) DeserializeJSONResponseBody(v interface{}) {
	decoder := json.NewDecoder(this.ResponseRecorder.Result().Body)
	err := decoder.Decode(v)
	if err != nil {
		log.Panicln(err)
	}
}
func (this *HTTPFixture) AssertJSONResponse(expectedStatus int, expectedBody interface{}) {
	this.AssertResponseStatusCode(expectedStatus)

	var actualBody interface{}
	this.DeserializeJSONResponseBody(&actualBody)
	if !reflect.DeepEqual(expectedBody, actualBody) {
		this.reportFailure(expectedBody, actualBody)
	}
}
func (this *HTTPFixture) AssertResponseStatusCode(expected int) {
	actual := this.ResponseRecorder.Result().StatusCode
	if actual == expected {
		return
	}
	this.t.Helper()
	this.reportFailure(expected, actual)
}
func (this *HTTPFixture) reportFailure(expected, actual interface{}) {
	this.t.Helper()
	this.t.Errorf("\n"+
		"Expected: %#v\n"+
		"Actual:   %#v",
		expected,
		actual,
	)
}
