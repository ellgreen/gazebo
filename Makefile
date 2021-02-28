CC = gcc
CFLAGS = -I./include -g -Wall -Wextra -Wpedantic

C_SRC = $(wildcard src/*.c)
C_OBJ = $(C_SRC:.c=.o)

PROG = gazebo

$(PROG): $(C_OBJ)
	$(CC) $(CFLAGS) $(C_OBJ) -o $(PROG)

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	rm -f $(PROG) src/*.o
