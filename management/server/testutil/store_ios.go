//go:build ios
// +build ios

package testutil

func CreatePGDB() (func(), error) { return func() {}, nil }

func CreateMyDB() (func(), error) { return func() {}, nil }
