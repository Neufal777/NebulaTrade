package trading

import (
	"log"
	"time"

	"github.com/NebulaTrade/config"
	"github.com/NebulaTrade/exchanges"
	"github.com/NebulaTrade/mathnebula"
	"github.com/NebulaTrade/utils"
	"github.com/NebulaTrade/wallet"
)

//BeforeBuyingCrypto - before buying X crypto
func BeforeBuyingCrypto(currency string) float64 {

	//store the prices for x time
	var allPrices []float64
	for i := 0; i <= 29; i++ {

		CurrentPrice := exchanges.BinancePrice(currency)
		allPrices = append(allPrices, utils.StringToFloat(CurrentPrice.Price))
		time.Sleep(1 * time.Second)
	}

	//Once we have all the last 30 prices

	var total float64
	for _, p := range allPrices {

		total += p
	}

	priceToBuy := mathnebula.ToFixed(total/30, 8)

	return priceToBuy
}

//RecurrentBuy -
func RecurrentBuy() {

	/*
		if the status=="BUY"
		 - check last buy difference
		 - if current price is lower, new buy
	*/

	w := wallet.ReadWallet()
	last := w.LastBuy
	current := exchanges.BinancePrice(config.CURRENCY)

	diff := last - utils.StringToFloat(current.Price)

	if diff >= config.PROFIT {

		/*
			Buy more assets
		*/
		ammountToBuy := w.Available / utils.StringToFloat(current.Price)
		truncateAmmountToBuy := mathnebula.ToFixed((ammountToBuy), 7)
		ammountstring := utils.FloatToString(truncateAmmountToBuy)
		exchanges.ExecuteBuyOrderCURRENCY(ammountstring[:len(ammountstring)-13], utils.AnyTypeToString(current.Price), &w)

	} else {

		log.Println("waiting for price to drop")
	}
}
