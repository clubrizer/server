package handlers

import (
	"github.com/clubrizer/services/users/internal/util/appconfig"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func Test_getRefreshTokenCookie(t *testing.T) {
	expiresAt := time.Now()

	type args struct {
		token     string
		expiresAt time.Time
		jwtConfig appconfig.Jwt
	}
	tests := []struct {
		name string
		args args
		want *http.Cookie
	}{
		{
			name: "SameSite = none",
			args: args{
				token:     "token",
				expiresAt: expiresAt,
				jwtConfig: appconfig.Jwt{
					RefreshToken: appconfig.JwtRefreshTokenConfig{
						Cookie: appconfig.JwtRefreshTokenCookieConfig{
							Name:     "refresh-token",
							SameSite: "none",
							HttpOnly: true,
							Secure:   true,
						},
					},
				},
			},
			want: &http.Cookie{
				Name:     "refresh-token",
				Value:    "token",
				Expires:  expiresAt,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteNoneMode,
			},
		},
		{
			name: "SameSite = lax",
			args: args{
				token:     "token",
				expiresAt: expiresAt,
				jwtConfig: appconfig.Jwt{
					RefreshToken: appconfig.JwtRefreshTokenConfig{
						Cookie: appconfig.JwtRefreshTokenCookieConfig{
							Name:     "refresh-token",
							SameSite: "lax",
							HttpOnly: false,
							Secure:   true,
						},
					},
				},
			},
			want: &http.Cookie{
				Name:     "refresh-token",
				Value:    "token",
				Expires:  expiresAt,
				HttpOnly: false,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			},
		},
		{
			name: "SameSite = strict",
			args: args{
				token:     "token",
				expiresAt: expiresAt,
				jwtConfig: appconfig.Jwt{
					RefreshToken: appconfig.JwtRefreshTokenConfig{
						Cookie: appconfig.JwtRefreshTokenCookieConfig{
							Name:     "refresh-token",
							SameSite: "strict",
							HttpOnly: false,
							Secure:   false,
						},
					},
				},
			},
			want: &http.Cookie{
				Name:     "refresh-token",
				Value:    "token",
				Expires:  expiresAt,
				HttpOnly: false,
				Secure:   false,
				SameSite: http.SameSiteStrictMode,
			},
		},
		{
			name: "SameSite = default",
			args: args{
				token:     "token",
				expiresAt: expiresAt,
				jwtConfig: appconfig.Jwt{
					RefreshToken: appconfig.JwtRefreshTokenConfig{
						Cookie: appconfig.JwtRefreshTokenCookieConfig{
							Name:     "refresh-token",
							SameSite: "something invalid",
							HttpOnly: true,
							Secure:   true,
						},
					},
				},
			},
			want: &http.Cookie{
				Name:     "refresh-token",
				Value:    "token",
				Expires:  expiresAt,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteDefaultMode,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRefreshTokenCookie(tt.args.token, tt.args.expiresAt, tt.args.jwtConfig); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRefreshTokenCookie() = %v, want %v", got, tt.want)
			}
		})
	}
}
