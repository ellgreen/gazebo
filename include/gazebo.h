#ifndef GAZEBO_H
#define GAZEBO_H

#include <assert.h>
#include <inttypes.h>
#include <stdio.h>
#include <stdlib.h>

#define _G_DEBUG(...)                                                   \
    do {                                                                \
        fprintf(stderr, "%s:%s(%s) :: ", __FILE__, __LINE__, __FUNC__); \
        fprintf(stderr, __VA_ARGS__);                                   \
    } while (0)

#define G_ASSERT(condition) assert(condition)

#endif /* GAZEBO_H */
