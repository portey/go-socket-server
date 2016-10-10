package server

const ACTION_AUTH string = "AUTH";

type Message struct {
	Action string `json:"action"`
	Data   interface{} `json:"data"`
}

func (message *Message) String() string {
	return "Action: " + message.Action
}

func (message *Message) Attribute(attribute string) string {
	return "test"
}
