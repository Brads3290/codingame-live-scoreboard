package ddb

import (
	"codingame-live-scoreboard/constants"
	"codingame-live-scoreboard/schema"
	"codingame-live-scoreboard/schema/dbschema"
)

func UpdateDynamoDbFromScoreData(evtGuid string, data *schema.ScoreData) error {

	// If data does not contain any rounds, we have nothing to update
	if len(data.Rounds) == 0 {
		return nil
	}

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

			// Check if we're already adding this player
			for _, p := range newPlayers {
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
				LanguageUsed:      player.LanguageUsed,
			}

			newResults = append(newResults, result)
		}
	}

	err = AddResultsToDynamoDb(newResults)
	if err != nil {
		return err
	}

	// Set 'Last Updated' for event
	err = MarkEventUpdatedNow(evtGuid)
	if err != nil {
		return err
	}

	return nil
}

func DeleteEvent(eventId string) error {

	// Delete the event
	err := DeleteKeysFromDynamoDb(constants.DB_TABLE_EVENTS, map[string]interface{}{
		"Event_ID": eventId,
	})

	if err != nil {
		return err
	}

	// Get all the rounds associated with the event
	rounds, err := GetAllRoundsForEvent(eventId)
	if err != nil {
		return err
	}

	// Delete the rounds
	err = DeleteRoundRecords(rounds)
	if err != nil {
		return err
	}

	// Extract the roundIds
	roundIds := make([]string, 0)
	for _, v := range rounds {
		roundIds = append(roundIds, v.RoundId)
	}

	// Get all the results from the rounds that we just deleted
	results, err := GetAllResultsForRounds(roundIds...)
	if err != nil {
		return err
	}

	// Delete result records
	err = DeleteResultRecords(results)
	if err != nil {
		return err
	}

	// Get all players associated with the event
	players, err := GetAllPlayersInEvent(eventId)
	if err != nil {
		return err
	}

	// Delete player records
	err = DeletePlayerRecords(players)
	if err != nil {
		return err
	}

	return nil
}

func DeleteRound(eventId string, roundId string) error {

	// Delete the round
	err := DeleteKeysFromDynamoDb(constants.DB_TABLE_ROUNDS, map[string]interface{}{
		"Event_ID": eventId,
		"Round_ID": roundId,
	})

	if err != nil {
		return err
	}

	// Query results related to this round and delete those too
	res, err := GetAllResultsForRound(roundId)
	if err != nil {
		return err
	}

	err = DeleteResultRecords(res)
	if err != nil {
		return err
	}

	return nil
}

func getPlayerDataFromChan(c chan interface{}) ([]dbschema.PlayerModel, error) {
	res := <-c

	switch rt := res.(type) {
	case []dbschema.PlayerModel:
		return rt, nil
	case error:
		return nil, rt
	default:
		panic("unexpected type returned from chan")
	}
}
