package main

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// isAuthed checks if the user is authenticated and returns true if the user is allowed to proceed with their action,
// However, it returns false if the user is not allowed to proceed with their action.
func isAuthed(c *gin.Context) bool {
	// Read from headers if token provided

	// If the token is valid, return true
	authTokenHeader := c.Request.Header.Get("Authorization")
	// fmt.Println(authTokenHeader)
	authTokens := returnAuthTokens()

	if authTokenHeader == "" {

		// log
		log.Warnf("AUTH: ACCESS DENIED. Reason: No auth token provided by user IP: %s", c.ClientIP())

		// deny user
		return false
	}

	for _, token := range authTokens {
		if token == authTokenHeader { // Token is valid

			// log
			log.Warnf("AUTH: ACCESS GRANTED. Reason: Auth token provided by user IP: %s", c.ClientIP())

			// allow user to continue
			return true
		} else {
			// log
			log.Warnf("AUTH: ACCESS DENIED. Reason: Invalid auth token provided by user IP: %s", c.ClientIP())

			// deny user
			return false
		}
	}

	return false

}

// returnAuthTokens returns every available auth token from memory
func returnAuthTokens() []string {
	// Read from memory
	return authTokens
}

// @Summary Create a session token
// @Description Create a session token to interact with endpoints that require authentication
// @Produce json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 200 "Session created"
// @Failure 401 "No/Incorrect username OR password provided"
// @Failure 500
// @Tags apiAuthGroup
// @Router /auth/create/session/ [post]
// createSession creates a session token and stores it for the user to authenticate later on.
// It returns a boolean to indicate if the session was created successfully and a token for the user to use.
// If the session was not created successfully, it returns false and the token string will be empty. Furthermore, it
// returns a string to indicate the reason for the failure of the session creation
// If the session was created successfully, it returns true and the token string is the session token.
func createSession(c *gin.Context) (success bool, failReason string, token string) {
	// Create a session for the user
	// todo: make this use a more secure method of password storage
	// todo: make sessions store the time created
	// todo: store sessions on disk, not in memory

	// log
	log.Infof("AUTH: Creating session for user IP: %s", c.ClientIP())

	// get the password
	username := c.PostForm("username")
	password := c.PostForm("password")

	// check if username OR password is blank
	if password == "" || username == "" {
		// log
		log.Warnf("AUTH: ACCESS DENIED. Reason: No username OR password provided by user IP: %s", c.ClientIP())

		if username == "" {
			c.JSON(401, gin.H{
				"status": "error",
				"reason": "No username provided",
			})
		} else {
			c.JSON(401, gin.H{
				"status": "error",
				"reason": "No password provided",
			})
		}

		// deny user
		return false, "No username OR password provided", ""
	}

	// check if the password is correct
	config, err := ConfigParser()
	if err != nil {
		// log
		log.Warnf("AUTH: ACCESS DENIED. Reason: Internal error: Config file could not be parsed. Error: %s", err)
		c.JSON(500, gin.H{
			"status": "error",
			"reason": "Internal error: Config file could not be parsed",
		})
		return false, "Internal error: Config file could not be parsed", ""
	}

	// get the real values
	configUser := config.Auth.Username
	configPass := config.Auth.Password

	if password == configPass && username == configUser {

		// log
		log.Infof("AUTH: ACCESS GRANTED. Reason: Correct username AND password provided by user IP: %s", c.ClientIP())

		// make the whole token
		token = "begin:session:ip=" + c.ClientIP() + ":token=" + GenerateSecureToken() + ":end"

		// add session to memory
		authTokens = append(authTokens, token)

		log.Infof("AUTH: Session created for user IP: %s", c.ClientIP())
		c.JSON(200, gin.H{
			"status": "success",
			"reason": "Session created",
			"token":  token,
		})

		return true, "", token
	} else {

		// log
		log.Warnf("AUTH: ACCESS DENIED. Reason: Incorrect credentials provided by user IP: %s", c.ClientIP())

		// send gin response
		c.JSON(401, gin.H{
			"status": "error",
			"reason": "Incorrect username OR password provided",
		})

		// deny user
		return false, "Incorrect password", ""
	}
}

// GenerateSecureToken generates a secure token for the user to use
func GenerateSecureToken() string {

	// todo: have this verified to be secure

	// always make tokens 40 characters long
	length := 40

	log.Infof("AUTH: Generating secure token. Length: %d", length)

	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}

	return hex.EncodeToString(b)
}

// verifySession checks if the passed session token is valid and returns true if it is valid.
func verifySession(c *gin.Context) bool {

	// todo: does this have any role? disabling for now.
	c.JSON(http.StatusNotImplemented, gin.H{"msg": "This endpoint is not yet implemented and is disabled for security."})
	return false

	// Verify the session token
	passedToken := c.PostForm("token")

	if passedToken == "" {
		// log
		log.Warnf("AUTH: ACCESS DENIED. Reason: No auth token provided by user IP: %s", c.ClientIP())

		// send gin response
		c.JSON(401, gin.H{
			"status": "error",
			"reason": "No auth token provided",
		})

		// deny user
		return false
	}

	// log
	log.Infof("AUTH: Verifying session for user IP: %s", c.ClientIP())

	// check if the token is valid
	authTokens := returnAuthTokens()
	for _, token := range authTokens {
		if token == passedToken {
			// log
			log.Infof("AUTH: TOKEN IS VALID. Reason: A valid session token was provided by user IP: %s", c.ClientIP())

			// send gin response
			c.JSON(200, gin.H{
				"status": "success",
				"reason": "A valid session token was provided",
			})

			// allow user to continue
			return true
		} else {
			// log
			log.Warnf("AUTH: TOKEN IS INVALID. Reason: An invalid session token was provided by user IP: %s", c.ClientIP())

			// send gin response
			c.JSON(401, gin.H{
				"status": "error",
				"reason": "An invalid session token was provided",
			})

			// deny user
			return false
		}
	}

	log.Errorf("AUTH: Unknown error. Function finished without auth result.")
	return false
}
