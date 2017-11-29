package main

import (
	"Adventure/AdvServer/network"
)

func main() {
	// resp := &network.SyncPlayerBaseInfoNtf{
	// 	PlayerID:           1,
	// 	GameZoonID:         1,
	// 	IsSupperMan:        true,
	// 	PlatformType:       1,
	// 	Viplevel:           1,
	// 	TotalRechargeIngot: 1,
	// }
	// byteDatas := MsgMarshal(resp)
	// fmt.Println("byteDatas :", len(byteDatas), "datas :", byteDatas)

	//testData := []byte{116, 97, 107, 111, 110, 49, 0, 0, 0, 0, 0, 56, 48, 49, 48, 48, 48, 50, 48, 56, 0}
	//network.LoginServerPlatform(nil, testData)

	listener := network.NewTCPListenter()
	listener.StartAccept()
}
