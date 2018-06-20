package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/homike/cuttletest/framework"
	RB "github.com/homike/cuttletest/robot"
)

func SignalProc() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT /*, syscall.SIGUSR1, syscall.SIGUSR2*/)

	for {
		<-ch
		os.Exit(0)
		return
	}
}

func RunEmpty(robot *RB.Robot, stepId int) {

}

var arrActFunc = []framework.RunCaseInfo{
	{RunEmpty, 1},
}

func main() {
	go SignalProc()

	robots := FanInRobot()

	DoTest(robots)
}

var curStartTime int64
var nextStartTime int64

type Robot struct {
	RobotIndex int
	SendCnt    int
	RecvCnt    int
}

func FanInRobot() chan *Robot {

	curStartTime = time.Now().UnixNano() / 1000000
	nextStartTime = curStartTime + (int64)(1000*2)
	robots := make(chan *Robot, 1)

	for i := 0; i < 1; i++ {
		go func(index int) {

			robot := &Robot{
				RobotIndex: index,
				SendCnt:    0,
				RecvCnt:    0,
			}

			robots <- robot
		}(i)
	}

	return robots
}

func DoTest(robots chan *Robot) {

	var count sync.WaitGroup
	timeStart := time.Now().Unix()

	count.Add(1000)

	go func() {
		count.Wait()
		close(robots)
	}()

	for r := range robots {
		r := r

		go func() {
			// Connect
			conn, err := net.Dial("tcp", "127.0.0.1:9110")
			if err != nil {
				fmt.Println("connect error", err.Error())
				return
			}
			//fmt.Println(r.RobotIndex, " connect success")

			defer func() {
				fmt.Println("robotIndex: ", r.RobotIndex, " Exit")
				conn.Close()
				count.Done()
			}()

			bufReader := bufio.NewReader(conn)
			for {
				Write(conn, 1)

				r.SendCnt++
				_, _, err := Read(bufReader)
				if err != nil {
					fmt.Println("gate message read error")
					return
				}

				r.RecvCnt++
				//fmt.Println("robotIndex1: ", r.RobotIndex, "recv count: ", r.RecvCnt)
				//time.Sleep(1 * time.Millisecond)

				//if r.RecvCnt >= 10 {
				//return
				//}
			}
		}()
	}

	timeEnd := time.Now().Unix()
	fmt.Println("cost time: ", timeEnd-timeStart)
}

func Write(conn net.Conn, msgID uint16) {
	//message := []byte{}
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, uint32(10))
	message := bytesBuffer.Bytes()

	writer := bufio.NewWriter(conn)
	binary.Write(writer, binary.LittleEndian, uint32(len(message)+6))
	binary.Write(writer, binary.LittleEndian, msgID)
	binary.Write(writer, binary.LittleEndian, message)
	writer.Flush()
}

func Read(bufReader *bufio.Reader) (uint16, []byte, error) {
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
	bodySize := headerSize - 6
	bodyData := make([]byte, bodySize)
	err = binary.Read(bufReader, binary.LittleEndian, &bodyData)
	if err != nil {
		log.Println("read body error")
		return 0, nil, err
	}

	return msgID, bodyData, nil
}
