#include "greeter.h"
#include <stdio.h>

//fun_ptr is a pointer to a function that takes a string and returns void
int (*fun_ptr)(const char *name);


int my_fun (const char *name) {
    printf ("Hi there %s\n", name);
    return 0;
}

int callPyCb(const char *name) {
    if (fun_ptr != NULL) {
        fun_ptr (name);
    } else {
        printf ("NULL if fun_ptr\n");
    }
    return 0;
}

int registerPyCb (int (*fun_ptrcb)(const char *name)) {
    fun_ptr = fun_ptrcb;
    printf ("registered a callback funtion....\n");
    return 0;
}



