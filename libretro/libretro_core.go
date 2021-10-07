package libretro // Core is an instance of a dynamically loaded libretro core

import "unsafe"

// Core is an instance of a dynamically loaded libretro core
type Core struct {
	handle DlHandle

	symRetroInit                    unsafe.Pointer
	symRetroDeinit                  unsafe.Pointer
	symRetroAPIVersion              unsafe.Pointer
	symRetroGetSystemInfo           unsafe.Pointer
	symRetroGetSystemAVInfo         unsafe.Pointer
	symRetroSetEnvironment          unsafe.Pointer
	symRetroSetVideoRefresh         unsafe.Pointer
	symRetroSetControllerPortDevice unsafe.Pointer
	symRetroSetInputPoll            unsafe.Pointer
	symRetroSetInputState           unsafe.Pointer
	symRetroSetAudioSample          unsafe.Pointer
	symRetroSetAudioSampleBatch     unsafe.Pointer
	symRetroRun                     unsafe.Pointer
	symRetroReset                   unsafe.Pointer
	symRetroLoadGame                unsafe.Pointer
	symRetroUnloadGame              unsafe.Pointer
	symRetroSerializeSize           unsafe.Pointer
	symRetroSerialize               unsafe.Pointer
	symRetroUnserialize             unsafe.Pointer
	symRetroGetMemorySize           unsafe.Pointer
	symRetroGetMemoryData           unsafe.Pointer

	AudioCallback       *AudioCallback
	FrameTimeCallback   *FrameTimeCallback
	DiskControlCallback *DiskControlCallback

	MemoryMap []MemoryDescriptor
}
