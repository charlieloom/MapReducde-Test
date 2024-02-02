package Model

type QueryMsg struct {
	Condition Condition `json:"condition"`
	Offset    int       `json:"offset"`
	Limit     int       `json:"limit"`
	Row       int       `json:"row"`
	File      string    `json:"file"`
}
