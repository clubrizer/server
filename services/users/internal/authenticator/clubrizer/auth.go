// Package clubrizer handles authentication and authorization of Clubrizer users.
package clubrizer

import (
	"github.com/clubrizer/server/pkg/storageutils"
	"github.com/clubrizer/services/users/internal/authenticator/google"
	"github.com/clubrizer/services/users/internal/storage"
	"github.com/clubrizer/services/users/internal/util/appconfig"
)

type userRepository interface {
	GetFromExternalID(issuer string, externalID string) (*storage.User, error)
	Create(user *google.User, isAdmin bool) (*storage.User, error)
	GetFromID(id int64) (*storage.User, error)
}

// Authenticator allows users to authenticate & authorize against Clubrizer.
type Authenticator struct {
	r          userRepository
	initConfig appconfig.Init
}

// NewAuthenticator creates a new [Authenticator] with the given parameters.
func NewAuthenticator(initConfig appconfig.Init, r userRepository) *Authenticator {
	return &Authenticator{r, initConfig}
}

// Authenticate gets the Google user that is currently trying to log in and either returns the existing Clubrizer user
// for that Google user or registers a new Clubrizer user if no matching user exists.
func (a Authenticator) Authenticate(googleUser *google.User) (*User, error) {
	// Get internal user if it exists
	user, err := a.r.GetFromExternalID(googleUser.Issuer, googleUser.ID)
	if err == nil {
		return mapToUser(user), nil
	} else if err != storageutils.ErrNotFound {
		return nil, err
	} // else err == storageutils.ErrNotFound -> user does not exist yet -> create it

	// User does not exist -> create it
	user, err = a.r.Create(googleUser, a.initConfig.AdminEmail == googleUser.Email)
	if err != nil {
		return nil, err
	}

	return mapToUser(user), nil
}

// Get gets the user with the given userID.
func (a Authenticator) Get(userID int64) (*User, error) {
	user, err := a.r.GetFromID(userID)
	if err != nil {
		return nil, err
	}
	return mapToUser(user), nil
}

func mapToUser(user *storage.User) *User {
	u := &User{
		ID:      user.ID,
		IsAdmin: user.IsAdmin,
	}

	for _, teamClaim := range user.TeamClaims {
		u.TeamClaims = append(u.TeamClaims, TeamClaim{
			Name:          teamClaim.Team.Name,
			ApprovalState: teamClaim.ApprovalState,
			Role:          teamClaim.Role,
		})
	}

	return u
}
