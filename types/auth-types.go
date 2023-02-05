package types

type SigninTokens struct {
	Access_token  string `json:"access_token"`
	Refresh_token string `json:"refresh_token"`
}

type GoogleOauthToken struct {
	Access_token string `json:"access_token"`
	Id_token     string `json:"id_token"`
}

type GoogleUserResult struct {
	Id             string
	Email          string
	Verified_email bool
	Name           string
	Given_name     string
	Family_name    string
	Picture        string
	Locale         string
}
