/* Windows link driver. GitHub Actions cannot exec scripts/cc-static-jpeg.
   Drops -ljpeg so the link uses the libjpeg.a path in CGO_LDFLAGS. */
#include <process.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main(int argc, char **argv) {
	const char *cc = getenv("REAL_CC");
	char **out;
	int i, n;

	if (cc == NULL || cc[0] == '\0') {
		cc = "gcc";
	}
	out = calloc((size_t)argc + 1, sizeof(char *));
	if (out == NULL) {
		return 127;
	}
	n = 0;
	out[n++] = (char *)cc;
	for (i = 1; i < argc; i++) {
		if (strcmp(argv[i], "-ljpeg") == 0) {
			continue;
		}
		out[n++] = argv[i];
	}
	out[n] = NULL;
	_spawnvp(_P_OVERLAY, cc, (const char *const *)out);
	perror(cc);
	return 127;
}
