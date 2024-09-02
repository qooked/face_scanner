package usecase

import (
	"context"
	"faceScanner/internal/constants"
	scannerErrors "faceScanner/internal/errors"
	"faceScanner/internal/models"
	"faceScanner/pkg/fileManager"
	"fmt"
	"log/slog"
	"sync"
)

type FaceScannerUsecase struct {
	repository            Repository
	tevianRequestProvider TevianRequestProvider
}

type TevianRequestProvider interface {
	ProvideRequest(image []byte) (tevianApiResponse models.TevianApiResponse, err error)
	GetResponse(body []byte) (tevianApiResponse models.TevianApiResponse, err error)
}

type Repository interface {
	ExtendFaceScannerTask(ctx context.Context, task models.ExtendFaceScannerTaskParamsRepository) (err error)
	GetFaceScannerTask(ctx context.Context, taskUUID string) (task models.GetFaceScannerTaskResponseRepository, err error)
	ChangeFaceScannerTaskStatus(ctx context.Context, taskUUID string, statusID int) (err error)
	DeleteFaceScannerTask(ctx context.Context, taskUUID string) (err error)
	CreateFaceScannerTask(ctx context.Context, task models.CreateFaceScannerTaskParamsRepository) (err error)
	UpdateTaskImageInfo(ctx context.Context, task models.UpdateTaskImageInfoParamsRepository) (err error)
}

func NewFaceScannerUsecase(repository Repository, tevianRequestProvider TevianRequestProvider) *FaceScannerUsecase {
	return &FaceScannerUsecase{
		tevianRequestProvider: tevianRequestProvider,
		repository:            repository,
	}
}

func (uc *FaceScannerUsecase) ExtendFaceScannerTask(ctx context.Context, task models.ExtendFaceScannerTaskUsecase) (err error) {
	taskRepo, err := uc.repository.GetFaceScannerTask(ctx, task.TaskUUID)
	if err != nil {
		err = fmt.Errorf("uc.repository.GetFaceScannerTask(...): %w", err)
		return err
	}

	if taskRepo.Status != constants.StatusNew {
		return scannerErrors.ErrTaskAlreadyStarted
	}

	fileName, err := fileManager.SaveFile(task.ImageUUID, task.Image)
	if err != nil {
		err = fmt.Errorf("fileManager.SaveFile(...): %w", err)
		return err
	}

	err = uc.repository.ExtendFaceScannerTask(ctx, models.ExtendFaceScannerTaskParamsRepository{
		TaskUUID:  task.TaskUUID,
		Image:     task.Image,
		ImageUUID: task.ImageUUID,
		FileName:  fileName,
	})
	if err != nil {
		err = fmt.Errorf("uc.repository.ExtendFaceScannerTask(...): %w", err)
		return err
	}

	return nil
}

func (uc *FaceScannerUsecase) GetFaceScannerTask(ctx context.Context, taskUUID string) (task models.GetFaceScannerTaskResponseUsecase, err error) {
	var (
		totalMaleCount, totalFemaleCount int
		totalMaleAge, totalFemaleAge     float64
	)
	taskRepo, err := uc.repository.GetFaceScannerTask(ctx, taskUUID)
	if err != nil {
		err = fmt.Errorf("uc.repository.GetFaceScannerTask(...): %w", err)
		return task, err
	}

	for _, storedData := range taskRepo.ImagesData {
		if !storedData.ApiResponse.Valid {
			continue
		}

		apiResponse, err := uc.tevianRequestProvider.GetResponse([]byte(storedData.ApiResponse.String))
		if err != nil {
			err = fmt.Errorf("uc.tevianRequestProvider.GetResponse(...): %w", err)
			return task, err
		}

		maleCount, femaleCount := apiResponse.GetMaleFemaleCount()
		totalMaleCount += maleCount
		totalFemaleCount += femaleCount
		var facesInPicture []models.Face

		for _, apiData := range apiResponse.Data {
			totalMaleAge += apiData.GetMaleAgeOrZero()
			totalFemaleAge += apiData.GetFemaleAgeOrZero()

			facesInPicture = append(facesInPicture, models.Face{
				BoundingBox: models.BoundingBox{
					X: apiData.BBox.X,
					Y: apiData.BBox.Y,
					W: apiData.BBox.Width,
					H: apiData.BBox.Height,
				},
				Sex: apiData.Demographics.Gender,
				Age: apiData.Demographics.Age.Mean,
			})
		}

		task.ImagesData = append(task.ImagesData, models.SingleTaskPictureUsecase{
			ApiResponse: apiResponse.BodyRaw,
			ImageUUID:   storedData.ImageUUID,
			Faces:       facesInPicture,
			FileName:    storedData.FileName,
		})
	}

	task.MaleFemaleCount = totalMaleCount + totalFemaleCount
	task.FacesCount = totalMaleCount + totalFemaleCount
	task.AverageMaleAge = totalMaleAge / float64(totalMaleCount)
	task.AverageFemaleAge = totalFemaleAge / float64(totalFemaleCount)
	task.TaskUUID = taskUUID
	task.Status = taskRepo.Status

	return task, nil
}

func (uc *FaceScannerUsecase) StartFaceScannerTask(ctx context.Context, taskUUID string) (err error) {
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

				err = uc.repository.UpdateTaskImageInfo(ctx, models.UpdateTaskImageInfoParamsRepository{
					ApiResponse: apiResponse.BodyRaw,
					TaskUUID:    taskUUID,
				})
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
func (uc *FaceScannerUsecase) DeleteFaceScannerTask(ctx context.Context, taskUUID string) (err error) {

	taskRepo, err := uc.repository.GetFaceScannerTask(ctx, taskUUID)
	if err != nil {
		err = fmt.Errorf("uc.repository.GetFaceScannerTask(...): %w", err)
		return err
	}
	for _, storedData := range taskRepo.ImagesData {
		err = fileManager.DeleteFile(storedData.ImageUUID)
		if err != nil {
			err = fmt.Errorf("fileManager.DeleteFile(...): %w", err)
			return err
		}
	}
	err = uc.repository.DeleteFaceScannerTask(ctx, taskUUID)
	if err != nil {
		err = fmt.Errorf("uc.repository.DeleteFaceScannerTask(...): %w", err)
		return err
	}

	return nil
}
func (uc *FaceScannerUsecase) CreateFaceScannerTask(ctx context.Context, task models.CreateFaceScannerTaskParamsUsecase) (err error) {
	fileName, err := fileManager.SaveFile(task.ImageUUID, task.Image)
	if err != nil {
		err = fmt.Errorf("fileManager.SaveFile(...): %w", err)
		return err
	}

	err = uc.repository.CreateFaceScannerTask(ctx, models.CreateFaceScannerTaskParamsRepository{
		TaskUUID:  task.TaskUUID,
		Image:     task.Image,
		ImageUUID: task.ImageUUID,
		FileName:  fileName,
	})
	if err != nil {
		err = fmt.Errorf("uc.repository.CreateFaceScannerTask(...): %w", err)
		return err
	}

	return nil
}
