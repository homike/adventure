package msghandler

import (
	"adventure/advserver/network"
	"Adventure/common/structs"
	"fmt"
	"time"
)

// 3
func GetSystemTime(client *network.TCPClient, msgBody []byte) {

	timeNow := time.Now().Unix()
	resp := &structs.GetSystemTimeResp{
		Time: timeNow,
	}
	fmt.Println("czx@@@ GetSystemTime: ", timeNow)

	client.Write(structs.Protocol_GetSystemTime_Resp, resp)
}

// 1413
func SyncUserGuidRecords(client *network.TCPClient) {
	fmt.Println("czx@@@ SyncUserGuidRecords:")

	records := []structs.GuildRecord{}
	for i := 0; i < 2; i++ {
		records = append(records, structs.GuildRecord{
			UserGuidTypes: uint8(i + 2),
			TriggerCount:  int32(i + 3),
		})
	}
	resp := &structs.SyncUserGuidRecordsNtf{
		Records: records,
	}

	client.Write(structs.Protocol_SyncUserGuidRecords_Ntf, resp)
}
