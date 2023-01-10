package _project

import (
	"context"
	"fiber/Tools/mongodb"
	"reflect"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TPostData struct {
	TITLE string
	OWNER string
	STYLE string
	HTML  string
}

type TDeleteData struct {
	ID    string
	OWNER string
}

type TPutData struct {
	ID    string
	TITLE string
	STYLE string
	HTML  string
}

var Get = func(c *fiber.Ctx) error {
	data, err := CallGetData(c.Params("params"))
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	return c.Status(200).JSON(data)
}

var Post = func(c *fiber.Ctx) error {
	p := TPostData{}
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	values := reflect.ValueOf(p)
	for i := 0; i < values.NumField(); i++ {
		if values.Field(i).String() == "" {
			return c.Status(400).JSON("Please send all require data.")
		}
	}
	data, err := CallPostData(p)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	return c.Status(200).JSON(data)
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
	result, err := CallPutData(p)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	return c.Status(200).JSON(result)
}

func CallGetData(owner string) ([]map[string]interface{}, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("project")
	filter := bson.M{"owner": owner}
	opts := options.Find().SetSort(bson.M{"updatedAt": -1})
	cursor, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	var results []bson.D
	err = cursor.All(context.TODO(), &results)
	dataArray := make([]map[string]interface{}, 0, len(results))
	for _, result := range results {
		dataMap := map[string]interface{}{}
		for _, k := range result {
			dataMap[k.Key] = k.Value
		}
		dataArray = append(dataArray, dataMap)
	}
	return dataArray, err
}

func CallPostData(p TPostData) (*mongo.InsertOneResult, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("project")
	insertData := bson.M{
		"title":     p.TITLE,
		"owner":     p.OWNER,
		"style":     p.STYLE,
		"html":      p.HTML,
		"component": [0]map[string]string{},
		"createdAt": time.Now(),
		"updatedAt": time.Now(),
	}
	result, err := coll.InsertOne(context.TODO(), insertData)
	return result, err
}

func CallPutData(data TPutData) (*mongo.UpdateResult, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("project")

	updateData := bson.M{"updatedAt": time.Now()}
	values := reflect.ValueOf(data)
	for i := 0; i < values.NumField(); i++ {
		dataName := strings.ToLower(values.Type().Field(i).Name)
		if values.Field(i).Interface() != "" && dataName != "id" {
			updateData[dataName] = values.Field(i).Interface()
		}
	}

	id, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.M{"$set": updateData}
	updateResult, err := coll.UpdateOne(context.TODO(), filter, update)
	return updateResult, err
}

func CallDeleteData(data TDeleteData) (*mongo.DeleteResult, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("project")
	id, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"$and": [2]primitive.M{bson.M{"_id": id}, bson.M{"owner": data.OWNER}}}
	result, err := coll.DeleteOne(context.TODO(), filter)
	return result, err
}
