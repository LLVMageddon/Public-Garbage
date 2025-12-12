#include "include/log.h"
#include <stdarg.h>
#include <stdio.h>

static FILE *log_file = NULL;
static logging_level current_level = LOG_LEVEL_INFO;

static const char *level_names[] = {
    "DEBUG",
    "INFO",
    "WARN",
    "ERROR",
    "FATAL"
};


int log_init(const char *path, logging_level level){
    if(path){
        log_file = fopen(path, "a");
        if(!log_file){
            return -1;
        }
    }
    else{
        log_file = stderr;
    }

    if(level){
        current_level = level;
    }

    return 0;
}

void log_shutdown(void){

    if(log_file && log_file != stderr){
        fflush(log_file);
        fclose(log_file);
        log_file = NULL;
    }

}

void log_set_level(logging_level level){
    current_level = level;
}

void log_log(logging_level level, const char *fmt, ...){

    if (level < current_level){
        return;
    }

    FILE *file = log_file ? log_file : stderr;
    fprintf(file, "[DATETIME] %-5s: ", level_names[level]);//TODO: find a way to get datetime

    va_list ap;
    va_start(ap, fmt);
    vfprintf(file, fmt, ap);
    va_end(ap);

    fprintf(file,"\n");
    fflush(file);
}
