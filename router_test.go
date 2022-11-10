package jago

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJagoStatic(t *testing.T) {
	g := New()
	loadJagoRoutes(g, parseAPI)
	r := g.router
	c := g.NewContext(nil, nil)
	r.find("/1/users/", "GET", c)
	assert.Equal(t, "/1/users", c.Path())
}

func TestJagoParam1(t *testing.T) {
	g := New()
	loadJagoRoutes(g, parseAPI)
	r := g.router
	c := g.NewContext(nil, nil)
	r.find("/1/classes/a/Obj", "GET", c)
	assert.Equal(t, "/1/classes/:className/:objectId", c.Path())
	assert.Equal(t, "a", c.Param("className"))
	assert.Equal(t, "Obj", c.Param("objectId"))
}

func TestJagoParam2(t *testing.T) {
	g := New()
	loadJagoRoutes(g, parseAPI)
	r := g.router
	c := g.NewContext(nil, nil)
	r.find("/1/classes/a/Obj", "POST", c)
	assert.Equal(t, "", c.Path())
}

func TestJagoParam3(t *testing.T) {
	g := New()
	loadJagoRoutes(g, parseAPI)
	r := g.router
	c := g.NewContext(nil, nil)
	r.find("/1/classes/Category/Item", "GET", c)
	assert.Equal(t, "/1/:type/Category/Item", c.Path())
	assert.Equal(t, "classes", c.Param("type"))
}

func TestJagoWildcard1(t *testing.T) {
	g := New()
	loadJagoRoutes(g, parseAPI)
	r := g.router
	c := g.NewContext(nil, nil)
	r.find("/1/functions/", "GET", c)
	assert.Equal(t, "/1/functions/*", c.Path())
}

func TestJagoWildcard2(t *testing.T) {
	g := New()
	loadJagoRoutes(g, parseAPI)
	r := g.router
	c := g.NewContext(nil, nil)
	r.find("/1/functions/funcA/hello-world", "GET", c)
	assert.Equal(t, "/1/functions/*", c.Path())
}
