CFLAGS = -I/usr/local/include/ggml -I/usr/local/include

.PHONY: install
install:
	$(CC) $(CFLAGS) -shared -o libcencodec.so -c cencodec.cpp
	install -v libcencodec.so /usr/local/lib/
