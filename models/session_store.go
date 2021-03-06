package models

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/anogowski/gamebase/Godeps/_workspace/src/github.com/lib/pq"
)

type SessionStore interface {
	Find(id string) (*Session, error)
	Delete(sess *Session) error
	Save(sess *Session) error
}

var GlobalSessionStore SessionStore

type PostgresSessionStore struct {
	db *sql.DB
}

func NewPostgresSessionStore() *PostgresSessionStore {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Error opening the database: %q", err)
	}
	sessstore := PostgresSessionStore{db: db}
	if _, err := sessstore.db.Exec("CREATE TABLE IF NOT EXISTS sessions (id VARCHAR(30) PRIMARY KEY, userid VARCHAR(30), expiry TIMESTAMP)"); err != nil {
		log.Fatal("Error creating sessions table: %q", err)
	}
	return &sessstore
}
func (this *PostgresSessionStore) Find(id string) (*Session, error) {
	row := this.db.QueryRow("SELECT id,userid,expiry FROM sessions WHERE id='" + id + "'")
	sess := Session{}
	err := row.Scan(&sess.ID, &sess.UserID, &sess.Expiry)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return &sess, err
	}
	return &sess, nil
}
func (this *PostgresSessionStore) Delete(sess *Session) error {
	if _, err := this.db.Exec("DELETE FROM sessions WHERE id='" + sess.ID + "'"); err != nil {
		return err
	}
	return nil
}
func (this *PostgresSessionStore) Save(sess *Session) error {
	row := this.db.QueryRow("SELECT id FROM sessions WHERE id='" + sess.ID + "'")
	s := &Session{}
	if err := row.Scan(s.ID); err == sql.ErrNoRows {
		if _, err := this.db.Exec("INSERT INTO sessions VALUES('" + sess.ID + "', '" + sess.UserID + "', '" + sess.Expiry.Format("2006-01-02 15:04:05") + "')"); err != nil {
			return err
		}
	} else {
		if _, err := this.db.Exec("UPDATE sessions SET userid='" + sess.UserID + "', expiry='" + sess.Expiry.Format("2006-01-02 15:04:05") + "' WHERE id='" + sess.ID + "'"); err != nil {
			return err
		}
	}
	return nil
}
