package config

/*
	Here we have the configuration
	- PROFIT PER TRANS && CURRENCY
*/

var (
	PROFIT   = 0.000001
	CURRENCY = "CURRENCY"
	COUNTER  = 100 //how much we wait to buy new stock [in seconds] #26 minutes
)

//go run main.go -currency=MITHBNB -limit=2 -profit=0.00001
//go run main.go -currency=XLMBNB -limit=2 -profit=0.00001
