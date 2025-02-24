package configurations

import (
	"github.com/nicholasjackson/env"
)

var appConfig *appConfigs
var goUri *string = env.String("GO_URI", false, "0.0.0.0:9091", "Bind address for the app server")

type appConfigs struct {
	appURI string
}

func NewAppConfig() (*appConfigs, error) {
	// return if the instance already exists
	if appConfig != nil {
		return appConfig, nil
	}
	// creat new instance
	if err := env.Parse(); err != nil {
		return nil, err
	}

	appConfig = &appConfigs{
		appURI: *goUri,
	}
	return appConfig, nil
}

// gets the URI
func (apconfig *appConfigs) GetAppURI() *string {
	return &apconfig.appURI
}
