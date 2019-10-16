package main

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var testLog = shim.NewLogger("MultiOrgs-Chaincode_test")

// checkInit
func checkInit(t *testing.T, stub *shim.MockStub) {
	// Arguments for mocking Init
	Args := [][]byte{[]byte("init"), []byte("Org1"), []byte("500"), []byte("Org2"), []byte("600")}

	response := stub.MockInit("001", Args)
	testLog.Info("Response is: ", response)
	fmt.Println("End of InitCC")
}

//checkInvoke
func checkInvoke(t *testing.T, stub *shim.MockStub, function string) {
	var args [][]byte
	if function == "invoke" {
		args = [][]byte{[]byte("invoke"), []byte("invoke"), []byte("Org1"), []byte("Org2"), []byte("10")}
	} else if function == "createUser" {
		args = [][]byte{[]byte("invoke"), []byte("createUser"), []byte("Alpha"), []byte("alpha@gmail.com"), []byte("123456789"), []byte("TestCreateUser"), []byte("Org1"), []byte("false"), []byte("50")}
	}
	response := stub.MockInvoke("002", args)
	//stub.
	if response.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(response.Message))
		t.FailNow()
	}
	fmt.Println("Invoke successfully", string(response.Message))
	testLog.Info("Response in checkInvoke is: ", response)
}

//checkQuery
func checkQuery(t *testing.T, stub *shim.MockStub, args [][]byte) {
	response := stub.MockInvoke("1", args)
	if response.Status != shim.OK {
		testLog.Info("Query", args[1], "failed", string(response.Message))
		t.FailNow()
	}
	if response.Payload == nil {
		testLog.Info("Query", args[1], "failed to get value")
		t.FailNow()
	}
	payload := string(response.Payload)

	testLog.Info("Query value", args[1], "is", payload, "as expected")
}

//checkState
func checkState(t *testing.T, stub *shim.MockStub, name string) {
	bytes := stub.State[name]
	if bytes == nil {
		testLog.Info("State", name, "failed to get value")
		t.FailNow()
	}
	testLog.Info("State value", name, "is", string(bytes), "as expected")
}

// ##################### Test Cases ###################### //
func TestInit(t *testing.T) {
	cc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex01", cc)
	fmt.Println("================Test initLedger==========================")
	checkInit(t, stub)
	fmt.Println("================Test invokeLedger==========================")
	//args := [][]byte{[]byte("invoke"), []byte("Org1"), []byte("Org2"), []byte("10")}
	checkInvoke(t, stub, "invoke")
	fmt.Println("================End ======================================")
	fmt.Println("")

}

func TestCreateUser(t *testing.T) {
	fmt.Println("\n ###### Inside createUserTest ")
	cc := new(SimpleChaincode)
	stub := shim.NewMockStub("ex02", cc)
	fmt.Println("\n ================Test Create User ==========================")
	checkInvoke(t, stub, "createUser")

	fmt.Println("================End ======================================")
	fmt.Println("")
}
