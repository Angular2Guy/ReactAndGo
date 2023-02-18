package aubody

type AppUserRequest struct {
	Username     string
	Password     string
	Latitude     float64
	Longitude    float64
	SearchRadius float64
	TargetDiesel string
	TargetE10    string
	TargetE5     string
}
