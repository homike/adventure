package msghandler

import (
	"adventure/advserver/gamedata"
	"adventure/advserver/sessions"
	"adventure/common/structs"
	"adventure/common/util"
)

func InitMessageMine() {
	handlers := map[uint16]ProcessFunc{
		uint16(structs.Protocol_SyncMiningMap_Req):  SyncMiningMapReq,
		uint16(structs.Protocol_ResetMiningMap_Req): ResetMiningMapReq,
		uint16(structs.Protocol_DigRequest_Req): DigReq,
	}

	for k, v := range handlers {
		MapFunc[k] = v
	}
}

func SyncMiningMapReq(sess *sessions.Session, msgBody []byte) {
	sess.SyncMineMapData()
}

func ResetMiningMapReq(sess *sessions.Session, msgBody []byte) {
	req := structs.ResetMiningMapReq{}
	sess.UnMarshal(msgBody, &req)

	resp := &structs.ResetMiningMapResp{
		Ret: structs.AdventureRet_Failed,
	}

	/////////////////////////////////////////////Data Check////////////////////////////////////////
	shopPrice := int32(0)
	seconds := util.TimeSubNow(sess.PlayerData.MineMap.UserData.LastResetDate) - 1
	if seconds < gamedata.AllTemplates.GlobalData.MineMapResetDelayTime {
		if !req.UserIgnot {
			sess.Send(structs.Protocol_ResetMiningMap_Resp, resp)
			return
		}
		shopT, ok := gamedata.AllTemplates.ShopTemplates[int32(structs.ShopEnum_MiningMapReset)]
		if !ok {
			sess.Send(structs.Protocol_ResetMiningMap_Resp, resp)
			return
		}
		if sess.PlayerData.Res.Ingot < shopT.ShopPrice {
			sess.Send(structs.Protocol_ResetMiningMap_Resp, resp)
			return
		}
	}

	///////////////////////////////////////////Logic Process///////////////////////////////////////
	if shopPrice > 0 {
		sess.IgnotChange(shopPrice, 0)
	}

	userMine := sess.PlayerData.MineMap
	userMine.UserData.DigNode.NodeID = -1
	userMine.UserData.DigNode.X = int8(gamedata.AllTemplates.GlobalData.MineMapWidth / 2)
	userMine.UserData.DigNode.Y = 0
	userMine.UserData.Boss.Status = structs.BossStatus_NoAppear
	userMine.UserData.StatueLv = 1
	userMine.UserData.StatueCnt = 0

	userMine.MapData.BossIDs = []int32{}
	userMine.ResetMineMap(false)

	//  同步客户端巨魔数量
	sess.Send(structs.Protocol_UpdateStatueLevelAndCount_Ntf, &structs.UpdateStatueLevelAndCountNtf{userMine.UserData.StatueLv, userMine.UserData.StatueCnt})

	// 重置完成
	resp.Ret = structs.AdventureRet_Success
	sess.Send(structs.Protocol_ResetMiningMap_Resp, resp)

	// 同步地图数据
	sess.SyncMineMapData()
}

func DigReq(sess *sessions.Session, msgBody []byte) {
	req := structs.DigReq{}
	sess.UnMarshal(msgBody, &req)

	resp := &structs.DigResp{
		Ret: structs.AdventureRet_Failed,
	}

	/////////////////////////////////////////////Data Check////////////////////////////////////////

}
public void DigRequest(NetState netstate, byte x, byte y, bool isOnlyExpandSight = false)
{
	var player = (Player)netstate.Player;
	int nodeId = -1;
	BlockNode node = IsExistBlock(player, x, y, out nodeId);
	
	//校验当前开采点是否有效，规则是非空地块且与其相邻的地块中至少有一个是空地块
	int needDigNum = 0;
	DigResult result = DigResult.Normal;
	if (isOnlyExpandSight)
	{
		result = DigResult.OnlyExpandSight;
	}
	else
	{
		result = CheckDigNode(player, nodeId, x, y, out needDigNum);
	}
	if (result == DigResult.Normal || result == DigResult.OnlyExpandSight)
	{
		//累计挖掘点数
		TileTemplate tileTemplate = Templates.GetTileTemplate(nodeId);

		//检测挖矿成就(力量点)
		GameController.Achievement.CheckAchievements(player, ConditionTypes.CollectPoint, (int)CollectPointTypes.Strength, tileTemplate.AddDigPoint);

		//扣除矿镐耐久度
		GameController.MiningMap.AddMiningPickNum(player, -needDigNum, ResouceChangeType.挖矿_挖掘);

		List<MapTileNode> waitingAddNodeList = new List<MapTileNode>();
		List<MapTileNode> waitingVisibleNodeList = new List<MapTileNode>();

		//对于无资源的上一次开采的地块，直接设置为走过状态
		SetNodeToDigged(player, player.MiningMap.DiggingNode.NodeId, player.MiningMap.DiggingNode.X, player.MiningMap.DiggingNode.Y);

		//设置正在开采地块
		SetCurrentDig(player, nodeId, x, y);
		int expandSight = GlobalTemplate.MiningMap.DefaultSight + player.MiningMap.ExpandSight;

		//更新开采深度
		int tempDepth = Math.Min(y + expandSight, GlobalTemplate.MiningMap.MiningMapHeight);
		if (tempDepth > player.MiningMap.MaxMiningDepth)
		{
			player.MiningMap.MaxMiningDepth = tempDepth;
		}

		//检测是否有超时的boss，有则回收
		CheckOverTimeBoss(player);

		//根据该点，产生周围可视点集合
		for (int i = -expandSight; i <= expandSight; i++)
		{
			for (int j = -expandSight; j <= expandSight; j++)
			{
				if (Math.Abs(i) + Math.Abs(j) <= expandSight)// 
				{
					int intX = i + x;
					int intY = j + y;
					byte newX = (byte)intX;
					byte newY = (byte)intY;
					if (intX >= 0 && intX < GlobalTemplate.MiningMap.MiningMapWidth && intY >= 0 && intY < GlobalTemplate.MiningMap.MiningMapHeight)
					{
						//检测新的点是否已经存在于可视资源/地块队列
						int curNodeId = -1;
						BlockNode curNode = IsExistBlock(player, newX, newY, out curNodeId);
						if (curNode == null)
						{
							//为当前地块生成一个资源或地砖块
							bool isRes = false;
							int layerId = 0;
							curNodeId = RandomResOrBlock(player, newX, newY, out isRes, out layerId);
							if (curNodeId != -1)
							{
								waitingAddNodeList.Add(new MapTileNode() { NodeId = (byte)curNodeId, X = newX, Y = newY, LayerId = (byte)layerId, IsRes = isRes, IsVisible = true, IsAdd = true });
							}
						}
						else
						{
							//已存在，但需要显示的地块
							if (!curNode.IsVisible && curNodeId != -1)
							{
								bool isRes = true;
								if (curNodeId >= 80 && curNodeId < 100)
								{
									isRes = false;
								}
								waitingVisibleNodeList.Add(new MapTileNode() { NodeId = (byte)curNodeId, X = curNode.X, Y = curNode.Y, IsRes = isRes, IsVisible = true, IsAdd = false });
							}
						}
					}
				}
			}
		}

		//检测是否有资源产生，有则判断是否有连续资源出现
		//  最终是已线的方式连续随机，而不是一块的连续随机
		int len = waitingAddNodeList.Count;
		for (int i = 0; i < len; i++)
		{
			MapTileNode tileNode = waitingAddNodeList[i];
			if (tileNode.IsRes)
			{
				//取得该资源所在地层
				MapLayerTemplate layer = Templates.GetMapLayerTemplate(tileNode.LayerId);//mapLayerTemplate[node.LayerId];

				//随机连续资源数量
				int index = layer.TileTypeList.IndexOf(tileNode.NodeId);
				int topCount = layer.TileContinuationCountTop[index];
				int count = rnd.Next(topCount);
				if (count > 0)
				{
					//以当前资源位置为起点，向两侧或上下同化
					bool isHorizontal = false;
					int direction = 1;
					if (rnd.Next(100) > 50)
					{
						isHorizontal = true;
					}
					if (rnd.Next(100) > 50)
					{
						direction = -1;
					}
					for (int j = 1; j <= count; j++)
					{
						byte newX = 0;
						byte newY = 0;
						int delta = rnd.Next(100) > 50 ? 1 : -1;
						if (isHorizontal)
						{
							//水平方向同化
							newX = (byte)(node.X + j * direction);
							newY = (byte)(node.Y + delta);
						}
						else
						{
							//垂直方向同化
							newX = (byte)(node.X + delta);
							newY = (byte)(node.Y + j * direction);
						}
						if (newX >= 0 && newX < GlobalTemplate.MiningMap.MiningMapWidth && newY >= 0 && newY < GlobalTemplate.MiningMap.MiningMapHeight)
						{
							//检测新的点是否已经存在于可视资源/地块队列
							int curNodeId = -1;
							BlockNode curNode = IsExistBlock(player, newX, newY, out curNodeId);
							if (curNode == null)
							{
								MapTileNode oldNode = waitingAddNodeList.FirstOrDefault(o => o.X == newX && o.Y == newY);
								if (oldNode == null)
								{
									oldNode = new MapTileNode { NodeId = tileNode.NodeId, X = newX, Y = newY, IsRes = true, IsVisible = Math.Abs(newX) + Math.Abs(newY) <= expandSight, IsAdd = true };
									waitingAddNodeList.Add(oldNode);
								}
							}
						}
					}
				}
			}
		}

		//将处理过的待添加资源/地块加入到队列
		foreach (var tileNode in waitingAddNodeList)
		{
			AddBlock(player, tileNode.NodeId, tileNode.X, tileNode.Y, tileNode.IsVisible);
		}
		//将待可视的地块变为可视状态
		foreach(var tileNode in waitingVisibleNodeList)
		{
			UpdateBlock(player, tileNode.NodeId, tileNode.X, tileNode.Y, true);
		}
		//合并待增加及待更新的地块列表
		waitingAddNodeList.AddRange(waitingVisibleNodeList);
		//同步客户端，添加可视地块
		ClientProxy.MiningMap.DigResult(netstate, result, waitingAddNodeList.ToArray(), player.MiningMap.DiggingNode, player.MiningMap.MaxMiningDepth);

		DataWriter.MiningMap.Digmap(player, x, y, nodeId);
	}
	else
	{
		//通知客户端开采申请结果
		ClientProxy.MiningMap.DigResult(netstate, result, new MapTileNode[0], player.MiningMap.DiggingNode, player.MiningMap.MaxMiningDepth);
	}
}
