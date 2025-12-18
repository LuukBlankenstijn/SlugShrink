package types

type AuthStrategyType string

const (
	AuthStrategyAuthless AuthStrategyType = "authless"
	AuthStrategyBasic    AuthStrategyType = "basic"
)

var AuthStrategyFactories = map[AuthStrategyType]func() AuthStrategy{
	AuthStrategyAuthless: func() AuthStrategy { return &AuthlessStrategy{} },
	AuthStrategyBasic:    func() AuthStrategy { return &BasicStrategy{} },
}
