package httpunit

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"testing"
)

type DumpHandler struct {
	t T

	inner http.Handler
	dump  *bytes.Buffer
}

func NewDumpHandler(t T, inner http.Handler) *DumpHandler {
	return &DumpHandler{
		t: t,

		inner: inner,
		dump:  new(bytes.Buffer),
	}
}

func (this *DumpHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	this.dumpRequest(request)
	this.inner.ServeHTTP(response, request)
	this.dumpResponse(response.(*httptest.ResponseRecorder).Result())
}
func (this *DumpHandler) dumpRequest(request *http.Request) {
	requestDump, _ := httputil.DumpRequest(request, true)
	fmt.Fprintf(this.dump, "REQUEST DUMP:\n%s\n\n", formatDump(">", string(requestDump)))
}
func (this *DumpHandler) dumpResponse(response *http.Response) {
	responseDump, _ := httputil.DumpResponse(response, true)
	fmt.Fprintf(this.dump, "RESPONSE DUMP:\n%s\n\n", formatDump("<", string(responseDump)))
}
func formatDump(prefix, dump string) string {
	prefix = "\n" + prefix + " "
	lines := strings.Split(strings.TrimSpace(dump), "\n")
	return prefix + strings.Join(lines, prefix)
}

func (this *DumpHandler) Teardown() {
	if testing.Verbose() || this.t.Failed() {
		this.t.Log("\n"+this.dump.String())
	}
}

