package apischema

type ApiScoreboardModel struct {
	EventId string               `json:"event_id"`
	Scores  []ApiScoreboardEntry `json:"scores"`
}

type ApiScoreboardEntry struct {
	Player ApiPlayerModel    `json:"player"`
	Score  ApiAggregateScore `json:"score"`
}

type ApiAggregateScore struct {
	Points int `json:"event_points"`
}
