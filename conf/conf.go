package conf

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/comail/colog"
)

var (
	// Name of the project
	ProjectName string

	// Commit SHA and version for the current build, set by the
	// compile process.
	Version  string
	Revision string
)

type Config struct {
	// Current environment (e.g. 'production', 'debug').  This is set from
	// an environment variable.
	Environment string `json:"-"`

	// Web configuration
	Host          string `json:"host"`
	Port          uint16 `json:"port"`
	SessionSecret string `json:"session_secret"`

	// DB configuration
	DbType string `json:"dbtype"`
	DbConn string `json:"dbconn"`
}

func (c *Config) HostString() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *Config) IsDebug() bool {
	return c.Environment == "debug"
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production" || c.Environment == "prod"
}

var (
	ConfigPath      = filepath.Join(".", "config.json")
	configPathGiven bool

	C = &Config{}
)

func init() {
	if ProjectName == "" {
		panic("no project name set - did you use the Makefile to build?")
	}

	// Set defaults
	C.Environment = "debug"
	C.Host = "localhost"
	C.Port = 3001
	C.DbType = "sqlite3"
	C.DbConn = ":memory:"

	// Load the environment
	C.Environment = os.Getenv("ENVIRONMENT")

	// Register colog to handle the standard logger output
	colog.Register()
	colog.SetFlags(0)
	colog.ParseFields(true)

	// Set up logger
	if C.IsDebug() {
		colog.SetMinLevel(colog.LDebug)
	} else {
		colog.SetMinLevel(colog.LInfo)
		colog.SetFormatter(&colog.JSONFormatter{})
	}

	// Generate a random session secret.
	buf := make([]byte, 20)
	if _, err := rand.Read(buf); err != nil {
		log.Printf("error: could not generate random secret: %s\n", err)
		os.Exit(1)
		return
	}
	C.SessionSecret = hex.EncodeToString(buf)

	// Let the user override the config file path.
	if cpath := os.Getenv(strings.ToUpper(ProjectName) + "_CONFIG_PATH"); cpath != "" {
		ConfigPath = cpath
		configPathGiven = true
	}

	// Read the configuration file, if present.
	f, err := os.Open(ConfigPath)
	if err != nil {
		// We don't print an error if the user did not give a config path, and
		// the default config file does not exist.
		if !configPathGiven && os.IsNotExist(err) {
			// Do nothing
		} else {
			log.Printf("error: could not read configuration file `%s`: %s\n", ConfigPath, err)
		}
		return
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(C); err != nil {
		log.Printf("error: could not decode configuration file: %s\n", err)
	}
}
