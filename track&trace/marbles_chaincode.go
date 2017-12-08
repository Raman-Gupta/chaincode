/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

// ====CHAINCODE EXECUTION SAMPLES (CLI) ==================

// ==== Invoke marbles ====
// peer chaincode invoke -C myc1 -n marbles -c '{"Args":["initSurgicalKit","marble1","blue","35","tom"]}'
// peer chaincode invoke -C myc1 -n marbles -c '{"Args":["initSurgicalKit","marble2","red","50","tom"]}'
// peer chaincode invoke -C myc1 -n marbles -c '{"Args":["initSurgicalKit","marble3","blue","70","tom"]}'
// peer chaincode invoke -C myc1 -n marbles -c '{"Args":["transferMarble","marble2","jerry"]}'
// peer chaincode invoke -C myc1 -n marbles -c '{"Args":["transferMarblesBasedOnColor","blue","jerry"]}'
// peer chaincode invoke -C myc1 -n marbles -c '{"Args":["delete","marble1"]}'

// ==== Query marbles ====
// peer chaincode query -C myc1 -n marbles -c '{"Args":["readSurgicalKit","marble1"]}'
// peer chaincode query -C myc1 -n marbles -c '{"Args":["getMarblesByRange","marble1","marble3"]}'
// peer chaincode query -C myc1 -n marbles -c '{"Args":["getHistoryForMarble","marble1"]}'

// Rich Query (Only supported if CouchDB is used as state database):
//   peer chaincode query -C myc1 -n marbles -c '{"Args":["queryMarblesByOwner","tom"]}'
//   peer chaincode query -C myc1 -n marbles -c '{"Args":["queryMarbles","{\"selector\":{\"owner\":\"tom\"}}"]}'

//The following examples demonstrate creating indexes on CouchDB
//Example hostname:port configurations
//
//Docker or vagrant environments:
// http://couchdb:5984/
//
//Inside couchdb docker container
// http://127.0.0.1:5984/

// Index for chaincodeid, docType, owner.
// Note that docType and owner fields must be prefixed with the "data" wrapper
// chaincodeid must be added for all queries
//
// Definition for use with Fauxton interface
// {"index":{"fields":["chaincodeid","data.docType","data.owner"]},"ddoc":"indexOwnerDoc", "name":"indexOwner","type":"json"}
//
// example curl definition for use with command line
// curl -i -X POST -H "Content-Type: application/json" -d "{\"index\":{\"fields\":[\"chaincodeid\",\"data.docType\",\"data.owner\"]},\"name\":\"indexOwner\",\"ddoc\":\"indexOwnerDoc\",\"type\":\"json\"}" http://hostname:port/myc1/_index
//

// Index for chaincodeid, docType, owner, size (descending order).
// Note that docType, owner and size fields must be prefixed with the "data" wrapper
// chaincodeid must be added for all queries
//
// Definition for use with Fauxton interface
// {"index":{"fields":[{"data.size":"desc"},{"chaincodeid":"desc"},{"data.docType":"desc"},{"data.owner":"desc"}]},"ddoc":"indexSizeSortDoc", "name":"indexSizeSortDesc","type":"json"}
//
// example curl definition for use with command line
// curl -i -X POST -H "Content-Type: application/json" -d "{\"index\":{\"fields\":[{\"data.size\":\"desc\"},{\"chaincodeid\":\"desc\"},{\"data.docType\":\"desc\"},{\"data.owner\":\"desc\"}]},\"ddoc\":\"indexSizeSortDoc\", \"name\":\"indexSizeSortDesc\",\"type\":\"json\"}" http://hostname:port/myc1/_index

// Rich Query with index design doc and index name specified (Only supported if CouchDB is used as state database):
//   peer chaincode query -C myc1 -n marbles -c '{"Args":["queryMarbles","{\"selector\":{\"docType\":\"marble\",\"owner\":\"tom\"}, \"use_index\":[\"_design/indexOwnerDoc\", \"indexOwner\"]}"]}'

// Rich Query with index design doc specified only (Only supported if CouchDB is used as state database):
//   peer chaincode query -C myc1 -n marbles -c '{"Args":["queryMarbles","{\"selector\":{\"docType\":{\"$eq\":\"marble\"},\"owner\":{\"$eq\":\"tom\"},\"size\":{\"$gt\":0}},\"fields\":[\"docType\",\"owner\",\"size\"],\"sort\":[{\"size\":\"desc\"}],\"use_index\":\"_design/indexSizeSortDoc\"}"]}'

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type surgicalkit struct {
	ObjectType		string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	AssetID    		string `json:"assetID"`    //the fieldtags are needed to keep case from bouncing around
	Payload			string `json:"Payload"`
	Compliant       bool   `json:"compliant"`
	// St         		string `json:"st"`
	// Location   		string `json:"location"`
	// Ownerrole       string `json:"ownerrole"`
	// Ownername 		string `json:"ownerrole"`
	// AssetType		string `json:"assetType`
	// AssetModel		string `json:"assetModel"`
	// FirmwareVersion	string `json:"firmwareVersion"`
	// BuildNumber		string `json:"buildNumber"`
	// BuildBy			string `json:"buildBy"`
	// BatteryLevel	string `json:"batteryLevel"`
	// Ownerid			string `json:"ownerid"`
	// Alt				string `json:"Alt"`
	// Lon				string `json:"Lon"`
	// Lat				string `json:"Lat"`
	// CSts			string `json:"CSts"`
	// DSts			string `json:"DSts"`
	// Tmpr			string `json:"Tmpr"`
	// Hmdt			string `json:"Hmdt"`
	// Cycle 			int    `json:"cycle"`
	// Elst			string `json:"elst"`
	// 
	
}

type customEvent struct {
    Type        string `json:"type"`
    Description string `json:"description"`
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	if function == "initSurgicalKit" { //create a new marble
		return t.initSurgicalKit(stub, args)
	}else if function == "delete" { //delete a marble
		return t.delete(stub, args)
	} else if function == "readSurgicalKit" { //read a marble
		return t.readSurgicalKit(stub, args)
	}else if function == "getHistoryForSurgicalKit" { //get history of values for a marble
		return t.getHistoryForSurgicalKit(stub, args)
	} else if function == "getSurgicalKitByRange" { //get marbles based on range query
		return t.getSurgicalKitByRange(stub, args)
	}else if function == "transferSurgicalKit" { //change owner of a specific marble
		return t.transferSurgicalKit(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
// initSurgicalKit - create a new marble, store into chaincode state
// ============================================================
func (t *SimpleChaincode) initSurgicalKit(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	
	//   0       1       2     3
	// "asdf", "blue", "35", "bob"
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	
	fmt.Println(len(args))
	// ==== Input sanitation ====
	fmt.Println("- start init marble")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	// if len(args[2]) <= 0 {
	// 	return shim.Error("3rd argument must be a non-empty string")
	// }
	// if len(args[3]) <= 0 {
	// 	return shim.Error("4th argument must be a non-empty string")
	// }
	// if len(args[4]) <= 0 {
	// 	return shim.Error("5th argument must be a non-empty string")
	// }
	// if len(args[5]) <= 0 {
	// 	return shim.Error("6th argument must be a non-empty string")
	// }
	// if len(args[6]) <= 0 {
	// 	return shim.Error("7th argument must be a non-empty string")
	// }
	// if len(args[7]) <= 0 {
	// 	return shim.Error("8th argument must be a non-empty string")
	// }
	// if len(args[8]) <= 0 {
	// 	return shim.Error("9th argument must be a non-empty string")
	// }
	// if len(args[9]) <= 0 {
	// 	return shim.Error("10th argument must be a non-empty string")
	// }
	// if len(args[10]) <= 0 {
	// 	return shim.Error("11th argument must be a non-empty string")
	// }
	// if len(args[11]) <= 0 {
	// 	return shim.Error("12th argument must be a non-empty string")
	// }
	// if len(args[12]) <= 0 {
	// 	return shim.Error("13th argument must be a non-empty string")
	// }
	// if len(args[13]) <= 0 {
	// 	return shim.Error("14th argument must be a non-empty string")
	// }
	// if len(args[14]) <= 0 {
	// 	return shim.Error("15th argument must be a non-empty string")
	// }
	// if len(args[15]) <= 0 {
	// 	return shim.Error("16th argument must be a non-empty string")
	// }
	// if len(args[16]) <= 0 {
	// 	return shim.Error("17th argument must be a non-empty string")
	// }
	// if len(args[17]) <= 0 {
	// 	return shim.Error("18th argument must be a non-empty string")
	// }
	// if len(args[18]) <= 0 {
	// 	return shim.Error("19th argument must be a non-empty string")
	// }
	// if len(args[19]) <= 0 {
	// 	return shim.Error("20th argument must be a array string")
	// }
	// if len(args[20]) <= 0 {
	// 	return shim.Error("21th argument must be a non-empty string")
	// }
	
	AssetID   		:= args[0]
	Payload	      	:= args[1]
	// Location  		:= strings.ToLower(args[2])
	// Ownerrole		:= strings.ToLower(args[3])
	// Ownername 		:= strings.ToLower(args[4])
	// AssetType 		:= strings.ToLower(args[5])
	// AssetModel		:= strings.ToLower(args[6])
	// FirmwareVersion := strings.ToLower(args[7])
	// BuildNumber 	:= strings.ToLower(args[8])
	// BuildBy 		:= strings.ToLower(args[9])
	// BatteryLevel 	:= strings.ToLower(args[10])
	// Ownerid 		:= strings.ToLower(args[11])
	// Alt 			:= strings.ToLower(args[12])
	// Lon 			:= strings.ToLower(args[13])
	// Lat 			:= strings.ToLower(args[14])
	// CSts 			:= strings.ToLower(args[15])
	// DSts 			:= strings.ToLower(args[16])
	// Tmpr 			:= strings.ToLower(args[17])
	// Hmdt 			:= strings.ToLower(args[18])
	// Elst			:= strings.ToLower(args[19])
	// Cycle, err 		:= strconv.Atoi(args[20])
	// if err != nil {
	// 	return shim.Error("Cycle field value must be numeric")
	// }
	fmt.Println("AssetID--"+AssetID)
	fmt.Println("Payload--"+Payload)
	

	// ==== Check if kit already exists ====
	kitAsBytes, err := stub.GetState(AssetID)
	if err != nil {
		return shim.Error("Failed to get kit: " + err.Error())
	} else if kitAsBytes != nil {
		fmt.Println("This Kit already exists: " + AssetID)
		return shim.Error("This Kit already exists: " + AssetID)
	}
	// ==== Create kit object and kit to JSON ====
	objectType := "surgicalkit"
	surgicalkit := &surgicalkit{ObjectType:objectType, AssetID: AssetID, Payload:Payload, Compliant:true}
	surgicalKitJSONasBytes, err := json.Marshal(surgicalkit)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println(string(surgicalKitJSONasBytes))
	// === Save marble to state ===
	err = stub.PutState(AssetID, surgicalKitJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	//  ==== Index the marble to enable color-based range queries, e.g. return all blue marbles ====
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on indexName~color~name.
	//  This will enable very efficient state range queries based on composite keys matching indexName~color~*
	indexName := "assetId~name"
	kitStatusIndexKey, err := stub.CreateCompositeKey(indexName, []string{surgicalkit.AssetID, surgicalkit.AssetID})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(kitStatusIndexKey, value)

	// ==== Marble saved and indexed. Return success ====
	fmt.Println("- end init surgicalkit")
	return shim.Success(nil)
}

// ===============================================
// readSurgicalKit - read a marble from chaincode state
// ===============================================
func (t *SimpleChaincode) readSurgicalKit(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var AssetID, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the surgical kit to query")
	}

	AssetID = args[0]
	valAsbytes, err := stub.GetState(AssetID) //get the surgical kit from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + AssetID + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Surgical Kit does not exist: " + AssetID + "\"}"
		return shim.Error(jsonResp)
	}
	return shim.Success(valAsbytes)
}

// ==================================================
// delete - remove a marble key/value pair from state
// ==================================================
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var jsonResp string
	var surgicalKitJSON surgicalkit
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	surgicalKitName := args[0]

	// to maintain the color~name index, we need to read the marble first and get its color
	valAsbytes, err := stub.GetState(surgicalKitName) //get the marble from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + surgicalKitName + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"surgicalKitName does not exist: " + surgicalKitName + "\"}"
		return shim.Error(jsonResp)
	}

	err = json.Unmarshal([]byte(valAsbytes), &surgicalKitJSON)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to decode JSON of: " + surgicalKitName + "\"}"
		return shim.Error(jsonResp)
	}

	err = stub.DelState(surgicalKitName) //remove the marble from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}

	// maintain the index
	indexName := "assetId~name"
	stNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{surgicalKitJSON.AssetID, surgicalKitJSON.AssetID})
	if err != nil {
		return shim.Error(err.Error())
	}

	//  Delete index entry to state.
	err = stub.DelState(stNameIndexKey)
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}
	return shim.Success(nil)
}

// // ===========================================================
// // transfer a marble by setting a new owner name on the marble
// // ===========================================================
func (t *SimpleChaincode) transferSurgicalKit(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//   0       1
	// "name", "bob"
	if len(args) <1 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	var kitstatus string
	AssetID   		:= args[0]
	Payload	      	:= args[1]
	fmt.Println("- start updation ", AssetID, Payload)
    
	surgicalkitAsBytes, err := stub.GetState(AssetID)
	if err != nil {
		return shim.Error("Failed to get surgicalkit:" + err.Error())
	} else if surgicalkitAsBytes == nil {
		return shim.Error("surgicalkit does not exist")
	}

	surgicalkitToTransfer := surgicalkit{}
	err = json.Unmarshal(surgicalkitAsBytes, &surgicalkitToTransfer) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	surgicalkitToTransfer.AssetID 			= AssetID //change the owner
	surgicalkitToTransfer.Payload 			= Payload
	var j =[]byte(Payload)
	// a map container to decode the JSON structure into
    c := make(map[string]interface{})

    // unmarschal JSON
     json.Unmarshal(j, &c)
	for s, v := range c {
		switch vv :=v.(type) {
			case string :
						if(s=="st"){
							fmt.Println("status field found in payload--"+vv)
							kitstatus = vv
						}
		}
	}
	fmt.Println("kitstatus------"+kitstatus)

	if kitstatus == "tampered" {
		surgicalkitToTransfer.Compliant = false
	}else{
		surgicalkitToTransfer.Compliant = true
	}
	surgicalkitJSONasBytes, _ := json.Marshal(surgicalkitToTransfer)
	err = stub.PutState(AssetID, surgicalkitJSONasBytes) //rewrite the marble
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("- end transfersurgicalkit(success)")
	return shim.Success(nil)
}

// // ===========================================================================================
// // getMarblesByRange performs a range query based on the start and end keys provided.

// // Read-only function results are not typically submitted to ordering. If the read-only
// // results are submitted to ordering, or if the query is used in an update transaction
// // and submitted to ordering, then the committing peers will re-execute to guarantee that
// // result sets are stable between endorsement time and commit time. The transaction is
// // invalidated by the committing peers if the result set has changed between endorsement
// // time and commit time.
// // Therefore, range queries are a safe option for performing update transactions based on query results.
// // ===========================================================================================
func (t *SimpleChaincode) getSurgicalKitByRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	startKey := args[0]
	endKey := args[1]

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getSurgicalKitByRange queryResult:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// ==== Example: GetStateByPartialCompositeKey/RangeQuery =========================================
// transferMarblesBasedOnColor will transfer marbles of a given color to a certain new owner.
// Uses a GetStateByPartialCompositeKey (range query) against color~name 'index'.
// Committing peers will re-execute range queries to guarantee that result sets are stable
// between endorsement time and commit time. The transaction is invalidated by the
// committing peers if the result set has changed between endorsement time and commit time.
// Therefore, range queries are a safe option for performing update transactions based on query results.
// ===========================================================================================
// func (t *SimpleChaincode) transferMarblesBasedOnColor(stub shim.ChaincodeStubInterface, args []string) pb.Response {

// 	//   0       1
// 	// "color", "bob"
// 	if len(args) < 2 {
// 		return shim.Error("Incorrect number of arguments. Expecting 2")
// 	}

// 	color := args[0]
// 	newOwner := strings.ToLower(args[1])
// 	fmt.Println("- start transferMarblesBasedOnColor ", color, newOwner)

// 	// Query the color~name index by color
// 	// This will execute a key range query on all keys starting with 'color'
// 	coloredMarbleResultsIterator, err := stub.GetStateByPartialCompositeKey("color~name", []string{color})
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}
// 	defer coloredMarbleResultsIterator.Close()

// 	// Iterate through result set and for each marble found, transfer to newOwner
// 	var i int
// 	for i = 0; coloredMarbleResultsIterator.HasNext(); i++ {
// 		// Note that we don't get the value (2nd return variable), we'll just get the marble name from the composite key
// 		responseRange, err := coloredMarbleResultsIterator.Next()
// 		if err != nil {
// 			return shim.Error(err.Error())
// 		}

// 		// get the color and name from color~name composite key
// 		objectType, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
// 		if err != nil {
// 			return shim.Error(err.Error())
// 		}
// 		returnedColor := compositeKeyParts[0]
// 		returnedMarbleName := compositeKeyParts[1]
// 		fmt.Printf("- found a marble from index:%s color:%s name:%s\n", objectType, returnedColor, returnedMarbleName)

// 		// Now call the transfer function for the found marble.
// 		// Re-use the same function that is used to transfer individual marbles
// 		response := t.transferMarble(stub, []string{returnedMarbleName, newOwner})
// 		// if the transfer failed break out of loop and return error
// 		if response.Status != shim.OK {
// 			return shim.Error("Transfer failed: " + response.Message)
// 		}
// 	}

// 	responsePayload := fmt.Sprintf("Transferred %d %s marbles to %s", i, color, newOwner)
// 	fmt.Println("- end transferMarblesBasedOnColor: " + responsePayload)
// 	return shim.Success([]byte(responsePayload))
// }

// // =======Rich queries =========================================================================
// // Two examples of rich queries are provided below (parameterized query and ad hoc query).
// // Rich queries pass a query string to the state database.
// // Rich queries are only supported by state database implementations
// //  that support rich query (e.g. CouchDB).
// // The query string is in the syntax of the underlying state database.
// // With rich queries there is no guarantee that the result set hasn't changed between
// //  endorsement time and commit time, aka 'phantom reads'.
// // Therefore, rich queries should not be used in update transactions, unless the
// // application handles the possibility of result set changes between endorsement and commit time.
// // Rich queries can be used for point-in-time queries against a peer.
// // ============================================================================================

// // ===== Example: Parameterized rich query =================================================
// // queryMarblesByOwner queries for marbles based on a passed in owner.
// // This is an example of a parameterized query where the query logic is baked into the chaincode,
// // and accepting a single query parameter (owner).
// // Only available on state databases that support rich query (e.g. CouchDB)
// // =========================================================================================
// func (t *SimpleChaincode) queryMarblesByOwner(stub shim.ChaincodeStubInterface, args []string) pb.Response {

// 	//   0
// 	// "bob"
// 	if len(args) < 1 {
// 		return shim.Error("Incorrect number of arguments. Expecting 1")
// 	}

// 	owner := strings.ToLower(args[0])

// 	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"marble\",\"owner\":\"%s\"}}", owner)

// 	queryResults, err := getQueryResultForQueryString(stub, queryString)
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}
// 	return shim.Success(queryResults)
// }

// // ===== Example: Ad hoc rich query ========================================================
// // queryMarbles uses a query string to perform a query for marbles.
// // Query string matching state database syntax is passed in and executed as is.
// // Supports ad hoc queries that can be defined at runtime by the client.
// // If this is not desired, follow the queryMarblesForOwner example for parameterized queries.
// // Only available on state databases that support rich query (e.g. CouchDB)
// // =========================================================================================
// func (t *SimpleChaincode) queryMarbles(stub shim.ChaincodeStubInterface, args []string) pb.Response {

// 	//   0
// 	// "queryString"
// 	if len(args) < 1 {
// 		return shim.Error("Incorrect number of arguments. Expecting 1")
// 	}

// 	queryString := args[0]

// 	queryResults, err := getQueryResultForQueryString(stub, queryString)
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}
// 	return shim.Success(queryResults)
// }

// // =========================================================================================
// // getQueryResultForQueryString executes the passed in query string.
// // Result set is built and returned as a byte array containing the JSON results.
// // =========================================================================================
// func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

// 	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

// 	resultsIterator, err := stub.GetQueryResult(queryString)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resultsIterator.Close()

// 	// buffer is a JSON array containing QueryRecords
// 	var buffer bytes.Buffer
// 	buffer.WriteString("[")

// 	bArrayMemberAlreadyWritten := false
// 	for resultsIterator.HasNext() {
// 		queryResponse, err := resultsIterator.Next()
// 		if err != nil {
// 			return nil, err
// 		}
// 		// Add a comma before array members, suppress it for the first array member
// 		if bArrayMemberAlreadyWritten == true {
// 			buffer.WriteString(",")
// 		}
// 		buffer.WriteString("{\"Key\":")
// 		buffer.WriteString("\"")
// 		buffer.WriteString(queryResponse.Key)
// 		buffer.WriteString("\"")

// 		buffer.WriteString(", \"Record\":")
// 		// Record is a JSON object, so we write as-is
// 		buffer.WriteString(string(queryResponse.Value))
// 		buffer.WriteString("}")
// 		bArrayMemberAlreadyWritten = true
// 	}
// 	buffer.WriteString("]")

// 	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

// 	return buffer.Bytes(), nil
// }

func (t *SimpleChaincode) getHistoryForSurgicalKit(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	surgicalKitName := args[0]

	fmt.Printf("- start getHistoryForSurgicalKit: %s\n", surgicalKitName)

	resultsIterator, err := stub.GetHistoryForKey(surgicalKitName)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the surgicalkit
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForSurgicalKit returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
