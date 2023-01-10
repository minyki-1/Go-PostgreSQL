package _project_component

import (
	"context"
	"fiber/Tools/mongodb"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TPutData struct {
	ID    string
	NAME  string
	HTML  string
	STYLE string
}

type TDeleteData struct {
	ID      string
	COMP_ID string
}

var Delete = func(c *fiber.Ctx) error {
	p := TDeleteData{}
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	values := reflect.ValueOf(p)
	for i := 0; i < values.NumField(); i++ {
		if values.Field(i).String() == "" {
			return c.Status(400).JSON("Please send all require data.")
		}
	}
	result, err := CallDeleteData(p)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	return c.Status(200).JSON(result)
}

var Put = func(c *fiber.Ctx) error {
	p := TPutData{}
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	values := reflect.ValueOf(p)
	for i := 0; i < values.NumField(); i++ {
		if values.Field(i).String() == "" {
			return c.Status(400).JSON("Please send all require data.")
		}
	}
	result, err := CallPutData(p)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	return c.Status(200).JSON(result)
}

func CallPutData(data TPutData) (*mongo.UpdateResult, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("project")

	updateData := map[string]map[string]string{"component": {"id": primitive.NewObjectID().Hex(), "name": data.NAME, "html": data.HTML, "style": data.STYLE}}

	id, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": id}
	update := bson.M{"$push": updateData}
	updateResult, err := coll.UpdateOne(context.TODO(), filter, update)
	return updateResult, err
}

func CallDeleteData(data TDeleteData) (*mongo.UpdateResult, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("project")
	id, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": id}
	update := bson.M{"$pull": bson.M{"component": bson.M{"id": data.COMP_ID}}}
	updateResult, err := coll.UpdateOne(context.TODO(), filter, update)
	return updateResult, err
}
