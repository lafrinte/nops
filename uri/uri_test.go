package uri

import (
	"testing"
)

func TestURI(t *testing.T) {
	cases := []struct {
		Name   string
		Uri    *URI
		Expect string
	}{
		{
			Name: "full option",
			Uri: New(
				WithSchema("postgres"),
				WithUser("postgres"),
				WithPass("@!12358"),
				WithHost("177.19.0.1"),
				WithPort(8909),
				WithPath("postgres"),
				WithQuery([]string{"pool_size", "32"}, []string{"connection_timeout", "5"}),
				WithFragment("fragment-111"),
			),
			Expect: "postgres://postgres:%40%2112358@177.19.0.1:8909/postgres?connection_timeout=5&pool_size=32#fragment-111",
		},
		{
			Name: "without host/port/path/query/fragment",
			Uri: New(
				WithSchema("postgres"),
				WithUser("postgres"),
				WithPass("@!12358"),
			),
			Expect: "postgres://postgres:%40%2112358@localhost",
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if c.Uri.String() != c.Expect {
				t.Fatalf("expect -> %s got -> %s", c.Expect, c.Uri.String())
			}
		})
	}
}
