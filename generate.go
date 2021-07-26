package main

import (
	stdlog "log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func setupLogging() {
	log = &logrus.Logger{
		// Output to stdout instead of the default stderr
		// Can be any io.Writer, see below for File example
		Out: ansicolor.NewAnsiColorWriter(os.Stdout),

		// Log as JSON instead of the default ASCII formatter.
		Formatter: &logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		},
		Hooks: make(logrus.LevelHooks),

		// Only log the warning severity or above.
		Level: logrus.InfoLevel,
	}

	w := logrus.New().Writer()
	//defer w.Close()
	stdlog.SetOutput(w)
}

func main() {
	setupLogging()
	log.Info("Generating version.go...")
	f, err := os.Create(path.Join("..", "cmd", "version.go"))
	die(err)
	defer f.Close()

	hostname, err := os.Hostname()
	if err != nil {
		log.Warnf("Unable to get hostname. %v", err)
		hostname = "unknown"
	}

	username := "unknown"
	currentUser, err := user.Current()
	if err != nil {
		log.Warnf("Unable to get username. %v", err)
	} else {
		if currentUser.Username != "" {
			username = currentUser.Username
		}
	}

	version := "unknown"
	versionBytes, err := exec.Command("git", "describe", "--tags").Output()
	if err != nil {
		log.Warnf("Unable to get git revision/tag based version string. %v", err)
	} else {
		version = strings.TrimSpace(string(versionBytes))
	}

	gitStatusBytes, err := exec.Command("git", "status", "--porcelain").Output()
	gitTainted := "unknown"
	if err != nil {
		log.Warnf("Unable to get git status. %v", err)
	} else {
		gitStatus := strings.TrimSpace(string(gitStatusBytes))
		if len(gitStatus) > 0 {
			gitTainted = "tainted"
		} else {
			gitTainted = "clean"
		}
	}

	version = version + "/" + gitTainted

	now := time.Now()

	packageTemplate.Execute(f, struct {
		Timestamp     time.Time
		UnixTimestamp int64
		Username      string
		Hostname      string
		Version       string
	}{
		Timestamp:     now,
		UnixTimestamp: now.Unix(),
		Username:      username,
		Hostname:      hostname,
		Version:       version,
	})
	log.Info("done")
}

var packageTemplate = template.Must(template.New("").Parse(`// Code generated by "go generate"; DO NOT EDIT.
// This file was generated by
// {{ .Username }}@{{ .Hostname }}
// at
// {{ .Timestamp }}
package main

import (
    "fmt"
    "io"
    "time"
)

type VersionInfo struct {
    Timestamp time.Time
    Username  string
    Hostname  string
    Version   string
}

var versionInfo = &VersionInfo{
    Timestamp: time.Unix({{ .UnixTimestamp }}, 0),
    Username:  {{ printf "%q" .Username }},
    Hostname:  {{ printf "%q" .Hostname }},
    Version:   {{ printf "%q" .Version }},
}

func (v *VersionInfo) Print(w io.Writer) {
    fmt.Fprintf(w, "Version:  %s\n", v.Version)
    fmt.Fprintf(w, "Built by: %s@%s\n", v.Username, v.Hostname)
    fmt.Fprintf(w, "Built at: %s\n", v.Timestamp)
}
`))

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
