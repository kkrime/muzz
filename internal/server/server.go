package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"muzz/internal/config"
	"muzz/internal/model"
	"muzz/internal/service"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/request"
	"github.com/julienschmidt/httprouter"
)

type Server interface {
	Run() error
	Login() httprouter.Handle
	CreateUser() httprouter.Handle
	Discover() httprouter.Handle
	Swipe() httprouter.Handle
}

type server struct {
	service service.Service
	router  *httprouter.Router
}

func NewServer(config *config.Config) (Server, error) {
	servive, err := service.NewService(config)
	if err != nil {
		return nil, err
	}

	router := httprouter.New()
	server := &server{
		service: servive,
		router:  router,
	}

	router.POST("/login", server.Login())
	router.GET("/user/create", server.CreateUser())
	router.GET("/discover", server.Discover())
	router.POST("/swipe", server.Swipe())

	return server, nil
}

func (s *server) Run() error {
	return http.ListenAndServe(":8080", s.router)
}

func (s *server) Login() httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {

		var result model.Result
		var login model.Login

		err := json.NewDecoder(r.Body).Decode(&login)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			result.Error = "bad input"
			json.NewEncoder(w).Encode(&result)
			return
		}
		// TODO validate input

		token, err := s.service.Login(r.Context(), &login)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			result.Error = err.Error()
			json.NewEncoder(w).Encode(&result)
			return
		}

		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			result.Error = "bad username or password"
			json.NewEncoder(w).Encode(&result)
			return
		}

		out := model.Token{
			Token: token,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&out)
	}
}

func (s *server) CreateUser() httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {

		var result model.Result

		res, err := s.service.CreateUser(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			result.Error = err.Error()
		} else {
			w.WriteHeader(http.StatusCreated)
			result.Result = res
		}

		json.NewEncoder(w).Encode(&result)
	}
}

func (s *server) Discover() httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {
		var result model.Result

		userID, err := authorize(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			result.Error = err.Error()
			json.NewEncoder(w).Encode(&result)
			return
		}

		res, err := s.service.Discover(r.Context(), userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			result.Error = err.Error()
		} else {
			w.WriteHeader(http.StatusOK)
			result.Result = res
		}

		json.NewEncoder(w).Encode(&result)
	}
}

func (s *server) Swipe() httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {
		var result model.Result

		userID, err := authorize(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			result.Error = err.Error()
			json.NewEncoder(w).Encode(&result)
			return
		}

		var swipe model.Swipe
		err = json.NewDecoder(r.Body).Decode(&swipe)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			result.Error = "bad input"
			json.NewEncoder(w).Encode(&result)
			return
		}
		// TODO validate input

		res, err := s.service.Swipe(r.Context(), userID, swipe.UserID, swipe.SwipeRight)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			result.Error = err.Error()
		} else {
			w.WriteHeader(http.StatusOK)
			result.Result = res
		}

		json.NewEncoder(w).Encode(&result)
	}
}

func authorize(r *http.Request) (int, error) {
	var claims = jwt.MapClaims{}
	token, err := request.ParseFromRequestWithClaims(r, request.AuthorizationHeaderExtractor, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	// for some reason jwt packages ints as float64
	userID, ok := claims["userID"].(float64)
	if !ok || userID == 0 {
		// log error
		return 0, fmt.Errorf("unable to get userID from claim. token:  %v", token)
	}

	return int(userID), nil
}