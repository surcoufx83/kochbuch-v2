package services

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

var (
	ncRedirectUrl   string
	ncHost          string
	ncClientId      string
	ncSecret        string
	ncAuthEndpoint  string
	ncTokenEndpoint string

	statesMutex sync.RWMutex
	statesCache map[int]dbNextcloudState
)

type NextcloudStatus struct {
	Installed       bool   `json:"installed"`
	Maintenance     bool   `json:"maintenance"`
	NeedsDbUpgrade  bool   `json:"needsDbUpgrade"`
	Version         string `json:"version"`
	VersionString   string `json:"versionstring"`
	Edition         string `json:"edition"`
	ProductName     string `json:"productname"`
	ExtendedSupport bool   `json:"extendedSupport"`
}

type dbNextcloudState struct {
	Id      int       `db:"id"`
	State   string    `db:"state"`
	Ipagent string    `db:"ipagent"`
	Created time.Time `db:"created"`
	Until   time.Time `db:"until"`
}

func NcConnect() {
	ncRedirectUrl = os.Getenv("NC_RedirUrl")
	ncHost = os.Getenv("NC_Host")
	ncClientId = os.Getenv("NC_ClientId")
	ncSecret = os.Getenv("NC_ClientSecret")

	log.Println("Checking Nextcloud instance " + ncHost)

	if ncRedirectUrl == "" {
		log.Fatal("Missing ENV variable NC_RedirUrl")
	}
	ncRedirectUrl = ncRedirectUrl + "/login/oauth2"

	if ncHost == "" {
		log.Fatal("Missing ENV variable NC_Host")
	}

	if ncClientId == "" {
		log.Fatal("Missing ENV variable NC_ClientId")
	}

	if ncSecret == "" {
		log.Fatal("Missing ENV variable NC_ClientSecret")
	}

	resp, err := http.Get("https://" + ncHost + "/status.php")
	if err != nil {
		log.Fatalf("Failed loading Nextcloud status: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed loading Nextcloud status! Response status code is %v", resp.StatusCode)
	}

	var ncStatus NextcloudStatus
	if err := json.NewDecoder(resp.Body).Decode(&ncStatus); err != nil {
		log.Fatalf("Failed to decode Nextcloud status JSON: %v", err)
	}

	_, err = url.Parse("https://" + ncHost + "/")
	if err != nil {
		log.Fatalf("Failed parsing ENV variable NC_Host: %v", err)
	}

	_, err = url.Parse(ncRedirectUrl)
	if err != nil {
		log.Fatalf("Failed parsing ENV variable NC_RedirUrl: %v", err)
	}

	log.Printf("Nextcloud instance found with version %v", ncStatus.VersionString)
	ncLoadStatesCache()

}

func ncLoadStatesCache() {
	query := "SELECT * FROM `user_login_states` WHERE `until` > current_timestamp()"
	var states []dbNextcloudState

	err := Db.Select(&states, query)
	if err != nil {
		log.Fatalf("Failed to load user_login_states: %v", err)
	}

	// Build cache
	statesMutex.Lock()
	statesCache = make(map[int]dbNextcloudState)
	for _, state := range states {
		statesCache[state.Id] = state
	}
	statesMutex.Unlock()

	log.Printf("Loaded %d login states into cache", len(statesCache))

	query = "DELETE FROM `user_login_states` WHERE `until` < current_timestamp()"

	result, err := Db.Exec(query)
	if err != nil {
		log.Printf("Failed to delete outdated user_login_states: %v", err)
	}

	rows, _ := result.RowsAffected()
	if rows > 0 {
		log.Printf("Deleted %d outdated items from user_login_states", len(statesCache))
	}

}
