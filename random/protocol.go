package random

import "github.com/TheThingsNetwork/ttn/api/protocol"

func ProtocolTxConfiguration() *protocol.TxConfiguration {
	return LorawanProtocolTxConfiguration()
}

func ProtocolRxMetadata() *protocol.RxMetadata {
	return LorawanProtocolRxMetadata()
}
