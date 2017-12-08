package sessions

import (
	"adventure/advserver/model"
	"adventure/advserver/network"
	"sync"
)

var SessionMgr *SessionManager

func Init() {
	SessionMgr = &SessionManager{
		Sessions: make(map[uint]*Session),
	}
}

type SessionManager struct {
	sync.RWMutex
	Sessions map[uint]*Session
}

func (mgr *SessionManager) CreateSession(player *model.Player, conn *network.TCPClient) {
	sess := NewSession(player, conn)
	mgr.Sessions[player.AccountID] = sess
}
