package server

import (
	"log"
	"time"

	nbpeer "github.com/netbirdio/netbird/management/server/peer"
	"gorm.io/gorm"
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

func (row *SetupKey) BeforeSave(tx *gorm.DB) (err error) {

	if row.CreatedAt.IsZero() {
		row.CreatedAt = GetDefaultTimezone()
	}

	if row.ExpiresAt.IsZero() {
		row.ExpiresAt = GetDefaultTimezone()
	}

	if row.UpdatedAt.IsZero() {
		row.UpdatedAt = GetDefaultTimezone()
	}

	if row.LastUsed.IsZero() {
		row.LastUsed = GetDefaultTimezone()
	}

	return nil
}
func (row *Account) BeforeSave(tx *gorm.DB) (err error) {

	if row.CreatedAt.IsZero() {
		row.CreatedAt = GetDefaultTimezone()
	}

	// maybe pointer*
	for _, key := range row.Peers {
		if key.LastLogin.IsZero() {
			key.LastLogin = GetDefaultTimezone()
		}
	}
	for _, key := range row.PeersG {
		if key.LastLogin.IsZero() {
			key.LastLogin = GetDefaultTimezone()
		}
	}

	for _, key := range row.UsersG {
		if key.LastLogin.IsZero() {
			key.LastLogin = GetDefaultTimezone()
		}
	}

	for _, key := range row.Users {
		if key.LastLogin.IsZero() {
			key.LastLogin = GetDefaultTimezone()
		}

		for _, pat := range key.PATs {
			if pat.LastUsed.IsZero() {
				pat.LastUsed = GetDefaultTimezone()
			}
		}
	}

	for _, key := range row.SetupKeys {
		if key.LastUsed.IsZero() {
			key.LastUsed = GetDefaultTimezone()
		}
	}

	for _, key := range row.SetupKeysG {
		if key.LastUsed.IsZero() {
			key.LastUsed = GetDefaultTimezone()
		}
	}

	return nil
}

func (row *PersonalAccessToken) BeforeSave(tx *gorm.DB) (err error) {

	if row.CreatedAt.IsZero() {
		row.CreatedAt = GetDefaultTimezone()
	}

	if row.ExpirationDate.IsZero() {
		row.ExpirationDate = GetDefaultTimezone()
	}

	if row.LastUsed.IsZero() {
		row.LastUsed = GetDefaultTimezone()
	}

	return nil
}

func (row *User) BeforeSave(tx *gorm.DB) (err error) {

	if row.CreatedAt.IsZero() {
		row.CreatedAt = GetDefaultTimezone()
	}

	for _, pat := range row.PATs {
		if pat.LastUsed.IsZero() {
			pat.LastUsed = GetDefaultTimezone()
		}
	}

	for _, pat := range row.PATsG {
		if pat.LastUsed.IsZero() {
			pat.LastUsed = GetDefaultTimezone()
		}
	}

	if row.LastLogin.IsZero() {
		log.Printf("row.LastLogin - User - 1: %s", row.LastLogin)

		row.LastLogin = GetDefaultTimezone()

		log.Printf("row.LastLogin - User - 2: %s", row.LastLogin)
	} else {
		log.Printf("row.LastLogin - User - 3: %s", row.LastLogin)
	}

	return nil
}

func (row *UserInfo) BeforeSave(tx *gorm.DB) (err error) {

	log.Printf("row.LastLogin - UserInfo - 0: %s", row.LastLogin)

	if row.LastLogin.IsZero() {
		log.Printf("row.LastLogin - UserInfo - 1: %s", row.LastLogin)

		row.LastLogin = GetDefaultTimezone()

		log.Printf("row.LastLogin - UserInfo - 2: %s", row.LastLogin)
	} else {
		log.Printf("row.LastLogin - UserInfo - 3: %s", row.LastLogin)
	}

	return nil
}

func (row *PeerChildren) BeforeSave(tx *gorm.DB) (err error) {

	if row.CreatedAt.IsZero() {
		row.CreatedAt = GetDefaultTimezone()
	}

	if row.LastLogin.IsZero() {
		log.Printf("row.LastLogin - PeerChildren - 1: %s", row.LastLogin)

		row.LastLogin = GetDefaultTimezone()

		log.Printf("row.LastLogin - PeerChildren - 2: %s", row.LastLogin)
	} else {
		log.Printf("row.LastLogin - PeerChildren - 3: %s", row.LastLogin)
	}

	return nil
}

func (row *PeerStatusChildren) BeforeSave(tx *gorm.DB) (err error) {

	if row.LastSeen.IsZero() {
		row.LastSeen = GetDefaultTimezone()
	}

	return nil
}

func (row *UserInfo) BeforeCreate(tx *gorm.DB) (err error) {

	log.Printf("row.LastLogin - BeforeCreate - UserInfo - 0: %s", row.LastLogin)

	if row.LastLogin.IsZero() {
		log.Printf("row.LastLogin - BeforeCreate - UserInfo - 1: %s", row.LastLogin)

		row.LastLogin = GetDefaultTimezone()

		log.Printf("row.LastLogin - BeforeCreate -UserInfo - 2: %s", row.LastLogin)
	} else {
		log.Printf("row.LastLogin - BeforeCreate -UserInfo - 3: %s", row.LastLogin)
	}

	return nil
}
func (row *PeerChildren) BeforeCreate(tx *gorm.DB) (err error) {

	log.Printf("row.LastLogin - BeforeCreate - PeerChildren - 0: %s", row.LastLogin)

	if row.LastLogin.IsZero() {
		log.Printf("row.LastLogin - BeforeCreate - PeerChildren - 1: %s", row.LastLogin)

		row.LastLogin = GetDefaultTimezone()

		log.Printf("row.LastLogin - BeforeCreate -PeerChildren - 2: %s", row.LastLogin)
	} else {
		log.Printf("row.LastLogin - BeforeCreate -PeerChildren - 3: %s", row.LastLogin)
	}

	return nil
}
func (row *User) BeforeCreate(tx *gorm.DB) (err error) {

	log.Printf("row.LastLogin - BeforeCreate - LastLogin - 0: %s", row.LastLogin)

	if row.LastLogin.IsZero() {
		log.Printf("row.LastLogin - BeforeCreate - LastLogin - 1: %s", row.LastLogin)

		row.LastLogin = GetDefaultTimezone()

		log.Printf("row.LastLogin - BeforeCreate - LastLogin - 2: %s", row.LastLogin)
	} else {
		log.Printf("row.LastLogin - BeforeCreate - LastLogin - 3: %s", row.LastLogin)
	}

	return nil
}

func (row *UserInfo) BeforeUpdate(tx *gorm.DB) (err error) {

	log.Printf("row.LastLogin - BeforeUpdate - UserInfo - 0: %s", row.LastLogin)

	if row.LastLogin.IsZero() {
		log.Printf("row.LastLogin - BeforeUpdate - UserInfo - 1: %s", row.LastLogin)

		row.LastLogin = GetDefaultTimezone()

		log.Printf("row.LastLogin - BeforeUpdate -UserInfo - 2: %s", row.LastLogin)
	} else {
		log.Printf("row.LastLogin - BeforeUpdate -UserInfo - 3: %s", row.LastLogin)
	}

	return nil
}
func (row *PeerChildren) BeforeUpdate(tx *gorm.DB) (err error) {

	log.Printf("row.LastLogin - BeforeUpdate - PeerChildren - 0: %s", row.LastLogin)

	if row.LastLogin.IsZero() {
		log.Printf("row.LastLogin - BeforeUpdate - PeerChildren - 1: %s", row.LastLogin)

		row.LastLogin = GetDefaultTimezone()

		log.Printf("row.LastLogin - BeforeUpdate -PeerChildren - 2: %s", row.LastLogin)
	} else {
		log.Printf("row.LastLogin - BeforeUpdate -PeerChildren - 3: %s", row.LastLogin)
	}

	return nil
}
func (row *User) BeforeUpdate(tx *gorm.DB) (err error) {

	log.Printf("row.LastLogin - BeforeUpdate - LastLogin - 0: %s", row.LastLogin)

	if row.LastLogin.IsZero() {
		log.Printf("row.LastLogin - BeforeUpdate - LastLogin - 1: %s", row.LastLogin)

		row.LastLogin = GetDefaultTimezone()

		log.Printf("row.LastLogin - BeforeUpdate - LastLogin - 2: %s", row.LastLogin)
	} else {
		log.Printf("row.LastLogin - BeforeUpdate - LastLogin - 3: %s", row.LastLogin)
	}

	return nil
}
