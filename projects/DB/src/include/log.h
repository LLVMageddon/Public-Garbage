#ifndef LOG_H
#define LOG_H

/* 
TODO:

I'll reinvent the wheel, I always wanted to create a logging library for C. 

REFERENCE: LOG4J

Logging Levels:
    1. FATAL
    2. ERROR
    3. WARN
    4. INFO
    5. DEBUG
    6. TRACE
    7. ALL
    8. OFF

Logging Functions
    1. logger.error(String) // I'll use this type of function call, I'll try to use macros for each logging level. e.g: LOG_ERROR(char[])
    2. logger.error(String, value) // IGNORE THIS
    3. logger.error(String, exception) // IGNORE THIS
    ...

Other Logging Functions
    1. A init function
    2. A shutdown function
    3. A way to set logging level
    3. The actual log function (which would be mapped ot a macro)

*/   
typedef enum{
    LOG_LEVEL_DEBUG = 0,
    LOG_LEVEL_INFO,
    LOG_LEVEL_WARN,
    LOG_LEVEL_ERROR,
    LOG_LEVEL_FATAL

} logging_level;

int log_init(const char *path, logging_level level);
void log_shutdown(void);
void log_set_level(logging_level level);
void log_log(logging_level level, const char *fmt, ...);

#define LOG_DEBUG(...) log_log(LOG_LEVEL_DEBUG, __VA_ARGS__)
#define LOG_INFO(...) log_log(LOG_LEVEL_INFO, __VA_ARGS__)
#define LOG_WARN(...) log_log(LOG_LEVEL_WARN, __VA_ARGS__)
#define LOG_ERROR(...) log_log(LOG_LEVEL_ERROR, __VA_ARGS__)
#define LOG_FATAL(...) do { log_log(LOG_LEVEL_FATAL, __VA_ARGS__); abort();} while(0)

#endif 
