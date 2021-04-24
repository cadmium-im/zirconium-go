package models

import "testing"

func TestNewEntityIDFromString(t *testing.T) {
	str := "@user@cadmium.org"

	eid, err := NewEntityIDFromString(str)
	if err != nil {
		t.Fatal("error must be null")
	}
	if eid.Attr != "" || eid.Type != UsernameType || eid.ServerPart != "cadmium.org" || eid.LocalPart != "user" {
		t.FailNow()
	}
}

func TestNewEntityIDFromStringWithAttr(t *testing.T) {
	str := "%msisdn:18002003040@cadmium.org"

	eid, err := NewEntityIDFromString(str)
	if err != nil {
		t.Fatal("error must be null")
	}
	if eid.Attr != "msisdn" || eid.Type != ThirdPIDType || eid.ServerPart != "cadmium.org" || eid.LocalPart != "18002003040" {
		t.FailNow()
	}
}

func TestNewEntityIDFromStringWithEmailAttr(t *testing.T) {
	str := "%email:abslimit_netwhood.online@cadmium.org"

	eid, err := NewEntityIDFromString(str)
	if err != nil {
		t.Fatal("error must be null")
	}
	if eid.Attr != "email" || eid.Type != ThirdPIDType || eid.ServerPart != "cadmium.org" || eid.LocalPart != "abslimit_netwhood.online" {
		t.Fatal(eid.String())
	}
}

func TestNewEntityIDFromStringWithOnlyServerPart(t *testing.T) {
	str := "cadmium.org"

	eid, err := NewEntityIDFromString(str)
	if err != nil {
		t.Fatal("error must be null")
	}
	if !eid.OnlyServerPart && eid.ServerPart != "cadmium.org" {
		t.Fatal(eid.String())
	}
}
