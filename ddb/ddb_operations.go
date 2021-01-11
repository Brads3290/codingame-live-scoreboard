package ddb

import (
	"codingame-live-scoreboard/constants"
	"codingame-live-scoreboard/schema/dbschema"
	"codingame-live-scoreboard/schema/errors"
)

func GetEvent(evtGuid string) (*dbschema.EventModel, error) {
	var em dbschema.EventModel
	err := GetItemFromDynamoDb(constants.DB_TABLE_EVENTS, &em, map[string]interface{}{
		"Event_ID": evtGuid,
	})

	if err != nil {
		return nil, err
	}

	return &em, nil
}

func SetEventUpdating(evtGuid string, isUpdating bool) error {
	err := UpdateItemAttrsInDynamoDb(constants.DB_TABLE_EVENTS, map[string]interface{}{
		"Event_ID": evtGuid,
	}, map[string]interface{}{
		"Is_Updating": isUpdating,
	})

	if err != nil {
		return err
	}

	return nil
}

func SetRoundActive(evtGuid string, roundId string, isActive bool) error {
	err := UpdateItemAttrsInDynamoDb(constants.DB_TABLE_ROUNDS, map[string]interface{}{
		"Event_ID": evtGuid,
		"Round_ID": roundId,
	}, map[string]interface{}{
		"Is_Active": isActive,
	})

	if err != nil {
		return err
	}

	return nil
}

func GetAllRoundsForEvent(evtGuid string) ([]dbschema.RoundModel, error) {
	var rms []dbschema.RoundModel

	err := QueryItemsFromDynamoDb(constants.DB_TABLE_ROUNDS, &rms, map[string]interface{}{
		"Event_ID": evtGuid,
	})

	if err != nil {
		return nil, err
	}

	return rms, nil
}

func GetActiveRoundsForEvent(evtGuid string) ([]dbschema.RoundModel, error) {
	var rms []dbschema.RoundModel

	err := QueryItemsFromDynamoDbWithFilter(constants.DB_TABLE_ROUNDS, &rms, map[string]interface{}{
		"Event_ID": evtGuid,
	}, map[string]interface{}{
		"Is_Active": true,
	})

	if err != nil {
		return nil, err
	}

	return rms, nil
}

func GetAllPlayersInEvent(evtGuid string) ([]dbschema.PlayerModel, error) {
	var sl []dbschema.PlayerModel

	err := QueryItemsFromDynamoDb(constants.DB_TABLE_PLAYERS, &sl, map[string]interface{}{
		"Event_ID": evtGuid,
	})

	if err != nil {
		return nil, err
	}

	return sl, nil
}

func GetAllResultsForRound(roundId string) ([]dbschema.ResultModel, error) {
	var rms []dbschema.ResultModel

	err := QueryItemsFromDynamoDb(constants.DB_TABLE_RESULTS, &rms, map[string]interface{}{
		"Round_ID": roundId,
	})

	if err != nil {
		return nil, err
	}

	return rms, nil
}

func GetAllResultsForRounds(roundIds ...string) ([]dbschema.ResultModel, error) {

	type threadResult struct {
		err error
		ls  []dbschema.ResultModel
	}

	chRes := make(chan threadResult)

	for _, v := range roundIds {
		go func(roundId string) {
			resInner, err := GetAllResultsForRound(roundId)
			chRes <- threadResult{err, resInner}
		}(v)
	}

	results := make([]dbschema.ResultModel, 0)
	errs := make([]error, 0)
	for i := 0; i < len(roundIds); i++ {
		r := <-chRes

		if r.err != nil {
			errs = append(errs, r.err)
		} else if r.ls != nil {
			results = append(results, r.ls...)
		}
	}

	if len(errs) > 0 {
		return nil, errors.NewComposite(errs...)
	}

	return results, nil
}

func AddPlayersToEvent(evtGuid string, players []dbschema.PlayerModel) error {

	// Ensure that the players are in this event
	newPlayers := make([]dbschema.PlayerModel, len(players))
	for i, v := range players {
		v.EventId = evtGuid
		newPlayers[i] = v
	}

	err := BatchPutItemsToDynamoDb(constants.DB_TABLE_PLAYERS, newPlayers)
	if err != nil {
		return err
	}

	return nil
}

func AddResultsToDynamoDb(results []dbschema.ResultModel) error {
	err := BatchPutItemsToDynamoDb(constants.DB_TABLE_RESULTS, results)
	if err != nil {
		return err
	}

	return nil
}

//func UpdateDatabaseWithScoreData(evtGuid string, data *schema.ScoreData) error {
//
//	// Start by updating the rounds
//
//	//
//	//	1. Add new rounds that didn't exist before
//	//
//
//	activeRounds, err := GetActiveRoundsForEvent(evtGuid)
//	if err != nil {
//		return err
//	}
//
//	// For any rounds in the schema.ScoreData not in the activeRounds list,
//	// add to DynamoDB
//	for _, ar := range activeRounds {
//
//		found := false
//		for _, d := range data.Rounds {
//			if d.RoundId == ar.RoundId {
//				found = true
//				break
//			}
//		}
//
//		if !found {
//			err = PutItemToDynamoDb(constants.DB_TABLE_ROUNDS, ar)
//			if err != nil {
//				return err
//			}
//		}
//	}
//
//	//
//	// 2. Add new players to database
//	//
//
//	// Extract players as *schema.Player
//	roundPlayers := make([]*schema.PlayerData, 0)
//	for _, ar := range data.Rounds {
//		for _, p := range ar.Players {
//
//			// Does this already exist?
//			found := false
//			for _, dp := range roundPlayers {
//				if dp.Name == p.Name {
//					found = true
//					break
//				}
//			}
//
//			if found {
//				continue
//			}
//
//			var newDataPlayer schema.PlayerData
//			newDataPlayer.Name = p.Name
//
//			roundPlayers = append(roundPlayers, &newDataPlayer)
//		}
//	}
//
//	databasePlayers, err := GetAllPlayersInEvent(evtGuid)
//	if err != nil {
//		return err
//	}
//
//	for _, rp := range roundPlayers {
//
//		// If the round player is not in the database, they need adding to the database
//		found := false
//		for _, ap := range databasePlayers {
//			if ap.Name == rp.Name {
//				found = true
//				break
//			}
//		}
//
//		if !found {
//			uid, err := uuid.NewRandom()
//			if err != nil {
//				return err
//			}
//
//			rp.PlayerId = uid.String()
//
//			err = PutItemToDynamoDb(constants.DB_TABLE_PLAYERS, rp)
//			if err != nil {
//				return err
//			}
//
//			databasePlayers = append(databasePlayers, rp)
//		}
//	}
//
//	//
//	// 3. Add results to database
//	//
//
//	resultList := make([]*schema.ResultData, 0)
//	for _, ar := range data.Rounds {
//		for _, p := range ar.Players {
//
//			var dp *schema.PlayerData
//			for _, dpi := range databasePlayers {
//				if dpi.Name == p.Name {
//					dp = dpi
//					break
//				}
//			}
//
//			if dp == nil {
//				log.Printf("WARN: Player not found: " + p.Name)
//				continue
//			}
//
//			var rd schema.ResultData
//			rd.RoundId = ar.RoundId
//			rd.PlayerId = dp.PlayerId
//			rd.Score = p.Score
//			rd.Rank = p.Rank
//			rd.Status = p.SessionStatus
//
//			resultList = append(resultList, &rd)
//		}
//	}
//
//	chErr := make(chan error)
//	for _, rl := range resultList {
//		go func(rlInner *schema.ResultData) {
//			err := UpdateItemInDynamoDb(constants.DB_TABLE_RESULTS, rlInner, "Round_ID", rlInner.RoundId, "Player_ID", rlInner.PlayerId)
//			chErr <- err
//		}(rl)
//	}
//
//	errs := make([]error, 0)
//	for i := 0; i < len(resultList); i++ {
//		err = <-chErr
//
//		if err != nil {
//			errs = append(errs, err)
//		}
//	}
//
//	if len(errs) > 0 {
//		err = errors.NewComposite(errs...)
//		return err
//	}
//
//	// Then update the event to change last_update to now
//	w, err := codezone_util.NewGenericWritable("Last_Update", time.Now())
//	if err != nil {
//		return err
//	}
//
//	err = UpdateItemInDynamoDb(constants.DB_TABLE_EVENTS, w, "Event_ID", evtGuid)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
