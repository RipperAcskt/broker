package models

type Queue struct {
	ID            string
	Messages      []*Message
	UserConnected bool
}
