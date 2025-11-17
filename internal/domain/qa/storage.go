package qa

type Storage interface {
    // Questions
    GetAllQuestions() ([]Question, error)
    CreateQuestion(q *Question) (uint64, error)
    GetQuestionWithAnswers(id uint64) (*Question, []Answer, error)
    DeleteQuestion(id uint64) error

    // Answers
    CreateAnswer(a *Answer) (uint64, error)
    GetAnswer(id uint64) (*Answer, error)
    DeleteAnswer(id uint64) error
}