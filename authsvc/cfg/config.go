package cfg

import (
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	log "github.com/parthoshuvo/authsvc/log4u"
	usrTable "github.com/parthoshuvo/authsvc/table/user"
	"github.com/parthoshuvo/authsvc/token"
)

const defaultConfigFilePath = "authsvc.json"
const defaultLogLevel = "DEBUG"

// Config holds configuration data.
type Config struct {
	configData *configData
	logFile    *os.File
	logDebug   bool
	appName    string
}

// ServerDef defines a server address and port.
type ServerDef struct {
	Bind string
	Port int
}

// DBDef database definition
type DBDef struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}

// TokenDBDef database defintion
type TokenDBDef struct {
	Host     string
	Port     int
	Password string
	Database int
}

type SmtpServerDef struct {
	Host string
	Port int
	From usrTable.Email
}

// logDef defines logging
type logDef struct {
	Filename string
	Level    string
}

// configData defines the authsvc configuration file structure.
type configData struct {
	Name        string
	Description string
	AllowCORS   bool
	Server      *ServerDef
	DB          *DBDef
	TokenDB     *TokenDBDef
	JWTDef      *token.JWTDef
	SmtpServer  *SmtpServerDef
	Logging     *logDef
	Indent      bool
}

// NewConfig creates the application configuration.
func NewConfig(version string) *Config {
	cd := loadConfig()
	lf := configureLogging(cd.Logging.Filename, cd.Logging.Level)
	ld := cd.Logging.isDebug()
	an := cd.appName(version)
	return &Config{cd, lf, ld, an}
}

// AppName provides the application name and version.
func (c *Config) AppName() string {
	return c.appName
}

// AllowCORS determines whether cross origin calls are allowed.
func (c *Config) AllowCORS() bool {
	return c.configData.AllowCORS
}

// Server returns the address and port to use for this service.
func (c *Config) Server() *ServerDef {
	return c.configData.Server
}

// DbDef returns the database definition.
func (c *Config) DbDef() *DBDef {
	return c.configData.DB
}

// TokenDBDef returns the token database definition.
func (c *Config) TokenDBDef() *TokenDBDef {
	return c.configData.TokenDB
}

// JWTDef returns JWT token configuration of access and refresh tokens
func (c *Config) JWTDef() *token.JWTDef {
	return c.configData.JWTDef
}

// SmtpServerDef returns SMTP mail server definition
func (c *Config) SmtpServerDef() *SmtpServerDef {
	return c.configData.SmtpServer
}

// IsLogDebug indicates whether debug logging is wanted.
func (c *Config) IsLogDebug() bool {
	return c.logDebug
}

// Indent determines whether JSON and XML renderers indent output.
func (c *Config) Indent() bool {
	return c.configData.Indent
}

// CloseLog closes the log file.
func (c *Config) CloseLog() {
	if c.logFile != nil {
		c.logFile.Close()
	}
}

// HomePage renders the authsvc configuration.
func (c *Config) HomePage() string {
	return "<html>" +
		"<head><title>" + c.configData.Name + " Service</title></head>" +
		"<body><dl>" +
		render("name", c.configData.Name) +
		render("description", c.configData.Description) +
		render("version", c.AppName()) +
		render("server", c.Server().String()) +
		render("log file", c.configData.Logging.Filename) +
		render("log level", c.configData.Logging.Level) +
		render("indent", strconv.FormatBool(c.configData.Indent)) +
		"</dl></body>" +
		"</html>"
}

func render(label, data string) string {
	return fmt.Sprintf("<dt><b>%s</b></dt><dd>%s</dd>", html.EscapeString(label), html.EscapeString(data))
}

func loadConfig() *configData {
	flag.Parse()
	configFilePath := flag.Arg(0)
	if configFilePath == "" {
		configFilePath = defaultConfigFilePath
	}
	log.Infof("Reading configuration from %s\n", configFilePath)
	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("failed to read file %s: %v", configFilePath, err)
	}
	var cfgData configData
	if err := json.Unmarshal(data, &cfgData); err != nil {
		log.Fatalf("failed to unmarshal from file %s: %v", configFilePath, err)
	}
	return &cfgData
}

func configureLogging(filename, level string) *os.File {
	var logFile *os.File
	var err error
	if filename == "" {
		log.SetLevel(defaultLogLevel)
	} else {
		logFile, err = os.Create(filename)
		if err != nil {
			log.Fatalf("failed to create file %s: %v", filename, err)
		}
		logger := io.MultiWriter(os.Stderr, logFile)
		log.SetLevel(level)
		log.SetOutput(logger)
	}
	return logFile
}

func (cd *configData) appName(version string) string {
	return fmt.Sprintf("%s/%s", cd.Name, version)
}

func (dd *DBDef) String() string {
	return fmt.Sprintf("%s:%s:%d:%s", dd.User, dd.Host, dd.Port, dd.Database)
}

func (sd *ServerDef) String() string {
	return fmt.Sprintf("%s:%d", sd.Bind, sd.Port)
}

func (ld *logDef) isDebug() bool {
	return strings.EqualFold(ld.Level, "DEBUG")
}
