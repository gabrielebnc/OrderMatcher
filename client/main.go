package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/gabrielebnc/OrderMatcher/client/transport"
)

func randomOrder(traderId string) *Order {

	var tId [32]rune
	var orderType rune
	var side rune
	var symbol [8]rune

	if rand.Int()%2 == 0 {
		side = 'B'
	} else {
		side = 'S'
	}

	if rand.Int()%2 == 0 {
		orderType = 'L'
	} else {
		orderType = 'M'
	}

	for i, char := range "0000AAPL" {
		symbol[i] = char
	}

	for i, char := range traderId {
		tId[i] = char
	}

	return &Order{
		TraderId:  tId,
		OrderType: orderType,
		Side:      side,
		Symbol:    symbol,
		Quantity:  int16(rand.Int()),
		Price:     rand.Int31(),
		OrderTime: time.Now(),
	}
}

func main() {
	serverAddr := ":3000"
	traderId := "00000000000000000000000002141255"

	client := transport.NewTCPClient()

	client.Dial(serverAddr)

	go func() {
		for msg := range client.Consume() {
			fmt.Printf("msg from %v: %v\n", msg.From(), string(msg.Payload()))
		}
	}()

	fmt.Println(randomOrder(traderId))

	for i := 0; i < 12; i++ {
		data, err := json.Marshal(*randomOrder(traderId))

		if err != nil {
			continue
		}

		client.SendMessage(serverAddr, data)
	}

	client.CloseConnection(serverAddr)
}
