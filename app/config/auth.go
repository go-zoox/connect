package config

// Auth ...
type Auth struct {
	// Mode is the auth connect mode
	//	1. oauth2
	//	2. password
	//	3. none
	//  4. openid
	Mode     string `config:"mode"`
	Provider string `config:"provider"`
	//
	IgnorePaths []string `config:"ignore_paths"`
	//
	AllowUsernames []string `config:"allow_usernames"`
}

// AuthPassword ...
type AuthPassword struct {
	Mode    string            `config:"mode"`
	Local   AuthPasswordLocal `config:"local"`
	Service string            `config:"service"`
}

// AuthPasswordLocal ...
type AuthPasswordLocal struct {
	Username string `config:"username"`
	Password string `config:"password"`
}

// AuthOAuth2 ...
type AuthOAuth2 struct {
	Name         string `config:"name"`
	ClientID     string `config:"client_id"`
	ClientSecret string `config:"client_secret"`
	RedirectURI  string `config:"redirect_uri"`
	Scope        string `config:"scope"`
}
