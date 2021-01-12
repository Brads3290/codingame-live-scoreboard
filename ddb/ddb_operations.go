package ddb

import (
	"codingame-live-scoreboard/constants"
	"codingame-live-scoreboard/schema/dbschema"
	"codingame-live-scoreboard/schema/errors"
	"time"
)

func GetEvents() ([]dbschema.EventModel, error) {
	var eventsList []dbschema.EventModel
	err := ScanItemsFromDynamoDb(constants.DB_TABLE_EVENTS, &eventsList)

	if err != nil {
		return nil, err
	}

	return eventsList, nil
}

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

func MarkEventUpdatedNow(evtGuid string) error {
	err := UpdateItemAttrsInDynamoDb(constants.DB_TABLE_EVENTS, map[string]interface{}{
		"Event_ID": evtGuid,
	}, map[string]interface{}{
		"Last_Updated": time.Now(),
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

func DeleteRoundRecords(results []dbschema.RoundModel) error {
	err := BatchDeleteItemsFromDynamoDb(constants.DB_TABLE_ROUNDS, results)
	if err != nil {
		return err
	}

	return nil
}

func DeleteResultRecords(results []dbschema.ResultModel) error {
	err := BatchDeleteItemsFromDynamoDb(constants.DB_TABLE_RESULTS, results)
	if err != nil {
		return err
	}

	return nil
}

func DeletePlayerRecords(players []dbschema.PlayerModel) error {
	err := BatchDeleteItemsFromDynamoDb(constants.DB_TABLE_PLAYERS, players)
	if err != nil {
		return err
	}

	return nil
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
