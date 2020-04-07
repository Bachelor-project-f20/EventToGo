package eventToGo

type BrokerType int

const (
	NATS BrokerType = iota
	SNS
)
