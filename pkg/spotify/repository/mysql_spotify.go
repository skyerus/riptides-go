package repository

import (
	"database/sql"
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/models"
	"github.com/skyerus/riptides-go/pkg/spotify"
)

type mysqlSpotifyRepository struct {
	Conn *sql.DB
}

func NewMysqlSpotifyRepository(Conn *sql.DB) spotify.Repository {
	return &mysqlSpotifyRepository{Conn}
}

func (mysql mysqlSpotifyRepository) CredentialsExist(user *models.User) (bool, customError.Error) {
	var exists bool
	err := mysql.Conn.QueryRow("SELECT EXISTS(SELECT 1 FROM spotify_credentials WHERE user_id = ?)", user.ID).Scan(&exists)
	if err != nil {
		return exists, customError.NewGenericHttpError(err)
	}
	return exists, nil
}
