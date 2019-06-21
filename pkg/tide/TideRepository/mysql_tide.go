package TideRepository

import (
	"database/sql"
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/models"
	"github.com/skyerus/riptides-go/pkg/tide"
	"time"
)

type mysqlTideRepository struct {
	Conn *sql.DB
}

func NewMysqlTideRepository(conn *sql.DB) tide.Repository {
	return &mysqlTideRepository{conn}
}

func (mysql mysqlTideRepository) CreateTide(user *models.User, tide *models.Tide) customError.Error {
	stmtIns, err := mysql.Conn.Prepare("INSERT INTO tide (user_id, name, date_created, about) VALUES(?, ?, ?, ?)")
	if err != nil {
		return customError.NewGenericHttpError(err)
	}
	defer stmtIns.Close()

	res, err := stmtIns.Exec(user.ID, tide.Name, time.Now(), tide.About)
	if err != nil {
		return customError.NewGenericHttpError(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return customError.NewGenericHttpError(err)
	}
	tide.ID = int(id)

	return nil
}

func (mysql mysqlTideRepository) CreateTideGenre(tide *models.Tide, genre *models.Genre) customError.Error {
	stmtIns, err := mysql.Conn.Prepare("INSERT INTO tide_genre (tide_id, genre_id) VALUES(?, ?)")
	if err != nil {
		return customError.NewGenericHttpError(err)
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(tide.ID, genre.ID)
	if err != nil {
		return customError.NewGenericHttpError(err)
	}

	return nil
}

func (mysql mysqlTideRepository) GetTag(tag *models.Tag) bool {
	err := mysql.Conn.QueryRow("SELECT id FROM tag WHERE name = ?", tag.Name).Scan(&tag.ID)
	if err != nil {
		return false
	}
	return true
}

func (mysql mysqlTideRepository) CreateTag(tag *models.Tag) customError.Error {
	stmtIns, err := mysql.Conn.Prepare("INSERT INTO tag (name) VALUES(?)")
	if err != nil {
		return customError.NewGenericHttpError(err)
	}
	defer stmtIns.Close()

	res, err := stmtIns.Exec(tag.Name)
	if err != nil {
		return customError.NewGenericHttpError(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return customError.NewGenericHttpError(err)
	}
	tag.ID = int(id)

	return nil
}

func (mysql mysqlTideRepository) CreateTideTag(tide *models.Tide, tag *models.Tag) customError.Error {
	stmtIns, err := mysql.Conn.Prepare("INSERT INTO tide_tag (tide_id, tag_id) VALUES(?, ?)")
	if err != nil {
		return customError.NewGenericHttpError(err)
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(tide.ID, tag.ID)
	if err != nil {
		return customError.NewGenericHttpError(err)
	}

	return nil
}

func (mysql mysqlTideRepository) GetGenres() ([]models.Genre, customError.Error) {
	var genres []models.Genre
	results, err := mysql.Conn.Query("SELECT * FROM genre")
	if err != nil {
		return genres, customError.NewGenericHttpError(err)
	}
	defer results.Close()

	for results.Next() {
		var genre models.Genre
		err = results.Scan(&genre.ID, &genre.Name)
		if err != nil {
			return genres, customError.NewGenericHttpError(err)
		}
		genres = append(genres, genre)
	}

	return genres, nil
}

func (mysql mysqlTideRepository) GetTides(orderBy string, offset int, limit int) ([]models.Tide, customError.Error) {
	var tides []models.Tide
	results, err := mysql.Conn.Query("SELECT tide.*, u.username, u.bio, u.avatar, a.count FROM tide " +
		"LEFT JOIN ( SELECT tide_id, COUNT(*) as count FROM tide_participant as tp GROUP BY tide_id ORDER BY 2 DESC) AS a ON tide.id = a.tide_id " +
		"LEFT JOIN user as u ON u.id = tide.user_id" +
		"GROUP BY tide.id ORDER BY a.count DESC, tide." + orderBy + " DESC LIMIT " + string(offset) + ", " + string(limit))
	if err != nil {
		return tides, customError.NewGenericHttpError(err)
	}
	defer results.Close()

	for results.Next() {
		var Tide models.Tide
		var user models.User
		err = results.Scan(&Tide.ID, &Tide.Name, &Tide.DateCreated, &Tide.About, &user.Username, &user.Bio, &user.Avatar, &Tide.ParticipantCount)
		if err != nil {
			return tides, customError.NewGenericHttpError(err)
		}
		Tide.User = user

		var customErr customError.Error
		Tide.Participants, customErr = mysql.GetTideParticipants(&Tide, 20, 0)
		if customErr != nil {
			return tides, customErr
		}

		Tide.Genres, customErr = mysql.GetTideGenres(&Tide, 20, 0)
		if customErr != nil {
			return tides, customErr
		}

		Tide.Tags, customErr = mysql.GetTideTags(&Tide, 20, 0)
		if customErr != nil {
			return tides, customErr
		}

		tides = append(tides, Tide)
	}

	return tides, nil
}

func (mysql mysqlTideRepository) GetTideParticipants(tide *models.Tide, limit int, offset int) ([]models.User, customError.Error) {
	var users []models.User
	results, err := mysql.Conn.Query("SELECT u.username, u.bio, u.avatar FROM user as u, tide as t " +
		"INNER JOIN tide_participant AS tp ON t.id = tp.tide_id WHERE t.id = " + string(tide.ID) + " AND u.id = tp.user_id " +
		"LIMIT " + string(offset) + ", " + string(limit))
	if err != nil {
		return users, customError.NewGenericHttpError(err)
	}
	defer results.Close()

	for results.Next() {
		var user models.User
		err = results.Scan(&user.Username, &user.Bio, &user.Avatar)
		if err != nil {
			return users, customError.NewGenericHttpError(err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (mysql mysqlTideRepository) GetTideGenres(tide *models.Tide, limit int, offset int) ([]models.Genre, customError.Error) {
	var genres []models.Genre
	results, err := mysql.Conn.Query("SELECT g.name FROM genre as g, tide as t " +
		"INNER JOIN tide_genre as tg on t.id = tg.tide_id WHERE t.id = " + string(tide.ID) + " AND g.id = tg.genre_id " +
		"LIMIT " + string(offset) + ", " + string(limit))
	if err != nil {
		return genres, customError.NewGenericHttpError(err)
	}
	defer results.Close()

	for results.Next() {
		var genre models.Genre
		err = results.Scan(&genre.Name)
		if err != nil {
			return genres, customError.NewGenericHttpError(err)
		}

		genres = append(genres, genre)
	}

	return genres, nil
}

func (mysql mysqlTideRepository) GetTideTags(tide *models.Tide, limit int, offset int) ([]models.Tag, customError.Error) {
	var tags []models.Tag
	results, err := mysql.Conn.Query("SELECT tag.name FROM tag, tide AS t " +
		"INNER JOIN tide_tag as tt ON t.id = tt.tide_id WHERE t.id = " + string(tide.ID) + " AND tag.id = tt.tag_id " +
		"LIMIT " + string(offset) + ", " + string(limit))
	if err != nil {
		return tags, customError.NewGenericHttpError(err)
	}
	defer results.Close()

	for results.Next() {
		var tag models.Tag
		err = results.Scan(&tag.Name)
		if err != nil {
			return tags, customError.NewGenericHttpError(err)
		}

		tags = append(tags, tag)
	}

	return tags, nil
}
