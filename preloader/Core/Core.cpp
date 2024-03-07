#include "Core.h"

void Core::Init(HINSTANCE hinstDll) {
#if SHOULD_ALLOC_CONSOLE
    AllocConsole();
    freopen_s(reinterpret_cast<FILE**>(stdin), "CONIN$", "r", stdin);
    freopen_s(reinterpret_cast<FILE**>(stdout), "CONOUT$", "w", stdout);
    freopen_s(reinterpret_cast<FILE**>(stderr), "CONOUT$", "w", stderr);
    HANDLE InputHandle = GetStdHandle(STD_INPUT_HANDLE);
    HANDLE OutputHandle = GetStdHandle(STD_OUTPUT_HANDLE);

    SetStdHandle(STD_OUTPUT_HANDLE, OutputHandle);
    SetStdHandle(STD_ERROR_HANDLE, OutputHandle);
#endif

    std::cout << "Preloader: Core initialized.\r\n";
}

// Copyright 2023 Ethan Cordray <ben@zehow.lt>
// data = start of memory
// max_length = how much memory is there to search for. usually size of the module in bytes
// pattern = you pattern in IDA format (55 ? ? ? ? 8B)
// offset = how many bytes away is the address you're looking for from the start of the pattern
uint8_t* Core::FindPattern(uint8_t* data, size_t max_length, const char* pattern, uintptr_t offset = 0)
{
    std::vector<int> bytes;
    char* start = const_cast<char*>(pattern);
    char* end = start + strlen(pattern);

    for (char* current = start; current < end; ++current)
    {
        if (*current == '?')
        {
            ++current;
            if (*current == '?')
                ++current;
            bytes.push_back(-1);
        }
        else
        {
            bytes.push_back(strtoul(current, &current, 16));
        }
    }

    size_t pattern_size = bytes.size();
    const int* bytes_data = bytes.data();

    for (size_t i = 0u; i < max_length - pattern_size; ++i)
    {
        bool found = true;
        for (size_t j = 0u; j < pattern_size; ++j)
        {
            if (data[i + j] != bytes_data[j] && bytes_data[j] != -1)
            {
                found = false;
                break;
            }
        }

        if (found)
            return (uint8_t*)(data + i + offset);

    }

    return nullptr;
}