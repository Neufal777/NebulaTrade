package wallet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/NebulaTrade/config"
)

//Order - defines struct of an order
type Order struct {
	Symbol        string
	OrderID       int64
	ClientOrderID string
	Price         string
	Ammount       string
	Status        string //BUY OR SELL ORDER
	Type          string //LIMIT, STOP, ETC
	Active        int
}

//Wallet - wallet of each user
type Wallet struct {
	Symbol       string
	Balance      float64
	Available    float64
	Ammount      float64
	Status       string
	LastBuy      float64
	LastSell     float64
	Transactions int32
	Timer        int
	Limit        float64
	Profit       float64
	OrdNum       int
	Orders       []Order
}

//WriteInWallet - updates the information
func (wallet *Wallet) WriteInWallet() {

	file, err := json.MarshalIndent(wallet, "", " ")

	if err != nil {
		log.Panic(err)
	}
	_ = ioutil.WriteFile(config.CURRENCY+".json", file, 0644)
}

//ReadWallet - gets all the information from the wallet and returns a wallet struct
func ReadWallet() Wallet {

	// Open our jsonFile
	jsonFile, err := os.Open(config.CURRENCY + ".json")

	if err != nil {
		fmt.Println(err)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	fileValue, _ := ioutil.ReadAll(jsonFile)

	var wallet Wallet

	json.Unmarshal(fileValue, &wallet)

	return wallet

}

/*
	Here is where we get all the data
	from OUR current wallet (wallet.json file)
*/

//GetLastSell - gets the price of your last sell
func GetLastSell() float64 {

	w := ReadWallet()

	/*
		1- Get the data from the last sell
		2- Convert it from string to float64
	*/

	return w.LastSell
}

//GetLastBuy - gets the price of your last buy
func GetLastBuy() float64 {

	w := ReadWallet()

	/*
		1- Get the data from the last Buy
		2- Convert it from string to float64
	*/

	return w.LastBuy
}

//GetStatus - gets the status of your last buy
func GetStatus() string {

	w := ReadWallet()

	/*
		--- Get the Status (ready to BUY or SELL)
	*/
	return w.Status
}
