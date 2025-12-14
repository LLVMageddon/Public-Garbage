#define _POSIX_C_SOURCE 200809L

#include "include/log.h"
#include <stdio.h>
#include <pthread.h>
#include <time.h>
#include <string.h>
#include <stdlib.h>

void get_timestamp(char *buffer, size_t n);

static FILE *log_file = NULL;
static logging_level current_level = LOG_LEVEL_INFO;
static pthread_mutex_t log_lock = PTHREAD_MUTEX_INITIALIZER;// FROM log4J Threads are required

static const char *level_names[] = {
    "DEBUG",
    "INFO",
    "WARN",
    "ERROR",
    "FATAL"
};


int log_init(const char *path, logging_level level){
    pthread_mutex_lock(&log_lock);
    if(path){
        log_file = fopen(path, "a");
        if(!log_file){
            pthread_mutex_unlock(&log_lock);
            // perror("fopen");
            return -1;
        }
    }
    else{
        log_file = stderr;
    }

    if(level){
        current_level = level;
    }

    pthread_mutex_unlock(&log_lock);
    return 0;
}

void log_shutdown(void){

    pthread_mutex_lock(&log_lock);
    if(log_file && log_file != stderr){
        fflush(log_file);
        fclose(log_file);
        log_file = NULL;
    }

    pthread_mutex_unlock(&log_lock);

}

void log_set_level(logging_level level){

    pthread_mutex_lock(&log_lock);
    current_level = level;
    pthread_mutex_unlock(&log_lock);
}

void log_log(logging_level level, const char *fmt, ...){

    if (level < current_level){
        return;
    }

    char ts[64]; // Timestamp
    get_timestamp(ts, sizeof(ts));

    pthread_mutex_lock(&log_lock);
    FILE *file = log_file ? log_file : stderr;
    fprintf(file, "[%s] %-5s: ", ts, level_names[level]);//TODO: find a way to get datetime

    va_list ap;
    va_start(ap, fmt);
    vfprintf(file, fmt, ap);
    va_end(ap);

    fprintf(file,"\n");
    fflush(file);
    pthread_mutex_unlock(&log_lock);
}

void get_timestamp(char *buffer, size_t buf_len){
    /*
     This function was inspired by an old stackoverflow post.
    */
    struct timespec ts; // POSIX.1b structure for a time value that has nano seconds.
    clock_gettime(CLOCK_REALTIME, &ts);
    struct tm tm;
    localtime_r(&ts.tv_sec, &tm);
    snprintf(buffer, buf_len, 
             "%02d-%02d-%04d %02d:%02d:%02d.%03ld",// DD-MM-CCYY HH:MM:SS.NS 
             tm.tm_mday,
             tm.tm_mon + 1,
             tm.tm_year + 1900,
             tm.tm_hour,
             tm.tm_min,
             tm.tm_sec,
             ts.tv_nsec / 1000000

             );// This snprintf is similar to fprint, but instead of printing to screen you are able to print formatted output to a character array.
}
