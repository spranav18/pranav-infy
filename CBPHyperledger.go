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
	USD  float64 `json:"usdBalance"`
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
type TxnExchange struct {
	Initiator    string `json:"initiator"`
	Convertor    string `json:"convertor"`
	SellCurrency string `json:"sellCurrency"`
	SellQuantity string `json:"sellQuantity"`
	BuyCurrency  string `json:"buyCurrency"`
	BuyQuantity  string `json:"buyQuantity"`
	ExchangeRate string `json:"exchangeRate"`
	Remarks      string `json:"remarks"`
	ID           string `json:"id"`
	Time         string `json:"time"`
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
		USD:  3000.00,
		Euro: 3000.00,
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
		USD:  6000.00,
		Euro: 6000.00,
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
		USD:  10000.00,
		Euro: 10000.00,
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
		USD:  10000.00,
		Euro: 10000.00,
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
	err = stub.PutState("TxnExchange", blankBytes)
	if err != nil {
		fmt.Println("Failed to initialize TxnExchange key collection")
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
	} else if function == "loadWallet" {
		return t.loadWallet(stub, args)
	} else if function == "transfer" {
		return t.transfer(stub, args)
	} else if function == "exchangeCurrency" {
		return t.exchangeCurrency(stub, args)
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
	} else if function == "getAllTxnExchange" {
		return t.getAllTxnExchange(stub)    
	} else if function == "getAllTxnTransfer" {
		return t.getAllTxnTransfer(stub)
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
	usd, err := strconv.ParseFloat(args[3],64)

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

func (t *CrossBorderChainCode) exchangeCurrency(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("exchangeCurrency is running ")

	if len(args) != 7 {
		return nil, errors.New("Incorrect Number of arguments.Expecting 7 for exchange currency")
	}
	asset1 := args[0] //domestic currency usd or euro
	//asset2 := args[1]  //foreign currency usd or euro
	exchangeRate,err := strconv.ParseFloat(args[2],64)
	key1 := args[3]  //Entity1 ex: customer
	key2 := args[4]  //Entity2 ex: exchange counter
	qty, err := strconv.ParseFloat(args[5],64)

	bytes, err := stub.GetState(key1)
	if err != nil {
		return nil, errors.New("Failed to get state of " + key1)
	}
	if bytes == nil {
		return nil, errors.New("Entity not found")
	}
	customer := Entity{}
	err = json.Unmarshal(bytes, &customer)
	if err != nil {
		fmt.Println("Error Unmarshaling customerBytes")
		return nil, errors.New("Error Unmarshaling customerBytes")
	}

	bytes, err = stub.GetState(key2)
	if err != nil {
		return nil, errors.New("Failed to get state of " + key2)
	}
	if bytes == nil {
		return nil, errors.New("Entity not found")
	}
	exchangeCounter := Entity{}
	err = json.Unmarshal(bytes, &exchangeCounter)
	if err != nil {
		fmt.Println("Error Unmarshaling exchangeCounterBytes")
		return nil, errors.New("Error Unmarshaling exchangeCounterBytes")
	}
	sellQuantity:=qty*exchangeRate
		// Perform the transfer
		if asset1=="usd" {
			fmt.Println("usd transfer")
			if customer.USD >= qty {
				customer.USD = customer.USD - qty
				exchangeCounter.USD = exchangeCounter.USD + qty
				customer.Euro=customer.Euro+sellQuantity
				exchangeCounter.Euro=exchangeCounter.Euro-sellQuantity
				//args[4] = strconv.Itoa(product.Points * qty)
				fmt.Printf("customer USD = %d, exchangeCounter USD = %d\n", customer.USD, exchangeCounter.USD)
			} else {
				return nil, errors.New("Insufficient points to buy goods")
			}
		} else {
			if customer.Euro >= qty {
				customer.Euro = customer.Euro -qty
				exchangeCounter.Euro = exchangeCounter.Euro + qty
				customer.USD=customer.USD+sellQuantity
				exchangeCounter.USD=exchangeCounter.USD-sellQuantity
			//	args[4] = strconv.FormatFloat(product.Amount*float64(qty), 'E', -1, 64)
				fmt.Printf("customer Euro = %f, exchangeCounter Euro = %f\n", customer.Euro, exchangeCounter.Euro)
			} else {
				return nil, errors.New("Insufficient balance to buy goods")
			}
		}
		// Write the customer/entity1 state back to the ledger
		bytes, err = json.Marshal(customer)
		if err != nil {
			fmt.Println("Error marshaling customer")
			return nil, errors.New("Error marshaling customer")
		}
		err = stub.PutState(key1, bytes)
		if err != nil {
			return nil, err
		}

		// Write the exchangeCounter/entity2 state back to the ledger]
		bytes, err = json.Marshal(exchangeCounter)
		if err != nil {
			fmt.Println("Error marshaling exchangeCounter")
			return nil, errors.New("Error marshaling exchangeCounter")
		}
		err = stub.PutState(key2, bytes)
		if err != nil {
			return nil, err
		}
		// Write the product state back to the ledger
		args=append(args,fmt.Sprintf("%.6f", sellQuantity))
		args = append(args, stub.GetTxID())
		blockTime, err := stub.GetTxTimestamp()
		if err != nil {
			return nil, err
		}
		args = append(args, blockTime.String())
		t.putTxnExchange(stub, args)

	return nil, nil
}

func (t *CrossBorderChainCode) loadWallet(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("loadWallet is running ")

	if len(args) != 3 {
		return nil, errors.New("Incorrect Number of arguments.Expecting 3 for loadWallet")
	}

	asset := args[0] //usd or euro
	key := args[1]   //Entity ex: customer
	//amt, err := strconv.Atoi(args[2]) // points to be issued

	// GET the state of entity from the ledger
	bytes, err := stub.GetState(key)
	if err != nil {
		return nil, errors.New("Failed to get state of " + key)
	}

	entity := Entity{}
	err = json.Unmarshal(bytes, &entity)
	if err != nil {
		fmt.Println("Error Unmarshaling entity Bytes")
		return nil, errors.New("Error Unmarshaling entity Bytes")
	}

	// Perform the addition of assests
	if asset == "usd" {
		amt, err := strconv.ParseFloat(args[2],64)
		if err == nil {
			entity.USD = entity.USD + amt
			fmt.Println("entity USD Balance = ", entity.USD)
		}
	} else {
		amt, err := strconv.ParseFloat(args[2], 64)
		if err == nil {
			entity.Euro = entity.Euro + amt
			fmt.Println("entity Euro Balance = ", entity.Euro)
		}
	}

	// Write the state back to the ledger
	bytes, err = json.Marshal(entity)
	if err != nil {
		fmt.Println("Error marshaling entity")
		return nil, errors.New("Error marshaling entity")
	}
	err = stub.PutState(key, bytes)
	if err != nil {
		return nil, err
	}

	ID := stub.GetTxID()
	blockTime, err := stub.GetTxTimestamp()
	args = append(args, ID)
	args = append(args, blockTime.String())
	t.putTxnTopup(stub, args)

	return nil, nil
}

func (t *CrossBorderChainCode) transfer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("transfer is running ")

	if len(args) != 5 {
		return nil, errors.New("Incorrect Number of arguments.Expecting 5 for transfer")
	}

	key1 := args[0]   // fromEntity ex: customer
	key2 := args[1]  // toEntity ex: merchant
	asset := args[2] // usd or euro

	// GET the state of fromEntity from the ledger
	bytes, err := stub.GetState(key1)
	if err != nil {
		return nil, errors.New("Failed to get state of " + key1)
	}

	fromEntity := Entity{}
	err = json.Unmarshal(bytes, &fromEntity)
	if err != nil {
		fmt.Println("Error Unmarshaling entity Bytes")
		return nil, errors.New("Error Unmarshaling entity Bytes")
	}

	// GET the state of toEntity from the ledger
	bytes, err = stub.GetState(key2)
	if err != nil {
		return nil, errors.New("Failed to get state of " + key2)
	}

	toEntity := Entity{}
	err = json.Unmarshal(bytes, &toEntity)
	if err != nil {
		fmt.Println("Error Unmarshaling entity Bytes")
		return nil, errors.New("Error Unmarshaling entity Bytes")
	}

	// Perform transfer of assests
	if asset == "usd" {
		amt, err := strconv.ParseFloat(args[3],64)
		if err == nil {
			fromEntity.USD = fromEntity.USD - amt
			toEntity.USD = toEntity.USD + amt
			fmt.Println("from entity USD = ", fromEntity.USD)
		}
	} else {
		amt, err := strconv.ParseFloat(args[3], 64)
		if err == nil {
			fromEntity.Euro = fromEntity.Euro - amt
			toEntity.Euro = toEntity.Euro + amt
			fmt.Println("from entity Euro = ", fromEntity.Euro)
		}
	}

	// Write the state back to the ledger
	bytes, err = json.Marshal(fromEntity)
	if err != nil {
		fmt.Println("Error marshaling fromEntity")
		return nil, errors.New("Error marshaling fromEntity")
	}
	err = stub.PutState(key1, bytes)
	if err != nil {
		return nil, err
	}

	bytes, err = json.Marshal(toEntity)
	if err != nil {
		fmt.Println("Error marshaling toEntity")
		return nil, errors.New("Error marshaling toEntity")
	}
	err = stub.PutState(key2, bytes)
	if err != nil {
		return nil, err
	}

	ID := stub.GetTxID()
	blockTime, err := stub.GetTxTimestamp()
	args = append(args, ID)
	args = append(args, blockTime.String())
	t.putTxnTransfer(stub, args)

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
	entity := Entity{}
	err = json.Unmarshal(bytes, &entity)
	if err != nil {
		fmt.Println("Error Unmarshaling entityBytes")
		return nil, errors.New("Error Unmarshaling entityBytes")
	}
	bytes, err = json.Marshal(entity)
	if err != nil {
		fmt.Println("Error marshaling entity")
		return nil, errors.New("Error marshaling entity")
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

func (t *CrossBorderChainCode) putTxnExchange(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("putTxnExchange is running ")

	if len(args) != 10 {
		return nil, errors.New("Incorrect Number of arguments.Expecting 10 for putTxnExchange")
	}
	txn := TxnExchange{
		Initiator:    args[3],
		Convertor:    args[4],
		SellCurrency: args[0],
		SellQuantity: args[5],
		BuyCurrency:  args[1],
		BuyQuantity:  args[5],
		ExchangeRate: args[2],
		Remarks:      args[6] + " - " + args[5],
		ID:           args[8],
		Time:         args[9],
	}

	bytes, err := json.Marshal(txn)
	if err != nil {
		fmt.Println("Error marshaling TxnExchange")
		return nil, errors.New("Error marshaling TxnExchange")
	}

	err = stub.PutState(txn.ID, bytes)
	if err != nil {
		return nil, err
	}

	return t.appendKey(stub, "TxnExchange", txn.ID)
}

func (t *CrossBorderChainCode) getAllTxnExchange(stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("getAllTxnExchange is running ")

	var txns []TxnExchange

	// Get list of all the keys - TxnExchange
	keysBytes, err := stub.GetState("TxnExchange")
	if err != nil {
		fmt.Println("Error retrieving TxnExchange keys")
		return nil, errors.New("Error retrieving TxnExchange keys")
	}
	var keys []string
	err = json.Unmarshal(keysBytes, &keys)
	if err != nil {
		fmt.Println("Error unmarshalling TxnExchange key")
		return nil, errors.New("Error unmarshalling TxnExchange keys")
	}

	// Get each txn from "TxnExchange" keys
	for _, value := range keys {
		bytes, err := stub.GetState(value)

		var txn TxnExchange
		err = json.Unmarshal(bytes, &txn)
		if err != nil {
			fmt.Println("Error retrieving txn " + value)
			return nil, errors.New("Error retrieving txn " + value)
		}

		fmt.Println("Appending txn goods details " + value)
		txns = append(txns, txn)
	}

	bytes, err := json.Marshal(txns)
	if err != nil {
		fmt.Println("Error marshaling txns TxnExchange")
		return nil, errors.New("Error marshaling txns TxnExchange")
	}
	return bytes, nil
}

func (t *CrossBorderChainCode) putTxnTransfer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("putTxnTransfer is running ")

	if len(args) != 7 {
		return nil, errors.New("Incorrect Number of arguments.Expecting 8 for putTxnTransfer")
	}
	txn := TxnTransfer{
		Sender:   args[0],
		Receiver: args[1],
		Remarks:  args[4],
		ID:       args[5],
		Time:     args[6],
		Value:    args[3],
		Asset:    args[2],
	}

	bytes, err := json.Marshal(txn)
	if err != nil {
		fmt.Println("Error marshaling TxnTransfer")
		return nil, errors.New("Error marshaling TxnTransfer")
	}

	err = stub.PutState(txn.ID, bytes)
	if err != nil {
		return nil, err
	}

	return t.appendKey(stub, "TxnTransfer", txn.ID)
}

func (t *CrossBorderChainCode) getAllTxnTransfer(stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("getAllTxnTransfer is running ")

	var txns []TxnTransfer

	// Get list of all the keys - TxnTransfer
	keysBytes, err := stub.GetState("TxnTransfer")
	if err != nil {
		fmt.Println("Error retrieving TxnTransfer keys")
		return nil, errors.New("Error retrieving TxnTransfer keys")
	}
	var keys []string
	err = json.Unmarshal(keysBytes, &keys)
	if err != nil {
		fmt.Println("Error unmarshalling TxnTransfer key")
		return nil, errors.New("Error unmarshalling TxnTransfer keys")
	}

	// Get each txn from "TxnTransfer" keys
	for _, value := range keys {
		bytes, err := stub.GetState(value)

		var txn TxnTransfer
		err = json.Unmarshal(bytes, &txn)
		if err != nil {
			fmt.Println("Error retrieving txn " + value)
			return nil, errors.New("Error retrieving txn " + value)
		}

		fmt.Println("Appending txn goods details " + value)
		txns = append(txns, txn)
	}

	bytes, err := json.Marshal(txns)
	if err != nil {
		fmt.Println("Error marshaling txns TxnTransfer")
		return nil, errors.New("Error marshaling txns TxnTransfer")
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
