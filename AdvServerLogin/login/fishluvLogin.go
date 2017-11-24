package login

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type LoginRet struct {
	Error             string     `xml:"Error"`
	LoginTakon        string     `xml:"LoginTakon"`
	Notice            string     `xml:"Notice"`
	Version           string     `xml:"Version"`
	IsShowEveryTime   bool       `xml:"IsShowEveryTime"`
	NickName          string     `xml:"NickName"`
	PlatformAccountId string     `xml:"PlatformAccountId"`
	GameServers       AllServers `xml:"GameServers"`
}

type ServerStatus int

const (
	SERVER_STATUS_NORMAL ServerStatus = 1
)

type AllServers struct {
	Server []GameServer `xml:"GameServer"`
}
type GameServer struct {
	XMLName       xml.Name     `xml:"GameServer"`
	GameZoneId    int          `xml:"GameZoneId"`
	Name          string       `xml:"Name"`
	Status        ServerStatus `xml:"Status"`
	Host          string       `xml:"Host"`
	Port          int          `xml:"Port"`
	Recommend     bool         `xml:"Recommend"`
	CharacterName string       `xml:"CharacterName"`
}

func FishluvLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       //解析参数，默认是不会解析的
	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息

	UID := r.Form.Get("uid")
	phoneID := r.Form.Get("phoneid")

	fmt.Println("czx@@@ login", UID, phoneID)

	if UID == "" || phoneID == "" {
		fmt.Println("error openid")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	gsInfo1 := GameServer{
		GameZoneId:    1,
		Name:          "测试11",
		Status:        SERVER_STATUS_NORMAL,
		Host:          "127.0.0.1",
		Port:          9110,
		Recommend:     false,
		CharacterName: "角色1",
	}
	gsInfo2 := GameServer{
		GameZoneId:    2,
		Name:          "测试2",
		Status:        SERVER_STATUS_NORMAL,
		Host:          "127.0.0.1",
		Port:          9111,
		Recommend:     false,
		CharacterName: "角色2",
	}

	ret := LoginRet{
		LoginTakon:        "takon1",
		Notice:            "notice1",
		Version:           "version1",
		IsShowEveryTime:   false,
		NickName:          "nick1",
		PlatformAccountId: "platform1",
		GameServers: AllServers{
			Server: []GameServer{gsInfo1, gsInfo2},
		},
	}

	xmlret, err := xml.Marshal(ret)
	if err != nil {
		fmt.Println("error server info")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	fmt.Fprintf(w, string(xmlret))
}
