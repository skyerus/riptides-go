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

func (mysql mysqlUserRepository) Create(user models.User) customError.Error {
	stmtIns, err := mysql.Conn.Prepare("INSERT INTO user (username, email, password, salt, avatar, bio) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return customError.NewGenericHttpError(err)
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(user.Username, user.Email, user.Password, user.Salt, user.Avatar, user.Bio)
	if err != nil {
		return customError.NewGenericHttpError(err)
	}

	return nil
}

func (mysql mysqlUserRepository) DoesUserExistWithUsername(username string) (bool, error) {
	var exists bool
	err := mysql.Conn.QueryRow("SELECT EXISTS(SELECT 1 FROM user WHERE username = ?)", username).Scan(&exists)

	return exists, err
}

func (mysql mysqlUserRepository) DoesUserExistWithEmail(email string) (bool, error) {
	var exists bool
	err := mysql.Conn.QueryRow("SELECT EXISTS(SELECT 1 FROM user WHERE email = ?)", email).Scan(&exists)

	return exists, err
}

func (mysql mysqlUserRepository) Get(user *models.User) customError.Error {
	results, err := mysql.Conn.Query("SELECT * FROM user WHERE username = ?", user.Username)
	if err != nil {
		return customError.NewGenericHttpError(err)
	}

	for results.Next() {
		err = results.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Salt, &user.Avatar, &user.Bio)
		if err != nil {
			return customError.NewGenericHttpError(err)
		}
	}

	return nil
}


