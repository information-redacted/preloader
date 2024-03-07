#pragma once
#include <Windows.h>
#include <iostream>
#include <vector>

#ifndef SHOULD_ALLOC_CONSOLE
#define SHOULD_ALLOC_CONSOLE TRUE
#endif

class Core
{
public:
	static void Init(HINSTANCE hinstDLL);
	static uint8_t* FindPattern(uint8_t* data, size_t max_length, const char* pattern, uintptr_t offset);
private:

};