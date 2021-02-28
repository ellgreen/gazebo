CC = gcc
CFLAGS = -I./include -g -Wall -Wextra -Wpedantic
CFLAGS += -DG_DEBUG_MEM
CFLAGS += -DG_DEBUG_STRING
CFLAGS += -DG_DEBUG_FS

C_HED = $(shell find include -type f -name '*.h')
C_SRC = $(wildcard src/*.c)
C_OBJ = $(C_SRC:.c=.o)

PROG = gazebo

$(PROG): $(C_OBJ)
	$(CC) $(CFLAGS) $(C_OBJ) -o $(PROG)

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	rm -f $(PROG) src/*.o

compilecmd:
	make clean
	bear -- make

format:
	clang-format -i $(C_HED) $(C_SRC)
