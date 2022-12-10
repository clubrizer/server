package appconfig

// Auth contains all authentication & authorization related configuration.
type Auth struct {
	GoogleClientID string
	Jwt            Jwt
}

// Jwt contains all configuration related to Json Web Tokens (JWT).
type Jwt struct {
	Issuer       string
	AccessToken  JwtAccessTokenConfig
	RefreshToken JwtRefreshTokenConfig
}

// JwtAccessTokenConfig contains all configuration for the JWT access token.
type JwtAccessTokenConfig struct {
	Secret     string
	HeaderName string
}

// JwtRefreshTokenConfig contains all configuration for the JWT refresh token.
type JwtRefreshTokenConfig struct {
	Secret string
	Cookie JwtRefreshTokenCookieConfig
}

// JwtRefreshTokenCookieConfig contains all configuration for the JWT refresh token cookie.
type JwtRefreshTokenCookieConfig struct {
	Name     string
	SameSite string
	Secure   bool
	HTTPOnly bool
}
