package UserRepository

import (
	"database/sql"
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/models"
	"github.com/skyerus/riptides-go/pkg/user"
	"net/http"
	"time"
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
	results, err := mysql.Conn.Query("SELECT user.id, user.username, user.password, user.email, user.avatar, user.bio FROM user WHERE username = ?", user.Username)
	if err != nil {
		return customError.NewGenericHttpError(err)
	}
	defer results.Close()
	res := results.Next()
	if !res {
		return customError.NewHttpError(http.StatusBadRequest, "No user exists with this username", nil)
	}
	err = results.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Avatar, &user.Bio)
	if err != nil {
		return customError.NewGenericHttpError(err)
	}

	return nil
}

func (mysql mysqlUserRepository) GetFromId(id int) (models.User, customError.Error) {
	var User models.User
	results, err := mysql.Conn.Query("SELECT user.id, user.username, user.password, user.email, user.avatar, user.bio FROM user WHERE id = ?", id)
	if err != nil {
		return User, customError.NewGenericHttpError(err)
	}
	defer results.Close()
	res := results.Next()
	if !res {
		return User, customError.NewHttpError(http.StatusBadRequest, "No user exists with this id", nil)
	}
	err = results.Scan(&User.ID, &User.Username, &User.Password, &User.Email, &User.Avatar, &User.Bio)
	if err != nil {
		return User, customError.NewGenericHttpError(err)
	}

	return User, nil
}

func (mysql mysqlUserRepository) GetFollowing(user *models.User, offset int, limit int) ([]models.Following, customError.Error) {
	var users []models.Following
	query := "SELECT user.id, user.username, user.email, user.avatar, user.bio FROM user LEFT JOIN (SELECT following_id, follower_id FROM user_follow_user GROUP BY id ORDER BY date_created DESC) as f ON user.id = f.follower_id WHERE f.following_id = ? LIMIT ?, ?"
	results, err := mysql.Conn.Query(query, user.ID, offset, limit)
	if err != nil {
		return users, customError.NewGenericHttpError(err)
	}
	defer results.Close()

	for results.Next() {
		var u models.Following
		if err := results.Scan(&u.User.ID, &u.User.Username, &u.User.Email, &u.User.Avatar, &u.User.Bio); err != nil {
			return users, customError.NewGenericHttpError(err)
		}
		u.Following = true
		users = append(users, u)
	}

	return users, nil
}

func (mysql mysqlUserRepository) DoesUserFollow(currentUser *models.User, user *models.User) (bool, customError.Error) {
	var exists bool
	err := mysql.Conn.QueryRow("SELECT EXISTS(SELECT 1 FROM user_follow_user WHERE following_id = ? AND follower_id = ?)", currentUser.ID, user.ID).Scan(&exists)
	if err != nil {
		return exists, customError.NewGenericHttpError(err)
	}
	return exists, nil
}

func (mysql mysqlUserRepository) GetFollowers(user *models.User, offset int, limit int) ([]models.Following, customError.Error) {
	var users []models.Following
	query := "SELECT user.id, user.username, user.email, user.avatar, user.bio FROM user LEFT JOIN (SELECT following_id, follower_id FROM user_follow_user GROUP BY id ORDER BY date_created DESC) as f ON user.id = f.following_id WHERE f.follower_id = ? LIMIT ?, ?"
	results, err := mysql.Conn.Query(query, user.ID, offset, limit)
	if err != nil {
		return users, customError.NewGenericHttpError(err)
	}
	defer results.Close()

	for results.Next() {
		var u models.Following
		if err := results.Scan(&u.User.ID, &u.User.Username, &u.User.Email, &u.User.Avatar, &u.User.Bio); err != nil {
			return users, customError.NewGenericHttpError(err)
		}
		u.Following = true
		users = append(users, u)
	}

	return users, nil
}

func (mysql mysqlUserRepository) GetFollowingCount(user *models.User) (int, customError.Error) {
	var number int
	results, err := mysql.Conn.Query("SELECT COUNT(u.id) FROM riptides.user_follow_user as u WHERE u.following_id = ?", user.ID)
	if err != nil {
		return number, customError.NewGenericHttpError(err)
	}
	defer results.Close()
	res := results.Next()
	if !res {
		return number, customError.NewGenericHttpError(nil)
	}
	err = results.Scan(&number)
	if err != nil {
		return number, customError.NewGenericHttpError(nil)
	}

	return number, nil
}

func (mysql mysqlUserRepository) GetFollowerCount(user *models.User) (int, customError.Error) {
	var number int
	results, err := mysql.Conn.Query("SELECT COUNT(u.id) FROM riptides.user_follow_user as u WHERE u.follower_id = ?", user.ID)
	if err != nil {
		return number, customError.NewGenericHttpError(err)
	}
	defer results.Close()
	res := results.Next()
	if !res {
		return number, customError.NewGenericHttpError(nil)
	}
	err = results.Scan(&number)
	if err != nil {
		return number, customError.NewGenericHttpError(nil)
	}

	return number, nil
}

func (mysql mysqlUserRepository) Follow(currentUser *models.User, user *models.User) customError.Error {
	stmtIns, err := mysql.Conn.Prepare("INSERT INTO user_follow_user (following_id, follower_id, interest, date_created) VALUES(?, ?, ?, ?)")
	if err != nil {
		return customError.NewGenericHttpError(err)
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(currentUser.ID, user.ID, 0, time.Now())
	if err != nil {
		return customError.NewGenericHttpError(err)
	}

	return nil
}

func (mysql mysqlUserRepository) Unfollow(currentUser *models.User, user *models.User) customError.Error {
	stmtIns, err := mysql.Conn.Prepare("DELETE FROM user_follow_user WHERE following_id = ? AND follower_id = ?")
	if err != nil {
		return customError.NewGenericHttpError(err)
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(currentUser.ID, user.ID)
	if err != nil {
		return customError.NewGenericHttpError(err)
	}

	return nil
}
