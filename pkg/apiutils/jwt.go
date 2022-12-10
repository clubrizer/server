// Package apiutils provides helpers for interacting with the APIs of services.
package apiutils

import "github.com/golang-jwt/jwt/v4"

// JwtClaims are Clubrizer specific claims to be set on a JWT.
type JwtClaims struct {
	IsAdmin bool
	Teams   []JwtTeamClaim
	jwt.RegisteredClaims
}

// A JwtTeamClaim is a part of the JwtClaims that is used to identify a users role in a team.
type JwtTeamClaim struct {
	Name          string
	ApprovalState string
	Role          string
}
