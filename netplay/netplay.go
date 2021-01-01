package netplay

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"time"

	"github.com/libretro/ludo/input"
)

const NET_INPUT_DELAY = 3
const NET_INPUT_HISTORY_SIZE = 60
const NET_SEND_DELAY_FRAMES = 0
const NET_SEND_HISTORY_SIZE = 5

// Network code indicating the type of message.
const (
	MsgCodeHandshake   = 1 // Used when sending the hand shake.
	MsgCodePlayerInput = 2 // Sends part of the player's input buffer.
	MsgCodePing        = 3 // Used to tracking packet round trip time. Expect a "Pong" back.
	MsgCodePong        = 4 // Sent in reply to a Ping message for testing round trip time.
	MsgCodeSync        = 5 // Used to pass sync data
)

// Listen is used by the netplay host, listening address and port
var Listen bool

// Join is used by the netplay guest, address of the host
var Join bool

// Conn is the connection between two players
var Conn *net.UDPConn

var Enabled = false
var ConnectedToClient = false
var isServer = false
var confirmedTick = -1
var localSyncData = ""
var remoteSyncData = ""
var isStateDesynced = false
var localSyncDataTick = -1
var remoteSyncDataTick = -1
var LocalTickDelta = 0
var RemoteTickDelta = 0
var syncDataHistoryLocal = []string{}
var syncDataHistoryRemote = []string{}
var inputHistory = [][20]byte{}
var remoteInputHistory = [][20]byte{}
var toSendPackets = []struct {
	Packet []byte
	Time   time.Time
}{}
var clientIPPort net.Addr
var latency int64
var ConfirmedTick = 0
var TickSyncing = false
var TickOffset = 0
var LastSyncedTick = -1
var DesyncCheckRate = 20

// Init initialises a netplay session between two players
func Init() {
	if Listen { // Host mode
		Conn, err := net.ListenUDP("udp", &net.UDPAddr{
			Port: 8080,
		})
		if err != nil {
			log.Println("Netplay", err.Error())
			return
		}

		Conn.SetReadBuffer(1048576)

		log.Println("Netplay", "Listening.")

		msg := [2]byte{}
		Conn.Read(msg[:])
		log.Println(msg)
		log.Println("Netplay", "Player #2 is connected.")
		Enabled = true
		isServer = true
	} else if Join { // Guest mode
		var err error
		Conn, err = net.DialUDP("udp", nil, &net.UDPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 8080,
		})
		if err != nil {
			log.Println("Netplay", err.Error())
			return
		}
		log.Println("Netplay", "Connected.")
		Conn.Write([]byte("hi"))
		Enabled = true
		isServer = false
		clientIPPort = Conn.RemoteAddr()
	}
}

// Get input from the remote player for the passed in game tick.
func GetRemoteInputState(tick int) input.PlayerState {
	if tick > confirmedTick {
		// Repeat the last confirmed input when we don't have a confirmed tick
		tick = confirmedTick
	}
	return DecodeInput(remoteInputHistory[(NET_INPUT_HISTORY_SIZE+tick)%NET_INPUT_HISTORY_SIZE])
}

// Get input state for the local client
func GetLocalInputState(tick int) input.PlayerState {
	return DecodeInput(inputHistory[(NET_INPUT_HISTORY_SIZE+tick)%NET_INPUT_HISTORY_SIZE])
}

func GetLocalInputEncoded(tick int) [20]byte {
	return inputHistory[(NET_INPUT_HISTORY_SIZE+tick)%NET_INPUT_HISTORY_SIZE]
}

// Get the sync data which is used to check for game state desync between the clients.
func GetSyncDataLocal(tick int) string {
	index := (NET_INPUT_HISTORY_SIZE + tick) % NET_INPUT_HISTORY_SIZE
	return syncDataHistoryLocal[index]
}

// Get sync data from the remote client.
func GetSyncDataRemote(tick int) string {
	index := (NET_INPUT_HISTORY_SIZE + tick) % NET_INPUT_HISTORY_SIZE
	return syncDataHistoryRemote[index]
}

// Set sync data for a game tick
func SetLocalSyncData(tick int, syncData string) {
	if !isStateDesynced {
		localSyncData = syncData
		localSyncDataTick = tick
	}
}

// Check for a desync.
func DesyncCheck() (bool, int) {
	if localSyncDataTick < 0 {
		return false, 0
	}

	// When the local sync data does not match the remote data indicate a desync has occurred.
	if isStateDesynced || localSyncDataTick == remoteSyncDataTick {
		// print("Desync Check at: " .. localSyncDataTick)

		if localSyncData != remoteSyncData {
			isStateDesynced = true
			return true, localSyncDataTick
		}
	}

	return false, 0
}

// Send the inputState for the local player to the remote player for the given game tick.
func SendInputData(tick int) {
	// Don't send input data when not connect to another player's game client.
	if !(Enabled && ConnectedToClient) {
		return
	}

	SendPacket(MakeInputPacket(tick), 1)
}

func SetLocalInput(st input.PlayerState, tick int) {
	encodedInput := EncodeInput(st)
	inputHistory[(NET_INPUT_HISTORY_SIZE+tick)%NET_INPUT_HISTORY_SIZE] = encodedInput
}

func SetRemoteEncodedInput(encodedInput [20]byte, tick int) {
	remoteInputHistory[(NET_INPUT_HISTORY_SIZE+tick)%NET_INPUT_HISTORY_SIZE] = encodedInput
}

// Handles sending packets to the other client. Set duplicates to something > 0 to send more than once.
func SendPacket(packet []byte, duplicates int) {
	if duplicates == 0 {
		duplicates = 1
	}

	for i := 1; i < duplicates; i++ {
		if NET_SEND_DELAY_FRAMES > 0 {
			SendPacketWithDelay(packet)
		} else {
			SendPacketRaw(packet)
		}
	}
}

// Queues a packet to be sent later
func SendPacketWithDelay(packet []byte) {
	delayedPacket := struct {
		Packet []byte
		Time   time.Time
	}{
		Packet: packet,
		Time:   time.Now(),
	}
	toSendPackets = append(toSendPackets, delayedPacket)
}

// Send all packets which have been queued and who's delay time as elapsed.
func ProcessDelayedPackets() {
	newPacketList := []struct {
		Packet []byte
		Time   time.Time
	}{} // List of packets that haven't been sent yet.
	timeInterval := NET_SEND_DELAY_FRAMES / 60 // How much time must pass (converting from frames into seconds)

	for _, data := range toSendPackets {
		if (time.Now().Unix() - data.Time.Unix()) > int64(timeInterval) {
			SendPacketRaw(data.Packet) // Send packet when enough time as passed.
		} else {
			newPacketList = append(newPacketList, data) // Keep the packet if the not enough time as passed.
		}
	}
	toSendPackets = newPacketList
}

// Send a packet immediately
func SendPacketRaw(packet []byte) {
	Conn.WriteTo(packet, clientIPPort)
}

// Handles receiving packets from the other client.
func ReceivePacket() ([]byte, net.Addr, error) {
	data := []byte{}

	_, addr, err := Conn.ReadFrom(data)

	return data[:], addr, err
}

// Checks the queue for any incoming packets and process them.
func ReceiveData() {
	if !Enabled {
		return
	}

	// For now we'll process all packets every frame.
	for {
		data, addr, err := ReceivePacket()

		if err == nil {
			r := bytes.NewReader(data)
			var code byte
			binary.Read(r, binary.LittleEndian, &code)

			// Handshake code must be received by both game instances before a match can begin.
			if code == MsgCodeHandshake {
				if !ConnectedToClient {
					ConnectedToClient = true

					// The server needs to remember the address and port in order to send data to the other cilent.
					if true {
						// Server needs to the other the client address and ip to know where to send data.
						if isServer {
							clientIPPort = addr
						}
						log.Println("Received Handshake from: ", clientIPPort.String())
						// Send handshake to client.
						SendPacket(MakeHandshakePacket(), 5)
					}
				}
			} else if code == MsgCodePlayerInput {
				// Break apart the packet into its parts.
				//results := { love.data.unpack(INPUT_FORMAT_STRING, data, 1) } // Final parameter is the start position

				var tickDelta, receivedTick int
				binary.Read(r, binary.LittleEndian, &tickDelta)
				binary.Read(r, binary.LittleEndian, &receivedTick)

				// We only care about the latest tick delta, so make sure the confirmed frame is atleast the same or newer.
				// This would work better if we added a packet count.
				if receivedTick >= confirmedTick {
					RemoteTickDelta = tickDelta
				}

				if receivedTick > confirmedTick {
					if receivedTick-confirmedTick > NET_INPUT_DELAY {
						log.Println("Received packet with a tick too far ahead. Last: ", confirmedTick, "     Current: ", receivedTick)
					}

					confirmedTick = receivedTick

					// log.Println("Received Input: ", results[3+NET_SEND_HISTORY_SIZE], " @ ",  receivedTick)

					for offset := 0; offset < NET_SEND_HISTORY_SIZE-1; offset++ {
						var encodedInput [20]byte
						binary.Read(r, binary.LittleEndian, &encodedInput)
						// Save the input history sent in the packet.
						SetRemoteEncodedInput(encodedInput, receivedTick-offset)
					}
				}

				// NetLog("Received Tick: " .. receivedTick .. ",  Input: " .. remoteInputHistory[(confirmedTick % NET_INPUT_HISTORY_SIZE)+1])
			} else if code == MsgCodePing {
				var pingTime time.Time
				binary.Read(r, binary.LittleEndian, &pingTime)
				SendPacket(MakePongPacket(pingTime), 1)
			} else if code == MsgCodePong {
				var pongTime time.Time
				binary.Read(r, binary.LittleEndian, &pongTime)
				latency = time.Now().Unix() - pongTime.Unix()
				//print("Got pong message: " .. latency)
			} else if code == MsgCodeSync {
				var tick int
				var syncData string
				binary.Read(r, binary.LittleEndian, &tick)
				binary.Read(r, binary.LittleEndian, &syncData)
				// Ignore any tick that isn't more recent than the last sync data
				if !isStateDesynced && tick > remoteSyncDataTick {
					remoteSyncDataTick = tick
					remoteSyncData = syncData

					// Check for a desync
					DesyncCheck()
				}
			}
		} else {
			return
		}
	}
}

// Generate a packet containing information about player input.
func MakeInputPacket(tick int) []byte {
	// log.Println('[Packet] tick: ', tick, '      input: ', history[NET_SEND_HISTORY_SIZE])
	// data := love.data.pack("string", INPUT_FORMAT_STRING, MsgCodePlayerInput, LocalTickDelta, tick, unpack(history))
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, MsgCodePlayerInput)
	binary.Write(buf, binary.LittleEndian, LocalTickDelta)
	binary.Write(buf, binary.LittleEndian, tick)

	historyIndexStart := tick - NET_SEND_HISTORY_SIZE
	for i := 0; i < NET_SEND_HISTORY_SIZE; i++ {
		encodedInput := inputHistory[(NET_INPUT_HISTORY_SIZE+historyIndexStart+i)%NET_INPUT_HISTORY_SIZE]
		binary.Write(buf, binary.LittleEndian, encodedInput)
	}

	return buf.Bytes()
}

// Send a ping message in order to test network latency
func SendPingMessage() {
	SendPacket(MakePingPacket(time.Now()), 1)
}

// Make a ping packet
func MakePingPacket(t time.Time) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, MsgCodePing)
	binary.Write(buf, binary.LittleEndian, t)
	return buf.Bytes()
}

// Make pong packet
func MakePongPacket(t time.Time) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, MsgCodePong)
	binary.Write(buf, binary.LittleEndian, t)
	return buf.Bytes()
}

// Sends sync data
func SendSyncData() {
	SendPacket(MakeSyncDataPacket(localSyncDataTick, localSyncData), 5)
}

// Make a sync data packet
func MakeSyncDataPacket(tick int, syncData string) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, MsgCodeSync)
	binary.Write(buf, binary.LittleEndian, tick)
	binary.Write(buf, binary.LittleEndian, syncData)
	return buf.Bytes()
}

// Generate handshake packet for connecting with another client.
func MakeHandshakePacket() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, MsgCodeHandshake)
	return buf.Bytes()
}

// Encodes the player input state into a compact form for network transmission.
func EncodeInput(st input.PlayerState) [20]byte {
	netoutput := [20]byte{}
	for i, b := range st {
		if b {
			netoutput[i] = 1
		}
	}
	return netoutput
}

// Decodes the input from a packet generated by EncodeInput().
func DecodeInput(data [20]byte) input.PlayerState {
	st := input.PlayerState{}

	for i, b := range data {
		if b == 1 {
			st[i] = true
		}
	}
	return st
}
