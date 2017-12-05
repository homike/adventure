package network2

import (
	"Adventure/AdvServer/model"
)

type Session struct {
	AccountID  uint
	Connection *TCPClient
	UserToken  string
	PlayerData *model.Player
}

func NewSession(player *model.Player, conn *TCPClient) *Session {

	sess := &Session{
		AccountID:  player.AccountID,
		PlayerData: player,
		Connection: conn,
	}

	return sess
}
