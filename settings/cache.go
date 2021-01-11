package settings

import (
	"codingame-live-scoreboard/constants"
	"codingame-live-scoreboard/ddb"
	"codingame-live-scoreboard/schema/dbschema"
	"codingame-live-scoreboard/schema/errors"
	"github.com/sirupsen/logrus"
	"time"
)

type memCacheEntry struct {
	val     string
	expires time.Time
}

var memCache = make(map[string]memCacheEntry)
var logger = logrus.WithField(constants.PROGRAM_CONTEXT_FIELD, "SettingsCache")

func cacheGet(settingName string) (string, error) {

	// 1. Try the memory cache
	if mcVal, ok := memCache[settingName]; ok {

		// We have something, check if TTL is OK
		if mcVal.expires.After(time.Now()) {

			// Expiry is in the future; OK
			return mcVal.val, nil
		} else {

			// This has expired; remove it
			delete(memCache, settingName)
		}
	}

	// 2. If not present, or TTL expired, get from DB cache; we retrieve all settings
	// in bulk here, to refresh all of them
	var settingsList []dbschema.SettingModel
	err := ddb.ScanItemsFromDynamoDb(constants.DB_TABLE_SETTINGS, &settingsList)
	if err == nil {
		var result *string

		// Update the memory cache
		for _, sl := range settingsList {
			entry := memCacheEntry{
				val:     sl.SettingValue,
				expires: time.Now().Add(time.Duration(sl.TimeToLive) * time.Second),
			}

			memCache[sl.SettingName] = entry

			if sl.SettingName == settingName {
				result = &entry.val
			}
		}

		// Did we find a match?
		if result != nil {
			return *result, nil
		}
	} else {
		logger.Error("Failed to scan settings items from dynamodb: ", err)
	}

	// 3. Fall back on default
	if defVal, ok := settingDefaults[settingName]; ok {
		logger.Warn("Falling back on default value for setting: ", settingName)
		return defVal, nil
	}

	return "", errors.New("setting with name '" + settingName + "' could not be found.")
}
