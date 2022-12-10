package appconfig

type Auth struct {
	GoogleClientId string
	Jwt            Jwt
}

type Jwt struct {
	Issuer       string
	AccessToken  JwtAccessTokenConfig
	RefreshToken JwtRefreshTokenConfig
}

type JwtAccessTokenConfig struct {
	Secret     string
	HeaderName string
}

type JwtRefreshTokenConfig struct {
	Secret string
	Cookie JwtRefreshTokenCookieConfig
}

type JwtRefreshTokenCookieConfig struct {
	Name     string
	SameSite string
	Secure   bool
	HttpOnly bool
}
