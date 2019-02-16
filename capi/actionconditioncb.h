#ifndef _GREETER_H
#define _GREETER_H

int callPyCb(const char *name);
int registerPyCb (int (*fun_ptrcb)(const char *name));

#endif
