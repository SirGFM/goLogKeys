#include <errno.h>
#include <fcntl.h>
#include <linux/input.h>
#include <poll.h>
#include <unistd.h>
#include <sys/ioctl.h>

#include <stddef.h>
#include <stdlib.h>
#include <string.h>

#include "bindings.h"


struct input {
	/** The device's file descriptor. */
	int fd;
	/** Whether the input is currently being ignored. */
	int ignored;
};


struct context {
	/** Buffer of recently read, and parsed, events. */
	struct event *buf;
	/** Buffer of recently read raw events. */
	struct input_event *events;
	/** List of opened input devices. */
	struct input *inputs;
	/** List of file descriptors passed to poll(). */
	struct pollfd *poll_fds;
	/** Size of buf and events. */
	int events_len;
	/** Number of entries populated in inputs. */
	int inputs_len;
	/** Number of entries allocated in inputs. */
	int inputs_cap;
	/** Number of inputs ready to be read. */
	int ready;
};


void* new_context(int buf_size) {
	struct context *ctx = NULL;
	struct context tmp;

	memset(&tmp, sizeof(struct context), 0);

	tmp.inputs = NULL;
	tmp.poll_fds = NULL;

	tmp.buf = (struct event*)calloc(buf_size, sizeof(struct event));
	tmp.events = (struct input_event*)calloc(buf_size, sizeof(struct input_event));
	if (tmp.buf == NULL || tmp.events == NULL) {
		goto _err;
	}

	tmp.events_len = buf_size;

	ctx = (struct context*)calloc(1, sizeof(struct context));
	if (ctx == NULL) {
		goto _err;
	}
	memcpy(ctx, &tmp, sizeof(struct context));

	tmp.buf = NULL;
	tmp.events = NULL;

_err:
	free(tmp.buf);
	free(tmp.events);

	return ctx;
}


void close_context(void *c) {
	struct context *ctx = (struct context*)c;

	if (ctx == NULL) {
		return;
	}

	for (int i = 0; i < ctx->inputs_len; i++) {
		close(ctx->inputs[i].fd);
	}
	free(ctx->inputs);
	free(ctx->poll_fds);

	free(ctx->buf);
	free(ctx->events);

	free(ctx);
}


/**
 * is_keyboard_like checks whether the file descriptor handles all events required by keyboards,
 * returning 1 on success and 0 on failure.
 */
static int is_keyboard_like(int fd) {
	const unsigned long want_bits = (1 << EV_SYN) | (1 << EV_KEY) | (1 << EV_MSC);
	const unsigned long relevant_bits = want_bits | (1 << EV_REL);
	unsigned long got_bits;
	int rv;

	rv = ioctl(fd, EVIOCGBIT(0, sizeof(unsigned long)), &got_bits);
	if (rv < 0) {
		return 0;
	}

	return ((relevant_bits & got_bits) == want_bits);
}


int open_keyboard_fd(void *c, char *device) {
	struct context *ctx = (struct context*)c;
	int fd = -1;
	int rv = -1;

	fd = open(device, O_RDONLY | O_NONBLOCK);
	if (fd < 0) {
		goto _err;
	}

	if (!is_keyboard_like(fd)) {
		rv = 0;
		goto _err;
	}

	if (ctx->inputs_len == ctx->inputs_cap) {
		if (ctx->inputs_cap == 0) {
			ctx->inputs_cap = 1;
		}
		ctx->inputs_cap *= 2;

		ctx->inputs = (struct input*)reallocarray(ctx->inputs, ctx->inputs_cap, sizeof(struct input));
		ctx->poll_fds = (struct pollfd*)reallocarray(ctx->poll_fds, ctx->inputs_cap, sizeof(struct pollfd));
		if (ctx->inputs == NULL || ctx->poll_fds == NULL) {
			goto _err;
		}
	}

	ctx->inputs[ctx->inputs_len].fd = fd;
	ctx->inputs[ctx->inputs_len].ignored = 0;

	ctx->poll_fds[ctx->inputs_len].fd = fd;
	ctx->poll_fds[ctx->inputs_len].events = POLLIN;

	ctx->inputs_len++;
	fd = -1;

	rv = 0;
_err:
	if (fd != -1) {
		close(fd);
	}

	return rv;
}


int get_num_fd(void *c) {
	struct context *ctx = (struct context*)c;
	return ctx->inputs_len;
}


int wait_events(void *c) {
	struct context *ctx = (struct context*)c;
	int tmp;

	for (int i = 0; i < ctx->inputs_len; i++) {
		ctx->poll_fds[i].revents = 0;
	}

	ctx->ready = 0;
	tmp = poll(ctx->poll_fds, ctx->inputs_len, -1);
	if (tmp < 0) {
		if (errno != EINTR) {
			return -1;
		}
		tmp = 0;
	}

	ctx->ready = tmp;
	return ctx->ready;
}


int get_ready(void *c) {
	struct context *ctx = (struct context*)c;
	return ctx->ready;
}


int read_events(void *c, int idx, struct event **events) {
	struct context *ctx = (struct context*)c;
	size_t num;
	int parsed;

	if (
		ctx->ready == 0 ||
		idx > ctx->inputs_len ||
		ctx->poll_fds[idx].revents == 0
	) {
		return 0;
	}

	ctx->ready--;

	if ((ctx->poll_fds[idx].revents & POLLIN) == 0) {
		return 0;
	}

	num = read(ctx->poll_fds[idx].fd, ctx->events, ctx->events_len * sizeof(struct input_event));
	if (num == -1) {
		if (errno != EINTR) {
			return -1;
		}

		return 0;
	}
	num /= sizeof(struct input_event);

	parsed = 0;
	for (ssize_t i = 0; i < num; i++) {
		switch (ctx->events[i].type) {
		case EV_SYN:
			if (ctx->events[i].code == SYN_REPORT) {
				ctx->inputs[idx].ignored = 0;
			} else if (ctx->events[i].code == SYN_DROPPED) {
				ctx->inputs[idx].ignored = 1;
			}

			break;
		case EV_KEY:
			if (ctx->inputs[idx].ignored) {
				continue;
			}

			if (ctx->events[i].code == KEY_RESERVED || ctx->events[i].value == 2) {
				continue;
			}

			ctx->buf[parsed].key = ctx->events[i].code;
			if (ctx->events[i].value == 1) {
				ctx->buf[parsed].status = ST_PRESSED;
			} else {
				ctx->buf[parsed].status = ST_RELEASED;
			}

			parsed++;

			break;
		}
	}
	*events = ctx->buf;

	return parsed;
}
