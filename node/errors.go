package node

import "fmt"

type RPCError struct {
	Message string
	Method  string
	URL     string
}

func (e *RPCError) Error() string {
	return fmt.Sprintf("Node: %s, Method: %s, Message: %s", e.URL, e.Method, e.Message)
}
