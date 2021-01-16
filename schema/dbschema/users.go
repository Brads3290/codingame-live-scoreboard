package dbschema

import "time"

type UserModel struct {
	UserId       string    `ddb:"User_ID,key"`
	Username     string    `ddb:"Username"`
	PasswordHash string    `ddb:"Password_Hash"`
	PasswordSalt string    `ddb:"Password_Salt"`
	DateCreated  time.Time `ddb:"Date_Created"`
	LastLogin    time.Time `ddb:"Last_Login"`
}
