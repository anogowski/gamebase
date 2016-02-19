package models
import (
	_"github.com/lib/pq"
	"database/sql"
	"errors"
	"log"
	"os"
)

type UserStore interface{
	CreateUser(name, pass string)(*User, error)
	FindUser(id string)(*User,error)
	FindUserByName(name string)(*User, error)
	UpdateUser(user User)error
	Authenticate(name, pass string)(*User, error)
}

var GlobalUserStore UserStore

type PostgresUserStore struct{
	db *sql.DB
}
func NewPostgresUserStore() *PostgresUserStore{
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err!=nil{
		log.Fatal("Error opening the database: %q", err)
	}
	ustore := PostgresUserStore{db:db}
	if _, err :=ustore.db.Exec("CREATE TABLE IF NOT EXISTS users (id VARCHAR PRIMARY KEY, name VARCHAR, password VARCHAR, email VARCHAR)"); err!=nil{
		log.Fatal("Error creating users table: %q", err)
	}
	if _, err :=ustore.db.Exec("CREATE TABLE IF NOT EXISTS friends (id VARCHAR, friendid VARCHAR, PRIMARY KEY(id, friendid))"); err!=nil{
		log.Fatal("Error creating friends table: %q", err)
	}
	return &ustore
}
func (this *PostgresUserStore) CreateUser(name, pass string)(*User, error){
	user, err := this.FindUserByName(name)
	if err!=nil{
		return nil,err
	}
	if user!=nil{
		return nil,errors.New("Username already taken.")
	}
	user = NewUser(name, pass)
	if _,err = this.db.Exec("INSERT INTO users VALUES('"+user.UserId+"', '"+user.UserName+"', '"+user.Password+"', '"+user.Email+"')"); err!=nil{
		return user,err
	}
	return user,nil
}
func (this *PostgresUserStore) FindUser(id string)(*User, error){
	row := this.db.QueryRow("SELECT * FROM users WHERE id="+id)
	user := User{}
	err := row.Scan(&user.UserId, &user.UserName, &user.Password, &user.Email);
	switch{
		case err==sql.ErrNoRows:
			break;
		case err!=nil:
			return &user,err
	}
	return &user,nil
}
func (this *PostgresUserStore) FindUserByName(name string)(*User, error){
	row := this.db.QueryRow("SELECT id, name, password, email FROM users WHERE name='"+name+"'")
	user := User{}
	err := row.Scan(&user.UserId, &user.UserName, &user.Password, &user.Email);
	switch{
		case err==sql.ErrNoRows:
			return nil,nil
		case err!=nil:
			return &user,err
	}
	return &user,nil
}
func (this *PostgresUserStore) UpdateUser(user User)error{
	if _,err := this.db.Exec("UPDATE users SET name='"+user.UserName+"', password='"+user.Password+"', email='"+user.Email+"' WHERE id='"+user.UserId+"'"); err!=nil{
		return err
	}
	return nil
}
func (this *PostgresUserStore) Authenticate(name, pass string)(*User, error){
	user, err := this.FindUserByName(name)
	if err!=nil{
		return nil,err
	}
	if user==nil{
		return nil,errors.New("Incorrect username.")
	}
	if user.CheckPassword(pass){
		return user,nil
	}
	return nil,errors.New("Incorrect password.")
}
