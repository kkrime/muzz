package db

import (
	"context"
	"database/sql"
	"fmt"
	"muzz/internal/config"
	"muzz/internal/model"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
)

// helper interface
// had to make exported so I can make mock for testing
type Tx interface {
	Commit() error
}

type DB interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	Commit(t Tx) error

	CreateUser(
		ctx context.Context,
		firstName string,
		lastName string,
		email string,
		password []byte,
		gender string,
		dob time.Time) (int, error)
	Login(tx *sql.Tx, userID int, long float64, lat float64) error
	GetUserPassword(ctx context.Context, email string) (*model.UserPassword, error)

	Discover(ctx context.Context, UserId int) ([]model.Discover, error)
	Swipe(ctx context.Context, currentUserID int, theirUserID int, swipeRight bool) error
	Match(ctx context.Context, currentUserID int, theirUserID int) (bool, error)
}

type db struct {
	db *sqlx.DB
}

func NewDB(config *config.Config) (DB, error) {

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBname)

	db_, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db_.Ping()
	if err != nil {
		return nil, err
	}

	return &db{
		db: db_,
	}, nil
}

func (d *db) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return d.db.BeginTx(ctx, nil)
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
	// user doesn't exist
	if userPassword == nil {
		return nil, nil
	}

	return &userPassword[0], nil
}

// TODO pagination
func (d *db) Discover(ctx context.Context, userID int) ([]model.Discover, error) {
	var discover []model.Discover

	statement := `
	SELECT 
		id, name, gender, age, distance 
	FROM 
	(
		SELECT
			id, first_name || ' ' || last_name AS name, gender, DATE_PART('year', AGE(dob)) AS age,
				COALESCE(NULLIF(CEIL((
					SELECT 
						ST_DISTANCE
						(
							(SELECT location FROM logins WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1),
							(SELECT location FROM logins WHERE user_id = users.id ORDER BY created_at DESC LIMIT 1)
						) / 1609.34 
				)), 0), 1) AS distance,
			(SELECT COUNT(*) FROM swipes WHERE their_user_id = users.id AND swipe_right = TRUE) AS attractiveness,
			deleted_at
		FROM
			users
	) AS results
	WHERE
		deleted_at IS NULL
		AND
		-- filter by gender
		gender =  (
								CASE (SELECT gender FROM users WHERE id = $1)
									WHEN 'M'::gender THEN 'F'::gender
									WHEN 'F'::gender THEN 'M'::gender
								END
							) 
		AND
		-- TODO make age configurable
		-- filter by age
		age >= 21 AND age <= 35
		AND 
		distance IS NOT NULL -- corner case for users who have signed up but not logged in yet
	  AND
	  -- filter profiles already swipped on
	  id NOT IN (SELECT their_user_id FROM swipes WHERE user_id = $1)
	-- order by attractiveness
	ORDER BY attractiveness DESC
	-- TODO make LIMIT configurable
	LIMIT 20
	;`

	err := d.db.SelectContext(ctx, &discover, statement, userID)
	if err != nil {
		return nil, err
	}

	return discover, nil
}

func (d *db) Login(tx *sql.Tx, userID int, long float64, lat float64) error {

	statement := `
		INSERT INTO
			logins
			(
				user_id,
				location
			)
		VALUES
			(
				$1,
				ST_SetSRID(ST_MakePoint($2, $3), 4326)::GEOGRAPHY
			)
	;`

	r, err := tx.Query(statement, userID, long, lat)
	if err != nil {
		return err

	}
	r.Close()

	return nil
}

func (d *db) Swipe(ctx context.Context, currentUserID int, theirUserID int, swipeRight bool) error {

	statement := `
		INSERT INTO
			swipes
			(
				user_id,
				their_user_id,
				swipe_right
			)
		VALUES
			(
				$1,
				$2,
				$3
			)
	;`

	r, err := d.db.Query(statement, currentUserID, theirUserID, swipeRight)
	if err != nil {
		return err
	}
	r.Close()

	return err

}

func (d *db) Match(ctx context.Context, currentUserID int, theirUserID int) (bool, error) {

	var match bool

	statement := `
		SELECT 
		( 
			COALESCE
			( 
				(SELECT TRUE FROM swipes WHERE user_id = $1 AND their_user_id = $2 AND swipe_right = TRUE LIMIT 1) 
				AND 
				(SELECT TRUE FROM swipes WHERE user_id = $2 AND their_user_id = $1 AND swipe_right = TRUE LIMIT 1)
			, 'f')
			)	AS match
		;`

	err := d.db.QueryRowContext(ctx, statement, currentUserID, theirUserID).Scan(&match)

	return match, err

}

func (d *db) Commit(t Tx) error {
	return t.Commit()
}
