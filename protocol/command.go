package protocol

import "github.com/sirupsen/logrus"

const (
	cmdVersion    = "version"
	cmdVerack     = "verack"
	commandLength = 12
)

var commands = map[string][commandLength]byte{
	cmdVersion: newCommand(cmdVersion),
	cmdVerack:  newCommand(cmdVerack),
}

func newCommand(command string) [commandLength]byte {
	l := len(command)
	if l > commandLength {
		logrus.Fatalf("command %s is too long", command)
	}

	var packed [commandLength]byte
	buf := make([]byte, commandLength-l)
	copy(packed[:], append([]byte(command), buf...)[:])

	return packed
}
