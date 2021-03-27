package remotePinning

import (
	"fmt"
	"strconv"
	"time"

	db "github.com/alexanderschau/ipfs-pinning-service/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getPins(c *gin.Context, user string) ([]PinStatus, error) {
	results := []PinStatus{}
	cid := c.Request.URL.Query().Get("cid")
	name := c.Request.URL.Query().Get("name")
	status := c.Request.URL.Query().Get("status")
	before := c.Request.URL.Query().Get("before")
	after := c.Request.URL.Query().Get("after")
	limit := c.Request.URL.Query().Get("limit")

	filter := bson.M{
		"owner": user,
	}
	opt := options.Find()

	if cid != "" {
		filter["cid"] = cid
	}

	if name != "" {
		filter["name"] = name
	}

	if status != "" {
		if status == "pinned" {
			filter["$where"] = fmt.Sprintf("this.pinned.length >= 1")
		} else {
			filter["pinned"] = bson.M{"$size": 0}
		}
	}

	if before != "" || after != "" {
		res := bson.M{}

		if before != "" {
			t, err := time.Parse(time.RFC3339, before)

			if err != nil {
				return results, err
			}

			oid, err := primitive.ObjectIDFromHex(fmt.Sprintf("%X0000000000000000", t.Unix()))

			if err != nil {
				return results, err
			}

			res["$lt"] = oid
		}

		if after != "" {
			t, err := time.Parse(time.RFC3339, after)

			if err != nil {
				return results, err
			}

			oid, err := primitive.ObjectIDFromHex(fmt.Sprintf("%X0000000000000000", t.Unix()))

			if err != nil {
				return results, err
			}

			res["$gt"] = oid
		}

		filter["_id"] = res
	}

	if limit != "" {
		limitInt, err := strconv.ParseInt(limit, 10, 64)

		if err != nil {
			return results, err
		}

		opt.SetLimit(limitInt)
	}

	//run search
	res, err := db.Pins.Find(db.Ctx, filter, opt)

	if err != nil {
		return results, err
	}

	var pins []db.Pin
	err = res.All(db.Ctx, &pins)

	if err != nil {
		return results, err
	}

	for _, pin := range pins {
		r, err := pinToPinStatus(pin)
		if err != nil {
			return results, err
		}
		results = append(results, r)
	}

	return results, nil
}

func pinToPinStatus(pin db.Pin) (PinStatus, error) {
	objID, err := primitive.ObjectIDFromHex(pin.ObjectID)

	if err != nil {
		return PinStatus{}, err
	}

	status := PINNING
	if len(pin.Pinned) > 0 {
		status = PINNED
	}

	return PinStatus{
		Requestid: pin.ObjectID,
		Status:    status,
		Created:   objID.Timestamp(),
		Pin: Pin{
			Cid:  pin.Cid,
			Name: pin.Name,
		},
		Delegates: []string{},
	}, nil
}
