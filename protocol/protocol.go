package protocol

const (
	// Version ...
	Version = 70015

	// SrvNodeNetwork This node can be asked for full blocks instead of just headers.
	SrvNodeNetwork = 1
)

// NewUserAgent ...
func NewUserAgent(userAgent string) VarStr {
	return newVarStr(userAgent)
}
