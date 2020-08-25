package internal_test

import (
	"fmt"
	"testing"

	"github.com/jjeffcaii/rattle/internal"
	"github.com/stretchr/testify/assert"
)

func TestNewPathTrie(t *testing.T) {
	pt := internal.NewPathTrie()
	err := pt.AddPath("/users/:uid/orders/:orderId", 31)
	assert.NoError(t, err)
	n, ok := pt.Load("/users/7777/orders/8888")
	assert.True(t, ok)
	assert.True(t, n.IsLeaf())
	v, ok := n.Value()
	assert.True(t, ok)
	fmt.Println("value:", v)
	va, v, ok := pt.Find("/users/7777/orders/8888")
	assert.Equal(t, "7777", va.GetOrDefault("uid", ""))
	assert.Equal(t, "8888", va.GetOrDefault("orderId", ""))

	_, ok = pt.Load("//users/")
	assert.False(t, ok)
	n, ok = pt.Load("/users")
	assert.True(t, ok)
	assert.False(t, n.IsLeaf())
}
