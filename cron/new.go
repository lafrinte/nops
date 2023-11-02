package cron

import (
	"github.com/gorhill/cronexpr"
)

func New(expr string) *ExprParser {
	return &ExprParser{
		expr:       expr,
		expression: cronexpr.MustParse(expr),
	}
}
