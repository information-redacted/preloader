#pragma once
#include <Windows.h>
#include <stdio.h>
#include <tchar.h>

BOOL real_dll_init(void);
BOOL real_dll_free(void);
extern "C" FARPROC resolve_export_proc(size_t index);