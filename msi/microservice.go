package msi

// #cgo CFLAGS: -I/usr/include/irods
// #cgo CXXFLAGS: -I/usr/include/irods -I/opt/irods-externals/boost1.60.0-0/include -I/opt/irods-externals/clang3.8-0/include/c++/v1 -nostdinc++ -std=c++14
// #cgo LDFLAGS: -lirods_server -lirods_common -lpthread -lc++ -lc++abi
/*
#include <stdlib.h>
#include "call_microservice.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

var rei unsafe.Pointer

func Configure(ruleExecInfo unsafe.Pointer) {
	rei = ruleExecInfo
}

func Call(msiName string, params ...interface{}) error {
	if rei == nil {
		return fmt.Errorf("Unable to call %v, ruleExecInfo is nil, please set using msi.Configure", msiName)
	}

	numParams := C.int(len(params))

	cParams := C.NewParamList(numParams)
	defer C.free(unsafe.Pointer(cParams))

	for inx, param := range params {
		var msParam *Param

		switch p := param.(type) {
		case string:
			msParam = NewParam(STR_MS_T)

			msParam.SetString(p)
		case int:
			msParam = NewParam(INT_MS_T)

			msParam.SetInt(p)
		case int64:
			msParam = NewParam(INT_MS_T)

			msParam.SetInt(int(p))
		case *Param:
			msParam = p
		case nil:
			msParam = new(Param)
		default:
			return fmt.Errorf("Parameter with unknown type passed in params")
		}

		C.SetMsParamListItem(cParams, C.int(inx), msParam.ptr)
	}

	var callInfo C.msiCallInfo_t
	size := unsafe.Sizeof(callInfo)
	C.bzero(unsafe.Pointer(&callInfo), C.size_t(size))

	cMsiName := C.CString(msiName)
	defer C.free(unsafe.Pointer(cMsiName))

	callInfo.microserviceName = cMsiName
	callInfo.params = cParams
	callInfo.paramsLen = numParams
	callInfo.rei = rei

	var errStr *C.char
	if status := C.call_microservice(&callInfo, &errStr); status < 0 {
		return fmt.Errorf("Error in call_microservice: %v", C.GoString(errStr))
	}

	return nil

}