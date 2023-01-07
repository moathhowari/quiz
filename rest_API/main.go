package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux" //for handling request routers
)

// defines a struct that stores question details
type Question struct {
	Text         string    `json:"Text"`
	Options      [3]string `json:"options"`
	SubmittedANS int       `json:"ANS"`
	CorrectANS   int       `json:"-"`
}

// defines a struct that stores exam details.
type Exam struct {
	UserID    int        `json:"ID"`
	Questions []Question `json:"question"`
}

// defines a struct that stores submitted answers, to be handled.
type Answers struct {
	UserID    int   `json:"id"`
	PickedANS []int `json:"ANS"`
}

// making array of exams for all users quiz[userID], it contains array struct of Question..
var quiz []Exam

// since we have no database in the task i used this function to add some questions as JSON
func initJSONExam() {

	// for 10 users!
	for i := 0; i < 10; i++ {
		var dataQuestions = []Question{
			Question{Text: "3*4", Options: [3]string{"10", "12", "15"}, CorrectANS: 2, SubmittedANS: 0},
			Question{Text: "6*6", Options: [3]string{"10", "36", "15"}, CorrectANS: 2, SubmittedANS: 0},
			Question{Text: "5*5", Options: [3]string{"10", "12", "25"}, CorrectANS: 3, SubmittedANS: 0},
			Question{Text: "4*4", Options: [3]string{"16", "12", "15"}, CorrectANS: 1, SubmittedANS: 0},
			Question{Text: "3*3", Options: [3]string{"10", "12", "9"}, CorrectANS: 3, SubmittedANS: 0},
			Question{Text: "2*2", Options: [3]string{"10", "4", "15"}, CorrectANS: 2, SubmittedANS: 0},
			Question{Text: "1*1", Options: [3]string{"10", "1", "15"}, CorrectANS: 2, SubmittedANS: 0},
			Question{Text: "3*0", Options: [3]string{"10", "12", "0"}, CorrectANS: 3, SubmittedANS: 0},
			Question{Text: "9*4", Options: [3]string{"10", "12", "36"}, CorrectANS: 3, SubmittedANS: 0},
			Question{Text: "3*9", Options: [3]string{"10", "27", "15"}, CorrectANS: 2, SubmittedANS: 0},
		}
		//quiz[i] = Exam{UserID: i, Questions: dataQuestions}
		quiz = append(quiz, Exam{UserID: i, Questions: dataQuestions})
	}
}

// returns JSON data of all users including all questions details
func getQuizData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quiz)
}

// returns JSON data of all questions for the given user by passing id of the user
func getQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	//converts id to int for array usage..
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		panic(err)
	}
	json.NewEncoder(w).Encode(quiz[id])
}

// adds answer into JSON Array ("insert" alternative in database)
func submitAnswer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		panic(err)
	}
	QNumber, err := strconv.Atoi(r.FormValue("q"))
	if err != nil {
		panic(err)
	}
	submittedAns, err := strconv.Atoi(r.FormValue("ans"))
	if err != nil {
		panic(err)
	}
	println(userID)
	log.Println(QNumber)
	log.Println(submittedAns)
	quiz[userID].Questions[QNumber].SubmittedANS = submittedAns
}

// adds all answers into JSON Array ("insert" alternative in database)
func submitAnswers(w http.ResponseWriter, r *http.Request) {
	/*{
		"id":"1",
		"ANS":[2,2,3,2,1,2,3,1,2,3]
	  }*/
	//reads the submitted answer as json like above
	var answers Answers
	var correctCounter int
	err := json.NewDecoder(r.Body).Decode(&answers)
	if err != nil {
		log.Panic(err)
	}
	//stores inside quiz variable the submitted answers
	for i := 0; i < 10; i++ {
		quiz[answers.UserID].Questions[i].SubmittedANS = answers.PickedANS[i]
		if quiz[answers.UserID].Questions[i].CorrectANS == answers.PickedANS[i] {
			correctCounter++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(correctCounter)
}

func main() {
	//adds questions and JSON data, without database..
	initJSONExam()

	//to hadle requests, restAPI paths
	router := mux.NewRouter()
	router.HandleFunc("/quizdata", getQuizData).Methods("GET")
	router.HandleFunc("/questions/{id}", getQuestions).Methods("GET")
	router.HandleFunc("/submitanswer", submitAnswer).Methods("POST")
	router.HandleFunc("/submitanswers", submitAnswers).Methods("POST")
	//start the Listening
	log.Fatal(http.ListenAndServe(":5000", router))

}
