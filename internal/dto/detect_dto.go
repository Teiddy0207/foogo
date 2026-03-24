package dto

type AnalyzeFoodRequest struct {
	ObjectKey string `json:"object_key"`
}

type FoodItem struct {
	Name        string  `json:"name"`
	Confidence  float32 `json:"confidence"`
	CaloriesEst float32 `json:"calories_est"`
}

type AnalyzeFoodResponse struct {
	Items []FoodItem `json:"items"`
	Note  string     `json:"note"`
}
