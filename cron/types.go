package cron

import (
	"github.com/gorhill/cronexpr"
	"time"
)

type ExprParser struct {
	expr       string
	expression *cronexpr.Expression
}

func (e *ExprParser) Next(t time.Time) time.Time {
	return e.expression.Next(t)
}

func (e *ExprParser) NextN(t time.Time, N uint) []time.Time {
	return e.expression.NextN(t, N)
}

func (e *ExprParser) GetExpr() string {
	return e.expr
}
