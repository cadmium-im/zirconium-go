package models

import (
	"fmt"
	"strings"
)

type EntityIDType string

const (
	UsernameType  EntityIDType = "@"
	RoomAliasType EntityIDType = "#"
	RoomIDType    EntityIDType = "!"
)

type EntityID struct {
	EntityIDType EntityIDType
	LocalPart    string
	ServerPart   string
}

func NewEntityID(entityID string) *EntityID {
	eID := &EntityID{}
	switch EntityIDType(string(entityID[0])) {
	case UsernameType:
		{
			eID.EntityIDType = UsernameType
		}
	case RoomAliasType:
		{
			eID.EntityIDType = RoomAliasType
		}
	case RoomIDType:
		{
			eID.EntityIDType = RoomIDType
		}
	}
	localAndServerPart := strings.Split(entityID, "@")
	if len(localAndServerPart) == 3 {
		localAndServerPart = localAndServerPart[1:]
	}
	eID.LocalPart = localAndServerPart[0]
	eID.ServerPart = localAndServerPart[1]

	return eID
}

func (eID *EntityID) String() string {
	return fmt.Sprintf("%s%s@%s", eID.EntityIDType, eID.LocalPart, eID.ServerPart)
}
