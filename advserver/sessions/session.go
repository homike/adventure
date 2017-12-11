package sessions

import (
	"adventure/advserver/model"
	"adventure/advserver/network"
)

type Session struct {
	AccountID  uint
	UserToken  string
	PlayerData *model.Player
	Connection *network.TCPClient
}

func NewSession(conn *network.TCPClient) *Session {

	sess := &Session{
		Connection: conn,
	}

	return sess
}

func (sess *Session) SetPlayer(player *model.Player) error {
	sess.Connection.AccountID = player.AccountID
	sess.PlayerData = player
	sess.AccountID = player.AccountID
	return nil
}

func (sess *Session) Send(msgID uint16, msgStruct interface{}) error {
	sess.Connection.Write(msgID, msgStruct)
	return nil
}

func (sess *Session) UnMarshal(msgBody []byte, msgStruct interface{}) {
	sess.Connection.UnMarshal(msgBody, msgStruct)
}
