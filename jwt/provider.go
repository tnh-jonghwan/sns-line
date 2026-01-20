package jwt

import "go.uber.org/fx"

// ProvideAccessToken - Access Token Provider
func ProvideAccessToken() string {
	return GetAccessToken()
}

// Module - JWT 모듈
var Module = fx.Options(
	fx.Provide(ProvideAccessToken),
)
