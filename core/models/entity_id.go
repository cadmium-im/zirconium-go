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
	ThirdPIDType  EntityIDType = "%"
	OtherType     EntityIDType = "&"
)

type EntityID struct {
	Type           EntityIDType
	LocalPart      string
	ServerPart     string
	Attr           string
	OnlyServerPart bool
}

func NewEntityIDFromString(entityID string) (*EntityID, error) {
	eid := &EntityID{}
	typ := string(entityID[0])
	withAttr := false
	localAndServerPart := strings.Split(entityID, "@")
	if localAndServerPart[0] == "" {
		localAndServerPart = localAndServerPart[1:]
	}
	if len(localAndServerPart) == 1 {
		eid.ServerPart = localAndServerPart[0]
		eid.OnlyServerPart = true
		return eid, nil
	}

	switch EntityIDType(typ) {
	case UsernameType:
		fallthrough
	case ChatAliasType:
		fallthrough
	case ChatIDType:
		{
			eid.Type = EntityIDType(typ)
		}
	case ThirdPIDType:
		fallthrough
	case OtherType:
		{
			eid.Type = EntityIDType(typ)
			withAttr = true
		}
	default:
		return nil, fmt.Errorf("invalid entity id type: %s", typ)
	}

	if !withAttr {
		eid.LocalPart = localAndServerPart[0]
		eid.ServerPart = localAndServerPart[1]
	} else {
		attrAndLocal := strings.Split(localAndServerPart[0], ":")
		attr := attrAndLocal[0][1:]
		eid.Attr = attr
		eid.LocalPart = attrAndLocal[1]
		eid.ServerPart = localAndServerPart[len(localAndServerPart)-1]
	}

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
			eid.Type = EntityIDType(typ)
		}
	default:
		return nil, fmt.Errorf("invalid entity id type: %s", typ)
	}

	eid.LocalPart = localPart
	eid.ServerPart = serverPart

	return eid, nil
}

func NewEntityIDWithAttr(typ, attr, localPart, serverPart string) (*EntityID, error) {
	eid := &EntityID{}

	switch EntityIDType(typ) {
	case OtherType:
		fallthrough
	case ThirdPIDType:
		{
			eid.Type = EntityIDType(typ)
		}
	default:
		return nil, fmt.Errorf("invalid entity id type: %s", typ)
	}

	eid.Attr = attr
	eid.LocalPart = localPart
	eid.ServerPart = serverPart

	return eid, nil
}

func (eID *EntityID) String() string {
	if eID.Attr != "" {
		return fmt.Sprintf("%s%s:%s@%s", eID.Type, eID.Attr, eID.LocalPart, eID.ServerPart)
	}
	return fmt.Sprintf("%s%s@%s", eID.Type, eID.LocalPart, eID.ServerPart)
}
