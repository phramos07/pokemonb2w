package control

import (
	"encoding/json"
	"net/http"
	"pokemonb2w/internal/model"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

// An AppInfoResponse response model
//
// This is used for returning a response the App's info.
//
// swagger:response appinfoResponse
type AppInfoResponse struct {
	//in:body
	AppInfo model.AppInfo `json:"app-info"`
}

// swagger:operation GET /app-info appinfo getAppInfo
// ---
// summary: Retrieves APP info.
// description: JSON with app info will be returned.
// responses:
//   '200':
//     "$ref": "#/responses/appinfoResponse"
func getAppInfo(w http.ResponseWriter, r *http.Request) {

	nameVar := viper.GetString("app.deploy.herokuAppNameVar")
	releaseVar := viper.GetString("app.deploy.herokuBuildVar")
	releasedAtVar := viper.GetString("app.deploy.herokuCreatedAtVar")

	appInfo := model.AppInfo{
		Name:          viper.GetString(nameVar),
		ReleaseNumber: viper.GetString(releaseVar),
		ReleasedAt:    viper.GetString(releasedAtVar),
	}

	err := json.NewEncoder(w).Encode(appInfo)

	if err != nil {
		panic(err.Error())
	}
}

// AddAppInfoRoute adds routes from path app-info to the main API router
func AddAppInfoRoute(r *mux.Router) {
	r.HandleFunc(apiInfoPath, getAppInfo).Methods(http.MethodGet, http.MethodOptions)
}
