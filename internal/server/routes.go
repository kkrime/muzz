package server

import (
	"encoding/json"
	"muzz/internal/config"
	"muzz/internal/model"
	"muzz/internal/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Server interface {
	Run() error
	Login() httprouter.Handle
	CreateUser() httprouter.Handle
	Discover() httprouter.Handle
}

type server struct {
	service service.Service
	router  *httprouter.Router
}

func NewServer(config *config.DBConfig) (Server, error) {
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

		success, token, err := s.service.Login(r.Context(), login.Email, login.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			result.Error = err.Error()
			json.NewEncoder(w).Encode(&result)
			return
		}

		if !success {
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

		res, err := s.service.Discover(r.Context(), 1)
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
