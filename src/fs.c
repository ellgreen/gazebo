#include <gazebo.h>
#include <gazebo/fs.h>

#ifdef G_DEBUG_FS
#define G_DEBUG(...) _G_DEBUG(__VA_ARGS__)
#else
#define G_DEBUG(...)
#endif

int fs_read_file(char* path, g_string* dest)
{
    G_ASSERT(path && dest);

    FILE* fp;
    char  buffer[1024] = {0};

    if (!(fp = fopen(path, "r"))) {
        G_DEBUG("failed to open %s for reading", path);
        return 0;
    }

    while (!feof(fp)) {
        fread(buffer, 1, 1023, fp);
        g_string_append_chars(dest, buffer);
    }

    fclose(fp);

    return 1;
}
