#include "../include/log.h"
#include <stdlib.h>
#include <stdio.h>



static int log_test_to_screen(void){

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

    return 0;
}

static int log_test_to_file(void){

    if(log_init("DBP.log", LOG_LEVEL_DEBUG) != 0){
        return -1;
    }
        LOG_INFO("LOGGING INITIALIZED");
        LOG_DEBUG("log debug");
        LOG_INFO("log info");
        LOG_WARN("log warn");
        LOG_ERROR("log error");
        // LOG_FATAL("log fatal");

        log_shutdown();

        return 0;
    }

    int run_log_unit_tests(void){
        return 
        (log_test_to_screen() + 1 )+
        log_test_to_file();

    }
