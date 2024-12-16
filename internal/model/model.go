package model

type AirPollutionResponse struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	List []struct {
		Main struct {
			AQI int `json:"aqi"` // Индекс качества воздуха (1 - хорошо, 5 - плохо)
		} `json:"main"`
		Components struct {
			CO    float64 `json:"co"`    // Углекислый газ
			NO    float64 `json:"no"`    // Окись азота
			NO2   float64 `json:"no2"`   // Диоксид азота
			O3    float64 `json:"o3"`    // Озон
			SO2   float64 `json:"so2"`   // Диоксид серы
			PM2_5 float64 `json:"pm2_5"` // Мелкие частицы (меньше 2.5 мкм)
			PM10  float64 `json:"pm10"`  // Крупные частицы (до 10 мкм)
			NH3   float64 `json:"nh3"`   // Аммиак
		} `json:"components"`
		Dt int64 `json:"dt"` // Метка времени
	} `json:"list"`
}
