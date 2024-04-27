package models

type AuthClaim struct {
	Sub           string   `json:"sub"`
	CognitoGroups []string `json:"cognito:groups"`
	Iss           string   `json:"iss"`
	Version       int      `json:"version"`
	ClientId      string   `json:"client_id"`
	TokenUse      string   `json:"token_use"`
	Scope         string   `json:"scope"`
	AuthTime      int      `json:"auth_time"`
	Exp           int      `json:"exp"`
	Iat           int      `json:"iat"`
	Jti           string   `json:"jti"`
	Username      string   `json:"username"`
}
