package remotePinning

import (
	db "github.com/alexanderschau/ipfs-pinning-service/database"
	"go.mongodb.org/mongo-driver/mongo"
)

func addPin(inputData Pin, user string) (*mongo.InsertOneResult, error) {
	res, err := db.Pins.InsertOne(db.Ctx, db.Pin{
		Cid:   inputData.Cid,
		Name:  inputData.Name,
		Owner: user,
	})
	return res, err
}
