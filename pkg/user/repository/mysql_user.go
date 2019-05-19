package repository

import (
	"database/sql"
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/models"
	"github.com/skyerus/riptides-go/pkg/user"
)

type mysqlUserRepository struct {
	Conn *sql.DB
}

func NewMysqlUserRepository(Conn *sql.DB) user.Repository {
	return &mysqlUserRepository{Conn}
}

func (mysql mysqlUserRepository) Create(user models.User, c chan string, m chan map[string]bool, e chan error) customError.Error {
	stmtIns, err := mysql.Conn.Prepare("INSERT INTO user (username, email, password, salt, avatar, bio) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return customError.NewGenericHttpError(err)
	}
	defer stmtIns.Close()
	var exists map[string]bool

	select {
	case user.Password = <-c:
		exists = <-m
	case err := <-e:
		return customError.NewGenericHttpError(err)
	case exists = <-m:
		user.Password = <-c
	}
	if exists["username"] {
		return customError.NewHttpError(409, "A user already exists with this username", nil)
	}
	if exists["email"] {
		return customError.NewHttpError(409, "A user already exists with this email", nil)
	}
	_, err = stmtIns.Exec(user.Username, user.Email, user.Password, user.Salt, user.Avatar, user.Bio)
	if err != nil {
		return customError.NewGenericHttpError(err)
	}

	return nil
}

func (mysql mysqlUserRepository) DoesUserExist(user models.User, m chan map[string]bool, e chan error) {
	var exists bool
	existsMap := make(map[string]bool, 2)
	existsMap["username"] = false
	existsMap["email"] = false

	err := mysql.Conn.QueryRow("SELECT EXISTS(SELECT 1 FROM user WHERE username = ?)", user.Username).Scan(&exists)
	if err != nil {
		e <- err
		return
	}
	if exists {
		existsMap["username"] = true
		m <- existsMap
		return
	}

	err = mysql.Conn.QueryRow("SELECT EXISTS(SELECT 1 FROM user WHERE email = ?)", user.Email).Scan(&exists)
	if err != nil {
		e <- err
		return
	}
	if exists {
		existsMap["email"] = true
	}

	m <- existsMap
	return
}


