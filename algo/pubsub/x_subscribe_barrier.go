package pubsub

type xSubscribeBarrier[T any] struct {
	subscriberCapacity uint64
	subscribers        []Subscriber[T]
}
