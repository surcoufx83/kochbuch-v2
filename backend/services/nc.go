package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kochbuch-v2-backend/types"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	ValidDomains map[string]bool

	ncRedirectUrl     string
	ncHost            string
	ncClientId        string
	ncSecret          string
	ncAuthEndpoint    string
	ncTokenEndpoint   string
	ncProfileEndpoint string

	statesMutex sync.RWMutex
	statesCache map[string]dbNextcloudState

	userMutex sync.RWMutex
	userCache map[uint32]*types.UserProfile

	groupsMutex     sync.RWMutex
	groupsCache     map[uint16]types.Group
	groupNamesCache map[string]uint16
)

type AppParams struct {
	LoginUrl string `json:"loginUrl"`
	Session  string `json:"session"`
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
	Id           int              `db:"id"`
	UserId       sql.NullInt32    `db:"userid"`
	State        string           `db:"state"`
	Ip           string           `db:"remoteaddr"`
	Ua           string           `db:"useragent"`
	Created      time.Time        `db:"created"`
	Until        types.NullTime   `db:"until"`
	Granted      types.NullTime   `db:"granted"`
	AccessToken  types.NullString `db:"accesstoken"`
	RefreshToken types.NullString `db:"refreshtoken"`
	Expires      types.NullTime   `db:"expires"`
}

// Nextcloud structure when querying user profile
type nextcloudUserProfileResponse struct {
	OCS struct {
		Meta struct {
			Status     string `json:"status"`
			StatusCode int    `json:"statuscode"`
			Message    string `json:"message"`
		} `json:"meta"`
		Data nextcloudUserProfile `json:"data"`
	} `json:"ocs"`
}

type nextcloudUserProfile struct {
	Id          string   `json:"id"`
	Enabled     bool     `json:"enabled"`
	Email       string   `json:"email"`
	DisplayName string   `json:"displayname"`
	Groups      []string `json:"groups"`
}

type dbUserGroup struct {
	Id      uint32    `db:"id"`
	UserId  uint32    `db:"userid"`
	GroupId uint16    `db:"groupid"`
	Granted time.Time `db:"granted"`
}

func createNewGroup(name string) (int, error) {
	log.Printf("  > Registering new group: %v", name)
	query := "INSERT INTO `groups`(`displayname`, `ncname`) VALUES(?, ?)"

	stmt, err := dbPrepareStmt("createNewGroup", query)
	if err != nil {
		log.Printf("Failed to prepare INSERT INTO groups query: %v", err)
		return 0, err
	}

	result, err := stmt.Exec(name, name)
	if err != nil {
		log.Printf("Failed to INSERT INTO groups: %v", err)
		return 0, err
	}

	groupid, err := result.LastInsertId()
	if err != nil {
		log.Printf("Failed to get id from INSERT INTO groups: %v", err)
		return 0, err
	}

	return int(groupid), nil
}

func generateState(remoteAddr string, userAgent string) (state string, err error) {
	log.Printf("Generate new state for %v", remoteAddr)
	query := "INSERT INTO `user_login_states`(`remoteaddr`, `useragent`) VALUES(?, ?)"

	stmt, err := dbPrepareStmt("generateState", query)
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

func GetApplicationParams(conn *websocket.Conn, state string) (newstate string, params AppParams, err error) {

	if state == "" {
		// No session cookie
		state, err = generateState(conn.RemoteAddr().String(), "wss://")
		if err != nil {
			return "", AppParams{}, err
		}
	}

	statesMutex.RLock()
	clientstate, exists := statesCache[state]
	statesMutex.RUnlock()

	if exists {
		state = clientstate.State
	} else {
		log.Println("Cookie state not in cache")
		state, err = generateState(conn.RemoteAddr().String(), "wss://")
		if err != nil {
			return "", AppParams{}, err
		}
	}

	value := AppParams{
		LoginUrl: strings.Replace(ncAuthEndpoint, "_state_", state, -1),
		Session:  state,
	}
	return state, value, nil
}

func GetSelf(c *gin.Context) (int, string, *types.UserProfile, error) {
	state, err := c.Cookie("session")

	if state == "" || err != nil {
		// No session cookie
		return http.StatusBadRequest, "", &types.UserProfile{}, err
	}

	return GetSelfByState(state)
}

func GetSelfByState(state string) (int, string, *types.UserProfile, error) {
	if state == "" {
		// No session cookie
		return http.StatusBadRequest, "", &types.UserProfile{}, errors.New("invalid request")
	}

	statesMutex.RLock()
	clientstate, exists := statesCache[state]
	statesMutex.RUnlock()
	if !exists {
		log.Printf("Cookie not found in cache: %v", state)
		return http.StatusBadRequest, "", &types.UserProfile{}, nil
	}

	if !clientstate.UserId.Valid || clientstate.UserId.Int32 == 0 {
		log.Printf("Cookie does not has loggedin user: %v", state)
		return http.StatusUnauthorized, "", &types.UserProfile{}, nil
	}

	user, found := userCache[uint32(clientstate.UserId.Int32)]
	if !found {
		return http.StatusUnauthorized, "", &types.UserProfile{}, nil
	}

	return http.StatusOK, state, user, nil
}

func GetUser(id int) (int, *types.UserProfile) {
	user, found := userCache[uint32(id)]
	if !found {
		return http.StatusUnauthorized, &types.UserProfile{}
	}
	return http.StatusOK, user
}

func loadUserProfile(state dbNextcloudState) (*types.UserProfile, error) {
	log.Printf("Getting user profile from Nextcloud instance for state %v", state.State)

	if !state.AccessToken.Valid {
		log.Printf("  - skipped due to invalid access_token %v", state.AccessToken)
		return &types.UserProfile{}, nil
	}

	req, err := http.NewRequest("GET", ncProfileEndpoint, nil)
	if err != nil {
		log.Printf("Failed preparing request for user profile: %v", err)
		return &types.UserProfile{}, err
	}

	// Set required headers
	req.Header.Set("Authorization", "Bearer "+state.AccessToken.String)
	req.Header.Set("OCS-APIRequest", "true")

	log.Printf("Request: %v", req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed sending request for user profile: %v", err)
		return &types.UserProfile{}, err
	}
	defer resp.Body.Close()

	// Check for a successful HTTP status code.
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Unexpected status for user profile: %d  body = %s", resp.StatusCode, string(body))
		return &types.UserProfile{}, err
	}

	var profileResp nextcloudUserProfileResponse
	if err := json.NewDecoder(resp.Body).Decode(&profileResp); err != nil {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Failed to parse user profile: %d  body = %s", resp.StatusCode, string(body))
		return &types.UserProfile{}, err
	}

	log.Printf("User profile: %v", profileResp)
	var userProfile *types.UserProfile
	insertedOne := false

	for _, groupname := range profileResp.OCS.Data.Groups {
		if groupNamesCache[groupname] == 0 {
			_, err := createNewGroup(groupname)
			if err != nil {
				log.Printf("Failed creating new group: %v", err)
				return &types.UserProfile{}, err
			}
			insertedOne = true
		}
	}
	if insertedOne {
		ncLoadGroupsCache()
	}

	tx, err := Db.Begin()
	if err != nil {
		log.Printf("Failed starting transaction: %v", err)
		return &types.UserProfile{}, err
	}

	defer ncLoadUserCache()
	defer ncLoadStatesCache()

	for _, user := range userCache {
		if user.UserName == profileResp.OCS.Data.Id {
			userProfile = user
			break
		}
	}

	if userProfile.Id == 0 {
		userProfile, err = registerUserProfile(tx, profileResp)
		if err != nil || userProfile.Id == 0 {
			_ = tx.Rollback()
			return &types.UserProfile{}, err
		}
	}

	registerUserGroups(tx, userProfile, profileResp.OCS.Data.Groups)
	if err != nil {
		_ = tx.Rollback()
		return &types.UserProfile{}, err
	}

	query := "UPDATE `user_login_states` SET `userid` = ? WHERE `id` = ? LIMIT 1"

	stmt, err := dbPrepareStmt("registerUserGroups", query)
	if err != nil {
		log.Printf("Failed to prepare user_login_states query: %v", err)
		_ = tx.Rollback()
		return &types.UserProfile{}, err
	}
	stmt = tx.Stmt(stmt)

	_, err = stmt.Exec(userProfile.Id, state.Id)
	if err != nil {
		log.Printf("Failed to update user_login_states: %v", err)
		_ = tx.Rollback()
		return &types.UserProfile{}, err
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		log.Printf("Failed to commit transaction: %v", err)
		return &types.UserProfile{}, err
	}

	log.Printf("User profile: %v", userProfile)
	go OnUserProfileUpdated(state.State, userProfile)

	return userProfile, nil

}

func Logout(conn *wsConnection) (int, error) {
	code, state, _, err := GetSelfByState(conn.ConnectionParams.Session)
	if err != nil || code != http.StatusOK || state == "" {
		return http.StatusAccepted, nil
	}

	statesMutex.Lock()
	defer statesMutex.Unlock()

	log.Printf("  > Removing session: %v", state)

	query := "DELETE FROM `user_login_states` WHERE `state` = ? LIMIT 1"

	stmt, err := dbPrepareStmt("Logout", query)
	if err != nil {
		log.Printf("Failed to prepare DELETE FROM user_login_states query: %v", err)
		return http.StatusInternalServerError, err
	}

	result, err := stmt.Exec(state)
	if err != nil {
		log.Printf("Failed to DELETE FROM user_login_states: %v", err)
		return http.StatusInternalServerError, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to get rowcount from DELETE FROM user_login_states: %v", err)
	} else if count != 0 {
		log.Printf("Rowcount == 0 from DELETE FROM user_login_states: %v", err)
	}

	delete(statesCache, state)

	return http.StatusAccepted, nil

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

	profileUrl := *baseurl
	profileUrl.Path, _ = url.JoinPath(profileUrl.Path, "ocs/v2.php/cloud/user")
	params = url.Values{}
	params.Add("format", "json")
	profileUrl.RawQuery = params.Encode()
	ncProfileEndpoint = profileUrl.String()

	log.Printf("Nextcloud instance found with version %v", ncStatus.VersionString)
	log.Printf("  - AuthEndpoint = %v", ncAuthEndpoint)
	log.Printf("  - TokenEndpoint = %v", ncTokenEndpoint)
	log.Printf("  - ProfileEndpoint = %v", ncProfileEndpoint)
	log.Printf("  - RedirectUrl = %v", ncRedirectUrl)

	log.Println("")

	ncLoadGroupsCache()
	ncLoadUserCache()
	ncLoadStatesCache()

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

func ncLoadGroupsCache() {
	// log.Println("Loading registered groups")
	query := "SELECT * FROM `groups`"
	var groups []types.Group

	err := Db.Select(&groups, query)
	if err != nil {
		log.Fatalf("Failed to load groups: %v", err)
	}

	// Build cache
	groupsMutex.Lock()
	groupsCache = make(map[uint16]types.Group)
	groupNamesCache = make(map[string]uint16)
	for _, group := range groups {
		groupsCache[group.Id] = group
		groupNamesCache[group.Name] = group.Id
		// log.Printf("  - Loaded %v", group.Name)
	}
	groupsMutex.Unlock()

	log.Printf("Loaded %d groups into cache", len(groupsCache))

}

func ncLoadStatesCache() {
	//log.Println("Loading existing states (cookies)")
	query := "SELECT * FROM `user_login_states` WHERE `until` > current_timestamp() OR `granted` IS NOT NULL"
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
		// log.Printf("  - Loaded %v", state.State)
	}
	statesMutex.Unlock()

	log.Printf("Loaded %d login states into cache", len(statesCache))

	go ncDeleteExpiredStates()

}

func ncLoadUserCache() {
	//log.Println("Loading registered users")
	query := "SELECT * FROM `users`"
	var profiles []types.UserProfile

	err := Db.Select(&profiles, query)
	if err != nil {
		log.Fatalf("Failed to load users: %v", err)
	}

	// Build cache
	userMutex.Lock()
	userCache = make(map[uint32]*types.UserProfile)
	for _, profile := range profiles {
		profile.SimpleProfile = types.NullUserProfileSimple{
			Profile: types.UserProfileSimple{
				Id:          &profile.Id,
				DisplayName: &profile.DisplayName,
			},
			Valid: true,
		}
		userCache[profile.Id] = &profile
		// log.Printf("  - Loaded %v", profile.UserName)
	}

	log.Printf("Loaded %d users into cache", len(userCache))

	ncLoadUserGroupsCache()
	loadCollections()

	userMutex.Unlock()
}

func ncLoadUserGroupsCache() {
	//log.Print("Loading usergroups")
	query := "SELECT * FROM `user_groups`"
	var groups []dbUserGroup

	err := Db.Select(&groups, query)
	if err != nil {
		log.Fatalf("Failed to load user_groups: %v", err)
	}

	for _, group := range groups {
		user := userCache[group.UserId]
		user.Groups = append(user.Groups, groupsCache[group.GroupId])
		userCache[group.UserId] = user
		// log.Printf("  - %s > %s", user.UserName, groupsCache[group.GroupId].Name)
	}

	log.Printf("Loaded %d user groups into cache", len(groups))
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

	clientstate.AccessToken = types.NullString{String: tokenResp.AccessToken, Valid: true}
	clientstate.Granted = types.NullTime{Time: time.Now(), Valid: true}
	clientstate.Expires = types.NullTime{Time: time.Now().Add(time.Second * time.Duration(tokenResp.ExpiresIn)), Valid: true}
	clientstate.RefreshToken = types.NullString{String: tokenResp.RefreshToken, Valid: true}
	clientstate.Until = types.NullTime{Valid: false}

	query := "UPDATE `user_login_states` SET `until` = NULL, `granted` = ?, `accesstoken` = ?, `refreshtoken` = ?, `expires` = ? WHERE `id` = ? LIMIT 1"

	stmt, err := dbPrepareStmt("NcLoginCallback", query)
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

func registerUserGroups(tx *sql.Tx, profile *types.UserProfile, groups []string) (bool, error) {
	log.Printf("  > Saving user groups: %v", groups)

	stmt, err := dbPrepareStmt("registerUserGroups_Delete", "DELETE FROM `user_groups` WHERE `userid` = ?")
	if err != nil {
		log.Printf("Failed to prepare DELETE FROM user_groups query: %v", err)
		return false, err
	}
	stmt = tx.Stmt(stmt)

	_, err = stmt.Exec(profile.Id)
	if err != nil {
		log.Printf("Failed to execute DELETE FROM user_groups query: %v", err)
		return false, err
	}

	for _, groupname := range groups {
		groupid := groupNamesCache[groupname]
		if groupid == 0 {
			log.Printf("Failed to get group: %v", err)
			return false, err
		}
		group := groupsCache[groupid]

		stmt, err := dbPrepareStmt("registerUserGroups_Insert", "INSERT INTO `user_groups`(`userid`, `groupid`) VALUES(?, ?)")
		if err != nil {
			log.Printf("Failed to prepare INSERT INTO user_groups query: %v", err)
			return false, err
		}
		stmt = tx.Stmt(stmt)

		_, err = stmt.Exec(profile.Id, group.Id)
		if err != nil {
			log.Printf("Failed to execute INSERT INTO user_groups query: %v", err)
			return false, err
		}

	}

	return true, nil

}

func registerUserProfile(tx *sql.Tx, profile nextcloudUserProfileResponse) (*types.UserProfile, error) {
	log.Printf("  > Creating new user profile: %v", profile.OCS.Data.Id)

	query := "INSERT INTO `users`(`cloudid`, `clouddisplayname`, `cloudenabled`, `cloudsync`, `cloudsync_status`, `enabled`) VALUES(?, ?, ?, CURRENT_TIMESTAMP(), 200, ?)"

	stmt, err := dbPrepareStmt("registerUserProfile", query)
	if err != nil {
		log.Printf("Failed to prepare INSERT INTO users query: %v", err)
		return &types.UserProfile{}, err
	}
	stmt = tx.Stmt(stmt)

	result, err := stmt.Exec(profile.OCS.Data.Id, profile.OCS.Data.DisplayName, profile.OCS.Data.Enabled, profile.OCS.Data.Enabled)
	if err != nil {
		log.Printf("Failed to INSERT INTO users: %v", err)
		return &types.UserProfile{}, err
	}

	userid, err := result.LastInsertId()
	if err != nil {
		log.Printf("Failed to get id from INSERT INTO users: %v", err)
		return &types.UserProfile{}, err
	}

	user, found := userCache[uint32(userid)]
	if !found {
		return &types.UserProfile{}, errors.New("not found")
	}

	return user, nil
}

func setUserModified(user *types.UserProfile, time time.Time) {
	fn := fmt.Sprintf("setUserModified(%d, %s, %s)", user.Id, user.DisplayName, time)
	log.Print(fn)
	query := "UPDATE `users` SET `modified` = ? WHERE `user_id` = ?"

	stmt, err := dbPrepareStmt("setUserModified", query)
	if err != nil {
		log.Printf("%s: Failed to prepare stmt: %v", fn, err)
		return
	}

	_, err = stmt.Exec(time, user.Id)
	if err != nil {
		log.Printf("%s: Failed to exec stmt: %v", fn, err)
		return
	}

	user.Modified = time
	go wsWelcomeAgain(user)
}
