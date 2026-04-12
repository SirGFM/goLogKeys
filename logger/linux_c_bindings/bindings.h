#ifndef BINDINGS_H
#define BINDINGS_H

#include <stdint.h>

enum status {
	ST_RELEASED = 0
	, ST_PRESSED
};

struct event {
	uint16_t key;
	enum status status;
};

/**
 * new_context initializes a new context for monitoring keyboard events.
 *
 * buf_size indicates how many items may be read in a single read().
 *
 * On success, a valid pointer is returned. Otherwise, NULL is returned.
 */
extern void* new_context(int buf_size);

/**
 * close_context releases the context and any resources associated with it.
 */
extern void close_context(void *c);

/**
 * open_keyboard_fd opens a file descriptor for a keyboard device.
 *
 * The file descriptor is opened for read-only in nonblocking mode
 * and is added to the context.
 *
 * On success, 0 is returned. Otherwise, -1 is returned.
 */
extern int open_keyboard_fd(void *c, char *device);

/**
 * get_num_fd retrieves the number of file descriptors that are open.
 */
extern int get_num_fd(void *c);

/**
 * wait_events blocks until any file descriptor is ready to be read.
 *
 * On failure, -1 is returned. Otherwise,
 * the number of file descriptors with data is returned.
 */
extern int wait_events(void *c);

/**
 * get_ready retrieves the number of file descriptors that are ready to be read.
 */
extern int get_ready(void *c);

/**
 * read_events read events from the specified file descriptors if it's ready.
 *
 * On failure, -1 is returned.
 * Otherwise, the number of read items is returned.
 * A return of 0 indicates that this source wasn't ready.
 *
 * If any event is ready, events shall be filled with an internal buffer
 * with all events, which shall last until the next read_events() call.
 */
extern int read_events(void *c, int idx, struct event **events);

#endif /* BINDINGS_H */
