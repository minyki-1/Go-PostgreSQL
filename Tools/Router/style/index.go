package _style

import (
	"context"
	"fiber/Tools/mongodb"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Get = func(c *fiber.Ctx) error {
	data, err := CallGetData(c.Params("params"))
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	return c.Status(200).JSON(data["style"])
}

func CallGetData(id string) (primitive.M, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("project")
	var result bson.M
	projectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": projectId}
	err = coll.FindOne(context.TODO(), filter).Decode(&result)
	return result, err
}
