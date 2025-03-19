package goma

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand/v2"
	"net"
	"sync"
	"time"
)

const HOMA_SOCKET_PATH string = "/tmp/homa.sock"
const HOMA_MESSAGE_HEADER_LENGTH uint64 = 24
const HOMA_MESSAGE_MAX_LENGTH uint64 = 524_288_000

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
	conn   net.Conn
	id     uint32
	sendMu sync.Mutex
	recvMu sync.Mutex
}

func NewHomaSocket(id uint32) (*homaSocket, error) {
	randomSleep()
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

	return &homaSocket{
		conn,
		id,
		sync.Mutex{},
		sync.Mutex{},
	}, nil
}

func (h *homaSocket) SendTo(content []byte, sourceAddress string, destinationAddress string, id uint32) error {
	randomSleep()
	h.sendMu.Lock()
	defer h.sendMu.Unlock()
	sourceAddressBytes := [4]byte{}
	_, err := fmt.Sscanf(sourceAddress, "%d.%d.%d.%d", &sourceAddressBytes[0], &sourceAddressBytes[1], &sourceAddressBytes[2], &sourceAddressBytes[3])
	if err != nil {
		h.conn.Close()
		return err
	}

	destinationAddressBytes := [4]byte{}
	_, err = fmt.Sscanf(destinationAddress, "%d.%d.%d.%d", &destinationAddressBytes[0], &destinationAddressBytes[1], &destinationAddressBytes[2], &destinationAddressBytes[3])
	if err != nil {
		h.conn.Close()
		return err
	}

	homaMessage := &homaMessage{
		sourceAddress:      destinationAddressBytes,
		destinationAddress: sourceAddressBytes,
		sourceId:           h.id,
		destinationId:      id,
		length:             uint64(len(content)),
		content:            content,
	}

	homaMessageBytes, err := homaMessage.marshal()
	if err != nil {
		h.conn.Close()
		return err
	}

	sent := 0
	for sent < len(homaMessageBytes) {
		n, err := h.conn.Write(homaMessageBytes[sent:])
		if err != nil {
			h.conn.Close()
			return errors.New("message not sent")
		}
		sent += n
	}
	return nil
}

func (h *homaSocket) RecvFrom() ([]byte, string, string, uint32, error) {
	randomSleep()
	h.recvMu.Lock()
	defer h.recvMu.Unlock()
	sizeBytes := make([]byte, 8)
	n, err := h.conn.Read(sizeBytes)
	if err != nil || n < 8 {
		h.conn.Close()
		return nil, "", "", 0, errors.New("message size not received")
	}
	size := binary.LittleEndian.Uint64(sizeBytes)
	if size > HOMA_MESSAGE_MAX_LENGTH {
		h.conn.Close()
		return nil, "", "", 0, errors.New("message exceeded maximum size")
	}

	homaMessageBytes := make([]byte, size)
	var received uint64 = 0
	for received < size {
		n, err = h.conn.Read(homaMessageBytes[received:])
		if err != nil {
			h.conn.Close()
			return nil, "", "", 0, errors.New("message not read")
		}
		received += uint64(n)
	}

	homaMessage := &homaMessage{}
	err = homaMessage.unmarshal(homaMessageBytes)
	if err != nil {
		h.conn.Close()
		return nil, "", "", 0, err
	}

	sourceAddress := fmt.Sprintf(
		"%d.%d.%d.%d",
		homaMessage.sourceAddress[0],
		homaMessage.sourceAddress[1],
		homaMessage.sourceAddress[2],
		homaMessage.sourceAddress[3],
	)

	destinationAddress := fmt.Sprintf(
		"%d.%d.%d.%d",
		homaMessage.destinationAddress[0],
		homaMessage.destinationAddress[1],
		homaMessage.destinationAddress[2],
		homaMessage.destinationAddress[3],
	)

	return homaMessage.content, sourceAddress, destinationAddress, homaMessage.sourceId, nil
}

func (h *homaSocket) Close() error {
	return h.conn.Close()
}

func randomSleep() {
	time.Sleep(time.Duration(rand.IntN(1000)) * time.Nanosecond)
}
