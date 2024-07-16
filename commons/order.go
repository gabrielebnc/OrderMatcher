package commons

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Order from the Client's perspective,
// it contains less information than the Server's one
type Order struct {
	TraderId  [32]rune
	OrderType rune
	Side      rune
	Symbol    [8]rune
	Quantity  int16
	Price     int32
	OrderTime int64
}

func (o Order) Serialize() []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, o)
	if err != nil {
		fmt.Println("ERROR (Order Serialize):", err)
		return []byte{0}
	}
	fmt.Printf("Bytes: %v\n", buf.Bytes())

	return buf.Bytes()
}
