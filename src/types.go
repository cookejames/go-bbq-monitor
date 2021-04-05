package main

type TemperatureState struct {
	State struct {
		Reported struct {
			Temperature1 int16 `json:"temperature1,omitempty"`
			Temperature2 int16 `json:"temperature2,omitempty"`
			Temperature3 int16 `json:"temperature3,omitempty"`
			Temperature4 int16 `json:"temperature4,omitempty"`
		} `json:"reported"`
	} `json:"state"`
}

type ConnectionState struct {
	State struct {
		Reported struct {
			Connection string `json:"connection"`
		} `json:"reported"`
	} `json:"state"`
}
