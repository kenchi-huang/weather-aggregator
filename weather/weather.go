package weather

type Weather struct {
	Temperature   float64
	City          string
	Provider      string
	Humidity      float64
	Precipitation float64
}

type Provider interface {
	GetWeather(lat float64, lon float64) (*Weather, error)
}
