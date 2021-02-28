#ifndef GAZEBO_H
#define GAZEBO_H

#include <assert.h>
#include <inttypes.h>
#include <stdio.h>
#include <stdlib.h>

#include <gazebo/fs.h>
#include <gazebo/malloc.h>
#include <gazebo/string.h>

#define _G_DEBUG(...)                                                              \
    do {                                                                           \
        fprintf(stderr, "%15.15s:%-4d %10.10s :: ", __FILE__, __LINE__, __func__); \
        fprintf(stderr, __VA_ARGS__);                                              \
        fprintf(stderr, "\n");                                                     \
    } while (0)

#define G_ASSERT(condition) assert(condition)
#define G_UNREACHED() G_ASSERT(0)

#endif /* GAZEBO_H */
