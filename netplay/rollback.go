package netplay

import (
	"log"

	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/state"
)

var saved struct {
	GameState []byte
	Inputs    [input.MaxPlayers][input.MaxFrames]input.PlayerState
	Tick      int64
}

func serialize() {
	s := state.Global.Core.SerializeSize()
	var err error
	saved.GameState, err = state.Global.Core.Serialize(s)
	if err != nil {
		log.Println(err)
	}

	saved.Inputs = input.Serialize()
	saved.Tick = state.Global.Tick
}

func unserialize() {
	if len(saved.GameState) == 0 {
		log.Println("Trying to unserialize a savestate of len 0")
		return
	}

	s := state.Global.Core.SerializeSize()
	err := state.Global.Core.Unserialize(saved.GameState, s)
	if err != nil {
		log.Println(err)
	}
	input.Unserialize(saved.Inputs)
	state.Global.Tick = saved.Tick
}

// handleRollbacks will rollback if needed.
func handleRollbacks() {
	lastGameTick := state.Global.Tick - 1
	// The input needed to resync state is available so rollback.
	// lastSyncedTick keeps track of the lastest synced game tick.
	// When the tick count for the inputs we have is more than the number of synced ticks it's possible to rerun those
	// game updates with a rollback.

	if lastGameTick >= 0 && lastGameTick > (lastSyncedTick+1) && confirmedTick > lastSyncedTick {

		// The number of frames that's elasped since the game has been out of sync.
		// Rerun rollbackFrames number of updates.
		rollbackFrames := lastGameTick - lastSyncedTick

		log.Println("Rollback", rollbackFrames, "frames")

		// Disable audio because audio is blocking
		state.Global.FastForward = true

		// Must revert back to the last known synced game frame.
		unserialize()

		for i := int64(0); i < rollbackFrames; i++ {
			// Get input from the input history buffer.
			// The network system can predict input after the last confirmed tick (for the remote player).
			input.SetState(input.LocalPlayerPort, getLocalInputState(state.Global.Tick))
			input.SetState(input.RemotePlayerPort, getRemoteInputState(state.Global.Tick))

			lastRolledBackGameTick := state.Global.Tick
			gameUpdate()
			state.Global.Tick++

			// Confirm that we are indeed still synced
			if lastRolledBackGameTick <= confirmedTick {
				log.Println("Saving after a rollback")

				serialize()

				lastSyncedTick = lastRolledBackGameTick

				// Confirm the game clients are in sync
				checkSync()
			}
		}

		// Enable audio again
		state.Global.FastForward = false
	}
}
