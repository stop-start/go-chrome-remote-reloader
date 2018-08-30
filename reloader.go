package reloader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"golang.org/x/net/websocket"
)

// Parameters  struct ...
type Parameters struct {
	IgnoreCache bool `json:"ignoreCache"`
}

// RefreshJSON struct ...
type RefreshJSON struct {
	ID     uint16     `json:"id"`
	Method string     `json:"method"`
	Params Parameters `json:"params"`
}

// ChromeTab is the json returned from chrome remote debugging api
type ChromeTab struct {
	ID                   string `json:"id"`
	Title                string `json:"title"`
	Type                 string `json:"type"`
	URL                  string `json:"url"`
	WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
	Description          string `json:"description"`
}

// RemoteConfig is the configuration used to start chrome remote debugging
type RemoteConfig struct {
	// ExecName is the name of the chrome executable (e.g "google-chrome", "chromium")
	ExecName string
	// Port is the chrome remote debugging port, usually 9222
	Port int
	// UserDataDir is the folder for chrome user profile
	UserDataDir string
	// OriginAddr is the address to be reloaded
	OriginAddr string
	// OriginPort is the port for the address to be reloaded
	OriginPort int
}

// RemoteChrome starts a new chrome remote debugging session
func RemoteChrome(rc *RemoteConfig) error {
	return rc.remoteChrome()
}

// RemoteChromeDefault starts a new chrome remote debugging session with default port and user data directory
func RemoteChromeDefault() error {
	rc := &RemoteConfig{
		ExecName:    "chromium", //TODO: change to google-chrome
		Port:        9222,
		UserDataDir: "~/.chrome-remote-profile",
		OriginAddr:  "localhost",
		OriginPort:  8080,
	}
	return rc.remoteChrome()
}

func (rc *RemoteConfig) remoteChrome() error {
	cmd := fmt.Sprintf("%s --remote-debugging-port=%d --user-data-dir=%s http://%s:%d",
		rc.ExecName, rc.Port, rc.UserDataDir, rc.OriginAddr, rc.OriginPort)

	go func() {

	}()

	return nil
}

// ReloadTab reloads one chrome tab
func ReloadTab() error {
	apiURL := fmt.Sprintf("http://localhost:%d/json", 9222)
	resp, err := http.Get(apiURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var tabs []ChromeTab
	buffer, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(buffer, &tabs)
	if err != nil {
		return err
	}

	//for _, tab := range tabs {
	url := tabs[0].WebSocketDebuggerURL
	origin := tabs[0].URL
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		return err
	}

	jsonStruct := &RefreshJSON{ID: 0, Method: "Page.reload", Params: Parameters{IgnoreCache: true}}
	jsonString, _ := json.Marshal(jsonStruct)
	_, err = ws.Write(jsonString)
	if err != nil {
		return err
	}
	//}

	//  func Post(url, contentType string, body io.Reader) (resp *Response, err error)
	// _, err = http.Post(url.String(), "application/json", strings.NewReader(string(jsonString)))
	return nil
}
