package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricHandle(t *testing.T) {
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
			name: "wrong request path",
			url:  "/invalid-method/counter/test-name/10",
			want: want{
				code:        404,
				response:    "Wrong request path\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "wrong request path",
			url:  "/update",
			want: want{
				code:        404,
				response:    "Wrong request path\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "wrong metric value string",
			url:  "/update/counter/test-name/invalid-type",
			want: want{
				code:        400,
				response:    "Wrong metric value\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "wrong metric value float",
			url:  "/update/counter/test-name/0.01",
			want: want{
				code:        400,
				response:    "Wrong metric value\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "invalid metric type",
			url:  "/update/invalid-type/test-name/10",
			want: want{
				code:        501,
				response:    "Wrong metric type\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "positive case",
			url:  "/update/counter/test-name/10",
			want: want{
				code:        200,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.url, bytes.NewBufferString(""))

			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := http.HandlerFunc(MetricHandle)
			// запускаем сервер
			h.ServeHTTP(w, request)
			result := w.Result()

			fmt.Println(result.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.code, result.StatusCode)
			if result.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, result.Header.Get("Content-Type"))
			}
			// assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			response, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.response, string(response))
		})
	}
}
