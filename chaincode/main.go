/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("MultiOrgs-Chaincode")

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// Init method: one time initialisation
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Init")
	function, args := stub.GetFunctionAndParameters()
	logger.Infof("Invoke is running " + function)

	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var err error

	fmt.Println("Value of Args is: ", args)
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	if function != "init" {
		return shim.Error("Unknown function call")
	}
	test := strings.Join(args, ", ")
	fmt.Println("Value of args is: ", test)

	// Initialize the chaincode
	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	B = args[2]
	Bval, err = strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	fmt.Printf("A = %s, Aval = %d, B = %s, Bval = %d\n", A, Aval, B, Bval)

	// Write the state to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Invoke method: used to send the request to the various custom methods
// All the future requests with name "invoke" will arive here
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	//fmt.Println("Inside Invoke")
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("Invoke is running " + function)

	if args[0] == "createUser" {

		fmt.Println("Create User Function Called")
		return t.createUser(stub, args)

	} else if args[0] == "updateUserData" {

		fmt.Println("Update User Data Function Called")
		return t.updateUserData(stub, args)

	} else if args[0] == "readUser" {

		fmt.Println("Read User Function Called")
		return t.readUser(stub, args)

	} else if args[0] == "readAllUser" {
		fmt.Println("Read All User Function Called")
		return t.readAllUser(stub, args)

	} else if args[0] == "readHistory" {
		fmt.Println("Read History Data Function Called")
		return t.readHistory(stub, args)

	} else if args[0] == "deleteUser" {
		fmt.Println("Delete User Function Called")
		return t.deleteUser(stub, args)

	} else if args[0] == "invoke" {
		logger.Infof("Inside Invoke/invoke")
		// Make payment of X units from A to B
		return t.invoke(stub, args)

	} else if args[0] == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)

	} else if args[0] == "query" {
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)

	} else if args[0] == "testInvoke" {
		eventID := args[1]
		fmt.Println(" #####  Test Event  - " + eventID)
		err := stub.SetEvent(eventID, []byte{})
		if err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(nil)
	}

	return shim.Error("Invalid invoke function name. Available methods are createUser, updateUserData, readUser, readAllUser, readHistory, deleteUser, invoke, delete, query")
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var X int          // Transaction value
	var err error

	fmt.Println("Value of args in invoke is: ", strings.Join(args, ", "))
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	A = args[0]
	B = args[1]

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Bvalbytes == nil {
		return shim.Error("Entity not found")
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	// Perform the execution
	X, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	Aval = Aval - X
	Bval = Bval + X
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state back to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

// query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("Inside query method")
	var A string // Entities
	var err error

	fmt.Println("\value of args in query method is: ", strings.Join(args, ", "), "and length is: ", len(args))
	fmt.Println("Value of arg[0] is: ", args[0], " and arg[1] is: ", args[1])
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[1]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
