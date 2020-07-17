package httpunit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strings"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func NewHTTPFixture(inner *gunit.Fixture) *HTTPFixture {
	return &HTTPFixture{
		Fixture:        inner,
		RequestMethod:  http.MethodGet,
		RequestURL:     url.URL{Path: "/"},
		RequestHeaders: make(http.Header),
		RequestContext: context.Background(),
		Dump:           new(bytes.Buffer),
	}
}

type HTTPFixture struct {
	*gunit.Fixture

	RequestMethod  string
	RequestURL     url.URL
	RequestBody    string
	RequestHeaders http.Header
	RequestContext context.Context

	Response     *http.Response
	ResponseBody string

	Dump *bytes.Buffer
}

func (this *HTTPFixture) Teardown() {
	if this.Failed() || testing.Verbose() {
		this.Println(this.Dump.String())
	}
}

func (this *HTTPFixture) SetQueryStringParameter(key, value string) {
	query := this.RequestURL.Query()
	query.Set(key, value)
	this.RequestURL.RawQuery = query.Encode()
}
func (this *HTTPFixture) SetRequestContextValue(key, value interface{}) {
	this.RequestContext = context.WithValue(this.RequestContext, key, value)
}
func (this *HTTPFixture) SerializeJSONRequestBody(body interface{}) {
	raw, _ := json.Marshal(body)
	this.RequestBody = string(raw)
}

func (this *HTTPFixture) Serve(handler http.Handler) {
	request := this.buildRequest()
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	this.collectResponse(recorder)
}
func (this *HTTPFixture) buildRequest() *http.Request {
	request := httptest.NewRequest(
		this.RequestMethod,
		this.RequestURL.String(),
		strings.NewReader(this.RequestBody),
	)
	request.Header = this.RequestHeaders
	request = request.WithContext(this.RequestContext)
	this.dumpRequest(request)
	return request
}
func (this *HTTPFixture) collectResponse(recorder *httptest.ResponseRecorder) {
	response := recorder.Result()
	this.dumpResponse(response)
	body, _ := ioutil.ReadAll(response.Body)
	this.ResponseBody = string(body)
	this.Response = response
}
func (this *HTTPFixture) dumpRequest(request *http.Request) {
	requestDump, _ := httputil.DumpRequest(request, true)
	fmt.Fprintf(this.Dump, "REQUEST DUMP:\n%s\n\n", formatDump(">", string(requestDump)))
}
func (this *HTTPFixture) dumpResponse(response *http.Response) {
	responseDump, _ := httputil.DumpResponse(response, true)
	fmt.Fprintf(this.Dump, "RESPONSE DUMP:\n%s\n\n", formatDump("<", string(responseDump)))
}
func formatDump(prefix, dump string) string {
	prefix = "\n" + prefix + " "
	lines := strings.Split(strings.TrimSpace(dump), "\n")
	return prefix + strings.Join(lines, prefix)
}

func (this *HTTPFixture) DeserializeJSONResponseBody() (actual interface{}) {
	err := json.Unmarshal([]byte(this.ResponseBody), &actual)
	if err != nil {
		log.Panicln("JSON UNMARSHAL:", err)
	}
	return actual
}

func (this *HTTPFixture) AssertJSONResponse(expectedStatus int, expectedBody interface{}) {
	this.So(this.Response.StatusCode, should.Equal, expectedStatus)
	this.So(this.Response.Header.Get("Content-Type"), should.Equal, "application/json; charset=utf-8")
	this.So(this.DeserializeJSONResponseBody(), should.Resemble, expectedBody)
}
func (this *HTTPFixture) AssertInternalServerError() {
	this.So(this.Response.StatusCode, should.Equal, http.StatusInternalServerError)
	this.So(this.Response.Header.Get("Content-Type"), should.Equal, "text/plain; charset=utf-8")
	this.So(strings.TrimSpace(this.ResponseBody), should.Equal, "Internal Server Error")
}
