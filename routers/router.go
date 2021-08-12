package router

import (
	"encoding/base64"
	"encoding/json"
	"image"
	"io/ioutil"
	"log"
	"net/http"

	models "github.com/brwillian/api-rest/models"
	"gocv.io/x/gocv"
)

func GetClassificacao(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var msg models.Message
	json.Unmarshal(reqBody, &msg)

	base64Text := make([]byte, base64.StdEncoding.DecodedLen(len(msg.Image)))
	base64.StdEncoding.Decode(base64Text, []byte(msg.Image))

	mat, _ := gocv.IMDecode(base64Text, gocv.IMReadUnchanged)

	net := gocv.ReadNet("data/classificador.weights", "data/classificador.cfg")
	defer net.Close()

	OutputNames := GetOutputsNames(&net)

	detectClass := Detect(&net, mat, 0.6, 0.4, OutputNames)

	log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(detectClass)
}
func Detect(net *gocv.Net, src gocv.Mat, scoreThreshold float32, nmsThreshold float32, OutputNames []string) models.Response {
	blob := gocv.BlobFromImage(src, 1/255.0, image.Pt(320, 320), gocv.NewScalar(0, 0, 0, 0), true, false)
	net.SetInput(blob, "")
	probs := net.ForwardLayers(OutputNames)
	boxes, confidences, classIds := PostProcess(src, &probs)

	indices := make([]int, 6)
	if len(boxes) == 0 {
		return models.Response{}
	}
	gocv.NMSBoxes(boxes, confidences, scoreThreshold, nmsThreshold, indices)

	var result models.Response

	for _, idx := range indices {
		if idx == 0 {
			continue
		}

		detection := models.Detections{
			Id:         classIds[idx],
			Confidence: confidences[idx],
			Boxes: models.Boxes{
				X: boxes[idx].Min.X,
				Y: boxes[idx].Min.Y,
				W: boxes[idx].Max.X + boxes[idx].Min.X,
				H: boxes[idx].Max.Y + boxes[idx].Min.Y,
			},
		}

		result.Detections = append(result.Detections, detection)
	}

	return result
}
func GetOutputsNames(net *gocv.Net) []string {
	var outputLayers []string
	for _, i := range net.GetUnconnectedOutLayers() {
		layer := net.GetLayer(i)
		layerName := layer.GetName()
		if layerName != "_input" {
			outputLayers = append(outputLayers, layerName)
		}
	}
	return outputLayers
}
func PostProcess(frame gocv.Mat, outs *[]gocv.Mat) ([]image.Rectangle, []float32, []int) {
	var classIds []int
	var confidences []float32
	var boxes []image.Rectangle
	for _, out := range *outs {

		data, _ := out.DataPtrFloat32()
		for i := 0; i < out.Rows(); i, data = i+1, data[out.Cols():] {

			scoresCol := out.RowRange(i, i+1)

			scores := scoresCol.ColRange(5, out.Cols())
			_, confidence, _, classIDPoint := gocv.MinMaxLoc(scores)
			if confidence > 0.75 {

				centerX := int(data[0] * float32(frame.Cols()))
				centerY := int(data[1] * float32(frame.Rows()))
				width := int(data[2] * float32(frame.Cols()))
				height := int(data[3] * float32(frame.Rows()))

				left := centerX - width/2
				top := centerY - height/2
				classIds = append(classIds, classIDPoint.X)
				confidences = append(confidences, float32(confidence))
				boxes = append(boxes, image.Rect(left, top, width, height))
			}
		}
	}
	return boxes, confidences, classIds
}
func GetVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}
