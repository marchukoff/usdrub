package main

type CBR struct {
	Date         string `json:"Date"`
	PreviousDate string `json:"PreviousDate"`
	PreviousURL  string `json:"PreviousURL"`
	Timestamp    string `json:"Timestamp"`
	Valute map[string]Currency `json:"Valute"`
}

type Currency struct {
	CharCode string  `json:"CharCode"`
	ID       string  `json:"ID"`
	Name     string  `json:"Name"`
	Nominal  int     `json:"Nominal"`
	NumCode  string  `json:"NumCode"`
	Previous float64 `json:"Previous"`
	Value    float64 `json:"Value"`
}

type Money struct {
	CharCode string  `json:"CharCode"`
	Value    float64 `json:"Value"`
}