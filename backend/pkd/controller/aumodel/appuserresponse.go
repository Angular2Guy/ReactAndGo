package aubody

type AppUserResponse struct {
	Token        string
	Message      string
	Longitude    float64
	Latitude     float64
	SearchRadius float64
	TargetDiesel int
	TargetE5     int
	TargetE10    int
}
