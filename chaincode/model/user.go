package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
)

const (
	// CO1 : custom org
	CO1 = "Org1"
	// CO2 : custom org
	CO2 = "Org2"
)

const (
	// COLLECTION_KEY : used to fetch the users
	COLLECTION_KEY = "email"
	// ADMIN : role
	ADMIN = "admin"
)

// GetCustomOrgName : get org name
func GetCustomOrgName(org string) string {
	fmt.Println(" ### GetCustomOrgName = " + org)
	if strings.EqualFold(org, "org1") {
		return CO1
	} else if strings.EqualFold(org, "org2") {
		return CO2
	}
	return CO1
}

// User : storing user info
type User struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
	//	Age         string `json:"age"`
	//	Salary      string `json:"salary"`
	Owner string `json:"owner"`
	//	ShareAccess string `json:"shareAccess"`
	Role string `json:"role"`
	Time string `json:"time"`
	//	Remarks     string `json:"remarks"`
}

// HistoryData : get the txn history
type HistoryData struct {
	EmailKey        string `json:"emailKey"`
	TxID            string `json:"txId"`
	QueryCreator    string `json:"creator"`
	Query           string `json:"query"`
	QueryCreatorOrg string `json:"queryCreatorOrg"`
	Time            string `json:"time"`
	Remarks         string `json:"remarks"`
}

// IsAdmin : if user is admin or not
func IsAdmin(role string) bool {
	return strings.EqualFold(role, ADMIN)
}

// GetTime : get time
func GetTime(timestamp *timestamp.Timestamp) string {

	t := time.Unix(timestamp.GetSeconds(), 0)

	return t.Format(time.RFC1123)
}
