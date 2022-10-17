package lesson01

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func Run() {
	problemsFile := flag.String("problems", "lesson01/problems.csv", "Path of the CSV file to source the problems from.")
	quizTimeout := flag.Int("timeout", 30, "Timeout for quiz in seconds.")
	shuffle := flag.Bool("shuffle", false, "To list questions in a shuffled manner.")
	flag.Parse()
	questions, err := loadQuestionsFromCSV(*problemsFile)
	if err != nil {
		log.Fatalf("Error in loading - %v", err)
	}
	repo := NewRepository(questions)
	quiz := NewQuiz(repo, *shuffle)
	fmt.Printf("Quiz with time duration [%d] seconds and problems source as [%s] with shuffle [%s]. Press enter to start!\n",
		*quizTimeout, *problemsFile, strconv.FormatBool(*shuffle))
	fmt.Scanln()
	quizComplete := make(chan error)
	go func() {
		quizComplete <- startQuiz(&quiz)
	}()
	timer := time.After(time.Duration(*quizTimeout) * time.Second)
	select {
	case <-quizComplete:
		break
	case <-timer:
		fmt.Printf("\n\n!!! Time Over !!!\n\n")
		break
	}
	fmt.Printf("Result:\n%s\n", quiz.Result().Json())
}

func startQuiz(quiz *Quiz) error {
	var ans string
	var question *Question
	var err error
	for {
		question, err = quiz.Next()
		if err != nil {
			break
		}
		fmt.Printf("Q: %s\n", question.text)
		fmt.Scanf("%s", &ans)
		quiz.AddAnswer(question.id, strings.Trim(ans, " "))
		fmt.Printf("A: %s\n", ans)
		ans = ""
	}
	return nil
}

func loadQuestionsFromCSV(fileName string) ([]Question, error) {
	questions := []Question{}
	file, err := os.Open(fileName)
	if err != nil {
		return questions, err
	}
	reader := csv.NewReader(file)
	qid := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return questions, err
		}
		questions = append(questions, Question{
			id:   qid,
			text: record[0],
			answer: Answer{
				value: record[1],
			},
		})
		qid += 1
	}
	return questions, nil
}
