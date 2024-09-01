package models

type TaskResponse struct {
	TaskUUID string
	Status   string
	Images   []Image
	Stats    []Stat
}
type Image struct {
	ImageName string
	Faces     []Face
}

type Face struct {
	BoundingBox []byte
	Sex         string
	Age         int
}

type Stat struct {
	FacesCount       int
	PeopleCount      int
	AverageMaleAge   int
	AverageFemaleAge int
}
