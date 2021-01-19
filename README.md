
# Crypto trading bot

**Disclaimer!**

Use this at your OWN risk! this is for educational purposes, and even it works, the author of this repository is not responsable for any loss for bad usage or market crash ! 

> At this moment it just works for binance exchange.

### How to use it?

 1. Change API credentials :
  ```bash
binanceaccount>config.go
```
 2. Run:
   ```bash
 go run main.go -currency=CURRENCY_TO_TRADE [ex. XLMBNB] -profit=DIFF_TO_SELL_BUY [ex. 0.00002]
```


 *Make sure you have funds in your corresponding binance wallets [in this case BNB available], if not, an "No funds available" error message will be shown*
 

```bash
go get github.com/Neufal777/InstagramExtractor
```

### How it works ?


*This bot is based on the usual strategy, **buy low, sell high**. with the difference that if, for example, you buy a cryptocurrency at **€ 0.20** at the time of execution, if after the purchase it continues to drop to **€ 0.19** or **€ 0.18**, what you will do to avoid losing the opportunity is to buy a little more, since in the first purchase only **33.33%** of our balance is bought, another **66.66%** remains since it is divided into **3** purchases (or more purchases or less, that already depends on the configuration that each one wants).  Once the purchases are made, what the bot will do is find the moment at which to sell X quantity at Y price to obtain a profit.  

Once one of the acquisitions is sold, another buy order at a lower price is set.

I leave you a drawing of how it would works*


![Strategy](https://www.IMAGE..com/)


### Generated file:
This function returns a struct called `User` with the data
```bash
{

"Symbol": "XLMBNB",
"Balance": 3.12093805,
"FirstBuy": 0.00602400,
"Available": 1.0403126833333334,
"Ammount": 0.968,
"Status": "BUY MORE",
"LastBuy": 0.00600000,
"LastSell": 0.00622400,
"Transactions": 1,
"Timer": 443,
"Limit": 1,
"Profit": 0.00002,
"OrdNum": 1,
"Orders": [
			{
				"Symbol": "XLMBNB",
				"OrderID": 81301675,
				"ClientOrderID": "pBhC0avrSFb9ZzmXpJPygm",
				"Price": "0.00602400",
				"Ammount": "166.00000000",
				"Status": "BUY",
				"Type": "LIMIT",
				"Active": 0
			},
			{
				"Symbol": "XLMBNB",
				"OrderID": 81301879,
				"ClientOrderID": "x0waz7EvrPxmvdApI4XQDS",
				"Price": "0.00600000",
				"Ammount": "167.00000000",
				"Status": "BUY",
				"Type": "LIMIT",
				"Active": 0
			}
		etc..
	]
}

```

Note: The Orders with *Active : 0*, means that the ammount was already sold and profit taken.

#### Conclusion..
_It is finished? of course NOT!, there's a looot to do and some things to fix but at the moment is completly functional and you can use it, but remember.. at your own RISK!  :)_
_Note: This was just for learning, that means that will not be mantained, feel free to add anything :)_
