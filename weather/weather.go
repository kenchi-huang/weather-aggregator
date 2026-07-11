package weather

type Weather struct {
	City     string
	Provider string
	Current  CurrentWeather
	Daily    []ForecastDay
}

type CurrentWeather struct {
	Temperature   float64
	Humidity      float64
	Precipitation float64
	MinTemp       float64
	MaxTemp       float64
	ConditionCode int
}

type ForecastDay struct {
	Date          string
	MinTemp       float64
	MaxTemp       float64
	ConditionCode int
	Hourly        []ForecastHour
}

type ForecastHour struct {
	Time          string
	Temperature   float64
	ConditionCode int
}

type Provider interface {
	GetWeather(lat float64, lon float64) (*Weather, error)
}

func getIconClassFromCode(code int) string {
	switch {
	case code == 0:
		return "fa-solid fa-sun" // ☀️
	case code == 1 || code == 2 || code == 3:
		return "fa-solid fa-cloud" // ☁️
	case code == 45 || code == 48:
		return "fa-solid fa-smog" // 🌫️
	case code >= 51 && code <= 67:
		return "fa-solid fa-cloud-rain" // 🌧️
	case code >= 71 && code <= 86:
		return "fa-solid fa-snowflake" // ❄️
	case code >= 95:
		return "fa-solid fa-cloud-bolt" // ⛈️
	default:
		return "fa-solid fa-temperature-half" // 🌡️
	}
}

// Update the methods to return IconClass()
func (c CurrentWeather) IconClass() string { return getIconClassFromCode(c.ConditionCode) }
func (f ForecastDay) IconClass() string    { return getIconClassFromCode(f.ConditionCode) }
func (h ForecastHour) IconClass() string   { return getIconClassFromCode(h.ConditionCode) }
