package nekoutils

import "runtime"

var Windows_Protect_BindInterfaceIndex func() uint32

func Windows_ShouldApplyWindowsProtect() bool {
	if runtime.GOOS != "windows" {
		return false
	}

	return Windows_Protect_BindInterfaceIndex != nil
}
