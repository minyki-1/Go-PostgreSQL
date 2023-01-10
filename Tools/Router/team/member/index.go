package _team_member

import (
	"context"
	"fiber/Tools/mongodb"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TPutData struct {
	ID     string
	MEMBER map[string]string
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

	coll := client.Database("hvData").Collection("team")
	teamId, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": teamId}
	var result bson.M
	err = coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	newMemberList := map[string]string{}
	for key1, value1 := range data.MEMBER {
		inMember := false
		for key2 := range result["member"].(primitive.M) {
			if key1 == key2 {
				inMember = true
			}
		}
		if !inMember {
			newMemberList[key1] = value1
		}
	}
	for key1, value1 := range result["member"].(primitive.M) {
		inMember := false
		for key2 := range data.MEMBER {
			if key1 == key2 {
				inMember = true
			}
		}
		if !inMember {
			newMemberList[key1] = value1.(string)
		}
	}
	id, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return nil, err
	}
	filter = bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"member": newMemberList}}
	updateResult, err := coll.UpdateOne(context.TODO(), filter, update)
	return updateResult, err
}
