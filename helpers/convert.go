package helpers

import (
	"gopkg.in/mgo.v2/bson"
	"errors"
	"strings"
)

var (
	ErrMalformedString = errors.New("malformed string")
)

func StringToObjectIds(s string) ([]bson.ObjectId, error) {
	var ids []bson.ObjectId

	for _, id := range strings.Split(s, ",") {
		if !bson.IsObjectIdHex(id) {
			return nil, ErrMalformedString
		}

		ids = append(ids, bson.ObjectIdHex(id))
	}

	return ids, nil
}

func StringToObjectId(s string) (bson.ObjectId, error) {
	if !bson.IsObjectIdHex(s) {
		return "", ErrMalformedString
	}

	return bson.ObjectIdHex(s), nil
}
