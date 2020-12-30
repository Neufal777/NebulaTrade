package trading

import (
	"fmt"
	"strconv"

	"github.com/ttacon/chalk"

	"github.com/NebulaTrade/console"
	"github.com/NebulaTrade/exchanges"
	"github.com/NebulaTrade/utils"
	"github.com/NebulaTrade/wallet"
)

const (
	//PROFITPERTRANS - profit we want for each transfer
	PROFITPERTRANS = 0.000002
)

//DecisionMakeBuy - where the decisions of buying or selling is made
func DecisionMakeBuy(w *wallet.Wallet) {

	//Check at how much we sold and if the actual price is lower
	buyActualPrice := exchanges.BinancePrice(exchanges.MITHBNB)
	lastSellFloat := wallet.GetLastSell()
	lastPriceFloat := utils.StringToFloat(buyActualPrice.Price)

	/*
		difference betwen last sell and actual
		to check if the price dropped
	*/
	difference := lastSellFloat - lastPriceFloat
	fmt.Println(chalk.Bold.TextStyle("Waiting for price drop to buy..."), chalk.Red)

	/*
		Showing information about the account
		 - Ammount of crypto
		 - Actual Balance
		 - Transaction...
	*/
	console.InformationDisplayConsole()
	if difference >= PROFITPERTRANS {

		/*
			EXECUTE BUY ORDER
		*/
		//execorder.BuyOrder()
		/*
			change last sell file with updated info
			with the lastPriceFloat
		*/

		// lastPriceStringToWrite := utils.FloatToString(lastPriceFloat)
		// utils.WriteFile(lastPriceStringToWrite, "core/lastBuy.txt")

		w.LastBuy = lastPriceFloat
		w.Status = "SELL"
		w.Transactions++
		w.Ammount = w.Balance / lastPriceFloat

		w.WriteInWallet()
		//change status from buying to selling
		//_ = ioutil.WriteFile("core/status.txt", []byte("SELL"), 0)

		/*
			- Displaying information
			- Details about the *wallet.Wallet
		*/

		fmt.Println(chalk.Bold.TextStyle("BUY ORDER EXECUTED!"), chalk.Green)
		console.InformationDisplayConsole()

	}
}

//DecisionMakeSell - once we've bought new crypto, we wait to sell them
func DecisionMakeSell() {

	/*
		Get latest data from the coin
		To evaluate when to sell the coin
		and change status from buying to selling
	*/

	w := wallet.ReadWallet()
	data := exchanges.BinancePrice(exchanges.MITHBNB)

	currentPriceFloat, _ := strconv.ParseFloat(data.Price, 32)
	/*
		Information about the last BUY
	*/

	differenceToSell := currentPriceFloat - w.LastBuy

	/*
		Displaying information in the console
	*/
	fmt.Println(chalk.Bold.TextStyle("Waiting to sell.."), chalk.Green)
	console.InformationDisplayConsole()

	if differenceToSell >= PROFITPERTRANS {

		/*
			EXECUTE SELL ORDER
		*/
		//execorder.SellOrder()

		/*
			Change 2 files:
				- Last Sell
				- Ststus to BUY
		*/

		w.Balance = w.Ammount * currentPriceFloat
		w.Ammount = 0
		w.Status = "BUY"
		w.Transactions++
		w.LastSell = currentPriceFloat

		w.WriteInWallet()

		fmt.Println(chalk.Bold.TextStyle("SOLD!"), chalk.Green)

	}

}

//ExecuteMarket - workflow of the program
func ExecuteMarket(w *wallet.Wallet) {

	/*
		Check the status (BUY OR SELL)
	*/

	actualStatusString := wallet.GetStatus()

	/*
		depending on the status,
		we execute buy or sell orders
	*/

	switch actualStatusString {
	case "BUY":
		DecisionMakeBuy(w)
	case "SELL":
		DecisionMakeSell()
	}

}
