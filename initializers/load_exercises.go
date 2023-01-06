package initializers

import (
	"encoding/json"
	"hypegym-backend/models"
	"io/ioutil"
	"log"
)

func LoadExercises() {
	content, err := ioutil.ReadFile("./exercises.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var exercises []models.Exercise
	err = json.Unmarshal(content, &exercises)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	DB.Create(&exercises)
}
