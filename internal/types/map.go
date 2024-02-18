package types

import "fmt"

type Map struct {
	Keys TypeDef
	Vals TypeDef
}

func (s *Map) Type(uctx UsageCtx) string {
	return fmt.Sprintf("map[%s]%s", s.Keys.Type(uctx), s.Vals.Type(uctx))
}
