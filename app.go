package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
)

type App struct {
	Router     *mux.Router
	session    *mgo.Session
	collection *mgo.Collection
}

func (a *App) Initialize(serverIp, dbName, collectionName string) {
	session, err := mgo.Dial(serverIp)
	a.session = session
	a.collection = a.session.DB(dbName).C(collectionName)

	if err != nil {
		panic(err)
	}
	//defer session.Close()

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/movie", a.addMovie).Methods("POST")
	a.Router.HandleFunc("/movie/{id:[a-zA-Z0-9]*}", a.getMovie).Methods("GET")
	a.Router.HandleFunc("/movie/{id:[a-zA-Z0-9]*}", a.updateMovie).Methods("PUT")
	a.Router.HandleFunc("/movie/{id:[a-zA-Z0-9]*}", a.deleteMovie).Methods("DELETE")
	a.Router.HandleFunc("/movies/", a.getMovies).Methods("GET")
}

func (a *App) Run(add string) {}

func (a *App) addMovie(w http.ResponseWriter, r *http.Request) {
	var movie Movie
	postBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "parsing error")
		return
	}
	json.Unmarshal(postBody, &movie)

	defer r.Body.Close()

	if _, err := movie.AddMovie(a); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, movie)
}

func (a *App) getMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	movie := Movie{ID: bson.ObjectIdHex(vars["id"])}
	if err := movie.GetMovie(a); err != nil {
		switch err {
		case mgo.ErrNotFound:
			respondWithError(w, http.StatusNotFound, err.Error())
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}

		return
	}

	respondWithJSON(w, http.StatusOK, movie)
}

func (a *App) updateMovie(w http.ResponseWriter, r *http.Request) {
	var movie Movie
	postBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "parsing error")
		return
	}
	json.Unmarshal(postBody, &movie)

	defer r.Body.Close()

	if err := movie.UpdateMovie(a); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, movie)
}

func (a *App) deleteMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	movie := &Movie{ID: bson.ObjectIdHex(vars["id"])}

	if err := movie.DeleteMovie(a); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, movie)
}

func (a *App) getMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := GetMovies(a)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, movies)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
