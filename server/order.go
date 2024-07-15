package main

import "time"

type Order struct {
	TraderId  [32]rune
	OrderType rune
	Side      rune
	Symbol    [8]rune
	Quantity  int16
	Price     int32

	//This is the time contained in the package, sent by the client
	ClientOrderTime time.Time
	//The time the server received the Order, will be different from ClientOrderTime
	ServerOrdertime time.Time
}
