package db

type Pin struct {
	Cid      string   `json:"cid"`
	Clusters []string `json:"clusters"`
	Pinned   []string `json:"pinned"`
	Name     string   `json:"name"`
	Owner    string   `json:"owner"`
	ObjectID string   `bson:"_id,omitempty"`
}
