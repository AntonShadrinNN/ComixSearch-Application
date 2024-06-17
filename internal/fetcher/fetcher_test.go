package fetcher

import (
	"bytes"
	"comixsearch/internal/fetcher/mocks"
	"comixsearch/internal/models"
	"context"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetComicesCount(t *testing.T) {
	urlArchive, urlComices := "", ""
	someErr := fmt.Errorf("some error")

	tests := []struct {
		name     string
		jsonBody string
		num      int64
		err      error
	}{
		{
			name:     "Request error",
			jsonBody: "",
			num:      0,
			err:      someErr,
		},
		{
			name: "No Errors",
			jsonBody: `{
				"num": 20
			}
			`,
			num: 20,
			err: nil,
		},
		{
			name: "Json error",
			jsonBody: `
				"num": 20
			
			`,
			num: 0,
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewClient(t)
			fetcher := NewFetcher(urlArchive, urlComices, mockClient)
			resp := http.Response{
				Body: io.NopCloser(bytes.NewBufferString(tt.jsonBody)),
			}
			mockClient.On("Do", mock.Anything).Return(&resp, tt.err)

			num, _ := fetcher.getComicesCount()
			assert.Equal(t, tt.num, num)
		})
	}
}

func TestGetData(t *testing.T) {
	urlArchive, urlComices := "", ""
	someErr := fmt.Errorf("some error")

	tests := []struct {
		name     string
		jsonBody string
		comic    models.Comic
		err      error
		jsonErr  error
	}{
		{
			name: "No error",
			jsonBody: `{
				"num": 1,
				"title": "title",
				"transcript": "transcript",
				"img": "img"
			}`,
			comic: models.Comic{
				Id:      1,
				Title:   "title",
				Content: "transcript",
				Link:    "img",
			},
			err:     nil,
			jsonErr: nil,
		},
		{
			name: "Request error",
			jsonBody: `{
				"num": 1,
				"title": "title",
				"transcript": "transcript",
				"img": "img"
			}`,
			comic:   models.Comic{},
			err:     someErr,
			jsonErr: nil,
		},
		{
			name: "Json parse error",
			jsonBody: `{
				"num": 1
				"tit
			}`,
			comic:   models.Comic{},
			err:     nil,
			jsonErr: someErr,
		},
	}

	ctx := context.Background()
	maxProc := runtime.NumCPU()
	lastId := 0

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewClient(t)
			fetcher := NewFetcher(urlArchive, urlComices, mockClient)
			resp := http.Response{
				Body: io.NopCloser(bytes.NewBufferString(tt.jsonBody)),
			}
			mockClient.On("Do", mock.Anything).Return(&resp, tt.err)

			data, _ := fetcher.GetData(ctx, maxProc, int64(lastId))

			for _, c := range data {
				assert.Equal(t, tt.comic, c)
			}
		})
	}
}
