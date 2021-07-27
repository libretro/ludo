package libretro

/*
#include "libretro.h"
#include <stdbool.h>
#include <stdarg.h>
#include <stdio.h>
#include <pthread.h>

#ifdef __APPLE__
#include <mach/semaphore.h>
#include <mach/task.h>
#include <mach/mach_init.h>
#define SEM_T semaphore_t
#define SEM_INIT(x) semaphore_create(mach_task_self(), &x, SYNC_POLICY_FIFO, 0);
#define SEM_POST(x) semaphore_signal(x)
#define SEM_WAIT(x) semaphore_wait(x)
#else
#include <semaphore.h>
#define SEM_T sem_t
#define SEM_INIT(x) sem_init(&x, 0, 0)
#define SEM_POST(x) sem_post(&x)
#define SEM_WAIT(x) sem_wait(&x)
#endif

#if 1
#define print_sema(...) (printf(__VA_ARGS__))
#else
#define print_sema(...) do {} while (0)
#endif

enum {
	CMD_F,
	CMD_SERIALIZE,
};

struct thread_cmd_t {
	int   cmd;
	void* f;
	void* arg1;
	void* arg2;
	void* arg3;
	void* arg4;
	void* res;
};

struct thread_cmd_t s_job;
static pthread_t s_thread;
static SEM_T s_sem_do;
static SEM_T s_sem_done;
static bool s_use_thread = false;

void* emu_thread_loop(void *a0) {
	print_sema("begin thread\n");

	SEM_POST(s_sem_done);

	print_sema("signal thread\n");

	while (1) {
		print_sema("wait do\n");
		SEM_WAIT(s_sem_do);

		print_sema("do\n");
		switch (s_job.cmd) {
		case CMD_F:
			((void (*)(void))s_job.f)();
			break;
		case CMD_SERIALIZE: {
			bool res = ((bool (*)(void*, size_t))s_job.f)(s_job.arg1, *(size_t*)s_job.arg2);
			*(bool*)s_job.res = res;
			break;
		}
		default:
			break;
		}

		print_sema("signal done\n");
		SEM_POST(s_sem_done);
	}
}

void thread_sync() {
	// Fire the job
	print_sema("signal do\n");
	SEM_POST(s_sem_do);

	// Wait the result
	print_sema("wait done\n");
	SEM_WAIT(s_sem_done);

	print_sema("done\n");
}

void run_wrapper(void *f) {
	if (s_use_thread) {
		s_job.cmd = CMD_F;
		s_job.f = f;
		thread_sync();
	} else {
		((void (*)(void))f)();
	}
}

void cothread_init() {
	s_use_thread = true;

	SEM_INIT(s_sem_do);
	SEM_INIT(s_sem_done);

	print_sema("create thread\n");
	pthread_create(&s_thread, NULL, emu_thread_loop, NULL);

	print_sema("wait thread\n");
	SEM_WAIT(s_sem_done);
}

void bridge_retro_init(void *f) {
	run_wrapper(f);
}

void bridge_retro_deinit(void *f) {
	run_wrapper(f);
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
	if (s_use_thread) {
		bool res;
		s_job.cmd = CMD_SERIALIZE;
		s_job.f = f;
		s_job.arg1 = data;
		s_job.arg2 = &size;
		s_job.res  = &res;

		thread_sync();

		return s_job.res;
	} else {
		return ((bool (*)(void*, size_t))f)(data, size);
	}
}

bool bridge_retro_unserialize(void *f, void *data, size_t size) {
	if (s_use_thread) {
		bool res;
		s_job.cmd = CMD_SERIALIZE; // Same command format for both serialize & unserialize
		s_job.f = f;
		s_job.arg1 = data;
		s_job.arg2 = &size;
		s_job.res  = &res;

		thread_sync();

		return s_job.res;
	} else {
		return ((bool (*)(void*, size_t))f)(data, size);
	}
}

void bridge_retro_unload_game(void *f) {
	run_wrapper(f);
}

void bridge_retro_run(void *f) {
	run_wrapper(f);
}

void bridge_retro_reset(void *f) {
	run_wrapper(f);
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
