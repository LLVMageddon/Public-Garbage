#include <stdio.h>
#include <stdlib.h>
#include "include/log.h"
int main()
{
    printf("DBP Application\n");

    // extern int run_log_unit_tests(void);

    if(log_init(NULL, LOG_LEVEL_DEBUG) != 0){
        return -1;
    }

    LOG_INFO("LOGGING INITIALIZED");
    LOG_DEBUG("log debug");
    LOG_INFO("log info");
    LOG_WARN("log warn");
    LOG_ERROR("log error");
    // LOG_FATAL("log fatal");

    log_shutdown();


    if(log_init("DBP.log", LOG_LEVEL_DEBUG) != 0){
        return -1;
    }

    LOG_INFO("LOGGING INITIALIZED");
    LOG_DEBUG("log debug");
    LOG_INFO("log info");
    LOG_WARN("log warn");
    LOG_ERROR("log error");
    LOG_FATAL("log fatal");

    log_shutdown();

    return 0;
}
