#include "runtime.h"

int64 ·Id(void) {
	return g->goid;
}
