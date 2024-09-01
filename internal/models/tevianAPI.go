package models

type TevianApiResponse struct {
	BodyRaw    string     `json:"-"`
	Data       []FaceData `json:"data"`
	Rotation   int        `json:"rotation"`
	StatusCode int        `json:"status_code"`
}

type FaceData struct {
	Attributes   Attributes   `json:"attributes"`
	BBox         BBox         `json:"bbox"`
	Demographics Demographics `json:"demographics"`
	Landmarks    []Landmark   `json:"landmarks"`
	Liveness     int          `json:"liveness"`
	Masks        Masks        `json:"masks"`
	Quality      Quality      `json:"quality"`
	Score        float64      `json:"score"`
}

type BBox struct {
	Height int `json:"height"`
	Width  int `json:"width"`
	X      int `json:"x"`
	Y      int `json:"y"`
}

type Attributes struct {
	FacialHair string `json:"facial_hair"`
	Glasses    string `json:"glasses"`
	HairColor  string `json:"hair_color"`
	HairType   string `json:"hair_type"`
	Headwear   string `json:"headwear"`
}

type Demographics struct {
	Age       Age    `json:"age"`
	Ethnicity string `json:"ethnicity"`
	Gender    string `json:"gender"`
}

type Age struct {
	Mean     int `json:"mean"`
	Variance int `json:"variance"`
}

type Landmark struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Masks struct {
	FullFaceMask  int `json:"full_face_mask"`
	LowerFaceMask int `json:"lower_face_mask"`
	NoMask        int `json:"no_mask"`
	OtherMask     int `json:"other_mask"`
}

type Quality struct {
	Blurriness    int `json:"blurriness"`
	Overexposure  int `json:"overexposure"`
	Underexposure int `json:"underexposure"`
}
