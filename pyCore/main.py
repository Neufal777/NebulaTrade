import requests

x = requests.get('https://api-pub.bitfinex.com/v2/tickers?symbols=tBTCUSD')
print(x.content[2:])