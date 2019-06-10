package repository

import (
	"database/sql"
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/models"
	"github.com/skyerus/riptides-go/pkg/user"
	"net/http"
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
	defer results.Close()
	if err != nil {
		return customError.NewGenericHttpError(err)
	}
	res := results.Next()
	if !res {
		return customError.NewHttpError(http.StatusBadRequest, "No user exists with this username", nil)
	}
	err = results.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Salt, &user.Avatar, &user.Bio)
	if err != nil {
		return customError.NewGenericHttpError(err)
	}

	return nil
}

func (mysql mysqlUserRepository) GetFollowing(user *models.User, offset int, limit int) ([]models.Following, customError.Error) {
	users := make([]models.Following, limit - offset)
	query := "SELECT user.* FROM user LEFT JOIN (SELECT following_id, follower_id FROM user_follow_user GROUP BY id ORDER BY date_created DESC) as f ON user.id = f.follower_id WHERE f.following_id = ? LIMIT ?, ?"
	results, err := mysql.Conn.Query(query, user.ID, offset, limit)
	if err != nil {
		return users, customError.NewGenericHttpError(err)
	}
	defer results.Close()

	for results.Next() {
		var u models.Following
		if err := results.Scan(&u.User.ID, &u.User.Username, &u.User.Email, &u.User.Password, &u.User.Salt, &u.User.Avatar, &u.User.Bio); err != nil {
			return users, customError.NewGenericHttpError(err)
		}
		u.Following = false
		users = append(users, u)
	}

	return users, nil
}


