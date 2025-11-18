package qa

type Service interface {
    // Questions
    GetAllQuestions() ([]Question, error)
    CreateQuestion(q Question) (*Question, error)
    GetQuestionWithAnswers(id uint64) (*Question, []Answer, error)
    DeleteQuestion(id uint64) error

    // Answers
    CreateAnswer(a Answer) (uint64, error)
    GetAnswer(id uint64) (*Answer, error)
    DeleteAnswer(id uint64) error
}

type service struct {
    storage Storage
}

func NewService(storage Storage) Service {
    return &service{storage: storage}
}

func (s *service) GetAllQuestions() ([]Question, error) {
    return s.storage.GetAllQuestions()
}

func (s *service) CreateQuestion(q Question) (*Question, error) {
    return s.storage.CreateQuestion(q)
}

func (s *service) GetQuestionWithAnswers(id uint64) (*Question, []Answer, error) {
    return s.storage.GetQuestionWithAnswers(id)
}

func (s *service) DeleteQuestion(id uint64) error {
    return s.storage.DeleteQuestion(id)
}

func (s *service) CreateAnswer(a Answer) (uint64, error) {
    return s.storage.CreateAnswer(a)
}

func (s *service) GetAnswer(id uint64) (*Answer, error) {
    return s.storage.GetAnswer(id)
}

func (s *service) DeleteAnswer(id uint64) error {
    return s.storage.DeleteAnswer(id)
}
