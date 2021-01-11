package ddb

import (
	"codingame-live-scoreboard/schema"
	"codingame-live-scoreboard/schema/dbschema"
)

func UpdateDynamoDbFromScoreData(evtGuid string, data *schema.ScoreData) error {

	// Start by pulling some data that we need on separate threads
	chanPlayerResults := make(chan interface{}, 1)
	go func() {
		players, err := GetAllPlayersInEvent(evtGuid)
		if err != nil {
			chanPlayerResults <- err
			return
		}

		chanPlayerResults <- players
	}()

	// For rounds that are no longer active, update the database
	for i := range data.Rounds {
		if data.Rounds[i].Finished {
			err := SetRoundActive(evtGuid, data.Rounds[i].RoundId, false)
			if err != nil {
				return err
			}
		}
	}

	// Get the player data from chanPlayerResults
	players, err := getPlayerDataFromChan(chanPlayerResults)
	if err != nil {
		return err
	}

	// Check for any players in 'data' that do not already exist in 'players', and add them
	newPlayers := make([]dbschema.PlayerModel, 0)
	for _, round := range data.Rounds {
		for _, player := range round.Players {

			// Find player in "players"
			found := false
			for _, p := range players {
				if p.PlayerId == player.PlayerId {
					found = true
					break
				}
			}

			if found {
				continue
			}

			newPlayer := dbschema.PlayerModel{
				PlayerId: player.PlayerId,
				Name:     player.Name,
			}

			newPlayers = append(newPlayers, newPlayer)
		}
	}

	// Add the new players to dynamodb
	err = AddPlayersToEvent(evtGuid, newPlayers)
	if err != nil {
		return err
	}

	// Add the results to dynamodb, overwriting if necessary
	newResults := make([]dbschema.ResultModel, 0)
	for _, round := range data.Rounds {
		for _, player := range round.Players {

			result := dbschema.ResultModel{
				RoundId:           round.RoundId,
				PlayerId:          player.PlayerId,
				PlayerRoundStatus: player.SessionStatus,
				PlayerRoundRank:   player.Rank,
				PlayerRoundScore:  player.Score,
			}

			newResults = append(newResults, result)
		}
	}

	err = AddResultsToDynamoDb(newResults)
	if err != nil {
		return err
	}

	return nil
}

func getPlayerDataFromChan(c chan interface{}) ([]*dbschema.PlayerModel, error) {
	res := <-c

	switch rt := res.(type) {
	case []*dbschema.PlayerModel:
		return rt, nil
	case error:
		return nil, rt
	default:
		panic("unexpected type returned from chan")
	}
}
