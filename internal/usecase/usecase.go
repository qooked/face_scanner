package usecase

import (
	"context"
	"faceScanner/internal/constants"
	scannerErrors "faceScanner/internal/errors"
	"faceScanner/internal/models"
	"fmt"
	"log/slog"
	"sync"
)

type Usecase struct {
	repository            Repository
	tevianRequestProvider TevianRequestProvider
}

type TevianRequestProvider interface {
	ProvideRequest(image []byte) (tevianApiResponse models.TevianApiResponse, err error)
}

type Repository interface {
	ExtendFaceScannerTask(ctx context.Context, task models.ExtendFaceScannerTaskParamsRepository) (err error)
	GetFaceScannerTask(ctx context.Context, taskUUID string) (task models.GetFaceScannerTaskResponseRepository, err error)
	ChangeFaceScannerTaskStatus(ctx context.Context, taskUUID string, statusID int) (err error)
	DeleteFaceScannerTask(ctx context.Context, taskUUID string) (err error)
	CreateFaceScannerTask(ctx context.Context, task models.CreateFaceScannerTaskParamsRepository) (err error)
	UpdateTaskImageInfo(ctx context.Context, apiResponse string, taskUUID string) (err error)
}

func New(repository Repository, tevianRequestProvider TevianRequestProvider) *Usecase {
	return &Usecase{
		tevianRequestProvider: tevianRequestProvider,
		repository:            repository,
	}
}

func (uc *Usecase) ExtendFaceScannerTask(ctx context.Context, task models.ExtendFaceScannerTaskUsecase) (err error) {
	taskRepo, err := uc.repository.GetFaceScannerTask(ctx, task.TaskUUID)
	if err != nil {
		err = fmt.Errorf("uc.repository.GetFaceScannerTask(...): %w", err)
		return err
	}
	if taskRepo.Status != constants.StatusNew {
		return scannerErrors.ErrTaskAlreadyStarted
	}

	err = uc.repository.ExtendFaceScannerTask(ctx, models.ExtendFaceScannerTaskParamsRepository{
		TaskUUID: task.TaskUUID,
		Image:    task.Image,
	})
	if err != nil {
		err = fmt.Errorf("uc.repository.ExtendFaceScannerTask(...): %w", err)
		return err
	}

	return nil
}

func (uc *Usecase) GetFaceScannerTask(ctx context.Context, taskUUID string) (task models.GetFaceScannerTaskResponseUsecase, err error) {
	taskRepo, err := uc.repository.GetFaceScannerTask(ctx, taskUUID)
	if err != nil {
		err = fmt.Errorf("uc.repository.GetFaceScannerTask(...): %w", err)
		return task, err
	}

	return taskRepo.ToUsecase(), nil
}

func (uc *Usecase) StartFaceScannerTask(ctx context.Context, taskUUID string) (err error) {
	var (
		wg           sync.WaitGroup
		statusID     = constants.StatusPending
		responseChan = make(chan models.TevianApiResponse)
		errorChan    = make(chan error)
	)

	taskRepo, err := uc.repository.GetFaceScannerTask(ctx, taskUUID)
	if err != nil {
		err = fmt.Errorf("uc.repository.GetFaceScannerTask(...): %w", err)
		return err
	}
	taskUC := taskRepo.ToUsecase()
	if taskUC.Status != constants.StatusNew {
		return scannerErrors.ErrTaskAlreadyStarted
	}

	err = uc.repository.ChangeFaceScannerTaskStatus(ctx, taskUUID, statusID)
	if err != nil {
		err = fmt.Errorf("uc.repository.ChangeFaceScannerTaskStatus(...): %w", err)
		return err
	}

	for i := 0; i < len(taskUC.ImagesData); i++ {
		wg.Add(1)

		go func(errChan chan error, respChan chan models.TevianApiResponse, image []byte, wg *sync.WaitGroup) {
			defer wg.Done()

			tevianApiResponse, err := uc.tevianRequestProvider.ProvideRequest(image)
			if err != nil {
				err = fmt.Errorf("uc.tevianRequestProvider.ProvideRequest(...): %w", err)
				errorChan <- err
				return
			}

			respChan <- tevianApiResponse
			return
		}(errorChan, responseChan, taskUC.ImagesData[i].ImageData, &wg)
	}

	go func(chan error) {
		var errors []error
		for i := 0; i < len(taskUC.ImagesData); i++ {
			select {
			case <-ctx.Done():
				slog.Info("ctx.Done()")
				return
			case err := <-errorChan:
				if err != nil {
					slog.Error(err.Error())
					errors = append(errors, err)
				}
			case apiResponse := <-responseChan:
				err = uc.repository.UpdateTaskImageInfo(ctx, apiResponse.BodyRaw, taskUUID)
				if err != nil {
					slog.Error(err.Error())
				}
			}
		}
		wg.Wait()
		switch len(errors) {
		case 0:
			statusID = constants.StatusSuccess
		case len(errors):
			statusID = constants.StatusFailed
		default:
			statusID = constants.StatusPartiallySuccess
		}

		err = uc.repository.ChangeFaceScannerTaskStatus(ctx, taskUUID, statusID)
		if err != nil {
			err = fmt.Errorf("uc.repository.ChangeFaceScannerTaskStatus(...): %w", err)
			slog.Error(err.Error())
		}
	}(errorChan)

	return err
}
func (uc *Usecase) DeleteFaceScannerTask(ctx context.Context, taskUUID string) (err error) {
	err = uc.repository.DeleteFaceScannerTask(ctx, taskUUID)
	if err != nil {
		err = fmt.Errorf("uc.repository.DeleteFaceScannerTask(...): %w", err)
		return err
	}

	return nil
}
func (uc *Usecase) CreateFaceScannerTask(ctx context.Context, task models.CreateFaceScannerTaskParamsUsecase) (err error) {
	err = uc.repository.CreateFaceScannerTask(ctx, models.CreateFaceScannerTaskParamsRepository{
		TaskUUID: task.TaskUUID,
		Image:    task.Image,
	})
	if err != nil {
		err = fmt.Errorf("uc.repository.CreateFaceScannerTask(...): %w", err)
		return err
	}

	return nil
}
