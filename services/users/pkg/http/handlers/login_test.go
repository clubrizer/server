package handlers_test

import (
	"errors"
	"github.com/clubrizer/services/users/internal/authenticator/clubrizer"
	"github.com/clubrizer/services/users/internal/authenticator/google"
	"github.com/clubrizer/services/users/internal/util/appconfig"
	"github.com/clubrizer/services/users/pkg/http/handlers"
	mocks "github.com/clubrizer/services/users/pkg/mocks/http/handlers"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestLogin_Authenticate(t *testing.T) {
	request, err := http.NewRequest(http.MethodPost, "/login", nil)
	if err != nil {
		t.Errorf("failed to create http request: %v", err)
	}

	googleUser := &google.User{Id: "123"}
	user := &clubrizer.User{Id: 1}
	jwtConfig := appconfig.Jwt{
		AccessToken:  appconfig.JwtAccessTokenConfig{HeaderName: "access-token"},
		RefreshToken: appconfig.JwtRefreshTokenConfig{Cookie: appconfig.JwtRefreshTokenCookieConfig{Name: "refresh-token"}},
	}

	accessToken := "my access token"
	refreshToken := "my refresh token"

	type fields struct {
		gAuth     func(m *mocks.GoogleAuthenticator)
		auth      func(m *mocks.Authenticator)
		tokener   func(m *mocks.Tokener)
		jwtConfig appconfig.Jwt
	}
	tests := []struct {
		name             string
		fields           fields
		request          *http.Request
		wantStatus       int
		wantAccessToken  string
		wantRefreshToken string
	}{
		{
			name: "no google user",
			fields: fields{
				gAuth: func(m *mocks.GoogleAuthenticator) {
					m.EXPECT().GetUserFromContext(request.Context()).Times(1).Return(nil, false)
				},
			},
			request:    request,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "failed to authenticate",
			fields: fields{
				gAuth: func(m *mocks.GoogleAuthenticator) {
					m.EXPECT().GetUserFromContext(request.Context()).Times(1).Return(googleUser, true)
				},
				auth: func(m *mocks.Authenticator) {
					m.EXPECT().Authenticate(googleUser).Times(1).Return(nil, errors.New("hi, I'm an error"))
				},
			},
			request:    request,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "failed to generate access token",
			fields: fields{
				gAuth: func(m *mocks.GoogleAuthenticator) {
					m.EXPECT().GetUserFromContext(request.Context()).Once().Return(googleUser, true)
				},
				auth: func(m *mocks.Authenticator) {
					m.EXPECT().Authenticate(googleUser).Once().Return(user, nil)
				},
				tokener: func(m *mocks.Tokener) {
					m.EXPECT().GenerateAccessToken(user).Once().Return("", errors.New("error"))
				},
				jwtConfig: jwtConfig,
			},
			request:    request,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "failed to generate refresh token",
			fields: fields{
				gAuth: func(m *mocks.GoogleAuthenticator) {
					m.EXPECT().GetUserFromContext(request.Context()).Once().Return(googleUser, true)
				},
				auth: func(m *mocks.Authenticator) {
					m.EXPECT().Authenticate(googleUser).Once().Return(user, nil)
				},
				tokener: func(m *mocks.Tokener) {
					m.EXPECT().GenerateAccessToken(user).Once().Return("access token", nil)
					m.EXPECT().GenerateRefreshToken(user).Once().Return("", time.Now(), errors.New("error"))
				},
				jwtConfig: jwtConfig,
			},
			request:    request,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "ok",
			fields: fields{
				gAuth: func(m *mocks.GoogleAuthenticator) {
					m.EXPECT().GetUserFromContext(request.Context()).Once().Return(googleUser, true)
				},
				auth: func(m *mocks.Authenticator) {
					m.EXPECT().Authenticate(googleUser).Once().Return(user, nil)
				},
				tokener: func(m *mocks.Tokener) {
					m.EXPECT().GenerateAccessToken(user).Once().Return(accessToken, nil)
					m.EXPECT().GenerateRefreshToken(user).Once().Return(refreshToken, time.Now(), nil)
				},
				jwtConfig: jwtConfig,
			},
			request:          request,
			wantStatus:       http.StatusOK,
			wantAccessToken:  accessToken,
			wantRefreshToken: refreshToken,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gAuth := mocks.NewGoogleAuthenticator(t)
			if tt.fields.gAuth != nil {
				tt.fields.gAuth(gAuth)
			}
			auth := mocks.NewAuthenticator(t)
			if tt.fields.auth != nil {
				tt.fields.auth(auth)
			}
			tokener := mocks.NewTokener(t)
			if tt.fields.tokener != nil {
				tt.fields.tokener(tokener)
			}

			recorder := httptest.NewRecorder()
			h := handlers.NewLoginHandler(tt.fields.jwtConfig, gAuth, auth, tokener)
			h.Authenticate()(recorder, tt.request)

			if !reflect.DeepEqual(recorder.Code, tt.wantStatus) {
				t.Errorf("Authenticate() status = %d, want %d", recorder.Code, tt.wantStatus)
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
