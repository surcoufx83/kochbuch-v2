package api

type loginCallbackObject struct {
	State string `json:"state" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

/* // Handles the /me endpoint to query own users profile
func GetMyProfile(c *gin.Context) {
	code, _, profile, _ := services.GetSelf(c)
	if code != http.StatusOK {
		c.String(code, "")
	} else {
		c.JSON(code, profile)
	}
}

func PostLogout(c *gin.Context) {
	code, _ := services.Logout(c)
	c.Status(code)
}

// Handles the /login endpoint containing the Nextcloud Login code
// synchronously retrieve access_token
func PostOauth2Login(c *gin.Context) {
	var payload loginCallbackObject

	err := c.BindJSON(&payload)
	if err != nil {
		log.Printf("Failed to parse input: %v", err)
		c.String(http.StatusBadRequest, "")
		return
	}

	state, _, err := services.GetApplicationParams(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	if state != payload.State {
		log.Printf("Requested state does not match cookie state: %v    vs   %v", state, payload.State)
		c.String(http.StatusBadRequest, "")
		return
	}

	r, err := services.NcLoginCallback(payload.State, payload.Code)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	if r {
		c.String(http.StatusAccepted, "")
	} else {
		c.String(http.StatusInternalServerError, "")
	}

} */
