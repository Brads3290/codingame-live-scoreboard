package apishared

import (
	"codingame-live-scoreboard/ddb"
	"codingame-live-scoreboard/schema/apischema"
	"codingame-live-scoreboard/schema/dbschema"
)

func GetApiScoreboardData(evtGuid string) (*apischema.ApiScoreboardModel, error) {
	chPlayers := make(chan interface{}, 1)
	go func() {
		players, err := ddb.GetAllPlayersInEvent(evtGuid)
		if err != nil {
			chPlayers <- err
			return
		}

		chPlayers <- players
	}()

	chResults := make(chan interface{}, 1)
	go func() {
		rounds, err := ddb.GetAllRoundsForEvent(evtGuid)
		if err != nil {
			chResults <- err
			return
		}

		roundIds := make([]string, len(rounds))
		for i, v := range rounds {
			roundIds[i] = v.RoundId
		}

		results, err := ddb.GetAllResultsForRounds(roundIds...)
		if err != nil {
			chResults <- err
			return
		}

		chResults <- results
	}()

	outPlayers := <-chPlayers
	outResults := <-chResults

	var players []dbschema.PlayerModel
	switch pt := outPlayers.(type) {
	case error:
		return nil, pt
	case []dbschema.PlayerModel:
		players = pt
	default:
		panic("unexpected type")
	}

	var results []dbschema.ResultModel
	switch rt := outResults.(type) {
	case error:
		return nil, rt
	case []dbschema.ResultModel:
		results = rt
	default:
		panic("unexpected type")
	}

	// Calculate all the players' event scores
	playerScores := CalculatePlayerScoresForEvent(results)

	// For each player, retrieve event score and add them to the scoreboard return object
	var sb apischema.ApiScoreboardModel
	sb.EventId = evtGuid
	sb.Scores = make([]apischema.ApiScoreboardEntry, 0)

	for _, v := range players {
		playerScore, ok := playerScores[v.PlayerId]
		if !ok {
			// Player has no results..
			continue
		}

		sbe := apischema.ApiScoreboardEntry{
			Player: apischema.ApiPlayerModel{
				PlayerId: v.PlayerId,
				Name:     v.Name,
			},
			Score: apischema.ApiAggregateScore{
				Points: playerScore,
			},
		}

		sb.Scores = append(sb.Scores, sbe)
	}

	return &sb, nil
}

func CalculatePlayerScoresForEvent(results []dbschema.ResultModel) map[string]int {
	type playerRoundResult struct {
		rank int
	}

	// How many players were in each round?
	roundPlayerResults := make(map[string]map[string]playerRoundResult)
	for _, r := range results {
		if _, ok := roundPlayerResults[r.RoundId]; !ok {
			roundPlayerResults[r.RoundId] = make(map[string]playerRoundResult)
		}

		roundPlayerResults[r.RoundId][r.PlayerId] = playerRoundResult{rank: r.PlayerRoundRank}
	}

	// For each round, calculate the points each player gets and add it to their total
	playerPoints := make(map[string]int)
	for _, playerResults := range roundPlayerResults {
		participants := len(playerResults)

		for playerId, playerResult := range playerResults {
			if _, ok := playerPoints[playerId]; !ok {
				playerPoints[playerId] = 0
			}

			playerPoints[playerId] += participants - playerResult.rank + 1
		}
	}

	return playerPoints
}
