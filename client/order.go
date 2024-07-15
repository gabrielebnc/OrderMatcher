package main

import "time"

// Order from the Client's perspective,
// it contains less information than the Server's one
type Order struct {
	TraderId  [32]rune
	OrderType rune
	Side      rune
	Symbol    [8]rune
	Quantity  int16
	Price     int32
	OrderTime time.Time
}
