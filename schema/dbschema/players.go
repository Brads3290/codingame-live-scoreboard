package dbschema

type PlayerModel struct {
	EventId  string `ddb:"Event_ID,key"`
	PlayerId string `ddb:"Player_ID,key"`
	Name     string `ddb:"Player_Name"`
}
