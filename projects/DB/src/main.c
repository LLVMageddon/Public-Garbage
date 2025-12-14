#include <stdio.h>
#include <stdlib.h>
#include "include/log.h"
#include "tests/test_log.c"


int tests();
int run_log_tests();

int main(){
    printf("DBP Application\n");
    int test = tests();
    if(test != 0){
        return -1;
    }
    return 0;
}


int tests(){
    return  run_log_tests();
}

int run_log_tests(){

    extern int run_log_unit_tests(void);
    int log_test1 = run_log_unit_tests();

    if(log_test1 != 0){
        printf("LOGGING TEST PASSED\n");
    }
    else{

        printf("LOGGING TEST FAILED\n");
    }
    return log_test1;
}
