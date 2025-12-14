#ifndef MEM_H
#define MEM_H

#include <stddef.h>
#include <stdint.h>

// Memory layout = {[HEADER][PAYLOAD][FOOTER]}

#define MEM_MAGIC_HEAD 0xDEADBEEFDEADBEEFULL
#define MEM_MAGIC_FOOT 0xBEEFDEADBEEFDEADULL

typedef struct mem_hdr{
    uint64_t magic;
    size_t size;
    const char *tag;
} mem_hdr_t;

typedef struct mem_ftr{
    uint64_t magic;
} mem_ftr_t;

void mem_init(void);
void mem_shutdown(void);

// Wrappers
void *mem_alloc(size_t size, const char *tag);
void *mem_calloc(size_t nmemb, size_t size, const char *tag);
void *mem_realoc(void *ptr, size_t size, const char *tag);
void *mem_free(void *ptr);

size_t mem_allocated_count(void);
size_t mem_allocated_bytes(void);

#endif
