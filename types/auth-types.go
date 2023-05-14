package types

// maybe to move it to adquate place
type JwtTokens struct {
	Access_token  string `json:"access_token"`
	Refresh_token string `json:"refresh_token"`
}
