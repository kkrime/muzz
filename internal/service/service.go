package service

import (
	"context"
	"database/sql"
	"fmt"
	"muzz/internal/config"
	"muzz/internal/db"
	"muzz/internal/model"
	"muzz/internal/security"
	"regexp"
	"time"

	"github.com/0x6flab/namegenerator"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/exp/rand"
)

var nameRegex *regexp.Regexp
var minAge = 21
var maxAge = 35

func init() {
	nameRegex = regexp.MustCompile(`^(\w+)-(\w+)$`)
}

type Service interface {
	CreateUser(ctx context.Context) (*model.CreatedUser, error)
	Login(ctx context.Context, login *model.Login) (string, error)
	Discover(ctx context.Context, userId int) ([]model.Discover, error)
	Swipe(ctx context.Context, currentUserID int, theirUserID int, swipeRight bool) (*model.Match, error)
}

type service struct {
	db       db.DB
	security security.Security
}

func NewService(config *config.Config) (Service, error) {
	db, err := db.NewDB(config)
	if err != nil {
		return nil, err
	}

	security := security.NewSecurity(config.TokenKey)

	return &service{
		db:       db,
		security: security,
	}, nil
}

func (s *service) Login(ctx context.Context, login *model.Login) (string, error) {

	ctx.Deadline()
	userPassword, err := s.db.GetUserPassword(ctx, login.Email)
	if err != nil {
		return "", err
	}

	// user not found
	if userPassword == nil {
		return "", nil
	}
	userID := userPassword.ID

	// NOTE; I'm writing to the login table using a go routine to "improve performance"
	// realistically, this isn't going to make a world of difference, as the commit still needs
	// to go over the network. However, this does give me an opportunity to use goroutin/channel
	// and demonstrate my understanding of go and performance.
	chann := make(chan error, 1)
	var tx *sql.Tx
	txCtx, cancelTxCtx := context.WithCancel(ctx)
	go func() {
		tx, err = s.db.BeginTx(txCtx)
		if err != nil {
			chann <- err
			return
		}
		chann <- s.db.Login(tx, userID, login.Long, login.Lat)
	}()

	defer func() {
		if err != nil {
			// will have the same effect as tx.Rollback()
			cancelTxCtx()
		}
	}()

	err = s.security.ComparePassword(userPassword.Password, login.Password)
	if err != nil {
		return "", err
	}

	token, err := s.security.CreateToken(map[string]any{
		"userID": userID,
	})
	if err != nil {
		return "", err
	}

	err = <-chann
	if err != nil {
		return "", err
	}

	err = s.db.Commit(tx)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) CreateUser(ctx context.Context) (*model.CreatedUser, error) {

	// gender
	gender := namegenerator.Gender(time.Now().Unix() & int64(1))

	generator := namegenerator.NewGenerator()
	namegenerator.WithGender(gender)

	// name
	name := generator.Generate()

	regexRes := nameRegex.FindAllStringSubmatch(name, -1)
	if regexRes == nil {
		return nil, fmt.Errorf("error generating name")
	}

	firstName := regexRes[0][1]
	lastName := regexRes[0][2]

	// email
	email := firstName + lastName + "@muzz.com"

	// dob
	age := rand.Intn(maxAge-minAge) + minAge
	dob := time.Now().AddDate(-age, 0, 0)

	// password
	password, err := password.Generate(8, 1, 1, false, false)
	if err != nil {
		return nil, err
	}
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	hashedPassword, err := s.security.GeneratePassword(password)
	if err != nil {
		return nil, err
	}

	// write to db
	var gender_ string
	if gender == namegenerator.Male {
		gender_ = "M"
	} else {
		gender_ = "F"
	}
	id, err := s.db.CreateUser(ctx, firstName, lastName, email, hashedPassword, gender_, dob)
	if err != nil {
		return nil, err
	}

	return &model.CreatedUser{
		Id:       id,
		Email:    email,
		Password: password,
		Name:     firstName + " " + lastName,
		Gender:   gender_,
		Age:      age,
	}, nil
}

func (s *service) Discover(ctx context.Context, userID int) ([]model.Discover, error) {
	return s.db.Discover(ctx, userID)
}

func (s *service) Swipe(ctx context.Context, currentUserID int, theirUserID int, swipeRight bool) (*model.Match, error) {
	var out model.Match

	// sql transactions are atomic, so no need to do any locking to avoid race conditions
	// where two usrs swipe right for each other at the same tiem
	err := s.db.Swipe(ctx, currentUserID, theirUserID, swipeRight)
	if err != nil {
		return nil, err
	}

	match, err := s.db.Match(ctx, currentUserID, theirUserID)
	if err != nil {
		return nil, err
	}

	out.Matched = match
	if match {
		out.MatchID = theirUserID
	}

	return &out, nil
}
