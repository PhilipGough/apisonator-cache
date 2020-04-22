package apisonator

import (
	"net/http"
	"net/url"

	"github.com/3scale/3scale-go-client/threescale"
	"github.com/3scale/3scale-go-client/threescale/api"
	convert "github.com/3scale/3scale-go-client/threescale/http"
	"github.com/3scale/3scale-istio-adapter/pkg/threescale/backend"
)

const (
	rejectionHeader = "3scale-rejection-reason"
)

type Server struct {
	router       *http.ServeMux
	upstreamPeer threescale.Client
}

func NewServer(upstream string) (*Server, error) {
	backend, err := backend.NewBackend(upstream, nil)
	if err != nil {
		return nil, err
	}

	s := &Server{
		router:       http.NewServeMux(),
		upstreamPeer: backend,
	}
	s.registerRoutes()
	return s, nil
}

func (s *Server) registerRoutes() {
	s.router.HandleFunc("/transactions/authorize.xml", s.handleAuthorize())
	s.router.HandleFunc("/transactions/authrep.xml", s.handleAuthRep())
}

func (s *Server) handleAuthorize() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		clientReq := convertRequest(request)
		upstreamResp, err := s.upstreamPeer.Authorize(clientReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		if !upstreamResp.Authorized {
			statusCode := convert.CodeToStatusCode(upstreamResp.ErrorCode)
			if statusCode == 0 {
				statusCode = http.StatusNotImplemented
			}
			w.WriteHeader(statusCode)
			w.Header().Set(rejectionHeader, upstreamResp.ErrorCode)
			return

		}
		w.WriteHeader(http.StatusOK)
		return
	}
}

func (s *Server) handleAuthRep() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		clientReq := convertRequest(request)
		upstreamResp, err := s.upstreamPeer.Authorize(clientReq)
		if err != nil {
			// using a generic bad gateway error here
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		if !upstreamResp.Authorized {
			statusCode := convert.CodeToStatusCode(upstreamResp.ErrorCode)
			if statusCode == 0 {
				statusCode = http.StatusNotImplemented
			}
			w.WriteHeader(statusCode)
			w.Header().Set(rejectionHeader, upstreamResp.ErrorCode)
			return

		}
		w.WriteHeader(http.StatusOK)
		return
	}
}

func convertRequest(request *http.Request) threescale.Request {
	values, _ := url.ParseQuery(request.URL.RawQuery)
	return threescale.Request{
		Auth:    getClientAuthFromRequest(values),
		Service: api.Service(values.Get("service_id")),
		Transactions: []api.Transaction{
			{
				Metrics: nil,
				Params: api.Params{
					AppID:    values.Get("app_id"),
					AppKey:   values.Get("app_key"),
					Referrer: values.Get("referrer"),
					UserID:   values.Get("user_id"),
					UserKey:  values.Get("user_key"),
				},
			},
		},
	}

}

func getClientAuthFromRequest(values url.Values) api.ClientAuth {
	if value := values.Get("service_token"); value != "" {
		return api.ClientAuth{
			Type:  api.ServiceToken,
			Value: value,
		}
	}
	return api.ClientAuth{
		Type:  api.ProviderKey,
		Value: values.Get("provider_key"),
	}
}
