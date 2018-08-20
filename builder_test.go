package httpxrouter

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestRoutesBuilderFixture(t *testing.T) {
	gunit.Run(new(RoutesBuilderFixture), t)
}

type RoutesBuilderFixture struct {
	*gunit.Fixture
	response *httptest.ResponseRecorder
	handler  http.Handler
}

func (this *RoutesBuilderFixture) Serve(method, path string) string {
	request := httptest.NewRequest(method, path, nil)
	this.response = httptest.NewRecorder()
	this.handler.ServeHTTP(this.response, request)
	return strings.TrimSpace(this.response.Body.String())
}

func (this *RoutesBuilderFixture) TestNotFound() {
	this.handler = New(GET("/hi", NewFakeHandler("I'm lonely")))
	this.Serve("GET", "/not-there")
	this.So(this.response.Code, should.Equal, http.StatusNotFound)
}

func (this *RoutesBuilderFixture) TestMethodNotAllowed() {
	this.handler = New(GET("/hi", NewFakeHandler("noop")))
	this.Serve("PATCH", "/hi")
	this.So(this.response.Code, should.Equal, http.StatusMethodNotAllowed)
}

func (this *RoutesBuilderFixture) TestAllMethodsOnASingleRoute() {
	this.handler = New(
		GET("/hi", NewFakeHandler("GET")),
		PUT("/hi", NewFakeHandler("PUT")),
		HEAD("/hi", NewFakeHandler("HEAD")),
		POST("/hi", NewFakeHandler("POST")),
		DELETE("/hi", NewFakeHandler("DELETE")),
		OPTIONS("/hi", NewFakeHandler("OPTIONS")),
	)

	this.So(this.Serve("PUT", "/hi"), should.Equal, "PUT")
	this.So(this.Serve("GET", "/hi"), should.Equal, "GET")
	this.So(this.Serve("POST", "/hi"), should.Equal, "POST")
	this.So(this.Serve("HEAD", "/hi"), should.Equal, "HEAD")
	this.So(this.Serve("DELETE", "/hi"), should.Equal, "DELETE")
	this.So(this.Serve("OPTIONS", "/hi"), should.Equal, "OPTIONS")
}

func (this *RoutesBuilderFixture) TestCompound() {
	this.handler = New(
		Compound(
			GET("/hi|/bye", NewFakeHandler("GET")),
			PUT("/hi|/bye", NewFakeHandler("PUT")),
		),
	)

	this.So(this.Serve("PUT", "/hi"), should.Equal, "PUT")
	this.So(this.Serve("GET", "/hi"), should.Equal, "GET")

	this.So(this.Serve("PUT", "/bye"), should.Equal, "PUT")
	this.So(this.Serve("GET", "/bye"), should.Equal, "GET")
}

func (this *RoutesBuilderFixture) TestAllMethodsOnAMultipleRoutes() {
	this.handler = New(
		GET("/hi|/bye", NewFakeHandler("GET")),
		PUT("/hi|/bye", NewFakeHandler("PUT")),
		HEAD("/hi|/bye", NewFakeHandler("HEAD")),
		POST("/hi|/bye", NewFakeHandler("POST")),
		DELETE("/hi|/bye", NewFakeHandler("DELETE")),
		OPTIONS("/hi|/bye", NewFakeHandler("OPTIONS")),
	)

	this.So(this.Serve("PUT", "/hi"), should.Equal, "PUT")
	this.So(this.Serve("GET", "/hi"), should.Equal, "GET")
	this.So(this.Serve("POST", "/hi"), should.Equal, "POST")
	this.So(this.Serve("HEAD", "/hi"), should.Equal, "HEAD")
	this.So(this.Serve("DELETE", "/hi"), should.Equal, "DELETE")
	this.So(this.Serve("OPTIONS", "/hi"), should.Equal, "OPTIONS")

	this.So(this.Serve("PUT", "/bye"), should.Equal, "PUT")
	this.So(this.Serve("GET", "/bye"), should.Equal, "GET")
	this.So(this.Serve("POST", "/bye"), should.Equal, "POST")
	this.So(this.Serve("HEAD", "/bye"), should.Equal, "HEAD")
	this.So(this.Serve("DELETE", "/bye"), should.Equal, "DELETE")
	this.So(this.Serve("OPTIONS", "/bye"), should.Equal, "OPTIONS")
}

func (this *RoutesBuilderFixture) TestAllMethodsOnRoute() {
	this.handler = New(Register("GET|PUT", "/hi|/bye", NewFakeHandler("GET and PUT")))

	this.So(this.Serve("GET", "/hi"), should.Equal, "GET and PUT")
	this.So(this.Serve("GET", "/bye"), should.Equal, "GET and PUT")
	this.So(this.Serve("PUT", "/hi"), should.Equal, "GET and PUT")
	this.So(this.Serve("PUT", "/bye"), should.Equal, "GET and PUT")
}

func (this *RoutesBuilderFixture) TestRegisteredHandlersAreInvokedBeforeRouteHandlers() {
	this.handler = New(Prepend(NewFakeHandler("A "), NewFakeHandler("B ")),
		GET("/hi", NewFakeHandler("GET")),
		PUT("/hi", NewFakeHandler("PUT")),
		HEAD("/hi", NewFakeHandler("HEAD")),
		POST("/hi", NewFakeHandler("POST")),
		DELETE("/hi", NewFakeHandler("DELETE")),
		OPTIONS("/hi", NewFakeHandler("OPTIONS")),
	)

	this.So(this.Serve("PUT", "/hi"), should.Equal, "A B PUT")
	this.So(this.Serve("GET", "/hi"), should.Equal, "A B GET")
	this.So(this.Serve("POST", "/hi"), should.Equal, "A B POST")
	this.So(this.Serve("HEAD", "/hi"), should.Equal, "A B HEAD")
	this.So(this.Serve("DELETE", "/hi"), should.Equal, "A B DELETE")
	this.So(this.Serve("OPTIONS", "/hi"), should.Equal, "A B OPTIONS")
}

func (this *RoutesBuilderFixture) TestRoutesMayReceiveMultipleNestingHandlers() {
	this.handler = New(
		GET("/hi", NewFakeHandler("GET"), NewFakeHandler("GET")),
		PUT("/hi", NewFakeHandler("PUT"), NewFakeHandler("PUT")),
		HEAD("/hi", NewFakeHandler("HEAD"), NewFakeHandler("HEAD")),
		POST("/hi", NewFakeHandler("POST"), NewFakeHandler("POST")),
		DELETE("/hi", NewFakeHandler("DELETE"), NewFakeHandler("DELETE")),
		OPTIONS("/hi", NewFakeHandler("OPTIONS"), NewFakeHandler("OPTIONS")),
	)

	this.So(this.Serve("PUT", "/hi"), should.Equal, "PUT"+"PUT")
	this.So(this.Serve("GET", "/hi"), should.Equal, "GET"+"GET")
	this.So(this.Serve("POST", "/hi"), should.Equal, "POST"+"POST")
	this.So(this.Serve("HEAD", "/hi"), should.Equal, "HEAD"+"HEAD")
	this.So(this.Serve("DELETE", "/hi"), should.Equal, "DELETE"+"DELETE")
	this.So(this.Serve("OPTIONS", "/hi"), should.Equal, "OPTIONS"+"OPTIONS")
}

func (this *RoutesBuilderFixture) Test404NotFoundCallsRegisteredHandler() {
	this.handler = New(NotFound(NewFakeNotFoundHelperHandler("Not Here")))

	this.So(this.Serve("GET", "/bad-request"), should.Equal, "Not Here")
	this.So(this.response.Code, should.Equal, http.StatusNotFound)
}

func (this *RoutesBuilderFixture) Test405MethodNotAllowedCallsRegisteredHandler() {
	this.handler = New(
		GET("/hi", NewFakeHandler("GET"), NewFakeHandler("GET")),
		MethodNotAllowed(NewFakeNotFoundHelperHandler("Registered Handler")),
	)

	this.So(this.Serve("PATCH", "/hi"), should.Equal, "Registered Handler")
	this.So(this.response.Code, should.Equal, http.StatusNotFound)
}

func (this *RoutesBuilderFixture) TestURLParamsSavedInRequestContext() {
	handler := NewFakeHandler("GET")
	this.handler = New(GET("/hi/:id/there", handler))

	this.Serve("GET", "/hi/smarty/there")
	this.So(stringID(handler.request), should.Equal, "smarty")
}

func stringID(request *http.Request) string {
	context := request.Context()
	opaque := context.Value(httprouter.ParamsKey)
	params, _ := opaque.(httprouter.Params)
	return params.ByName("id")
}

func (this *RoutesBuilderFixture) TestPanicHandler() {
	handler := &HandlerThatPanics{expected: errors.New("expected panic")}
	this.handler = New(
		Register("GET", "/", handler),
		Panic(handler.Recover),
	)

	this.So(this.Serve("GET", "/"), should.Equal, "Panic Handled")
	this.So(handler.actual, should.Equal, handler.expected)
}

type HandlerThatPanics struct {
	expected error
	actual   interface{}
}

func (this *HandlerThatPanics) Install(handler http.Handler) {
}
func (this *HandlerThatPanics) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	panic(this.expected)
}
func (this *HandlerThatPanics) Recover(response http.ResponseWriter, _ *http.Request, err interface{}) {
	this.actual = err
	response.Write([]byte("Panic Handled"))
}

/////////////////////////////////////////////////////////////////

func TestHandlerChainFixture(t *testing.T) {
	gunit.Run(new(HandlerChainFixture), t)
}

type HandlerChainFixture struct {
	*gunit.Fixture
	request  *http.Request
	response *httptest.ResponseRecorder
}

func (this *HandlerChainFixture) Setup() {
	this.request = httptest.NewRequest("GET", "/", nil)
	this.response = httptest.NewRecorder()
}

func (this *HandlerChainFixture) Serve(handler http.Handler) string {
	handler.ServeHTTP(this.response, this.request)
	return strings.TrimSpace(this.response.Body.String())
}

func (this *HandlerChainFixture) TestHandlersAreChainedTogetherInTheCorrectOrder() {
	handler := chainN([]NestingHandler{
		NewFakeHandler("1. sanitize "),
		NewFakeHandler("2. authentication "),
		NewFakeHandler("3. application "),
	})

	this.So(this.Serve(handler), should.Equal, strings.TrimSpace(`
1. sanitize 2. authentication 3. application`))
}

//////////////////////////////////////////////////////////////////////////

type FakeHandler struct {
	id      string
	inner   http.Handler
	request *http.Request
}

func NewFakeHandler(id string) *FakeHandler {
	return &FakeHandler{id: id}
}

func (this *FakeHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	this.request = request
	fmt.Fprint(response, this.id)
	if this.inner != nil {
		this.inner.ServeHTTP(response, request)
	}
}

func (this *FakeHandler) Install(inner http.Handler) {
	this.inner = inner
}

////////////////////////////////////////////////////////////////////////////

type FakeNotFoundHelperHandler struct {
	value string
}

func NewFakeNotFoundHelperHandler(id string) *FakeNotFoundHelperHandler {
	return &FakeNotFoundHelperHandler{value: id}
}

func (this *FakeNotFoundHelperHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusNotFound)
	fmt.Fprint(response, this.value)
}
