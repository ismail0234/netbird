package server

import (
	//"log"
	"time"

	nbpeer "github.com/netbirdio/netbird/management/server/peer"
	// "gorm.io/gorm"
)

type PeerChildren struct {
	*nbpeer.Peer
}

type PeerStatusChildren struct {
	*nbpeer.PeerStatus
}

func GetDefaultTimezone() time.Time {
	return time.Date(1, 1, 1, 1, 1, 1, 0, time.Local)
}

