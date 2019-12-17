package primitive

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var r = mux.NewRouter()

type tempStruct struct {
	err    error
	osFile os.File
	cmd    []byte
}

func (m *tempStruct) mockTempfile(prefix, ext string, location string) (*os.File, error) {
	return nil, m.err
}

func (m *tempStruct) dummyTransform(image io.Reader, ext string, numShapes int, opts ...func() []string) (io.Reader, error) {
	return bytes.NewBuffer(nil), m.err
}

func (m *tempStruct) dummyExeCmd(args ...string) ([]byte, error) {
	return nil, m.err
}

func TestUpload(t *testing.T) {
	t.Run("it checks in-valid image file scenario: bad request", func(t *testing.T) {
		payload := strings.NewReader("------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"image\"; filename=\"test_images.png\"\r\nContent-Type: image/png\r\n\r\n\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW--")
		req, err := http.NewRequest("POST", "/upload", payload)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		r.HandleFunc("/upload", Upload).Methods("POST")
		r.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	t.Run("it redirects to modify with found status", func(t *testing.T) {
		payload := strings.NewReader("------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"image\"; filename=\"invalid_images.png\"\r\nContent-Type: image/png\r\n\r\n\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW--")
		req, err := http.NewRequest("POST", "/upload", payload)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		req.Header.Add("content-type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
		req.Header.Add("Accept", "*/*")
		r.HandleFunc("/upload", Upload).Methods("POST")
		r.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusFound {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusFound)
		}
	})

	t.Run("it mocks temp file error", func(t *testing.T) {
		// this mocks out the function that mockTempfile() calls
		m := &tempStruct{err: errors.New("failed")}
		mocktempfile = m.mockTempfile
		payload := strings.NewReader("------WebKitFormBoundary7MA4YWxkTrZu0gW\r\nContent-Disposition: form-data; name=\"image\"; filename=\"test_images.png\"\r\nContent-Type: image/png\r\n\r\n\r\n------WebKitFormBoundary7MA4YWxkTrZu0gW--")
		req, err := http.NewRequest("POST", "/upload", payload)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		req.Header.Add("content-type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
		req.Header.Add("Accept", "*/*")
		r.HandleFunc("/upload", Upload).Methods("POST")
		r.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}

		defer func() {
			mocktempfile = tempfile // set back original func at end of test
		}()
	})

}

func TestModify(t *testing.T) {
	var m = &tempStruct{}
	t.Run("it checks for an invalid file error", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/modify/invalid.png", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		r.HandleFunc("/modify/{id}", Modify).Methods("GET")
		r.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	t.Run("it renders to mode choices", func(t *testing.T) {
		mockTransform = m.dummyTransform
		req, err := http.NewRequest("GET", "/modify/test_images.png", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		r.HandleFunc("/modify/{id}", Modify).Methods("GET")
		r.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	})

	t.Run("it checks invalid mode error", func(t *testing.T) {
		mockTransform = m.dummyTransform
		req, err := http.NewRequest("GET", "/modify/test_images.png?mode=invalidnum", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		r.HandleFunc("/modify/{id}", Modify).Methods("GET")
		r.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	t.Run("it renders renderNumShapeChoices on valid mode", func(t *testing.T) {
		mockTransform = m.dummyTransform
		req, err := http.NewRequest("GET", "/modify/test_images.png?mode=1", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		r.HandleFunc("/modify/{id}", Modify).Methods("GET")
		r.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	})

	t.Run("it redirects to image path", func(t *testing.T) {
		mockTransform = m.dummyTransform
		req, err := http.NewRequest("GET", "/modify/test_images.png?mode=1&n=1", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		r.HandleFunc("/modify/{id}", Modify).Methods("GET")
		r.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusFound {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusFound)
		}
	})

	t.Run("it checks for the errors when num shapes not prsent", func(t *testing.T) {
		// this mocks out the function that mockTempfile() calls
		mTmp := &tempStruct{err: errors.New("failed")}
		mocktempfile = mTmp.mockTempfile
		req, err := http.NewRequest("GET", "/modify/test_images.png?&mode=1", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		r.HandleFunc("/modify/{id}", Modify).Methods("GET")
		r.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})

	t.Run("it checks for the errors when mode is not present", func(t *testing.T) {
		// this mocks out the function that mockTempfile() calls
		m = &tempStruct{err: errors.New("failed")}
		mockTransform = m.dummyTransform
		req, err := http.NewRequest("GET", "/modify/test_images.png?&n=1", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		r.HandleFunc("/modify/{id}", Modify).Methods("GET")
		r.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})

	defer func() {
		mocktempfile = tempfile // set back original func at end of test
		mockTransform = Transform
	}()
}

func TestTransform(t *testing.T) {
	r, _ := os.Open("./img/test_images.png")
	invalidR, _ := os.Open("./invalid")
	t.Run("it returns byte buffer", func(t *testing.T) {
		out, _ := Transform(r, ".png", 1, WithMode(ModePolygon))
		assert.Equal(t, reflect.TypeOf(out).String(), "*bytes.Buffer")
	})

	t.Run("it mocks temp file error", func(t *testing.T) {
		// this mocks out the function that mockTempfile() calls
		m := &tempStruct{err: errors.New("failed")}
		mocktPrimativeTmpfile = m.mockTempfile
		out, _ := Transform(r, "invalid", 1, WithMode(ModePolygon))
		assert.Equal(t, out, nil)

		defer func() {
			mocktPrimativeTmpfile = tempfile // set back original func at end of test
		}()

	})

	t.Run("it railses file copy error", func(t *testing.T) {
		out, _ := Transform(invalidR, ".png", 1, WithMode(ModePolygon))
		assert.Equal(t, out, nil)
	})

	t.Run("it railses primitive error", func(t *testing.T) {
		m := &tempStruct{err: errors.New("failed")}
		mockExeCmd = m.dummyExeCmd
		out, _ := Transform(r, ".png", 1, WithMode(ModePolygon))
		assert.Equal(t, out, nil)
	})
}
