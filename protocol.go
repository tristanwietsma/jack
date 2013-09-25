package jackdb

import (
	"bytes"
	"fmt"
)

const (
	// EOM denotes the end of a message
	EOM = byte(0)

	// SEP denotes a delimitor between message units
	SEP = byte(30)

	// GET denotes a get-value-at-key command
	GET = byte('g')

	// SET denotes a set-key-to-value command
	SET = byte('s')

	// DEL denotes a delete-key command
	DEL = byte('d')

	// PUB denotes a publish command
	PUB = byte('P')

	// SUB denotes a subscribe command
	SUB = byte('S')

	// SUCCESS denotes a successful client command
	SUCCESS = byte('1')

	// FAIL denotes a failed client command
	FAIL = byte('0')
)

// COMMANDS is a list of all valid commands.
// The order is important: commands with index less than 3 (GET, DEL, SUB) require one additional argument (a key).
// Commands with index greater than or equal to 3 (SET, PUB) require two arguments (key, value).
// The index is used in Parse.
var COMMANDS = []byte{GET, DEL, SUB, SET, PUB}

// ProtocolError is the error type for a malformed message.
type ProtocolError struct {
	desc string
}

func (e ProtocolError) Error() string {
	return fmt.Sprintf("ProtocolError: %s", e.desc)
}

// Message is the server-side object representation of a client command.
type Message struct {
	cmd      byte
	key, arg []byte
}

// Bytes returns a byte slice representation of Message per protocol.
func (c *Message) Bytes() []byte {
	var b []byte
	if c.arg != nil {
		b = bytes.Join([][]byte{[]byte{c.cmd}, c.key, c.arg}, []byte{SEP})
	} else {
		b = bytes.Join([][]byte{[]byte{c.cmd}, c.key}, []byte{SEP})
	}
	return append(b, EOM)
}

// Parse accepts an incoming byte buffer and returns a Message (and error if malformed).
func Parse(b []byte) (*Message, error) {

	c := Message{}

	end := bytes.IndexByte(b, EOM)
	if end < 0 {
		err := ProtocolError{"Message is missing EOM byte."}
		return &c, err
	}

	units := bytes.Split(b[:end], []byte{SEP})
	numUnits := len(units)

	// command
	cid := bytes.IndexByte(COMMANDS, units[0][0])
	if cid < 0 {
		err := ProtocolError{"Message invokes a nonexistent command."}
		return &c, err
	}
	c.cmd = units[0][0]

	// key
	if numUnits < 2 {
		err := ProtocolError{"Message is missing key."}
		return &c, err
	}
	c.key = units[1]

	// GET, DEL, SUB
	if cid <= 2 && numUnits == 2 {
		return &c, nil
	}

	// SET, PUB
	if cid > 2 && numUnits == 3 {
		c.arg = units[2]
		return &c, nil
	}

	err := ProtocolError{"Message is nonsense."}
	return &c, err
}
