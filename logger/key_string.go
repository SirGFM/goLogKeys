// Code generated by "stringer -type Key,KeyState logger/logger.go"; DO NOT EDIT.

package logger

import "strconv"

const _Key_name = "BackspaceTabReturnEscSpacePageUpPageDownEndHomeLeftUpRightDownInsertDeleteKey0Key1Key2Key3Key4Key5Key6Key7Key8Key9ABCDEFGHIJKLMNOPQRSTUVWXYZNumpad0Numpad1Numpad2Numpad3Numpad4Numpad5Numpad6Numpad7Numpad8Numpad9F1F2F3F4F5F6F7F8F9F10F11F12LShiftRShiftLCtrlRCtrlLAltRAltCommaPeriod"

var _Key_index = [...]uint16{0, 9, 12, 18, 21, 26, 32, 40, 43, 47, 51, 53, 58, 62, 68, 74, 78, 82, 86, 90, 94, 98, 102, 106, 110, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 147, 154, 161, 168, 175, 182, 189, 196, 203, 210, 212, 214, 216, 218, 220, 222, 224, 226, 228, 231, 234, 237, 243, 249, 254, 259, 263, 267, 272, 278}

func (i Key) String() string {
	i -= 1
	if i >= Key(len(_Key_index)-1) {
		return "Key(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _Key_name[_Key_index[i]:_Key_index[i+1]]
}

const _KeyState_name = "ReleasedPressed"

var _KeyState_index = [...]uint8{0, 8, 15}

func (i KeyState) String() string {
	if i >= KeyState(len(_KeyState_index)-1) {
		return "KeyState(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _KeyState_name[_KeyState_index[i]:_KeyState_index[i+1]]
}