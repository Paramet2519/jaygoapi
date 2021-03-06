package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

var todos = map[int]*Todo{
	1: &Todo{ID: 1, Title: "pay phone bills", Status: "active"},
}

func getTodoByIdHandler(c echo.Context) error {
	var id int
	err := echo.PathParamsBinder(c).Int("id", &id).BindError()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	t, ok := todos[id]
	if !ok {
		return c.JSON(http.StatusOK, map[int]string{})
	}
	return c.JSON(http.StatusOK, t)
}

func createTodosHandler(e echo.Context) error {
	t := Todo{}
	if err := e.Bind(&t); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	ID := len(todos)
	ID++
	todos[t.ID] = &t
	return e.JSON(http.StatusCreated, "Create todos")
}
func getTodosHandler(c echo.Context) error {
	items := []*Todo{}
	for _, item := range todos {
		items = append(items, item)
	}
	return c.JSON(http.StatusOK, items)
}
func helloHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Paramet",
	})
}
func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/hello", helloHandler)
	e.GET("Todos", getTodosHandler)
	e.GET("Todos:id", getTodoByIdHandler)
	e.POST("/Todos", createTodosHandler)
	port := os.Getenv("PORT")
	log.Println("port", port)
	e.Start(":" + port)

}
