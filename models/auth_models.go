package models

type AuthClaim struct {
	Sub           string   `mapstructure:"sub"`
	CognitoGroups []string `mapstructure:"cognito:groups"`
	Iss           string   `mapstructure:"iss"`
	Version       int      `mapstructure:"version"`
	ClientId      string   `mapstructure:"client_id"`
	TokenUse      string   `mapstructure:"token_use"`
	Scope         string   `mapstructure:"scope"`
	AuthTime      int      `mapstructure:"auth_time"`
	Exp           int      `mapstructure:"exp"`
	Iat           int      `mapstructure:"iat"`
	Jti           string   `mapstructure:"jti"`
	Username      string   `mapstructure:"username"`
}
