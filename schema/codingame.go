package schema

type CodinGameClashReportResponse struct {
	PublicHandle string                    `json:"publicHandle"`
	Mode         string                    `json:"mode"`
	Finished     bool                      `json:"finished"`
	Started      bool                      `json:"started"`
	Players      []CodinGamePlayerResponse `json:"players"`
}

type CodinGamePlayerResponse struct {
	Nickname      string `json:"codingamerNickname"`
	RoundScore    int    `json:"score"`
	DurationMs    int    `json:"duration"`
	SessionStatus string `json:"testSessionStatus"`
	Rank          int    `json:"rank"`
}
