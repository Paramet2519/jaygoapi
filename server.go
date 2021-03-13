package main

import (
	"log"
	"net/http"
	"os"
	"database/sql"
	
		
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
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

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()

	t, ok := db.Prepare("SELECT id, title, status FROM todos where id=$1")
	if ok != nil {
		return c.JSON(http.StatusOK, map[int]string{})
	}

	
	row := t.QueryRow(id)
	
	var title, status string

	err = row.Scan(&id, &title, &status)
	if err != nil {
		log.Fatal("can't Scan row into variables", err)
	}

	todo := &Todo{
		ID: id,
		Title: title,
		Status: status,
	}

	return c.JSON(http.StatusOK, todo)
}

func createTodosHandler(e echo.Context) error {
	t := Todo{}
	if err := e.Bind(&t); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	ID := len(todos)
	ID++
	t.ID = ID
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
	e.GET("/Todos", getTodosHandler)
	e.GET("/Todos/:id", getTodoByIdHandler)
	e.POST("/Todos", createTodosHandler)

	
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()

	
	port := os.Getenv("PORT")
	log.Println("port", port)
	e.Start(":" + port)

}
