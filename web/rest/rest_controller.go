package rest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/MultiOrgs/sdk/org"
)

// PORT : local port
var PORT = "5050"

// RestApp : calling Setup struct
type RestApp struct {
	OrgSetup *org.Setup
}

func hash(s string) string {
	//h := sha1.New()
	h := sha256.New()
	h.Write([]byte(s))
	sha256Hash := hex.EncodeToString(h.Sum(nil))

	return sha256Hash
}

func respondJSON(w http.ResponseWriter, payload interface{}) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}
