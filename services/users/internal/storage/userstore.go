// Package storage is responsible for handling all database interactions.
package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/clubrizer/server/pkg/config"
	"github.com/clubrizer/server/pkg/log"
	"github.com/clubrizer/server/pkg/storageutils"
	"github.com/clubrizer/services/users/internal/authenticator/google"
	"github.com/clubrizer/services/users/internal/util/enums/approvalstate"
	roleEnum "github.com/clubrizer/services/users/internal/util/enums/role"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userDb interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// A UserEditor is responsible for editing users in the database.
type UserEditor struct {
	db userDb
}

// NewUserEditor creates a new [UserEditor].
func NewUserEditor(postgresConfig config.Postgres) *UserEditor {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s",
		postgresConfig.User,
		postgresConfig.Password,
		postgresConfig.Url,
	)

	connection, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatal(err, "Unable to create database Connection pool")
	}
	log.Info("Created database connection pool")
	return &UserEditor{connection}
}

// GetFromID queries a user by its id.
func (s UserEditor) GetFromID(id int64) (*User, error) {
	rows, err := s.db.Query(
		context.Background(),
		`select u.id, u.given_name, u.family_name, u.email, u.picture, u.is_admin,
				t.name, t.display_name, t.description, t.logo,
				tc.role, tc.approval_state
			from users u
				left join team_claims tc on u.id = tc.user_id
				left join teams t on tc.team_id = t.id
			where u.id = $1`,
		id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, storageutils.ErrNotFound
		}
		log.Error(err, "Failed to query user with team claims [id %d]", id)
		return nil, storageutils.ErrUnknown
	}
	defer rows.Close()

	return scanToUser(rows, fmt.Sprintf("with id [%d]", id))
}

// GetFromExternalID queries a user by its external (e.g. google) id.
func (s UserEditor) GetFromExternalID(issuer string, externalID string) (*User, error) {
	rows, err := s.db.Query(
		context.Background(),
		`select u.id, u.given_name, u.family_name, u.email, u.picture, u.is_admin,
				t.name, t.display_name, t.description, t.logo,
				tc.role, tc.approval_state
			from users u
				left join team_claims tc on u.id = tc.user_id
				left join teams t on tc.team_id = t.id
			where u.issuer = $1 and u.external_id = $2`,
		issuer,
		externalID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, storageutils.ErrNotFound
		}
		log.Error(err, "Failed to query user with team claims [issuer: %s, externalID: %s]", issuer, externalID)
		return nil, storageutils.ErrUnknown
	}
	defer rows.Close()

	return scanToUser(rows, fmt.Sprintf("with team claims [issuer: %s, externalID: %s]", issuer, externalID))
}

// Create creates a new user and stores it in the database.
func (s UserEditor) Create(googleUser *google.User, isAdmin bool) (*User, error) {
	log.Info("Creating user [%s]", getUserDescription(googleUser))

	var userID int64
	err := s.db.QueryRow(
		context.Background(),
		"insert into users (given_name, family_name, issuer, external_id, email, picture, is_admin) values ($1, $2, $3, $4, $5, $6, $7) returning id",
		googleUser.GivenName,
		googleUser.FamilyName,
		googleUser.Issuer,
		googleUser.ID,
		googleUser.Email,
		googleUser.Picture,
		isAdmin,
	).Scan(&userID)
	if err != nil {
		log.Error(err, "Failed to create user [%s]", getUserDescription(googleUser))
		return nil, storageutils.ErrUnknown
	}

	return &User{
		ID:         userID,
		GivenName:  googleUser.GivenName,
		FamilyName: googleUser.FamilyName,
		Email:      googleUser.Email,
		Picture:    googleUser.Picture,
		IsAdmin:    isAdmin,
		TeamClaims: nil,
	}, nil
}

func scanToUser(rows pgx.Rows, userInfo string) (*User, error) {
	var id int64
	var givenName string
	var familyName string
	var email string
	var picture sql.NullString
	var isAdmin bool
	var teamClaims []TeamClaim

	numRows := 0
	for rows.Next() {
		var teamName sql.NullString
		var teamDisplayName sql.NullString
		var teamDescription sql.NullString
		var teamLogo sql.NullString
		var approvalState sql.NullString
		var role sql.NullString

		err := rows.Scan(
			&id, &givenName, &familyName, &email, &picture, &isAdmin,
			&teamName, &teamDisplayName, &teamDescription, &teamLogo,
			&role, &approvalState,
		)
		if err != nil {
			log.Error(err, "Failed to scan row for user %s", userInfo)
			return nil, storageutils.ErrScanFailed
		}

		if teamName.Valid {
			approvalStateValue, ok := approvalstate.FromString(getStringValue(approvalState))
			roleValue, ok := roleEnum.FromString(getStringValue(role))
			if !ok {
				log.Error(err, "Failed to parse enums for user %s", userInfo)
				return nil, storageutils.ErrScanFailed
			}
			teamClaims = append(teamClaims, TeamClaim{
				Team: Team{
					Name:        getStringValue(teamName),
					DisplayName: getStringValue(teamDisplayName),
					Description: getStringValue(teamDescription),
					Logo:        getStringValue(teamLogo),
				},
				ApprovalState: approvalStateValue,
				Role:          roleValue,
			})
		}
		numRows++
	}

	if numRows == 0 {
		log.Info("No rows found for user %s", userInfo)
		return nil, storageutils.ErrNotFound
	}

	return &User{
		ID:         id,
		GivenName:  givenName,
		FamilyName: familyName,
		Email:      email,
		Picture:    getStringValue(picture),
		IsAdmin:    isAdmin,
		TeamClaims: teamClaims,
	}, nil
}

func getUserDescription(user *google.User) string {
	return fmt.Sprintf("'%s %s' (external id: %s, issuer: %s)", user.GivenName, user.FamilyName, user.ID, user.Issuer)
}

func getStringValue(input sql.NullString) string {
	if input.Valid {
		return input.String
	}
	return ""
}
