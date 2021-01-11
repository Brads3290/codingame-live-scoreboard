package codingame

import (
	"bytes"
	"codingame-live-scoreboard/constants"
	"codingame-live-scoreboard/ddb"
	"codingame-live-scoreboard/schema"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"path"
	"time"
)

var c = http.Client{Timeout: 30 * time.Second}

func GetCodinGameData(evtGuid string) (*schema.ScoreData, error) {
	activeRounds, err := ddb.GetActiveRoundsForEvent(evtGuid)
	if err != nil {
		return nil, err
	}

	resultChan := make(chan interface{})
	for _, ar := range activeRounds {

		// TODO: Pass context to thread and implement cancellation
		go getCodinGameDataOnThread(ar.RoundId, resultChan)
	}

	roundData := make([]schema.RoundData, 0)

	threadsCreated := len(activeRounds)
	for i := 0; i < threadsCreated; i++ {
		select {
		case res := <-resultChan:
			switch tres := res.(type) {
			case error:
				return nil, tres
			case schema.RoundData:
				tres.EventId = evtGuid
				roundData = append(roundData, tres)
			default:
				return nil, errors.New("invalid type returned from codingame thread")
			}
		case <-time.After(35 * time.Second):
			return nil, errors.New("timeout after 35 seconds")
		}
	}

	var sd schema.ScoreData
	sd.EventId = evtGuid
	sd.ActiveRounds = roundData

	return &sd, nil
}

func getCodinGameDataOnThread(roundId string, ce chan interface{}) {
	httpUrl := path.Join(constants.CODINGAME_BASE_URL, constants.CODINGAME_CLASHREPORT_PATH)

	jsonData := []string{roundId}
	b, err := json.Marshal(jsonData)
	if err != nil {
		ce <- err
		return
	}

	res, err := c.Post(httpUrl, "application/json; charset=utf-8", bytes.NewReader(b))
	if err != nil {
		ce <- err
		return
	}

	defer res.Body.Close()

	br, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ce <- err
		return
	}

	var resStruct schema.CodinGameClashReportResponse
	err = json.Unmarshal(br, &resStruct)

	var retStruct schema.RoundData
	retStruct.RoundId = resStruct.PublicHandle
	retStruct.Mode = resStruct.Mode
	retStruct.Players = make([]schema.PlayerRoundData, 0)

	for _, v := range resStruct.Players {
		var retPlayer schema.PlayerRoundData
		retPlayer.Name = v.Nickname
		retPlayer.Rank = v.Rank
		retPlayer.Score = v.RoundScore
		retPlayer.SessionStatus = v.SessionStatus

		retStruct.Players = append(retStruct.Players, retPlayer)
	}

	ce <- retStruct
}
