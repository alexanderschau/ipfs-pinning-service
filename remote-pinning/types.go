package remotePinning

import "time"

// Pin - Pin object
type Pin struct {

	// Content Identifier (CID) to be pinned recursively
	Cid string `json:"cid"`

	// Optional name for pinned data; can be used for lookups later
	Name string `json:"name,omitempty"`

	// Optional list of multiaddrs known to provide the data
	Origins []string `json:"origins,omitempty"`

	// Optional metadata for pin object
	Meta map[string]string `json:"meta,omitempty"`
}

// Status : Status a pin object can have at a pinning service
type Status string

// List of Status
const (
	QUEUED  Status = "queued"
	PINNING Status = "pinning"
	PINNED  Status = "pinned"
	FAILED  Status = "failed"
)

// PinStatus - Pin object with status
type PinStatus struct {

	// Globally unique identifier of the pin request; can be used to check the status of ongoing pinning, or pin removal
	Requestid string `json:"requestid"`

	Status Status `json:"status"`

	// Immutable timestamp indicating when a pin request entered a pinning service; can be used for filtering results and pagination
	Created time.Time `json:"created"`

	Pin Pin `json:"pin"`

	// List of multiaddrs designated by pinning service for transferring any new data from external peers
	Delegates []string `json:"delegates"`

	// Optional info for PinStatus response
	Info map[string]string `json:"info,omitempty"`
}

// PinResults - Response used for listing pin objects matching request
type PinResults struct {

	// The total number of pin objects that exist for passed query filters
	Count int32 `json:"count"`

	// An array of PinStatus results
	Results []PinStatus `json:"results"`
}
