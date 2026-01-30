package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello world",
			expected: []string{"hello", "world"},
		}, {
			input:    "  Ah, hello  there General Kenobi ",
			expected: []string{"ah,", "hello", "there", "general", "kenobi"},
		},
	}
	for _, c := range cases {
		slice := cleanInput(c.input)
		if len(slice) != len(c.expected) {
			t.Errorf("cleanInput Error: len(slice) does not match \ngot - %v\nexp - %v", slice, c.expected)
			t.FailNow()

		}
		for i := range slice {
			word := slice[i]
			e_word := c.expected[i]
			if word != e_word {
				t.Errorf("cleanInput Error: Words do not match \n got - %v\n exp - %v", word, e_word)
				t.Failed()
			}
		}
	}

}
