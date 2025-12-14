#include "include/mem.h"
#include <pthread.h>
#include <stdlib.h>
#include <string.h>
#include "include/log.h"

// I haven't even started and I'm already thinking C is not the right language for this project.
// Maybe I should use go, I'll see how this MEMORY MANAGEMENT goes before I make the change.

static pthread_mutex_t mem_lock = PTHREAD_MUTEX_INITIALIZER;
static size_t alloc_count = 0;
static size_t alloc_bytes = 0;

void mem_init(void){}

void mem_shutdown(void){
    pthread_mutex_lock(&mem_lock);
    if (alloc_count != 0){
        LOG_WARN("Memory leak detected: %zu allocations, %zu bytes still allocated.",
                 alloc_count, alloc_bytes);
    }
    else{
        LOG_INFO("No memory leaks detected.");
    }

    pthread_mutex_unlock(&mem_lock);
}

// Wrappers
void *mem_alloc(size_t size, const char *tag){
    size_t total = sizeof(mem_hdr_t) + size + sizeof(mem_ftr_t); // {[hdr][pay][ftr]}
    mem_hdr_t *hdr = (mem_hdr_t *)malloc(total);
    if(!hdr){
        return NULL;
    }
    hdr->magic = MEM_MAGIC_HEAD;
    hdr->size = size;
    hdr->tag = tag;
    void *payload = (void *)(hdr +1);
    mem_ftr_t *ftr = (mem_ftr_t *)((char *)payload + size);
    ftr->magic = MEM_MAGIC_FOOT;

    pthread_mutex_lock(&mem_lock);
    alloc_count +=1;
    alloc_bytes += size;
    pthread_mutex_unlock(&mem_lock);

    return payload;

}

void *mem_calloc(size_t nmemb, size_t size, const char *tag){
    size_t total = nmemb * size;
    void *p = mem_alloc(total, tag);
    if(p){
        memset(p, 0, total);
    }
    return p;
}

void *mem_realoc(void *ptr, size_t size, const char *tag){
    if(!ptr){
        return mem_alloc(size, tag);
    }    
    mem_hdr_t *old = (mem_hdr_t *)ptr - 1;
    if(old->magic != MEM_MAGIC_HEAD){
        LOG_ERROR("mem_alloc: invalid header magic(possible double free or overflow)");
        return NULL;
    }

    size_t old_size = old->size;
    size_t total = sizeof(mem_hdr_t) + size + sizeof(mem_ftr_t);
    mem_hdr_t *newp = (mem_hdr_t *)realloc(old, total);
    if(!newp){
        return NULL;
    }
    newp->size = size;
    void *payload = (void *)(newp + 1);
    mem_ftr_t *ftr = (mem_ftr_t *)((char *)payload + size);
    ftr->magic = MEM_MAGIC_FOOT;

    pthread_mutex_lock(&mem_lock);
    alloc_bytes += (size - old_size);
    pthread_mutex_unlock(&mem_lock);

    return payload;
}


void *mem_free(void *ptr){
    if(!ptr){
        return NULL;
    }
    mem_hdr_t *hdr = (mem_hdr_t *)ptr - 1;
    mem_ftr_t *ftr = (mem_ftr_t *)((char *)ptr + hdr->size);

    if(hdr->magic != MEM_MAGIC_HEAD){
        LOG_ERROR("");
    }
    else if(ftr->magic != MEM_MAGIC_FOOT){
        LOG_ERROR("");
    }

    pthread_mutex_lock(&mem_lock);
    alloc_count -= 1;
    alloc_bytes -= hdr->size;
    pthread_mutex_unlock(&mem_lock);

#if defined(DEBUG) || defined(_DEBUG)
    memset(hdr + 1, 0xA5, hdr->size);
#endif
    free(hdr);
}

size_t mem_allocated_count(void){
    pthread_mutex_lock(&mem_lock);
    size_t c = alloc_count;
    pthread_mutex_unlock(&mem_lock);
    return c;
}


size_t mem_allocated_bytes(void){

    pthread_mutex_lock(&mem_lock);
    size_t c = alloc_bytes;
    pthread_mutex_unlock(&mem_lock);
    return c;
}

