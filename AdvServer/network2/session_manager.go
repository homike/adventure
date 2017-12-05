package network2

import (
	"Adventure/AdvServer/model"
	"sync"
)

var SessionMgr *SessionManager

func init() {
	SessionMgr = &SessionManager{
		Sessions: make(map[uint]*Session),
	}
}

type SessionManager struct {
	sync.RWMutex
	Sessions map[uint]*Session
}

func (mgr *SessionManager) CreateSession(player *model.Player, conn *TCPClient) {
	sess := NewSession(player, conn)
	mgr.Sessions[player.AccountID] = sess
}
