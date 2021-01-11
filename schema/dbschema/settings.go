package dbschema

type SettingModel struct {
	SettingName  string `ddb:"Setting_Name,key"`
	SettingValue string `ddb:"Value"`
	TimeToLive   int    `ddb:"TTL"`
}
