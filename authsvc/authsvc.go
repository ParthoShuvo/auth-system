package main

import (
	"net/http"

	"github.com/parthoshuvo/authsvc/cfg"
	"github.com/parthoshuvo/authsvc/db"
	log "github.com/parthoshuvo/authsvc/log4u"
	"github.com/parthoshuvo/authsvc/render"
	"github.com/parthoshuvo/authsvc/resource"
	"github.com/parthoshuvo/authsvc/route"
	usrTable "github.com/parthoshuvo/authsvc/table/user"
	"github.com/parthoshuvo/authsvc/uc/user"
	"github.com/parthoshuvo/authsvc/validator"
)

func main() {
	config := cfg.NewConfig(version)
	defer config.CloseLog()
	log.Infof("Starting %s on %s\n", config.AppName(), config.Server())

	audb := db.NewAuthDB(config.GetDbDef())
	defer audb.Close()
	validate := validator.New()
	rndr := render.NewJSONRenderer(config.Indent())

	rb := route.NewRouteBuilder(config.AllowCORS(), new(resource.DefaultProtector), config.AppName(), config.IsLogDebug())
	rb.Add("Home", http.MethodGet, "/", resource.HomeHandler(config.HomePage()))

	aurb := rb.SubrouteBuilder("/auth")
	aurs := resource.NewAuthResource(user.NewHandler(usrTable.NewTable(audb)), rndr, validate)
	aurb.AddSafe("UserLogin", http.MethodPost, "/login", aurs.UserLogin())

	log.Fatal(http.ListenAndServe(config.Server().String(), rb.Router()))
}
