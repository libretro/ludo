package main

/*
#include <stdbool.h>

bool coreEnvironment_cgo(unsigned cmd, void *data) {
	bool coreEnvironment(unsigned, void*);
	return coreEnvironment(cmd, data);
}
*/
import "C"
