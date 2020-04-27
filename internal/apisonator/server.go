package apisonator

import (
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

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

func NewServer(upstream string, stop chan struct{}) (*Server, error) {
	backend, err := backend.NewBackend(upstream, nil)
	if err != nil {
		return nil, err
	}

	ticker := time.NewTicker(time.Second * 15)
	go func() {
		for {
			select {
			case <-ticker.C:
				backend.Flush()
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()

	s := &Server{
		router:       http.NewServeMux(),
		upstreamPeer: backend,
	}
	s.registerRoutes()
	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
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
			w.Header().Set(rejectionHeader, upstreamResp.ErrorCode)
			w.WriteHeader(statusCode)
			return

		}
		w.WriteHeader(http.StatusOK)
		return
	}
}

func (s *Server) handleAuthRep() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		clientReq := convertRequest(request)
		upstreamResp, err := s.upstreamPeer.AuthRep(clientReq)
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
			w.Header().Set(rejectionHeader, upstreamResp.ErrorCode)
			w.WriteHeader(statusCode)
			return

		}
		w.WriteHeader(http.StatusOK)
		return
	}
}

// convertRequest returns a 3scale request from a http request
func convertRequest(request *http.Request) threescale.Request {
	q, _ := url.QueryUnescape(request.URL.RawQuery)
	values, _ := url.ParseQuery(q)

	return threescale.Request{
		Auth:    getClientAuthFromRequest(values),
		Service: api.Service(values.Get("service_id")),
		Transactions: []api.Transaction{
			{
				Metrics: parseMetricsFromValues(values),
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

var rgx = regexp.MustCompile(`\[(.*?)\]`)

func parseMetricsFromValues(values url.Values) api.Metrics {
	metrics := api.Metrics{}
	for k, v := range values {
		if strings.HasPrefix(k, "usage[") {
			key := rgx.FindStringSubmatch(k)
			if len(key) != 2 {
				continue
			}
			intVal, _ := strconv.Atoi(v[0])
			metrics.Add(key[1], intVal)
		}
	}
	return metrics
}
