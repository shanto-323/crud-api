package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
)

// ERROR HANDLING
type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func sendError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Message: message, Code: code})
}

// DB
// DB model
type Member struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	Active            bool   `json:"active"`
	Subscription_type string `json:"subscription_type"`
	Join_date         string `json:"join_date"`
}

var db *sql.DB

func init() {
	var err error
	//use your own usernaem , password & dbname
	dsn := "user1:12345678@tcp(localhost:3306)/mydb"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	fmt.Println("Database connection successful!")
}

// Main function
func main() {
	router := chi.NewRouter()

	router.Get("/home/members", GetAllData)
	router.Get("/home/members/{id}", GetDataById)
	router.Post("/home/members", PostData)
	router.Put("/home/members/{id}", UpadateData)
	router.Delete("/home/members/{id}", Delete)

	fmt.Print("server is running on port 8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Print(err)
		return
	}
}

// CONTROLLERS
// GET ALL DATA
func GetAllData(w http.ResponseWriter, r *http.Request) {
	query := `
        SELECT * FROM members
  `
	rows, err := db.Query(query)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var members []Member
	for rows.Next() {
		var member Member
		err = rows.Scan(&member.Id, &member.Name, &member.Active, &member.Subscription_type, &member.Join_date)
		members = append(members, member)
	}

	if err = rows.Err(); err != nil {
		sendError(w, http.StatusBadGateway, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(members); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
}

// GET DATA BY ID
func GetDataById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	query := `
        SELECT * FROM members
        WHERE id = ?
  `
	row := db.QueryRow(query, id)

	var member Member
	err := row.Scan(&member.Id, &member.Name, &member.Active, &member.Subscription_type, &member.Join_date)

	if err = row.Err(); err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(member); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
}

// POST DATA
func PostData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var member *Member
	err := json.NewDecoder(r.Body).Decode(&member)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	query := `
      INSERT INTO members (name, active, subscription_type, join_date) 
      VALUES (?, ?, ?, ?) 
  `
	_, err = db.Exec(query, member.Name, member.Active, member.Subscription_type, member.Join_date)

	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(member); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
}

// UPDATE DATA
func UpadateData(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var member *Member
	err := json.NewDecoder(r.Body).Decode(&member)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	query := `
      UPDATE members SET
      name = ?, active = ?, subscription_type = ?, join_date = ?
      WHERE id = ?
  `
	_, err = db.Exec(query, member.Name, member.Active, member.Subscription_type, member.Join_date, id)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	json.NewEncoder(w).Encode(member)
}

// DELETE DATA
func Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	query := `
      DELETE FROM members
      WHERE id = ?
  `
	_, err := db.Exec(query, id)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Write([]byte("Id: " + id + " delted successfully"))
}
