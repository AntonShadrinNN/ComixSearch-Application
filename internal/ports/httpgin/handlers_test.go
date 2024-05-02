package httpgin

import (
	"bytes"
	"comixsearch/internal/app/mocks"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ComicsData struct {
	Comices map[string]string `json:"comices"`
	Error   interface{}       `json:"error"`
}

func TestSearch(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		jsonBody       string
		data           map[string]string
		searchErr      error
		expectedStatus int
	}{
		{
			name:        "No errors",
			queryParams: map[string]string{},
			jsonBody: `{
				"keywords": "fish"
			}`,
			data:           map[string]string{"fish": "https://imgs.xkcd.com/comics/fish.png"},
			expectedStatus: http.StatusOK,
			searchErr:      nil,
		},
		// {
		// 	name:        "Json error",
		// 	queryParams: map[string]string{"limit": "10"},
		// 	jsonBody: `
		// 		keywords fish
		// 	}`,
		// 	data:           map[string]string{},
		// 	expectedStatus: http.StatusBadRequest,
		// 	searchErr:      nil,
		// },
		{
			name:        "Bad query parameters",
			queryParams: map[string]string{"limit": ""},
			jsonBody: `{
			"keywords": "fish"
			}`,
			data:           map[string]string{},
			expectedStatus: http.StatusBadRequest,
			searchErr:      nil,
		},
		{
			name:        "Search error",
			queryParams: map[string]string{"limit": "10"},
			jsonBody: `{
			"keywords": "fish"
			}`,
			data:           map[string]string{},
			expectedStatus: http.StatusInternalServerError,
			searchErr:      fmt.Errorf("some error"),
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockApp := mocks.NewApp(t)
			mockApp.On("Search", mock.Anything, mock.Anything, mock.Anything).Return(tt.data, tt.searchErr).Maybe()
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req, _ := http.NewRequest(http.MethodPost, "/api/v1/search", io.NopCloser(bytes.NewBufferString(tt.jsonBody)))

			q := req.URL.Query()
			for k, v := range tt.queryParams {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()

			c.Request = req

			search(ctx, mockApp)(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			// var data ComicsData

			// json.NewDecoder(w.Body).Decode(&data)
		})
	}
}
