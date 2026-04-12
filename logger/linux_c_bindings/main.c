#ifndef NO_TEST_MAIN

#include <stdio.h>

#include "bindings.h"


int main(int argc, char *argv[]) {
	void *ctx;

	if (argc < 2) {
		printf("Usage: %s <input-device-1> ...\n", argv[0]);
		return 1;
	}

	ctx = new_context(128);

	for (; argc > 1; argc--) {
		int rv;

		rv = open_keyboard_fd(ctx, argv[argc - 1]);
		if (rv == -1) {
			printf("failed to open input '%s'\n", argv[argc - 1]);
			goto _err;
		}
	}

	printf("got %d inputs\n", get_num_fd(ctx));

	if (get_num_fd(ctx) > 0) {
		while (1) {
			int ready;

			ready = wait_events(ctx);
			if (ready == -1) {
				printf("failed to wait for events\n");
				goto _err;
			}
			printf("ready fd: %d\n", ready);

			for (int i = 0; i < get_num_fd(ctx) && get_ready(ctx) > 0; i++) {
				struct event *events;
				int num;

				num = read_events(ctx, i, &events);
				if (num == -1) {
					printf("failed to get events\n");
					goto _err;
				}
				printf("num events: %d\n", num);

				for (int j = 0; j < num; j++) {
					printf("key: %d state: %d\n", events[j].key, events[j].status);
				}
			}
		}
	}

_err:
	close_context(ctx);
	return 0;
}

#endif /* NO_TEST_MAIN*/
