package _test

import (
	"errors"
	_component "fiber/Tools/Router/component"
	_component_like "fiber/Tools/Router/component/like"
	_project "fiber/Tools/Router/project"
	_project_component "fiber/Tools/Router/project/component"
	_team "fiber/Tools/Router/team"
	_team_member "fiber/Tools/Router/team/member"
	"math/rand"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Test = func(c *fiber.Ctx) error {
	err := TestCompoent()
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	err = TestProject()
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	err = TestTeam()
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	return c.Status(200).JSON("All Test OK!")
}

func TestProject() error {
	postData := _project.TPostData{
		TITLE: strconv.Itoa(rand.Int()),
		OWNER: strconv.Itoa(rand.Int()),
		STYLE: strconv.Itoa(rand.Int()),
		HTML:  strconv.Itoa(rand.Int()),
	}
	postResult, err := _project.CallPostData(postData)
	insertedID := postResult.InsertedID.(primitive.ObjectID).Hex()
	if err != nil {
		return err
	} else if insertedID == "" {
		return errors.New("cannot create project")
	}
	putData := _project.TPutData{
		ID:    insertedID,
		TITLE: strconv.Itoa(rand.Int()),
		STYLE: strconv.Itoa(rand.Int()),
		HTML:  strconv.Itoa(rand.Int()),
	}
	putResult, err := _project.CallPutData(putData)
	if err != nil {
		return err
	} else if putResult.ModifiedCount != 1 {
		return errors.New("cannot update project")
	}

	putComponentData := _project_component.TPutData{
		ID:    insertedID,
		NAME:  strconv.Itoa(rand.Int()),
		STYLE: strconv.Itoa(rand.Int()),
		HTML:  strconv.Itoa(rand.Int()),
	}
	putResult, err = _project_component.CallPutData(putComponentData)
	if err != nil {
		return err
	} else if putResult.ModifiedCount != 1 {
		return errors.New("cannot update project component")
	}

	getResult, err := _project.CallGetData(postData.OWNER)
	if err != nil {
		return err
	}
	for _, v := range getResult {
		if v["_id"].(primitive.ObjectID).Hex() == insertedID {
			deleteData := _project.TDeleteData{
				ID:    insertedID,
				OWNER: postData.OWNER,
			}
			deleteResult, err := _project.CallDeleteData(deleteData)
			if err != nil {
				return err
			}
			if deleteResult.DeletedCount == 0 {
				return errors.New("cannot delete project")
			}
			return nil
		}
	}
	return errors.New("cannot get project")
}

func TestCompoent() error {
	postData := _component.TPostData{
		NAME:  strconv.Itoa(rand.Int()),
		HTML:  strconv.Itoa(rand.Int()),
		STYLE: strconv.Itoa(rand.Int()),
		MAKER: strconv.Itoa(rand.Int()),
	}
	postResult, err := _component.CallPostData(postData)
	insertedID := postResult.InsertedID.(primitive.ObjectID).Hex()
	if err != nil {
		return err
	} else if insertedID == "" {
		return errors.New("cannot create compoent")
	}
	putData := _component.TPutData{
		NAME:  strconv.Itoa(rand.Int()),
		HTML:  strconv.Itoa(rand.Int()),
		STYLE: strconv.Itoa(rand.Int()),
		ID:    insertedID,
	}
	putResult, err := _component.CallPutData(putData)
	if err != nil {
		return err
	} else if putResult.ModifiedCount != 1 {
		return errors.New("cannot update compoent")
	}

	putLikeData := _component_like.TPutData{
		COMP_ID: insertedID,
		USER_ID: strconv.Itoa(rand.Int()),
	}
	putLikeResult, err := _component_like.CallPutData(putLikeData)
	if err != nil {
		return err
	} else if putLikeResult.ModifiedCount != 1 {
		return errors.New("cannot update compoent like")
	}
	getResult, err := _component.CallGetData(putData.NAME, 0, 0)
	if err != nil {
		return err
	}
	for _, v := range getResult {
		if v["_id"].(primitive.ObjectID).Hex() == insertedID {
			deleteData := _component.TDeleteData{
				ID:    insertedID,
				MAKER: postData.MAKER,
			}
			deleteResult, err := _component.CallDeleteData(deleteData)
			if err != nil {
				return err
			} else if deleteResult.DeletedCount == 0 {
				return errors.New("cannot delete compoent")
			}
			return nil
		}
	}
	return errors.New("cannot get compoent")
}

func TestTeam() error {
	masterId := "639d7567a035ad845d18bb42"
	postData := _team.TPostData{
		NAME:   strconv.Itoa(rand.Int()),
		MEMBER: map[string]string{masterId: "master"},
	}
	postResult, err := _team.CallPostData(postData)
	insertedID := postResult.InsertedID.(primitive.ObjectID).Hex()
	if err != nil {
		return err
	} else if insertedID == "" {
		return errors.New("cannot create team")
	}
	putData := _team.TPutData{
		ID:   insertedID,
		NAME: strconv.Itoa(rand.Int()),
	}
	putResult, err := _team.CallPutData(putData)
	if err != nil {
		return err
	} else if putResult.ModifiedCount != 1 {
		return errors.New("cannot update team")
	}
	putTeamData := _team_member.TPutData{
		ID:     insertedID,
		MEMBER: map[string]string{"63ab394f95c018523cbb004a": "reader"},
	}
	putResult, err = _team_member.CallPutData(putTeamData)
	if err != nil {
		return err
	} else if putResult.ModifiedCount != 1 {
		return errors.New("cannot update team member")
	}

	getResult, err := _team.CallGetData(masterId)
	if err != nil {
		return err
	}

	for _, v := range getResult {
		if v["_id"].(primitive.ObjectID).Hex() == insertedID {
			deleteData := _team.TDeleteData{
				ID:     insertedID,
				MASTER: masterId,
			}
			deleteResult, err := _team.CallDeleteData(deleteData)
			if err != nil {
				return err
			} else if deleteResult.DeletedCount == 0 {
				return errors.New("cannot delete team")
			}
			return nil
		}
	}
	return errors.New("cannot get team")
}
