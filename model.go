package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

type Movie struct {
	ID        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name      string        `json:"name" bson:"name"`
	Year      string        `json:"year" bson:"year"`
	Directors []string      `json:"directors" bson:"directors"`
	Writers   []string      `json:"writers" bson:"writers"`
	BoxOffice BoxOffice     `json:"boxOffice" bson:"boxOffice"`
}

type BoxOffice struct {
	Budget uint64 `json:"budget" bson:"budget"`
	Gross  uint64 `json:"gross" bson:"gross"`
}

//영화를 가져오는 함수
func (m *Movie) GetMovie(db *App) error {

	//var movie Movie
	err := db.collection.Find(bson.M{"_id": m.ID}).One(&m)
	if err != nil {
		return err
	}

	return nil
}

//영화 정보를 입력함
func (m *Movie) AddMovie(db *App) (string, error) {
	//var movie Movie
	//postBody, _ := ioutil.ReadAll(r.Body)
	m.ID = bson.NewObjectId()
	var name = db.collection.Name
	fmt.Println(name)
	err := db.collection.Insert(m)
	if err != nil {
		return "", err
	}

	return m.ID.Hex(), nil
}

func (m *Movie) UpdateMovie(db *App) error {
	fmt.Printf("updateMovie : %v\n", m)
	err := db.collection.Update(bson.M{"_id": bson.ObjectIdHex(m.ID.Hex())}, bson.M{"$set": &m})
	if err != nil {
		return err
	}

	return nil
}

func (m *Movie) DeleteMovie(db *App) error {
	err := db.collection.Remove(bson.M{"_id": bson.ObjectIdHex(m.ID.Hex())})
	if err != nil {
		return err
	}

	return nil
}

//func (m *Movie) GetMovie(db *App) error {
//	//vars := mux.Vars(r)
//	db.collection.Database
//}

//영화 정보를 가져옴
func GetMovies(db *App) ([]Movie, error) {
	var movies []Movie

	err := db.collection.Find(nil).All(&movies)
	if err != nil {
		return nil, err
	}

	return movies, nil
}
