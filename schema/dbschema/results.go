package dbschema

type ResultModel struct {
	RoundId           string `ddb:"Round_ID,key"`
	PlayerId          string `ddb:"Player_ID,key"`
	PlayerRoundStatus string `ddb:"Player_Round_Status"`
	PlayerRoundRank   int    `ddb:"Player_Round_Rank"`
	PlayerRoundScore  int    `ddb:"Player_Round_Score"`
}
