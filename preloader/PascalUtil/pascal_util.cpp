#include "pascal_util.h"

char* PascalUtil::MakeStringA(const std::string& input) {
    size_t totalSize = sizeof(size_t) + input.length() + 1;
    size_t length = input.length();
    char* buffer = (char*)malloc(sizeof(char) * totalSize);
    if (buffer == nullptr) {
        return nullptr;
    }

    std::memcpy(buffer, &length, sizeof(size_t));
    std::memcpy((void*)((char*)buffer + sizeof(size_t)), input.c_str(), input.length());

    return (buffer + sizeof(size_t)); // to free, -sizeof(size_t) it
}

wchar_t* PascalUtil::MakeStringW(const std::wstring& input) {
    size_t totalSize = sizeof(size_t) + (input.size() * sizeof(wchar_t)) + 1;
    size_t length = input.size() * sizeof(wchar_t);
    wchar_t* buffer = (wchar_t*)malloc(totalSize);
    if (buffer == nullptr) {
        return nullptr;
    }
    std::memcpy((void*)((char*)buffer), &length, sizeof(size_t));
    std::memcpy((void*)((char*)buffer + sizeof(size_t)), input.c_str(), input.size() * sizeof(wchar_t));

    return (wchar_t*)((char*)buffer + sizeof(size_t)); // to free, -sizeof(size_t) it
}

void PascalUtil::Free(void* ptr) {
    void* pPtr = NULL;
    if (sizeof(size_t) == 8) {
        pPtr = (void*)((uint64_t)ptr - 8);
    }
    else {
        pPtr = (void*)((uint32_t)ptr - 4);
    }

    if (pPtr != NULL) free(pPtr);
}