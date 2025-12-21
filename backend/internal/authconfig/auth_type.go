package authconfig

type AuthStrategyType string

const (
	AuthStrategyAuthless AuthStrategyType = "authless"
	AuthStrategyBasic    AuthStrategyType = "basic"
	AuthStrategyProxy    AuthStrategyType = "proxy"
)

var AuthStrategyFactories = map[AuthStrategyType]func() AuthStrategy{
	AuthStrategyAuthless: func() AuthStrategy { return &AuthlessStrategy{} },
	AuthStrategyBasic:    func() AuthStrategy { return &BasicStrategy{} },
	AuthStrategyProxy:    func() AuthStrategy { return &ProxyStrategy{} },
}
