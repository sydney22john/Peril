package pubsub

type SimpleQueueType int

const (
	Durable = iota + 1
	Transient
)
