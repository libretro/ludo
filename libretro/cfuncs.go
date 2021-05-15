package libretro

/*
#include "libretro.h"
#include <stdbool.h>
#include <stdarg.h>
#include <stdio.h>

void bridge_retro_init(void *f) {
	return ((void (*)(void))f)();
}

void bridge_retro_deinit(void *f) {
	return ((void (*)(void))f)();
}

unsigned bridge_retro_api_version(void *f) {
	return ((unsigned (*)(void))f)();
}

void bridge_retro_frame_time_callback(retro_frame_time_callback_t f, retro_usec_t usec) {
	f(usec);
}

void bridge_retro_audio_callback(retro_audio_callback_t f) {
	f();
}

void bridge_retro_audio_set_state(retro_audio_set_state_callback_t f, bool state) {
	f(state);
}

void bridge_retro_get_system_info(void *f, struct retro_system_info *si) {
  return ((void (*)(struct retro_system_info *))f)(si);
}

void bridge_retro_get_system_av_info(void *f, struct retro_system_av_info *si) {
  return ((void (*)(struct retro_system_av_info *))f)(si);
}

bool bridge_retro_set_environment(void *f, void *callback) {
	return ((bool (*)(retro_environment_t))f)((retro_environment_t)callback);
}

void bridge_retro_set_video_refresh(void *f, void *callback) {
	((bool (*)(retro_video_refresh_t))f)((retro_video_refresh_t)callback);
}

void bridge_retro_set_controller_port_device(void *f, unsigned port, unsigned device) {
	return ((void (*)(unsigned, unsigned))f)(port, device);
}

void bridge_retro_set_input_poll(void *f, void *callback) {
	((bool (*)(retro_input_poll_t))f)((retro_input_poll_t)callback);
}

void bridge_retro_set_input_state(void *f, void *callback) {
	((bool (*)(retro_input_state_t))f)((retro_input_state_t)callback);
}

void bridge_retro_set_audio_sample(void *f, void *callback) {
	((bool (*)(retro_audio_sample_t))f)((retro_audio_sample_t)callback);
}

void bridge_retro_set_audio_sample_batch(void *f, void *callback) {
	((bool (*)(retro_audio_sample_batch_t))f)((retro_audio_sample_batch_t)callback);
}

bool bridge_retro_load_game(void *f, struct retro_game_info *gi) {
  return ((bool (*)(struct retro_game_info *))f)(gi);
}

size_t bridge_retro_serialize_size(void *f) {
  return ((size_t (*)(void))f)();
}

bool bridge_retro_serialize(void *f, void *data, size_t size) {
  return ((bool (*)(void*, size_t))f)(data, size);
}

bool bridge_retro_unserialize(void *f, void *data, size_t size) {
  return ((bool (*)(void*, size_t))f)(data, size);
}

void bridge_retro_unload_game(void *f) {
	return ((void (*)(void))f)();
}

void bridge_retro_run(void *f) {
	return ((void (*)(void))f)();
}

void bridge_retro_reset(void *f) {
	return ((void (*)(void))f)();
}

size_t bridge_retro_get_memory_size(void *f, unsigned id) {
	return ((size_t (*)(unsigned))f)(id);
}

void* bridge_retro_get_memory_data(void *f, unsigned id) {
	return ((void* (*)(unsigned))f)(id);
}

void bridge_retro_set_eject_state(retro_set_eject_state_t f, bool state) {
	f(state);
}

bool bridge_retro_get_eject_state(retro_get_eject_state_t f) {
	return ((bool (*)())f)();
}

unsigned bridge_retro_get_image_index(retro_get_image_index_t f) {
	return ((unsigned (*)())f)();
}

void bridge_retro_set_image_index(retro_set_image_index_t f, unsigned index) {
	f(index);
}

unsigned bridge_retro_get_num_images(retro_get_num_images_t f) {
	return ((unsigned (*)())f)();
}

bool coreEnvironment_cgo(unsigned cmd, void *data) {
	bool coreEnvironment(unsigned, void*);
	return coreEnvironment(cmd, data);
}

void coreVideoRefresh_cgo(void *data, unsigned width, unsigned height, size_t pitch) {
	void coreVideoRefresh(void*, unsigned, unsigned, size_t);
	return coreVideoRefresh(data, width, height, pitch);
}

void coreInputPoll_cgo() {
	void coreInputPoll();
	return coreInputPoll();
}

int16_t coreInputState_cgo(unsigned port, unsigned device, unsigned index, unsigned id) {
	int16_t coreInputState(unsigned, unsigned, unsigned, unsigned);
	return coreInputState(port, device, index, id);
}

void coreAudioSample_cgo(int16_t left, int16_t right) {
	void coreAudioSample(int16_t, int16_t);
	coreAudioSample(left, right);
}

size_t coreAudioSampleBatch_cgo(const int16_t *data, size_t frames) {
	size_t coreAudioSampleBatch(const int16_t*, size_t);
	return coreAudioSampleBatch(data, frames);
}

void coreLog_cgo(enum retro_log_level level, const char *fmt, ...) {
	char msg[4096] = {0};
	va_list va;
	va_start(va, fmt);
	vsnprintf(msg, sizeof(msg), fmt, va);
	va_end(va);

	void coreLog(enum retro_log_level level, const char *msg);
	coreLog(level, msg);
}

int64_t coreGetTimeUsec_cgo() {
	uint64_t coreGetTimeUsec();
	return coreGetTimeUsec();
}

*/
import "C"
