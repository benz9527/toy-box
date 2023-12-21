package pubsub

type xPublishBarrier[T any] struct {
	capacity uint64
	cache    []T
}
