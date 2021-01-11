package schema

type ScoreData struct {
	EventId string
	Rounds  []RoundData `json:"active_rounds"`
}

type RoundData struct {
	EventId  string
	RoundId  string            `json:"round_id"`
	Mode     string            `json:"mode"`
	Finished bool              `json:"finished"`
	Players  []PlayerRoundData `json:"players"`
}

type PlayerRoundData struct {
	PlayerId      string `json:"player_id"`
	Name          string `json:"name"`
	Rank          int    `json:"rank"`
	Score         int    `json:"score"`
	SessionStatus string `json:"session_status"`
}

type PlayerData struct {
	EventId  string
	PlayerId string
	Name     string
}

type ResultData struct {
	RoundId  string
	PlayerId string
	Status   string
	Rank     int
	Score    int
}
