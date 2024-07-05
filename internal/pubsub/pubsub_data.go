package pubsub

type SimpleQueueType int

const (
	Durable SimpleQueueType = iota + 1
	Transient
)

type AckType int

const (
	Ack AckType = iota + 1
	NackRequeue
	NackDiscard
)
