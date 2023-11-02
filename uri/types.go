package uri

import (
	"net/url"
	"strconv"
	"strings"
)

type URI struct {
	Scheme   string
	User     string
	Password string
	Host     string
	Port     int
	Path     string
	RawQuery string
	Fragment string
	u        *url.URL
}

func (u *URI) String() string {
	if u.u == nil {
		u.parsing()
	}

	return u.u.String()
}

func (u *URI) parsing() {
	bu := url.URL{
		Scheme:   u.Scheme,
		Path:     u.Path,
		User:     url.UserPassword(u.User, u.Password),
		RawQuery: u.RawQuery,
		Fragment: u.Fragment,
	}

	if u.Host == "" {
		u.Host = "localhost"
	}

	if u.Port > 0 {
		bu.Host = strings.Join([]string{u.Host, strconv.Itoa(u.Port)}, ":")
	} else {
		bu.Host = u.Host
	}

	u.u = &bu
}
