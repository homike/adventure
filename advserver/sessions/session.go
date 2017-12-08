package sessions

import (
	"adventure/advserver/model"
	"adventure/advserver/network"
)

type Session struct {
	AccountID  uint
	Connection *network.TCPClient
	UserToken  string
	PlayerData *model.Player
}

func NewSession(player *model.Player, conn *network.TCPClient) *Session {

	sess := &Session{
		AccountID:  player.AccountID,
		PlayerData: player,
		Connection: conn,
	}

	return sess
}
