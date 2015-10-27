package conf

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	// General configuration
	Debug bool `json:"debug"`

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

var (
	ConfigPath = filepath.Join(".", "config.json")
	C          = &Config{}
)

func init() {
	if ProjectName == "" {
		panic("no project name set - did you use the Makefile to build?")
	}

	// Set defaults
	C.Debug = false
	C.Host = "0.0.0.0"
	C.Port = 3001
	C.DbType = "sqlite3"
	C.DbConn = ":memory:"

	// Generate a random session secret.
	buf := make([]byte, 20)
	if _, err := rand.Read(buf); err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not generate random secret: %s\n", err)
		os.Exit(1)
		return
	}
	C.SessionSecret = hex.EncodeToString(buf)

	// Let the user override the config file path.
	if cpath := os.Getenv(strings.ToUpper(ProjectName) + "_CONFIG_PATH"); cpath != "" {
		ConfigPath = cpath
	}

	// Read the configuration file, if present.
	f, err := os.Open(ConfigPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not read configuration file `%s`: %s\n", ConfigPath, err)
		return
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(C); err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not decode configuration file: %s\n", err)
	}
}
