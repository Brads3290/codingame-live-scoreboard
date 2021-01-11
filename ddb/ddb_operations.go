package ddb

import (
	"codingame-live-scoreboard/constants"
	"codingame-live-scoreboard/schema/dbschema"
)

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
//		for _, d := range data.ActiveRounds {
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
//	for _, ar := range data.ActiveRounds {
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
//	for _, ar := range data.ActiveRounds {
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
