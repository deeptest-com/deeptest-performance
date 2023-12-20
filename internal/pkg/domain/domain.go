package domain

type Metrics struct {
	Name      string `json:"name"`
	Value     string `gorm:"type:text" json:"value"`
	Timestamp string `json:"timestamp"`
}
