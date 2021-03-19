package remotePinning

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	db "github.com/alexanderschau/ipfs-pinning-service/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func sendErr(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": gin.H{
			"reason":  "INTERNAL_SERVER_ERROR",
			"details": fmt.Sprint(err),
		},
	})
}

//PinsGet - get all pins
func PinsGet(c *gin.Context) {
	Cid := c.Request.URL.Query().Get("cid")

	results := []PinStatus{}

	if Cid != "" {
		res := db.Collection.FindOne(db.Ctx, bson.M{
			"cid":   Cid,
			"owner": "example",
		})

		var pin db.Pin
		err := res.Decode(&pin)

		if err != nil {
			sendErr(c, err)
			return
		}

		objID, err := primitive.ObjectIDFromHex(pin.ObjectID)

		if err != nil {
			sendErr(c, err)
			return
		}

		results = append(results, PinStatus{
			Requestid: pin.ObjectID,
			Status:    PINNED,
			Created:   objID.Timestamp(),
			Pin: Pin{
				Cid:  Cid,
				Name: pin.Name,
			},
			Delegates: []string{},
		})
	}

	if Cid == "" {
		res, err := db.Collection.Find(db.Ctx, bson.M{
			"owner": "example",
		})

		if err != nil {
			sendErr(c, err)
			return
		}

		var pins []db.Pin
		err = res.All(db.Ctx, &pins)

		if err != nil {
			sendErr(c, err)
			return
		}

		for _, pin := range pins {
			objID, err := primitive.ObjectIDFromHex(pin.ObjectID)

			if err != nil {
				sendErr(c, err)
				return
			}

			results = append(results, PinStatus{
				Requestid: pin.ObjectID,
				Status:    PINNED,
				Created:   objID.Timestamp(),
				Pin: Pin{
					Cid:  pin.Cid,
					Name: pin.Name,
				},
				Delegates: []string{},
			})
		}
	}

	c.JSON(http.StatusOK, PinResults{
		Count:   int32(len(results)),
		Results: results,
	})
}

//PinsPost - add new pin
func PinsPost(c *gin.Context) {
	//accessToken := strings.Split(c.Request.Header.Get("Authorization"), " ")[1]
	jsonData, _ := ioutil.ReadAll(c.Request.Body)
	var inputData Pin
	err := json.Unmarshal(jsonData, &inputData)
	if err != nil {
		sendErr(c, err)
		return
	}

	res, err := db.Collection.InsertOne(db.Ctx, db.Pin{
		Cid:   inputData.Cid,
		Name:  inputData.Name,
		Owner: "example",
	})
	if err != nil {
		sendErr(c, err)
		return
	}
	rid := res.InsertedID.(primitive.ObjectID).Hex()

	c.JSON(http.StatusAccepted, PinStatus{
		Requestid: rid,
		Status:    PINNING,
		Created:   res.InsertedID.(primitive.ObjectID).Timestamp(),
		Pin:       inputData,
		Delegates: []string{},
	})
}

//PinsRequestidDelete - delete pin
func PinsRequestidDelete(c *gin.Context) {
	requestID, _ := c.Params.Get("requestid")

	objID, err := primitive.ObjectIDFromHex(requestID)
	if err != nil {
		sendErr(c, err)
		return
	}

	_, err = db.Collection.DeleteOne(db.Ctx, bson.M{
		"_id":   objID,
		"owner": "example",
	})

	if err != nil {
		sendErr(c, err)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{})
}

//PinsRequestidGet - get pin by requestID
func PinsRequestidGet(c *gin.Context) {
	requestID, _ := c.Params.Get("requestid")

	objID, err := primitive.ObjectIDFromHex(requestID)

	if err != nil {
		sendErr(c, err)
		return
	}

	res := db.Collection.FindOne(db.Ctx, bson.M{
		"_id": objID,
	})

	var pin db.Pin

	err = res.Decode(&pin)

	if err != nil {
		fmt.Println(err)
		sendErr(c, err)
		return
	}

	c.JSON(http.StatusOK, PinStatus{
		Requestid: requestID,
		Status:    PINNED,
		Created:   objID.Timestamp(),
		Pin: Pin{
			Cid:  pin.Cid,
			Name: pin.Name,
		},
		Delegates: []string{},
	})
}

//PinsRequestidPost - update pin
func PinsRequestidPost(c *gin.Context) {
	requestID, _ := c.Params.Get("requestid")

	jsonData, _ := ioutil.ReadAll(c.Request.Body)
	var inputData Pin
	err := json.Unmarshal(jsonData, &inputData)

	if err != nil {
		sendErr(c, err)
		return
	}

	//remove old pin
	_, err = db.Collection.DeleteOne(db.Ctx, bson.M{
		"_id": requestID,
	})

	if err != nil {
		sendErr(c, err)
		return
	}

	//pin new one
	res, err := db.Collection.InsertOne(db.Ctx, db.Pin{
		Cid:   inputData.Cid,
		Name:  inputData.Name,
		Owner: "example",
	})

	if err != nil {
		sendErr(c, err)
		return
	}

	rid := res.InsertedID.(primitive.ObjectID).Hex()

	c.JSON(http.StatusAccepted, PinStatus{
		Requestid: rid,
		Status:    PINNING,
		Created:   res.InsertedID.(primitive.ObjectID).Timestamp(),
		Pin:       inputData,
		Delegates: []string{},
	})
}
