package schema

type StatsResponse struct {
	Round    *RoundStats    `json:"rounds"`
	Language *LanguageStats `json:"language"`
}

// Round stats

type RoundStats struct {
	NumberOfRounds int                          `json:"number_of_rounds"`
	PlayerStats    map[string]*PlayerRoundStats `json:"player_stats"`
}

type PlayerRoundStats struct {
	Username     string `json:"username"`
	RoundsPlayed int    `json:"rounds_played"`
}

// Language stats
type LanguageStats struct {
	Popularity []LanguagePopularityStat `json:"popularity"`
}

type LanguagePopularityStat struct {
	Name string `json:"name"`
	Rank int    `json:"rank"`
	Uses int    `json:"uses"`
}
