package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/JamesYYang/jago"
)

type (
	Employee struct {
		Name string `json:"Name"`
		Age  int    `json:"Age"`
		ID   string `json:"ID"`
	}
)

// /users/new
// /users/:id
// /users/:id/address/:address
// /users/:id/report/download
func main() {
	jago := jago.New()
	g := jago.Group("/users", middlewareOne)
	g.Get("/new", middlewareSingle, TestHandler)
	g.Get("/:id", QuerySingleEmployee)
	g.Get("/:id/address/:address", TestHandler)
	g.Get("/:id/report/download", TestHandler)
	g.Get("/*", TestHandler)
	g.Get("/info/*", TestHandler)

	gOrder := jago.Group("/orders", middlewareTwo)
	gOrder.Get("/new", middlewareSingle, TestHandler)
	gOrder.Get("/:id", TestHandler)
	gOrder.Get("/:id/items/:item-number", TestHandler)
	gOrder.Get("/:id/report/download", TestHandler)
	gOrder.Get("/*", TestHandler)
	gOrder.Get("/shipping/*", TestHandler)
	server := &http.Server{
		Handler: jago,
		Addr:    ":8080",
	}

	log.Fatal(server.ListenAndServe())
}

func TestHandler(c jago.Context) error {
	result := "hello Jago\n"
	result += fmt.Sprintf("your path: %s \n", c.Request().URL)
	result += fmt.Sprintf("match path: %s \n", c.Path())
	result += fmt.Sprintf("params id: %s \n", c.Param("id"))
	result += fmt.Sprintf("params address: %s \n", c.Param("address"))
	result += fmt.Sprintf("params item-number: %s \n", c.Param("item-number"))
	return c.String(http.StatusOK, result)
}

func middlewareOne(c jago.Context) error {
	log.Printf("this is first middleware")
	return c.Next()
}

func middlewareTwo(c jago.Context) error {
	log.Printf("this is second middleware")
	return c.Next()
}

func middlewareSingle(c jago.Context) error {
	log.Printf("this is single middleware")
	return c.Next()
}

func QuerySingleEmployee(c jago.Context) error {
	e1 := Employee{
		Name: "Jago",
		Age:  18,
		ID:   "Jago-001",
	}
	return c.JSON(http.StatusOK, e1)
}
