package _team

import (
	"context"
	"errors"
	"fiber/Tools/mongodb"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TPostData struct {
	NAME   string
	MEMBER map[string]string
}

type TDeleteData struct {
	ID     string
	MASTER string
}

type TPutData struct {
	ID   string
	NAME string
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
	result, err := CallPostData(p)
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

func CallGetData(memberId string) ([]map[string]interface{}, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("team")
	filter := bson.M{
		"$or": [4]primitive.M{
			bson.M{"member." + memberId: "master"},
			bson.M{"member." + memberId: "manager"},
			bson.M{"member." + memberId: "maker"},
			bson.M{"member." + memberId: "reader"},
		},
	}
	var cursor *mongo.Cursor
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var teamResults []map[string]interface{}
	err = cursor.All(context.TODO(), &teamResults)
	if err != nil {
		return nil, err
	}

	coll = client.Database("hvData").Collection("user")
	memberList := []map[string]interface{}{}
	for _, v := range teamResults {
		for l := range v["member"].(map[string]interface{}) {
			existMember := false
			for _, x := range memberList {
				if x["_id"].(primitive.ObjectID).Hex() == l {
					existMember = true
				}
			}
			if !existMember {
				objID, err := primitive.ObjectIDFromHex(l)
				if err != nil {
					return nil, err
				}
				memberList = append(memberList, map[string]interface{}{"_id": objID})
			}
		}
	}
	if len(memberList) == 0 {
		return nil, errors.New("cannot found member")
	}
	filter = bson.M{"$or": memberList}
	cursor, err = coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var memberResults []map[string]interface{}
	err = cursor.All(context.TODO(), &memberResults)

	var newTeamList []map[string]interface{}
	for _, result := range teamResults {
		newTeam := map[string]interface{}{}
		newMemberList := []map[string]interface{}{}
		for k, v := range result {
			if k == "member" {
				for l, w := range v.(map[string]interface{}) {
					for _, x := range memberResults {
						if l == x["_id"].(primitive.ObjectID).Hex() {
							newMember := x
							newMember["role"] = w.(string)
							newMemberList = append(newMemberList, newMember)
						}
					}
				}
			} else {
				newTeam[k] = v
			}
		}
		newTeam["member"] = newMemberList
		newTeamList = append(newTeamList, newTeam)
	}
	return newTeamList, err
}

func CallPostData(data TPostData) (*mongo.InsertOneResult, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("team")
	insertData := bson.M{
		"name":   data.NAME,
		"member": data.MEMBER,
	}
	result, err := coll.InsertOne(context.TODO(), insertData)
	return result, err
}

func CallDeleteData(data TDeleteData) (*mongo.DeleteResult, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("team")
	id, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"$and": [2]primitive.M{bson.M{"_id": id}, bson.M{"member." + data.MASTER: "master"}}}
	result, err := coll.DeleteOne(context.TODO(), filter)
	return result, err
}

func CallPutData(data TPutData) (*mongo.UpdateResult, error) {
	client := mongodb.GetMongoClient()
	coll := client.Database("hvData").Collection("team")
	updateData := bson.M{"name": data.NAME}
	id, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.M{"$set": updateData}
	updateResult, err := coll.UpdateOne(context.TODO(), filter, update)
	return updateResult, err
}
