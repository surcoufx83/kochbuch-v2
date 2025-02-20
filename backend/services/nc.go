package services

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	ValidDomains map[string]bool

	ncRedirectUrl   string
	ncHost          string
	ncClientId      string
	ncSecret        string
	ncAuthEndpoint  string
	ncTokenEndpoint string

	statesMutex sync.RWMutex
	statesCache map[string]dbNextcloudState
)

type AppParams struct {
	LoginUrl string `json:"loginUrl"`
}

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

type nextcloudTokenResponse struct {
	AccessToken  string `json:"access_token" binding:"required"`
	TokenType    string `json:"token_type" binding:"required"`
	ExpiresIn    int    `json:"expires_in" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type dbNextcloudState struct {
	Id           int            `db:"id"`
	State        string         `db:"state"`
	Ip           string         `db:"remoteaddr"`
	Ua           string         `db:"useragent"`
	Created      time.Time      `db:"created"`
	Until        sql.NullTime   `db:"until"`
	Granted      sql.NullTime   `db:"granted"`
	AccessToken  sql.NullString `db:"accesstoken"`
	RefreshToken sql.NullString `db:"refreshtoken"`
	Expires      sql.NullTime   `db:"expires"`
}

func generateState(remoteAddr string, userAgent string) (state string, err error) {
	log.Printf("Generate new state for %v", remoteAddr)

	query := "INSERT INTO `user_login_states`(`remoteaddr`, `useragent`) VALUES(?, ?)"
	stmt, err := Db.Prepare(query)
	if err != nil {
		log.Printf("Failed to prepare user_login_states query: %v", err)
		return "", err
	}

	result, err := stmt.Exec(remoteAddr, userAgent)
	if err != nil {
		log.Printf("Failed to insert into user_login_states: %v", err)
		return "", err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Failed to get insert id from user_login_states: %v", err)
		return "", err
	}

	ncLoadStatesCache()
	statesMutex.RLock()
	defer statesMutex.RUnlock()

	for _, val := range statesCache {
		if val.Id == int(id) {
			return val.State, nil
		}
	}

	return "", nil
}

func GetApplicationParams(c *gin.Context) (newstate string, params AppParams, err error) {
	state, err := c.Cookie("session")

	if state == "" || err != nil {
		// No session cookie
		state, err = generateState(c.Request.RemoteAddr, c.Request.Header.Get("user-agent"))
		if err != nil {
			return "", AppParams{}, err
		}
	}

	statesMutex.RLock()
	clientstate, exists := statesCache[state]
	statesMutex.RUnlock()
	log.Printf("Clientstate: %v", clientstate)
	log.Printf("exists: %v", exists)
	log.Printf("Cookie: %v", state)
	if exists {
		log.Printf("Cookie found in cache: %v", clientstate.State)
		state = clientstate.State
	} else {
		log.Println("Cookie state not in cache")
		state, err = generateState(c.Request.RemoteAddr, c.Request.Header.Get("user-agent"))
		if err != nil {
			return "", AppParams{}, err
		}
	}

	value := AppParams{
		LoginUrl: strings.Replace(ncAuthEndpoint, "_state_", state, -1),
	}
	return state, value, nil
}

func NcConnect() {
	domains := os.Getenv("KB_Domains")
	if domains == "" {
		log.Fatal("Missing ENV variable KB_Domains")
	}

	ValidDomains = make(map[string]bool)
	for _, domain := range strings.Split(domains, ",") {
		ValidDomains[domain] = true
	}
	log.Printf("Registered valid domains %v", ValidDomains)

	ncRedirectUrl = os.Getenv("NC_RedirUrl")
	ncHost = os.Getenv("NC_Host")
	ncClientId = os.Getenv("NC_ClientId")
	ncSecret = os.Getenv("NC_ClientSecret")

	log.Println("Checking Nextcloud instance " + ncHost)

	if ncRedirectUrl == "" {
		log.Fatal("Missing ENV variable NC_RedirUrl")
	}
	ncRedirectUrl = strings.TrimSuffix(ncRedirectUrl, "/") + "/login/oauth2"

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

	baseurl, err := url.Parse("https://" + ncHost + "/")
	if err != nil {
		log.Fatalf("Failed parsing ENV variable NC_Host: %v", err)
	}

	_, err = url.Parse(ncRedirectUrl)
	if err != nil {
		log.Fatalf("Failed parsing ENV variable NC_RedirUrl: %v", err)
	}

	authUrl := *baseurl
	authUrl.Path, _ = url.JoinPath(authUrl.Path, "index.php/apps/oauth2/authorize")
	params := url.Values{}
	params.Add("state", "_state_")
	params.Add("scope", "")
	params.Add("response_type", "code")
	params.Add("approval_prompt", "auto")
	params.Add("redirect_uri", ncRedirectUrl)
	params.Add("client_id", ncClientId)
	authUrl.RawQuery = params.Encode()
	ncAuthEndpoint = authUrl.String()

	tokenUrl := *baseurl
	tokenUrl.Path, _ = url.JoinPath(tokenUrl.Path, "index.php/apps/oauth2/api/v1/token")
	ncTokenEndpoint = tokenUrl.String()

	log.Printf("Nextcloud instance found with version %v", ncStatus.VersionString)
	log.Printf("  - AuthEndpoint = %v", ncAuthEndpoint)
	log.Printf("  - TokenEndpoint = %v", ncTokenEndpoint)
	log.Printf("  - RedirectUrl = %v", ncRedirectUrl)
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
	statesCache = make(map[string]dbNextcloudState)
	for _, state := range states {
		statesCache[state.State] = state
	}
	statesMutex.Unlock()

	log.Printf("Loaded %d login states into cache", len(statesCache))

	go ncDeleteExpiredStates()

}

func ncDeleteExpiredStates() {
	query := "DELETE FROM `user_login_states` WHERE `until` IS NOT NULL AND `until` < current_timestamp()"

	result, err := Db.Exec(query)
	if err != nil {
		log.Printf("Failed to delete outdated user_login_states: %v", err)
	}

	rows, _ := result.RowsAffected()
	if rows > 0 {
		log.Printf("Deleted %d outdated items from user_login_states", len(statesCache))
	}
}

func NcLoginCallback(state string, code string) (bool, error) {
	statesMutex.RLock()
	clientstate, exists := statesCache[state]
	statesMutex.RUnlock()
	if !exists {
		log.Printf("Failed getting state object from cache %v", state)
		return false, nil
	}

	tokenResp, err := ncLoginCallbackSendRequest(code)
	if err != nil {
		return false, err
	}

	clientstate.AccessToken = sql.NullString{String: tokenResp.AccessToken, Valid: true}
	clientstate.Granted = sql.NullTime{Time: time.Now(), Valid: true}
	clientstate.Expires = sql.NullTime{Time: time.Now().Add(time.Second * time.Duration(tokenResp.ExpiresIn)), Valid: true}
	clientstate.RefreshToken = sql.NullString{String: tokenResp.RefreshToken, Valid: true}
	clientstate.Until = sql.NullTime{Valid: false}

	query := "UPDATE `user_login_states` SET `until` = NULL, `granted` = ?, `accesstoken` = ?, `refreshtoken` = ?, `expires` = ? WHERE `id` = ? LIMIT 1"

	stmt, err := Db.Prepare(query)
	if err != nil {
		log.Printf("Failed to prepare user_login_states query: %v", err)
		return false, err
	}

	result, err := stmt.Exec(clientstate.Granted, clientstate.AccessToken, clientstate.RefreshToken, clientstate.Expires, clientstate.Id)
	if err != nil {
		log.Printf("Failed to update user_login_states: %v", err)
		return false, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to get rows count from result: %v", err)
		return false, err
	}
	if count == 0 {
		log.Println("Update did not change any rows")
		return false, err
	}

	statesMutex.Lock()
	statesCache[state] = clientstate
	statesMutex.Unlock()

	go loadUserProfile(clientstate)

	return true, nil
}

func ncLoginCallbackSendRequest(code string) (nextcloudTokenResponse, error) {
	// Prepare request payload
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", ncRedirectUrl)

	req, err := http.NewRequest("POST", ncTokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("Failed retrieving access_token: %v", err)
		return nextcloudTokenResponse{}, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(ncClientId, ncSecret)

	start := time.Now()

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nextcloudTokenResponse{}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read token response body: %v", err)
		return nextcloudTokenResponse{}, err
	}

	// Decode the JSON response into nextcloudTokenResponse
	var tokenResp nextcloudTokenResponse
	if err := json.Unmarshal(bodyBytes, &tokenResp); err != nil {
		log.Printf("Failed to decode token response: %v", err)
		log.Println(string(bodyBytes))
		return nextcloudTokenResponse{}, err
	}

	elapsed := time.Since(start)
	// Extract seconds and milliseconds
	secs := int64(elapsed / time.Second)
	msecs := int64((elapsed % time.Second) / time.Millisecond)
	log.Printf("Nextcloud token response: %v", tokenResp)
	log.Printf("                    took: %d.%d", secs, msecs)
	return tokenResp, nil
}

func loadUserProfile(state dbNextcloudState) {
	log.Printf("Getting user profile from Nextcloud instance for state %v", state.State)

}

func revokeState(id int) {
	log.Printf("Revoke outdated state id %d", id)
}
