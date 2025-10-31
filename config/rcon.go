package config

import (
	"log"
	"os"

	"github.com/gorcon/rcon"
)

func InitRCON() {
	conn, err := rcon.Dial(os.Getenv("RCON_SRV"), os.Getenv("RCON_PORT"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
}
