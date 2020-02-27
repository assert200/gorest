package gorest

import (
	"gotest.tools/assert"
	"testing"
	"time"
)

func TestChannelReturnsNilOnTimeout(t *testing.T) {

	rt := NewRestTestChannel(1 * time.Second)

	result := rt.Read()
	assert.Assert(t, result == nil, "result is not nil")
}

func TestChannelReturnsOnRange(t *testing.T) {

	rt := NewRestTestChannel(1 * time.Second)

	rt.Write(&RestTest{})

	count := 0
	for rt.Next() {
		_ = rt.Read()
		count++
	}

	assert.Assert(t, count == 1, "count was not 1")

}