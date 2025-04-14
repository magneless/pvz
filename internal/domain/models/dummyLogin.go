package models

type Role string

func (r Role) IsValid() bool {
	return r == "moderator" || r == "employee"
}

type DummyLogin struct {
	Role Role
}