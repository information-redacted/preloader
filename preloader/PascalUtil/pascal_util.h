#pragma once
#include <Windows.h>
#include <string>

class PascalUtil {
public:
	static char* MakeStringA(const std::string& input);
	static wchar_t* MakeStringW(const std::wstring& input);
	static void     Free(void* ptr);
};