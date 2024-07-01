package service

import (
	"context"
	"fmt"
	"muzz/internal/client"
	"muzz/internal/config"
	"muzz/internal/db"
	"muzz/internal/model"
	"regexp"
	"time"

	"github.com/0x6flab/namegenerator"
	"github.com/golang-jwt/jwt"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
)

var nameRegex *regexp.Regexp
var minAge = 18
var maxAge = 70

func init() {
	nameRegex = regexp.MustCompile(`^(\w+)-(\w+)$`)
}

type Service interface {
	CreateUser(ctx context.Context) (*model.CreatedUser, error)
	Login(ctx context.Context, email string, password string) (bool, string, error)
}

type service struct {
	db         db.DB
	httpClient client.HttpClient
}

func NewService(config *config.DBConfig) (Service, error) {
	db, err := db.NewDB(config)
	if err != nil {
		return nil, err
	}

	httpClient := client.NewHttpClient()

	return &service{
		db:         db,
		httpClient: httpClient,
	}, nil
}

func (s *service) Login(ctx context.Context, email string, password string) (bool, string, error) {

	userPassword, err := s.db.GetUserPassword(ctx, email)
	if err != nil {
		return false, "", err
	}

	// user not found
	if userPassword == nil {
		return false, "", nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(userPassword.Password), []byte(password))
	if err != nil {
		return false, "", nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userPassword.ID,
		"nbf":    time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// TODO add secret to .env
	tokenString, err := token.SignedString([]byte("secret"))

	fmt.Println(tokenString, err)

	return true, tokenString, nil
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)
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
