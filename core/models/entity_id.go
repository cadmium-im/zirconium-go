package models

import (
	"fmt"
	"strings"
)

type EntityIDType string

const (
	UsernameType  EntityIDType = "@"
	ChatAliasType EntityIDType = "#"
	ChatIDType    EntityIDType = "!"
)

type EntityID struct {
	EntityIDType EntityIDType
	LocalPart    string
	ServerPart   string
}

func NewEntityIDFromString(entityID string) (*EntityID, error) {
	eid := &EntityID{}
	typ := string(entityID[0])
	switch EntityIDType(typ) {
	case UsernameType:
		fallthrough
	case ChatAliasType:
		fallthrough
	case ChatIDType:
		{
			eid.EntityIDType = EntityIDType(typ)
		}
	default:
		return nil, fmt.Errorf("invalid entity id type: %s", typ)
	}

	localAndServerPart := strings.Split(entityID, "@")
	if len(localAndServerPart) == 3 {
		localAndServerPart = localAndServerPart[1:]
	}
	eid.LocalPart = localAndServerPart[0]
	eid.ServerPart = localAndServerPart[1]

	return eid, nil
}

func NewEntityID(typ, localPart, serverPart string) (*EntityID, error) {
	eid := &EntityID{}

	switch EntityIDType(typ) {
	case UsernameType:
		fallthrough
	case ChatAliasType:
		fallthrough
	case ChatIDType:
		{
			eid.EntityIDType = EntityIDType(typ)
		}
	default:
		return nil, fmt.Errorf("invalid entity id type: %s", typ)
	}

	eid.LocalPart = localPart
	eid.ServerPart = serverPart

	return eid, nil
}

func (eID *EntityID) String() string {
	return fmt.Sprintf("%s%s@%s", eID.EntityIDType, eID.LocalPart, eID.ServerPart)
}
