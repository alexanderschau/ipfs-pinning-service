package remotePinning

import (
	"os"
	"strings"

	db "github.com/alexanderschau/ipfs-pinning-service/database"
	"go.mongodb.org/mongo-driver/mongo"
)

func addPin(inputData Pin, user string) (*mongo.InsertOneResult, error) {
	res, err := db.Pins.InsertOne(db.Ctx, db.Pin{
		Cid:      inputData.Cid,
		Name:     inputData.Name,
		Owner:    user,
		Clusters: strings.Split(os.Getenv("STANDARD_CLUSTER"), ","),
		Pinned:   []string{},
	})
	return res, err
}
