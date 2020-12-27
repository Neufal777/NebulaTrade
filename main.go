package main

import (
	"time"

	"github.com/NebulaTrade/core"
)

const (
	XRP = "xrpeur"
	BTC = "btceur"
)

func main() {

	for {
		core.ExecuteMarket()
		time.Sleep(2 * time.Second)
	}

}
