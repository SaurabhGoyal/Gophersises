package lesson01

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Quiz struct {
	mu      sync.Mutex
	repo    Repository
	answers map[int]string
	correct int
	next    int
}

func (q *Quiz) AddAnswer(qid int, ans string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if ans == q.answers[qid] {
		return
	}
	_, answered := q.answers[qid]
	q.answers[qid] = ans
	if ans == q.repo.index[qid].answer.value && !answered {
		q.correct += 1
	}
}

var ErrNoMoreQuestions = errors.New("no more questions")

func (q *Quiz) Next() (*Question, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.next == len(q.repo.questions) {
		return &Question{}, ErrNoMoreQuestions
	}
	nextQuestion := q.repo.questions[q.next]
	q.next += 1
	return nextQuestion, nil
}

type Result struct {
	Answers   map[string]string `json:"answers"`
	Correct   int               `json:"correct"`
	Incorrect int               `json:"incorrect"`
}

func (r Result) Json() string {
	res, err := json.Marshal(r)
	if err != nil {
		log.Fatalf("Error in formatting result - %v", err)
	}
	return string(res)
}

func (q *Quiz) Result() Result {
	q.mu.Lock()
	defer q.mu.Unlock()
	answers := map[string]string{}
	for qid, ans := range q.answers {
		answers[q.repo.index[qid].text] = ans
	}
	return Result{
		Answers:   answers,
		Correct:   q.correct,
		Incorrect: len(answers) - q.correct,
	}
}

func NewQuiz(repo Repository, shuffle bool) Quiz {
	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(repo.questions), func(i, j int) {
			repo.questions[i], repo.questions[j] = repo.questions[j], repo.questions[i]
		})
	}
	return Quiz{
		repo:    repo,
		answers: map[int]string{},
		correct: 0,
		next:    0,
	}
}
