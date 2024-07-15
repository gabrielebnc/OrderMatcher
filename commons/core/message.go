package core

type Message struct {
	from    string
	payload []byte
}

func (msg *Message) From() string {
	return msg.from
}

func (msg *Message) Payload() []byte {
	return msg.payload
}

func NewMessage(from string, payload []byte) *Message {
	return &Message{
		from:    from,
		payload: payload,
	}
}
