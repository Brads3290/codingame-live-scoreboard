package dbschema

type RoundModel struct {
	EventId string `ddb:"Event_ID,key" json:"event_id"`
	RoundId string `ddb:"Round_ID,key" json:"round_id"`
	Active  bool   `ddb:"Is_Active" json:"is_active"`
}
