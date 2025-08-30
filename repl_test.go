package main

import(
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	}{
		{
			input : " hello world ", 
			expected : []string{"hello", "world"},
		},
		{
			input : "Charmander bUlbasAur PIKACHU", 
			expected : []string{"charmander","bulbasaur","pikachu"},
		},
		{
			input : "hello", 
			expected : []string{"hello"},
		},
		{
			input : "  hello  ", 
			expected : []string{"hello"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("clean: %s; got: %v; want: %v", c.input, actual, c.expected)
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Match error at index %d; got: %s; want: %s", i, word, expectedWord)
			}
		}
	}

}
