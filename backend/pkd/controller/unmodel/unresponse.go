package unbody

import "time"

type UnResponse struct {
	Timestamp time.Time
	UserUuid  string
	Title     string
	Message   string
	DataJson  string
}
