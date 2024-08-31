package repository

import (
	"context"
	"faceScanner/internal/models"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}
func (h *Repository) ExtendFaceScannerTask(ctx context.Context, task models.TaskParams) (err error) {
	return nil
}

func (h *Repository) GetFaceScannerTask(ctx context.Context) (task models.TaskResponse, err error) {
	return task, nil
}
func (h *Repository) StartFaceScannerTask(ctx context.Context, taskUUID string) (err error) {
	return nil
}
func (h *Repository) DeleteFaceScannerTask(ctx context.Context, taskUUID string) (err error) {
	return nil
}
func (h *Repository) CreateFaceScannerTask(ctx context.Context, task models.TaskParams) (err error) {
	return nil
}
