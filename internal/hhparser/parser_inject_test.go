package hhparser

import (
	"testing"
)

func TestInjectCount(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
		wantErr  bool
	}{
		{
			name:     "valid input with value",
			input:    `{"value":123,}`,
			expected: 123,
			wantErr:  false,
		},
		{
			name:     "valid input with value and extra fields",
			input:    `{"value":456, "otherField": "test"}`,
			expected: 456,
			wantErr:  false,
		},
		{
			name:     "input without value field",
			input:    `{"otherField":789}`,
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "input with non-integer value",
			input:    `{"value":"abc"}`,
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "empty input",
			input:    ``,
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "invalid JSON input",
			input:    `{"value":123`,
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count, err := injectCount(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("injectCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if count != tt.expected {
				t.Errorf("injectCount() = %v, want %v", count, tt.expected)
			}
		})
	}
}
