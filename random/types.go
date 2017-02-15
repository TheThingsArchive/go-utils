package random

import (
	"fmt"
	"math/rand"

	"github.com/TheThingsNetwork/ttn/core/types"
	brocaar "github.com/brocaar/lorawan"
)

// Seed is a wrapper around rand.Seed
func Seed(seed int64) {
	rand.Seed(seed)
}

const idChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// ID returns randomly generated ID
func ID(n ...int) string {
	var nVal int
	if len(n) > 0 {
		nVal = n[0]
	} else {
		nVal = 2 + rand.Intn(61)
	}
	b := make([]byte, nVal)
	for i := range b {
		b[i] = idChars[rand.Intn(len(idChars))]
	}
	return string(b)
}

func Bool() bool {
	if rand.Int31()%2 == 0 {
		return true
	}
	return false
}

func Byte() byte {
	return byte(rand.Intn(255))
}

func random2byteArray() [2]byte {
	var bytes [2]byte
	for i := range bytes {
		bytes[i] = Byte()
	}
	return bytes
}
func random3byteArray() [3]byte {
	var bytes [3]byte
	for i := range bytes {
		bytes[i] = Byte()
	}
	return bytes
}
func random4byteArray() [4]byte {
	var bytes [4]byte
	for i := range bytes {
		bytes[i] = Byte()
	}
	return bytes
}
func random8byteArray() [8]byte {
	var bytes [8]byte
	for i := range bytes {
		bytes[i] = Byte()
	}
	return bytes
}

func ByteSlice(n int) []byte {
	bytes := make([]byte, n)
	for i := range bytes {
		bytes[i] = Byte()
	}
	return bytes
}

func brocaarDevNonce() [2]byte {
	return random2byteArray()
}
func brocaarAppNonce() [3]byte {
	return random3byteArray()
}
func brocaarNetID() brocaar.NetID {
	return brocaar.NetID(random3byteArray())
}
func brocaarDevAddr() brocaar.DevAddr {
	return brocaar.DevAddr(random4byteArray())
}
func brocaarEUI64() brocaar.EUI64 {
	return brocaar.EUI64(random8byteArray())
}
func brocaarDevEUI() brocaar.EUI64 {
	return brocaarEUI64()
}
func brocaarAppEUI() brocaar.EUI64 {
	return brocaarEUI64()
}

func DevNonce() types.DevNonce {
	return types.DevNonce(brocaarDevNonce())
}
func AppNonce() types.AppNonce {
	return types.AppNonce(brocaarAppNonce())
}
func NetID() types.NetID {
	return types.NetID(brocaarNetID())
}
func DevAddr() types.DevAddr {
	return types.DevAddr(brocaarDevAddr())
}
func EUI64() types.EUI64 {
	return types.EUI64(brocaarEUI64())
}
func DevEUI() types.DevEUI {
	return types.DevEUI(brocaarDevEUI())
}
func AppEUI() types.AppEUI {
	return types.AppEUI(brocaarAppEUI())
}

func DataRate() string {
	return fmt.Sprintf("SF%dBW%d", 7+rand.Intn(5), 125*(1+rand.Intn(1)*(1+rand.Intn(1))))
}

func CodingRate() string {
	first := 1 + rand.Intn(7)
	second := first + rand.Intn(9-first)
	return fmt.Sprintf("%d/%d", first, second)
}
