package network2

const (
	Protocol_Test_Req                = 1
	Protocol_Test_Resp               = 2
	Protocol_GetSystemTime_Req       = 3
	Protocol_GetSystemTime_Resp      = 4
	Protocol_LoginServerResult_Ntf   = 1001
	Protocol_CreatePlayer_Req        = 1002
	Protocol_CreatePlayer_Resp       = 1003
	Protocol_SyncLoginDataFinish_Ntf = 1006
	Protocol_LoginServerPlatform_Req = 1007
	Protocol_SyncPlayerBaseInfo_Ntf  = 1008
	Protocol_NameExists_Req          = 1009
	Protocol_NameExists_Resp         = 1010
	Protocol_SyncUserGuidRecords_Ntf = 1413
)

type GetSystemTimeReq struct {
}

type GetSystemTimeResp struct {
	Time int64
}

type LoginServerResultNtf struct {
	Result         int32 // 0: Success
	IsCreatePlayer bool
}

type CreatePlayerReq struct {
	PlayerName     string
	HeroTemplateId int32
}

type CreatePlayerResp struct {
	Result int32 // 0: Success
}

type SyncLoginDataFinishNtf struct {
}

type LoginServerPlatformReq struct {
	Takon     string
	Version   int32
	ChannelID string
}

type SyncPlayerBaseInfoNtf struct {
	PlayerID           int32
	GameZoonID         int32 // 游戏分区ID
	IsSupperMan        bool  // 是否是GM
	PlatformType       int32 // 平台类型
	Viplevel           int32
	TotalRechargeIngot int32
}

type NameExistsReq struct {
	Name string
}

// 如果存在，则返回一个新名字，如果和传入的名字一样，则说明没有重名
type NameExistsResp struct {
	Name string
}

type GuildRecord struct {
	UserGuidTypes uint8
	TriggerCount  int32
}

type SyncUserGuidRecordsNtf struct {
	Records []GuildRecord
}
