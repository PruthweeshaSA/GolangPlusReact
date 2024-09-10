package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	var greeting = "Hey VSauce!"
	var myName = "Michael"
	fmt.Println(greeting, myName, "here.")

	app := fiber.New()

	todos := []Todo{}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"msg": "Hey Vsauce! Michael here. I'm supposed to show todo list somewhere here ... Or am I?"})
	})

	app.Get("/api/todos", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(todos)
	})

	// Create a TODO
	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{}

		if err := c.BodyParser(todo); err != nil {
			return err
		}

		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Todo body is required"})
		}

		todo.ID = len(todos)
		todos = append(todos, *todo)
		return c.Status(201).JSON(todo)
	})

	// Update a TODO
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos[i].Completed = !todos[i].Completed
				return c.Status(200).JSON(todos[i])
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "Requested Todo not found"})
	})

	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		id_num, _ := strconv.Atoi(id)
		if id_num < len(todos) {
			todos = append(todos[:id_num], todos[id_num+1:]...)
			return c.Status(200).JSON(todos)
		}

		return c.Status(404).JSON(fiber.Map{"error": "Requested index not found"})
	})

	log.Fatal(app.Listen(":4000"))

}
