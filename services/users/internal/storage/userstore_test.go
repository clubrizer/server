package storage_test

import (
	"context"
	"database/sql"
	"errors"
	"github.com/clubrizer/server/pkg/storageutils"
	"github.com/clubrizer/services/users/internal/authenticator/google"
	"github.com/clubrizer/services/users/internal/storage"
	"github.com/clubrizer/services/users/internal/util/enums/approvalstate"
	"github.com/clubrizer/services/users/internal/util/enums/role"
	"reflect"
	"testing"

	"github.com/pashagolub/pgxmock/v2"
)

func TestUserEditor_Create(t *testing.T) {
	u := &google.User{
		Issuer:     "google",
		ID:         "123",
		GivenName:  "john",
		FamilyName: "doe",
		Email:      "hi@john.doe",
	}

	type prepareFields struct {
		db func(m pgxmock.PgxConnIface)
	}
	type args struct {
		googleUser *google.User
		isAdmin    bool
	}
	tests := []struct {
		name          string
		prepareFields prepareFields
		args          args
		want          *storage.User
		wantErr       error
	}{
		{
			name: "create user",
			prepareFields: prepareFields{db: func(m pgxmock.PgxConnIface) {
				m.
					ExpectQuery("^insert into users (.+) values (.+) returning id$").
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(int64(1)))
			}},
			args: args{
				googleUser: u,
				isAdmin:    false,
			},
			want: &storage.User{
				ID:         1,
				GivenName:  u.GivenName,
				FamilyName: u.FamilyName,
				Email:      u.Email,
				Picture:    u.Picture,
				IsAdmin:    false,
				TeamClaims: nil,
			},
		},
		{
			name: "create admin",
			prepareFields: prepareFields{db: func(m pgxmock.PgxConnIface) {
				m.
					ExpectQuery("^insert into users (.+) values (.+) returning id$").
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(int64(1)))
			}},
			args: args{
				googleUser: u,
				isAdmin:    true,
			},
			want: &storage.User{
				ID:         1,
				GivenName:  u.GivenName,
				FamilyName: u.FamilyName,
				Email:      u.Email,
				Picture:    u.Picture,
				IsAdmin:    true,
				TeamClaims: nil,
			},
		},
		{
			name: "query fails",
			prepareFields: prepareFields{db: func(m pgxmock.PgxConnIface) {
				m.
					ExpectQuery("^insert into users (.+) values (.+) returning id$").
					WillReturnError(errors.New("err"))
			}},
			args: args{
				googleUser: u,
				isAdmin:    false,
			},
			wantErr: storageutils.ErrUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := pgxmock.NewConn()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close(context.Background())

			if tt.prepareFields.db != nil {
				tt.prepareFields.db(db)
			}

			s := storage.NewTestUserEditor(db)

			got, err := s.Create(tt.args.googleUser, tt.args.isAdmin)
			if err != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserEditor_GetFromExternalID(t *testing.T) {
	issuer := "google"
	externalID := "123"
	u := storage.User{
		ID:         1,
		GivenName:  "john",
		FamilyName: "doe",
		Email:      "john@doe.com",
		IsAdmin:    false,
		TeamClaims: []storage.TeamClaim{
			{
				Team: storage.Team{
					Name:        "team1",
					DisplayName: "Team 1",
					Description: "first team",
					Logo:        "super fancy logo url",
				},
				ApprovalState: approvalstate.Approved,
				Role:          role.Admin,
			},
		},
	}

	type prepareFields struct {
		db func(m pgxmock.PgxConnIface)
	}
	type args struct {
		issuer     string
		externalID string
	}
	expectedSQL := "^select (.+) from users (.+) where u.issuer = \\$1 and u.external_id = \\$2$"
	tests := []struct {
		name          string
		prepareFields prepareFields
		args          args
		want          *storage.User
		wantErr       error
	}{
		{
			name: "get user",
			prepareFields: prepareFields{db: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(expectedSQL).WithArgs(issuer, externalID).WillReturnRows(getRowsForUser(&u))
			}},
			args: args{
				issuer:     issuer,
				externalID: externalID,
			},
			want: &u,
		},
		{
			name: "not found",
			prepareFields: prepareFields{db: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(expectedSQL).WithArgs(issuer, externalID).WillReturnError(sql.ErrNoRows)
			}},
			args: args{
				issuer:     issuer,
				externalID: externalID,
			},
			wantErr: storageutils.ErrNotFound,
		},
		{
			name: "not found (scan)",
			prepareFields: prepareFields{db: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(expectedSQL).WithArgs(issuer, externalID).WillReturnRows(getRowsForUser(nil))
			}},
			args: args{
				issuer:     issuer,
				externalID: externalID,
			},
			wantErr: storageutils.ErrNotFound,
		},
		{
			name: "user query failed",
			prepareFields: prepareFields{db: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(expectedSQL).WithArgs(issuer, externalID).WillReturnError(errors.New("unknown"))
			}},
			args: args{
				issuer:     issuer,
				externalID: externalID,
			},
			wantErr: storageutils.ErrUnknown,
		},
		{
			name: "empty rows",
			prepareFields: prepareFields{db: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(expectedSQL).WithArgs(issuer, externalID).WillReturnRows(pgxmock.NewRows([]string{}))
			}},
			args: args{
				issuer:     issuer,
				externalID: externalID,
			},
			wantErr: storageutils.ErrNotFound,
		},
		{
			name: "scan failed",
			prepareFields: prepareFields{db: func(m pgxmock.PgxConnIface) {
				m.
					ExpectQuery(expectedSQL).
					WithArgs(issuer, externalID).
					WillReturnRows(pgxmock.NewRows([]string{"unknown"}).AddRow("just an invalid value"))
			}},
			args: args{
				issuer:     issuer,
				externalID: externalID,
			},
			wantErr: storageutils.ErrScanFailed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := pgxmock.NewConn()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close(context.Background())

			if tt.prepareFields.db != nil {
				tt.prepareFields.db(db)
			}

			s := storage.NewTestUserEditor(db)

			got, err := s.GetFromExternalID(tt.args.issuer, tt.args.externalID)
			if err != tt.wantErr {
				t.Errorf("GetFromExternalID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFromExternalID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserEditor_GetFromId(t *testing.T) {
	id := int64(324)
	u := storage.User{
		ID:         id,
		GivenName:  "john",
		FamilyName: "doe",
		Email:      "john@doe.com",
		IsAdmin:    false,
		TeamClaims: []storage.TeamClaim{
			{
				Team: storage.Team{
					Name:        "team1",
					DisplayName: "Team 1",
					Description: "first team",
					Logo:        "super fancy logo url",
				},
				ApprovalState: approvalstate.Approved,
				Role:          role.Admin,
			},
		},
	}

	type prepareFields struct {
		db func(m pgxmock.PgxConnIface)
	}
	type args struct {
		id int64
	}
	expectedSQL := "^select (.+) from users (.+) where u.id = \\$1$"
	tests := []struct {
		name          string
		prepareFields prepareFields
		args          args
		want          *storage.User
		wantErr       error
	}{
		{
			name: "get user",
			prepareFields: prepareFields{db: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(expectedSQL).WithArgs(id).WillReturnRows(getRowsForUser(&u))
			}},
			args: args{id: id},
			want: &u,
		},
		{
			name: "not found",
			prepareFields: prepareFields{db: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(expectedSQL).WithArgs(id).WillReturnError(sql.ErrNoRows)
			}},
			args:    args{id: id},
			wantErr: storageutils.ErrNotFound,
		},
		{
			name: "not found (scan)",
			prepareFields: prepareFields{db: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(expectedSQL).WithArgs(id).WillReturnRows(getRowsForUser(nil))
			}},
			args:    args{id: id},
			wantErr: storageutils.ErrNotFound,
		},
		{
			name: "user query failed",
			prepareFields: prepareFields{db: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(expectedSQL).WithArgs(id).WillReturnError(errors.New("unknown"))
			}},
			args:    args{id: id},
			wantErr: storageutils.ErrUnknown,
		},
		{
			name: "empty rows",
			prepareFields: prepareFields{db: func(m pgxmock.PgxConnIface) {
				m.ExpectQuery(expectedSQL).WithArgs(id).WillReturnRows(pgxmock.NewRows([]string{}))
			}},
			args:    args{id: id},
			wantErr: storageutils.ErrNotFound,
		},
		{
			name: "scan failed",
			prepareFields: prepareFields{db: func(m pgxmock.PgxConnIface) {
				m.
					ExpectQuery(expectedSQL).
					WithArgs(id).
					WillReturnRows(pgxmock.NewRows([]string{"unknown"}).AddRow("just an invalid value"))
			}},
			args:    args{id: id},
			wantErr: storageutils.ErrScanFailed,
		},
		{
			name: "invalid role",
			prepareFields: prepareFields{db: func(m pgxmock.PgxConnIface) {
				u.TeamClaims[0].Role = "invalid!"
				m.ExpectQuery(expectedSQL).WithArgs(id).WillReturnRows(getRowsForUser(&u))
			}},
			args:    args{id: id},
			wantErr: storageutils.ErrScanFailed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := pgxmock.NewConn()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close(context.Background())

			if tt.prepareFields.db != nil {
				tt.prepareFields.db(db)
			}

			s := storage.NewTestUserEditor(db)

			got, err := s.GetFromID(tt.args.id)
			if err != tt.wantErr {
				t.Errorf("GetFromID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFromID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getStringValue(t *testing.T) {
	type args struct {
		input sql.NullString
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid string",
			args: args{input: sql.NullString{
				String: "valid string :)",
				Valid:  true,
			}},
			want: "valid string :)",
		},
		{
			name: "invalid string",
			args: args{input: sql.NullString{
				String: "invalid string :(",
				Valid:  false,
			}},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := storage.GetStringValue(tt.args.input); got != tt.want {
				t.Errorf("getStringValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getRowsForUser(user *storage.User) *pgxmock.Rows {
	rows := pgxmock.NewRows([]string{
		"id", "given_name", "family_name", "email", "picture", "is_admin",
		"name", "display_name", "description", "logo",
		"role", "approval_state",
	})
	if user == nil {
		return rows
	}
	if user.TeamClaims == nil {
		rows.AddRow(
			user.ID, user.GivenName, user.FamilyName, user.Email, user.Picture, user.IsAdmin,
			nil, nil, nil, nil,
			nil, nil,
		)
	}
	for _, teamClaim := range user.TeamClaims {
		rows.AddRow(
			user.ID, user.GivenName, user.FamilyName, user.Email, user.Picture, user.IsAdmin,
			teamClaim.Team.Name, teamClaim.Team.DisplayName, teamClaim.Team.Description, teamClaim.Team.Logo,
			teamClaim.Role, teamClaim.ApprovalState,
		)
	}

	return rows
}
