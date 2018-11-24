/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Trade Finance Use Case - WORK IN  PROGRESS
 */

package main


import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the trade agreement
type TradeAgreement struct {
	Amount			int		`json:"amount"`
	Goods	string		`json:"Goods"`
	Status			string		`json:"status"`
}
// Define the letter of credit
type LetterOfCredit struct {
	Id			string		`json:"id"`
	ExpirationDate		string		`json:"expirationDate"`
	Beneficiary		string		`json:"beneficiary"`
	Amount			int		`json:"amount"`
	Status			string		`json:"status"`
}


func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "acceptTrade" {
		return s.acceptTrade(APIstub, args)
	} else if function == "requestTrade" {
		return s.requestTrade(APIstub)
	} else if function == "requestLC" {
		return s.createCar(APIstub, args)
	} else if function == "issueLC" {
		return s.issueLC(APIstub)
	} else if function == "requestPayment" {
		return s.requestPayment(APIstub, args)
	}else if function == "makePayment" {
		return s.makePayment(APIstub, args)
	}else if function == "getTradeStatus" {
		return s.getTradeStatus(APIstub, args)
	}else if function == "getLCStatus" {
		return s.getLCStatus(APIstub, args)
	}else if function == "getLCStatus" {
		return s.getLCStatus(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) requestTrade(APIstub shim.ChaincodeStubInterface,args []string) sc.Response {

	TradeAgreementId := args[0];
	Amount, err := strconv.Atoi(args[1]);
	if err != nil {
		return shim.Error("No Amount")


	TradeAgreement := TradeAgreement{Amount: Amount, Goods: args[2], Status: "Requested"}
	tradeAgreementBytes, err := json.Marshal(TradeAgreement)

  APIstub.PutState(TradeAgreementId,tradeAgreementBytes)
	fmt.Println("Trade Requested", TradeAgreement)

	return shim.Success(nil)
}

func (s *SmartContract) acceptTrade(APIstub shim.ChaincodeStubInterface) sc.Response {


	return shim.Success(nil)
}

func (s *SmartContract) requestLC(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {


	return shim.Success(nil)
}

func (s *SmartContract) issueLC(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {


	return shim.Success(nil)
}

func (s *SmartContract) requestLC(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {


	return shim.Success(nil)
}



// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
