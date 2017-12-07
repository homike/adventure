package network

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"reflect"
	"strconv"
)

//var MsgParserSingleton *MsgParser

// func init() {
// 	MsgParserSingleton = NewMsgParser()
// }

// --------------
// | len | data |
// --------------
type MsgParser struct {
	MsgProcessor *Processor
	MsgLen       int
	Endian       struct{}
}

func NewMsgParser() *MsgParser {
	return &MsgParser{
		MsgProcessor: newProcessor(),
		MsgLen:       6,
		Endian:       binary.LittleEndian,
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

	bodySize := headerSize - uint32(m.MsgLen)
	bodyData := make([]byte, bodySize)
	err = binary.Read(bufReader, binary.LittleEndian, &bodyData)
	if err != nil {
		log.Println("read body error")
		return 0, nil, err
	}

	return msgID, bodyData, nil
}

func (m *MsgParser) Write(msgID uint16, msgStruct interface{}) []byte {
	message := m.MsgProcessor.Marshal(msgStruct)

	w := bytes.NewBuffer([]byte{})
	binary.Write(w, binary.LittleEndian, uint32(len(message)+m.MsgLen))
	binary.Write(w, binary.LittleEndian, msgID)
	binary.Write(w, binary.LittleEndian, message)

	return w.Bytes()
	//client.Write(w.Bytes())
}

// self-defined protocol processor
type Processor struct {
}

func newProcessor() *Processor {
	return &Processor{}
}

func (p *Processor) UnMarshal(msgBody []byte, msgStruct interface{}) {
	readIndex := 0
	v := reflect.ValueOf(msgStruct).Elem()
	vType := v.Type()
	for i := 0; i < v.NumField(); i++ {
		vf := v.Field(i)
		tf := vType.Field(i)

		fmt.Println(tf.Name, vf.Kind())
		switch vf.Kind() {
		case reflect.String:
			for i := readIndex; i < len(msgBody); i++ {
				if msgBody[i] == byte(0) {
					fmt.Println(readIndex, "string :", msgBody[readIndex:i])
					canSetValue := reflect.ValueOf(string(msgBody[readIndex:i]))
					vf.Set(canSetValue)
					readIndex = i + 1
					break
				}
			}

		case reflect.Int32:
			fmt.Println(readIndex, "int :", msgBody[readIndex:readIndex+4])
			intValue, err := strconv.Atoi(string(msgBody[readIndex : readIndex+4]))
			if err != nil {
			}
			canSetValue := reflect.ValueOf(int32(intValue))
			vf.Set(canSetValue)
			readIndex = readIndex + 4
		default:
		}
	}
}

func (p *Processor) Marshal(msgStruct interface{}) []byte {

	bytesBuffer := bytes.NewBuffer([]byte{})

	v := reflect.ValueOf(msgStruct).Elem()
	for i := 0; i < v.NumField(); i++ {
		vf := v.Field(i)

		vfBytes := marshal(vf)
		binary.Write(bytesBuffer, binary.LittleEndian, vfBytes)
	}

	return bytesBuffer.Bytes()
}

func marshal(v reflect.Value) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})

	switch v.Kind() {
	case reflect.String:
		binary.Write(bytesBuffer, binary.LittleEndian, v.String())
		binary.Write(bytesBuffer, binary.LittleEndian, byte(0))

	case reflect.Uint8:
		binary.Write(bytesBuffer, binary.LittleEndian, uint8(v.Uint()))

	case reflect.Int32:
		binary.Write(bytesBuffer, binary.LittleEndian, int32(v.Int()))

	case reflect.Int64:
		binary.Write(bytesBuffer, binary.LittleEndian, v.Int())

	case reflect.Bool:
		b := 0
		if v.Bool() {
			b = 1
		}
		binary.Write(bytesBuffer, binary.LittleEndian, uint8(b))

	case reflect.Slice:
		binary.Write(bytesBuffer, binary.LittleEndian, int32(v.Len()))
		for j := 0; j < v.Len(); j++ {
			data := v.Slice(j, j+1).Index(0)
			sliceBytes := marshal(data)
			binary.Write(bytesBuffer, binary.LittleEndian, sliceBytes)
		}

	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			vf := v.Field(i)
			vfBytes := marshal(vf)
			binary.Write(bytesBuffer, binary.LittleEndian, vfBytes)
		}

	default:
		binary.Write(bytesBuffer, binary.LittleEndian, v.Bytes())
	}

	return bytesBuffer.Bytes()
}
