package netplay

import (
	"hash/crc32"
	"log"
	"os"

	"github.com/libretro/ludo/state"
)

const detectDesyncs = true
const desyncCheckRate = int64(10)

// Gets the sync data to confirm the client game states are in sync
func gameGetSyncData() uint32 {
	s := state.Global.Core.SerializeSize()
	bytes, err := state.Global.Core.Serialize(s)
	if err != nil {
		log.Println(err)
		return 0
	}

	return crc32.ChecksumIEEE(bytes)
}

// Checks whether or not a game state desync has occurred between the local and remote clients.
func checkSync() {
	if !detectDesyncs {
		return
	}

	if lastSyncedTick < 0 {
		return
	}

	// Check desyncs at a fixed rate.
	if (lastSyncedTick % desyncCheckRate) != 0 {
		return
	}

	// Generate the data we'll send to the other player for testing that their game state is in sync.
	setLocalSyncData(lastSyncedTick, gameGetSyncData())

	// Send sync data everytime we've applied from the remote player to a game frame.
	sendSyncData()

	desynced, desyncFrame := isDesynced()
	if !desynced {
		return
	}

	// Detect when the sync data doesn't match then halt the game
	log.Println("Desync detected at tick: ", desyncFrame)

	os.Exit(0)
}

// Check for a desync.
func isDesynced() (bool, int64) {
	if localSyncDataTick < 0 {
		return false, 0
	}

	// When the local sync data does not match the remote data indicate a desync has occurred.
	if isStateDesynced || localSyncDataTick == remoteSyncDataTick {
		//log.Println("Desync Check at: ", localSyncDataTick)

		if localSyncData != remoteSyncData {
			log.Println(localSyncDataTick, localSyncData, remoteSyncData)
			isStateDesynced = true
			return true, localSyncDataTick
		}
	}

	return false, 0
}

// Set sync data for a game tick
func setLocalSyncData(tick int64, syncData uint32) {
	if !isStateDesynced {
		// log.Println("setLocalSyncData", tick, syncData)
		localSyncData = syncData
		localSyncDataTick = tick
	}
}
