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
	s := state.Core.SerializeSize()
	var err error
	saved.GameState, err = state.Core.Serialize(s)
	if err != nil {
		log.Println(err)
	}

	saved.Inputs = input.Serialize()
	saved.Tick = state.Tick
}

func unserialize() {
	if len(saved.GameState) == 0 {
		log.Println("Trying to unserialize a savestate of len 0")
		return
	}

	s := state.Core.SerializeSize()
	err := state.Core.Unserialize(saved.GameState, s)
	if err != nil {
		log.Println(err)
	}
	input.Unserialize(saved.Inputs)
	state.Tick = saved.Tick
}

// handleRollbacks will rollback if needed.
func handleRollbacks() {
	lastGameTick := state.Tick - 1
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
		state.FastForward = true

		// Must revert back to the last known synced game frame.
		unserialize()

		for i := int64(0); i < rollbackFrames; i++ {
			// Get input from the input history buffer.
			// The network system can predict input after the last confirmed tick (for the remote player).
			input.SetState(input.LocalPlayerPort, getLocalInputState(state.Tick))
			input.SetState(input.RemotePlayerPort, getRemoteInputState(state.Tick))

			lastRolledBackGameTick := state.Tick
			gameUpdate()
			state.Tick++

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
		state.FastForward = false
	}
}
