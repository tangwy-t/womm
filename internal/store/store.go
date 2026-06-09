package store

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type Store struct{ db *sql.DB }

type BadgeState struct {
	GitHubUser string     `json:"github_user"`
	BadgeID    string     `json:"badge_id"`
	Unlocked   bool       `json:"unlocked"`
	UnlockedAt *time.Time `json:"unlocked_at,omitempty"`
	Source     string     `json:"source"`
}

type CacheEntry struct {
	GitHubUser string
	BadgeID    string
	Result     bool
	RawData    string
	ExpiresAt  time.Time
}

const schema = `
CREATE TABLE IF NOT EXISTS badge_state (
    github_user TEXT NOT NULL,
    badge_id    TEXT NOT NULL,
    unlocked    INTEGER NOT NULL DEFAULT 0,
    unlocked_at TEXT,
    source      TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (github_user, badge_id)
);
CREATE TABLE IF NOT EXISTS cert_cache (
    github_user TEXT NOT NULL,
    badge_id    TEXT NOT NULL,
    result      INTEGER NOT NULL,
    raw_data    TEXT NOT NULL DEFAULT '',
    expires_at  TEXT NOT NULL,
    PRIMARY KEY (github_user, badge_id)
);`

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) ClaimBadge(user, badgeID string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := s.db.Exec(
		`INSERT INTO badge_state (github_user,badge_id,unlocked,unlocked_at,source)
		 VALUES(?,?,1,?,'claimed')
		 ON CONFLICT(github_user,badge_id) DO UPDATE SET unlocked=1,unlocked_at=?,source='claimed'`,
		user, badgeID, now, now)
	return err
}

func (s *Store) CertifyBadge(user, badgeID string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := s.db.Exec(
		`INSERT INTO badge_state (github_user,badge_id,unlocked,unlocked_at,source)
		 VALUES(?,?,1,?,'certified')
		 ON CONFLICT(github_user,badge_id) DO UPDATE SET unlocked=1,unlocked_at=?,source='certified'`,
		user, badgeID, now, now)
	return err
}

func (s *Store) IsUnlocked(user, badgeID string) (bool, error) {
	var unlocked int
	err := s.db.QueryRow(
		`SELECT unlocked FROM badge_state WHERE github_user=? AND badge_id=?`,
		user, badgeID).Scan(&unlocked)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return unlocked == 1, nil
}

func (s *Store) GetUserBadges(user string) ([]BadgeState, error) {
	rows, err := s.db.Query(
		`SELECT github_user,badge_id,unlocked,unlocked_at,source FROM badge_state
		 WHERE github_user=? AND unlocked=1`, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var states []BadgeState
	for rows.Next() {
		var bs BadgeState
		var ua sql.NullString
		if err := rows.Scan(&bs.GitHubUser, &bs.BadgeID, &bs.Unlocked, &ua, &bs.Source); err != nil {
			return nil, err
		}
		if ua.Valid {
			t, _ := time.Parse(time.RFC3339, ua.String)
			bs.UnlockedAt = &t
		}
		states = append(states, bs)
	}
	return states, rows.Err()
}

func (s *Store) SetCertCache(user, badgeID string, result bool, rawData string, ttl time.Duration) error {
	expires := time.Now().UTC().Add(ttl).Format(time.RFC3339)
	ri := 0
	if result {
		ri = 1
	}
	_, err := s.db.Exec(
		`INSERT INTO cert_cache (github_user,badge_id,result,raw_data,expires_at)
		 VALUES(?,?,?,?,?)
		 ON CONFLICT(github_user,badge_id) DO UPDATE SET result=?,raw_data=?,expires_at=?`,
		user, badgeID, ri, rawData, expires, ri, rawData, expires)
	return err
}

func (s *Store) GetCertCache(user, badgeID string) (*CacheEntry, bool, error) {
	var result int
	var rawData, expiresStr string
	err := s.db.QueryRow(
		`SELECT result,raw_data,expires_at FROM cert_cache WHERE github_user=? AND badge_id=?`,
		user, badgeID).Scan(&result, &rawData, &expiresStr)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	expires, _ := time.Parse(time.RFC3339, expiresStr)
	if time.Now().UTC().After(expires) {
		return nil, false, nil
	}
	return &CacheEntry{
		GitHubUser: user, BadgeID: badgeID,
		Result: result == 1, RawData: rawData, ExpiresAt: expires,
	}, true, nil
}
