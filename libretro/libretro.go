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
void bridge_retro_hw_context_reset(retro_hw_context_reset_t f);
uintptr_t bridge_retro_hw_get_current_framebuffer(retro_hw_get_current_framebuffer_t f);
void bridge_retro_hw_context_destroy(retro_hw_context_reset_t f);
void bridge_retro_audio_callback(retro_audio_callback_t f);
void bridge_retro_audio_set_state(retro_audio_set_state_callback_t f, bool state);

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
	"sync"
	"unsafe"
)

// GameGeometry represents the geometry of a game, with its aspect ratio, with and height
type GameGeometry struct {
	AspectRatio float64
	BaseWidth   int
	BaseHeight  int
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
	s := (**C.char)(&v.value)
	*s = C.CString(val)
}

// FrameTimeCallback stores the frame time callback itself and the reference time
type FrameTimeCallback struct {
	Callback  func(int64)
	Reference int64
}

// HWRenderCallback sets an interface to let a libretro core render with
// hardware acceleration.
type HWRenderCallback struct {
	HWContextType              uint
	ContextReset               func()
	GetCurrentFramebuffer      func() uintptr
	GetProcAddress             func()
	Depth                      bool
	Stencil                    bool
	BottomLeftOrigin           bool
	VersionMajor, VersionMinor uint
	CacheContext               bool
	ContextDestroy             func()
	DebugContext               bool
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

// DeviceJoypad represents the RetroPad. It is essentially a Super Nintendo
// controller, but with additional L2/R2/L3/R3 buttons, similar to a
// PS1 DualShock.
const DeviceJoypad = uint32(C.RETRO_DEVICE_JOYPAD)

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

// Environment callback API. See libretro.h for details
const (
	EnvironmentGetUsername          = uint32(C.RETRO_ENVIRONMENT_GET_USERNAME)
	EnvironmentGetLogInterface      = uint32(C.RETRO_ENVIRONMENT_GET_LOG_INTERFACE)
	EnvironmentGetCanDupe           = uint32(C.RETRO_ENVIRONMENT_GET_CAN_DUPE)
	EnvironmentSetPixelFormat       = uint32(C.RETRO_ENVIRONMENT_SET_PIXEL_FORMAT)
	EnvironmentGetSystemDirectory   = uint32(C.RETRO_ENVIRONMENT_GET_SYSTEM_DIRECTORY)
	EnvironmentGetSaveDirectory     = uint32(C.RETRO_ENVIRONMENT_GET_SAVE_DIRECTORY)
	EnvironmentShutdown             = uint32(C.RETRO_ENVIRONMENT_SHUTDOWN)
	EnvironmentGetVariable          = uint32(C.RETRO_ENVIRONMENT_GET_VARIABLE)
	EnvironmentSetVariables         = uint32(C.RETRO_ENVIRONMENT_SET_VARIABLES)
	EnvironmentGetVariableUpdate    = uint32(C.RETRO_ENVIRONMENT_GET_VARIABLE_UPDATE)
	EnvironmentGetPerfInterface     = uint32(C.RETRO_ENVIRONMENT_GET_PERF_INTERFACE)
	EnvironmentSetFrameTimeCallback = uint32(C.RETRO_ENVIRONMENT_SET_FRAME_TIME_CALLBACK)
	EnvironmentSetAudioCallback     = uint32(C.RETRO_ENVIRONMENT_SET_AUDIO_CALLBACK)
	EnvironmentSetHWRenderer        = uint32(C.RETRO_ENVIRONMENT_SET_HW_RENDER)
)

// Debug levels
const (
	LogLevelDebug = uint32(C.RETRO_LOG_DEBUG)
	LogLevelInfo  = uint32(C.RETRO_LOG_INFO)
	LogLevelWarn  = uint32(C.RETRO_LOG_WARN)
	LogLevelError = uint32(C.RETRO_LOG_ERROR)
	LogLevelDummy = uint32(C.RETRO_LOG_DUMMY)
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

var mu sync.Mutex

// Load dynamically loads a libretro core at the given path and returns a Core instance
func Load(sofile string) (Core, error) {
	core := Core{}

	mu.Lock()
	var err error
	core.handle, err = DlOpen(sofile)
	if err != nil {
		return core, err
	}

	core.symRetroInit = core.DlSym("retro_init")
	core.symRetroDeinit = core.DlSym("retro_deinit")
	core.symRetroAPIVersion = core.DlSym("retro_api_version")
	core.symRetroGetSystemInfo = core.DlSym("retro_get_system_info")
	core.symRetroGetSystemAVInfo = core.DlSym("retro_get_system_av_info")
	core.symRetroSetEnvironment = core.DlSym("retro_set_environment")
	core.symRetroSetVideoRefresh = core.DlSym("retro_set_video_refresh")
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
	mu.Unlock()

	return core, nil
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
	return C.size_t(audioSampleBatch(C.GoBytes(buf, C.int(4096)), int32(frames)))
}

//export coreLog
func coreLog(level C.enum_retro_log_level, msg *C.char) {
	log(uint32(level), C.GoString(msg))
}

//export coreGetTimeUsec
func coreGetTimeUsec() C.uint64_t {
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

// GetVariables is an environment callback helper that returns the list of Variables needed by a core
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

// SetFrameTimeCallback is an environment callback helper to set the FrameTimeCallback
func SetFrameTimeCallback(data unsafe.Pointer) FrameTimeCallback {
	c := *(*C.struct_retro_frame_time_callback)(data)
	ftc := FrameTimeCallback{}
	ftc.Reference = int64(c.reference)
	ftc.Callback = func(usec int64) {
		C.bridge_retro_frame_time_callback(c.callback, C.retro_usec_t(usec))
	}
	return ftc
}

// SetHWRenderCallback is an environment callback helper to set the HWRenderCallback
func SetHWRenderCallback(data unsafe.Pointer) HWRenderCallback {
	c := *(*C.struct_retro_hw_render_callback)(data)
	hwrc := HWRenderCallback{}
	hwrc.HWContextType = uint(c.context_type)
	hwrc.ContextReset = func() {
		C.bridge_retro_hw_context_reset(c.context_reset)
	}
	hwrc.GetCurrentFramebuffer = func() uintptr {
		return uintptr(C.bridge_retro_hw_get_current_framebuffer(c.get_current_framebuffer))
	}
	hwrc.Depth = bool(c.depth)
	hwrc.Stencil = bool(c.stencil)
	hwrc.BottomLeftOrigin = bool(c.bottom_left_origin)
	hwrc.VersionMajor = uint(c.version_major)
	hwrc.VersionMinor = uint(c.version_minor)
	hwrc.CacheContext = bool(c.cache_context)
	hwrc.ContextDestroy = func() {
		C.bridge_retro_hw_context_destroy(c.context_destroy)
	}
	hwrc.DebugContext = bool(c.debug_context)
	return hwrc
}

// SetAudioCallback is an environment callback helper to set the AudioCallback
func SetAudioCallback(data unsafe.Pointer) AudioCallback {
	c := *(*C.struct_retro_audio_callback)(data)
	auc := AudioCallback{}
	auc.Callback = func() {
		C.bridge_retro_audio_callback(c.callback)
	}
	auc.SetState = func(state bool) {
		C.bridge_retro_audio_set_state(c.set_state, C.bool(state))
	}
	return auc
}
