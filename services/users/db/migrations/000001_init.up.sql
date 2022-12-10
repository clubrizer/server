CREATE SEQUENCE IF NOT EXISTS seq_teams;
CREATE TABLE IF NOT EXISTS teams
(
    id           BIGINT              NOT NULL DEFAULT nextval('seq_teams') PRIMARY KEY,
    created_at   TIMESTAMP           NOT NULL DEFAULT now(),
    name         VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(255) UNIQUE NOT NULL,
    description  VARCHAR(1000)                DEFAULT NULL,
    logo         VARCHAR(255)                 DEFAULT NULL
);

CREATE SEQUENCE IF NOT EXISTS seq_users;
CREATE TABLE IF NOT EXISTS users
(
    id          BIGINT       NOT NULL DEFAULT nextval('seq_users') PRIMARY KEY,
    created_at  TIMESTAMP    NOT NULL DEFAULT now(),
    given_name  VARCHAR(255) NOT NULL,
    family_name VARCHAR(255) NOT NULL,
    issuer      VARCHAR(255) NOT NULL,
    external_id VARCHAR(255) NOT NULL,
    email       VARCHAR(255) NOT NULL,
    picture     VARCHAR(255)          DEFAULT NULL,
    is_admin    BOOLEAN      NOT NULL DEFAULT FALSE
);

CREATE UNIQUE INDEX users_external_id_index ON users (external_id);
CREATE UNIQUE INDEX users_email_index ON users (email);

CREATE TYPE APPROVAL_STATE AS ENUM ('pending', 'approved', 'declined');
CREATE TYPE TEAM_ROLE AS ENUM ('member', 'admin');

CREATE SEQUENCE IF NOT EXISTS seq_team_claims;
CREATE TABLE IF NOT EXISTS team_claims
(
    id             BIGINT         NOT NULL      DEFAULT nextval('seq_team_claims') PRIMARY KEY,
    created_at     TIMESTAMP      NOT NULL      DEFAULT now(),
    team_id        BIGINT         NOT NULL REFERENCES teams (id),
    user_id        BIGINT         NOT NULL REFERENCES users (id),
    approval_state APPROVAL_STATE NOT NULL      DEFAULT 'pending',
    role           TEAM_ROLE      NOT NULL      DEFAULT 'member',
    reviewed_by    BIGINT REFERENCES users (id) DEFAULT NULL,
    reviewed_at    TIMESTAMP                    DEFAULT NULL
);

