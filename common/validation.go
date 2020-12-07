package common

import (
	"log"
	"strconv"
)

func IsValidPort(port string) bool {
	pInt, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		log.Printf("info: failed to parse port string=%s. err=%s\n", port, err)
		return false
	}

	if pInt < 1024 {
		log.Printf("info: invalid port range. should be in between 1024 and 65535, but provided %s\n", port)
		return false
	}

	return true
}
