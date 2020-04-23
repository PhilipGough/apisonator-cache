package apisonator

import (
	"errors"
	"fmt"
	"github.com/3scale/3scale-istio-adapter/pkg/threescale/backend"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/3scale/3scale-go-client/threescale"
	"github.com/3scale/3scale-go-client/threescale/api"
	convert "github.com/3scale/3scale-go-client/threescale/http"
)

const (
	application  = "app"
	service      = "service"
	user         = "user"
	serviceToken = "st"
	providerKey  = "pk"
)

func Test_handleAuthorize(t *testing.T) {
	tests := []struct {
		name              string
		queryBuilder      func() url.Values
		upstream          *mockApisonator
		expectCode        int
		expectHeaderValue string
	}{
		{
			name: "Test upstream error results in failure",
			queryBuilder: func() url.Values {
				values := url.Values{}
				values.Add("app_id", application)
				values.Add("service_id", service)
				values.Add("service_token", serviceToken)
				return values
			},
			upstream: &mockApisonator{
				callback: func(request threescale.Request) {
					expect := threescale.Request{
						Auth: api.ClientAuth{
							Type:  api.ServiceToken,
							Value: serviceToken,
						},
						Service: service,
						Transactions: []api.Transaction{
							{
								Params: api.Params{
									AppID: application,
								},
								Metrics: api.Metrics{},
							},
						},
					}
					equals(t, expect, request)
				},
				err: errors.New("some error"),
			},
			expectCode: http.StatusBadGateway,
		},
		{
			name: "Test upstream denial fails to authorize",
			queryBuilder: func() url.Values {
				values := url.Values{}
				values.Add("user_key", user)
				values.Add("service_id", service)
				values.Add("provider_key", providerKey)
				return values
			},
			upstream: &mockApisonator{
				callback: func(request threescale.Request) {
					expect := threescale.Request{
						Auth: api.ClientAuth{
							Type:  api.ProviderKey,
							Value: providerKey,
						},
						Service: service,
						Transactions: []api.Transaction{
							{
								Params: api.Params{
									UserKey: user,
								},
								Metrics: api.Metrics{},
							},
						},
					}
					equals(t, expect, request)
				},
				result: &threescale.AuthorizeResult{
					Authorized: false,
					ErrorCode:  "service_id_invalid",
				},
			},
			expectCode:        convert.CodeToStatusCode("service_id_invalid"),
			expectHeaderValue: "service_id_invalid",
		},
		{
			name: "Test success",
			queryBuilder: func() url.Values {
				values := url.Values{}
				values.Add("user_key", user)
				values.Add("service_id", service)
				values.Add("provider_key", providerKey)
				return values
			},
			upstream: &mockApisonator{
				callback: func(request threescale.Request) {
					expect := threescale.Request{
						Auth: api.ClientAuth{
							Type:  api.ProviderKey,
							Value: providerKey,
						},
						Service: service,
						Transactions: []api.Transaction{
							{
								Params: api.Params{
									UserKey: user,
								},
								Metrics: api.Metrics{},
							},
						},
					}
					equals(t, expect, request)
				},
				result: &threescale.AuthorizeResult{Authorized: true},
			},
			expectCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/transactions/authorize.xml", nil)
			req.URL.RawQuery = test.queryBuilder().Encode()

			srv := newTestServer(t, test.upstream)
			w := httptest.NewRecorder()
			srv.router.ServeHTTP(w, req)

			equals(t, test.expectCode, w.Code)
			equals(t, test.expectHeaderValue, w.Header().Get(rejectionHeader))
		})
	}
}

func Test_handleAuthRep(t *testing.T) {
	tests := []struct {
		name              string
		queryBuilder      func() url.Values
		upstream          *mockApisonator
		expectCode        int
		expectHeaderValue string
	}{
		{
			name: "Test upstream error results in failure",
			queryBuilder: func() url.Values {
				values := url.Values{}
				values.Add("app_id", application)
				values.Add("service_id", service)
				values.Add("service_token", serviceToken)
				return values
			},
			upstream: &mockApisonator{
				callback: func(request threescale.Request) {
					expect := threescale.Request{
						Auth: api.ClientAuth{
							Type:  api.ServiceToken,
							Value: serviceToken,
						},
						Service: service,
						Transactions: []api.Transaction{
							{
								Params: api.Params{
									AppID: application,
								},
								Metrics: api.Metrics{},
							},
						},
					}
					equals(t, expect, request)
				},
				err: errors.New("some error"),
			},
			expectCode: http.StatusBadGateway,
		},
		{
			name: "Test upstream denial fails to authorize",
			queryBuilder: func() url.Values {
				values := url.Values{}
				values.Add("user_key", user)
				values.Add("service_id", service)
				values.Add("provider_key", providerKey)
				return values
			},
			upstream: &mockApisonator{
				callback: func(request threescale.Request) {
					expect := threescale.Request{
						Auth: api.ClientAuth{
							Type:  api.ProviderKey,
							Value: providerKey,
						},
						Service: service,
						Transactions: []api.Transaction{
							{
								Params: api.Params{
									UserKey: user,
								},
								Metrics: api.Metrics{},
							},
						},
					}
					equals(t, expect, request)
				},
				result: &threescale.AuthorizeResult{
					Authorized: false,
					ErrorCode:  "service_id_invalid",
				},
			},
			expectCode:        convert.CodeToStatusCode("service_id_invalid"),
			expectHeaderValue: "service_id_invalid",
		},
		{
			name: "Test success",
			queryBuilder: func() url.Values {
				values := url.Values{}
				values.Add("user_key", user)
				values.Add("service_id", service)
				values.Add("provider_key", providerKey)
				return values
			},
			upstream: &mockApisonator{
				callback: func(request threescale.Request) {
					expect := threescale.Request{
						Auth: api.ClientAuth{
							Type:  api.ProviderKey,
							Value: providerKey,
						},
						Service: service,
						Transactions: []api.Transaction{
							{
								Params: api.Params{
									UserKey: user,
								},
								Metrics: api.Metrics{},
							},
						},
					}
					equals(t, expect, request)
				},
				result: &threescale.AuthorizeResult{Authorized: true},
			},
			expectCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/transactions/authrep.xml", nil)
			req.URL.RawQuery = test.queryBuilder().Encode()

			srv := newTestServer(t, test.upstream)
			w := httptest.NewRecorder()
			srv.router.ServeHTTP(w, req)

			equals(t, test.expectCode, w.Code)
			equals(t, test.expectHeaderValue, w.Header().Get(rejectionHeader))
		})
	}
}

func TestNewServer(t *testing.T) {
	srv, err := NewServer("http://3scale.net")
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	_, ok := srv.upstreamPeer.(*backend.Backend)
	if !ok {
		t.Error("unexpected result from type assertion")
	}

}

func Test_convertRequest(t *testing.T) {
	const userKey = "1234"
	encodedQuery := `service_token%3Dst%26service_id%3Dservice%26user_key%3D1234%26usage%5Bhits%5D%3D2%26usage%5Bmits%5D%3D3`
	// decodes to service_token=st&service_id=service&user_key=1234&usage[hits]=2&usage[mits]=3
	req := &http.Request{
		URL: &url.URL{
			RawQuery: encodedQuery,
		},
	}

	expect := threescale.Request{
		Auth: api.ClientAuth{
			Type:  api.ServiceToken,
			Value: serviceToken,
		},
		Service: service,
		Transactions: []api.Transaction{
			{
				Metrics: api.Metrics{
					"hits": 2,
					"mits": 3,
				},
				Params: api.Params{
					UserKey: userKey,
				},
			},
		},
	}
	result := convertRequest(req)
	equals(t, expect, result)
}

func newTestServer(t *testing.T, upstream *mockApisonator) *Server {
	t.Helper()
	srv := &Server{
		router:       http.NewServeMux(),
		upstreamPeer: upstream,
	}
	srv.registerRoutes()
	return srv
}

type mockApisonator struct {
	callback func(request threescale.Request)
	result   *threescale.AuthorizeResult
	err      error
}

func (m *mockApisonator) Authorize(request threescale.Request) (*threescale.AuthorizeResult, error) {
	return m.auth(request)
}

func (m *mockApisonator) AuthRep(request threescale.Request) (*threescale.AuthorizeResult, error) {
	return m.auth(request)
}

func (m mockApisonator) auth(request threescale.Request) (*threescale.AuthorizeResult, error) {
	if m.callback != nil {
		m.callback(request)
	}
	if m.err != nil {
		return nil, m.err
	}
	return m.result, nil
}

func (m mockApisonator) Report(request threescale.Request) (*threescale.ReportResult, error) {
	panic("implement me")
}

func (m mockApisonator) GetPeer() string {
	panic("implement me")
}

// equals fails the test if exp is not equal to act.
func equals(t *testing.T, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		t.FailNow()
	}
}
