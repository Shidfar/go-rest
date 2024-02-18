package types

import "time"

type FuncMeta struct {
	Name          string
	Arguments     Args
	Returns       Args
	RetryAttempts int
	Timeout       time.Duration
}
type UsageCtx string
type Args []Arg
type TypeDef interface {
	Type(useCtx UsageCtx) string
}

type Arg struct {
	Name string
	Type TypeDef
}
