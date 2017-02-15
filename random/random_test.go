package random

import (
	"testing"

	"github.com/TheThingsNetwork/ttn/api"
	"github.com/TheThingsNetwork/ttn/api/protocol/lorawan"
	s "github.com/smartystreets/assertions"
)

type PayloadUnmarshaller interface {
	UnmarshalPayload() error
}

func TestRandomizers(t *testing.T) {
	for name, msg := range map[string]interface{}{
		"Location": Location(),

		"GatewayStatus": GatewayStatus(),

		"LorawanProtocolRxMetadata(LORA)": LorawanProtocolRxMetadata(lorawan.Modulation_LORA),
		"LorawanRxMetadata(LORA)":         LorawanRxMetadata(lorawan.Modulation_LORA),
		"LorawanProtocolRxMetadata(FSK)":  LorawanProtocolRxMetadata(lorawan.Modulation_FSK),
		"LorawanRxMetadata(FSK)":          LorawanRxMetadata(lorawan.Modulation_FSK),

		"LorawanConfirmedUplink(LORA)":   LorawanConfirmedUplink(lorawan.Modulation_LORA),
		"LorawanUnconfirmedUplink(LORA)": LorawanUnconfirmedUplink(lorawan.Modulation_LORA),
		"LorawanJoinRequest(LORA)":       LorawanJoinRequest(lorawan.Modulation_LORA),
		"LorawanConfirmedUplink(FSK)":    LorawanConfirmedUplink(lorawan.Modulation_FSK),
		"LorawanUnconfirmedUplink(FSK)":  LorawanUnconfirmedUplink(lorawan.Modulation_FSK),
		"LorawanJoinRequest(FSK)":        LorawanJoinRequest(lorawan.Modulation_FSK),

		"LorawanGatewayUplink(LORA)": LorawanGatewayUplink(lorawan.Modulation_LORA),
		"LorawanGatewayUplink(FSK)":  LorawanGatewayUplink(lorawan.Modulation_FSK),

		"ProtocolRxMetadata()": ProtocolRxMetadata(),
		"GatewayRxMetadata()":  GatewayRxMetadata(),
		"GatewayUplink()":      GatewayUplink(),

		"LorawanProtocolTxConfiguration(LORA)": LorawanProtocolTxConfiguration(lorawan.Modulation_LORA),
		"LorawanTxConfiguration(LORA)":         LorawanTxConfiguration(lorawan.Modulation_LORA),
		"LorawanProtocolTxConfiguration(FSK)":  LorawanProtocolTxConfiguration(lorawan.Modulation_FSK),
		"LorawanTxConfiguration(FSK)":          LorawanTxConfiguration(lorawan.Modulation_FSK),

		"LorawanConfirmedDownlink(LORA)":   LorawanConfirmedDownlink(lorawan.Modulation_LORA),
		"LorawanUnconfirmedDownlink(LORA)": LorawanUnconfirmedDownlink(lorawan.Modulation_LORA),
		"LorawanJoinAccept(LORA)":          LorawanJoinAccept(lorawan.Modulation_LORA),
		"LorawanConfirmedDownlink(FSK)":    LorawanConfirmedDownlink(lorawan.Modulation_FSK),
		"LorawanUnconfirmedDownlink(FSK)":  LorawanUnconfirmedDownlink(lorawan.Modulation_FSK),
		"LorawanJoinAccept(FSK)":           LorawanJoinAccept(lorawan.Modulation_FSK),

		"LorawanGatewayDownlink(LORA)": LorawanGatewayDownlink(lorawan.Modulation_LORA),
		"LorawanGatewayDownlink(FSK)":  LorawanGatewayDownlink(lorawan.Modulation_FSK),

		"ProtocolTxConfiguration()": ProtocolTxConfiguration(),
		"GatewayTxConfiguration()":  GatewayTxConfiguration(),
		"GatewayDownlink()":         GatewayDownlink(),
	} {
		t.Run(name, func(t *testing.T) {
			if v, ok := msg.(api.Validator); ok {
				t.Run("Validate", func(t *testing.T) {
					a := s.New(t)
					a.So(v.Validate(), s.ShouldBeNil)
				})
			}
			if v, ok := msg.(PayloadUnmarshaller); ok {
				t.Run("UnmarshalPayload", func(t *testing.T) {
					a := s.New(t)
					a.So(v.UnmarshalPayload(), s.ShouldBeNil)
				})
			}
		})
	}
}
