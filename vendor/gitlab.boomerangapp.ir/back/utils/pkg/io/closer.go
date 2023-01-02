package io

import "io"

//SafeClose safe close nil io.Closer
func SafeClose(closer io.Closer) {
	if closer != nil {
		_ = closer.Close()
	}
}
