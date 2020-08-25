package rattle_test

import (
	"testing"

	"github.com/jjeffcaii/rattle"
	"github.com/stretchr/testify/assert"
)

func TestRouter_Route(t *testing.T) {
	router := rattle.NewRouter()

	h := func(context *rattle.Context) error {
		return nil
	}

	err := router.Route(rattle.GET, "/users/:userId", h)
	assert.NoError(t, err)
	err = router.Route(rattle.POST, "/users/:userId", h)
	assert.NoError(t, err)
	_ = router.Route(rattle.GET, "/products", h)
	_ = router.Route(rattle.GET, "/products/:productId", h)
	_ = router.Route(rattle.GET, "/users/:userId/orders/:orderId", h)

	params, _, ok := router.Get(rattle.GET, "/users/7777")
	assert.True(t, ok)
	assert.Equal(t, "7777", params.GetOrDefault("userId", ""))
	params, _, ok = router.Get(rattle.POST, "/users/8888")
	assert.True(t, ok)
	assert.Equal(t, "8888", params.GetOrDefault("userId", ""))

	params, _, ok = router.Get(rattle.GET, "/users/7777/orders/8888")
	assert.True(t, ok)
	assert.Equal(t, "7777", params.GetOrDefault("userId", ""))
	assert.Equal(t, "8888", params.GetOrDefault("orderId", ""))
}
