package bridge

// State enum stores the bridge state
type State uint16

const (
	// StateInit is the initial bridge state = 0
	StateInit State = 0 + iota
	// StateNew is the state after calling new bridge = 1
	StateNew
	// StateWaitForOffer is when waiting for an offer
	StateWaitForOffer
	// StateWaitForAnswer is waiting for an answer = 3
	StateWaitForAnswer
	// StateSendAnswer offer received, send answer = 4
	StateSendAnswer
	// StateWaitForRTC offer received and answer sent  = 5
	StateWaitForRTC
	// StateEstablishRTC offer sent and answer received  = 6
	StateEstablishRTC
	// StateOpenConnection when WebRTC is connected and open = 7
	StateOpenConnection
	// StateDisconnected when WebRTC goes from connected to diconnected and open = 8
	StateDisconnected
)

func (s State) String() string {
	switch s {
	case StateInit:
		return "StateInit"
	case StateNew:
		return "StateNew"
	case StateWaitForAnswer:
		return "StateWaitForAnswer"
	case StateSendAnswer:
		return "StateSendAnswer"
	case StateWaitForRTC:
		return "StateWaitForRTC"
	case StateEstablishRTC:
		return "StateEstablishRTC"
	case StateOpenConnection:
		return "StateOpenConnection"
	default:
		return "Undefined"
	}
}
