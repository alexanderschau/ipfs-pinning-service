package clusterController

import (
	"context"
	"log"
	"os"
	"os/signal"

	db "github.com/alexanderschau/ipfs-pinning-service/database"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/ipfs-cluster/api"
	"github.com/ipfs/ipfs-cluster/api/rest/client"
	"github.com/robfig/cron"
	"go.mongodb.org/mongo-driver/bson"
)

var clusterName = os.Getenv("CLUSTER_DOMAIN")

func StartController() {
	c := cron.New()
	c.AddFunc("*/2 * * * * *", func() {
		log.Println("Start runner")
		err := Controller()
		if err != nil {
			log.Println(err)
		}
		log.Println("Finish runner")
	})
	go c.Start()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}

func Controller() error {
	ctx := context.TODO()
	client.NewDefaultClient(&client.Config{})
	client, err := client.NewDefaultClient(&client.Config{Host: os.Getenv("CLUSTER_HOSTNAME"), Port: os.Getenv("CLUSTER_PORT")})

	if err != nil {
		return err
	}

	allocations, err := client.Allocations(ctx, api.DataType)
	if err != nil {
		return err
	}

	res, err := db.Pins.Find(db.Ctx, bson.M{"clusters": bson.M{"$all": []string{
		clusterName,
	}}})

	if err != nil {
		return err
	}

	var pins []db.Pin
	err = res.All(db.Ctx, &pins)

	if err != nil {
		return err
	}

	ruleCids := []string{}
	for _, pin := range pins {
		ruleCids = append(ruleCids, pin.Cid)
	}

	currentCids := []string{}
	for _, pin := range allocations {
		currentCids = append(currentCids, pin.Cid.String())
	}

	//pin all new ones
	for _, id := range ruleCids {
		if !contains(currentCids, id) {
			ID, err := cid.Decode(id)
			if err != nil {
				return err
			}

			_, err = client.Pin(ctx, ID, api.PinOptions{})
			if err != nil {
				return err
			}
		}
	}

	//remove old ones
	for _, id := range currentCids {
		if !contains(ruleCids, id) {
			ID, err := cid.Decode(id)
			if err != nil {
				return err
			}

			_, err = client.Unpin(ctx, ID)
			if err != nil {
				return err
			}
		}
	}

	//update status list
	for _, pin := range pins {
		if !contains(pin.Pinned, clusterName) {
			id, err := cid.Decode(pin.Cid)

			if err != nil {
				return err
			}

			pinStatus, err := client.Status(ctx, id, false)

			if err != nil {
				return err
			}

			if checkStatus(pinStatus.PeerMap) {
				db.Pins.UpdateOne(db.Ctx, bson.M{
					"cid":    pin.Cid,
					"pinned": bson.M{"$nin": []string{clusterName}},
				}, bson.M{"$push": bson.M{"pinned": clusterName}})
			}

		}
	}

	return nil
}

func checkStatus(peers map[string]*api.PinInfoShort) bool {
	for _, peer := range peers {
		if peer.Status == api.TrackerStatusPinned {
			return true
		}
	}
	return false
}
