package httpunit_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/httpunit"
)

func TestHTTPFixtureFixture(t *testing.T) {
	gunit.Run(new(HTTPFixtureFixture), t)
}

type HTTPFixtureFixture struct {
	*gunit.Fixture
	*httpunit.HTTPFixture
}

func (this *HTTPFixtureFixture) Setup() {
	this.HTTPFixture = httpunit.NewHTTPFixture(this.Fixture)
}

func (this *HTTPFixtureFixture) TestRequestManipulation() {
	this.SetRequestContextValue("HELLO", "WORLD")
	this.So(this.RequestContext.Value("HELLO"), should.Equal, "WORLD")

	this.SetQueryStringParameter("hello", "world")
	this.So(this.RequestURL.Query().Get("hello"), should.Equal, "world")
}

func (this *HTTPFixtureFixture) TestAssertJSONResponse() {
	body := map[string]interface{}{"hello": "world"}
	this.SerializeJSONRequestBody(body)

	this.Serve(http.HandlerFunc(EchoHandler))

	this.So(this.ResponseBody, should.Equal, `{"hello":"world"}`)
	this.AssertJSONResponse(http.StatusOK, body)
}

func (this *HTTPFixtureFixture) TestAssertInternalServerError() {
	this.Serve(http.HandlerFunc(ServerErrorHandler))
	this.AssertInternalServerError()
}

func EchoHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = io.Copy(response, request.Body)
}

func ServerErrorHandler(response http.ResponseWriter, _ *http.Request) {
	http.Error(response, "Internal Server Error", http.StatusInternalServerError)
}
