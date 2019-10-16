package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// get the data from the ledger
func getDataFromLedger(stub shim.ChaincodeStubInterface, key string, result interface{}) error {

	resultAsByte, err := stub.GetState(key)
	if err != nil {
		return fmt.Errorf("Unable to get the data " + err.Error())
	}
	if resultAsByte == nil {
		return fmt.Errorf("the object doesn't exist in the ledger")
	}

	err = byteToObject(resultAsByte, result)
	if err != nil {
		return fmt.Errorf("unable to convert the result to object: %v", err)
	}
	return nil
}

// delete the data from the ledger
func deleteDataFromLedger(stub shim.ChaincodeStubInterface, key string) error {

	err := stub.DelState(key)
	if err != nil {
		return fmt.Errorf("unable to delete the object in the ledger: %v", err)
	}
	return nil
}

// marshelling the JSONobject to bytes
func objectToByte(object interface{}) ([]byte, error) {
	objectAsByte, err := json.Marshal(object)
	if err != nil {
		return nil, fmt.Errorf("unable convert the object to byte: %v", err)
	}
	return objectAsByte, nil
}

// unmarshelling the bytes to JSON object
func byteToObject(objectAsByte []byte, result interface{}) error {
	err := json.Unmarshal(objectAsByte, result)
	if err != nil {
		return fmt.Errorf("unable to convert the result to object: %v", err)
	}
	return nil
}
