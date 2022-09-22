package uuid

import (
	"fmt"
	"math/rand"
	"net"
	"sync/atomic"
	"time"
)

const (
	// 开始时间 （2020-05-03）
	epoch uint64 = 1588435200000

	workerIDBits = 10

	timestampBits = 41

	sequenceBits = 12

	maxWorkerID = -1 ^ (-1 << workerIDBits)

	// 掩码
	timestampAndSequenceMask uint64 = -1 ^ (-1 << (timestampBits + sequenceBits))
)

// highest 11 bit: workerID
// middle  41 bit: timestamp
// lowest  12 bit: sequence
var timestampAndSequence uint64

var workerID = generateWorkerID() << (timestampBits + sequenceBits)

func init() {
	timestamp := getNewestTimestamp()
	timestampWithSequence := timestamp << sequenceBits
	atomic.StoreUint64(&timestampAndSequence, timestampWithSequence)
}

func NextID() int64 {
	next := atomic.AddUint64(&timestampAndSequence, 1)
	timestampWithSequence := next & timestampAndSequenceMask
	return int64(uint64(workerID) | timestampWithSequence)
}

func NextStringId() string {
	id := NextID()
	return string(id)
}

func getNewestTimestamp() uint64 {
	return uint64(time.Now().UTC().UnixNano())/uint64(time.Millisecond/time.Nanosecond) - epoch
}

func generateWorkerID() int64 {
	id, err := generateWorkerIDBaseOnMac()
	if err != nil {
		id = generateRandomWorkerID()
	}
	return id
}

// use lowest 10 bit of available MAC as workerID
func generateWorkerIDBaseOnMac() (int64, error) {
	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		mac := iface.HardwareAddr

		return int64(int(rune(mac[4]&0b11)<<8) | int(mac[5]&0xFF)), nil
	}
	return 0, fmt.Errorf("no available mac found")
}

// randomly generate one as workerID
func generateRandomWorkerID() int64 {
	return rand.Int63n(maxWorkerID + 1)
}
