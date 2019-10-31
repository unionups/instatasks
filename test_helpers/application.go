package test_helpers

import (
	"github.com/fatih/color"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	// "sync"
)

func LogHttpData(req *http.Request, w *httptest.ResponseRecorder) {
	color.Green("ENDPOINT: %s  %s\n", req.Method, req.URL)
	color.Green("AUTH HEADER: %s\n", req.Header.Get("Authorization"))
	if req.Body != nil {
		b, _ := req.GetBody()
		rb, _ := ioutil.ReadAll(b)
		color.Green("REQUEST BODY: %s\n", rb)
	}
	color.Green("RSEPONSE CODE: %d\n", w.Code)
	color.Green("RSEPONSE BODY: %s\n", w.Body)
}
