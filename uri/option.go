package uri

import (
	"net/url"
)

type Option func(u *URI)

func WithSchema(val string) Option {
	return func(u *URI) {
		u.Scheme = val
	}
}

func WithUser(val string) Option {
	return func(u *URI) {
		u.User = val
	}
}

func WithPass(val string) Option {
	return func(u *URI) {
		u.Password = val
	}
}

func WithPort(val int) Option {
	return func(u *URI) {
		u.Port = val
	}
}

func WithPath(val string) Option {
	return func(u *URI) {
		u.Path = val
	}
}

func WithHost(val string) Option {
	return func(u *URI) {
		u.Host = val
	}
}

func WithFragment(val string) Option {
	return func(u *URI) {
		u.Fragment = val
	}
}

func WithQuery(kv ...[]string) Option {
	return func(u *URI) {
		q := url.Values{}
		for _, seg := range kv {
			q.Add(seg[0], seg[1])
		}

		u.RawQuery = q.Encode()
	}
}
