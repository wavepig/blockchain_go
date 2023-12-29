package server

import (
	"testing"
)

func TestBytes(t *testing.T) {
	command := "version"
	bytes := commandToBytes(command)
	t.Log("截取后：", bytes)
	toCommand := bytesToCommand(bytes)
	t.Log("转换回来后：", toCommand)

}
