package remotePinning

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

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
	//check auth
	check, user := authMiddleware(c)
	if !check {
		return
	}

	Cid := c.Request.URL.Query().Get("cid")

	results := []PinStatus{}

	if Cid != "" {
		res := db.Pins.FindOne(db.Ctx, bson.M{
			"cid":   Cid,
			"owner": user,
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
		res, err := db.Pins.Find(db.Ctx, bson.M{
			"owner": user,
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

			status := PINNING
			if len(pin.Pinned) > 0 {
				status = PINNED
			}

			results = append(results, PinStatus{
				Requestid: pin.ObjectID,
				Status:    status,
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
	//check auth
	check, user := authMiddleware(c)
	if !check {
		return
	}

	jsonData, _ := ioutil.ReadAll(c.Request.Body)
	var inputData Pin
	err := json.Unmarshal(jsonData, &inputData)
	if err != nil {
		sendErr(c, err)
		return
	}

	res, err := addPin(inputData, user)
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
		Delegates: strings.Split(os.Getenv("STANDARD_CLUSTER"), ","),
	})
}

//PinsRequestidDelete - delete pin
func PinsRequestidDelete(c *gin.Context) {
	//check auth
	check, user := authMiddleware(c)
	if !check {
		return
	}

	requestID, _ := c.Params.Get("requestid")

	objID, err := primitive.ObjectIDFromHex(requestID)
	if err != nil {
		sendErr(c, err)
		return
	}

	_, err = db.Pins.DeleteOne(db.Ctx, bson.M{
		"_id":   objID,
		"owner": user,
	})

	if err != nil {
		sendErr(c, err)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{})
}

//PinsRequestidGet - get pin by requestID
func PinsRequestidGet(c *gin.Context) {
	//check auth
	check, user := authMiddleware(c)
	if !check {
		return
	}

	requestID, _ := c.Params.Get("requestid")

	objID, err := primitive.ObjectIDFromHex(requestID)

	if err != nil {
		sendErr(c, err)
		return
	}

	res := db.Pins.FindOne(db.Ctx, bson.M{
		"_id":   objID,
		"owner": user,
	})

	var pin db.Pin

	err = res.Decode(&pin)

	if err != nil {
		fmt.Println(err)
		sendErr(c, err)
		return
	}

	status := PINNING
	if len(pin.Pinned) > 0 {
		status = PINNED
	}

	c.JSON(http.StatusOK, PinStatus{
		Requestid: requestID,
		Status:    status,
		Created:   objID.Timestamp(),
		Pin: Pin{
			Cid:  pin.Cid,
			Name: pin.Name,
		},
		Delegates: pin.Clusters,
	})
}

//PinsRequestidPost - update pin
func PinsRequestidPost(c *gin.Context) {
	//check auth
	check, user := authMiddleware(c)
	if !check {
		return
	}

	requestID, _ := c.Params.Get("requestid")

	jsonData, _ := ioutil.ReadAll(c.Request.Body)
	var inputData Pin
	err := json.Unmarshal(jsonData, &inputData)

	if err != nil {
		sendErr(c, err)
		return
	}

	//remove old pin
	_, err = db.Pins.DeleteOne(db.Ctx, bson.M{
		"_id":   requestID,
		"owner": user,
	})

	if err != nil {
		sendErr(c, err)
		return
	}

	//pin new one
	res, err := addPin(inputData, user)

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
		Delegates: strings.Split(os.Getenv("STANDARD_CLUSTER"), ","),
	})
}
