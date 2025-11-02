package config

import (
	"log"
	"os"

	"github.com/gorcon/rcon"
)

func InitRCON() (*rcon.Conn) {
	conn, err := rcon.Dial(os.Getenv("RCON_SRV"), os.Getenv("RCON_PASSWORD"))
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
