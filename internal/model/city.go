package model

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lon"`
}

type City struct {
	CityID   string   `json:"-"`
	Name     string   `json:"name"`
	Region   string   `json:"region"`
	Location Location `json:"location"`
}
