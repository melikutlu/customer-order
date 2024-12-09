package helpers

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"tesodev-korpes/shared/model"
)

func QueryParams(c echo.Context) model.QueryParams {
	return model.QueryParams{
		FirstName:      c.QueryParam("first_name"),
		AgeGreaterThan: c.QueryParam("agt"),
		AgeLessThan:    c.QueryParam("alt"),
		Limit:          c.QueryParam("limit"),
		Offset:         c.QueryParam("offset"),
	}
}

func CreateFilter(params model.QueryParams) (bson.M, error) {
	filter := bson.M{}

	if params.FirstName != "" {
		filter["first_name"] = bson.M{"$regex": params.FirstName, "$options": "i"}
	}
	if params.AgeGreaterThan != "" {
		ageGt, err := strconv.Atoi(params.AgeGreaterThan)
		if err != nil {
			return nil, err
		}
		filter["age"] = bson.M{"$gt": ageGt}
	}
	if params.AgeLessThan != "" {
		ageLt, err := strconv.Atoi(params.AgeLessThan)
		if err != nil {
			return nil, err
		}
		if _, ok := filter["age"]; ok {
			filter["age"].(bson.M)["$lt"] = ageLt
		} else {
			filter["age"] = bson.M{"$lt": ageLt}
		}
	}

	return filter, nil
}
