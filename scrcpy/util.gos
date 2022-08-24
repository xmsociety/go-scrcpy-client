package scrcpy

import (
	"reflect"
	"unsafe"

	/*
		#cgo darwin LDFLAGS: -lavformat -lavutil -lavcodec
		#include <libavutil/avutil.h>
		#include <libavformat/avformat.h>
		static void libav_init() {
			av_register_all();
			av_log_set_level(AV_LOG_DEBUG);
		}
	*/
	"C"
)

func init() {
	C.libav_init()
}

func fromCPtr(buf unsafe.Pointer, size int) (ret []uint8) {
	hdr := (*reflect.SliceHeader)((unsafe.Pointer(&ret)))
	hdr.Cap = size
	hdr.Len = size
	hdr.Data = uintptr(buf)
	return
}
