package main

import (
	"net/http"

	"github.com/parthoshuvo/authsvc/cache"
	"github.com/parthoshuvo/authsvc/cfg"
	"github.com/parthoshuvo/authsvc/db"
	"github.com/parthoshuvo/authsvc/email"
	log "github.com/parthoshuvo/authsvc/log4u"
	"github.com/parthoshuvo/authsvc/render"
	"github.com/parthoshuvo/authsvc/resource"
	"github.com/parthoshuvo/authsvc/route"
	permTable "github.com/parthoshuvo/authsvc/table/permission"
	roleTable "github.com/parthoshuvo/authsvc/table/role"
	usrTable "github.com/parthoshuvo/authsvc/table/user"
	toknSvc "github.com/parthoshuvo/authsvc/token"
	"github.com/parthoshuvo/authsvc/uc/adm"
	"github.com/parthoshuvo/authsvc/uc/permission"
	"github.com/parthoshuvo/authsvc/uc/role"
	"github.com/parthoshuvo/authsvc/uc/token"
	"github.com/parthoshuvo/authsvc/uc/user"
	"github.com/parthoshuvo/authsvc/validator"
)

func main() {
	config := cfg.NewConfig(version)
	defer config.CloseLog()

	audb := db.NewAuthDB(config.DbDef())
	defer audb.Close()
	tdb := cache.NewTokenDB(config.TokenDBDef())
	defer tdb.Close()
	emailClient := email.NewEmailClient(config.SmtpServerDef())
	defer emailClient.Close()

	validate := validator.New()
	rndr := render.NewJSONRenderer(config.Indent())

	rb := route.NewRouteBuilder(config.AllowCORS(), new(resource.DefaultProtector), config.AppName(), config.IsLogDebug())
	rb.Add("Home", http.MethodGet, "/", resource.HomeHandler(config.HomePage()))

	usrHndlr := user.NewHandler(usrTable.NewTable(audb))
	toknHndlr := token.NewHandler(toknSvc.NewService(config.JWTDef(), tdb))
	roleHndlr := role.NewHandler(roleTable.NewTable(audb))
	permHndlr := permission.NewHandler(permTable.NewTable(audb))

	aurb := rb.SubrouteBuilder("/auth")
	aurs := resource.NewAuthResource(usrHndlr, toknHndlr, rndr, validate, emailClient)
	aurb.Add("LoginUser", http.MethodPost, "/login", aurs.UserLogin())
	aurb.Add("RegisterUser", http.MethodPost, "/register", aurs.UserRegistration())
	aurb.Add("VerifyEmail", http.MethodGet, "/email_verification", aurs.EmailVerifier())

	trb := aurb.SubrouteBuilder("/token")
	trs := resource.NewTokenResource(toknHndlr, adm.NewHandler(usrHndlr, roleHndlr, permHndlr), usrHndlr, rndr)
	trb.Add("VerifyAccessToken", http.MethodPost, "/verify", trs.AccessTokenVerifier())
	trb.Add("GenerateTokenPair", http.MethodPost, "/refresh", trs.TokenPairGenerator())

	log.Infof("Starting %s on %s\n", config.AppName(), config.Server())
	log.Fatal(http.ListenAndServe(config.Server().String(), rb.Router()))
}
