package libretro

/*
#cgo LDFLAGS: -ldl
#include "libretro.h"
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <dlfcn.h>

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
void bridge_retro_unload_game(void *f);
void bridge_retro_run(void *f);

bool coreEnvironment_cgo(unsigned cmd, void *data);
void coreVideoRefresh_cgo(void *data, unsigned width, unsigned height, size_t pitch);
void coreInputPoll_cgo();
void coreAudioSample_cgo(int16_t left, int16_t right);
size_t coreAudioSampleBatch_cgo(const int16_t *data, size_t frames);
int16_t coreInputState_cgo(unsigned port, unsigned device, unsigned index, unsigned id);
void coreLog_cgo(enum retro_log_level level, const char *msg);
*/
import "C"
import (
	"errors"
	"sync"
	"unsafe"
)

type Core struct {
	handle unsafe.Pointer

	symRetroInit                unsafe.Pointer
	symRetroDeinit              unsafe.Pointer
	symRetroAPIVersion          unsafe.Pointer
	symRetroGetSystemInfo       unsafe.Pointer
	symRetroGetSystemAVInfo     unsafe.Pointer
	symRetroSetEnvironment      unsafe.Pointer
	symRetroSetVideoRefresh     unsafe.Pointer
	symRetroSetInputPoll        unsafe.Pointer
	symRetroSetInputState       unsafe.Pointer
	symRetroSetAudioSample      unsafe.Pointer
	symRetroSetAudioSampleBatch unsafe.Pointer
	symRetroRun                 unsafe.Pointer
	symRetroLoadGame            unsafe.Pointer
	symRetroUnloadGame          unsafe.Pointer

	videoRefresh videoRefreshFunc
}

type GameGeometry struct {
	AspectRatio float64
	BaseWidth   int
	BaseHeight  int
}

type GameInfo struct {
	Path string
	Size int64
	Data unsafe.Pointer
}

type SystemInfo struct {
	LibraryName     string
	LibraryVersion  string
	ValidExtensions string
	NeedFullpath    bool
	BlockExtract    bool
}

type SystemTiming struct {
	FPS        float64
	SampleRate float64
}

type SystemAVInfo struct {
	Geometry GameGeometry
	Timing   SystemTiming
}

type Variable struct {
	Key   string
	Value string
}

const (
	PixelFormat0RGB1555 = uint32(C.RETRO_PIXEL_FORMAT_0RGB1555)
	PixelFormatXRGB8888 = uint32(C.RETRO_PIXEL_FORMAT_XRGB8888)
	PixelFormatRGB565   = uint32(C.RETRO_PIXEL_FORMAT_RGB565)
)

const (
	DeviceJoypad = uint32(C.RETRO_DEVICE_JOYPAD)

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

const (
	EnvironmentGetUsername        = uint32(C.RETRO_ENVIRONMENT_GET_USERNAME)
	EnvironmentGetLogInterface    = uint32(C.RETRO_ENVIRONMENT_GET_LOG_INTERFACE)
	EnvironmentGetCanDupe         = uint32(C.RETRO_ENVIRONMENT_GET_CAN_DUPE)
	EnvironmentSetPixelFormat     = uint32(C.RETRO_ENVIRONMENT_SET_PIXEL_FORMAT)
	EnvironmentGetSystemDirectory = uint32(C.RETRO_ENVIRONMENT_GET_SYSTEM_DIRECTORY)
	EnvironmentGetSaveDirectory   = uint32(C.RETRO_ENVIRONMENT_GET_SAVE_DIRECTORY)
	EnvironmentShutdown           = uint32(C.RETRO_ENVIRONMENT_SHUTDOWN)
	EnvironmentGetVariable        = uint32(C.RETRO_ENVIRONMENT_GET_VARIABLE)
)

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
)

var (
	environment      environmentFunc
	videoRefresh     videoRefreshFunc
	audioSample      audioSampleFunc
	audioSampleBatch audioSampleBatchFunc
	inputPoll        inputPollFunc
	inputState       inputStateFunc
	log              logFunc
)

var mu sync.Mutex

func Load(sofile string) (Core, error) {
	core := Core{}

	mu.Lock()
	core.handle = C.dlopen(C.CString(sofile), C.RTLD_NOW)
	if core.handle == nil {
		return core, errors.New("dlopen failed")
	}

	core.symRetroInit = C.dlsym(core.handle, C.CString("retro_init"))
	core.symRetroDeinit = C.dlsym(core.handle, C.CString("retro_deinit"))
	core.symRetroAPIVersion = C.dlsym(core.handle, C.CString("retro_api_version"))
	core.symRetroGetSystemInfo = C.dlsym(core.handle, C.CString("retro_get_system_info"))
	core.symRetroGetSystemAVInfo = C.dlsym(core.handle, C.CString("retro_get_system_av_info"))
	core.symRetroSetEnvironment = C.dlsym(core.handle, C.CString("retro_set_environment"))
	core.symRetroSetVideoRefresh = C.dlsym(core.handle, C.CString("retro_set_video_refresh"))
	core.symRetroSetInputPoll = C.dlsym(core.handle, C.CString("retro_set_input_poll"))
	core.symRetroSetInputState = C.dlsym(core.handle, C.CString("retro_set_input_state"))
	core.symRetroSetAudioSample = C.dlsym(core.handle, C.CString("retro_set_audio_sample"))
	core.symRetroSetAudioSampleBatch = C.dlsym(core.handle, C.CString("retro_set_audio_sample_batch"))
	core.symRetroRun = C.dlsym(core.handle, C.CString("retro_run"))
	core.symRetroLoadGame = C.dlsym(core.handle, C.CString("retro_load_game"))
	core.symRetroUnloadGame = C.dlsym(core.handle, C.CString("retro_unload_game"))
	mu.Unlock()

	C.bridge_retro_set_environment(core.symRetroSetEnvironment, C.coreEnvironment_cgo)
	C.bridge_retro_set_video_refresh(core.symRetroSetVideoRefresh, C.coreVideoRefresh_cgo)
	C.bridge_retro_set_input_poll(core.symRetroSetInputPoll, C.coreInputPoll_cgo)
	C.bridge_retro_set_input_state(core.symRetroSetInputState, C.coreInputState_cgo)
	C.bridge_retro_set_audio_sample(core.symRetroSetAudioSample, C.coreAudioSample_cgo)
	C.bridge_retro_set_audio_sample_batch(core.symRetroSetAudioSampleBatch, C.coreAudioSampleBatch_cgo)

	return core, nil
}

func (core *Core) Init() {
	C.bridge_retro_init(core.symRetroInit)
}

func (core *Core) APIVersion() uint {
	return uint(C.bridge_retro_api_version(core.symRetroAPIVersion))
}

func (core *Core) Deinit() {
	C.bridge_retro_deinit(core.symRetroDeinit)
}

func (core *Core) Run() {
	C.bridge_retro_run(core.symRetroRun)
}

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

func (core *Core) LoadGame(gi GameInfo) bool {
	rgi := C.struct_retro_game_info{}
	rgi.path = C.CString(gi.Path)
	rgi.size = C.size_t(gi.Size)
	rgi.data = gi.Data
	return bool(C.bridge_retro_load_game(core.symRetroLoadGame, &rgi))
}

func (core *Core) UnloadGame() {
	C.bridge_retro_unload_game(core.symRetroUnloadGame)
}

func (core *Core) SetEnvironment(f environmentFunc) {
	environment = f
}

func (core *Core) SetVideoRefresh(f videoRefreshFunc) {
	videoRefresh = f
}

func (core *Core) SetAudioSample(f audioSampleFunc) {
	audioSample = f
}

func (core *Core) SetAudioSampleBatch(f audioSampleBatchFunc) {
	audioSampleBatch = f
}

func (core *Core) SetInputPoll(f inputPollFunc) {
	inputPoll = f
}

func (core *Core) SetInputState(f inputStateFunc) {
	inputState = f
}

func (core *Core) BindLogCallback(data unsafe.Pointer, f logFunc) {
	log = f
	cb := (*C.struct_retro_log_callback)(data)
	cb.log = (C.retro_log_printf_t)(C.coreLog_cgo)
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

func (gi *GameInfo) SetData(bytes []byte) {
	cstr := C.CString(string(bytes))
	gi.Data = unsafe.Pointer(cstr)
}

// Environment helpers

func GetPixelFormat(data unsafe.Pointer) uint32 {
	return *(*C.enum_retro_pixel_format)(data)
}

func GetVariable(data unsafe.Pointer) Variable {
	variable := (*C.struct_retro_variable)(data)
	return Variable{
		Key:   C.GoString(variable.key),
		Value: C.GoString(variable.value),
	}
}

func SetBool(data unsafe.Pointer, val bool) {
	b := (*C.bool)(data)
	*b = C.bool(val)
}

func SetString(data unsafe.Pointer, val string) {
	s := (**C.char)(data)
	*s = C.CString(val)
}
