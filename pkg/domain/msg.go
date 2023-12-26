package _domain

type MqMsg struct {
	Namespace string `json:"namespace"`
	Room      string `json:"room"`
	Event     string `json:"event"`
	Content   string `json:"content"`
}
