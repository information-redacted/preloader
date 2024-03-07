#include "proxy.h"
#include "gen/header.h"

static HMODULE hmod = NULL;
static LPCSTR export_names[] = {
    COMMON_EXPORT_NAMES
#if defined(_X86_) && defined(X86_EXPORT_NAMES)
    , X86_EXPORT_NAMES
#endif
#if defined(_WIN64) && defined(X64_EXPORT_NAMES)
    , X64_EXPORT_NAMES
#endif
};
static FARPROC export_procs[_countof(export_names)] = { 0 };

BOOL real_dll_init(void) {
    if (!hmod) {
        TCHAR path[MAX_PATH];
#ifdef _WIN64
        GetSystemDirectory(path, _countof(path));
#else
        GetSystemWow64Directory(path, _countof(path));
#endif
        _tcscat_s(path, _countof(path), _T(DLL_FNAME));
        hmod = LoadLibrary(path);
    }
    return (hmod != NULL);
}

BOOL real_dll_free(void) {
    if (!hmod)
        return FALSE;

    BOOL result = FreeLibrary(hmod);
    if (result)
        hmod = NULL;

    return result;
}

extern "C" FARPROC resolve_export_proc(size_t index) {
    index = index - 1; // Ordinals are offset by 1.
    if (index < _countof(export_names)) {
        if (hmod && export_procs[index])
            return export_procs[index];

        if (real_dll_init()) {
            export_procs[index] = GetProcAddress(hmod, export_names[index]);
            return export_procs[index];
        }
    }
    return NULL;
}
