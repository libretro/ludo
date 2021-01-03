package netplay

import (
	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/state"
)

const ROLLBACK_TEST_ENABLED = false
const NET_ROLLBACK_MAX_FRAMES = 10

var tickSyncing = false
var tickOffset = float64(0)
var lastConfirmedTick int64
var syncedLastUpdate = false

func Update(inputPoll, gameUpdate func()) {
	lastGameTick := state.Global.Tick
	shouldUpdate := false

	if ROLLBACK_TEST_ENABLED {
		shouldUpdate = true
	}

	// The network is update first
	if enabled {
		// First get any data that has been sent from the other client
		receiveData()

		// Send any packets that have been queued
		processDelayedPackets()

		if connectedToClient {
			// First we assume that the game can be updated, sync checks below can halt updates
			shouldUpdate = true

			if state.Global.ForcePause {
				shouldUpdate = false
			}

			// Run any rollbacks that can be processed before the next game update
			HandleRollbacks(gameUpdate)

			// Calculate the difference between remote game tick and the local. This will be used for syncing.
			// We don't use the latest local tick, but the tick for the latest input sent to the remote client.
			localTickDelta = lastGameTick - confirmedTick

			//timeSyncGraphTable[1+(lastGameTick%60)*2+1] = -1 * (localTickDelta - RemoteTickDelta) * GRAPH_UNIT_SCALE

			// Only do time sync check when the previous confirmed tick from the remote client hasn't been used yet.
			if confirmedTick > lastConfirmedTick {

				lastConfirmedTick = confirmedTick

				// Prevent updating the game when the tick difference is greater on this end.
				// This allows the game deltas to be off by 2 frames. Our timing is only accurate to one frame so any slight increase in network latency
				// would cause the game to constantly hold. You could increase this tolerance, but this would increase the advantage for one player over the other.

				// Only calculate time sync frames when we are not currently time syncing.
				if tickSyncing == false {
					// Calculate tick offset using the clock synchronization algorithm.
					// See https://en.wikipedia.org/wiki/Network_Time_Protocol#Clock_synchronization_algorithm
					tickOffset = (float64(localTickDelta) - float64(remoteTickDelta)) / 2.0
					// log.Println(tickOffset)

					// Only sync when the tick difference is more than one frame.
					if tickOffset >= 1 {
						tickSyncing = true
					}
				}

				if tickSyncing && syncedLastUpdate == false {
					shouldUpdate = false
					syncedLastUpdate = true

					tickOffset = tickOffset - 1

					// Stop time syncing when the tick difference is less than 1 so we don't overshoot
					if tickOffset < 1 {
						tickSyncing = false
					}
				} else {
					syncedLastUpdate = false
				}
			}

			// Only halt the game update based on exceeding the rollback window when the game updated hasn't previously been stopped by time sync code
			if shouldUpdate {
				// We allow the game to run for NET_ROLLBACK_MAX_FRAMES updates without having input for the current frame.
				// Once the game can no longer update, it will wait until the other player's client can catch up.
				if lastGameTick <= (confirmedTick + NET_ROLLBACK_MAX_FRAMES) {
					shouldUpdate = true
				} else {
					shouldUpdate = false
				}
			}
		}
	}

	if shouldUpdate {
		// Test rollbacks
		// TestRollbacks()

		// Poll inputs for this frame. In network mode the network manager will handle updating player command buffers.
		// updateCommandBuffers := !enabled
		// input.Poll(updateCommandBuffers)
		input.Poll()

		// Network manager will handle updating inputs.
		if enabled {
			// Update local input history
			sendInput := input.GetLatest(input.LocalPlayerPort)
			setLocalInput(sendInput, lastGameTick+NET_INPUT_DELAY)
			// log.Println(sendInput, lastGameTick+NET_INPUT_DELAY)

			// Set the input state fo[r the current tick for the remote player's character.
			input.SetState(input.LocalPlayerPort, GetLocalInputState(lastGameTick))
			input.SetState(input.RemotePlayerPort, GetRemoteInputState(lastGameTick))
		}

		// Increment the tick count only when the game actually updates.
		gameUpdate()

		state.Global.Tick++

		// Save stage after an update if testing rollbacks
		if ROLLBACK_TEST_ENABLED {
			// Save local input history for this game tick
			setLocalInput(input.GetLatest(input.LocalPlayerPort), lastGameTick)
		}

		if enabled {
			// Check whether or not the game state is confirmed to be in sync.
			// Since we previously rolled back, it's safe to set the lastSyncedTick here since we know any previous frames will be synced.
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
		// if shouldUpdate then
		// 	PacketLog("Sending Input: " .. GetLocalInputEncoded(lastGameTick + NET_INPUT_DELAY) .. ' @ ' .. lastGameTick + NET_INPUT_DELAY  )
		// end

		// Send this player's input state. We when NET_INPUT_DELAY frames ahead.
		// Note: This input comes from the last game update, so we subtract 1 to set the correct tick.
		sendInputData(state.Global.Tick - 1 + NET_INPUT_DELAY)

		// Send ping so we can test network latency.
		sendPingMessage()
	}
}
