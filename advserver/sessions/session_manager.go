package sessions

import (
	"adventure/advserver/log"
	"adventure/advserver/model"
	"adventure/advserver/network"
	"adventure/common/clog"
	"errors"
	"fmt"
	"sync"
)

var SessionMgr *SessionManager
var logger *clog.Logger

func Init() error {
	SessionMgr = &SessionManager{
		Sessions: make(map[uint]*Session),
	}

	logger = log.GetLogger()
	return nil
}

type SessionManager struct {
	sync.RWMutex
	Sessions map[uint]*Session
}

func (mgr *SessionManager) CreateSession(player *model.Player, conn *network.TCPClient) {
	sess := NewSession(conn)
	sess.SetPlayer(player)
	mgr.Sessions[player.AccountID] = sess
}

func (mgr *SessionManager) AddSession(sess *Session) {
	if sess == nil {
		return
	}
	mgr.Sessions[sess.PlayerData.AccountID] = sess
}

func (mgr *SessionManager) FindSession(AccountID uint) (*Session, error) {
	v, ok := mgr.Sessions[AccountID]
	if !ok {
		return nil, errors.New(fmt.Sprintf("%v Session not exist", AccountID))
	}
	return v, nil
}
