/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start",
	Long:  `start longer`,
	Run: func(cmd *cobra.Command, args []string) {

		//according to the task we have many users
		println("please enter user ID")
		var userID string

		//reads user ID from user..
		fmt.Scanln(&userID)
		getQuesions(userID) //stors inside exam struct that defines globaly, the JSON data

		displayExam()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

type Question struct {
	Text         string    `json:"Text"`
	Options      [3]string `json:"options"`
	SubmittedANS int       `json:"ANS"`
	CorrectANS   int       `json:"-"`
}
type Exam struct {
	UserID    int        `json:"ID"`
	Questions []Question `json:"question"`
}

// defines a struct that stores submitted answers, to be handled.
type Answers struct {
	UserID    int     `json:"id"`
	PickedANS [10]int `json:"ANS"`
}

var answers Answers
var exam Exam

func getQuesions(usesrID string) {
	url := "http://localhost:5000/questions/" + usesrID
	request, _ := http.NewRequest(
		http.MethodGet, //method
		url,            //url
		nil,            //body
	)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("User-Agent", "CLI Muath Hawari")

	response, _ := http.DefaultClient.Do(request)

	body, _ := ioutil.ReadAll(response.Body)

	_ = json.Unmarshal(body, &exam)

	println("userid:" + strconv.Itoa(exam.UserID))
}

func displayExam() {

	//loop for 10 questions,, can be changed based on the lenth value
	for i := 0; i < 10; i++ {
		println("----------------------------------------------------------------")

		//prints questions that we got from getQuesions() and stored inside exam struct
		println("Question " + strconv.Itoa(i+1) + ":")
		println(exam.Questions[i].Text)

		//prints options for the current quesion.
		println("Options: ")
		for j := 0; j < 3; j++ {
			println(strconv.Itoa(j+1) + ": ANS(" + exam.Questions[i].Options[j] + ")")
		}

		//reads from user the pciked answer
		var option int
		fmt.Scanln(&option)
		//stores pciked answer inside some struct for buisness logic..
		exam.Questions[i].SubmittedANS = option
		answers.UserID = exam.UserID
		answers.PickedANS[i] = option
	}

	//posts answers struct as json
	postAnswers()
}

func postAnswers() {

	jsonValue, _ := json.Marshal(answers)
	url := "http://localhost:5000/submitanswers"

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println(err)
		return
	}

	//var to store the recieved correct answers from server
	var correctCounter int
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Could not read response body. %v", err)
	}
	err = json.Unmarshal(body, &correctCounter)
	if err != nil {
		log.Printf("%v", err)
	}
	print("the number of correct answers is: ")
	println(correctCounter)

}
