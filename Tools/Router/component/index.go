package _component

import (
	"context"
	"fiber/Tools/mongodb"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TPostData struct { // * Create
	NAME  string
	HTML  string
	STYLE string
	MAKER string
}

type TPutData struct { // * Update
	ID    string
	NAME  string
	HTML  string
	STYLE string
}

type TDeleteData struct { // * Delete
	ID    string
	MAKER string
}

var Get = func(c *fiber.Ctx) error {
	keyword := strings.ToLower(c.Query("keyword"))
	limit, _ := strconv.ParseInt(c.Query("limit"), 10, 32)
	skip, _ := strconv.ParseInt(c.Query("skip"), 10, 32)
	data, err := CallGetData(keyword, limit, skip)

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
	result, err := CallPostData(p)
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

func CallGetData(keyword string, limit int64, skip int64) ([]map[string]interface{}, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("component")
	filter := bson.M{"keyword": bson.M{"$regex": keyword}}
	opts := options.Find().SetSort(bson.M{"likeCount": -1}).SetLimit(limit).SetSkip(skip)

	cursor, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	var results []map[string]interface{}
	if err := cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	return results, err
}

func CallPostData(data TPostData) (*mongo.InsertOneResult, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("component")
	insertData := bson.M{
		"name":      data.NAME,
		"keyword":   strings.ToLower(data.NAME),
		"html":      data.HTML,
		"style":     data.STYLE,
		"maker":     data.MAKER,
		"like":      [0]string{},
		"likeCount": 0,
		"createdAt": time.Now(),
		"updatedAt": time.Now(),
	}
	result, err := coll.InsertOne(context.TODO(), insertData)
	return result, err
}

func CallPutData(data TPutData) (*mongo.UpdateResult, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("component")
	updateData := bson.M{"updatedAt": time.Now()}
	values := reflect.ValueOf(data)
	for i := 0; i < values.NumField(); i++ {
		dataName := strings.ToLower(values.Type().Field(i).Name)
		if values.Field(i).Interface() != "" && dataName != "id" {
			if dataName == "name" {
				updateData["keyword"] = strings.ToLower(values.Field(i).String())
			}
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
	coll := client.Database("hvData").Collection("component")
	id, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"$and": [2]primitive.M{bson.M{"_id": id}, bson.M{"maker": data.MAKER}}}
	result, err := coll.DeleteOne(context.TODO(), filter)
	return result, err
}
