package helpers

import (
	"fmt"

	"github.com/gosnmp/gosnmp"
	"github.com/slayercat/GoSNMPServer"
	"github.com/slayercat/GoSNMPServer/mibImps/dismanEventMib"
	"github.com/slayercat/GoSNMPServer/mibImps/ucdMib"
)

var BasicOID = []*GoSNMPServer.PDUValueControlItem{
	{
		OID:      "1.3.6.1.2.1.1.1.0",
		Type:     gosnmp.OctetString,
		OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1OctetStringWrap("RouterOS CHR"), nil },
		Document: "RouterOS CHR",
	},
	{
		OID:  "1.3.6.1.2.1.1.2.0",
		Type: gosnmp.ObjectIdentifier,
		OnGet: func() (value interface{}, err error) {
			return GoSNMPServer.Asn1ObjectIdentifierWrap("1.3.6.1.4.1.14988.1"), nil
		},
		Document: "OID",
	},
	{
		OID:      "1.3.6.1.2.1.1.4.0",
		Type:     gosnmp.OctetString,
		OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1OctetStringWrap("chr"), nil },
		Document: "chr",
	},
	{
		OID:      "1.3.6.1.2.1.1.5.0",
		Type:     gosnmp.OctetString,
		OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1OctetStringWrap("MikroTik"), nil },
		Document: "MikroTik",
	},
	{
		OID:      "1.3.6.1.2.1.1.6.0",
		Type:     gosnmp.OctetString,
		OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1OctetStringWrap(""), nil },
		Document: "",
	},
	{
		OID:      "1.3.6.1.2.1.1.7.0",
		Type:     gosnmp.Integer,
		OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(78), nil },
		Document: "",
	},
}

func (InfoProc *infoProcess) GetStatusInterface(index int) int {
	InfoProc.MikrotikInfo.GetNetwatch()

	switch InfoProc.MikrotikInfo.GetNetwatch()[index].status {
	case "up":
		return 1
	case "down":
		return 2
	default:
		return 4
	}
}

// generate oids
func (InfoProc *infoProcess) GenerateOIDs() []*GoSNMPServer.PDUValueControlItem {
	customOIDs := []*GoSNMPServer.PDUValueControlItem{}
	customOIDs = append(customOIDs, BasicOID...)

	netwaches := InfoProc.MikrotikInfo.GetNetwatch()
	CurrentNetWatches = netwaches

	// Get uptime (hardware)
	upt := dismanEventMib.All()
	customOIDs = append(customOIDs, upt...)

	if netwaches != nil {

		objNumber := []*GoSNMPServer.PDUValueControlItem{
			{
				OID:  "1.3.6.1.2.1.2.1.0",
				Type: gosnmp.Integer,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1IntegerWrap((len(netwaches))), nil
				},
				Document: "netwatchIndex",
			},
			{
				OID:  "1.3.6.1.2.1.55.1.3.0",
				Type: gosnmp.Integer,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1IntegerWrap((len(netwaches))), nil
				},
				Document: "netwatchIndex",
			},
		}
		netWatches := []*GoSNMPServer.PDUValueControlItem{}
		for id, item := range netwaches {
			cid := id + 1
			ifName := item.name
			ifOperStatus, ifHighSpeed := func() (uint, uint) {
				switch item.status {
				case "up":
					return uint(1000000000), 1000
				case "down":
					return uint(0), 0
				default:
					return uint(0), 0
				}
			}()
			thisNetwatchID := []*GoSNMPServer.PDUValueControlItem{
				{
					OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.1.%d", cid),
					Type:     gosnmp.Integer,
					OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(cid), nil },
					Document: "netwatchIndex",
				},
				{
					OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.2.%d", cid),
					Type:     gosnmp.OctetString,
					OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1OctetStringWrap(ifName), nil },
					Document: "IFDESCR",
				},
				{
					OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.3.%d", cid),
					Type:     gosnmp.Integer,
					OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(6), nil },
					Document: "IANAifType",
				},
				{
					OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.4.%d", cid),
					Type:     gosnmp.Integer,
					OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1500), nil },
					Document: "IFTYPE",
				},
				{
					OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.5.%d", cid),
					Type:     gosnmp.Gauge32,
					OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Gauge32Wrap(ifOperStatus), nil },
					Document: "IFOPERSTATUS",
				},
				{
					OID:      fmt.Sprintf("1.3.6.1.2.1.31.1.1.1.6.%d", cid),
					Type:     gosnmp.Counter64,
					OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter64Wrap(0), nil },
					Document: "ifHCInOctets",
				},
				{
					OID:      fmt.Sprintf("1.3.6.1.2.1.31.1.1.1.10.%d", cid),
					Type:     gosnmp.Counter64,
					OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter64Wrap(15736), nil },
					Document: "ifHCOutOctets",
				},
				{
					OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.13.%d", cid),
					Type:     gosnmp.Counter32,
					OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
					Document: "ifHCOutOctets",
				},
				{
					OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.14.%d", cid),
					Type:     gosnmp.Counter32,
					OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
					Document: "ifInErrors",
				},
				{
					OID:      fmt.Sprintf("1.3.6.1.2.1.31.1.1.1.15.%d", cid),
					Type:     gosnmp.Counter32,
					OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(ifHighSpeed), nil },
					Document: "ifHighSpeed",
				},
				{
					OID:      fmt.Sprintf("1.3.6.1.2.1.31.1.1.1.18.%d", cid),
					Type:     gosnmp.OctetString,
					OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1OctetStringWrap(""), nil },
					Document: "IFALIAS",
				},
				{
					OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.19.%d", cid),
					Type:     gosnmp.Counter32,
					OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
					Document: "ifOutDiscards",
				},
				{
					OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.20.%d", cid),
					Type:     gosnmp.Counter32,
					OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
					Document: "ifOutErrors",
				},
				{
					OID:      fmt.Sprintf("1.3.6.1.2.1.31.1.1.1.1.%d", cid),
					Type:     gosnmp.OctetString,
					OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1OctetStringWrap(ifName), nil },
					Document: "IFDESCR",
				},
				{
					OID:  fmt.Sprintf("1.3.6.1.2.1.2.2.1.7.%d", cid),
					Type: gosnmp.Integer,
					OnGet: func() (value interface{}, err error) {
						return GoSNMPServer.Asn1IntegerWrap(InfoProc.GetStatusInterface(id)), nil
					},
					Document: "IFADMINSTATUS",
				},
				{
					OID:  fmt.Sprintf("1.3.6.1.2.1.2.2.1.8.%d", cid),
					Type: gosnmp.Integer,
					OnGet: func() (value interface{}, err error) {
						return GoSNMPServer.Asn1IntegerWrap(InfoProc.GetStatusInterface(id)), nil
					},
					Document: "IFOPERSTATUS",
				},
				{
					OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.18.%d", cid),
					Type:     gosnmp.OctetString,
					OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1OctetStringWrap(""), nil },
					Document: "IFALIAS",
				},
				{
					OID:      fmt.Sprintf("1.3.6.1.2.1.55.1.5.1.2.%d", cid),
					Type:     gosnmp.OctetString,
					OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1OctetStringWrap(ifName), nil },
					Document: "IFDESCR",
				},
				{
					OID:      fmt.Sprintf("1.3.6.1.2.1.55.1.5.1.4.%d", cid),
					Type:     gosnmp.Gauge32,
					OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Gauge32Wrap(1500), nil },
					Document: "IFALIAS",
				},
				{
					OID:  fmt.Sprintf("1.3.6.1.2.1.55.1.5.1.10.%d", cid),
					Type: gosnmp.Integer,
					OnGet: func() (value interface{}, err error) {
						return GoSNMPServer.Asn1IntegerWrap(InfoProc.GetStatusInterface(id)), nil
					},
					Document: "IFALIAS",
				},
				{
					OID:  fmt.Sprintf("1.3.6.1.2.1.55.1.5.1.9.%d", cid),
					Type: gosnmp.Integer,
					OnGet: func() (value interface{}, err error) {
						return GoSNMPServer.Asn1IntegerWrap(InfoProc.GetStatusInterface(id)), nil
					},
					Document: "DODINTERNET",
				},
			}
			netWatches = append(netWatches, thisNetwatchID...)
		}

		customOIDs = append(customOIDs, objNumber...)
		customOIDs = append(customOIDs, netWatches...)

	}

	customOIDs = append(customOIDs, ucdMib.DiskUsageOIDs()...)
	customOIDs = append(customOIDs, []*GoSNMPServer.PDUValueControlItem{
		{
			OID:  "1.3.6.1.2.1.9999.1.1.1.1.0",
			Type: gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1OctetStringWrap("RouterOS DHCP server"), nil
			},
			Document: "RouterOS DHCP server",
		}, {
			OID:  "1.3.6.1.2.1.9999.1.1.1.2.0",
			Type: gosnmp.ObjectIdentifier,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1ObjectIdentifierWrap("1.3.6.1.4.1.14988.1"), nil
			},
			Document: "End of SNMP",
		}}...)

	return customOIDs
}
