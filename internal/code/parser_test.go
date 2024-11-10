package code

import (
	"reflect"
	"testing"
)

func TestParseCodeBlocks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected []Block
	}{
		{
			name:  "Multiple languages, selects first language (python)",
			input: "```python\nprint('Hello')\n```\n```bash\necho 'World'\n```",
			expected: []Block{
				{Language: "python", Code: "print('Hello')"},
			},
		},
		{
			name:     "No code blocks",
			input:    "This is just plain text",
			expected: []Block{},
		},
		{
			name:  "Empty code block",
			input: "```python\n```",
			expected: []Block{
				{Language: "python", Code: ""},
			},
		},
		{
			name:  "Unclosed code block",
			input: "```python\nprint('Unclosed')",
			expected: []Block{
				{Language: "python", Code: "print('Unclosed')"},
			},
		},
		{
			name:  "Multiple languages, selects first language (bash)",
			input: "```bash\necho 'World'\n```\n```python\nprint('Hello')\n```",
			expected: []Block{
				{Language: "bash", Code: "echo 'World'"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := ParseCodeBlocks(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseCodeBlocks() = %v, want %v", result, tt.expected)
			}
		})
	}
}
