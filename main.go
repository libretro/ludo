package main

import (
	"flag"
	"log"
	"math"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/libretro/ludo/audio"
	"github.com/libretro/ludo/core"
	"github.com/libretro/ludo/history"
	"github.com/libretro/ludo/input"
	"github.com/libretro/ludo/menu"
	"github.com/libretro/ludo/netplay"
	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/playlists"
	"github.com/libretro/ludo/scanner"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/video"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

const ROLLBACK_TEST_ENABLED = false
const NET_ROLLBACK_MAX_FRAMES = 10
const TICK_RATE = 1.0 / 60.0
const MAX_FRAME_SKIP = 25

func update() {
	lastGameTick := state.Global.Tick
	updateGame := false

	if ROLLBACK_TEST_ENABLED {
		updateGame = true
	}

	// The network is update first
	if netplay.Enabled {
		// First get any data that has been sent from the other client
		netplay.ReceiveData()

		// Send any packets that have been queued
		netplay.ProcessDelayedPackets()

		if netplay.ConnectedToClient {
			// First we assume that the game can be updated, sync checks below can halt updates
			updateGame = true

			if state.Global.ForcePause {
				updateGame = false
			}

			// Run any rollbacks that can be processed before the next game update
			netplay.HandleRollbacks(gameUpdate)

			// Calculate the difference between remote game tick and the local. This will be used for syncing.
			// We don't use the latest local tick, but the tick for the latest input sent to the remote client.
			netplay.LocalTickDelta = lastGameTick - netplay.ConfirmedTick

			//timeSyncGraphTable[1+(lastGameTick%60)*2+1] = -1 * (netplay.LocalTickDelta - netplay.RemoteTickDelta) * GRAPH_UNIT_SCALE

			// Only do time sync check when the previous confirmed tick from the remote client hasn't been used yet.
			if netplay.ConfirmedTick > state.Global.LastConfirmedTick {

				state.Global.LastConfirmedTick = netplay.ConfirmedTick

				// Prevent updating the game when the tick difference is greater on this end.
				// This allows the game deltas to be off by 2 frames. Our timing is only accurate to one frame so any slight increase in network latency
				// would cause the game to constantly hold. You could increase this tolerance, but this would increase the advantage for one player over the other.

				// Only calculate time sync frames when we are not currently time syncing.
				if netplay.TickSyncing == false {
					// Calculate tick offset using the clock synchronization algorithm.
					// See https://en.wikipedia.org/wiki/Network_Time_Protocol#Clock_synchronization_algorithm
					netplay.TickOffset = (float64(netplay.LocalTickDelta) - float64(netplay.RemoteTickDelta)) / 2.0
					// log.Println(netplay.TickOffset)

					// Only sync when the tick difference is more than one frame.
					if netplay.TickOffset >= 1 {
						netplay.TickSyncing = true
					}
				}

				if netplay.TickSyncing && state.Global.SyncedLastUpdate == false {
					updateGame = false
					state.Global.SyncedLastUpdate = true

					netplay.TickOffset = netplay.TickOffset - 1

					// Stop time syncing when the tick difference is less than 1 so we don't overshoot
					if netplay.TickOffset < 1 {
						netplay.TickSyncing = false
					}
				} else {
					state.Global.SyncedLastUpdate = false
				}
			}

			// Only halt the game update based on exceeding the rollback window when the game updated hasn't previously been stopped by time sync code
			if updateGame {
				// We allow the game to run for NET_ROLLBACK_MAX_FRAMES updates without having input for the current frame.
				// Once the game can no longer update, it will wait until the other player's client can catch up.
				if lastGameTick <= (netplay.ConfirmedTick + NET_ROLLBACK_MAX_FRAMES) {
					updateGame = true
				} else {
					updateGame = false
				}
			}
		}
	}

	if updateGame {
		// Test rollbacks
		// TestRollbacks()

		// Poll inputs for this frame. In network mode the network manager will handle updating player command buffers.
		// updateCommandBuffers := !netplay.Enabled
		// input.Poll(updateCommandBuffers)
		input.Poll()

		// Network manager will handle updating inputs.
		if netplay.Enabled {
			// Update local input history
			sendInput := input.GetLatest(input.LocalPlayerPort)
			netplay.SetLocalInput(sendInput, lastGameTick+netplay.NET_INPUT_DELAY)
			// log.Println(sendInput, lastGameTick+netplay.NET_INPUT_DELAY)

			// Set the input state fo[r the current tick for the remote player's character.
			input.SetState(input.LocalPlayerPort, netplay.GetLocalInputState(lastGameTick))
			input.SetState(input.RemotePlayerPort, netplay.GetRemoteInputState(lastGameTick))
		}

		// Increment the tick count only when the game actually updates.
		gameUpdate()

		state.Global.Tick++

		// Save stage after an update if testing rollbacks
		if ROLLBACK_TEST_ENABLED {
			// Save local input history for this game tick
			netplay.SetLocalInput(input.GetLatest(input.LocalPlayerPort), lastGameTick)
		}

		if netplay.Enabled {
			// Check whether or not the game state is confirmed to be in sync.
			// Since we previously rolled back, it's safe to set the LastSyncedTick here since we know any previous frames will be synced.
			if netplay.LastSyncedTick+1 == lastGameTick && lastGameTick <= netplay.ConfirmedTick {
				// Increment the synced tick number if we have inputs
				netplay.LastSyncedTick = lastGameTick

				// Applied the remote player's input, so this game frame should synced.
				netplay.Serialize()

				// Confirm the game clients are in sync
				netplay.CheckSync()
			}
		}
	}

	// Since our input is update in gameupdate() we want to send the input as soon as possible.
	// Previously this as happening before the gameupdate() and adding uneeded latency.
	if netplay.Enabled && netplay.ConnectedToClient {
		// if updateGame then
		// 	PacketLog("Sending Input: " .. netplay.GetLocalInputEncoded(lastGameTick + NET_INPUT_DELAY) .. ' @ ' .. lastGameTick + NET_INPUT_DELAY  )
		// end

		// Send this player's input state. We when NET_INPUT_DELAY frames ahead.
		// Note: This input comes from the last game update, so we subtract 1 to set the correct tick.
		netplay.SendInputData(state.Global.Tick - 1 + netplay.NET_INPUT_DELAY)

		// Send ping so we can test network latency.
		netplay.SendPingMessage()
	}
}

func runLoop(vid *video.Video, m *menu.Menu) {
	currTime := time.Now()
	prevTime := time.Now()
	lag := float64(0)
	for !vid.Window.ShouldClose() {
		currTime = time.Now()
		dt := float64(currTime.Sub(prevTime)) / 1000000000

		glfw.PollEvents()
		m.ProcessHotkeys()
		vid.ResizeViewport()
		m.UpdatePalette()

		state.Global.ForcePause = vid.Window.GetKey(glfw.KeySpace) == glfw.Press

		// Cap number of Frames that can be skipped so lag doesn't accumulate
		lag = math.Min(lag+dt, TICK_RATE*MAX_FRAME_SKIP)

		for lag >= TICK_RATE {
			update()
			lag -= TICK_RATE
		}

		vid.Render()
		glfw.SwapInterval(0)
		vid.Window.SwapBuffers()
		prevTime = currTime
	}
}

func gameUpdate() {
	if state.Global.CoreRunning {
		// if input.LocalPlayerPort == 0 {
		// 	log.Println("----> updating", state.Global.Tick, gameGetSyncData(), netplay.GetLocalInputState(state.Global.Tick))
		// } else {
		// 	log.Println("----> updating", state.Global.Tick, gameGetSyncData(), netplay.GetRemoteInputState(state.Global.Tick))
		// }
		state.Global.Core.Run()
		//log.Println("----> done updating")
		if state.Global.Core.FrameTimeCallback != nil {
			state.Global.Core.FrameTimeCallback.Callback(state.Global.Core.FrameTimeCallback.Reference)
		}
		if state.Global.Core.AudioCallback != nil {
			state.Global.Core.AudioCallback.Callback()
		}
	}
}

func main() {
	err := settings.Load()
	if err != nil {
		log.Println("[Settings]: Loading failed:", err)
		log.Println("[Settings]: Using default settings")
	}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flag.StringVar(&state.Global.CorePath, "L", "", "Path to the libretro core")
	flag.BoolVar(&state.Global.Verbose, "v", false, "Verbose logs")
	flag.BoolVar(&state.Global.LudOS, "ludos", false, "Expose the features related to LudOS")
	flag.BoolVar(&netplay.Listen, "listen", false, "For the netplay server")
	flag.BoolVar(&netplay.Join, "join", false, "For the netplay client")
	flag.Parse()
	args := flag.Args()

	var gamePath string
	if len(args) > 0 {
		gamePath = args[0]
	}

	if err := glfw.Init(); err != nil {
		log.Fatalln("Failed to initialize glfw", err)
	}
	defer glfw.Terminate()

	state.Global.DB, err = scanner.LoadDB(settings.Current.DatabaseDirectory)
	if err != nil {
		log.Println("Can't load game database:", err)
	}

	playlists.Load()

	history.Load()

	vid := video.Init(settings.Current.VideoFullscreen)

	audio.Init()

	m := menu.Init(vid)

	core.Init(vid)

	input.Init(vid)

	if len(state.Global.CorePath) > 0 {
		err := core.Load(state.Global.CorePath)
		if err != nil {
			panic(err)
		}
	}

	if len(gamePath) > 0 {
		err := core.LoadGame(gamePath)
		if err != nil {
			ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
		} else {
			m.WarpToQuickMenu()
		}
	}

	// No game running? display the menu
	state.Global.MenuActive = !state.Global.CoreRunning

	runLoop(vid, m)

	// Unload and deinit in the core.
	core.Unload()
}
