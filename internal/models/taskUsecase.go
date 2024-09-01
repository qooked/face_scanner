package models

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
type GetFaceScannerTaskResponseUsecase struct {
	TaskUUID         string
	Status           int
	PictureName      string
	ImagesData       []SingleTaskPictureUsecase
	FacesCount       int
	MaleFemaleCount  int
	AverageMaleAge   float64
	AverageFemaleAge float64
}

type SingleTaskPictureUsecase struct {
	ApiResponse string
	ImageData   []byte
	ImageUUID   string
	Faces       []Face
	FileName    string
}

type Face struct {
	BoundingBox `json:"boundingBox"`
	Sex         string  `json:"sex"`
	Age         float64 `json:"age"`
}

type BoundingBox struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}
