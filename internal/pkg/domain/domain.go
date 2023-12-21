package domain

type Task struct {
	Uuid string `json:"uuid,omitempty"`
	Vus  int    `json:"vus,omitempty"`
	Dur  int    `json:"dur,omitempty"`

	VuNo     int      `json:"vuNo,omitempty"`
	Scenario Scenario `json:"scenario,omitempty"`
}

type Scenario struct {
	Name       string   `json:"name"`
	Processors []string `json:"processors"`
	Dur        int      `json:"dur,omitempty"`
}

type Metrics struct {
	Name      string `json:"name"`
	Value     string `gorm:"type:text" json:"value"`
	Timestamp string `json:"timestamp"`
}
