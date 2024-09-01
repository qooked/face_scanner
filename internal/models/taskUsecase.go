package models

type GetFaceScannerTaskResponseUsecase struct {
	TaskUUID   string
	Status     int
	ImagesData []SingleTaskPictureUsecase
}

type SingleTaskPictureUsecase struct {
	ImageData   []byte
	ApiResponse string
	ImageUUID   string
}
type ExtendFaceScannerTaskUsecase struct {
	TaskUUID  string
	Image     []byte
	ImageUUID string
}

type CreateFaceScannerTaskParamsUsecase struct {
	TaskUUID  string
	Image     []byte
	ImageUUID string
}
