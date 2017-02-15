package random

import (
	"math/rand"

	"github.com/TheThingsNetwork/ttn/api/gateway"
	"github.com/TheThingsNetwork/ttn/api/router"
)

// Location returns randomly generated gateway location
func Location() (gps *gateway.GPSMetadata) {
	return &gateway.GPSMetadata{
		Longitude: rand.Float32(),
		Latitude:  rand.Float32(),
		Altitude:  rand.Int31(),
	}
}

func GatewayRxMetadata() *gateway.RxMetadata {
	return &gateway.RxMetadata{
		GatewayId:      ID(),
		Timestamp:      rand.Uint32(),
		Time:           rand.Int63(),
		RfChain:        rand.Uint32(),
		Channel:        rand.Uint32(),
		Frequency:      uint64(rand.Uint32()),
		Rssi:           rand.Float32(),
		Snr:            rand.Float32(),
		Gps:            Location(),
		GatewayTrusted: Bool(),
	}
}

func GatewayTxConfiguration() *gateway.TxConfiguration {
	return &gateway.TxConfiguration{
		Timestamp:             rand.Uint32(),
		RfChain:               rand.Uint32(),
		Frequency:             uint64(rand.Uint32()),
		Power:                 rand.Int31(),
		FrequencyDeviation:    rand.Uint32(),
		PolarizationInversion: Bool(),
	}
}

// GatewayStatus returns randomly generated gateway status
func GatewayStatus() (st *gateway.Status) {
	return &gateway.Status{
		Gps:            Location(),
		Timestamp:      rand.Uint32(),
		Time:           rand.Int63(),
		Ip:             []string{"test"},
		Platform:       "test",
		ContactEmail:   "test@test.test",
		Description:    "test",
		Region:         "test",
		Bridge:         "test",
		Router:         "test",
		Rtt:            rand.Uint32(),
		RxIn:           rand.Uint32(),
		RxOk:           rand.Uint32(),
		TxIn:           rand.Uint32(),
		TxOk:           rand.Uint32(),
		GatewayTrusted: Bool(),
	}
}

// GatewayUplink returns randomly generated gateway uplink
func GatewayUplink() *router.UplinkMessage {
	return LorawanGatewayUplink()
}

// GatewayDownlink returns randomly generated gateway downlink
func GatewayDownlink() *router.DownlinkMessage {
	return LorawanGatewayDownlink()
}
