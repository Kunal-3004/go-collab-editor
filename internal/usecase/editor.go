package usecase

import (
	"collab-editor/internal/domain"
	"log"
)

type Repository interface {
	GetByID(id string) (*domain.Document, error)
	Save(doc *domain.Document) error
}

type Publisher interface {
	Broadcast(roomID string, message interface{})
}

type EditorService struct {
	repo      Repository
	publisher Publisher
}

func NewEditorService(r Repository, p Publisher) *EditorService {
	return &EditorService{repo: r, publisher: p}
}

func (s *EditorService) ProcessEdit(roomID string, op domain.Operation) error {
	doc, err := s.repo.GetByID(roomID)
	if err != nil {
		doc = domain.NewDocument(roomID)
	}

	if err := doc.ApplyOperation(op); err != nil {
		return err
	}

	if err := s.repo.Save(doc); err != nil {
		return err
	}

	log.Printf("Broadcasting op to room %s: %v", roomID, op)
	s.publisher.Broadcast(roomID, op)

	return nil
}

func (s *EditorService) GetDocument(roomID string) (*domain.Document, error) {
	return s.repo.GetByID(roomID)
}
