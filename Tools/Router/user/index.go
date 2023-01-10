package _user

import (
	"context"
	"fiber/Tools/mongodb"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TUserData struct {
	UID   string
	NAME  string
	IMG   string
	EMAIL string
}

var Get = func(c *fiber.Ctx) error {
	data, err := CallGetData(c.Params("params"))
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	return c.Status(200).JSON(data)
}

var Post = func(c *fiber.Ctx) error {
	p := TUserData{}
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

func CallGetData(id string) (primitive.M, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("user")
	var result bson.M
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": userId}
	err = coll.FindOne(context.TODO(), filter).Decode(&result)
	return result, err
}

func CallPostData(TuserData TUserData) (interface{}, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("user")
	userId, err := primitive.ObjectIDFromHex(TuserData.UID)
	if err != nil {
		return nil, err
	}
	var findData bson.M
	insertData := bson.M{
		"name":  TuserData.NAME,
		"img":   TuserData.IMG,
		"email": TuserData.EMAIL,
		"_id":   userId,
	}
	if err := coll.FindOne(context.TODO(), bson.M{"_id": userId}).Decode(&findData); err == mongo.ErrNoDocuments {
		result, err := coll.InsertOne(context.TODO(), insertData)
		return result, err
	} else if err != nil {
		return nil, err
	}
	for key, value := range insertData {
		if findData[key] != value {
			filter := bson.D{{Key: "_id", Value: userId}}
			update := bson.M{"$set": insertData}
			updateResult, err := coll.UpdateOne(context.TODO(), filter, update)
			return updateResult, err
		}
	}
	return "Data already exist and fresh state", nil
}
