package rattle_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	rattle2 "github.com/jjeffcaii/rattle"
	"github.com/stretchr/testify/assert"
)

type FakeUser struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

func TestRattle(t *testing.T) {
	c := &rattle2.Config{
		Namespace: "",
		Token:     "",
		Bootstrap: "",
	}

	rattle, err := rattle2.NewRattle(c)
	assert.NoError(t, err)

	rattle.GET("/users/:id", func(c *rattle2.Context) error {
		return c.JSON(&FakeUser{
			ID:        1,
			Name:      "foobar",
			CreatedAt: time.Now(),
		})
	})

	err = rattle.Start(context.Background())
	assert.NoError(t, err)

	res, err := rattle.Request("fake-namespace").GET(context.Background(), "/users/1234", nil)
	assert.NoError(t, err)

	var u FakeUser
	err = res.Bind(&u)
	assert.NoError(t, err)
	fmt.Printf("response user: %v\n", u)
}
