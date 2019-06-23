package TideRepository

import (
	"database/sql"
	"github.com/skyerus/riptides-go/pkg/customError"
	"github.com/skyerus/riptides-go/pkg/models"
	"github.com/skyerus/riptides-go/pkg/tide"
	"net/http"
	"strconv"
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
		"LEFT JOIN (SELECT tide_id, COUNT(*) as count FROM tide_participant as tp GROUP BY tide_id ORDER BY 2 DESC) AS a ON tide.id = a.tide_id " +
		"LEFT JOIN user as u ON u.id = tide.user_id " +
		"GROUP BY tide.id ORDER BY a.count DESC, tide." + orderBy + " DESC LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(limit))
	if err != nil {
		return tides, customError.NewGenericHttpError(err)
	}
	defer results.Close()

	for results.Next() {
		var Tide models.Tide
		var user models.User
		err = results.Scan(&Tide.ID, &user.ID, &Tide.Name, &Tide.DateCreated, &Tide.About, &user.Username, &user.Bio, &user.Avatar, &Tide.ParticipantCount)
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
		"INNER JOIN tide_participant AS tp ON t.id = tp.tide_id WHERE t.id = " + strconv.Itoa(tide.ID) + " AND u.id = tp.user_id " +
		"LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(limit))
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
		"INNER JOIN tide_genre as tg on t.id = tg.tide_id WHERE t.id = " + strconv.Itoa(tide.ID) + " AND g.id = tg.genre_id " +
		"LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(limit))
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
		"INNER JOIN tide_tag as tt ON t.id = tt.tide_id WHERE t.id = " + strconv.Itoa(tide.ID) + " AND tag.id = tt.tag_id " +
		"LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(limit))
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

func (mysql mysqlTideRepository) FavoriteTide(tide *models.Tide, user *models.User) customError.Error {
	stmtIns, err := mysql.Conn.Prepare("INSERT INTO user_favorite_tide (user_id, tide_id, date_created) VALUES(?, ?, ?)")
	if err != nil {
		return customError.NewGenericHttpError(err)
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(user.ID, tide.ID, time.Now())
	if err != nil {
		return customError.NewGenericHttpError(err)
	}

	return nil
}

func (mysql mysqlTideRepository) UnfavoriteTide(tide *models.Tide, user *models.User) customError.Error {
	stmtIns, err := mysql.Conn.Prepare("DELETE FROM user_favorite_tide WHERE user_id = ? AND tide_id = ?")
	if err != nil {
		return customError.NewGenericHttpError(err)
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(user.ID, tide.ID)
	if err != nil {
		return customError.NewGenericHttpError(err)
	}

	return nil
}

func (mysql mysqlTideRepository) IsTideFavorited(tide *models.Tide, user *models.User) (bool, customError.Error) {
	var exists bool
	err := mysql.Conn.QueryRow("SELECT EXISTS(SELECT 1 FROM user_favorite_tide WHERE user_id = ? AND tide_id = ?)", user.ID, tide.ID).Scan(&exists)
	if err != nil {
		return exists, customError.NewGenericHttpError(err)
	}
	return exists, nil
}

func (mysql mysqlTideRepository) GetTide(id int) (models.Tide, customError.Error) {
	var Tide models.Tide
	results, err := mysql.Conn.Query("SELECT * FROM tide WHERE id = ?", id)
	if err != nil {
		return Tide, customError.NewGenericHttpError(err)
	}
	defer results.Close()

	res := results.Next()
	if !res {
		return Tide, customError.NewHttpError(http.StatusBadRequest, "No tide exists", nil)
	}
	err = results.Scan(&Tide.ID, &Tide.User.ID, &Tide.Name, &Tide.DateCreated, &Tide.About)
	if err != nil {
		return Tide, customError.NewGenericHttpError(err)
	}

	return Tide, nil
}

func (mysql mysqlTideRepository) GetFavoriteTides(user *models.User, offset int, limit int) ([]models.Tide, customError.Error) {
	var tides []models.Tide
	results, err := mysql.Conn.Query(`SELECT tide.*, u.username, u.bio, u.avatar, a.count, uftt.tide_id as uft FROM tide 
            LEFT JOIN 
                (
                    SELECT tide_id, date_created, user_id FROM user_favorite_tide ORDER BY date_created DESC
                ) AS uft
            ON tide.id = uft.tide_id AND uft.user_id = ` + strconv.Itoa(user.ID) + ` 
            LEFT JOIN 
            (
                SELECT tide_id, user_id FROM user_favorite_tide ORDER BY date_created DESC
            ) AS uftt
            ON tide.id = uftt.tide_id AND uftt.user_id = ` + strconv.Itoa(user.ID) + ` 
            LEFT JOIN 
                (
                    SELECT tide_id, COUNT(*) as count FROM tide_participant as tp GROUP BY tide_id ORDER BY 2 DESC
                ) AS a
            ON tide.id = a.tide_id
            LEFT JOIN user as u
            ON u.id = uft.user_id
            WHERE u.id = ` + strconv.Itoa(user.ID) + ` GROUP BY tide.id, uft.date_created, uftt.tide_id ORDER BY uft.date_created DESC 
            LIMIT ` + strconv.Itoa(offset) + `, ` + strconv.Itoa(limit))
	if err != nil {
		return tides, customError.NewGenericHttpError(err)
	}
	defer results.Close()

	for results.Next() {
		var Tide models.Tide
		var favoritedInt models.NullInt64
		err = results.Scan(&Tide.ID, &Tide.User.ID, &Tide.Name, &Tide.DateCreated, &Tide.About, &Tide.User.Username, &Tide.User.Bio, &Tide.User.Avatar, &Tide.ParticipantCount, &favoritedInt)
		if err != nil {
			return tides, customError.NewGenericHttpError(err)
		}

		if favoritedInt.Valid {
			Tide.Favorited = true
		}

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

func (mysql mysqlTideRepository) GetFavoriteTidesCount(user *models.User) (int, customError.Error) {
	var number int
	results, err := mysql.Conn.Query("SELECT COUNT(t.id) FROM user_favorite_tide  as t WHERE t.user_id = ?", user.ID)
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

func (mysql mysqlTideRepository) GetUserTides(user *models.User, offset int, limit int) ([]models.Tide, customError.Error) {
	var tides []models.Tide
	results, err := mysql.Conn.Query(`SELECT tide.*, a.count, uft.tide_id as uft FROM tide 
            LEFT JOIN 
                (
                    SELECT tide_id, COUNT(*) as count FROM tide_participant as tp GROUP BY tide_id ORDER BY 2 DESC
                ) AS a
            ON tide.id = a.tide_id AND tide.user_id = ` + strconv.Itoa(user.ID) + ` 
            LEFT JOIN 
            (
                SELECT tide_id, user_id FROM user_favorite_tide ORDER BY date_created DESC 
	 		) AS uft 
            ON uft.tide_id = tide.id AND uft.user_id = ` + strconv.Itoa(user.ID) + `
            WHERE tide.user_id = ` + strconv.Itoa(user.ID) + ` ORDER BY tide.date_created DESC 
            LIMIT ` + strconv.Itoa(offset) + `, ` + strconv.Itoa(limit))
	if err != nil {
		return tides, customError.NewGenericHttpError(err)
	}
	defer results.Close()

	for results.Next() {
		var Tide models.Tide
		var favoritedInt models.NullInt64
		err = results.Scan(&Tide.ID, &Tide.User.ID, &Tide.Name, &Tide.DateCreated, &Tide.About, &Tide.ParticipantCount, &favoritedInt)
		if err != nil {
			return tides, customError.NewGenericHttpError(err)
		}

		if favoritedInt.Valid {
			Tide.Favorited = true
		}

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

		Tide.User = *user
		tides = append(tides, Tide)
	}

	return tides, nil
}
