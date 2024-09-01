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
	ImageUUID   string      `db:"image_id"`
	FileName    string      `db:"file_name"`
}

type ExtendFaceScannerTaskParamsRepository struct {
	TaskUUID  string
	Image     []byte
	ImageUUID string
	FileName  string
}

type CreateFaceScannerTaskParamsRepository struct {
	TaskUUID  string
	Image     []byte
	ImageUUID string
	FileName  string
}

type UpdateTaskImageInfoParamsRepository struct {
	ApiResponse string
	TaskUUID    string
}

func (t *GetFaceScannerTaskResponseRepository) ToUsecase() GetFaceScannerTaskResponseUsecase {

	var imagesData []SingleTaskPictureUsecase
	for _, image := range t.ImagesData {

		imagesData = append(imagesData, SingleTaskPictureUsecase{
			ImageData:   image.ImageData,
			ApiResponse: image.ApiResponse.String,
			ImageUUID:   image.ImageUUID,
		})
	}
	return GetFaceScannerTaskResponseUsecase{
		TaskUUID:   t.TaskUUID,
		Status:     t.Status,
		ImagesData: imagesData,
	}
}
