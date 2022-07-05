package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vivalavoka/go-exporter/cmd/server/config"
	"github.com/vivalavoka/go-exporter/cmd/server/storage"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestMetricHandle(t *testing.T) {
	storage.New(config.Config{})
	r := chi.NewRouter()
	handlers := Handlers{}
	handlers.SetRoutes(r)
	ts := httptest.NewServer(r)
	defer ts.Close()

	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "invalid metric type",
			url:  "/update/invalid-type/test-name/10",
			want: want{
				code:        http.StatusNotImplemented,
				response:    "Wrong metric type\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "wrong counter value string",
			url:  "/update/counter/test-name/invalid-type",
			want: want{
				code:        http.StatusBadRequest,
				response:    "Wrong metric value\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "wrong counter value float",
			url:  "/update/counter/test-name/0.01",
			want: want{
				code:        http.StatusBadRequest,
				response:    "Wrong metric value\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "invalid gauge type string",
			url:  "/update/gauge/test-name/invalid-value",
			want: want{
				code:        http.StatusBadRequest,
				response:    "Wrong metric value\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "correct gauge",
			url:  "/update/gauge/test-gauge/0.01",
			want: want{
				code:        http.StatusOK,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "correct counter",
			url:  "/update/counter/test-counter/10",
			want: want{
				code:        http.StatusOK,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, body := testRequest(t, ts, http.MethodPost, tt.url)
			defer response.Body.Close()

			assert.Equal(t, tt.want.code, response.StatusCode)
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.response, string(body))
		})
	}
}
