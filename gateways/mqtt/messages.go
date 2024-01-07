package mqtt

type PacketType int

const (
	ConnectionPacket PacketType = iota
	PublishPacket
	SubscribePacket
	PingPacket
	UnknownPacket
)

const (
	ConnectionRequest = 0x10
	Connack           = 0x20
	ConnackLength     = 0x03
	PublishRequest    = 0x30
	Puback            = 0x40
	SubscribeRequest  = 0x80
	Suback            = 0x90
	PingReq           = 0xC0
	PingResp          = 0xD0
)

type MQTT interface {
	Decode() error
	Process() error
	Reply() ([]byte, error)
}

func RetrievePackageType(buf []byte, size int) PacketType {
	if size <= 1 {
		return UnknownPacket
	}

	req := buf[0]
	if req >= ConnectionRequest && req < Connack {
		return ConnectionPacket
	} else if req >= PublishRequest && req < Puback {
		return PublishPacket
	} else if req >= SubscribeRequest && req < Suback {
		return SubscribePacket
	} else if req >= PingReq && req < PingResp {
		return PingPacket
	}

	return UnknownPacket
}
