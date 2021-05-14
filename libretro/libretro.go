/*
Package libretro is a cgo binding for the libretro API.

Libretro is a simple but powerful development interface that allows for the easy creation of
emulators, games and multimedia applications that can plug straight into any libretro-compatible
frontend. This development interface is open to others so that they can run these pluggable emulator
and game cores also in their own programs or devices. */
package libretro

/*
#include "libretro.h"
#include <stdlib.h>
#include <stdio.h>
#include <string.h>

void bridge_retro_init(void *f);
void bridge_retro_deinit(void *f);
unsigned bridge_retro_api_version(void *f);
void bridge_retro_get_system_info(void *f, struct retro_system_info *si);
void bridge_retro_get_system_av_info(void *f, struct retro_system_av_info *si);
bool bridge_retro_set_environment(void *f, void *callback);
void bridge_retro_set_video_refresh(void *f, void *callback);
void bridge_retro_set_controller_port_device(void *f, unsigned port, unsigned device);
void bridge_retro_set_input_poll(void *f, void *callback);
void bridge_retro_set_input_state(void *f, void *callback);
void bridge_retro_set_audio_sample(void *f, void *callback);
void bridge_retro_set_audio_sample_batch(void *f, void *callback);
bool bridge_retro_load_game(void *f, struct retro_game_info *gi);
bool bridge_retro_serialize(void *f, void *data, size_t size);
bool bridge_retro_unserialize(void *f, void *data, size_t size);
size_t bridge_retro_serialize_size(void *f);
void bridge_retro_unload_game(void *f);
void bridge_retro_run(void *f);
void bridge_retro_reset(void *f);
void bridge_retro_frame_time_callback(retro_frame_time_callback_t f, retro_usec_t usec);
void bridge_retro_audio_callback(retro_audio_callback_t f);
void bridge_retro_audio_set_state(retro_audio_set_state_callback_t f, bool state);
size_t bridge_retro_get_memory_size(void *f, unsigned id);
void* bridge_retro_get_memory_data(void *f, unsigned id);

bool coreEnvironment_cgo(unsigned cmd, void *data);
void coreVideoRefresh_cgo(void *data, unsigned width, unsigned height, size_t pitch);
void coreInputPoll_cgo();
void coreAudioSample_cgo(int16_t left, int16_t right);
size_t coreAudioSampleBatch_cgo(const int16_t *data, size_t frames);
int16_t coreInputState_cgo(unsigned port, unsigned device, unsigned index, unsigned id);
void coreLog_cgo(enum retro_log_level level, const char *msg);
int64_t coreGetTimeUsec_cgo();
*/
import "C"
import (
	"errors"
	"strings"
	"unsafe"
)

// GameGeometry represents the geometry of a game, with its aspect ratio, with and height
type GameGeometry struct {
	AspectRatio float64
	BaseWidth   int
	BaseHeight  int
	MaxWidth    int
	MaxHeight   int
}

// GameInfo stores information about a ROM
type GameInfo struct {
	Path string
	Size int64
	Data unsafe.Pointer
}

// SystemInfo stores informations about the emulated system
type SystemInfo struct {
	LibraryName     string
	LibraryVersion  string
	ValidExtensions string
	NeedFullpath    bool
	BlockExtract    bool
}

// SystemTiming stores informations about the timing of the emulated system
type SystemTiming struct {
	FPS        float64
	SampleRate float64
}

// SystemAVInfo stores informations about the emulated system audio and video
type SystemAVInfo struct {
	Geometry GameGeometry
	Timing   SystemTiming
}

// Variable is a key value pair that represents a core option
type Variable C.struct_retro_variable

// Key returns the key of a Variable as a string
func (v *Variable) Key() string {
	return C.GoString(v.key)
}

// Desc returns the description of a Variable as a string
func (v *Variable) Desc() string {
	val := C.GoString(v.value)
	s := strings.Split(val, "; ")
	return s[0]
}

// Choices returns the list of possible choices for a given Variable
func (v *Variable) Choices() []string {
	val := C.GoString(v.value)
	s := strings.Split(val, "; ")
	return strings.Split(s[1], "|")
}

// SetValue sets the value of a Variable
func (v *Variable) SetValue(val string) {
	s := &v.value
	*s = C.CString(val)
}

// DefaultValue returns the default value of a Variable
func (v *Variable) DefaultValue() string {
	val := C.GoString(v.value)
	s := strings.Split(val, "; ")
	return strings.Split(s[1], "|")[0]
}

// CoreOptionValue represents the value of a core option in the version 1 of the core options API
type CoreOptionValue C.struct_retro_core_option_value

// Value is the expected option value
func (cov *CoreOptionValue) Value() string {
	return C.GoString(cov.value)
}

// Label is the human readable value label of the CoreOptionValue.
// If NULL, value itself will be displayed by the frontend
func (cov *CoreOptionValue) Label() string {
	return C.GoString(cov.label)
}

// CoreOptionDefinition represents a core option in the version 1 of the core options API
type CoreOptionDefinition C.struct_retro_core_option_definition

// Key returns the key of a CoreOptionDefinition as a string
func (cod *CoreOptionDefinition) Key() string {
	return C.GoString(cod.key)
}

// Desc returns the name of a CoreOptionDefinition as a string
func (cod *CoreOptionDefinition) Desc() string {
	return C.GoString(cod.desc)
}

// Info returns the detailed description of a CoreOptionDefinition as a string
func (cod *CoreOptionDefinition) Info() string {
	return C.GoString(cod.info)
}

// Values returns the possible values of a CoreOptionDefinition as a string
func (cod *CoreOptionDefinition) Values() []CoreOptionValue {
	values := []CoreOptionValue{}

	for i := 0; i < C.RETRO_NUM_CORE_OPTION_VALUES_MAX; i++ {
		v := (C.struct_retro_core_option_value)(cod.values[i])
		if v.value == nil {
			break
		}
		values = append(values, (CoreOptionValue)(v))
	}

	return values
}

// Choices returns the CoreOptionDefinition values as a string slice for compatibility with options v0
func (cod *CoreOptionDefinition) Choices() []string {
	choices := []string{}

	for i := 0; i < C.RETRO_NUM_CORE_OPTION_VALUES_MAX; i++ {
		v := (C.struct_retro_core_option_value)(cod.values[i])
		if v.value == nil {
			break
		}
		choices = append(choices, C.GoString(v.value))
	}

	return choices
}

// DefaultValue returns the default value of a CoreOptionDefinition as a string
func (cod *CoreOptionDefinition) DefaultValue() string {
	return C.GoString(cod.default_value)
}

// FrameTimeCallback stores the frame time callback itself and the reference time
type FrameTimeCallback struct {
	Callback  func(int64)
	Reference int64
}

// AudioCallback stores the audio callback itself and the SetState callback
type AudioCallback struct {
	Callback func()
	SetState func(bool)
}

// The pixel format the core must use to render into data.
// This format could differ from the format used in SET_PIXEL_FORMAT.
// Set by frontend in GET_CURRENT_SOFTWARE_FRAMEBUFFER.
const (
	PixelFormat0RGB1555 = uint32(C.RETRO_PIXEL_FORMAT_0RGB1555)
	PixelFormatXRGB8888 = uint32(C.RETRO_PIXEL_FORMAT_XRGB8888)
	PixelFormatRGB565   = uint32(C.RETRO_PIXEL_FORMAT_RGB565)
)

// Libretro's fundamental device abstractions.
//
// Libretro's input system consists of some standardized device types,
// such as a joypad (with/without analog), mouse, keyboard, lightgun
// and a pointer.
//
// The functionality of these devices are fixed, and individual cores
// map their own concept of a controller to libretro's abstractions.
// This makes it possible for frontends to map the abstract types to a
// real input device, and not having to worry about binding input
// correctly to arbitrary controller layouts.
const (
	// DeviceNone means that input is disabled.
	DeviceNone = uint32(C.RETRO_DEVICE_NONE)

	// DeviceJoypad represents the RetroPad. It is essentially a Super Nintendo
	// controller, but with additional L2/R2/L3/R3 buttons, similar to a
	// PS1 DualShock.
	DeviceJoypad = uint32(C.RETRO_DEVICE_JOYPAD)

	// DeviceMouse is a simple mouse, similar to Super Nintendo's mouse.
	// X and Y coordinates are reported relatively to last poll (poll callback).
	// It is up to the libretro implementation to keep track of where the mouse
	// pointer is supposed to be on the screen.
	// The frontend must make sure not to interfere with its own hardware
	// mouse pointer.
	DeviceMouse = uint32(C.RETRO_DEVICE_MOUSE)

	// DeviceKeyboard lets one poll for raw key pressed.
	// It is poll based, so input callback will return with the current
	// pressed state.
	// For event/text based keyboard input, see
	// RETRO_ENVIRONMENT_SET_KEYBOARD_CALLBACK.
	DeviceKeyboard = uint32(C.RETRO_DEVICE_KEYBOARD)

	// DeviceLightgun X/Y coordinates are reported relatively to last poll,
	// similar to mouse. */
	DeviceLightgun = uint32(C.RETRO_DEVICE_LIGHTGUN)

	// DeviceAnalog device is an extension to JOYPAD (RetroPad).
	// Similar to DualShock it adds two analog sticks.
	// This is treated as a separate device type as it returns values in the
	// full analog range of [-0x8000, 0x7fff]. Positive X axis is right.
	// Positive Y axis is down.
	// Only use ANALOG type when polling for analog values of the axes.
	DeviceAnalog = uint32(C.RETRO_DEVICE_ANALOG)
)

// Buttons for the RetroPad (JOYPAD).
// The placement of these is equivalent to placements on the
// Super Nintendo controller.
// L2/R2/L3/R3 buttons correspond to the PS1 DualShock.
const (
	DeviceIDJoypadB      = uint32(C.RETRO_DEVICE_ID_JOYPAD_B)
	DeviceIDJoypadY      = uint32(C.RETRO_DEVICE_ID_JOYPAD_Y)
	DeviceIDJoypadSelect = uint32(C.RETRO_DEVICE_ID_JOYPAD_SELECT)
	DeviceIDJoypadStart  = uint32(C.RETRO_DEVICE_ID_JOYPAD_START)
	DeviceIDJoypadUp     = uint32(C.RETRO_DEVICE_ID_JOYPAD_UP)
	DeviceIDJoypadDown   = uint32(C.RETRO_DEVICE_ID_JOYPAD_DOWN)
	DeviceIDJoypadLeft   = uint32(C.RETRO_DEVICE_ID_JOYPAD_LEFT)
	DeviceIDJoypadRight  = uint32(C.RETRO_DEVICE_ID_JOYPAD_RIGHT)
	DeviceIDJoypadA      = uint32(C.RETRO_DEVICE_ID_JOYPAD_A)
	DeviceIDJoypadX      = uint32(C.RETRO_DEVICE_ID_JOYPAD_X)
	DeviceIDJoypadL      = uint32(C.RETRO_DEVICE_ID_JOYPAD_L)
	DeviceIDJoypadR      = uint32(C.RETRO_DEVICE_ID_JOYPAD_R)
	DeviceIDJoypadL2     = uint32(C.RETRO_DEVICE_ID_JOYPAD_L2)
	DeviceIDJoypadR2     = uint32(C.RETRO_DEVICE_ID_JOYPAD_R2)
	DeviceIDJoypadL3     = uint32(C.RETRO_DEVICE_ID_JOYPAD_L3)
	DeviceIDJoypadR3     = uint32(C.RETRO_DEVICE_ID_JOYPAD_R3)
)

// Mask used to identify joypads
const (
	DeviceIDJoypadMask = uint32(C.RETRO_DEVICE_ID_JOYPAD_MASK)
)

// Index / Id values for ANALOG device.
const (
	DeviceIndexAnalogLeft   = uint32(C.RETRO_DEVICE_INDEX_ANALOG_LEFT)
	DeviceIndexAnalogRight  = uint32(C.RETRO_DEVICE_INDEX_ANALOG_RIGHT)
	DeviceIndexAnalogButton = uint32(C.RETRO_DEVICE_INDEX_ANALOG_BUTTON)
	DeviceIDAnalogX         = uint32(C.RETRO_DEVICE_ID_ANALOG_X)
	DeviceIDAnalogY         = uint32(C.RETRO_DEVICE_ID_ANALOG_Y)
)

// Environment callback API. See libretro.h for details
const (
	EnvironmentSetRotation                      = uint32(C.RETRO_ENVIRONMENT_SET_ROTATION)
	EnvironmentGetOverscan                      = uint32(C.RETRO_ENVIRONMENT_GET_OVERSCAN) // Deprecated
	EnvironmentGetCanDupe                       = uint32(C.RETRO_ENVIRONMENT_GET_CAN_DUPE)
	EnvironmentSetMessage                       = uint32(C.RETRO_ENVIRONMENT_SET_MESSAGE)
	EnvironmentShutdown                         = uint32(C.RETRO_ENVIRONMENT_SHUTDOWN)
	EnvironmentSetPerformanceLevel              = uint32(C.RETRO_ENVIRONMENT_SET_PERFORMANCE_LEVEL)
	EnvironmentGetSystemDirectory               = uint32(C.RETRO_ENVIRONMENT_GET_SYSTEM_DIRECTORY)
	EnvironmentSetPixelFormat                   = uint32(C.RETRO_ENVIRONMENT_SET_PIXEL_FORMAT)
	EnvironmentSetInputDescriptors              = uint32(C.RETRO_ENVIRONMENT_SET_INPUT_DESCRIPTORS)
	EnvironmentSetKeyboardCallback              = uint32(C.RETRO_ENVIRONMENT_SET_KEYBOARD_CALLBACK)
	EnvironmentSetDiskControlInterface          = uint32(C.RETRO_ENVIRONMENT_SET_DISK_CONTROL_INTERFACE)
	EnvironmentSetHWRender                      = uint32(C.RETRO_ENVIRONMENT_SET_HW_RENDER)
	EnvironmentGetVariable                      = uint32(C.RETRO_ENVIRONMENT_GET_VARIABLE)
	EnvironmentSetVariables                     = uint32(C.RETRO_ENVIRONMENT_SET_VARIABLES)
	EnvironmentGetVariableUpdate                = uint32(C.RETRO_ENVIRONMENT_GET_VARIABLE_UPDATE)
	EnvironmentSetSupportNoGame                 = uint32(C.RETRO_ENVIRONMENT_SET_SUPPORT_NO_GAME)
	EnvironmentGetLibretroPath                  = uint32(C.RETRO_ENVIRONMENT_GET_LIBRETRO_PATH)
	EnvironmentSetFrameTimeCallback             = uint32(C.RETRO_ENVIRONMENT_SET_FRAME_TIME_CALLBACK)
	EnvironmentSetAudioCallback                 = uint32(C.RETRO_ENVIRONMENT_SET_AUDIO_CALLBACK)
	EnvironmentGetRumbleInterface               = uint32(C.RETRO_ENVIRONMENT_GET_RUMBLE_INTERFACE)
	EnvironmentGetInputDeviceCapabilities       = uint32(C.RETRO_ENVIRONMENT_GET_INPUT_DEVICE_CAPABILITIES)
	EnvironmentGetSensorInterface               = uint32(C.RETRO_ENVIRONMENT_GET_SENSOR_INTERFACE)
	EnvironmentGetCameraInterface               = uint32(C.RETRO_ENVIRONMENT_GET_CAMERA_INTERFACE)
	EnvironmentGetLogInterface                  = uint32(C.RETRO_ENVIRONMENT_GET_LOG_INTERFACE)
	EnvironmentGetPerfInterface                 = uint32(C.RETRO_ENVIRONMENT_GET_PERF_INTERFACE)
	EnvironmentGetLocationInterface             = uint32(C.RETRO_ENVIRONMENT_GET_LOCATION_INTERFACE)
	EnvironmentGetCoreAssetDirectory            = uint32(C.RETRO_ENVIRONMENT_GET_CORE_ASSETS_DIRECTORY)
	EnvironmentGetSaveDirectory                 = uint32(C.RETRO_ENVIRONMENT_GET_SAVE_DIRECTORY)
	EnvironmentSetSystemAVInfo                  = uint32(C.RETRO_ENVIRONMENT_SET_SYSTEM_AV_INFO)
	EnvironmentSetProcAddressCallback           = uint32(C.RETRO_ENVIRONMENT_SET_PROC_ADDRESS_CALLBACK)
	EnvironmentSetSubsystemInfo                 = uint32(C.RETRO_ENVIRONMENT_SET_SUBSYSTEM_INFO)
	EnvironmentSetControllerInfo                = uint32(C.RETRO_ENVIRONMENT_SET_CONTROLLER_INFO)
	EnvironmentSetMemoryMaps                    = uint32(C.RETRO_ENVIRONMENT_SET_MEMORY_MAPS)
	EnvironmentSetGeometry                      = uint32(C.RETRO_ENVIRONMENT_SET_GEOMETRY)
	EnvironmentGetUsername                      = uint32(C.RETRO_ENVIRONMENT_GET_USERNAME)
	EnvironmentGetLanguage                      = uint32(C.RETRO_ENVIRONMENT_GET_LANGUAGE)
	EnvironmentGetCurrentSoftwareFramebuffer    = uint32(C.RETRO_ENVIRONMENT_GET_CURRENT_SOFTWARE_FRAMEBUFFER)
	EnvironmentGetHWRenderInterface             = uint32(C.RETRO_ENVIRONMENT_GET_HW_RENDER_INTERFACE)
	EnvironmentSetSupportAchievements           = uint32(C.RETRO_ENVIRONMENT_SET_SUPPORT_ACHIEVEMENTS)
	EnvironmentSetHWContextNegociationInterface = uint32(C.RETRO_ENVIRONMENT_SET_HW_RENDER_CONTEXT_NEGOTIATION_INTERFACE)
	EnvironmentSetSerializationQuirks           = uint32(C.RETRO_ENVIRONMENT_SET_SERIALIZATION_QUIRKS)
	EnvironmentSetHWSharedContext               = uint32(C.RETRO_ENVIRONMENT_SET_HW_SHARED_CONTEXT)
	EnvironmentGetVFSInterface                  = uint32(C.RETRO_ENVIRONMENT_GET_VFS_INTERFACE)
	EnvironmentGetLEDInterface                  = uint32(C.RETRO_ENVIRONMENT_GET_LED_INTERFACE)
	EnvironmentGetAudioVideoEnable              = uint32(C.RETRO_ENVIRONMENT_GET_AUDIO_VIDEO_ENABLE)
	EnvironmentGetMIDIInterface                 = uint32(C.RETRO_ENVIRONMENT_GET_MIDI_INTERFACE)
	EnvironmentGetFastforwarding                = uint32(C.RETRO_ENVIRONMENT_GET_FASTFORWARDING)
	EnvironmentGetTargetRefreshRate             = uint32(C.RETRO_ENVIRONMENT_GET_TARGET_REFRESH_RATE)
	EnvironmentGetInputBitmasks                 = uint32(C.RETRO_ENVIRONMENT_GET_INPUT_BITMASKS)
	EnvironmentGetCoreOptionsVersion            = uint32(C.RETRO_ENVIRONMENT_GET_CORE_OPTIONS_VERSION)
	EnvironmentSetCoreOptions                   = uint32(C.RETRO_ENVIRONMENT_SET_CORE_OPTIONS)
	EnvironmentSetCoreOptionsIntl               = uint32(C.RETRO_ENVIRONMENT_SET_CORE_OPTIONS_INTL)
	EnvironmentSetCoreOptionsDisplay            = uint32(C.RETRO_ENVIRONMENT_SET_CORE_OPTIONS_DISPLAY)
	EnvironmentGetPrefferedHWRender             = uint32(C.RETRO_ENVIRONMENT_GET_PREFERRED_HW_RENDER)
	EnvironmentGetDiskControlInterfaceVersion   = uint32(C.RETRO_ENVIRONMENT_GET_DISK_CONTROL_INTERFACE_VERSION)
	EnvironmentGetDiskControlExtInterface       = uint32(C.RETRO_ENVIRONMENT_SET_DISK_CONTROL_EXT_INTERFACE)
)

// Debug levels
const (
	LogLevelDebug = uint32(C.RETRO_LOG_DEBUG)
	LogLevelInfo  = uint32(C.RETRO_LOG_INFO)
	LogLevelWarn  = uint32(C.RETRO_LOG_WARN)
	LogLevelError = uint32(C.RETRO_LOG_ERROR)
	LogLevelDummy = uint32(C.RETRO_LOG_DUMMY)
)

// Memory constants
const (
	MemoryMask      = uint32(C.RETRO_MEMORY_MASK)
	MemorySaveRAM   = uint32(C.RETRO_MEMORY_SAVE_RAM)
	MemoryRTC       = uint32(C.RETRO_MEMORY_RTC)
	MemorySystemRAM = uint32(C.RETRO_MEMORY_SYSTEM_RAM)
	MemoryVideoRAM  = uint32(C.RETRO_MEMORY_VIDEO_RAM)
)

type (
	environmentFunc      func(uint32, unsafe.Pointer) bool
	videoRefreshFunc     func(unsafe.Pointer, int32, int32, int32)
	audioSampleFunc      func(int16, int16)
	audioSampleBatchFunc func([]byte, int32) int32
	inputPollFunc        func()
	inputStateFunc       func(uint, uint32, uint, uint) int16
	logFunc              func(uint32, string)
	getTimeUsecFunc      func() int64
)

var (
	environment      environmentFunc
	videoRefresh     videoRefreshFunc
	audioSample      audioSampleFunc
	audioSampleBatch audioSampleBatchFunc
	inputPoll        inputPollFunc
	inputState       inputStateFunc
	log              logFunc
	getTimeUsec      getTimeUsecFunc
)

// Load dynamically loads a libretro core at the given path and returns a Core instance
func Load(sofile string) (*Core, error) {
	core := Core{}

	var err error
	core.handle, err = DlOpen(sofile)
	if err != nil {
		return nil, err
	}

	core.symRetroInit = core.DlSym("retro_init")
	core.symRetroDeinit = core.DlSym("retro_deinit")
	core.symRetroAPIVersion = core.DlSym("retro_api_version")
	core.symRetroGetSystemInfo = core.DlSym("retro_get_system_info")
	core.symRetroGetSystemAVInfo = core.DlSym("retro_get_system_av_info")
	core.symRetroSetEnvironment = core.DlSym("retro_set_environment")
	core.symRetroSetVideoRefresh = core.DlSym("retro_set_video_refresh")
	core.symRetroSetControllerPortDevice = core.DlSym("retro_set_controller_port_device")
	core.symRetroSetInputPoll = core.DlSym("retro_set_input_poll")
	core.symRetroSetInputState = core.DlSym("retro_set_input_state")
	core.symRetroSetAudioSample = core.DlSym("retro_set_audio_sample")
	core.symRetroSetAudioSampleBatch = core.DlSym("retro_set_audio_sample_batch")
	core.symRetroRun = core.DlSym("retro_run")
	core.symRetroReset = core.DlSym("retro_reset")
	core.symRetroLoadGame = core.DlSym("retro_load_game")
	core.symRetroUnloadGame = core.DlSym("retro_unload_game")
	core.symRetroSerializeSize = core.DlSym("retro_serialize_size")
	core.symRetroSerialize = core.DlSym("retro_serialize")
	core.symRetroUnserialize = core.DlSym("retro_unserialize")
	core.symRetroGetMemorySize = core.DlSym("retro_get_memory_size")
	core.symRetroGetMemoryData = core.DlSym("retro_get_memory_data")

	return &core, nil
}

// Init takes care of the library global initialization
func (core *Core) Init() {
	C.bridge_retro_init(core.symRetroInit)
}

// APIVersion returns the RETRO_API_VERSION.
// Used to validate ABI compatibility when the API is revised.
func (core *Core) APIVersion() uint {
	return uint(C.bridge_retro_api_version(core.symRetroAPIVersion))
}

// Deinit takes care of the library global deinitialization
func (core *Core) Deinit() {
	C.bridge_retro_deinit(core.symRetroDeinit)
	environment = nil
	videoRefresh = nil
	audioSample = nil
	audioSampleBatch = nil
	inputPoll = nil
	inputState = nil
	log = nil
	getTimeUsec = nil
}

// Run runs the game for one video frame.
// During retro_run(), input_poll callback must be called at least once.
// If a frame is not rendered for reasons where a game "dropped" a frame,
// this still counts as a frame, and retro_run() should explicitly dupe
// a frame if GET_CAN_DUPE returns true.
// In this case, the video callback can take a NULL argument for data.
func (core *Core) Run() {
	C.bridge_retro_run(core.symRetroRun)
}

// Reset resets the current game.
func (core *Core) Reset() {
	C.bridge_retro_reset(core.symRetroReset)
}

// GetSystemInfo returns statically known system info. Pointers provided in *info
// must be statically allocated.
// Can be called at any time, even before retro_init().
func (core *Core) GetSystemInfo() SystemInfo {
	rsi := C.struct_retro_system_info{}
	C.bridge_retro_get_system_info(core.symRetroGetSystemInfo, &rsi)
	return SystemInfo{
		LibraryName:     C.GoString(rsi.library_name),
		LibraryVersion:  C.GoString(rsi.library_version),
		ValidExtensions: C.GoString(rsi.valid_extensions),
		NeedFullpath:    bool(rsi.need_fullpath),
		BlockExtract:    bool(rsi.block_extract),
	}
}

// GetSystemAVInfo returns information about system audio/video timings and geometry.
// Can be called only after retro_load_game() has successfully completed.
// NOTE: The implementation of this function might not initialize every
// variable if needed.
// E.g. geom.aspect_ratio might not be initialized if core doesn't
// desire a particular aspect ratio.
func (core *Core) GetSystemAVInfo() SystemAVInfo {
	avi := C.struct_retro_system_av_info{}
	C.bridge_retro_get_system_av_info(core.symRetroGetSystemAVInfo, &avi)
	return SystemAVInfo{
		Geometry: GameGeometry{
			AspectRatio: float64(avi.geometry.aspect_ratio),
			BaseWidth:   int(avi.geometry.base_width),
			BaseHeight:  int(avi.geometry.base_height),
			MaxWidth:    int(avi.geometry.max_width),
			MaxHeight:   int(avi.geometry.max_height),
		},
		Timing: SystemTiming{
			FPS:        float64(avi.timing.fps),
			SampleRate: float64(avi.timing.sample_rate),
		},
	}
}

// LoadGame loads a game
func (core *Core) LoadGame(gi GameInfo) bool {
	rgi := C.struct_retro_game_info{}
	rgi.path = C.CString(gi.Path)
	rgi.size = C.size_t(gi.Size)
	rgi.data = gi.Data
	return bool(C.bridge_retro_load_game(core.symRetroLoadGame, &rgi))
}

// SerializeSize returns the amount of data the implementation requires to serialize
// internal state (save states).
// Between calls to retro_load_game() and retro_unload_game(), the
// returned size is never allowed to be larger than a previous returned
// value, to ensure that the frontend can allocate a save state buffer once.
func (core *Core) SerializeSize() uint {
	return uint(C.bridge_retro_serialize_size(core.symRetroSerializeSize))
}

// Serialize serializes internal state and returns the state as a byte slice.
func (core *Core) Serialize(size uint) ([]byte, error) {
	data := C.malloc(C.size_t(size))
	ok := bool(C.bridge_retro_serialize(core.symRetroSerialize, data, C.size_t(size)))
	if !ok {
		return nil, errors.New("retro_serialize failed")
	}
	bytes := C.GoBytes(data, C.int(size))
	return bytes, nil
}

// Unserialize unserializes internal state from a byte slice.
func (core *Core) Unserialize(bytes []byte, size uint) error {
	if size <= 0 {
		return errors.New("retro_unserialize failed")
	}
	ok := bool(C.bridge_retro_unserialize(core.symRetroUnserialize, unsafe.Pointer(&bytes[0]), C.size_t(size)))
	if !ok {
		return errors.New("retro_unserialize failed")
	}
	return nil
}

// UnloadGame unloads a currently loaded game
func (core *Core) UnloadGame() {
	C.bridge_retro_unload_game(core.symRetroUnloadGame)
}

// SetEnvironment sets the environment callback.
// Must be called before Init
func (core *Core) SetEnvironment(f environmentFunc) {
	environment = f
	C.bridge_retro_set_environment(core.symRetroSetEnvironment, C.coreEnvironment_cgo)
}

// SetVideoRefresh sets the video refresh callback.
// Must be set before the first Run call
func (core *Core) SetVideoRefresh(f videoRefreshFunc) {
	videoRefresh = f
	C.bridge_retro_set_video_refresh(core.symRetroSetVideoRefresh, C.coreVideoRefresh_cgo)
}

// SetAudioSample sets the audio sample callback.
// Must be set before the first Run call
func (core *Core) SetAudioSample(f audioSampleFunc) {
	audioSample = f
	C.bridge_retro_set_audio_sample(core.symRetroSetAudioSample, C.coreAudioSample_cgo)
}

// SetAudioSampleBatch sets the audio sample batch callback.
// Must be set before the first Run call
func (core *Core) SetAudioSampleBatch(f audioSampleBatchFunc) {
	audioSampleBatch = f
	C.bridge_retro_set_audio_sample_batch(core.symRetroSetAudioSampleBatch, C.coreAudioSampleBatch_cgo)
}

// SetInputPoll sets the input poll callback.
// Must be set before the first Run call
func (core *Core) SetInputPoll(f inputPollFunc) {
	inputPoll = f
	C.bridge_retro_set_input_poll(core.symRetroSetInputPoll, C.coreInputPoll_cgo)
}

// SetInputState sets the input state callback.
// Must be set before the first Run call
func (core *Core) SetInputState(f inputStateFunc) {
	inputState = f
	C.bridge_retro_set_input_state(core.symRetroSetInputState, C.coreInputState_cgo)
}

// BindLogCallback binds f to the log callback
func (core *Core) BindLogCallback(data unsafe.Pointer, f logFunc) {
	log = f
	cb := (*C.struct_retro_log_callback)(data)
	cb.log = (C.retro_log_printf_t)(C.coreLog_cgo)
}

// BindPerfCallback binds f to the perf callback get_time_usec
func (core *Core) BindPerfCallback(data unsafe.Pointer, f getTimeUsecFunc) {
	getTimeUsec = f
	cb := (*C.struct_retro_perf_callback)(data)
	cb.get_time_usec = (C.retro_perf_get_time_usec_t)(C.coreGetTimeUsec_cgo)
}

// SetControllerPortDevice sets the device type attached to a controller port
func (core *Core) SetControllerPortDevice(port uint, device uint32) {
	C.bridge_retro_set_controller_port_device(core.symRetroSetControllerPortDevice, C.unsigned(port), C.unsigned(device))
}

//export coreEnvironment
func coreEnvironment(cmd C.unsigned, data unsafe.Pointer) bool {
	if environment == nil {
		return false
	}
	return environment(uint32(cmd), data)
}

//export coreVideoRefresh
func coreVideoRefresh(data unsafe.Pointer, width C.unsigned, height C.unsigned, pitch C.size_t) {
	if videoRefresh == nil {
		return
	}
	videoRefresh(data, int32(width), int32(height), int32(pitch))
}

//export coreInputPoll
func coreInputPoll() {
	if inputPoll == nil {
		return
	}
	inputPoll()
}

//export coreInputState
func coreInputState(port C.unsigned, device C.unsigned, index C.unsigned, id C.unsigned) C.int16_t {
	if inputState == nil {
		return 0
	}
	return C.int16_t(inputState(uint(port), uint32(device), uint(index), uint(id)))
}

//export coreAudioSample
func coreAudioSample(left C.int16_t, right C.int16_t) {
	if audioSample == nil {
		return
	}
	audioSample(int16(left), int16(right))
}

//export coreAudioSampleBatch
func coreAudioSampleBatch(buf unsafe.Pointer, frames C.size_t) C.size_t {
	if audioSampleBatch == nil {
		return 0
	}
	return C.size_t(audioSampleBatch(C.GoBytes(buf, C.int(4*int(frames))), int32(frames))) / 4
}

//export coreLog
func coreLog(level C.enum_retro_log_level, msg *C.char) {
	if log == nil {
		return
	}
	log(level, C.GoString(msg))
}

//export coreGetTimeUsec
func coreGetTimeUsec() C.uint64_t {
	if getTimeUsec == nil {
		return 0
	}
	return C.uint64_t(getTimeUsec())
}

// SetData is a setter for the data of a GameInfo type
func (gi *GameInfo) SetData(bytes []byte) {
	cstr := C.CString(string(bytes))
	gi.Data = unsafe.Pointer(cstr)
}

// Environment helpers

// GetPixelFormat is an environment callback helper that returns the pixel format.
// Should be used in the case of EnvironmentSetPixelFormat
func GetPixelFormat(data unsafe.Pointer) uint32 {
	return *(*C.enum_retro_pixel_format)(data)
}

// GetVariable is an environment callback helper that returns a Variable
func GetVariable(data unsafe.Pointer) *Variable {
	return (*Variable)(data)
}

// GetVariables is an environment callback helper that returns the list of Variable needed by a core
func GetVariables(data unsafe.Pointer) []Variable {
	var vars []Variable

	for {
		v := (*C.struct_retro_variable)(data)
		if v.key == nil || v.value == nil {
			break
		}
		vars = append(vars, *(*Variable)(v))
		data = unsafe.Pointer(uintptr(data) + unsafe.Sizeof(v.key) + unsafe.Sizeof(v.value))
	}

	return vars
}

// GetCoreOptionDefinitions is an environment callback helper that returns the list of CoreOptionDefinition needed by a core
func GetCoreOptionDefinitions(data unsafe.Pointer) []CoreOptionDefinition {
	var definitions []CoreOptionDefinition

	for {
		v := (*C.struct_retro_core_option_definition)(data)
		if v.key == nil {
			break
		}
		definitions = append(definitions, *(*CoreOptionDefinition)(v))
		data = unsafe.Pointer(uintptr(data) +
			unsafe.Sizeof(v.key) +
			unsafe.Sizeof(v.desc) +
			unsafe.Sizeof(v.info) +
			unsafe.Sizeof(v.values) +
			unsafe.Sizeof(v.default_value))
	}

	return definitions
}

// GetCoreOptionsIntl is an environment callback helper that returns the list of CoreOptionsIntl needed by a core
func GetCoreOptionsIntl(data unsafe.Pointer) []CoreOptionDefinition {
	var definitions []CoreOptionDefinition

	intl := (*C.struct_retro_core_options_intl)(data)
	uuss := unsafe.Pointer(intl.us)
	for {
		v := (*C.struct_retro_core_option_definition)(uuss)
		if v.key == nil {
			break
		}
		definitions = append(definitions, *(*CoreOptionDefinition)(v))
		uuss = unsafe.Pointer(uintptr(uuss) +
			unsafe.Sizeof(v.key) +
			unsafe.Sizeof(v.desc) +
			unsafe.Sizeof(v.info) +
			unsafe.Sizeof(v.values) +
			unsafe.Sizeof(v.default_value))
	}

	return definitions
}

// GetGeometry is an environment callback helper that returns the game geometry
// in EnvironmentSetGeometry.
func GetGeometry(data unsafe.Pointer) GameGeometry {
	geometry := (*C.struct_retro_game_geometry)(data)
	return GameGeometry{
		AspectRatio: float64(geometry.aspect_ratio),
		BaseWidth:   int(geometry.base_width),
		BaseHeight:  int(geometry.base_height),
		MaxWidth:    int(geometry.max_width),
		MaxHeight:   int(geometry.max_height),
	}
}

// GetSystemAVInfo is an environment callback helper that returns the game geometry
// in EnvironmentSetGeometry.
func GetSystemAVInfo(data unsafe.Pointer) SystemAVInfo {
	avi := (*C.struct_retro_system_av_info)(data)
	return SystemAVInfo{
		Geometry: GameGeometry{
			AspectRatio: float64(avi.geometry.aspect_ratio),
			BaseWidth:   int(avi.geometry.base_width),
			BaseHeight:  int(avi.geometry.base_height),
			MaxWidth:    int(avi.geometry.max_width),
			MaxHeight:   int(avi.geometry.max_height),
		},
		Timing: SystemTiming{
			FPS:        float64(avi.timing.fps),
			SampleRate: float64(avi.timing.sample_rate),
		},
	}
}

// SetBool is an environment callback helper to set a boolean
func SetBool(data unsafe.Pointer, val bool) {
	b := (*C.bool)(data)
	*b = C.bool(val)
}

// SetString is an environment callback helper to set a string
func SetString(data unsafe.Pointer, val string) {
	s := (**C.char)(data)
	*s = C.CString(val)
}

// SetUint is an environment callback helper to set a string
func SetUint(data unsafe.Pointer, val uint) {
	i := (*C.uint)(data)
	*i = C.uint(val)
}

// SetFrameTimeCallback is an environment callback helper to set the FrameTimeCallback
func (core *Core) SetFrameTimeCallback(data unsafe.Pointer) {
	c := *(*C.struct_retro_frame_time_callback)(data)
	ftc := &FrameTimeCallback{}
	ftc.Reference = int64(c.reference)
	ftc.Callback = func(usec int64) {
		C.bridge_retro_frame_time_callback(c.callback, C.retro_usec_t(usec))
	}
	core.FrameTimeCallback = ftc
}

// SetAudioCallback is an environment callback helper to set the AudioCallback
func (core *Core) SetAudioCallback(data unsafe.Pointer) {
	c := *(*C.struct_retro_audio_callback)(data)
	auc := &AudioCallback{}
	auc.Callback = func() {
		C.bridge_retro_audio_callback(c.callback)
	}
	auc.SetState = func(state bool) {
		C.bridge_retro_audio_set_state(c.set_state, C.bool(state))
	}
	core.AudioCallback = auc
}

// GetMemorySize returns the size of a region of the memory.
// See memory constants.
func (core *Core) GetMemorySize(id uint32) uint {
	return uint(C.bridge_retro_get_memory_size(core.symRetroGetMemorySize, C.unsigned(id)))
}

// GetMemoryData returns the size of a region of the memory.
// See memory constants.
func (core *Core) GetMemoryData(id uint32) unsafe.Pointer {
	return C.bridge_retro_get_memory_data(core.symRetroGetMemoryData, C.unsigned(id))
}
