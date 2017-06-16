package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var i int

//Entity - Structure for an entity like user, merchant, bank
type Entity struct {
	Type string  `json:"type"`
	Name string  `json:"name"`
	Euro float64 `json:"euroBalance"`
	USD  int     `json:"usdBalance"`
}

//TxnTopup - User transactions for adding USD or Euro
type TxnTopup struct {
	Initiator string `json:"initiator"`
	Remarks   string `json:"remarks"`
	ID        string `json:"id"`
	Time      string `json:"time"`
	Value     string `json:"value"`
	Asset     string `json:"asset"`
}

//TxnTransfer - User transactions for transfer of USD or Euro
type TxnTransfer struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Remarks  string `json:"remarks"`
	ID       string `json:"id"`
	Time     string `json:"time"`
	Value    string `json:"value"`
	Asset    string `json:"asset"`
}

//ExchangeCurrency - User transaction details for currency exchange
type ExchangeCurrency struct {
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Remarks   string `json:"remarks"`
	ID        string `json:"id"`
	Time      string `json:"time"`
	ValueFrom string `json:"valueFrom"`
	AssetFrom string `json:"assetFrom"`
	ValueTo   string `json:"valueTo"`
	AssetTo   string `json:"assetTo"`
}

// CrossBorderChainCode example simple Chaincode implementation
type CrossBorderChainCode struct {
}

func main() {
	err := shim.Start(new(CrossBorderChainCode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *CrossBorderChainCode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	key1 := args[0] //customer
	key2 := args[1] //regulator
	key3 := args[2] //receivingBank
	key4 := args[3] //exchangeCounter

	cust := Entity{
		Type: "customer",
		Name: key1,
		USD:  3000,
		Euro: 3000,
	}

	fmt.Println(cust)
	bytes, err := json.Marshal(cust)
	if err != nil {
		fmt.Println("Error marsalling")
		return nil, errors.New("Error marshalling")
	}
	fmt.Println(bytes)
	err = stub.PutState(key1, bytes)
	if err != nil {
		fmt.Println("Error writing state")
		return nil, err
	}

	regulator := Entity{
		Type: "regulator",
		Name: key2,
		USD:  6000,
		Euro: 6000,
	}
	fmt.Println(regulator)
	bytes, err = json.Marshal(regulator)
	if err != nil {
		fmt.Println("Error marsalling")
		return nil, errors.New("Error marshalling")
	}
	fmt.Println(bytes)
	err = stub.PutState(key2, bytes)
	if err != nil {
		fmt.Println("Error writing state")
		return nil, err
	}

	bank := Entity{
		Type: "receivingBank",
		Name: key3,
		USD:  10000,
		Euro: 10000,
	}
	fmt.Println(bank)
	bytes, err = json.Marshal(bank)
	if err != nil {
		fmt.Println("Error marsalling")
		return nil, errors.New("Error marshalling")
	}
	fmt.Println(bytes)
	err = stub.PutState(key3, bytes)
	if err != nil {
		fmt.Println("Error writing state")
		return nil, err
	}

	exchangeCounter := Entity{
		Type: "exchangeCounter",
		Name: key4,
		USD:  10000,
		Euro: 10000,
	}
	fmt.Println(exchangeCounter)
	bytes, err = json.Marshal(exchangeCounter)

	if err != nil {
		fmt.Println("Error marsalling")
		return nil, errors.New("Error marshalling")
	}
	fmt.Println(bytes)
	err = stub.PutState(key4, bytes)
	if err != nil {
		fmt.Println("Error writing state")
		return nil, err
	}

	// Initialize the collection of  keys for assets and various transactions
	fmt.Println("Initializing keys collection")
	var blank []string

	assets := []string{"USD", "Euro"}
	assetsBytes, _ := json.Marshal(&assets)

	err = stub.PutState("Assets", assetsBytes)
	if err != nil {
		fmt.Println("Failed to initialize Assets key collection")
	}

	blankBytes, _ := json.Marshal(&blank)

	err = stub.PutState("TxnTopup", blankBytes)
	if err != nil {
		fmt.Println("Failed to initialize TxnTopUp key collection")
	}
	err = stub.PutState("ExchangeCurrency", blankBytes)
	if err != nil {
		fmt.Println("Failed to initialize ExchangeCurrency key collection")
	}

	err = stub.PutState("TxnTransfer", blankBytes)
	if err != nil {
		fmt.Println("Failed to initialize TxnTransfer key collection")
	}

	fmt.Println("Initialization complete")

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *CrossBorderChainCode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions/transactions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	}

	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *CrossBorderChainCode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {
		return t.read(stub, args)
	} else if function == "getAllCurrencies" {
		return nil, nil //t.getAllCurrencies(stub)
	} else if function == "getAllTxnTopup" {
		return t.getAllTxnTopup(stub)
	} else if function == "getAllExchangeRecords" {
		return nil, nil //t.getAllExchangeRecords(stub)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *CrossBorderChainCode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("running write()")

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. expecting 3")
	}

	//writing a new customer to blockchain
	typeOf := args[0]
	name := args[1]
	euro, err := strconv.ParseFloat(args[2], 64)
	usd, err := strconv.Atoi(args[3])

	entity := Entity{
		Type: typeOf,
		Name: name,
		Euro: euro,
		USD:  usd,
	}

	fmt.Println(entity)
	bytes, err := json.Marshal(entity)
	if err != nil {
		fmt.Println("Error marsalling")
		return nil, errors.New("Error marshalling")
	}
	fmt.Println(bytes)
	err = stub.PutState(name, bytes)
	if err != nil {
		fmt.Println("Error writing state")
		return nil, err
	}

	return nil, nil
}

// read - query function to read key/value pair
func (t *CrossBorderChainCode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("read() is running")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. expecting 1")
	}

	key := args[0] // name of Entity

	bytes, err := stub.GetState(key)
	if err != nil {
		fmt.Println("Error retrieving " + key)
		return nil, errors.New("Error retrieving " + key)
	}
	customer := Entity{}
	err = json.Unmarshal(bytes, &customer)
	if err != nil {
		fmt.Println("Error Unmarshaling customerBytes")
		return nil, errors.New("Error Unmarshaling customerBytes")
	}
	bytes, err = json.Marshal(customer)
	if err != nil {
		fmt.Println("Error marshaling customer")
		return nil, errors.New("Error marshaling customer")
	}

	fmt.Println(bytes)
	return bytes, nil
}

func (t *CrossBorderChainCode) putTxnTopup(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("putTxnTopup is running ")

	if len(args) != 5 {
		return nil, errors.New("Incorrect Number of arguments.Expecting 5 for putTxnTopup")
	}
	txn := TxnTopup{
		Initiator: args[1],
		Remarks:   args[0] + " addedd",
		ID:        args[3],
		Time:      args[4],
		Value:     args[2],
		Asset:     args[0],
	}

	bytes, err := json.Marshal(txn)
	if err != nil {
		fmt.Println("Error marshaling TxnTopup")
		return nil, errors.New("Error marshaling TxnTopup")
	}

	err = stub.PutState(txn.ID, bytes)
	if err != nil {
		return nil, err
	}

	return t.appendKey(stub, "TxnTopup", txn.ID)
}

func (t *CrossBorderChainCode) getAllTxnTopup(stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("getAllTxnTopup is running ")

	var txns []TxnTopup

	// Get list of all the keys - TxnTopup
	keysBytes, err := stub.GetState("TxnTopup")
	if err != nil {
		fmt.Println("Error retrieving TxnTopup keys")
		return nil, errors.New("Error retrieving TxnTopup keys")
	}
	var keys []string
	err = json.Unmarshal(keysBytes, &keys)
	if err != nil {
		fmt.Println("Error unmarshalling TxnTopup key")
		return nil, errors.New("Error unmarshalling TxnTopup keys")
	}

	// Get each product txn "TxnTopup" keys
	for _, value := range keys {
		bytes, err := stub.GetState(value)

		var txn TxnTopup
		err = json.Unmarshal(bytes, &txn)
		if err != nil {
			fmt.Println("Error retrieving txn " + value)
			return nil, errors.New("Error retrieving txn " + value)
		}

		fmt.Println("Appending txn" + value)
		txns = append(txns, txn)
	}

	bytes, err := json.Marshal(txns)
	if err != nil {
		fmt.Println("Error marshaling txns topup")
		return nil, errors.New("Error marshaling txns topup")
	}
	return bytes, nil
}
func (t *CrossBorderChainCode) appendKey(stub shim.ChaincodeStubInterface, primeKey string, key string) ([]byte, error) {
	fmt.Println("appendKey is running " + primeKey + " " + key)

	bytes, err := stub.GetState(primeKey)
	if err != nil {
		return nil, err
	}
	var keys []string
	err = json.Unmarshal(bytes, &keys)
	if err != nil {
		return nil, err
	}
	keys = append(keys, key)
	bytes, err = json.Marshal(keys)
	if err != nil {
		fmt.Println("Error marshaling " + primeKey)
		return nil, errors.New("Error marshaling keys" + primeKey)
	}
	err = stub.PutState(primeKey, bytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
