package models

import (
	"database/sql"
	"errors"
	_ "gamebase/Godeps/_workspace/src/github.com/lib/pq"
	"log"
	"os"
)

type TagStore interface {
	CreateTag(name string) error
	GetTags()([]string, error)
	FindGamesByTag(tag string)([]Game, error)
}

var GlobalTagStore TagStore

type PostgresTagStore struct {
	db *sql.DB
}

func NewPostgresTagStore() *PostgresTagStore {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Error opening the database: %q", err)
	}
	tstore := PostgresTagStore{db: db}
	if _, err := tstore.db.Exec("CREATE TABLE IF NOT EXISTS labels (id VARCHAR(50) PRIMARY KEY)"); err != nil {
		log.Fatal("Error creating labels table: %q", err)
	}
	return &tstore
}
func (this *PostgresTagStore) CreateTag(name string) error {
	if _, err := this.db.Exec("INSERT INTO labels VALUES('" + name + "') IF NOT EXISTS"); err != nil {
		return err
	}
	return nil
}
func (this *PostgresTagStore) GetTags() ([]string, error) {
	rows, err := this.db.Query("SELECT * FROM labels")
	if err!=nil{
		return nil,err
	}
	defer rows.Close()
	ret := []string{}
	for rows.Next(){
		var name string
		err := rows.Scan(&name)
		switch {
		case err == sql.ErrNoRows:
			return ret, nil
		case err != nil:
			return ret, err
		default:
			ret = append(ret, name)
		}
	}
	return ret, nil
}
func (this *PostgresTagStore) FindGamesByTag(tag string)([]Game, error){
	return nil,errors.New("TODO: Not implemented")
}
