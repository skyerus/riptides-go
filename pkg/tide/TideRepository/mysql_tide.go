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

