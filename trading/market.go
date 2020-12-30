package trading

import (
	"fmt"
	"strconv"

	"github.com/ttacon/chalk"

	"github.com/NebulaTrade/exchanges"
	"github.com/NebulaTrade/utils"
	"github.com/NebulaTrade/wallet"
)

const (
	//PROFITPERTRANS - profit we want for each transfer
	PROFITPERTRANS = 0.00001
)

//DecisionMakeBuy - where the decisions of buying or selling is made
func DecisionMakeBuy(w *wallet.Wallet) {

	//Check at how much we sold and if the actual price is lower
	buyActualPrice := exchanges.BinancePrice("LUNABNB")
	lastSellFloat := wallet.GetLastSell()
	lastPriceFloat := utils.StringToFloat(buyActualPrice.Price)

	/*
		difference betwen last sell and actual
		to check if the price dropped
	*/
	difference := lastSellFloat - lastPriceFloat
	fmt.Println("Waiting for price drop to buy...", chalk.Red)
	fmt.Println("LIVE Price at: "+strconv.FormatFloat(lastPriceFloat, 'f', 12, 64), chalk.Green)
	fmt.Println("LAST Sell at: "+strconv.FormatFloat(lastSellFloat, 'f', 12, 64), chalk.Green)

	/*
		Showing information about the account
		 - Ammount of crypto
		 - Actual Balance
		 - Transaction
	*/

	fmt.Println("Ammount: ", w.Ammount, chalk.Green)
	fmt.Println("Actual Balance:", w.Balance, chalk.Green)
	fmt.Println("Transactions:", w.Transactions, chalk.Green)
	fmt.Println("-----------------------------", chalk.Red)

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

		fmt.Println("Buy order executed!", chalk.Green)

		fmt.Println("Ammount: ", w.Ammount, chalk.Green)
		fmt.Println("Actual Balance:", w.Balance, chalk.Green)
		fmt.Println("Transactions:", w.Transactions, "\n", chalk.Green)
		fmt.Println("-----------------------------", chalk.Red)

	} else {

		// fmt.Println("Actual Balance:", w.Balance, chalk.Green)
		// fmt.Println("Transactions:", w.Transactions, chalk.Green)
		// fmt.Println("Actual Ammount:", w.Ammount, chalk.Green)
		// fmt.Println("Waiting for the price to drop.. ", chalk.Green)
		// fmt.Println(, chalk.Green)

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
	data := exchanges.BinancePrice("LUNABNB")

	currentPriceFloat, _ := strconv.ParseFloat(data.Price, 32)
	/*
		Information about the last BUY
	*/
	lastBuyFloat := wallet.GetLastBuy()

	differenceToSell := currentPriceFloat - lastBuyFloat

	/*
		Displaying information in the console
	*/
	fmt.Println("Waiting to sell, price now at:", data.Price, chalk.Green)
	fmt.Println("Last Buy:", lastBuyFloat, chalk.Green)
	fmt.Println("Difference to sell:", differenceToSell, chalk.Green)
	fmt.Println("Ammount:", w.Ammount, chalk.Green)
	fmt.Println("Balance:", w.Balance, chalk.Green)
	fmt.Println("Transactions:", w.Transactions, chalk.Green)
	fmt.Println("-----------------------------", chalk.Red)

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

		w.LastSell = currentPriceFloat

		w.WriteInWallet()

		fmt.Println("SOLD!", chalk.Green)
		w.Transactions++
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
