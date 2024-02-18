package types

type Any struct {
}

func (t *Any) Type(uctx UsageCtx) string {
	return "any"
}
