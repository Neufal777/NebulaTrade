package trading

import (
	"log"
	"time"

	"github.com/NebulaTrade/exchanges"
	"github.com/NebulaTrade/mathnebula"
	"github.com/NebulaTrade/utils"
)

//BeforeBuyingCrypto - before buying X crypto
func BeforeBuyingCrypto(currency string) {

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

	log.Println(priceToBuy)
}
