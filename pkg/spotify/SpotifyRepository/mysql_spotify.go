package SpotifyRepository

import (
	"database/sql"
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/models"
	"github.com/skyerus/riptides-go/pkg/spotify"
	"github.com/skyerus/riptides-go/pkg/spotify/SpotifyHandler"
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

func (mysql mysqlSpotifyRepository) GetCredentials(user *models.User) (models.SpotifyCredentials, customError.Error) {
	var creds models.SpotifyCredentials
	err := mysql.Conn.QueryRow("SELECT access_token, refresh_token FROM spotify_credentials WHERE user_id = ?", user.ID).Scan(&creds.AccessToken, &creds.RefreshToken)
	if err != nil {
		return creds, customError.NewGenericHttpError(err)
	}
	return creds, nil
}

func (mysql mysqlSpotifyRepository) SaveCredentials(creds SpotifyHandler.Credentials, user *models.User) customError.Error {
	stmtIns, err := mysql.Conn.Prepare("INSERT INTO spotify_credentials (user_id, access_token, refresh_token) VALUES(?, ?, ?)")
	if err != nil {
		return customError.NewGenericHttpError(err)
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(user.ID, creds.AccessToken, creds.RefreshToken)
	if err != nil {
		return customError.NewGenericHttpError(err)
	}

	return nil
}

func (mysql mysqlSpotifyRepository) UpdateCredentials(creds SpotifyHandler.Credentials, user *models.User) customError.Error {
	stmtIns, err := mysql.Conn.Prepare("UPDATE spotify_credentials SET access_token = ? AND refresh_token = ? WHERE user_id = ?")
	if err != nil {
		return customError.NewGenericHttpError(err)
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(creds.AccessToken, creds.RefreshToken, user.ID)
	if err != nil {
		return customError.NewGenericHttpError(err)
	}
	return nil
}
