package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ApiServer struct {
	listenAddr string
	storage    Storage
}

type ApiError struct {
	Error error `json:"error"`
}

func MakeApi(addr string, s Storage) *ApiServer {
	return &ApiServer{
		listenAddr: addr,
		storage:    s,
	}
}

func (apiServer *ApiServer) Run() {
	router := chi.NewRouter()

	router.HandleFunc("/members", createHandlerFunc(apiServer.handleAccount))
	router.HandleFunc("/members/{id}", createHandlerFunc(apiServer.handleAccountById))

	fmt.Println("running api on port ", apiServer.listenAddr)
	err := http.ListenAndServe(apiServer.listenAddr, router)
	if err != nil {
		fmt.Println("error running api ", err)
	}
}

func (apiServer *ApiServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return apiServer.GetAllDataApi(w, r)
	}
	if r.Method == "POST" {
		return apiServer.CreateAccountApi(w, r)
	}
	return writeJson(w, http.StatusBadRequest, ApiError{Error: fmt.Errorf("method not found")})
}

func (apiServer *ApiServer) handleAccountById(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return apiServer.GetDataById(w, r)
	}
	if r.Method == "PUT" {
		return apiServer.UpdateDataByIdApi(w, r)
	}
	if r.Method == "DELETE" {
		return apiServer.DeleteDataByIdApi(w, r)
	}
	return writeJson(w, http.StatusBadRequest, ApiError{Error: fmt.Errorf("method not found")})
}

func createHandlerFunc(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			writeJson(w, http.StatusBadGateway, ApiError{Error: err})
		}
	}
}

func writeJson(w http.ResponseWriter, code int, massage any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(massage)
}

func (apiServer *ApiServer) CreateAccountApi(w http.ResponseWriter, r *http.Request) error {
	newMenberReq := &Member{}

	if err := json.NewDecoder(r.Body).Decode(newMenberReq); err != nil {
		return err
	}

	newMember := newMember(*newMenberReq)
	if err := apiServer.storage.CreateAccount(newMember); err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, newMember)
}

func (apiServer *ApiServer) GetAllDataApi(w http.ResponseWriter, r *http.Request) error {
	response, err := apiServer.storage.GetAllData()
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, response)
}

func (apiServer *ApiServer) GetDataById(w http.ResponseWriter, r *http.Request) error {
	id, err := GetId(r)
	if err != nil {
		writeJson(w, http.StatusMethodNotAllowed, err)
	}

	member := &Member{}
	member, err = apiServer.storage.GetAccountById(id)
	if err != nil {
		writeJson(w, http.StatusMethodNotAllowed, err)
	}
	return writeJson(w, http.StatusMethodNotAllowed, member)
}

func (apiServer *ApiServer) DeleteDataByIdApi(w http.ResponseWriter, r *http.Request) error {
	id, err := GetId(r)
	if err != nil {
		writeJson(w, http.StatusMethodNotAllowed, err)
	}

	err = apiServer.storage.DeleteAccount(id)
	if err != nil {
		writeJson(w, http.StatusMethodNotAllowed, err)
	}
	return writeJson(w, http.StatusMethodNotAllowed, "account delete")
}

func (apiServer *ApiServer) UpdateDataByIdApi(w http.ResponseWriter, r *http.Request) error {
	id, err := GetId(r)
	if err != nil {
		writeJson(w, http.StatusMethodNotAllowed, err)
	}

	member := &Member{}
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		writeJson(w, http.StatusMethodNotAllowed, err)
	}
	newMember, err := apiServer.storage.UpdateAccount(member, id)
	if err != nil {
		writeJson(w, http.StatusMethodNotAllowed, err)
	}
	return writeJson(w, http.StatusOK, newMember)
}

func GetId(r *http.Request) (int, error) {
	idReq := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idReq)
	if err != nil {
		return 0, err
	}
	return id, nil
}
