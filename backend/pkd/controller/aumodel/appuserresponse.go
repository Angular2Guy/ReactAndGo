package aubody

type AppUserResponse struct {
	Token        string
	Message      string
	Uuid         string
	Longitude    float64
	Latitude     float64
	SearchRadius float64
	TargetDiesel string
	TargetE5     string
	TargetE10    string
}
