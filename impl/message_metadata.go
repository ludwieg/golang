package impl

import (
	"bytes"
)

// MessageMetadata retains basic information about an in or outcoming message,
// such as the protocol version, message ID and type of package being sent.
type MessageMetadata struct {
	ProtocolVersion byte
	MessageID       byte
	PackageType     byte
}

func (m MessageMetadata) writeTo(buf *bytes.Buffer) error {
	buf.Write([]byte{m.ProtocolVersion, m.MessageID, m.PackageType})
	return nil
}
