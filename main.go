package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Completed bool               `json:"completed"`
	Body      string             `json:"body"`
}

var collection *mongo.Collection

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	MONGODB_URI := os.Getenv("mongodb_uri")
	client_options := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), client_options)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB Atlas")

	collection = client.Database("golang_db").Collection("todos")

	app := fiber.New()

	app.Get("/", getHomepage)
	app.Get("/api/todos", getTodos)
	app.Post("/api/todos/", postTodo)
	app.Patch("/api/todos/:id", patchTodo)
	app.Delete("/api/todos/:id", deleteTodo)

	port := os.Getenv("PORT")

	log.Fatal(app.Listen("0.0.0.0:" + port))

}

func getHomepage(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"msg": "Hey VSauce."})
}

func getTodos(c *fiber.Ctx) error {
	var todos []Todo

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return err
		}
		todos = append(todos, todo)
	}

	return c.JSON(todos)
}

func postTodo(c *fiber.Ctx) error {
	todo := new(Todo)

	if err := c.BodyParser(todo); err != nil {
		return err
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Post Todo Request body cannot be empty"})
	}

	insertResult, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}

	todo.ID = insertResult.InsertedID.(primitive.ObjectID)
	return c.Status(201).JSON(fiber.Map{"insert_result": insertResult, "todo inserted": todo})
}

func patchTodo(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"Completed": true}}

	updateResult, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return c.Status(201).JSON(fiber.Map{"completed": id, "update_result": updateResult})

}

func deleteTodo(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}

	filter := bson.M{"_id": id}

	deleteResult, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return c.Status(200).JSON(fiber.Map{"error": err})
	}

	return c.Status(201).JSON(fiber.Map{"delete_result": deleteResult})
}
