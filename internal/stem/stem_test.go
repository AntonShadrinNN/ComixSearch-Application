package stem

import (
	"comixsearch/internal/models"
	"context"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreproccess(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output []string
	}{
		{
			name:   "Special symbols",
			input:  "{example}*?",
			output: []string{"example"},
		},
		{
			name:   "Punctuation",
			input:  "example,.",
			output: []string{"example"},
		},
		{
			name:   "Several words",
			input:  "example1 example2",
			output: []string{"example1", "example2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := prepocess(tt.input)
			for i := 0; i < len(tokens); i++ {
				assert.Equal(t, tt.output[i], tokens[i])
			}
		})

	}
}
func TestProccess(t *testing.T) {
	tests := []struct {
		name   string
		lang   string
		input  string
		output string
	}{
		{
			name:   "No error",
			lang:   "english",
			input:  "landing",
			output: "land",
		},
		{
			name:   "Unknown language error",
			lang:   "ee",
			input:  "landing",
			output: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStem(tt.lang, true)
			stemmedData, _ := s.process(tt.input)
			assert.Equal(t, tt.output, stemmedData)
		})

	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		name   string
		lang   string
		input  []models.Comic
		output []models.Comic
	}{
		{
			name: "No error",
			lang: "english",
			input: []models.Comic{
				{
					Id:      1,
					Title:   "landing",
					Content: "landing",
					Link:    "",
				},
			},
			output: []models.Comic{
				{
					Id:      1,
					Title:   "land",
					Content: "land",
					Link:    "",
				},
			},
		},
		{
			name: "Unknown language error",
			lang: "ee",
			input: []models.Comic{
				{
					Id:      1,
					Title:   "landing",
					Content: "landing",
					Link:    "",
				},
			},
			output: []models.Comic{
				{
					Id:      1,
					Title:   "landing",
					Content: "landing",
					Link:    "",
				},
			},
		},
	}
	ctx := context.Background()
	maxProc := runtime.NumCPU()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStem(tt.lang, true)
			comices, _ := s.Normalize(ctx, tt.input, maxProc)
			for i := 0; i < len(comices); i++ {
				assert.Equal(t, tt.output[i], comices[i])
			}
		})

	}
}
