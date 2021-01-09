package codezone_util

import (
	"codezone-codingame-live-scoreboard/constants"
	"codezone-codingame-live-scoreboard/dbutils"
	"codezone-codingame-live-scoreboard/schema"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"strconv"
)

var sess = session.Must(session.NewSession())
var dynamodbClient = dynamodb.New(sess)

func createProcessedKey(keyVals []interface{}) (map[string]*dynamodb.AttributeValue, error) {
	if len(keyVals) % 2 != 0 {
		return nil, errors.New("keyVals must be provided as pairs of key/value")
	}

	processedKey := make(map[string]*dynamodb.AttributeValue)

	keyMap := make(map[string]interface{})
	for i := 0; i < len(keyVals); i += 2 {
		switch kt := keyVals[i].(type) {
		case string:
			keyMap[kt] = keyVals[i + 1]
		default:
			return nil, errors.New(fmt.Sprintf("key value at position %v is not a string", i))
		}
	}

	for k, v := range keyMap {
		var a dynamodb.AttributeValue

		switch vt := v.(type) {
		case string:
			a.SetS(vt)
			break
		case bool:
			a.SetBOOL(vt)
			break
		case int:
			a.SetS(strconv.Itoa(vt))
			break
		default:
			return nil, errors.New("unsupported partKey type")
		}

		processedKey[k] = &a
	}

	return processedKey, nil
}

func GetItemFromDynamoDb(tbl string, keyVals ...interface{}) (map[string]string, error) {
	processedKey, err := createProcessedKey(keyVals)
	if err != nil {
		return nil, err
	}

	consistentRead := false
	gii := &dynamodb.GetItemInput {
		ConsistentRead: &consistentRead,
		Key: processedKey,
		TableName: &tbl,
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
	processedKey, err := createProcessedKey(keyVals)
	if err != nil {
		return nil, err
	}

	consistentRead := false
	ri := make(map[string]*dynamodb.KeysAndAttributes)
	ri[tbl] = &dynamodb.KeysAndAttributes{
		ConsistentRead: &consistentRead,
		Keys: []map[string]*dynamodb.AttributeValue{processedKey},
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
		Item: dbWritable.ToDynamoDbMap(),
		TableName: &tableName,
	}

	_, err := dynamodbClient.PutItem(pii)
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

	for _, dp := range roundPlayers {

		// If the round player is not in the database, they need adding to the database
		found := false
		for _, ap := range databasePlayers {
			if ap.Name == dp.Name {
				found = true
				break
			}
		}

		if !found {
			uid, err := uuid.NewRandom()
			if err != nil {
				return err
			}

			dp.PlayerId = uid.String()

			err = PutItemToDynamoDb(constants.DB_TABLE_PLAYERS, dp)
			if err != nil {
				return err
			}
		}
	}

	//
	// 3. Add results to database
	//

	// Then update the event to change last_update to now

}
