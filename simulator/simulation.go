package simulator

import (
	"github.com/NebulaTrade/core"
)

//Wallet - contains wallet & info about it
type Wallet struct {
	Name    string
	Initial string
	Profit  string
	Balance string
}

//SimulateNebula -
func SimulateNebula() {

	for x := 1; x < 3; x++ {

		core.ExecuteMarket()
	}
}
