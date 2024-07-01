package db

import (
	"context"
	"fmt"
	"muzz/internal/config"
	"muzz/internal/model"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
)

type DB interface {
	CreateUser(
		ctc context.Context,
		firstName string,
		lastName string,
		email string,
		password []byte,
		gender string,
		dob time.Time) (int, error)

	GetUserPassword(ctx context.Context, email string) (*model.UserPassword, error)
}

type db struct {
	db *sqlx.DB
}

func NewDB(config *config.DBConfig) (DB, error) {

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBname)

	db_, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// dbLog := logger.CreateNewLogger()

	// db.DB = sqldblogger.OpenDriver(dsn, db.DB.Driver(), logrusadapter.New(dbLog),
	// 	sqldblogger.WithTimeFormat(sqldblogger.TimeFormatRFC3339),
	// 	sqldblogger.WithLogDriverErrorSkip(true),
	// 	sqldblogger.WithSQLQueryAsMessage(true))

	err = db_.Ping()
	if err != nil {
		return nil, err
	}

	return &db{
		db: db_,
	}, nil
}

func (d *db) CreateUser(
	ctx context.Context,
	firstName string,
	lastName string,
	email string,
	password []byte,
	gender string,
	dob time.Time) (int, error) {

	var id int

	statement := `
		INSERT INTO
			users
			(
				first_name,
				last_name,
				email,
				password,
				gender,
				dob
			)
		VALUES
			(
				$1,
				$2,
				$3,
				$4,
				$5,
				$6
			)
		RETURNING id
	;`

	// TODO duplicate check
	r := d.db.QueryRowxContext(ctx, statement, firstName, lastName, email, password, gender, dob)
	if err := r.Err(); err != nil {
		return 0, err
	}
	r.Scan(&id)

	return id, nil
}
func (d *db) GetUserPassword(ctx context.Context, email string) (*model.UserPassword, error) {

	var userPassword []model.UserPassword

	statement := `
		SELECT 
			id, password
		FROM
			users
		WHERE 
			email = $1
	;`

	err := d.db.SelectContext(ctx, &userPassword, statement, email)
	if err != nil {
		return nil, err
	}
	if userPassword == nil {
		return nil, nil
	}

	return &userPassword[0], nil

}
