package helpers

import (
	"fmt"
	"os"
	"sync"
	"time"

	"log"

	"github.com/slayercat/GoSNMPServer"
)

type infoProcess struct {
	ListenIP     string
	ListenPort   int
	prossing     *GoSNMPServer.SNMPServer
	MikrotikInfo *MikrotikInfo
}

var CurrentNetWatches []*BaseNetwach

func NewInfoProcess(lIP, mikrotikHost, mikrotikUser, mikrotikPass string, lPort int) *infoProcess {
	return &infoProcess{
		ListenIP:   lIP,
		ListenPort: lPort,
		MikrotikInfo: &MikrotikInfo{
			host:     mikrotikHost,
			username: mikrotikUser,
			password: mikrotikPass,
		},
	}
}

func (i *infoProcess) runServer() {
	master := GoSNMPServer.MasterAgent{
		SubAgents: []*GoSNMPServer.SubAgent{
			{
				CommunityIDs: []string{"public"},
				OIDs:         i.GenerateOIDs(),
			},
		},
	}

	server := GoSNMPServer.NewSNMPServer(master)
	host := fmt.Sprintf("%s:%d", i.ListenIP, i.ListenPort)
	err := server.ListenUDP("udp", host)
	if err != nil {
		log.Println("Error starting server:", err)
		os.Exit(1)
	}
	log.Printf("SNMP server is listening on %s", host)
	i.prossing = server
	if err := server.ServeForever(); err != nil {
		log.Println(err)
	}
	log.Println("SNMP server is shutting down")

}

func isEqualNetWaches(new, old *BaseNetwach) bool {
	return new.name == old.name && new.disabled == old.disabled
}

func diffNetWaches(list []*BaseNetwach) ([]*BaseNetwach, bool) {
	var diff []*BaseNetwach

	if len(CurrentNetWatches) != len(list) {
		return diff, true
	}

	for _, item1 := range CurrentNetWatches {
		found := false
		for _, item2 := range list {
			if isEqualNetWaches(item1, item2) {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, item1)
		}
	}

	if len(diff) != 0 {
		return diff, true
	}

	return nil, false
}
func (i *infoProcess) intervalMonitor(wg *sync.WaitGroup) {
	wg.Add(1)
	log.Println("start interval monitor...")
	go func(i *infoProcess) {
		for {
			time.Sleep(RepeatInterval)
			checkList := i.MikrotikInfo.GetNetwatch()
			if listDiff, ok := diffNetWaches(checkList); ok {
				i.prossing.Shutdown()
				log.Println("reload new data.")
				for _, netwach := range listDiff {
					log.Printf("Changed: Name = %s , disabled = %t", netwach.name, netwach.disabled)
				}
				go i.runServer()
			}
		}
	}(i)
}

func (i *infoProcess) Process() {
	wg := &sync.WaitGroup{}
	go i.runServer()
	i.intervalMonitor(wg)
	wg.Wait()
}
