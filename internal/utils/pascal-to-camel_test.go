package utils

import "testing"

func TestPascalToCamel(t *testing.T) {
	cases := map[string]string{
		"UserName":  "userName",
		"StoreID":   "storeId",
		"URL":       "url",
		"ParsedURL": "parsedUrl",
	}
	for in, want := range cases {
		if got := PascalToCamel(in); got != want {
			t.Errorf("%s → %s, want %s", in, got, want)
		}
	}
}
