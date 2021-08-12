package models

type Response struct {
	Detections []Detections `json:"Detections"`
}
type Boxes struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}
type Detections struct {
	Id         int     `json:"id"`
	Confidence float32 `json:"confidence"`
	Boxes      Boxes   `json:"boxes"`
}

type Message struct {
	Image string `json:"image"`
}
