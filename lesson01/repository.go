package lesson01

type Answer struct {
	value string
}

type Question struct {
	id     int
	text   string
	answer Answer
}

type Repository struct {
	questions []*Question
	index     map[int]*Question
}

func NewRepository(questions []Question) Repository {
	q := []*Question{}
	i := map[int]*Question{}
	for _, question := range questions {
		rQuestion := question
		q = append(q, &rQuestion)
		i[question.id] = &rQuestion
	}
	return Repository{
		questions: q,
		index:     i,
	}
}
