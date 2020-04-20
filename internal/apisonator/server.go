package apisonator

import (
	"encoding/xml"
	"net/http"
	"strconv"
)

type server struct {
	router *http.ServeMux
}

type authResponseXML struct {
	Name       xml.Name `xml:",any"`
	Authorized bool     `xml:"authorized,omitempty"`
	Reason     string   `xml:"reason,omitempty"`
	Code       string   `xml:"code,attr,omitempty"`
}

func newServer() *server {
	s := &server{
		router: http.NewServeMux(),
	}
	s.registerRoutes()
	return s
}

func (s *server) registerRoutes() {
	s.router.HandleFunc("/transactions/authorize.xml", s.handleAuthorize())
	s.router.HandleFunc("/transactions/authrep.xml", s.handleAuthRep())
}

func (s *server) handleAuthorize() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		encoder := xml.NewEncoder(w)
		resp := authResponseXML{
			Authorized: true,
			Code:       strconv.Itoa(http.StatusOK),
		}

		w.Header().Set("Content-Type", "application/xml")
		if err := encoder.Encode(&resp); err != nil {
			http.Error(w, "failed to encode to xml", http.StatusInternalServerError)
			return
		}
	}
}

func (s *server) handleAuthRep() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		encoder := xml.NewEncoder(w)
		resp := authResponseXML{
			Authorized: true,
			Code:       strconv.Itoa(http.StatusOK),
		}

		w.Header().Set("Content-Type", "application/xml")
		if err := encoder.Encode(&resp); err != nil {
			http.Error(w, "failed to encode to xml", http.StatusInternalServerError)
			return
		}
	}
}
