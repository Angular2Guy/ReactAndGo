package aufile

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type coordinateTuple [2]float64

type coordinateList []coordinateTuple

type plzPolygon struct {
	Typestr     string `json:"type"`
	Coordinates []coordinateList
}

type plzProperties struct {
	Plz        int32   `json:"plz,string"`
	Label      string  `json:"note"`
	Population int32   `json:"einwohner"`
	SquareKM   float32 `json:"qkm"`
}

type plzContainer struct {
	Typestr    string        `json:"type"`
	Properties plzProperties `json:"properties"`
	Geometry   plzPolygon    `json:"geometry"`
}

type PostCodeData struct {
	Label           string
	PostCode        int32
	Population      int32
	SquareKM        float32
	CenterLongitude float64
	CenterLatitude  float64
}

func UpdatePlzCoordinates(fileName string) {
	filePath := strings.TrimSpace(os.Getenv("PLZ_IMPORT_PATH"))
	log.Printf("File: %v%v", filePath, fileName)
	file, err := os.Open(fmt.Sprintf("%v%v", filePath, strings.TrimSpace(fileName)))
	defer file.Close()
	if err != nil {
		log.Printf("Failed to open file: %v, %v\n", fmt.Sprintf("%v%v", filePath, strings.TrimSpace(fileName)), err.Error())
		return
	}
	gzReader, err := gzip.NewReader(bufio.NewReader(file))
	if err != nil {
		log.Printf("Failed to create buffered gzip reader: %v, %v\n", fmt.Sprintf("%v%v", filePath, strings.TrimSpace(fileName)), err.Error())
		return
	}
	defer gzReader.Close()

	jsonDecoder := json.NewDecoder(gzReader)
	plzContainerNumber := 0

	jsonDecoder.Token()
	for jsonDecoder.More() {
		myPlzContainer := plzContainer{}
		jsonDecoder.Decode(&myPlzContainer)
		plzContainerNumber++
		myPostCode := createPostCode(&myPlzContainer)
		log.Printf("PostCode: %v\n", myPostCode)
	}
	jsonDecoder.Token()
	log.Printf("Number of postcodes: %v", plzContainerNumber)
}

func createPostCode(plzContainer *plzContainer) PostCodeData {
	postCodeData := PostCodeData{}
	postCodeData.Label = plzContainer.Properties.Label
	postCodeData.PostCode = plzContainer.Properties.Plz
	postCodeData.SquareKM = plzContainer.Properties.SquareKM
	postCodeData.Population = plzContainer.Properties.Population
	postCodeData.CenterLongitude, postCodeData.CenterLatitude = calcCentroid(plzContainer)
	return postCodeData
}

func calcCentroid(plzContainer *plzContainer) (float64, float64) {
	polygonArea := calcPolygonArea(plzContainer)
	//log.Printf("PolygonArea: %v", polygonArea)
	coordinateLists := plzContainer.Geometry.Coordinates
	centerLongitude := 0.0
	centerLatitude := 0.0
	for _, coordinateTuples := range coordinateLists {
		for index, coordinateTuple := range coordinateTuples {
			if index >= len(coordinateTuples)-1 {
				continue
			}
			centerLongitude += (coordinateTuple[0] + coordinateTuples[index+1][0]) * (coordinateTuple[0]*coordinateTuples[index+1][1] - coordinateTuples[index+1][0]*coordinateTuple[1])
			centerLatitude += (coordinateTuple[1] + coordinateTuples[index+1][1]) * (coordinateTuple[0]*coordinateTuples[index+1][1] - coordinateTuples[index+1][0]*coordinateTuple[1])
		}
	}
	centerLongitude = centerLongitude / (6 * polygonArea)
	centerLatitude = centerLatitude / (6 * polygonArea)
	return centerLongitude, centerLatitude
}

func calcPolygonArea(plzContainer *plzContainer) float64 {
	coordinateLists := plzContainer.Geometry.Coordinates
	polygonArea := 0.0
	for _, coordinateTuples := range coordinateLists {
		for index, coordinateTuple := range coordinateTuples {
			if index >= len(coordinateTuples)-1 {
				continue
			}
			polygonArea += coordinateTuple[0]*coordinateTuples[index+1][1] - coordinateTuples[index+1][0]*coordinateTuple[1]
		}
	}
	polygonArea = polygonArea / 2
	return polygonArea
}
