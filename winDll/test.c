#include <windows.h>
#include <winbase.h>
#include <winuser.h>
#include <stdio.h>
#include <stdint.h>
#include <unistd.h>

/**
 * Resources:
 *
 *  * https://msdn.microsoft.com/en-US/library/windows/desktop/ms644990(v=vs.85).aspx
 *  * https://msdn.microsoft.com/en-US/library/windows/desktop/ms644960(v=vs.85).aspx#installing_releasing
 */

typedef uint16_t __stdcall POP(void);
typedef void __stdcall WAIT(void);

int main(int argc, char *argv[]) {
    HOOKPROC pKeyCallback;
    POP *pPop;
    WAIT *pWait;
    static HINSTANCE hinstDLL;
    static HHOOK keyHook;
    uint32_t time;

    if (argc != 2) {
        printf("Not enough args!\n");
        printf("Usage: %s <path-to-DLL>\n", argv[0]);
        return 1;
    }

    hinstDLL = LoadLibraryA(argv[1]);
    if (hinstDLL == 0) {
        printf("Failed to load the DLL \"%s\"!\n", argv[1]);
        goto __exit;
    }
    pKeyCallback = (HOOKPROC)GetProcAddress(hinstDLL, "keyCallback");
    if (pKeyCallback == 0) {
        printf("Failed to get function's \"keyCallback\" address!\n");
        goto __ret;
    }

    keyHook = SetWindowsHookEx(
        WH_KEYBOARD_LL,
        pKeyCallback,
        hinstDLL,
        0);
    if (keyHook == 0) {
        printf("Failed to setup the hook\n");
        goto __ret;
    }

    pPop = (POP*)GetProcAddress(hinstDLL, "pop");
    if (pPop == 0) {
        printf("Failed to get function's \"pop\" address!\n");
        goto __ret2;
    }

    pWait = (WAIT*)GetProcAddress(hinstDLL, "wait");
    if (pWait == 0) {
        printf("Failed to get function's \"wait\" address!\n");
        goto __ret2;
    }

    printf("scanning...\n");
    (*pWait)();
    printf("printing...\n");
    do {
        uint16_t key;

        key = (*pPop)();
        while (key != 0x12) {
            printf("got key: %04x\n", (int)key);
            key = (*pPop)();
        }
    } while (0);

__ret2:
    UnhookWindowsHookEx(keyHook);
__ret:
    FreeLibrary(hinstDLL);
__exit:

    return 0;
}
