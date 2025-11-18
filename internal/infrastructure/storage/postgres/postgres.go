// Package postgres provides functionality for interacting with a PostgreSQL database.
package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"question-answer/internal/domain/qa"
	"question-answer/internal/infrastructure/storage/postgres/dto"
)

var (
	ErrOpenDB          = errors.New("failed to open database")
	ErrMigration       = errors.New("failed to run migrations")
	ErrGormOpen        = errors.New("failed to gorm open")
	ErrGetAllQuestions = errors.New("failed to get all questions")
	ErrCreateQuestion  = errors.New("failed to create question")
	ErrGetQuestion     = errors.New("failed to get question")
	ErrDeleteQuestion  = errors.New("failed to delete question")
	ErrCreateAnswer    = errors.New("failed to create answer")
	ErrGetAnswer       = errors.New("failed to get answer")
	ErrDeleteAnswer    = errors.New("failed to delete answer")
)

type PostgresStorage struct {
	db *gorm.DB
}

func New(cfg Config) (*PostgresStorage, error) {
	const op = "storage.postgres.NewStrorage"

	sqlDB, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, ErrOpenDB, err)
	}

	if err := goose.Up(sqlDB, cfg.MigrationsPath); err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, ErrMigration, err)
	}

	gormDB, err := gorm.Open(gormpg.New(gormpg.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, ErrGormOpen, err)
	}

	return &PostgresStorage{db: gormDB}, nil
}

func (s *PostgresStorage) GetAllQuestions() ([]qa.Question, error) {
	const op = "storage.postgres.GetAllQuestions"

	var dtos []pgdto.QuestionDTO

	if err := s.db.Order("id ASC").Find(&dtos).Error; err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, ErrGetAllQuestions, err)
	}

	res := make([]qa.Question, len(dtos))
	for i := range dtos {
		res[i] = pgdto.ToDomainQuestion(dtos[i])
	}

	return res, nil
}

func (s *PostgresStorage) CreateQuestion(q qa.Question) (uint64, error) {
	const op = "storage.postgres.CreateQuestion"

	dto := pgdto.ToDTOQuestion(q)

	if err := s.db.Create(&dto).Error; err != nil {
		return 0, fmt.Errorf("%s: %w: %w", op, ErrCreateQuestion, err)
	}

	return dto.ID, nil
}

func (s *PostgresStorage) GetQuestionWithAnswers(id uint64) (*qa.Question, []qa.Answer, error) {
	const op = "storage.postgres.GetQuestionWithAnswers"

	var qdto pgdto.QuestionDTO

	if err := s.db.First(&qdto, id).Error; err != nil {
		return nil, nil, fmt.Errorf("%s: %w: %w", op, ErrGetQuestion, err)
	}

	question := pgdto.ToDomainQuestion(qdto)

	var adtos []pgdto.AnswerDTO
	if err := s.db.Where("question_id = ?", id).
		Order("id ASC").
		Find(&adtos).Error; err != nil {

		return &question, nil, fmt.Errorf("%s: %w: %w", op, ErrGetQuestion, err)
	}

	answers := make([]qa.Answer, len(adtos))
	for i := range adtos {
		answers[i] = pgdto.ToDomainAnswer(adtos[i])
	}

	return &question, answers, nil
}

func (s *PostgresStorage) DeleteQuestion(id uint64) error {
	const op = "storage.postgres.DeleteQuestion"

	if err := s.db.Where("question_id = ?", id).
		Delete(&pgdto.AnswerDTO{}).Error; err != nil {

		return fmt.Errorf("%s: %w: %w", op, ErrDeleteQuestion, err)
	}

	if err := s.db.Delete(&pgdto.QuestionDTO{}, id).Error; err != nil {
		return fmt.Errorf("%s: %w: %w", op, ErrDeleteQuestion, err)
	}

	return nil
}

func (s *PostgresStorage) CreateAnswer(a qa.Answer) (uint64, error) {
	const op = "storage.postgres.CreateAnswer"

	dto := pgdto.ToDTOAnswer(a)

	if err := s.db.Create(&dto).Error; err != nil {
		return 0, fmt.Errorf("%s: %w: %w", op, ErrCreateAnswer, err)
	}

	return dto.ID, nil
}

func (s *PostgresStorage) GetAnswer(id uint64) (*qa.Answer, error) {
	const op = "storage.postgres.GetAnswer"

	var dto pgdto.AnswerDTO

	if err := s.db.First(&dto, id).Error; err != nil {
		return nil, fmt.Errorf("%s: %w: %w", op, ErrGetAnswer, err)
	}

	ans := pgdto.ToDomainAnswer(dto)
	return &ans, nil
}

func (s *PostgresStorage) DeleteAnswer(id uint64) error {
	const op = "storage.postgres.DeleteAnswer"

	if err := s.db.Delete(&pgdto.AnswerDTO{}, id).Error; err != nil {
		return fmt.Errorf("%s: %w: %w", op, ErrDeleteAnswer, err)
	}

	return nil
}
