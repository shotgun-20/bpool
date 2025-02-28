package main

import "errors"

/*
Simple thread-safe pool of byte buffers, with fixed pool capacity.
*/
type BPool struct {
	cap  int         // number of buffers
	size int         // size of a single buffer
	pool chan []byte // the pool itself
}

/*
Returns newly created BufferPool of "cap" buffers, each of size "size".
*/
func NewBPool(size, cap int) *BPool {
	pool := make(chan []byte, cap)
	for range cap {
		pool <- make([]byte, size)
	}
	bp := &BPool{cap: cap, size: size, pool: pool}
	return bp
}

/*
Returns a buffer from pool, if available. Blocking.
*/
func (bp *BPool) Get() []byte {
	return <-bp.pool
}

/*
Hands over an outbound channel of []byte buffers
*/
func (bp *BPool) GetC() <-chan []byte {
	return bp.pool
}

/*
Non-blocking version of get
*/
func (bp *BPool) GetNB() ([]byte, error) {
	select {
	case b, ok := <-bp.pool:
		if !ok {
			return nil, errors.New("pool died")
		}
		return b, nil
	default:
		return nil, errors.New("no buffer available in pool")
	}
}

/*
Returns buffer to the pool when it is not needed anymore.
Do not use buffer after it was recycled! Get() new buffer instead.
*/
func (bp *BPool) Recycle(b []byte) error {
	select {
	case bp.pool <- b[:]:
		return nil
	default:
		return errors.New("cannot recycle, pool is already at full cap")
	}
}
