package repositoryitems

import (
	"context"
	"fmt"

	"github.com/floffah/catena/internal/pkg/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repository *db.Queries
}

type CreateParams struct {
	Repository db.Repository
	Kind       db.RepositoryItemKind
	Title      string
	Body       *string
	AuthorID   pgtype.UUID
}

func NewService(conn db.DBTX) Service {
	return Service{
		repository: db.New(conn),
	}
}

func (s Service) Create(ctx context.Context, params CreateParams) (db.RepositoryItem, error) {
	number, err := s.repository.ReserveRepositoryItemNumber(ctx, params.Repository.ID)
	if err != nil {
		return db.RepositoryItem{}, err
	}

	item, err := s.repository.CreateRepositoryItem(ctx, db.CreateRepositoryItemParams{
		RepositoryID: params.Repository.ID,
		Number:       number,
		Kind:         params.Kind,
		Title:        params.Title,
		Body:         params.Body,
		AuthorID:     params.AuthorID,
	})
	if err != nil {
		return db.RepositoryItem{}, err
	}

	return item, nil
}

func Reference(repository db.Repository, item db.RepositoryItem) string {
	return ReferenceFromParts(repository.ItemPrefix, item.Number)
}

func ReferenceFromParts(prefix string, number int64) string {
	return fmt.Sprintf("%s-%d", prefix, number)
}
