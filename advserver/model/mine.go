package model

import (
	"adventure/advserver/gamedata"
	"adventure/common/structs"
	"math"
	"time"
)

type MineMap struct {
	MapData  *structs.MineMap      // 挖矿地图
	UserData *structs.UserMineData // 玩家挖矿数据
}

func NewUserMineData() *MineMap {
	mine := &MineMap{
		MapData: &structs.MineMap{
			NodeList: []*structs.NodeList{},
			BossIDs:  []int32{},
			DigCnt:   0,
		},
		UserData: &structs.UserMineData{
			ExpandSight: 0,
			StatueLv:    1,
			DigNode: &structs.DigBlockNode{
				NodeID: -1,
				X:      127,
			},
			MinePickLv:      1,
			LastResetDate:   0,
			LastRefreshDate: 0,
			Boss:            &structs.BossNode{},
			DigQueueIDs:     []int32{},
			DigProxys:       []structs.DigProxy{},
		},
	}

	mine.ResetMineMap(true)

	return mine
}

func (m *MineMap) ResetMineMap(resetDepth bool) {
	m.MapData.NodeList = []*structs.NodeList{}

	minePickT, ok := gamedata.AllTemplates.MinePickTemplates[m.UserData.MinePickLv]
	if !ok {
		return
	}

	maxDigDepth := int32(0)
	m.UserData.LvMinePickMax = minePickT.DigNum
	for _, v := range gamedata.AllTemplates.DefaultTileTemplates {
		for i := 0; i < len(v.ListX); i++ {
			m.AddBlock(v.NodeID, int8(v.ListX[i]), int8(v.ListY[i]), true)
			maxDigDepth = int32(math.Max(float64(maxDigDepth), float64(v.ListY[i])))
		}
	}

	if resetDepth {
		m.UserData.DigDepthMax = maxDigDepth
	}
	m.UserData.LastResetDate = time.Now().Unix()
}

// 添加资源
func (m *MineMap) AddBlock(nodeID int32, x, y int8, isVisible bool) bool {
	_, err := m.addBlock(nodeID, x, y, isVisible)
	if err != nil {
	}

	tileT, ok := gamedata.AllTemplates.TileTemplates[nodeID]
	if ok && tileT.TileType == structs.TileType_Statue {
		m.UserData.StatueCnt++
		return true
	}

	return false
}

func (m *MineMap) addBlock(nodeID int32, x, y int8, isVisible bool) (*structs.NodeList, error) {

	var node *structs.NodeList

	for _, v := range m.MapData.NodeList {
		if v.NodeID == nodeID {
			node = v
		}
	}

	if node == nil {
		node = &structs.NodeList{
			NodeID: 0,
			Nodes:  []*structs.BlockNode{},
		}
		m.MapData.NodeList = append(m.MapData.NodeList, node)
	}

	node.AddBlockNode(x, y, isVisible)

	return node, nil
}

func (m *MineMap) GetNodeListByType(nodeID int32) *structs.NodeList {
	for _, v := range m.MapData.NodeList {
		if v.NodeID == nodeID {
			return v
		}
	}

	return nil
}
