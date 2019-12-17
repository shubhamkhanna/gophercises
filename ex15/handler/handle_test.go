package handler

import (
	"net/http"
	"net/http/httptest"
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/mock"
)

var mux = http.NewServeMux()

func TestSourceCodeHandler(t *testing.T) {
	mux.HandleFunc("/debug/", SourceCodeHandler)
	t.Run("it checks valid scenario", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/debug/?line=24&path=%2Fhome%2Fshubham%2F.gvm%2Fgos%2Fgo1.13.4%2Fsrc%2Fruntime%2Fdebug%2Fstack.go", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := DevMiddleware(mux)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Something went worng! status is not OK")
		}
	})

	t.Run("it checks invalid line parameter", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/debug/?line=invalid&path=%2Fhome%2Fshubham%2F.gvm%2Fgos%2Fgo1.13.4%2Fsrc%2Fruntime%2Fdebug%2Fstack.go", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := DevMiddleware(mux)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Something went worng! status is not OK")
		}
	})

	t.Run("it checks invalid path parameter", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/debug/?line=24&path=invalid", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := DevMiddleware(mux)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("Something went worng! it should throw internal server error.")
		}
	})
}

func TestPanicDemo(t *testing.T) {
	mux.HandleFunc("/panic/", PanicDemo)
	t.Run("it checks valid scenario", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/panic/", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := DevMiddleware(mux)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("Something went worng! it should throw internal server error.")
		}
	})

	t.Run("it mocks makeLinks function output", func(t *testing.T) {
		testObject := new(MyTestObject)
		stack := debug.Stack()
		testObject.On("makeLinks", string(stack)).Return(string(stack))
		testObject.makeLinks(string(stack))
		testObject.AssertExpectations(t)
	})
}

type MyTestObject struct {
	mock.Mock
}

func (o *MyTestObject) makeLinks(stack string) string {
	args := o.Called(stack)
	return args.String()
}

func TestHello(t *testing.T) {
	mux.HandleFunc("/", Hello)
	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler := DevMiddleware(mux)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Something went worng! status is not OK")
	}
}
