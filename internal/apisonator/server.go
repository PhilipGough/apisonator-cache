package apisonator

import (
	"encoding/xml"
	"net/http"
	"strconv"

	"github.com/3scale/3scale-go-client/threescale"
	"github.com/3scale/3scale-go-client/threescale/api"
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

func convertRequest(request *http.Request) threescale.Request {
	return threescale.Request{
		Auth:    getClientAuthFromRequest(request),
		Service: api.Service(request.URL.Query().Get("service_id")),
		Transactions: []api.Transaction{
			{
				Metrics: nil,
				Params: api.Params{
					AppID:    request.URL.Query().Get("app_id"),
					AppKey:   request.URL.Query().Get("app_key"),
					Referrer: request.URL.Query().Get("referrer"),
					UserID:   request.URL.Query().Get("user_id"),
					UserKey:  request.URL.Query().Get("user_key"),
				},
			},
		},
	}

}

func getClientAuthFromRequest(request *http.Request) api.ClientAuth {
	if value := request.URL.Query().Get("service_token"); value != "" {
		return api.ClientAuth{
			Type:  api.ServiceToken,
			Value: value,
		}
	}
	return api.ClientAuth{
		Type:  api.ProviderKey,
		Value: request.URL.Query().Get("provider_key"),
	}
}
