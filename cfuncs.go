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

bool bridge_retro_load_game(void *f, struct retro_game_info *gi) {
  return ((bool (*)(struct retro_game_info *))f)(gi);
}

bool coreEnvironment_cgo(unsigned cmd, void *data) {
	bool coreEnvironment(unsigned, void*);
	return coreEnvironment(cmd, data);
}
*/
import "C"
