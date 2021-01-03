package netplay

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"time"

	"github.com/libretro/ludo/input"
)

const inputDelayFrames = 3
const inputHistorySize = int64(300)
const sendDelayFrames = 0
const sendHistorySize = 5

// Network code indicating the type of message.
const (
	MsgCodeHandshake   = byte(1) // Used when sending the hand shake.
	MsgCodePlayerInput = byte(2) // Sends part of the player's input buffer.
	MsgCodePing        = byte(3) // Used to tracking packet round trip time. Expect a "Pong" back.
	MsgCodePong        = byte(4) // Sent in reply to a Ping message for testing round trip time.
	MsgCodeSync        = byte(5) // Used to pass sync data
)

// Listen is used by the netplay host, listening address and port
var Listen bool

// Join is used by the netplay guest, address of the host
var Join bool

// Conn is the connection between two players
var Conn *net.UDPConn

var enabled = false
var connectedToClient = false
var isServer = false
var confirmedTick = int64(0)
var localSyncData = uint32(0)
var remoteSyncData = uint32(0)
var isStateDesynced = false
var localSyncDataTick = int64(-1)
var remoteSyncDataTick = int64(-1)
var localTickDelta = int64(0)
var remoteTickDelta = int64(0)
var inputHistory = [inputHistorySize]uint32{}
var remoteInputHistory = [inputHistorySize]uint32{}
var clientAddr net.Addr
var latency int64
var lastSyncedTick = int64(-1)

// var toSendPackets = []struct {
// 	Packet []byte
// 	Time   time.Time
// }{}
// var syncDataHistoryLocal = [inputHistorySize]uint32{}
// var syncDataHistoryRemote = [inputHistorySize]uint32{}

// Init initialises a netplay session between two players
func Init() {
	if Listen { // Host mode
		var err error
		Conn, err = net.ListenUDP("udp", &net.UDPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 1234,
		})
		if err != nil {
			log.Println("Netplay", err.Error())
			return
		}

		Conn.SetReadBuffer(1048576)

		enabled = true
		isServer = true

		input.InitializeBuffer(0)
		input.InitializeBuffer(1)
		input.LocalPlayerPort = 0
		input.RemotePlayerPort = 1

		log.Println("Netplay", "Listening.")

		buffer := make([]byte, 1024)
		_, addr, _ := Conn.ReadFrom(buffer)

		connectedToClient = true
		clientAddr = addr

		sendPacket(makeHandshakePacket(), 5)
	} else if Join { // Guest mode
		var err error
		Conn, err = net.ListenUDP("udp", &net.UDPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 1235,
		})
		if err != nil {
			log.Println("Netplay", err.Error())
			return
		}

		clientAddr = &net.UDPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 1234,
		}

		Conn.SetReadBuffer(1048576)

		enabled = true
		isServer = false

		input.InitializeBuffer(0)
		input.InitializeBuffer(1)
		input.LocalPlayerPort = 1
		input.RemotePlayerPort = 0

		log.Println("sending handshake")
		sendPacket(makeHandshakePacket(), 5)

		buffer := make([]byte, 1024)
		_, addr, _ := Conn.ReadFrom(buffer)

		connectedToClient = true
		clientAddr = addr
	}
}

// Get input from the remote player for the passed in game tick.
func getRemoteInputState(tick int64) input.PlayerState {
	if tick > confirmedTick {
		// Repeat the last confirmed input when we don't have a confirmed tick
		tick = confirmedTick
		log.Println("Predict:", confirmedTick, remoteInputHistory[(inputHistorySize+tick)%inputHistorySize])
	}
	return decodeInput(remoteInputHistory[(inputHistorySize+tick)%inputHistorySize])
}

// Get input state for the local client
func getLocalInputState(tick int64) input.PlayerState {
	return decodeInput(inputHistory[(inputHistorySize+tick)%inputHistorySize])
}

// func getLocalInputEncoded(tick int64) uint32 {
// 	return inputHistory[(inputHistorySize+tick)%inputHistorySize]
// }

// Get the sync data which is used to check for game state desync between the clients.
// func GetSyncDataLocal(tick int64) uint32 {
// 	index := (inputHistorySize + tick) % inputHistorySize
// 	return syncDataHistoryLocal[index]
// }

// Get sync data from the remote client.
// func GetSyncDataRemote(tick int64) uint32 {
// 	index := (inputHistorySize + tick) % inputHistorySize
// 	return syncDataHistoryRemote[index]
// }

// Send the inputState for the local player to the remote player for the given game tick.
func sendInputData(tick int64) {
	// Don't send input data when not connect to another player's game client.
	if !(enabled && connectedToClient) {
		return
	}

	//log.Println("Send input packet", tick)

	sendPacket(makeInputPacket(tick), 1)
}

func setLocalInput(st input.PlayerState, tick int64) {
	encodedInput := encodeInput(st)
	inputHistory[(inputHistorySize+tick)%inputHistorySize] = encodedInput
}

func setRemoteEncodedInput(encodedInput uint32, tick int64) {
	remoteInputHistory[(inputHistorySize+tick)%inputHistorySize] = encodedInput
}

// Handles sending packets to the other client. Set duplicates to something > 0 to send more than once.
func sendPacket(packet []byte, duplicates int) {
	if duplicates == 0 {
		duplicates = 1
	}

	for i := 0; i < duplicates; i++ {
		// if sendDelayFrames > 0 {
		// 	sendPacketWithDelay(packet)
		// } else {
		sendPacketRaw(packet)
		// }
	}
}

// Queues a packet to be sent later
// func sendPacketWithDelay(packet []byte) {
// 	delayedPacket := struct {
// 		Packet []byte
// 		Time   time.Time
// 	}{
// 		Packet: packet,
// 		Time:   time.Now(),
// 	}
// 	toSendPackets = append(toSendPackets, delayedPacket)
// }

// Send all packets which have been queued and who's delay time as elapsed.
// func processDelayedPackets() {
// 	newPacketList := []struct {
// 		Packet []byte
// 		Time   time.Time
// 	}{} // List of packets that haven't been sent yet.
// 	timeInterval := sendDelayFrames / 60 // How much time must pass (converting from frames into seconds)

// 	for _, data := range toSendPackets {
// 		if (time.Now().Unix() - data.Time.Unix()) > int64(timeInterval) {
// 			sendPacketRaw(data.Packet) // Send packet when enough time as passed.
// 		} else {
// 			newPacketList = append(newPacketList, data) // Keep the packet if the not enough time as passed.
// 		}
// 	}
// 	toSendPackets = newPacketList
// }

// Send a packet immediately
func sendPacketRaw(packet []byte) {
	_, err := Conn.WriteTo(packet, clientAddr)
	if err != nil {
		log.Println(err)
	}
}

// Handles receiving packets from the other client.
func receivePacket() (int, []byte, net.Addr, error) {
	buffer := make([]byte, 1024)
	Conn.SetReadDeadline(time.Now().Add(time.Millisecond))
	n, addr, err := Conn.ReadFrom(buffer)

	// if n > 0 {
	// 	log.Println("received", n, buffer[:n])
	// }

	return n, buffer[:n], addr, err
}

// Checks the queue for any incoming packets and process them.
func receiveData() {
	if !enabled {
		return
	}

	// For now we'll process all packets every frame.
	for {
		n, data, _, err := receivePacket()
		if err != nil {
			// log.Println(err)
			return
		}

		if n > 0 {
			r := bytes.NewReader(data)
			var code byte
			binary.Read(r, binary.LittleEndian, &code)

			// Handshake code must be received by both game instances before a match can begin.
			// if code == MsgCodeHandshake {
			// 	if !connectedToClient {
			// 		connectedToClient = true

			// 		// The server needs to remember the address and port in order to send data to the other cilent.
			// 		if true {
			// 			// Server needs to the other the client address and ip to know where to send data.
			// 			if isServer {
			// 				clientAddr = addr
			// 			}
			// 			log.Println("Received Handshake from: ", clientAddr.String())
			// 			// Send handshake to client.
			// 			sendPacket(makeHandshakePacket(), 5)
			// 		}
			// 	}
			// } else
			if code == MsgCodePlayerInput {
				// Break apart the packet into its parts.
				var tickDelta, receivedTick int64
				binary.Read(r, binary.LittleEndian, &tickDelta)
				binary.Read(r, binary.LittleEndian, &receivedTick)

				//log.Println("Received input", receivedTick)

				// We only care about the latest tick delta, so make sure the confirmed frame is atleast the same or newer.
				// This would work better if we added a packet count.
				if receivedTick >= confirmedTick {
					remoteTickDelta = tickDelta
				}

				if receivedTick > confirmedTick {
					if receivedTick-confirmedTick > inputDelayFrames {
						log.Println("Received packet with a tick too far ahead. Last: ", confirmedTick, "     Current: ", receivedTick)
					}

					confirmedTick = receivedTick

					// log.Println("----")
					// log.Println(confirmedTick)
					for offset := int64(sendHistorySize - 1); offset >= 0; offset-- {
						var encodedInput uint32
						binary.Read(r, binary.LittleEndian, &encodedInput)
						// Save the input history sent in the packet.
						setRemoteEncodedInput(encodedInput, receivedTick-offset)

						// log.Println(encodedInput, receivedTick-offset, offset)
					}
				}

				// NetLog("Received Tick: " .. receivedTick .. ",  Input: " .. remoteInputHistory[(confirmedTick % inputHistorySize)+1])
			} else if code == MsgCodePing {
				var pingTime int64
				binary.Read(r, binary.LittleEndian, &pingTime)
				sendPacket(makePongPacket(time.Unix(pingTime, 0)), 1)
			} else if code == MsgCodePong {
				var pongTime int64
				binary.Read(r, binary.LittleEndian, &pongTime)
				latency = time.Now().Unix() - pongTime
				//print("Got pong message: " .. latency)
			} else if code == MsgCodeSync {
				var tick int64
				var syncData uint32
				binary.Read(r, binary.LittleEndian, &tick)
				binary.Read(r, binary.LittleEndian, &syncData)
				// Ignore any tick that isn't more recent than the last sync data
				if !isStateDesynced && tick > remoteSyncDataTick {
					remoteSyncDataTick = tick
					remoteSyncData = syncData

					// Check for a desync
					isDesynced()
				}
			}
		}
	}
}

// Generate a packet containing information about player input.
func makeInputPacket(tick int64) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, MsgCodePlayerInput)
	binary.Write(buf, binary.LittleEndian, localTickDelta)
	binary.Write(buf, binary.LittleEndian, tick)

	historyIndexStart := tick - sendHistorySize + 1
	// log.Println("Make input", tick, historyIndexStart)
	for i := int64(0); i < sendHistorySize; i++ {
		encodedInput := inputHistory[(inputHistorySize+historyIndexStart+i)%inputHistorySize]
		binary.Write(buf, binary.LittleEndian, encodedInput)
		// log.Println((inputHistorySize + historyIndexStart + i) % inputHistorySize)
	}

	return buf.Bytes()
}

// Send a ping message in order to test network latency
func sendPingMessage() {
	sendPacket(makePingPacket(time.Now()), 1)
}

// Make a ping packet
func makePingPacket(t time.Time) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, MsgCodePing)
	binary.Write(buf, binary.LittleEndian, t.Unix())
	return buf.Bytes()
}

// Make pong packet
func makePongPacket(t time.Time) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, MsgCodePong)
	binary.Write(buf, binary.LittleEndian, t.Unix())
	return buf.Bytes()
}

// Sends sync data
func sendSyncData() {
	sendPacket(makeSyncDataPacket(localSyncDataTick, localSyncData), 5)
}

// Make a sync data packet
func makeSyncDataPacket(tick int64, syncData uint32) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, MsgCodeSync)
	binary.Write(buf, binary.LittleEndian, tick)
	err := binary.Write(buf, binary.LittleEndian, syncData)
	if err != nil {
		log.Println(err)
	}
	return buf.Bytes()
}

// Generate handshake packet for connecting with another client.
func makeHandshakePacket() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, MsgCodeHandshake)
	return buf.Bytes()
}

// Encodes the player input state into a compact form for network transmission.
func encodeInput(st input.PlayerState) uint32 {
	var out uint32
	for i, b := range st {
		if b {
			out |= (1 << i)
		}
	}
	return out
}

// Decodes the input from a packet generated by encodeInput().
func decodeInput(in uint32) input.PlayerState {
	st := input.PlayerState{}
	for i := range st {
		st[i] = in&(1<<i) > 0
	}
	return st
}
