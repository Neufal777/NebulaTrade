package console

import (
	"fmt"

	"github.com/NebulaTrade/config"
	"github.com/NebulaTrade/exchanges"
	"github.com/NebulaTrade/utils"
	"github.com/NebulaTrade/wallet"
	"github.com/ttacon/chalk"
)

//InformationDisplayConsole - displays recent information
func InformationDisplayConsole() {

	myWallet := wallet.ReadWallet()

	//get actual currency price
	latestCurrencyPrice := exchanges.BinancePrice(config.CURRENCY)
	currentPriceFloat := utils.StringToFloat(latestCurrencyPrice.Price)

	/*
		All the information available
	*/

	fmt.Println(chalk.Bold.TextStyle("Current Price:"), chalk.Bold.TextStyle(latestCurrencyPrice.Price), chalk.Green)
	fmt.Println(chalk.Bold.TextStyle("Last Buy:"), myWallet.LastBuy, chalk.Green)
	fmt.Println(chalk.Bold.TextStyle("Last Sell:"), myWallet.LastSell, chalk.Green)
	fmt.Println(chalk.Bold.TextStyle("Difference to sell:"), currentPriceFloat-myWallet.LastBuy, chalk.Green)
	fmt.Println(chalk.Bold.TextStyle("Ammount:"), myWallet.Ammount, chalk.Green)
	fmt.Println(chalk.Bold.TextStyle("Initial Balance:"), myWallet.Balance, chalk.Green)
	fmt.Println(chalk.Bold.TextStyle("Actual Balance:"), myWallet.Ammount*currentPriceFloat, chalk.Green)
	fmt.Println("Transactions:", myWallet.Transactions, chalk.Green)
	fmt.Println(chalk.Bold.TextStyle("-----------------------------"), chalk.Red)
}
