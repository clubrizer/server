// Package tokener is responsible for creating & validating JWT tokens.
package tokener

import (
	"errors"
	"fmt"
	"github.com/clubrizer/server/pkg/apiutils"
	"github.com/clubrizer/services/users/internal/authenticator/clubrizer"
	"github.com/clubrizer/services/users/internal/util/appconfig"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"time"
)

// A Generator is responsible for creating & validating JWT tokens.
type Generator struct {
	jwtConfig appconfig.Jwt
}

// NewGenerator creates a new [Generator]
func NewGenerator(jwtConfig appconfig.Jwt) *Generator {
	return &Generator{jwtConfig}
}

// ValidateRefreshTokenAndGetUserId validates if a given token is valid and if yes, returns the user associated to this
// token.
func (g Generator) ValidateRefreshTokenAndGetUserId(tokenString string) (int64, error) {
	token, err := g.getToken(tokenString)
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		userId, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			return 0, err
		}

		return userId, nil
	}

	return 0, errors.New("invalid JWT token")
}

// ValidateAccessToken checks if the given access token is valid.
func (g Generator) ValidateAccessToken(tokenString string) error {
	token, err := g.getToken(tokenString)
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return nil
	}

	return errors.New("invalid JWT token")
}

// GenerateAccessToken generates a new access token with user claims.
func (g Generator) GenerateAccessToken(user *clubrizer.User) (string, error) {
	expiresAt := jwt.NewNumericDate(time.Now().Add(15 * time.Minute))

	var teams []apiutils.JwtTeamClaim
	for _, c := range user.TeamClaims {
		teams = append(teams, apiutils.JwtTeamClaim{
			Name:          c.Name,
			ApprovalState: string(c.ApprovalState),
			Role:          string(c.Role),
		})
	}

	claims := apiutils.JwtClaims{
		IsAdmin: user.IsAdmin,
		Teams:   teams,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    g.jwtConfig.Issuer,
			Subject:   strconv.FormatInt(user.Id, 10),
			Audience:  []string{"ui", "services"},
			ExpiresAt: expiresAt,
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(g.jwtConfig.AccessToken.Secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// GenerateRefreshToken generates a new refresh token.
func (g Generator) GenerateRefreshToken(user *clubrizer.User) (string, time.Time, error) {
	expiresAt := jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Hour))

	claims := jwt.RegisteredClaims{
		Issuer:    g.jwtConfig.Issuer,
		Subject:   strconv.FormatInt(user.Id, 10),
		ExpiresAt: expiresAt,
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(g.jwtConfig.RefreshToken.Secret))
	if err != nil {
		return "", time.Now(), err
	}

	return signedToken, expiresAt.Time, nil
}

func (g Generator) getToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return errors.New(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"])), nil
		}

		return []byte(g.jwtConfig.RefreshToken.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
