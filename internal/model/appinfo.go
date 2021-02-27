package model

// AppInfo model.
//
//swagger:model appinfo
type AppInfo struct {
	// App name.
	//
	// required: true
	// x-order: 0
	Name string `json:"name"`

	// App version.
	//
	// required: true
	// x-order: 1
	ReleaseNumber string `json:"release"`

	// App release date.
	//
	// required: true
	// x-order: 2
	ReleasedAt string `json:"releasedAt"`
}
