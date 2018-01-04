package service

import (
	"adventure/advserver/sessions"
	"adventure/common/structs"
)

type Notice struct {
}

func (n *Notice) Broadcast(module string, textList []string, colorList []string, loops int32) {
	for _, v := range sessions.SessionMgr.Sessions {
		ntf := &structs.SystemAnnouncementNtf{
			Texts:           textList,
			Colors:          colorList,
			LoopCount:       loops,
			IsLeftDirection: true,
		}
		v.Send(structs.Protocol_SystemAnnouncement_Ntf, ntf)
	}
}

func (n *Notice) BroadcastRichTxt(module string, richText string, loops int32) {
	for _, v := range sessions.SessionMgr.Sessions {
		ntf := &structs.SystemAnnouncementRichNtf{
			Text:            richText,
			LoopCount:       loops,
			IsLeftDirection: true,
		}
		v.Send(structs.Protocol_SystemAnnouncementRich_Ntf, ntf)
	}
}
