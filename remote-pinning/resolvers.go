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

	results, err := getPins(c, user)

	if err != nil {
		sendErr(c, err)
		return
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

	fmt.Println(dnsaddrFormatter(strings.Split(os.Getenv("STANDARD_CLUSTER"), ",")))
	delegates, err := dnsaddrFormatter(strings.Split(os.Getenv("STANDARD_CLUSTER"), ","))
	if err != nil {
		sendErr(c, err)
		return
	}
	c.JSON(http.StatusAccepted, PinStatus{
		Requestid: rid,
		Status:    PINNING,
		Created:   res.InsertedID.(primitive.ObjectID).Timestamp(),
		Pin:       inputData,
		Delegates: delegates,
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

	delegates, err := dnsaddrFormatter(strings.Split(os.Getenv("STANDARD_CLUSTER"), ","))
	if err != nil {
		sendErr(c, err)
		return
	}

	c.JSON(http.StatusOK, PinStatus{
		Requestid: requestID,
		Status:    status,
		Created:   objID.Timestamp(),
		Pin: Pin{
			Cid:  pin.Cid,
			Name: pin.Name,
		},
		Delegates: delegates,
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

	delegates, err := dnsaddrFormatter(strings.Split(os.Getenv("STANDARD_CLUSTER"), ","))
	if err != nil {
		sendErr(c, err)
		return
	}

	c.JSON(http.StatusAccepted, PinStatus{
		Requestid: rid,
		Status:    PINNING,
		Created:   res.InsertedID.(primitive.ObjectID).Timestamp(),
		Pin:       inputData,
		Delegates: delegates,
	})
}
