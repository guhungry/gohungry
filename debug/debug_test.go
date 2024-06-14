package debug

import (
	"bytes"
	"math"
	"testing"
)

func TestDump(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected any
	}{
		{
			name:     "when nil",
			input:    nil,
			expected: "<nil>: <nil>" + lineBreak,
		},
		{
			name:     "when string",
			input:    "fdas",
			expected: "string: fdas" + lineBreak,
		},
		{
			name:     "when int",
			input:    1,
			expected: "int: 1" + lineBreak,
		},
		{
			name:     "when float(Pi)",
			input:    math.Pi,
			expected: "float64: 3.141592653589793" + lineBreak,
		},
		{
			name:     "when complex",
			input:    -23.44i,
			expected: "complex128: (0-23.44i)" + lineBreak,
		},
		{
			name:     "when boolean",
			input:    false,
			expected: "bool: false" + lineBreak,
		},
		{
			name: "when struct",
			input: struct {
				name string
				age  int
			}{
				name: "Bee",
				age:  12,
			},
			expected: "struct { name string; age int }: {name:Bee age:12}" + lineBreak,
		},
		{
			name:     "when map[int]string",
			input:    map[int]string{1: "one", 2: "two"},
			expected: "map[int]string: map[1:one 2:two]" + lineBreak,
		},
		{
			name:     "when []string",
			input:    []string{"A", "B", "C"},
			expected: "[]string: [A B C]" + lineBreak,
		},
	}

	for _, test := range tests {
		actual := getActual(test.input)
		if actual != test.expected {
			t.Errorf(`Dump("%v") actual="%s", expected=%s, case: %s`, test.input, actual, test.expected, test.name)
		}
	}
}

func getActual(input any) string {
	// Arrange
	backup := out // Replace standard out during test
	out = new(bytes.Buffer)
	defer func() { out = backup }()

	// Act
	Dump(input)

	// Assert
	return out.(*bytes.Buffer).String()
}
