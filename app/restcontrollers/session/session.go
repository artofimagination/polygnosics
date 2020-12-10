package session

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

var (
	// Store key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")), []byte(os.Getenv("SESSION_ENCRYPTION")))
)

func EncryptUserAndOrigin(userid uuid.UUID, origin string) (string, error) {
	data := make(map[string]interface{})
	data["userid"] = userid.String()
	data["origin"] = strings.Split(origin, ":")[0]

	binary, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(binary)

	return encoded, nil
}

func matchingUserAndOrigin(userid uuid.UUID, origin string, cookieData string) (bool, error) {
	data := make(map[string]interface{})
	data["userid"] = userid.String()
	data["origin"] = strings.Split(origin, ":")[0]

	binary, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	encoded := base64.StdEncoding.EncodeToString(binary)
	if encoded != cookieData {
		return false, nil
	}

	return true, nil
}

func IsAuthenticated(userID uuid.UUID, session *sessions.Session, r *http.Request) (bool, error) {
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		return false, nil
	}

	cookieKey, ok := session.Values["cookie_key"].(string)
	if !ok {
		return false, errors.New("Failed to decode cookie key")
	}

	match, err := matchingUserAndOrigin(userID, r.RemoteAddr, cookieKey)
	if err != nil {
		return false, err
	}

	return match, nil
}
