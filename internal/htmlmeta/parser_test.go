package htmlmeta

import (
	"context"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected MetaTag
		wantErr  bool
	}{
		{
			name:     "test_empty_reader",
			input:    "",
			expected: MetaTag{},
		},
		{
			name: "test_with_all_attributes",
			input: `<html>
						<head>
							<title>Test Title</title>
							<meta name="description" content="Test Description">
							<meta name="keywords" content="keyword1, keyword2, keyword3">
						</head>
						<body>
							<h1>Hello, World!</h1>
						</body>
					</html>`,
			expected: MetaTag{
				Title:       "Test Title",
				Description: "Test Description",
				Keywords:    []string{"keyword1", "keyword2", "keyword3"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				r := strings.NewReader(tt.input)
				got, err := Parse(context.Background(), r)

				if (err != nil) != tt.wantErr {
					t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if got.Title != tt.expected.Title {
					t.Errorf("Parse() got = %v, want %v", got.Title, tt.expected.Title)
				}

				if got.Description != tt.expected.Description {
					t.Errorf("Parse() got = %v, want %v", got.Description, tt.expected.Description)
				}

				if len(got.Keywords) != len(tt.expected.Keywords) {
					t.Errorf(
						"Parse() got = %v keywords, want %v keywords", len(got.Keywords), len(tt.expected.Keywords),
					)
					return
				}

				for i := range got.Keywords {
					if got.Keywords[i] != tt.expected.Keywords[i] {
						t.Errorf("Parse() got = %v, want %v", got.Keywords[i], tt.expected.Keywords[i])
					}
				}
			},
		)
	}
}
