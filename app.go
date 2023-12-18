package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type App struct {
	DB     *gorm.DB
	Router *mux.Router
}

func SendResponse(w http.ResponseWriter, statuscode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statuscode)
	json.NewEncoder(w).Encode(payload)
}
func SendError(w http.ResponseWriter, statuscode int, err string) {
	err_msg := map[string]string{"error": err}
	SendResponse(w, statuscode, err_msg)
}

func (app *App) getproducts(w http.ResponseWriter, r *http.Request) {
	products := GetAllBooks(app.DB)
	SendResponse(w, http.StatusOK, products)
}
func (app *App) getproduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	book := GetBookById(key, app.DB)

	SendResponse(w, http.StatusOK, book)
}

func (app *App) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p Book
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		SendError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}
	book, err := p.CreateBook(app.DB)
	if err != nil {
		SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	SendResponse(w, http.StatusCreated, book)

}
func (app *App) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		SendError(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	var p Book
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		SendError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	book := p.UpdateBook(key, app.DB)

	SendResponse(w, http.StatusOK, book)

}

func (app *App) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		SendError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	book := DeleteBook(key, app.DB)

	SendResponse(w, http.StatusOK, book)

}

func (app *App) Initialize() error {
	dsn := "root:S@hilkumar9873@tcp(127.0.0.1:3306)/BookStore?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	app.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	//app.DB.AutoMigrate(&Book{})
	app.Router = mux.NewRouter().StrictSlash(true)
	app.HandleRoute()
	return nil

}
func (app *App) Run(adress string) {
	log.Fatal(http.ListenAndServe(adress, app.Router))

}

func (app *App) HandleRoute() {
	app.Router.HandleFunc("/products", app.getproducts).Methods("GET")
	app.Router.HandleFunc("/products/{id}", app.getproduct).Methods("GET")
	app.Router.HandleFunc("/product", app.CreateProduct).Methods("Post")
	app.Router.HandleFunc("/product/{id}", app.UpdateProduct).Methods("PUT")
	app.Router.HandleFunc("/product/{id}", app.DeleteProduct).Methods("Delete")

}
