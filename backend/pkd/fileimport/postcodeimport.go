/*
  - Copyright 2022 Sven Loesekann
    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/
package aufile

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"os"
	appuser "react-and-go/pkd/appuser"
	gasstation "react-and-go/pkd/gasstation"
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

func UpdatePostCodeCoordinates(fileName string) {
	gzReader, file, err := createReader(fileName)
	defer gzReader.Close()
	defer file.Close()

	if err != nil {
		return
	}

	jsonDecoder := json.NewDecoder(gzReader)
	plzContainerNumber := 0
	result := []appuser.PostCodeData{}

	jsonDecoder.Token()
	for jsonDecoder.More() {
		myPlzContainer := plzContainer{}
		jsonDecoder.Decode(&myPlzContainer)
		plzContainerNumber++
		myPostCode := createPostCode(&myPlzContainer)
		result = append(result, myPostCode)
		//log.Printf("PostCode: %v\n", myPostCode)
	}
	jsonDecoder.Token()
	//log.Printf("Number of postcodes: %v\n", plzContainerNumber)
	appuser.ImportPostCodeData(result)
}

func UpdateStatesAndCounties(fileName string) {
	gzReader, file, err := createReader(fileName)
	if err != nil {
		return
	}

	defer gzReader.Close()
	defer file.Close()

	stateToAmount := make(map[string]int)
	plzToState := make(map[string]string)
	plzToCounty := make(map[string]string)
	lineId := 0
	scanner := bufio.NewScanner(gzReader)
	for scanner.Scan() {
		line := scanner.Text()
		lineTokens := strings.Split(line, ",")
		if lineId == 0 || len(lineTokens) < 5 {
			lineId += 1
			continue
		}
		//log.Printf(line)
		plzToCounty[lineTokens[3]] = lineTokens[4]
		plzToState[lineTokens[3]] = lineTokens[5]
		if _, ok := stateToAmount[lineTokens[5]]; ok {
			stateToAmount[lineTokens[5]] = stateToAmount[lineTokens[5]] + 1
		} else {
			stateToAmount[lineTokens[5]] = 1
		}
		lineId += 1
	}

	appuser.UpdateStatesCounties(plzToState, plzToCounty)
	gasstation.UpdateStatesCounties(plzToState, plzToCounty)
}

func createReader(fileName string) (*gzip.Reader, *os.File, error) {
	filePath := strings.TrimSpace(os.Getenv("PLZ_IMPORT_PATH"))
	log.Printf("File: %v%v", filePath, fileName)
	file, err := os.Open(fmt.Sprintf("%v%v", filePath, strings.TrimSpace(fileName)))
	if err != nil {
		log.Printf("Failed to open file: %v, %v\n", fmt.Sprintf("%v%v", filePath, strings.TrimSpace(fileName)), err.Error())
		return nil, nil, err
	}
	gzReader, err := gzip.NewReader(bufio.NewReader(file))
	if err != nil {
		file.Close()
		log.Printf("Failed to create buffered gzip reader: %v, %v\n", fmt.Sprintf("%v%v", filePath, strings.TrimSpace(fileName)), err.Error())
		return nil, nil, err
	}

	return gzReader, file, nil
}

func createPostCode(plzContainer *plzContainer) appuser.PostCodeData {
	postCodeData := appuser.PostCodeData{}
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
