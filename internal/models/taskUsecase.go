package models

type GetFaceScannerTaskResponseUsecase struct {
	TaskUUID   string
	Status     int
	ImagesData []SingleTaskPictureUsecase
}

type SingleTaskPictureUsecase struct {
	ImageData   []byte
	ApiResponse string
}
type ExtendFaceScannerTaskUsecase struct {
	TaskUUID string
	Image    []byte
}

type CreateFaceScannerTaskParamsUsecase struct {
	TaskUUID string
	Image    []byte
}
