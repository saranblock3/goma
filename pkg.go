package goma

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand/v2"
	"net"
	"os"
)

const HOMA_SOCKET_PATH string = "/tmp/homa.sock"
const HOMA_MESSAGE_HEADER_LENGTH uint64 = 32

type homaRegistrationMessage struct {
	id uint32
}

func (h *homaRegistrationMessage) marshal() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.LittleEndian, h.id); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type homaMessage struct {
	id                 uint64
	sourceAddress      [4]byte
	destinationAddress [4]byte
	sourceId           uint32
	destinationId      uint32
	length             uint64
	content            []byte
}

func (h *homaMessage) marshal() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.LittleEndian, HOMA_MESSAGE_HEADER_LENGTH+uint64(len(h.content))); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.LittleEndian, h.id); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, h.sourceAddress); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, h.destinationAddress); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, h.sourceId); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, h.destinationId); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, h.length); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, h.content); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (h *homaMessage) unmarshal(homaMessageBytes []byte) error {
	buf := bytes.NewReader(homaMessageBytes)

	if err := binary.Read(buf, binary.LittleEndian, &h.id); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &h.sourceAddress); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &h.destinationAddress); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &h.sourceId); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &h.destinationId); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &h.length); err != nil {
		return err
	}

	h.content = make([]byte, h.length)

	if err := binary.Read(buf, binary.LittleEndian, &h.content); err != nil {
		return err
	}

	return nil
}

type homaSocket struct {
	conn    net.Conn
	address [4]byte
	id      uint32
}

func NewHomaSocket(id uint32) (*homaSocket, error) {

	conn, err := net.Dial("unix", HOMA_SOCKET_PATH)
	if err != nil {
		return nil, err
	}

	registrationMessage := &homaRegistrationMessage{
		id,
	}
	registrationMessageBytes, err := registrationMessage.marshal()
	if err != nil {
		conn.Close()
		return nil, err
	}
	_, err = conn.Write(registrationMessageBytes)
	if err != nil {
		conn.Close()
		return nil, err
	}

	addressString, err := getLocalAddress()
	if err != nil {
		conn.Close()
		return nil, err
	}

	addressBytes := [4]byte{}
	_, err = fmt.Sscanf(addressString, "%d.%d.%d.%d", &addressBytes[0], &addressBytes[1], &addressBytes[2], &addressBytes[3])
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &homaSocket{
		conn,
		addressBytes,
		id,
	}, nil
}

func (h *homaSocket) WriteTo(content []byte, address string, id uint32) error {
	addressBytes := [4]byte{}
	_, err := fmt.Sscanf(address, "%d.%d.%d.%d", &addressBytes[0], &addressBytes[1], &addressBytes[2], &addressBytes[3])
	if err != nil {
		return err
	}
	messageId := rand.Uint64()
	fmt.Fprintln(os.Stderr, "-+-", messageId)

	homaMessage := &homaMessage{
		id:                 messageId,
		sourceAddress:      h.address,
		destinationAddress: addressBytes,
		sourceId:           h.id,
		destinationId:      id,
		length:             uint64(len(content)),
		content:            content,
	}

	homaMessageBytes, err := homaMessage.marshal()
	if err != nil {
		return err
	}

	_, err = h.conn.Write(homaMessageBytes)
	if err != nil {
		return err
	}
	return nil
}

func (h *homaSocket) Read() ([]byte, string, uint32, uint64, error) {
	sizeBytes := make([]byte, 8)
	n, err := h.conn.Read(sizeBytes)
	if err != nil || n < 8 {
		return nil, "", 0, 0, err
	}
	size := binary.LittleEndian.Uint64(sizeBytes)

	homaMessageBytes := make([]byte, size)
	n, err = h.conn.Read(homaMessageBytes)
	if err != nil || uint64(n) < size {
		return nil, "", 0, 0, err
	}

	homaMessage := &homaMessage{}
	err = homaMessage.unmarshal(homaMessageBytes)
	if err != nil {
		return nil, "", 0, 0, err
	}

	address := fmt.Sprintf(
		"%d.%d.%d.%d",
		homaMessage.sourceAddress[0],
		homaMessage.sourceAddress[1],
		homaMessage.sourceAddress[2],
		homaMessage.sourceAddress[3],
	)

	return homaMessage.content, address, homaMessage.sourceId, homaMessage.id, nil
}

func (h *homaSocket) Close() error {
	return h.conn.Close()
}
