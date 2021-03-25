// model.go

package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"sync/atomic"
)

var gradeSum [5]int32
var totalGrade int32

type question struct {
	ID       int    `json:"id"`
	Question string `json:"question"`
	Option0  string `json:"option0"`
	Option1  string `json:"option1"`
	Option2  string `json:"option2"`
}

type answer struct {
	ID     int `json:"id"`
	Answer int `json:"answer"`
}

type score struct {
	Result int     `json:"result"`
	Grade  float64 `json:"grade"`
}

type answers []answer

func getQuestions() ([]question, error) {
	var questions []question

	for i := 1; i <= 5; i++ {
		qdata := `{"Question": "The Question", "Option0": "Option A","Option1": "Option B", "Option2": "Option C"}`
		var obj question

		err := json.Unmarshal([]byte(qdata), &obj)
		obj.ID = rand.Intn(100)

		if err != nil {
			fmt.Println("error:", err)
		}
		questions = append(questions, obj)
	}
	return questions, nil
}

func (a answers) calculateResults() (score, error) {

	s := score{
		Result: 0,
		Grade:  0,
	}

	for _, ans := range a {
		res := math.Mod(float64(ans.ID), 3)
		if res == float64(ans.Answer) {
			s.Result++
		}
	}
	atomic.AddInt32(&gradeSum[s.Result], 1)
	atomic.AddInt32(&totalGrade, 1)
	s.Grade = calculateGrade(s.Result)
	return s, nil
}

func calculateGrade(r int) float64 {
	total := int32(0)
	for i := 0; i <= r; i++ {
		total += gradeSum[i]
	}
	return float64(total) / float64(totalGrade) * 100
}
