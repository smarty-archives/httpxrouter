# httpxrouter
--
    import "github.com/smartystreets/httpxrouter"


## Usage

#### func  New

```go
func New(options ...Option) http.Handler
```

#### type Option

```go
type Option func(*routesBuilder)
```


#### func  DELETE

```go
func DELETE(path string, handlers ...nestingHandler) Option
```

#### func  GET

```go
func GET(path string, handlers ...nestingHandler) Option
```

#### func  HEAD

```go
func HEAD(path string, handlers ...nestingHandler) Option
```

#### func  MethodNotAllowed

```go
func MethodNotAllowed(handler http.Handler) Option
```

#### func  NotFound

```go
func NotFound(handler http.Handler) Option
```

#### func  OPTIONS

```go
func OPTIONS(path string, handlers ...nestingHandler) Option
```

#### func  POST

```go
func POST(path string, handlers ...nestingHandler) Option
```

#### func  PUT

```go
func PUT(path string, handlers ...nestingHandler) Option
```

#### func  Panic

```go
func Panic(callback func(http.ResponseWriter, *http.Request, interface{})) Option
```

#### func  Prepend

```go
func Prepend(handlers ...nestingHandler) Option
```

#### func  Register

```go
func Register(methods, paths string, handlers ...nestingHandler) Option
```
