package netplay

import (
	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/state"
)

const rollbackTestEnabled = false
const rollbackMaxFrames = 10

var tickSyncing = false
var tickOffset = float64(0)
var lastConfirmedTick int64
var syncedLastUpdate = false

// Update processes a frame of the netplay, taking care of polling inputs and executing the game
func Update(inputPoll, gameUpdate func()) {
	lastGameTick := state.Global.Tick
	shouldUpdate := false

	if rollbackTestEnabled {
		shouldUpdate = true
	}

	// The network is update first
	if enabled {
		// First get any data that has been sent from the other client
		receiveData()

		if connectedToClient {
			// First we assume that the game can be updated, sync checks below can halt updates
			shouldUpdate = true

			if state.Global.ForcePause {
				shouldUpdate = false
			}

			// Run any rollbacks that can be processed before the next game update
			handleRollbacks(gameUpdate)

			// Calculate the difference between remote game tick and the local. This will be used for syncing.
			// We don't use the latest local tick, but the tick for the latest input sent to the remote client.
			localTickDelta = lastGameTick - confirmedTick

			// Only do time sync check when the previous confirmed tick from the remote client hasn't been used yet.
			if confirmedTick > lastConfirmedTick {

				lastConfirmedTick = confirmedTick

				// Prevent updating the game when the tick difference is greater on this end.
				// This allows the game deltas to be off by 2 frames.
				// Our timing is only accurate to one frame so any slight increase in network latency would cause the
				// game to constantly hold. You could increase this tolerance, but this would increase the advantage
				// for one player over the other.

				// Only calculate time sync frames when we are not currently time syncing.
				if tickSyncing == false {
					// Calculate tick offset using the clock synchronization algorithm.
					// See https://en.wikipedia.org/wiki/Network_Time_Protocol#Clock_synchronization_algorithm
					tickOffset = (float64(localTickDelta) - float64(remoteTickDelta)) / 2.0

					// Only sync when the tick difference is more than one frame.
					if tickOffset >= 1 {
						tickSyncing = true
					}
				}

				if tickSyncing && syncedLastUpdate == false {
					shouldUpdate = false
					syncedLastUpdate = true

					tickOffset--

					// Stop time syncing when the tick difference is less than 1 so we don't overshoot
					if tickOffset < 1 {
						tickSyncing = false
					}
				} else {
					syncedLastUpdate = false
				}
			}

			// Only halt the game update based on exceeding the rollback window when the game updated hasn't previously
			// been stopped by time sync code
			if shouldUpdate {
				// We allow the game to run for rollbackMaxFrames updates without having input for the current frame.
				// Once the game can no longer update, it will wait until the other player's client can catch up.
				shouldUpdate = lastGameTick <= (confirmedTick + rollbackMaxFrames)
			}
		}
	}

	if shouldUpdate {
		// Poll inputs for this frame.
		input.Poll()

		// Network manager will handle updating inputs.
		if enabled {
			// Update local input history
			sendInput := input.GetLatest(input.LocalPlayerPort)
			setLocalInput(sendInput, lastGameTick+inputDelayFrames)
			// log.Println(sendInput, lastGameTick+inputDelayFrames)

			// Set the input state fo[r the current tick for the remote player's character.
			input.SetState(input.LocalPlayerPort, getLocalInputState(lastGameTick))
			input.SetState(input.RemotePlayerPort, getRemoteInputState(lastGameTick))
		}

		// Increment the tick count only when the game actually updates.
		gameUpdate()

		state.Global.Tick++

		// Save stage after an update if testing rollbacks
		if rollbackTestEnabled {
			// Save local input history for this game tick
			setLocalInput(input.GetLatest(input.LocalPlayerPort), lastGameTick)
		}

		if enabled {
			// Check whether or not the game state is confirmed to be in sync.
			// Since we previously rolled back, it's safe to set the lastSyncedTick here since we know any previous
			// frames will be synced.
			if lastSyncedTick+1 == lastGameTick && lastGameTick <= confirmedTick {
				// Increment the synced tick number if we have inputs
				lastSyncedTick = lastGameTick

				// Applied the remote player's input, so this game frame should synced.
				serialize()

				// Confirm the game clients are in sync
				checkSync()
			}
		}
	}

	// Since our input is update in gameupdate() we want to send the input as soon as possible.
	// Previously this as happening before the gameupdate() and adding uneeded latency.
	if enabled && connectedToClient {
		// Send this player's input state. We when inputDelayFrames frames ahead.
		// Note: This input comes from the last game update, so we subtract 1 to set the correct tick.
		sendInputData(state.Global.Tick - 1 + inputDelayFrames)

		// Send ping so we can test network latency.
		sendPingMessage()
	}
}
