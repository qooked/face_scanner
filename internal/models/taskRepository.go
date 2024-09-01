package models

import "github.com/guregu/null"

type GetFaceScannerTaskResponseRepository struct {
	TaskUUID   string `db:"id"`
	Status     int    `db:"status"`
	ImagesData []SingleTaskPictureRepository
}

type SingleTaskPictureRepository struct {
	ImageData   []byte      `db:"image_data"`
	ApiResponse null.String `db:"api_response"`
}

type ExtendFaceScannerTaskParamsRepository struct {
	TaskUUID string
	Image    []byte
}

type CreateFaceScannerTaskParamsRepository struct {
	TaskUUID string
	Image    []byte
}

func (t *GetFaceScannerTaskResponseRepository) ToUsecase() GetFaceScannerTaskResponseUsecase {

	var imagesData []SingleTaskPictureUsecase
	for _, image := range t.ImagesData {

		imagesData = append(imagesData, SingleTaskPictureUsecase{
			ImageData:   image.ImageData,
			ApiResponse: image.ApiResponse.String,
		})
	}
	return GetFaceScannerTaskResponseUsecase{
		TaskUUID:   t.TaskUUID,
		Status:     t.Status,
		ImagesData: imagesData,
	}
}
