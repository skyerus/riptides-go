package repository

import (
	"database/sql"
	"github.com/skyerus/riptides-go/pkg/models"
	"github.com/skyerus/riptides-go/pkg/user"
)

type mysqlUserRepository struct {
	Conn *sql.DB
}

func NewMysqlUserRepository(Conn *sql.DB) user.Repository {
	return &mysqlUserRepository{Conn}
}

func (m mysqlUserRepository) Create(user models.User) (err error) {
	stmtIns, err := m.Conn.Prepare("INSERT INTO user (username, email, password, salt, avatar, bio) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(user.Username, user.Email, user.Password, user.Salt, user.Avatar, user.Bio)
	if err != nil {
		return err
	}

	return
}

func (mysqlUserRepository) DoesUserExist(user models.User) error {
	panic("implement me")
}


