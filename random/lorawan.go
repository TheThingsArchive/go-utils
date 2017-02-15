package random

import (
	"math/rand"

	"github.com/TheThingsNetwork/ttn/api/protocol"
	"github.com/TheThingsNetwork/ttn/api/protocol/lorawan"
	"github.com/TheThingsNetwork/ttn/api/router"
	brocaar "github.com/brocaar/lorawan"
)

func LorawanPayload(mType ...lorawan.MType) []byte {
	payload := &brocaar.PHYPayload{}
	payload.MHDR.Major = brocaar.LoRaWANR1

	if len(mType) > 0 {
		switch mType[0] {
		case lorawan.MType_JOIN_REQUEST:
			payload.MHDR.MType = brocaar.JoinRequest
		case lorawan.MType_JOIN_ACCEPT:
			payload.MHDR.MType = brocaar.JoinAccept
		case lorawan.MType_UNCONFIRMED_UP:
			payload.MHDR.MType = brocaar.UnconfirmedDataUp
		case lorawan.MType_UNCONFIRMED_DOWN:
			payload.MHDR.MType = brocaar.UnconfirmedDataDown
		case lorawan.MType_CONFIRMED_UP:
			payload.MHDR.MType = brocaar.ConfirmedDataUp
		case lorawan.MType_CONFIRMED_DOWN:
			payload.MHDR.MType = brocaar.ConfirmedDataDown
		}
	} else {
		switch rand.Intn(6) {
		case 0:
			payload.MHDR.MType = brocaar.JoinRequest
		case 1:
			payload.MHDR.MType = brocaar.JoinAccept
		case 2:
			payload.MHDR.MType = brocaar.UnconfirmedDataUp
		case 3:
			payload.MHDR.MType = brocaar.UnconfirmedDataDown
		case 4:
			payload.MHDR.MType = brocaar.ConfirmedDataUp
		case 5:
			payload.MHDR.MType = brocaar.ConfirmedDataDown
		}
	}

	switch payload.MHDR.MType {
	case brocaar.JoinRequest:
		payload.MACPayload = &brocaar.JoinRequestPayload{
			AppEUI:   brocaarAppEUI(),
			DevEUI:   brocaarDevEUI(),
			DevNonce: brocaarDevNonce(),
		}
	case brocaar.JoinAccept:
		payload.MACPayload = &brocaar.JoinAcceptPayload{
			AppNonce: brocaarAppNonce(),
			NetID:    brocaarNetID(),
			DevAddr:  brocaarDevAddr(),
			RXDelay:  uint8(rand.Intn(15)),
			DLSettings: brocaar.DLSettings{
				RX1DROffset: uint8(rand.Intn(7)),
				RX2DataRate: uint8(rand.Intn(15)),
			},
		}
	default:
		payload.MACPayload = &brocaar.MACPayload{
			FHDR: brocaar.FHDR{
				DevAddr: brocaarDevAddr(),
				FCtrl: brocaar.FCtrl{
					ADR:       Bool(),
					ADRACKReq: Bool(),
					ACK:       Bool(),
					FPending:  Bool(),
				},
				FCnt: rand.Uint32(),
			},
		}
	}

	b, err := payload.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return b
}

func LorawanUplinkPayload() []byte {
	switch rand.Intn(3) {
	case 0:
		return LorawanPayload(lorawan.MType_JOIN_REQUEST)
	case 1:
		return LorawanPayload(lorawan.MType_UNCONFIRMED_UP)
	default:
		return LorawanPayload(lorawan.MType_CONFIRMED_UP)
	}
}

func LorawanDownlinkPayload() []byte {
	switch rand.Intn(3) {
	case 0:
		return LorawanPayload(lorawan.MType_JOIN_ACCEPT)
	case 1:
		return LorawanPayload(lorawan.MType_UNCONFIRMED_DOWN)
	default:
		return LorawanPayload(lorawan.MType_CONFIRMED_DOWN)
	}
}

func LorawanProtocolRxMetadata(modulation ...lorawan.Modulation) *protocol.RxMetadata {
	return &protocol.RxMetadata{
		Protocol: &protocol.RxMetadata_Lorawan{
			Lorawan: LorawanRxMetadata(modulation...),
		},
	}
}
func LorawanRxMetadata(modulation ...lorawan.Modulation) *lorawan.Metadata {
	md := &lorawan.Metadata{
		FCnt:       rand.Uint32(),
		CodingRate: CodingRate(),
	}

	if len(modulation) == 1 {
		md.Modulation = modulation[0]
	} else {
		if rand.Int()%2 == 0 {
			md.Modulation = lorawan.Modulation_LORA
		} else {
			md.Modulation = lorawan.Modulation_FSK
		}
	}

	switch md.Modulation {
	case lorawan.Modulation_LORA:
		md.DataRate = DataRate()
	case lorawan.Modulation_FSK:
		md.BitRate = rand.Uint32()
	}
	return md
}

func LorawanProtocolTxConfiguration(modulation ...lorawan.Modulation) *protocol.TxConfiguration {
	return &protocol.TxConfiguration{
		Protocol: &protocol.TxConfiguration_Lorawan{
			Lorawan: LorawanTxConfiguration(modulation...),
		},
	}
}
func LorawanTxConfiguration(modulation ...lorawan.Modulation) *lorawan.TxConfiguration {
	conf := &lorawan.TxConfiguration{
		FCnt:       rand.Uint32(),
		CodingRate: CodingRate(),
	}

	if len(modulation) == 1 {
		conf.Modulation = modulation[0]
	} else {
		if rand.Int()%2 == 0 {
			conf.Modulation = lorawan.Modulation_LORA
		} else {
			conf.Modulation = lorawan.Modulation_FSK
		}
	}

	switch conf.Modulation {
	case lorawan.Modulation_LORA:
		conf.DataRate = DataRate()
	case lorawan.Modulation_FSK:
		conf.BitRate = rand.Uint32()
	}
	return conf
}

// LorawanGatewayUplink returns ly generated lorawan uplink message(join request, confirmed or unconfirmed uplink)
func LorawanGatewayUplink(modulation ...lorawan.Modulation) *router.UplinkMessage {
	switch rand.Intn(3) {
	case 0:
		return LorawanJoinRequest(modulation...)
	case 1:
		return LorawanConfirmedUplink(modulation...)
	default:
		return LorawanUnconfirmedUplink(modulation...)
	}
}

// LorawanJoinRequest returns ly generated lorawan join request
func LorawanJoinRequest(modulation ...lorawan.Modulation) *router.UplinkMessage {
	return &router.UplinkMessage{
		GatewayMetadata:  GatewayRxMetadata(),
		ProtocolMetadata: LorawanProtocolRxMetadata(modulation...),
		Payload:          LorawanPayload(lorawan.MType_JOIN_REQUEST),
	}
}

// LorawanConfirmedUplink returns ly generated confirmed lorawan uplink message
func LorawanConfirmedUplink(modulation ...lorawan.Modulation) *router.UplinkMessage {
	return &router.UplinkMessage{
		GatewayMetadata:  GatewayRxMetadata(),
		ProtocolMetadata: LorawanProtocolRxMetadata(modulation...),
		Payload:          LorawanPayload(lorawan.MType_CONFIRMED_UP),
	}
}

// LorawanUnconfirmedUplink returns ly generated unconfirmed lorawan uplink message
func LorawanUnconfirmedUplink(modulation ...lorawan.Modulation) *router.UplinkMessage {
	return &router.UplinkMessage{
		GatewayMetadata:  GatewayRxMetadata(),
		ProtocolMetadata: LorawanProtocolRxMetadata(modulation...),
		Payload:          LorawanPayload(lorawan.MType_UNCONFIRMED_UP),
	}
}

// LorawanGatewayDownlink returns ly generated lorawan downlink message(join request, confirmed or unconfirmed downlink)
func LorawanGatewayDownlink(modulation ...lorawan.Modulation) *router.DownlinkMessage {
	switch rand.Intn(3) {
	case 0:
		return LorawanJoinAccept(modulation...)
	case 1:
		return LorawanConfirmedDownlink(modulation...)
	default:
		return LorawanUnconfirmedDownlink(modulation...)
	}
}

// LorawanJoinAccept returns ly generated lorawan join request
func LorawanJoinAccept(modulation ...lorawan.Modulation) *router.DownlinkMessage {
	return &router.DownlinkMessage{
		GatewayConfiguration:  GatewayTxConfiguration(),
		ProtocolConfiguration: LorawanProtocolTxConfiguration(modulation...),
		Payload:               LorawanPayload(lorawan.MType_JOIN_ACCEPT),
	}
}

// LorawanConfirmedDownlink returns ly generated confirmed lorawan downlink message
func LorawanConfirmedDownlink(modulation ...lorawan.Modulation) *router.DownlinkMessage {
	return &router.DownlinkMessage{
		GatewayConfiguration:  GatewayTxConfiguration(),
		ProtocolConfiguration: LorawanProtocolTxConfiguration(modulation...),
		Payload:               LorawanPayload(lorawan.MType_CONFIRMED_DOWN),
	}
}

// LorawanUnconfirmedDownlink returns ly generated unconfirmed lorawan downlink message
func LorawanUnconfirmedDownlink(modulation ...lorawan.Modulation) *router.DownlinkMessage {
	return &router.DownlinkMessage{
		GatewayConfiguration:  GatewayTxConfiguration(),
		ProtocolConfiguration: LorawanProtocolTxConfiguration(modulation...),
		Payload:               LorawanPayload(lorawan.MType_UNCONFIRMED_DOWN),
	}
}
