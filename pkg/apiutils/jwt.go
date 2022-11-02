package apiutils

import "github.com/golang-jwt/jwt/v4"

type JwtClaims struct {
	IsAdmin bool
	Teams   []JwtTeamClaim
	jwt.RegisteredClaims
}

type JwtTeamClaim struct {
	Name          string
	ApprovalState string
	Role          string
}
