package cgschema

type ClashReportResponse struct {
	PublicHandle string           `json:"publicHandle"`
	Mode         string           `json:"mode"`
	Finished     bool             `json:"finished"`
	Started      bool             `json:"started"`
	Players      []PlayerResponse `json:"players"`
}

type PlayerResponse struct {
	PlayerId      int    `json:"codingamerId"`
	Nickname      string `json:"codingamerNickname"`
	RoundScore    int    `json:"score"`
	DurationMs    int    `json:"duration"`
	SessionStatus string `json:"testSessionStatus"`
	Rank          int    `json:"rank"`
}
