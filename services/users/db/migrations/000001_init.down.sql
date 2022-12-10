DROP TABLE IF EXISTS team_claims;
DROP SEQUENCE IF EXISTS seq_team_claims;

DROP TYPE IF EXISTS APPROVAL_STATE;
DROP TYPE IF EXISTS TEAM_ROLE;

DROP INDEX IF EXISTS users_external_id_index;
DROP INDEX IF EXISTS users_email_index;
DROP TABLE IF EXISTS users;
DROP SEQUENCE IF EXISTS seq_users;

DROP TABLE IF EXISTS teams;
DROP SEQUENCE IF EXISTS seq_teams;