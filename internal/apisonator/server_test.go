package apisonator

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handleAuthorize(t *testing.T) {
	srv := newServer()
	req := httptest.NewRequest("GET", "/transactions/authorize.xml", nil)
	w := httptest.NewRecorder()

	srv.router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("unexpected status code, wanted 200 but got %d", w.Code)
	}
	decoder := xml.NewDecoder(w.Result().Body)
	xml := &authResponseXML{}
	decoder.Decode(xml)

	if xml.Authorized != true {
		t.Error("expected to be authorized")
	}
}

func Test_handleAuthRep(t *testing.T) {
	srv := newServer()
	req := httptest.NewRequest("GET", "/transactions/authrep.xml", nil)
	w := httptest.NewRecorder()

	srv.router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("unexpected status code, wanted 200 but got %d", w.Code)
	}
	decoder := xml.NewDecoder(w.Result().Body)
	xml := &authResponseXML{}
	decoder.Decode(xml)

	if xml.Authorized != true {
		t.Error("expected to be authorized")
	}
}
