package models

type Exchange struct {
	ID     string
	Queues []*Queue
	Key    string
}
