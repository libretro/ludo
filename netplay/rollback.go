package netplay

import (
	"log"

	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/state"
)

var SAVESTATE = []byte{}
var INPUTSTATE = [input.MaxPlayers][input.MaxFrames]input.PlayerState{}
var TICK = int64(0)

func serialize() {
	//log.Println("Serialize")
	s := state.Global.Core.SerializeSize()
	bytes, err := state.Global.Core.Serialize(s)
	if err != nil {
		log.Println(err)
	}
	SAVESTATE = make([]byte, s)
	copy(SAVESTATE[:], bytes[:])

	INPUTSTATE = input.Serialize()
	TICK = state.Global.Tick
}

func unserialize() {
	log.Println("Unserialize")
	if len(SAVESTATE) == 0 {
		log.Println("Trying to unserialize a savestate of len 0")
		return
	}

	s := state.Global.Core.SerializeSize()
	err := state.Global.Core.Unserialize(SAVESTATE, s)
	if err != nil {
		log.Println(err)
	}
	input.Unserialize(INPUTSTATE)
	state.Global.Tick = TICK
}

// HandleRollbacks will rollback if needed.
func HandleRollbacks(gameUpdate func()) {
	lastGameTick := state.Global.Tick - 1
	// The input needed to resync state is available so rollback.
	// lastSyncedTick keeps track of the lastest synced game tick.
	// When the tick count for the inputs we have is more than the number of synced ticks it's possible to rerun those game updates
	// with a rollback.

	// The number of frames that's elasped since the game has been out of sync.
	// Rerun rollbackFrames number of updates.
	rollbackFrames := lastGameTick - lastSyncedTick

	// Update the graph indicating the number of rollback frames
	// rollbackGraphTable[ 1 + (lastGameTick % 60) * 2 + 1  ] = -1 * rollbackFrames * GRAPH_UNIT_SCALE

	if lastGameTick >= 0 && lastGameTick > (lastSyncedTick+1) && confirmedTick > lastSyncedTick {
		log.Println("Rollback", rollbackFrames, "frames")
		state.Global.FastForward = true

		// Must revert back to the last known synced game frame.
		unserialize()

		for i := int64(0); i < rollbackFrames; i++ {
			// Get input from the input history buffer. The network system will predict input after the last confirmed tick (for the remote player).
			input.SetState(input.LocalPlayerPort, GetLocalInputState(state.Global.Tick)) // Offset of 1 ensure it's used for the next game update.
			input.SetState(input.RemotePlayerPort, GetRemoteInputState(state.Global.Tick))

			lastRolledBackGameTick := state.Global.Tick
			gameUpdate()
			state.Global.Tick++

			// Confirm that we are indeed still synced
			if lastRolledBackGameTick <= confirmedTick {
				log.Println("Save after rollback")
				serialize()

				lastSyncedTick = lastRolledBackGameTick

				// Confirm the game clients are in sync
				checkSync()
			}
		}

		state.Global.FastForward = false
	}
}
