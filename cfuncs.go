package main

/*
#include "libretro.h"
#include <stdbool.h>

void bridge_retro_init(void *f) {
	return ((void (*)(void))f)();
}

void bridge_retro_deinit(void *f) {
	return ((void (*)(void))f)();
}

unsigned bridge_retro_api_version(void *f) {
	return ((unsigned (*)(void))f)();
}

bool bridge_retro_set_environment(void *f, void *callback) {
	return ((bool (*)(retro_environment_t))f)((retro_environment_t)callback);
}

void bridge_retro_set_video_refresh(void *f, void *callback) {
	((bool (*)(retro_video_refresh_t))f)((retro_video_refresh_t)callback);
}

void bridge_retro_set_input_poll(void *f, void *callback) {
	((bool (*)(retro_input_poll_t))f)((retro_input_poll_t)callback);
}

bool bridge_retro_load_game(void *f, struct retro_game_info *gi) {
  return ((bool (*)(struct retro_game_info *))f)(gi);
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

*/
import "C"
