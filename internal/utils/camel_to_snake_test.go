package utils

import "testing"

func TestCamelToSnake(t *testing.T) {
	cases := map[string]string{
		"CreatedAt": "created_at",
		"UserID":    "user_id",
		"URL":       "url",
		"ParsedURL": "parsed_url",
		"userId":    "user_id",
	}
	for in, want := range cases {
		if got := CamelToSnake(in); got != want {
			t.Errorf("%s → %s, want %s", in, got, want)
		}
	}
}
