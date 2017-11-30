package network

import (
	"bufio"
	"encoding/binary"
	"log"
)

// --------------
// | len | data |
// --------------
type MsgParser struct {
	MsgLen int
	Endian struct{}
}

func NewMsgParser() *MsgParser {
	return &MsgParser{
		MsgLen: 6,
		Endian: binary.LittleEndian,
	}
}

func (m *MsgParser) Read(bufReader *bufio.Reader) (uint16, []byte, error) {
	var headerSize uint32
	err := binary.Read(bufReader, binary.LittleEndian, &headerSize)
	if err != nil {
		log.Println("read headsize error")
		return 0, nil, err
	}

	var msgID uint16
	err = binary.Read(bufReader, binary.LittleEndian, &msgID)
	if err != nil {
		log.Println("read msgid error")
		return 0, nil, err
	}

	//fmt.Println("headerSize", headerSize, "msgID:", msgID)
	bodySize := headerSize - uint32(m.MsgLen)
	bodyData := make([]byte, bodySize)
	err = binary.Read(bufReader, binary.LittleEndian, &bodyData)
	if err != nil {
		log.Println("read body error")
		return 0, nil, err
	}

	return msgID, bodyData, nil
}

func (m *MsgParser) Write(w *bufio.Writer, msgID uint16, msgStruct interface{}) {
	message := MsgMarshal(msgStruct)

	binary.Write(w, binary.LittleEndian, uint32(len(message)+6))
	binary.Write(w, binary.LittleEndian, msgID)
	binary.Write(w, binary.LittleEndian, message)
	w.Flush()
}
