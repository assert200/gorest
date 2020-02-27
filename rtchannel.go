package gorest

import "time"

type RestTestChannel struct {
	ch chan *RestTest
	timeout time.Duration
}

func NewRestTestChannel(d time.Duration) *RestTestChannel {
	rt := &RestTestChannel{
		timeout: d,
	}
	rt.ch = make(chan *RestTest, 100)
	return rt
}

func (rt RestTestChannel) Read() *RestTest {

	select {
	case test := <-rt.ch:
		return test
	case <-time.After(rt.timeout):
		return nil
	}
}

func (rt RestTestChannel) Next() bool {
	return len(rt.ch) > 0
}

func (rt *RestTestChannel) Write(test *RestTest) {
	 rt.ch <- test
}

func (rt *RestTestChannel) Close() {
	close(rt.ch)
}
