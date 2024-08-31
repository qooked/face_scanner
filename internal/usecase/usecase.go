package usecase

import (
	"context"
	"faceScanner/internal/models"
)

type Usecase struct {
	repository Repository
}

type Repository interface {
	ExtendFaceScannerTask(ctx context.Context, task models.TaskParams) (err error)
	GetFaceScannerTask(ctx context.Context) (task models.TaskResponse, err error)
	StartFaceScannerTask(ctx context.Context, taskUUID string) (err error)
	DeleteFaceScannerTask(ctx context.Context, taskUUID string) (err error)
	CreateFaceScannerTask(ctx context.Context, task models.TaskParams) (err error)
}

func New(repository Repository) *Usecase {
	return &Usecase{
		repository: repository,
	}
}

func (uc *Usecase) ExtendFaceScannerTask(ctx context.Context, task models.TaskParams) (err error) {
	return err
}

func (uc *Usecase) GetFaceScannerTask(ctx context.Context) (task models.TaskResponse, err error) {
	return task, err
}
func (uc *Usecase) StartFaceScannerTask(ctx context.Context, taskUUID string) (err error) {
	return err
}
func (uc *Usecase) DeleteFaceScannerTask(ctx context.Context, taskUUID string) (err error) {
	return err
}
func (uc *Usecase) CreateFaceScannerTask(ctx context.Context, task models.TaskParams) (err error) {
	return err
}
