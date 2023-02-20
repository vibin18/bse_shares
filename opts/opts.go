package opts

import (
	"encoding/json"
	"log"
)

type Params struct {
	DbServer string `           long:"dbserver"      env:"DB_SERVER"  description:"Server name or IP of postgres db" default:"192.168.1.20"`
	DbPort   string `           long:"dbport"      env:"DB_PORT"  description:"Port number of db server" default:"5432"`
	DbName   string `           long:"dbname"      env:"DB_Name"  description:"Name of the db" default:"shares"`
	DbUser   string `           long:"dbuser"      env:"DB_User"  description:"Username of the db" default:"postgres"`
	DbPass   string `           long:"dbpass"      env:"DB_PASS"  description:"Password of the db"`
}

func (o *Params) GetJson() []byte {
	jsonBytes, err := json.Marshal(o)
	if err != nil {
		log.Panic(err)
	}
	return jsonBytes
}
