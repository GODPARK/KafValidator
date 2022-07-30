package util

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type MsgKey struct {
	key uuid.UUID
}

func NewUUID() *MsgKey {
	msgKey := &MsgKey{}

	myuuid := uuid.NewV4()
	msgKey.key = myuuid
	fmt.Printf("TEST KEY: %s\n", msgKey.KeyToString())
	return msgKey
}

func (msgKey *MsgKey) KeyToString() string {
	return msgKey.key.String()
}

func isKeyEqual(key1 string, key2 string) bool {
	return key1 == key2
}
