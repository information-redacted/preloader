#include <Windows.h>
#include "Core/Core.h"
#include "Proxy/proxy.h"

BOOL APIENTRY DllMain(HMODULE hModule, DWORD  ul_reason_for_call, LPVOID lpReserved)
{
    switch (ul_reason_for_call)
    {
    case DLL_PROCESS_ATTACH:
        DisableThreadLibraryCalls(hModule);
        Core::Init(hModule);

        break;
    case DLL_PROCESS_DETACH:
        if (!lpReserved)
            real_dll_free();
        break;
    default:
        break;
    }
    return TRUE;
}