package server

import (
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
	return time.Date(1, 1, 1, 1, 1, 1, 1, time.Local)
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

	if row.LastLogin.IsZero() {
		row.LastLogin = GetDefaultTimezone()
	}

	return nil
}

func (row *UserInfo) BeforeSave(tx *gorm.DB) (err error) {

	if row.LastLogin.IsZero() {
		row.LastLogin = GetDefaultTimezone()
	}

	return nil
}

func (row *PeerChildren) BeforeSave(tx *gorm.DB) (err error) {

	if row.CreatedAt.IsZero() {
		row.CreatedAt = GetDefaultTimezone()
	}

	if row.LastLogin.IsZero() {
		row.LastLogin = GetDefaultTimezone()
	}

	return nil
}

func (row *PeerStatusChildren) BeforeSave(tx *gorm.DB) (err error) {

	if row.LastSeen.IsZero() {
		row.LastSeen = GetDefaultTimezone()
	}

	return nil
}
