package _component_like

import (
	"context"
	"fiber/Tools/mongodb"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TPutData struct {
	COMP_ID string
	USER_ID string
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

func CallPutData(data TPutData) (*mongo.UpdateResult, error) {
	client := mongodb.GetMongoClient()

	coll := client.Database("hvData").Collection("component")
	compId, err := primitive.ObjectIDFromHex(data.COMP_ID)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	filter := bson.M{"_id": compId}
	err = coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	newLikeList := []string{data.USER_ID}
	for _, v := range result["like"].(primitive.A) {
		if v.(string) == data.USER_ID {
			newLikeList = newLikeList[1:]
		} else {
			newLikeList = append(newLikeList, v.(string))
		}
	}
	update := bson.M{
		"like":      newLikeList,
		"likeCount": len(newLikeList),
	}
	updateResult, err := coll.UpdateOne(context.TODO(), filter, bson.M{"$set": update})
	return updateResult, err
}
