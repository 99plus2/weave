package weave

import (
	"code.google.com/p/gopacket/pcap"
)

type PcapIO struct {
	handle *pcap.Handle
}

func NewPcapIO(ifName string, bufSz int) (pio PacketSourceSink, err error) {
	pio, err = newPcapIO(ifName, true, 65535, bufSz)
	return
}

func NewPcapO(ifName string) (po PacketSink, err error) {
	po, err = newPcapIO(ifName, false, 0, 0)
	return
}

func newPcapIO(ifName string, promisc bool, snaplen int, bufSz int) (handle *PcapIO, err error) {
	inactive, err := pcap.NewInactiveHandle(ifName)
	if err != nil {
		return
	}
	defer inactive.CleanUp()
	if err = inactive.SetPromisc(promisc); err != nil {
		return
	}
	if err = inactive.SetSnapLen(snaplen); err != nil {
		return
	}
	if err = inactive.SetTimeout(-1); err != nil {
		return
	}
	if err = inactive.SetImmediateMode(true); err != nil {
		return
	}
	if err = inactive.SetBufferSize(bufSz); err != nil {
		return
	}
	active, err := inactive.Activate()
	if err != nil {
		return
	}
	if err = active.SetDirection(pcap.DirectionIn); err != nil {
		return
	}
	return &PcapIO{handle: active}, nil
}

func (pi *PcapIO) ReadPacket() (data []byte, err error) {
	data, _, err = pi.handle.ZeroCopyReadPacketData()
	return
}

func (po *PcapIO) WritePacket(data []byte) error {
	return po.handle.WritePacketData(data)
}
