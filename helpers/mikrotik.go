package helpers

import (
	"log"
	"strconv"

	"github.com/go-routeros/routeros"
	"github.com/go-routeros/routeros/proto"
)

type MikrotikInfo struct {
	host     string
	username string
	password string
	obj      *routeros.Client
	Response *routeros.Reply
}

type BaseNetwach struct {
	name     string
	status   string
	disabled bool
}

type MikrotikNetwach struct {
	netwaches []*BaseNetwach
}

type IMapParser interface {
	MapStructs([]*proto.Sentence)
}

func (m *MikrotikInfo) MapParser(i IMapParser) {
	i.MapStructs(m.Response.Re)
}

func NewMikrotik(host, username, paswword string) *MikrotikInfo {
	return &MikrotikInfo{
		host:     host,
		username: username,
		password: paswword,
	}
}

func (m *MikrotikInfo) dial() error {
	client, err := routeros.Dial(m.host, m.username, m.password)
	if err != nil {
		log.Println(err)
		return err
	}
	m.obj = client
	return nil
}

func (m *MikrotikNetwach) MapStructs(dataResponse []*proto.Sentence) {
	for _, res := range dataResponse {
		boolVal, err := strconv.ParseBool(res.Map["disabled"])
		if err != nil {
			boolVal = false
		}

		m.netwaches = append(m.netwaches, &BaseNetwach{
			name:     res.Map["name"],
			status:   res.Map["status"],
			disabled: boolVal,
		})
	}
}

func (m *MikrotikInfo) GetNetwatch() []*BaseNetwach {
	if err := m.dial(); err != nil {
		return nil
	}
	defer m.obj.Close()
	res, err := m.obj.Run(NetWatchURL)
	if err != nil {
		log.Println(err)
	}
	m.Response = res
	net := MikrotikNetwach{}
	m.MapParser(&net)
	return net.netwaches
}
