package main

import (
	"net/http"

	"github.com/parthoshuvo/authsvc/cfg"
	log "github.com/parthoshuvo/authsvc/log4u"
	"github.com/parthoshuvo/authsvc/resource"
	"github.com/parthoshuvo/authsvc/route"
)

func main() {
	config := cfg.NewConfig(version)
	defer config.CloseLog()
	log.Infof("Starting %s on %s\n", config.AppName(), config.Server())
	// TODO uncomment the database connection here
	// replace ??? with appropriate name
	// ??? := db.NewDB(config.GetDbDef())
	// defer ???.Close()
	rb := route.NewRouteBuilder(config.AllowCORS(), new(resource.DefaultProtector), config.AppName(), config.IsLogDebug())
	rb.Add("Home", "GET", "/", resource.HomeHandler(config.HomePage()))

	log.Fatal(http.ListenAndServe(config.Server().String(), rb.Router()))
}
