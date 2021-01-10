package codezone_util

import (
	"codingame-live-scoreboard/constants"
	"codingame-live-scoreboard/schema"
	"codingame-live-scoreboard/schema/errors"
	"codingame-live-scoreboard/schema/shared_utils"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"log"
	"time"
)

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var dynamodbClient = dynamodb.New(sess)

// GetItemFromDynamoDb retrieves a single item from a given table based on key/value pairs given as variadic arguments.
// If no match is found, it will return nil.
func GetItemFromDynamoDb(tbl string, keyVals ...interface{}) (map[string]string, error) {
	processedKey, err := shared_utils.CreateKeyValuesFromList(keyVals)
	if err != nil {
		return nil, err
	}

	consistentRead := false
	gii := &dynamodb.GetItemInput{
		ConsistentRead: &consistentRead,
		Key:            processedKey,
		TableName:      &tbl,
	}

	res, err := dynamodbClient.GetItem(gii)
	if err != nil {
		return nil, err
	}

	retVal := make(map[string]string)
	if res.Item == nil {
		retVal = nil
	} else {
		for k, v := range res.Item {
			retVal[k] = v.String()
		}
	}

	return retVal, err
}

func BatchGetItemsFromDynamoDb(tbl string, keyVals ...interface{}) ([]map[string]string, error) {
	processedKey, err := shared_utils.CreateKeyValuesFromList(keyVals)
	if err != nil {
		return nil, err
	}

	consistentRead := false
	ri := make(map[string]*dynamodb.KeysAndAttributes)
	ri[tbl] = &dynamodb.KeysAndAttributes{
		ConsistentRead: &consistentRead,
		Keys:           []map[string]*dynamodb.AttributeValue{processedKey},
	}

	bgii := &dynamodb.BatchGetItemInput{
		RequestItems: ri,
	}

	res, err := dynamodbClient.BatchGetItem(bgii)
	if err != nil {
		return nil, err
	}

	ret := make([]map[string]string, 0)

	tblRes, ok := res.Responses[tbl]
	if !ok {
		return ret, nil
	}

	for _, v := range tblRes {
		retVal := make(map[string]string)
		for k, v2 := range v {
			retVal[k] = v2.String()
		}

		ret = append(ret, retVal)
	}

	return ret, nil
}

func PutItemToDynamoDb(tableName string, dbWritable schema.DynamoDbWritable) error {
	pii := &dynamodb.PutItemInput{
		Item:      dbWritable.ToDynamoDbMap(),
		TableName: &tableName,
	}

	_, err := dynamodbClient.PutItem(pii)
	if err != nil {
		return err
	}

	return nil
}

func UpdateItemInDynamoDb(tableName string, dbWritable schema.DynamoDbWritable, keyVals ...interface{}) error {
	keys, err := shared_utils.CreateKeyValuesFromList(keyVals)
	if err != nil {
		return err
	}

	attrValues := dbWritable.ToDynamoDbMap()
	attrValuesToWrite := make(map[string]*dynamodb.AttributeValue)

	for k, v := range attrValues {
		if _, ok := keys[k]; ok {
			continue
		}

		attrValuesToWrite[k] = v
	}

	uii := &dynamodb.UpdateItemInput{
		TableName:                 &tableName,
		Key:                       keys,
		ExpressionAttributeValues: attrValuesToWrite,
	}

	_, err = dynamodbClient.UpdateItem(uii)
	if err != nil {
		return err
	}

	return nil
}

func GetActiveRounds(evtGuid string) ([]*schema.RoundData, error) {
	activeRoundData, err := BatchGetItemsFromDynamoDb(constants.DB_TABLE_ROUNDS, "Event_ID", evtGuid)
	if err != nil {
		return nil, err
	}

	rdList := make([]*schema.RoundData, 0)

	for _, ard := range activeRoundData {
		var rd schema.RoundData
		rd.EventId = ard["Event_ID"]
		rd.RoundId = ard["Round_ID"]

		rdList = append(rdList, &rd)
	}

	return rdList, nil
}

func GetActivePlayers(evtGuid string) ([]*schema.PlayerData, error) {
	res, err := BatchGetItemsFromDynamoDb(constants.DB_TABLE_PLAYERS, "Event_ID", evtGuid)
	if err != nil {
		return nil, err
	}

	ret := make([]*schema.PlayerData, 0)
	for _, r := range res {
		var retVal schema.PlayerData
		retVal.Name = r["Name"]
		retVal.PlayerId = r["Player_ID"]

		ret = append(ret, &retVal)
	}

	return ret, nil
}

func UpdateDatabaseWithScoreData(evtGuid string, data *schema.ScoreData) error {

	// Start by updating the rounds

	//
	//	1. Add new rounds that didn't exist before
	//

	activeRounds, err := GetActiveRounds(evtGuid)
	if err != nil {
		return err
	}

	// For any rounds in the schema.ScoreData not in the activeRounds list,
	// add to DynamoDB
	for _, ar := range activeRounds {

		found := false
		for _, d := range data.ActiveRounds {
			if d.RoundId == ar.RoundId {
				found = true
				break
			}
		}

		if !found {
			err = PutItemToDynamoDb(constants.DB_TABLE_ROUNDS, ar)
			if err != nil {
				return err
			}
		}
	}

	//
	// 2. Add new players to database
	//

	// Extract players as *schema.Player
	roundPlayers := make([]*schema.PlayerData, 0)
	for _, ar := range data.ActiveRounds {
		for _, p := range ar.Players {

			// Does this already exist?
			found := false
			for _, dp := range roundPlayers {
				if dp.Name == p.Name {
					found = true
					break
				}
			}

			if found {
				continue
			}

			var newDataPlayer schema.PlayerData
			newDataPlayer.Name = p.Name

			roundPlayers = append(roundPlayers, &newDataPlayer)
		}
	}

	databasePlayers, err := GetActivePlayers(evtGuid)
	if err != nil {
		return err
	}

	for _, rp := range roundPlayers {

		// If the round player is not in the database, they need adding to the database
		found := false
		for _, ap := range databasePlayers {
			if ap.Name == rp.Name {
				found = true
				break
			}
		}

		if !found {
			uid, err := uuid.NewRandom()
			if err != nil {
				return err
			}

			rp.PlayerId = uid.String()

			err = PutItemToDynamoDb(constants.DB_TABLE_PLAYERS, rp)
			if err != nil {
				return err
			}

			databasePlayers = append(databasePlayers, rp)
		}
	}

	//
	// 3. Add results to database
	//

	resultList := make([]*schema.ResultData, 0)
	for _, ar := range data.ActiveRounds {
		for _, p := range ar.Players {

			var dp *schema.PlayerData
			for _, dpi := range databasePlayers {
				if dpi.Name == p.Name {
					dp = dpi
					break
				}
			}

			if dp == nil {
				log.Printf("WARN: Player not found: " + p.Name)
				continue
			}

			var rd schema.ResultData
			rd.RoundId = ar.RoundId
			rd.PlayerId = dp.PlayerId
			rd.Score = p.Score
			rd.Rank = p.Rank
			rd.Status = p.SessionStatus

			resultList = append(resultList, &rd)
		}
	}

	chErr := make(chan error)
	for _, rl := range resultList {
		go func(rlInner *schema.ResultData) {
			err := UpdateItemInDynamoDb(constants.DB_TABLE_RESULTS, rlInner, "Round_ID", rlInner.RoundId, "Player_ID", rlInner.PlayerId)
			chErr <- err
		}(rl)
	}

	errs := make([]error, 0)
	for i := 0; i < len(resultList); i++ {
		err = <-chErr

		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		err = errors.NewComposite(errs...)
		return err
	}

	// Then update the event to change last_update to now
	w, err := NewGenericWritable("Last_Update", time.Now())
	if err != nil {
		return err
	}

	err = UpdateItemInDynamoDb(constants.DB_TABLE_EVENTS, w, "Event_ID", evtGuid)
	if err != nil {
		return err
	}

	return nil
}
