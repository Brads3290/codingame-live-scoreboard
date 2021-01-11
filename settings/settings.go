package settings

import "strconv"

func GetString(settingName string) string {
	res, err := cacheGet(settingName)
	if err != nil {
		panic(err)
	}

	return res
}

func GetInt(settingName string) int {
	s := GetString(settingName)
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return i
}
