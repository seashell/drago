package conn

import (
	"fmt"
	"testing"
)

func TestCreateConn(t *testing.T) {

	conn := NewRPCConnection("127.0.0.1:8081", nil)

	if err := conn.Call("TestMethod", struct{}{}, nil); err != nil {
		fmt.Println(err)
	}

	t.Log("success!")

}
