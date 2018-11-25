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
	Amount		int				`json:"amount"`
	Goods			string		`json:"Goods"`
	Status		string		`json:"status"` //REQUESTED,ACCEPTED,SHIPPED,GOODS RECEIVED,PAYMENT REQUESTED,PAYMENT DONE
}
// Define the letter of credit
type LetterOfCredit struct {
	Id			string		`json:"id"`
	ExpirationDate		string		`json:"expirationDate"`
	Beneficiary		string		`json:"beneficiary"`
	Amount			int		`json:"amount"`
	Status			string		`json:"status"` //REQUESTED,ISSUED,ACCEPTED
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
	} else if function == "acceptLC" {
		return s.issueLC(APIstub)
	} else if function == "setShipmentStatus" {
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

	tradeAgreementId := args[0]
	Amount, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("No Amount")


	TradeAgreement := TradeAgreement{Amount: Amount, Goods: args[2], Status: "REQUESTED"}
	tradeAgreementBytes, err := json.Marshal(TradeAgreement)

  APIstub.PutState(TradeAgreementId,tradeAgreementBytes)
	fmt.Printf("Trade %s REQUESTED\n", args[0])

	return shim.Success(nil)
}

func (s *SmartContract) acceptTrade(APIstub shim.ChaincodeStubInterface,args []string) sc.Response {

	var tradeAgreement *TradeAgreement
	var tradeAgreementBytes []byte

	tradeAgreementId := args[0]
	tradeAgreementBytes, err = stub.GetState(tradeAgreementId)
	if err != nil {
		return shim.Error(err.Error())
	}
	if len(tradeAgreementBytes) == 0 {
	err = errors.New(fmt.Sprintf("No record found for trade ID %s", args[0]))
	return shim.Error(err.Error())
  }

	err = json.Unmarshal(tradeAgreementBytes, &tradeAgreement)
  if err != nil {
	 return shim.Error(err.Error())
  }

	if tradeAgreement.Status == "ACCEPTED" {
		fmt.Printf("Trade %s already accepted", args[0])
	} else {
		tradeAgreement.Status = "ACCEPTED"
		tradeAgreementBytes, err = json.Marshal(tradeAgreement)
		if err != nil {
			return shim.Error("Error marshaling trade agreement structure")
		}

		err = stub.PutState(tradeAgreementId, tradeAgreementBytes)
		if err != nil {
			return shim.Error(err.Error())
	}
}
fmt.Printf("Trade %s ACCEPTED\n", args[0])

return shim.Success(nil)
}

func (s *SmartContract) requestLC(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	var tradeAgreementId string
	var LCId string
	var tradeAgreementBytes, letterOfCreditBytes []byte
	var tradeAgreement TradeAgreement
	var letterOfCredit LetterOfCredit

	tradeAgreementId := args[0]
	LCId := args[1]
	expirationDate := args[2]
	beneficiary := args[3]

	tradeAgreementBytes, err = stub.GetState(tradeAgreementId)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = json.Unmarshal(tradeAgreementBytes, &tradeAgreement)
	if err != nil {
		return shim.Error(err.Error())
	}

	if tradeAgreement.Status != "ACCEPTED" {
		return shim.Error("Trade has not been ACCEPTED")
	}

	letterOfCredit = LetterOfCredit{LCId, expirationDate, beneficiary, tradeAgreement.Amount, "REQUESTED"}
	letterOfCreditBytes, err = json.Marshal(letterOfCredit)
	if err != nil {
		return shim.Error("Error marshaling letter of credit structure")
	}
	err = stub.PutState(LCId, letterOfCreditBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Letter of Credit %s REQUESTED\n", LCId)

	return shim.Success(nil)
}

func (s *SmartContract) issueLC(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var LCId string
	var letterOfCreditBytes []byte
	var letterOfCredit LetterOfCredit

	LCId := args[0]

	letterOfCreditBytes, err = stub.GetState(LCId)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = json.Unmarshal(letterOfCreditBytes, &letterOfCredit)
	if err != nil {
		return shim.Error(err.Error())
	}

	if letterOfCredit.Status != "REQUESTED" {
		return shim.Error("Letter of Credit %s Not Requested",LCId)
	}else if letterOfCredit.Status == "ACCEPTED" {
		return shim.Error("Letter of Credit %s  Already Accepted",LCId)

	letterOfCredit.Status = "ISSUED"
	letterOfCreditBytes, err = json.Marshal(letterOfCredit)
	if err != nil {
		return shim.Error("Error marshaling letter of credit structure")
	}
	err = stub.PutState(LCId, letterOfCreditBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Letter of Credit %s ISSUED\n", LCId)

	return shim.Success(nil)
}

func (s *SmartContract) acceptLC(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var LCId string
	var letterOfCreditBytes []byte
	var letterOfCredit LetterOfCredit

	LCId := args[0]

	letterOfCreditBytes, err = stub.GetState(LCId)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = json.Unmarshal(letterOfCreditBytes, &letterOfCredit)
	if err != nil {
		return shim.Error(err.Error())
	}

	if letterOfCredit.Status != "ISSUED" {
		return shim.Error("Letter of Credit %s Not Issued",LCId)
	}else if letterOfCredit.Status == "ACCEPTED" {
		return shim.Error("Letter of Credit %s  Already Accepted",LCId)
	}else{
		letterOfCredit.Status = "ACCEPTED"
		letterOfCreditBytes, err = json.Marshal(letterOfCredit)
		if err != nil {
		return shim.Error("Error marshaling letter of credit structure")
	}

	err = stub.PutState(LCId, letterOfCreditBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("Letter of Credit %s ACCEPTED\n", LCId)

	return shim.Success(nil)
}

func (s *SmartContract) setShipmentStatus(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	var tradeAgreementId string
	var LCId string
	var tradeAgreementBytes, letterOfCreditBytes []byte
	var tradeAgreement TradeAgreement
	var letterOfCredit LetterOfCredit

	tradeAgreementId := args[0]
	LCId := args[1]
	shipmentStatus := args[2]

	tradeAgreementBytes, err = stub.GetState(tradeAgreementId)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = json.Unmarshal(tradeAgreementBytes, &tradeAgreement)
	if err != nil {
		return shim.Error(err.Error())
	}

	letterOfCreditBytes, err = stub.GetState(LCId)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = json.Unmarshal(letterOfCreditBytes, &letterOfCredit)
	if err != nil {
		return shim.Error(err.Error())
	}

	if letterOfCredit.Status != "ACCEPTED"{
		return shim.Error("Cannot Start Shipping as Letter of Credit is not Accepted")
	}else{
		tradeAgreement.Status = shipmentStatus
		tradeAgreementBytes, err = json.Marshal(tradeAgreement)
		if err != nil {
		return shim.Error("Error marshaling trade agreement structure")
	}

	fmt.Printf("Shipment Status for TradeAgreement %s set to %s\n",tradeAgreementId,shipmentStatus)

	return shim.Success(nil)
}

func (s *SmartContract) requestPayment(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	var tradeAgreementId string
	var LCId string
	var tradeAgreementBytes, letterOfCreditBytes []byte
	var tradeAgreement TradeAgreement
	var letterOfCredit LetterOfCredit

	tradeAgreementId := args[0]
	LCId := args[1]

	tradeAgreementBytes, err = stub.GetState(tradeAgreementId)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = json.Unmarshal(tradeAgreementBytes, &tradeAgreement)
	if err != nil {
		return shim.Error(err.Error())
	}

	letterOfCreditBytes, err = stub.GetState(LCId)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = json.Unmarshal(letterOfCreditBytes, &letterOfCredit)
	if err != nil {
		return shim.Error(err.Error())
	}

	if letterOfCredit.Status != "ACCEPTED" && tradeAgreement.Status != "GOODS RECEIVED"{
		return shim.Error("Cannot Start Shipping as Letter of Credit is not Accepted")
	}else{
		tradeAgreement.Status = "PAYMENT REQUESTED"
		tradeAgreementBytes, err = json.Marshal(tradeAgreement)
		if err != nil {
		return shim.Error("Error marshaling trade agreement structure")
	}

	fmt.Printf("Payment Requested for TradeAgreement %s \n",tradeAgreementId)

	return shim.Success(nil)

}

func (s *SmartContract) makePayment(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	var tradeAgreementId string
	var LCId string
	var tradeAgreementBytes, letterOfCreditBytes []byte
	var tradeAgreement TradeAgreement
	var letterOfCredit LetterOfCredit

	tradeAgreementId := args[0]
	LCId := args[1]

	tradeAgreementBytes, err = stub.GetState(tradeAgreementId)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = json.Unmarshal(tradeAgreementBytes, &tradeAgreement)
	if err != nil {
		return shim.Error(err.Error())
	}

	letterOfCreditBytes, err = stub.GetState(LCId)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = json.Unmarshal(letterOfCreditBytes, &letterOfCredit)
	if err != nil {
		return shim.Error(err.Error())
	}

	if letterOfCredit.Status != "ACCEPTED" && tradeAgreement.Status != "PAYMENT REQUESTED"{
		return shim.Error("Cannot Start Shipping as Letter of Credit is not Accepted")
	}else{
		tradeAgreement.Status = "PAYMENT DONE"
		tradeAgreementBytes, err = json.Marshal(tradeAgreement)
		if err != nil {
		return shim.Error("Error marshaling trade agreement structure")
	}

	fmt.Printf("Payment Done for TradeAgreement %s \n",tradeAgreementId)

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
