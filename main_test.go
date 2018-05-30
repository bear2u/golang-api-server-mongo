package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cloudfunding-api-server"
	"github.com/google/go-cmp/cmp"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a main.App

/*
"192.168.99.100",
		"godb",
		"movies",
*/
var ip, dbName, collectionName = "192.168.99.100", "godb", "movies"

func TestMain(m *testing.M) {
	a = main.App{}
	a.Initialize(
		ip,
		dbName,
		collectionName,
	)

	//테이블 존재 테스팅
	//ensureTableExists()

	code := m.Run()

	//테이블 제거
	//clearTable()

	os.Exit(code)
}

func addMovie() {

}

func deleteMovie() {

}

//func TestEmptyTable(t *testing.T) {
//
//	req, _ := http.NewRequest("GET", "/movies", nil)
//	response := executeRequest(req)
//
//	checkResponseCode(t, http.StatusOK, response.Code)
//
//	if body := response.Body.String(); body != "[]" {
//		t.Errorf("Expected an empty array. Got %s", body)
//	}
//}
//
//func TestGetNotExistentProduct(t *testing.T) {
//	req, _ := http.NewRequest("GET", "/movie/11", nil)
//	response := executeRequest(req)
//
//	checkResponseCode(t, http.StatusNotFound, response.Code)
//
//	var m map[string]string
//	json.Unmarshal(response.Body.Bytes(), &m)
//	fmt.Println(m)
//	if m["error"] != "Product not found" {
//		t.Errorf("Expected the 'error' key of the response to be set to %s", m["error"])
//	}
//}

func TestAddMovie(t *testing.T) {

	darkNight := getMovieDummy()

	darkNightMarshel, err := json.Marshal(darkNight)

	if err != nil {
		panic(err)
	}

	//fmt.Println(string(darkNightMarshel))

	req, _ := http.NewRequest("POST", "/movie", bytes.NewBufferString(string(darkNightMarshel)))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	//var m map[string]interface{}
	var m main.Movie
	json.Unmarshal(response.Body.Bytes(), &m)

	fmt.Println("m: ", m)

	if darkNight.Name != m.Name {
		t.Errorf("Its not name equal %v", m.Name)
	}

	if darkNight.Year != m.Year {
		t.Errorf("Its not Year equal %v", m.Year)
	}

	if !cmp.Equal(darkNight.Directors, m.Directors) {
		t.Errorf("Its not Directors equal %v", m.Directors)
	}

	if !cmp.Equal(darkNight.Writers, m.Writers) {
		t.Errorf("Its not Writers equal %v", m.Writers)
	}

	if !cmp.Equal(darkNight.BoxOffice, m.BoxOffice) {
		t.Errorf("Its not BoxOffice equal %v", m.BoxOffice)
	}
}

func TestGetMovie(t *testing.T) {
	//id := bson.NewObjectId()
	id := addDommyMovie(getMovieDummy())
	req, _ := http.NewRequest("GET", "/movie/"+id, nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m main.Movie
	json.Unmarshal(response.Body.Bytes(), &m)
}

func TestUpdateMovie(t *testing.T) {
	//우선 값 추가
	movie := getMovieDummy()
	id := addDommyMovie(movie)

	req, _ := http.NewRequest("GET", "/movie/"+id, nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var orginalMovie main.Movie
	json.Unmarshal(response.Body.Bytes(), &orginalMovie)

	newMovie := getMovieDummy()
	newMovie.ID = bson.ObjectIdHex(id)
	newMovie.Name = "Star wars"

	fmt.Println(orginalMovie)
	fmt.Println(newMovie)

	NewMovieMarshel, err := json.Marshal(newMovie)

	if err != nil {
		panic(err)
	}

	//fmt.Println(string(darkNightMarshel))

	updateReq, _ := http.NewRequest("PUT", "/movie/"+id, bytes.NewBufferString(string(NewMovieMarshel)))
	updateResponse := executeRequest(updateReq)
	checkResponseCode(t, http.StatusOK, updateResponse.Code)

	checkReq, _ := http.NewRequest("GET", "/movie/"+id, nil)
	checkResponse := executeRequest(checkReq)
	checkResponseCode(t, http.StatusOK, checkResponse.Code)

	var checkMovie main.Movie
	json.Unmarshal(checkResponse.Body.Bytes(), &checkMovie)

	if orginalMovie.ID.Hex() != checkMovie.ID.Hex() {
		t.Errorf("Id should same %v, %v", orginalMovie.ID.Hex(), checkMovie.ID.Hex())
	}

	if orginalMovie.Name == checkMovie.Name {
		t.Errorf("Name should different %v", checkMovie.Name)
	}
}

func TestDeleteMovie(t *testing.T) {
	movie := getMovieDummy()
	id := addDommyMovie(movie)

	req, _ := http.NewRequest("GET", "/movie/"+id, nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/movie/"+id, nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/movie/"+id, nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestGetMovies(t *testing.T) {
	req, _ := http.NewRequest("GET", "/movies/", nil)
	response := executeRequest(req)

	var m []*main.Movie
	json.Unmarshal(response.Body.Bytes(), &m)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func addDommyMovie(m *main.Movie) string {
	id, _ := m.AddMovie(&a)
	return id
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func getMovieDummy() *main.Movie {
	return &main.Movie{
		Name:      "The Dark Knight",
		Year:      "2008",
		Directors: []string{"Christopher Nolan"},
		Writers:   []string{"Jonathan Nolan", "Christopher Nolan"},
		BoxOffice: main.BoxOffice{
			Budget: 185000000,
			Gross:  533316061,
		},
	}
}
