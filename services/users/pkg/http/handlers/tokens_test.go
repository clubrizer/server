package handlers_test

import (
	"errors"
	"github.com/clubrizer/services/users/internal/authenticator/clubrizer"
	"github.com/clubrizer/services/users/internal/util/appconfig"
	"github.com/clubrizer/services/users/pkg/http/handlers"
	mocks "github.com/clubrizer/services/users/pkg/mocks/http/handlers"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type fields struct {
	jwtConfig appconfig.Jwt
	auth      func(m *mocks.TokenAuthenticator)
	tokener   func(m *mocks.Tokener)
}

func TestToken_RefreshTokens(t *testing.T) {
	accessTokenHeaderName := "access-token"
	refreshTokenCookieName := "access-token"
	accessToken := "someAccessToken"
	refreshToken := "someRefreshToken"

	userID := int64(1)
	user := &clubrizer.User{ID: userID}
	jwtConfig := appconfig.Jwt{
		AccessToken: appconfig.JwtAccessTokenConfig{HeaderName: accessTokenHeaderName},
		RefreshToken: appconfig.JwtRefreshTokenConfig{
			Cookie: appconfig.JwtRefreshTokenCookieConfig{Name: refreshTokenCookieName},
		},
	}
	request, err := http.NewRequest(http.MethodPost, "/refresh", nil)
	request.Header.Set(accessTokenHeaderName, accessToken)
	if err != nil {
		t.Errorf("failed to create http request: %v", err)
	}
	request.AddCookie(&http.Cookie{Name: refreshTokenCookieName, Value: refreshToken})

	tests := []struct {
		name             string
		fields           fields
		request          *http.Request
		wantStatus       int
		wantAccessToken  string
		wantRefreshToken string
	}{
		{
			name: "ok",
			fields: fields{
				jwtConfig: jwtConfig,
				auth: func(m *mocks.TokenAuthenticator) {
					m.EXPECT().Get(userID).Once().Return(user, nil)
				},
				tokener: func(m *mocks.Tokener) {
					m.EXPECT().ValidateRefreshTokenAndGetUserID(refreshToken).Once().Return(userID, nil)
					m.EXPECT().GenerateAccessToken(user).Once().Return(accessToken, nil)
					m.EXPECT().GenerateRefreshToken(user).Once().Return(refreshToken, time.Now(), nil)
				},
			},
			request:          request,
			wantStatus:       http.StatusOK,
			wantAccessToken:  accessToken,
			wantRefreshToken: refreshToken,
		},
		{
			name: "invalid refresh token",
			fields: fields{
				jwtConfig: jwtConfig,
				tokener: func(m *mocks.Tokener) {
					m.EXPECT().ValidateRefreshTokenAndGetUserID(refreshToken).Once().Return(userID, errors.New("error"))
				},
			},
			request:    request,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "invalid user",
			fields: fields{
				jwtConfig: jwtConfig,
				auth: func(m *mocks.TokenAuthenticator) {
					m.EXPECT().Get(userID).Once().Return(nil, errors.New("error"))
				},
				tokener: func(m *mocks.Tokener) {
					m.EXPECT().ValidateRefreshTokenAndGetUserID(refreshToken).Once().Return(userID, nil)
				},
			},
			request:    request,
			wantStatus: http.StatusNotFound,
		},
		{
			name: "failed to generate access token",
			fields: fields{
				jwtConfig: jwtConfig,
				auth: func(m *mocks.TokenAuthenticator) {
					m.EXPECT().Get(userID).Once().Return(user, nil)
				},
				tokener: func(m *mocks.Tokener) {
					m.EXPECT().ValidateRefreshTokenAndGetUserID(refreshToken).Once().Return(userID, nil)
					m.EXPECT().GenerateAccessToken(user).Once().Return("", errors.New("error"))
				},
			},
			request:    request,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "failed to generate refresh token",
			fields: fields{
				jwtConfig: jwtConfig,
				auth: func(m *mocks.TokenAuthenticator) {
					m.EXPECT().Get(userID).Once().Return(user, nil)
				},
				tokener: func(m *mocks.Tokener) {
					m.EXPECT().ValidateRefreshTokenAndGetUserID(refreshToken).Once().Return(userID, nil)
					m.EXPECT().GenerateAccessToken(user).Once().Return(accessToken, nil)
					m.EXPECT().GenerateRefreshToken(user).Once().Return("", time.Now(), errors.New("error"))
				},
			},
			request:    request,
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := mocks.NewTokenAuthenticator(t)
			if tt.fields.auth != nil {
				tt.fields.auth(auth)
			}
			tokener := mocks.NewTokener(t)
			if tt.fields.tokener != nil {
				tt.fields.tokener(tokener)
			}

			recorder := httptest.NewRecorder()
			h := handlers.NewTokenHandler(tt.fields.jwtConfig, auth, tokener)
			h.RefreshTokens()(recorder, tt.request)

			if !reflect.DeepEqual(recorder.Code, tt.wantStatus) {
				t.Errorf("ValidateAccessToken() status = %d, want %d", recorder.Code, tt.wantStatus)
			}
			if len(tt.wantAccessToken) > 0 {
				token := recorder.Header().Get(tt.fields.jwtConfig.AccessToken.HeaderName)
				if !reflect.DeepEqual(token, tt.wantAccessToken) {
					t.Errorf("Authenticate() access token = %s, want %s", token, tt.wantAccessToken)
				}
			}

			if len(tt.wantRefreshToken) > 0 {
				token := recorder.Result().Cookies()[0].Value
				if !reflect.DeepEqual(token, tt.wantRefreshToken) {
					t.Errorf("Authenticate() refresh token = %s, want %s", token, tt.wantAccessToken)
				}
			}
		})
	}
}

func TestToken_ValidateAccessToken(t *testing.T) {
	accessTokenHeaderName := "access-token"
	accessToken := "someAccessToken"

	request, err := http.NewRequest(http.MethodPost, "/validate", nil)
	request.Header.Set(accessTokenHeaderName, accessToken)
	if err != nil {
		t.Errorf("failed to create http request: %v", err)
	}

	tests := []struct {
		name       string
		fields     fields
		request    *http.Request
		wantStatus int
	}{
		{
			name: "ok",
			fields: fields{
				jwtConfig: appconfig.Jwt{
					AccessToken: appconfig.JwtAccessTokenConfig{
						HeaderName: accessTokenHeaderName,
					},
				},
				tokener: func(m *mocks.Tokener) {
					m.EXPECT().ValidateAccessToken(accessToken).Once().Return(nil)
				},
			},
			request:    request,
			wantStatus: http.StatusOK,
		},
		{
			name: "no access token",
			fields: fields{
				jwtConfig: appconfig.Jwt{},
			},
			request:    request,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "invalid access token",
			fields: fields{
				jwtConfig: appconfig.Jwt{
					AccessToken: appconfig.JwtAccessTokenConfig{
						HeaderName: accessTokenHeaderName,
					},
				},
				tokener: func(m *mocks.Tokener) {
					m.EXPECT().ValidateAccessToken(accessToken).Once().Return(errors.New("hi"))
				},
			},
			request:    request,
			wantStatus: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := mocks.NewTokenAuthenticator(t)
			if tt.fields.auth != nil {
				tt.fields.auth(auth)
			}
			tokener := mocks.NewTokener(t)
			if tt.fields.tokener != nil {
				tt.fields.tokener(tokener)
			}

			recorder := httptest.NewRecorder()
			h := handlers.NewTokenHandler(tt.fields.jwtConfig, auth, tokener)
			h.ValidateAccessToken()(recorder, tt.request)

			if !reflect.DeepEqual(recorder.Code, tt.wantStatus) {
				t.Errorf("ValidateAccessToken() status = %d, want %d", recorder.Code, tt.wantStatus)
			}
		})
	}
}
