package main

//#include <stdint.h>
//#include <stdlib.h>
//#include <string.h>
//#include "moc.h"
import "C"
import (
	"strings"
	"time"
	"unsafe"

	custom_accounts_902b5fm "github.com/pavelkrolevets/MIR-pro/accounts"
	"github.com/therecipe/qt"
	std_core "github.com/therecipe/qt/core"
)

func cGoFreePacked(ptr unsafe.Pointer) { std_core.NewQByteArrayFromPointer(ptr).DestroyQByteArray() }
func cGoUnpackString(s C.struct_Moc_PackedString) string {
	defer cGoFreePacked(s.ptr)
	if int(s.len) == -1 {
		return C.GoString(s.data)
	}
	return C.GoStringN(s.data, C.int(s.len))
}
func cGoUnpackBytes(s C.struct_Moc_PackedString) []byte {
	defer cGoFreePacked(s.ptr)
	if int(s.len) == -1 {
		gs := C.GoString(s.data)
		return []byte(gs)
	}
	return C.GoBytes(unsafe.Pointer(s.data), C.int(s.len))
}
func unpackStringList(s string) []string {
	if len(s) == 0 {
		return make([]string, 0)
	}
	return strings.Split(s, "¡¦!")
}

type ApproveSignDataCtx_ITF interface {
	std_core.QObject_ITF
	ApproveSignDataCtx_PTR() *ApproveSignDataCtx
}

func (ptr *ApproveSignDataCtx) ApproveSignDataCtx_PTR() *ApproveSignDataCtx {
	return ptr
}

func (ptr *ApproveSignDataCtx) Pointer() unsafe.Pointer {
	if ptr != nil {
		return ptr.QObject_PTR().Pointer()
	}
	return nil
}

func (ptr *ApproveSignDataCtx) SetPointer(p unsafe.Pointer) {
	if ptr != nil {
		ptr.QObject_PTR().SetPointer(p)
	}
}

func PointerFromApproveSignDataCtx(ptr ApproveSignDataCtx_ITF) unsafe.Pointer {
	if ptr != nil {
		return ptr.ApproveSignDataCtx_PTR().Pointer()
	}
	return nil
}

func NewApproveSignDataCtxFromPointer(ptr unsafe.Pointer) (n *ApproveSignDataCtx) {
	if gPtr, ok := qt.Receive(ptr); !ok {
		n = new(ApproveSignDataCtx)
		n.SetPointer(ptr)
	} else {
		switch deduced := gPtr.(type) {
		case *ApproveSignDataCtx:
			n = deduced

		case *std_core.QObject:
			n = &ApproveSignDataCtx{QObject: *deduced}

		default:
			n = new(ApproveSignDataCtx)
			n.SetPointer(ptr)
		}
	}
	return
}

//export callbackApproveSignDataCtx687eda_Constructor
func callbackApproveSignDataCtx687eda_Constructor(ptr unsafe.Pointer) {
	this := NewApproveSignDataCtxFromPointer(ptr)
	qt.Register(ptr, this)
	this.ConnectClicked(this.clicked)
	this.ConnectOnBack(this.onBack)
	this.ConnectOnApprove(this.onApprove)
	this.ConnectOnReject(this.onReject)
	this.ConnectEdited(this.edited)
}

//export callbackApproveSignDataCtx687eda_Clicked
func callbackApproveSignDataCtx687eda_Clicked(ptr unsafe.Pointer, b C.int) {
	if signal := qt.GetSignal(ptr, "clicked"); signal != nil {
		(*(*func(int))(signal))(int(int32(b)))
	}

}

func (ptr *ApproveSignDataCtx) ConnectClicked(f func(b int)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "clicked") {
			C.ApproveSignDataCtx687eda_ConnectClicked(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "clicked")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "clicked"); signal != nil {
			f := func(b int) {
				(*(*func(int))(signal))(b)
				f(b)
			}
			qt.ConnectSignal(ptr.Pointer(), "clicked", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "clicked", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectClicked() {
	if ptr.Pointer() != nil {
		C.ApproveSignDataCtx687eda_DisconnectClicked(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "clicked")
	}
}

func (ptr *ApproveSignDataCtx) Clicked(b int) {
	if ptr.Pointer() != nil {
		C.ApproveSignDataCtx687eda_Clicked(ptr.Pointer(), C.int(int32(b)))
	}
}

//export callbackApproveSignDataCtx687eda_OnBack
func callbackApproveSignDataCtx687eda_OnBack(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "onBack"); signal != nil {
		(*(*func())(signal))()
	}

}

func (ptr *ApproveSignDataCtx) ConnectOnBack(f func()) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "onBack") {
			C.ApproveSignDataCtx687eda_ConnectOnBack(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "onBack")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "onBack"); signal != nil {
			f := func() {
				(*(*func())(signal))()
				f()
			}
			qt.ConnectSignal(ptr.Pointer(), "onBack", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "onBack", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectOnBack() {
	if ptr.Pointer() != nil {
		C.ApproveSignDataCtx687eda_DisconnectOnBack(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "onBack")
	}
}

func (ptr *ApproveSignDataCtx) OnBack() {
	if ptr.Pointer() != nil {
		C.ApproveSignDataCtx687eda_OnBack(ptr.Pointer())
	}
}

//export callbackApproveSignDataCtx687eda_OnApprove
func callbackApproveSignDataCtx687eda_OnApprove(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "onApprove"); signal != nil {
		(*(*func())(signal))()
	}

}

func (ptr *ApproveSignDataCtx) ConnectOnApprove(f func()) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "onApprove") {
			C.ApproveSignDataCtx687eda_ConnectOnApprove(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "onApprove")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "onApprove"); signal != nil {
			f := func() {
				(*(*func())(signal))()
				f()
			}
			qt.ConnectSignal(ptr.Pointer(), "onApprove", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "onApprove", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectOnApprove() {
	if ptr.Pointer() != nil {
		C.ApproveSignDataCtx687eda_DisconnectOnApprove(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "onApprove")
	}
}

func (ptr *ApproveSignDataCtx) OnApprove() {
	if ptr.Pointer() != nil {
		C.ApproveSignDataCtx687eda_OnApprove(ptr.Pointer())
	}
}

//export callbackApproveSignDataCtx687eda_OnReject
func callbackApproveSignDataCtx687eda_OnReject(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "onReject"); signal != nil {
		(*(*func())(signal))()
	}

}

func (ptr *ApproveSignDataCtx) ConnectOnReject(f func()) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "onReject") {
			C.ApproveSignDataCtx687eda_ConnectOnReject(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "onReject")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "onReject"); signal != nil {
			f := func() {
				(*(*func())(signal))()
				f()
			}
			qt.ConnectSignal(ptr.Pointer(), "onReject", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "onReject", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectOnReject() {
	if ptr.Pointer() != nil {
		C.ApproveSignDataCtx687eda_DisconnectOnReject(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "onReject")
	}
}

func (ptr *ApproveSignDataCtx) OnReject() {
	if ptr.Pointer() != nil {
		C.ApproveSignDataCtx687eda_OnReject(ptr.Pointer())
	}
}

//export callbackApproveSignDataCtx687eda_Edited
func callbackApproveSignDataCtx687eda_Edited(ptr unsafe.Pointer, b C.struct_Moc_PackedString, value C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "edited"); signal != nil {
		(*(*func(string, string))(signal))(cGoUnpackString(b), cGoUnpackString(value))
	}

}

func (ptr *ApproveSignDataCtx) ConnectEdited(f func(b string, value string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "edited") {
			C.ApproveSignDataCtx687eda_ConnectEdited(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "edited")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "edited"); signal != nil {
			f := func(b string, value string) {
				(*(*func(string, string))(signal))(b, value)
				f(b, value)
			}
			qt.ConnectSignal(ptr.Pointer(), "edited", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "edited", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectEdited() {
	if ptr.Pointer() != nil {
		C.ApproveSignDataCtx687eda_DisconnectEdited(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "edited")
	}
}

func (ptr *ApproveSignDataCtx) Edited(b string, value string) {
	if ptr.Pointer() != nil {
		var bC *C.char
		if b != "" {
			bC = C.CString(b)
			defer C.free(unsafe.Pointer(bC))
		}
		var valueC *C.char
		if value != "" {
			valueC = C.CString(value)
			defer C.free(unsafe.Pointer(valueC))
		}
		C.ApproveSignDataCtx687eda_Edited(ptr.Pointer(), C.struct_Moc_PackedString{data: bC, len: C.longlong(len(b))}, C.struct_Moc_PackedString{data: valueC, len: C.longlong(len(value))})
	}
}

//export callbackApproveSignDataCtx687eda_Remote
func callbackApproveSignDataCtx687eda_Remote(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "remote"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveSignDataCtxFromPointer(ptr).RemoteDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveSignDataCtx) ConnectRemote(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "remote"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "remote", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "remote", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectRemote() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "remote")
	}
}

func (ptr *ApproveSignDataCtx) Remote() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveSignDataCtx687eda_Remote(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveSignDataCtx) RemoteDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveSignDataCtx687eda_RemoteDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveSignDataCtx687eda_SetRemote
func callbackApproveSignDataCtx687eda_SetRemote(ptr unsafe.Pointer, remote C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setRemote"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(remote))
	} else {
		NewApproveSignDataCtxFromPointer(ptr).SetRemoteDefault(cGoUnpackString(remote))
	}
}

func (ptr *ApproveSignDataCtx) ConnectSetRemote(f func(remote string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setRemote"); signal != nil {
			f := func(remote string) {
				(*(*func(string))(signal))(remote)
				f(remote)
			}
			qt.ConnectSignal(ptr.Pointer(), "setRemote", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setRemote", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectSetRemote() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setRemote")
	}
}

func (ptr *ApproveSignDataCtx) SetRemote(remote string) {
	if ptr.Pointer() != nil {
		var remoteC *C.char
		if remote != "" {
			remoteC = C.CString(remote)
			defer C.free(unsafe.Pointer(remoteC))
		}
		C.ApproveSignDataCtx687eda_SetRemote(ptr.Pointer(), C.struct_Moc_PackedString{data: remoteC, len: C.longlong(len(remote))})
	}
}

func (ptr *ApproveSignDataCtx) SetRemoteDefault(remote string) {
	if ptr.Pointer() != nil {
		var remoteC *C.char
		if remote != "" {
			remoteC = C.CString(remote)
			defer C.free(unsafe.Pointer(remoteC))
		}
		C.ApproveSignDataCtx687eda_SetRemoteDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: remoteC, len: C.longlong(len(remote))})
	}
}

//export callbackApproveSignDataCtx687eda_RemoteChanged
func callbackApproveSignDataCtx687eda_RemoteChanged(ptr unsafe.Pointer, remote C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "remoteChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(remote))
	}

}

func (ptr *ApproveSignDataCtx) ConnectRemoteChanged(f func(remote string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "remoteChanged") {
			C.ApproveSignDataCtx687eda_ConnectRemoteChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "remoteChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "remoteChanged"); signal != nil {
			f := func(remote string) {
				(*(*func(string))(signal))(remote)
				f(remote)
			}
			qt.ConnectSignal(ptr.Pointer(), "remoteChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "remoteChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectRemoteChanged() {
	if ptr.Pointer() != nil {
		C.ApproveSignDataCtx687eda_DisconnectRemoteChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "remoteChanged")
	}
}

func (ptr *ApproveSignDataCtx) RemoteChanged(remote string) {
	if ptr.Pointer() != nil {
		var remoteC *C.char
		if remote != "" {
			remoteC = C.CString(remote)
			defer C.free(unsafe.Pointer(remoteC))
		}
		C.ApproveSignDataCtx687eda_RemoteChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: remoteC, len: C.longlong(len(remote))})
	}
}

//export callbackApproveSignDataCtx687eda_Transport
func callbackApproveSignDataCtx687eda_Transport(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "transport"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveSignDataCtxFromPointer(ptr).TransportDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveSignDataCtx) ConnectTransport(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "transport"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "transport", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "transport", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectTransport() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "transport")
	}
}

func (ptr *ApproveSignDataCtx) Transport() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveSignDataCtx687eda_Transport(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveSignDataCtx) TransportDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveSignDataCtx687eda_TransportDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveSignDataCtx687eda_SetTransport
func callbackApproveSignDataCtx687eda_SetTransport(ptr unsafe.Pointer, transport C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setTransport"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(transport))
	} else {
		NewApproveSignDataCtxFromPointer(ptr).SetTransportDefault(cGoUnpackString(transport))
	}
}

func (ptr *ApproveSignDataCtx) ConnectSetTransport(f func(transport string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setTransport"); signal != nil {
			f := func(transport string) {
				(*(*func(string))(signal))(transport)
				f(transport)
			}
			qt.ConnectSignal(ptr.Pointer(), "setTransport", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setTransport", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectSetTransport() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setTransport")
	}
}

func (ptr *ApproveSignDataCtx) SetTransport(transport string) {
	if ptr.Pointer() != nil {
		var transportC *C.char
		if transport != "" {
			transportC = C.CString(transport)
			defer C.free(unsafe.Pointer(transportC))
		}
		C.ApproveSignDataCtx687eda_SetTransport(ptr.Pointer(), C.struct_Moc_PackedString{data: transportC, len: C.longlong(len(transport))})
	}
}

func (ptr *ApproveSignDataCtx) SetTransportDefault(transport string) {
	if ptr.Pointer() != nil {
		var transportC *C.char
		if transport != "" {
			transportC = C.CString(transport)
			defer C.free(unsafe.Pointer(transportC))
		}
		C.ApproveSignDataCtx687eda_SetTransportDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: transportC, len: C.longlong(len(transport))})
	}
}

//export callbackApproveSignDataCtx687eda_TransportChanged
func callbackApproveSignDataCtx687eda_TransportChanged(ptr unsafe.Pointer, transport C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "transportChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(transport))
	}

}

func (ptr *ApproveSignDataCtx) ConnectTransportChanged(f func(transport string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "transportChanged") {
			C.ApproveSignDataCtx687eda_ConnectTransportChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "transportChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "transportChanged"); signal != nil {
			f := func(transport string) {
				(*(*func(string))(signal))(transport)
				f(transport)
			}
			qt.ConnectSignal(ptr.Pointer(), "transportChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "transportChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectTransportChanged() {
	if ptr.Pointer() != nil {
		C.ApproveSignDataCtx687eda_DisconnectTransportChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "transportChanged")
	}
}

func (ptr *ApproveSignDataCtx) TransportChanged(transport string) {
	if ptr.Pointer() != nil {
		var transportC *C.char
		if transport != "" {
			transportC = C.CString(transport)
			defer C.free(unsafe.Pointer(transportC))
		}
		C.ApproveSignDataCtx687eda_TransportChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: transportC, len: C.longlong(len(transport))})
	}
}

//export callbackApproveSignDataCtx687eda_Endpoint
func callbackApproveSignDataCtx687eda_Endpoint(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "endpoint"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveSignDataCtxFromPointer(ptr).EndpointDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveSignDataCtx) ConnectEndpoint(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "endpoint"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "endpoint", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "endpoint", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectEndpoint() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "endpoint")
	}
}

func (ptr *ApproveSignDataCtx) Endpoint() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveSignDataCtx687eda_Endpoint(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveSignDataCtx) EndpointDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveSignDataCtx687eda_EndpointDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveSignDataCtx687eda_SetEndpoint
func callbackApproveSignDataCtx687eda_SetEndpoint(ptr unsafe.Pointer, endpoint C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setEndpoint"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(endpoint))
	} else {
		NewApproveSignDataCtxFromPointer(ptr).SetEndpointDefault(cGoUnpackString(endpoint))
	}
}

func (ptr *ApproveSignDataCtx) ConnectSetEndpoint(f func(endpoint string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setEndpoint"); signal != nil {
			f := func(endpoint string) {
				(*(*func(string))(signal))(endpoint)
				f(endpoint)
			}
			qt.ConnectSignal(ptr.Pointer(), "setEndpoint", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setEndpoint", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectSetEndpoint() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setEndpoint")
	}
}

func (ptr *ApproveSignDataCtx) SetEndpoint(endpoint string) {
	if ptr.Pointer() != nil {
		var endpointC *C.char
		if endpoint != "" {
			endpointC = C.CString(endpoint)
			defer C.free(unsafe.Pointer(endpointC))
		}
		C.ApproveSignDataCtx687eda_SetEndpoint(ptr.Pointer(), C.struct_Moc_PackedString{data: endpointC, len: C.longlong(len(endpoint))})
	}
}

func (ptr *ApproveSignDataCtx) SetEndpointDefault(endpoint string) {
	if ptr.Pointer() != nil {
		var endpointC *C.char
		if endpoint != "" {
			endpointC = C.CString(endpoint)
			defer C.free(unsafe.Pointer(endpointC))
		}
		C.ApproveSignDataCtx687eda_SetEndpointDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: endpointC, len: C.longlong(len(endpoint))})
	}
}

//export callbackApproveSignDataCtx687eda_EndpointChanged
func callbackApproveSignDataCtx687eda_EndpointChanged(ptr unsafe.Pointer, endpoint C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "endpointChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(endpoint))
	}

}

func (ptr *ApproveSignDataCtx) ConnectEndpointChanged(f func(endpoint string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "endpointChanged") {
			C.ApproveSignDataCtx687eda_ConnectEndpointChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "endpointChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "endpointChanged"); signal != nil {
			f := func(endpoint string) {
				(*(*func(string))(signal))(endpoint)
				f(endpoint)
			}
			qt.ConnectSignal(ptr.Pointer(), "endpointChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "endpointChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectEndpointChanged() {
	if ptr.Pointer() != nil {
		C.ApproveSignDataCtx687eda_DisconnectEndpointChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "endpointChanged")
	}
}

func (ptr *ApproveSignDataCtx) EndpointChanged(endpoint string) {
	if ptr.Pointer() != nil {
		var endpointC *C.char
		if endpoint != "" {
			endpointC = C.CString(endpoint)
			defer C.free(unsafe.Pointer(endpointC))
		}
		C.ApproveSignDataCtx687eda_EndpointChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: endpointC, len: C.longlong(len(endpoint))})
	}
}

//export callbackApproveSignDataCtx687eda_From
func callbackApproveSignDataCtx687eda_From(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "from"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveSignDataCtxFromPointer(ptr).FromDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveSignDataCtx) ConnectFrom(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "from"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "from", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "from", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectFrom() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "from")
	}
}

func (ptr *ApproveSignDataCtx) From() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveSignDataCtx687eda_From(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveSignDataCtx) FromDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveSignDataCtx687eda_FromDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveSignDataCtx687eda_SetFrom
func callbackApproveSignDataCtx687eda_SetFrom(ptr unsafe.Pointer, from C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setFrom"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(from))
	} else {
		NewApproveSignDataCtxFromPointer(ptr).SetFromDefault(cGoUnpackString(from))
	}
}

func (ptr *ApproveSignDataCtx) ConnectSetFrom(f func(from string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setFrom"); signal != nil {
			f := func(from string) {
				(*(*func(string))(signal))(from)
				f(from)
			}
			qt.ConnectSignal(ptr.Pointer(), "setFrom", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setFrom", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectSetFrom() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setFrom")
	}
}

func (ptr *ApproveSignDataCtx) SetFrom(from string) {
	if ptr.Pointer() != nil {
		var fromC *C.char
		if from != "" {
			fromC = C.CString(from)
			defer C.free(unsafe.Pointer(fromC))
		}
		C.ApproveSignDataCtx687eda_SetFrom(ptr.Pointer(), C.struct_Moc_PackedString{data: fromC, len: C.longlong(len(from))})
	}
}

func (ptr *ApproveSignDataCtx) SetFromDefault(from string) {
	if ptr.Pointer() != nil {
		var fromC *C.char
		if from != "" {
			fromC = C.CString(from)
			defer C.free(unsafe.Pointer(fromC))
		}
		C.ApproveSignDataCtx687eda_SetFromDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: fromC, len: C.longlong(len(from))})
	}
}

//export callbackApproveSignDataCtx687eda_FromChanged
func callbackApproveSignDataCtx687eda_FromChanged(ptr unsafe.Pointer, from C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "fromChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(from))
	}

}

func (ptr *ApproveSignDataCtx) ConnectFromChanged(f func(from string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "fromChanged") {
			C.ApproveSignDataCtx687eda_ConnectFromChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "fromChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "fromChanged"); signal != nil {
			f := func(from string) {
				(*(*func(string))(signal))(from)
				f(from)
			}
			qt.ConnectSignal(ptr.Pointer(), "fromChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "fromChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectFromChanged() {
	if ptr.Pointer() != nil {
		C.ApproveSignDataCtx687eda_DisconnectFromChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "fromChanged")
	}
}

func (ptr *ApproveSignDataCtx) FromChanged(from string) {
	if ptr.Pointer() != nil {
		var fromC *C.char
		if from != "" {
			fromC = C.CString(from)
			defer C.free(unsafe.Pointer(fromC))
		}
		C.ApproveSignDataCtx687eda_FromChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: fromC, len: C.longlong(len(from))})
	}
}

//export callbackApproveSignDataCtx687eda_Message
func callbackApproveSignDataCtx687eda_Message(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "message"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveSignDataCtxFromPointer(ptr).MessageDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveSignDataCtx) ConnectMessage(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "message"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "message", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "message", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectMessage() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "message")
	}
}

func (ptr *ApproveSignDataCtx) Message() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveSignDataCtx687eda_Message(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveSignDataCtx) MessageDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveSignDataCtx687eda_MessageDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveSignDataCtx687eda_SetMessage
func callbackApproveSignDataCtx687eda_SetMessage(ptr unsafe.Pointer, message C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setMessage"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(message))
	} else {
		NewApproveSignDataCtxFromPointer(ptr).SetMessageDefault(cGoUnpackString(message))
	}
}

func (ptr *ApproveSignDataCtx) ConnectSetMessage(f func(message string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setMessage"); signal != nil {
			f := func(message string) {
				(*(*func(string))(signal))(message)
				f(message)
			}
			qt.ConnectSignal(ptr.Pointer(), "setMessage", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setMessage", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectSetMessage() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setMessage")
	}
}

func (ptr *ApproveSignDataCtx) SetMessage(message string) {
	if ptr.Pointer() != nil {
		var messageC *C.char
		if message != "" {
			messageC = C.CString(message)
			defer C.free(unsafe.Pointer(messageC))
		}
		C.ApproveSignDataCtx687eda_SetMessage(ptr.Pointer(), C.struct_Moc_PackedString{data: messageC, len: C.longlong(len(message))})
	}
}

func (ptr *ApproveSignDataCtx) SetMessageDefault(message string) {
	if ptr.Pointer() != nil {
		var messageC *C.char
		if message != "" {
			messageC = C.CString(message)
			defer C.free(unsafe.Pointer(messageC))
		}
		C.ApproveSignDataCtx687eda_SetMessageDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: messageC, len: C.longlong(len(message))})
	}
}

//export callbackApproveSignDataCtx687eda_MessageChanged
func callbackApproveSignDataCtx687eda_MessageChanged(ptr unsafe.Pointer, message C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "messageChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(message))
	}

}

func (ptr *ApproveSignDataCtx) ConnectMessageChanged(f func(message string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "messageChanged") {
			C.ApproveSignDataCtx687eda_ConnectMessageChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "messageChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "messageChanged"); signal != nil {
			f := func(message string) {
				(*(*func(string))(signal))(message)
				f(message)
			}
			qt.ConnectSignal(ptr.Pointer(), "messageChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "messageChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectMessageChanged() {
	if ptr.Pointer() != nil {
		C.ApproveSignDataCtx687eda_DisconnectMessageChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "messageChanged")
	}
}

func (ptr *ApproveSignDataCtx) MessageChanged(message string) {
	if ptr.Pointer() != nil {
		var messageC *C.char
		if message != "" {
			messageC = C.CString(message)
			defer C.free(unsafe.Pointer(messageC))
		}
		C.ApproveSignDataCtx687eda_MessageChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: messageC, len: C.longlong(len(message))})
	}
}

//export callbackApproveSignDataCtx687eda_RawData
func callbackApproveSignDataCtx687eda_RawData(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "rawData"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveSignDataCtxFromPointer(ptr).RawDataDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveSignDataCtx) ConnectRawData(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "rawData"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "rawData", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "rawData", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectRawData() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "rawData")
	}
}

func (ptr *ApproveSignDataCtx) RawData() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveSignDataCtx687eda_RawData(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveSignDataCtx) RawDataDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveSignDataCtx687eda_RawDataDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveSignDataCtx687eda_SetRawData
func callbackApproveSignDataCtx687eda_SetRawData(ptr unsafe.Pointer, rawData C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setRawData"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(rawData))
	} else {
		NewApproveSignDataCtxFromPointer(ptr).SetRawDataDefault(cGoUnpackString(rawData))
	}
}

func (ptr *ApproveSignDataCtx) ConnectSetRawData(f func(rawData string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setRawData"); signal != nil {
			f := func(rawData string) {
				(*(*func(string))(signal))(rawData)
				f(rawData)
			}
			qt.ConnectSignal(ptr.Pointer(), "setRawData", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setRawData", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectSetRawData() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setRawData")
	}
}

func (ptr *ApproveSignDataCtx) SetRawData(rawData string) {
	if ptr.Pointer() != nil {
		var rawDataC *C.char
		if rawData != "" {
			rawDataC = C.CString(rawData)
			defer C.free(unsafe.Pointer(rawDataC))
		}
		C.ApproveSignDataCtx687eda_SetRawData(ptr.Pointer(), C.struct_Moc_PackedString{data: rawDataC, len: C.longlong(len(rawData))})
	}
}

func (ptr *ApproveSignDataCtx) SetRawDataDefault(rawData string) {
	if ptr.Pointer() != nil {
		var rawDataC *C.char
		if rawData != "" {
			rawDataC = C.CString(rawData)
			defer C.free(unsafe.Pointer(rawDataC))
		}
		C.ApproveSignDataCtx687eda_SetRawDataDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: rawDataC, len: C.longlong(len(rawData))})
	}
}

//export callbackApproveSignDataCtx687eda_RawDataChanged
func callbackApproveSignDataCtx687eda_RawDataChanged(ptr unsafe.Pointer, rawData C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "rawDataChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(rawData))
	}

}

func (ptr *ApproveSignDataCtx) ConnectRawDataChanged(f func(rawData string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "rawDataChanged") {
			C.ApproveSignDataCtx687eda_ConnectRawDataChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "rawDataChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "rawDataChanged"); signal != nil {
			f := func(rawData string) {
				(*(*func(string))(signal))(rawData)
				f(rawData)
			}
			qt.ConnectSignal(ptr.Pointer(), "rawDataChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "rawDataChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectRawDataChanged() {
	if ptr.Pointer() != nil {
		C.ApproveSignDataCtx687eda_DisconnectRawDataChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "rawDataChanged")
	}
}

func (ptr *ApproveSignDataCtx) RawDataChanged(rawData string) {
	if ptr.Pointer() != nil {
		var rawDataC *C.char
		if rawData != "" {
			rawDataC = C.CString(rawData)
			defer C.free(unsafe.Pointer(rawDataC))
		}
		C.ApproveSignDataCtx687eda_RawDataChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: rawDataC, len: C.longlong(len(rawData))})
	}
}

//export callbackApproveSignDataCtx687eda_Hash
func callbackApproveSignDataCtx687eda_Hash(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "hash"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveSignDataCtxFromPointer(ptr).HashDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveSignDataCtx) ConnectHash(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "hash"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "hash", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "hash", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectHash() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "hash")
	}
}

func (ptr *ApproveSignDataCtx) Hash() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveSignDataCtx687eda_Hash(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveSignDataCtx) HashDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveSignDataCtx687eda_HashDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveSignDataCtx687eda_SetHash
func callbackApproveSignDataCtx687eda_SetHash(ptr unsafe.Pointer, hash C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setHash"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(hash))
	} else {
		NewApproveSignDataCtxFromPointer(ptr).SetHashDefault(cGoUnpackString(hash))
	}
}

func (ptr *ApproveSignDataCtx) ConnectSetHash(f func(hash string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setHash"); signal != nil {
			f := func(hash string) {
				(*(*func(string))(signal))(hash)
				f(hash)
			}
			qt.ConnectSignal(ptr.Pointer(), "setHash", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setHash", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectSetHash() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setHash")
	}
}

func (ptr *ApproveSignDataCtx) SetHash(hash string) {
	if ptr.Pointer() != nil {
		var hashC *C.char
		if hash != "" {
			hashC = C.CString(hash)
			defer C.free(unsafe.Pointer(hashC))
		}
		C.ApproveSignDataCtx687eda_SetHash(ptr.Pointer(), C.struct_Moc_PackedString{data: hashC, len: C.longlong(len(hash))})
	}
}

func (ptr *ApproveSignDataCtx) SetHashDefault(hash string) {
	if ptr.Pointer() != nil {
		var hashC *C.char
		if hash != "" {
			hashC = C.CString(hash)
			defer C.free(unsafe.Pointer(hashC))
		}
		C.ApproveSignDataCtx687eda_SetHashDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: hashC, len: C.longlong(len(hash))})
	}
}

//export callbackApproveSignDataCtx687eda_HashChanged
func callbackApproveSignDataCtx687eda_HashChanged(ptr unsafe.Pointer, hash C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "hashChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(hash))
	}

}

func (ptr *ApproveSignDataCtx) ConnectHashChanged(f func(hash string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "hashChanged") {
			C.ApproveSignDataCtx687eda_ConnectHashChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "hashChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "hashChanged"); signal != nil {
			f := func(hash string) {
				(*(*func(string))(signal))(hash)
				f(hash)
			}
			qt.ConnectSignal(ptr.Pointer(), "hashChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "hashChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectHashChanged() {
	if ptr.Pointer() != nil {
		C.ApproveSignDataCtx687eda_DisconnectHashChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "hashChanged")
	}
}

func (ptr *ApproveSignDataCtx) HashChanged(hash string) {
	if ptr.Pointer() != nil {
		var hashC *C.char
		if hash != "" {
			hashC = C.CString(hash)
			defer C.free(unsafe.Pointer(hashC))
		}
		C.ApproveSignDataCtx687eda_HashChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: hashC, len: C.longlong(len(hash))})
	}
}

//export callbackApproveSignDataCtx687eda_Password
func callbackApproveSignDataCtx687eda_Password(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "password"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveSignDataCtxFromPointer(ptr).PasswordDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveSignDataCtx) ConnectPassword(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "password"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "password", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "password", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectPassword() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "password")
	}
}

func (ptr *ApproveSignDataCtx) Password() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveSignDataCtx687eda_Password(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveSignDataCtx) PasswordDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveSignDataCtx687eda_PasswordDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveSignDataCtx687eda_SetPassword
func callbackApproveSignDataCtx687eda_SetPassword(ptr unsafe.Pointer, password C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setPassword"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(password))
	} else {
		NewApproveSignDataCtxFromPointer(ptr).SetPasswordDefault(cGoUnpackString(password))
	}
}

func (ptr *ApproveSignDataCtx) ConnectSetPassword(f func(password string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setPassword"); signal != nil {
			f := func(password string) {
				(*(*func(string))(signal))(password)
				f(password)
			}
			qt.ConnectSignal(ptr.Pointer(), "setPassword", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setPassword", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectSetPassword() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setPassword")
	}
}

func (ptr *ApproveSignDataCtx) SetPassword(password string) {
	if ptr.Pointer() != nil {
		var passwordC *C.char
		if password != "" {
			passwordC = C.CString(password)
			defer C.free(unsafe.Pointer(passwordC))
		}
		C.ApproveSignDataCtx687eda_SetPassword(ptr.Pointer(), C.struct_Moc_PackedString{data: passwordC, len: C.longlong(len(password))})
	}
}

func (ptr *ApproveSignDataCtx) SetPasswordDefault(password string) {
	if ptr.Pointer() != nil {
		var passwordC *C.char
		if password != "" {
			passwordC = C.CString(password)
			defer C.free(unsafe.Pointer(passwordC))
		}
		C.ApproveSignDataCtx687eda_SetPasswordDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: passwordC, len: C.longlong(len(password))})
	}
}

//export callbackApproveSignDataCtx687eda_PasswordChanged
func callbackApproveSignDataCtx687eda_PasswordChanged(ptr unsafe.Pointer, password C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "passwordChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(password))
	}

}

func (ptr *ApproveSignDataCtx) ConnectPasswordChanged(f func(password string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "passwordChanged") {
			C.ApproveSignDataCtx687eda_ConnectPasswordChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "passwordChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "passwordChanged"); signal != nil {
			f := func(password string) {
				(*(*func(string))(signal))(password)
				f(password)
			}
			qt.ConnectSignal(ptr.Pointer(), "passwordChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "passwordChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectPasswordChanged() {
	if ptr.Pointer() != nil {
		C.ApproveSignDataCtx687eda_DisconnectPasswordChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "passwordChanged")
	}
}

func (ptr *ApproveSignDataCtx) PasswordChanged(password string) {
	if ptr.Pointer() != nil {
		var passwordC *C.char
		if password != "" {
			passwordC = C.CString(password)
			defer C.free(unsafe.Pointer(passwordC))
		}
		C.ApproveSignDataCtx687eda_PasswordChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: passwordC, len: C.longlong(len(password))})
	}
}

//export callbackApproveSignDataCtx687eda_FromSrc
func callbackApproveSignDataCtx687eda_FromSrc(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "fromSrc"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveSignDataCtxFromPointer(ptr).FromSrcDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveSignDataCtx) ConnectFromSrc(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "fromSrc"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "fromSrc", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "fromSrc", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectFromSrc() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "fromSrc")
	}
}

func (ptr *ApproveSignDataCtx) FromSrc() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveSignDataCtx687eda_FromSrc(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveSignDataCtx) FromSrcDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveSignDataCtx687eda_FromSrcDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveSignDataCtx687eda_SetFromSrc
func callbackApproveSignDataCtx687eda_SetFromSrc(ptr unsafe.Pointer, fromSrc C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setFromSrc"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(fromSrc))
	} else {
		NewApproveSignDataCtxFromPointer(ptr).SetFromSrcDefault(cGoUnpackString(fromSrc))
	}
}

func (ptr *ApproveSignDataCtx) ConnectSetFromSrc(f func(fromSrc string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setFromSrc"); signal != nil {
			f := func(fromSrc string) {
				(*(*func(string))(signal))(fromSrc)
				f(fromSrc)
			}
			qt.ConnectSignal(ptr.Pointer(), "setFromSrc", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setFromSrc", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectSetFromSrc() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setFromSrc")
	}
}

func (ptr *ApproveSignDataCtx) SetFromSrc(fromSrc string) {
	if ptr.Pointer() != nil {
		var fromSrcC *C.char
		if fromSrc != "" {
			fromSrcC = C.CString(fromSrc)
			defer C.free(unsafe.Pointer(fromSrcC))
		}
		C.ApproveSignDataCtx687eda_SetFromSrc(ptr.Pointer(), C.struct_Moc_PackedString{data: fromSrcC, len: C.longlong(len(fromSrc))})
	}
}

func (ptr *ApproveSignDataCtx) SetFromSrcDefault(fromSrc string) {
	if ptr.Pointer() != nil {
		var fromSrcC *C.char
		if fromSrc != "" {
			fromSrcC = C.CString(fromSrc)
			defer C.free(unsafe.Pointer(fromSrcC))
		}
		C.ApproveSignDataCtx687eda_SetFromSrcDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: fromSrcC, len: C.longlong(len(fromSrc))})
	}
}

//export callbackApproveSignDataCtx687eda_FromSrcChanged
func callbackApproveSignDataCtx687eda_FromSrcChanged(ptr unsafe.Pointer, fromSrc C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "fromSrcChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(fromSrc))
	}

}

func (ptr *ApproveSignDataCtx) ConnectFromSrcChanged(f func(fromSrc string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "fromSrcChanged") {
			C.ApproveSignDataCtx687eda_ConnectFromSrcChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "fromSrcChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "fromSrcChanged"); signal != nil {
			f := func(fromSrc string) {
				(*(*func(string))(signal))(fromSrc)
				f(fromSrc)
			}
			qt.ConnectSignal(ptr.Pointer(), "fromSrcChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "fromSrcChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectFromSrcChanged() {
	if ptr.Pointer() != nil {
		C.ApproveSignDataCtx687eda_DisconnectFromSrcChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "fromSrcChanged")
	}
}

func (ptr *ApproveSignDataCtx) FromSrcChanged(fromSrc string) {
	if ptr.Pointer() != nil {
		var fromSrcC *C.char
		if fromSrc != "" {
			fromSrcC = C.CString(fromSrc)
			defer C.free(unsafe.Pointer(fromSrcC))
		}
		C.ApproveSignDataCtx687eda_FromSrcChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: fromSrcC, len: C.longlong(len(fromSrc))})
	}
}

func ApproveSignDataCtx_QRegisterMetaType() int {
	return int(int32(C.ApproveSignDataCtx687eda_ApproveSignDataCtx687eda_QRegisterMetaType()))
}

func (ptr *ApproveSignDataCtx) QRegisterMetaType() int {
	return int(int32(C.ApproveSignDataCtx687eda_ApproveSignDataCtx687eda_QRegisterMetaType()))
}

func ApproveSignDataCtx_QRegisterMetaType2(typeName string) int {
	var typeNameC *C.char
	if typeName != "" {
		typeNameC = C.CString(typeName)
		defer C.free(unsafe.Pointer(typeNameC))
	}
	return int(int32(C.ApproveSignDataCtx687eda_ApproveSignDataCtx687eda_QRegisterMetaType2(typeNameC)))
}

func (ptr *ApproveSignDataCtx) QRegisterMetaType2(typeName string) int {
	var typeNameC *C.char
	if typeName != "" {
		typeNameC = C.CString(typeName)
		defer C.free(unsafe.Pointer(typeNameC))
	}
	return int(int32(C.ApproveSignDataCtx687eda_ApproveSignDataCtx687eda_QRegisterMetaType2(typeNameC)))
}

func ApproveSignDataCtx_QmlRegisterType() int {
	return int(int32(C.ApproveSignDataCtx687eda_ApproveSignDataCtx687eda_QmlRegisterType()))
}

func (ptr *ApproveSignDataCtx) QmlRegisterType() int {
	return int(int32(C.ApproveSignDataCtx687eda_ApproveSignDataCtx687eda_QmlRegisterType()))
}

func ApproveSignDataCtx_QmlRegisterType2(uri string, versionMajor int, versionMinor int, qmlName string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	return int(int32(C.ApproveSignDataCtx687eda_ApproveSignDataCtx687eda_QmlRegisterType2(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC)))
}

func (ptr *ApproveSignDataCtx) QmlRegisterType2(uri string, versionMajor int, versionMinor int, qmlName string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	return int(int32(C.ApproveSignDataCtx687eda_ApproveSignDataCtx687eda_QmlRegisterType2(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC)))
}

func ApproveSignDataCtx_QmlRegisterUncreatableType(uri string, versionMajor int, versionMinor int, qmlName string, message string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	var messageC *C.char
	if message != "" {
		messageC = C.CString(message)
		defer C.free(unsafe.Pointer(messageC))
	}
	return int(int32(C.ApproveSignDataCtx687eda_ApproveSignDataCtx687eda_QmlRegisterUncreatableType(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC, C.struct_Moc_PackedString{data: messageC, len: C.longlong(len(message))})))
}

func (ptr *ApproveSignDataCtx) QmlRegisterUncreatableType(uri string, versionMajor int, versionMinor int, qmlName string, message string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	var messageC *C.char
	if message != "" {
		messageC = C.CString(message)
		defer C.free(unsafe.Pointer(messageC))
	}
	return int(int32(C.ApproveSignDataCtx687eda_ApproveSignDataCtx687eda_QmlRegisterUncreatableType(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC, C.struct_Moc_PackedString{data: messageC, len: C.longlong(len(message))})))
}

func (ptr *ApproveSignDataCtx) __children_atList(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQObjectFromPointer(C.ApproveSignDataCtx687eda___children_atList(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *ApproveSignDataCtx) __children_setList(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveSignDataCtx687eda___children_setList(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *ApproveSignDataCtx) __children_newList() unsafe.Pointer {
	return C.ApproveSignDataCtx687eda___children_newList(ptr.Pointer())
}

func (ptr *ApproveSignDataCtx) __dynamicPropertyNames_atList(i int) *std_core.QByteArray {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQByteArrayFromPointer(C.ApproveSignDataCtx687eda___dynamicPropertyNames_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QByteArray).DestroyQByteArray)
		return tmpValue
	}
	return nil
}

func (ptr *ApproveSignDataCtx) __dynamicPropertyNames_setList(i std_core.QByteArray_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveSignDataCtx687eda___dynamicPropertyNames_setList(ptr.Pointer(), std_core.PointerFromQByteArray(i))
	}
}

func (ptr *ApproveSignDataCtx) __dynamicPropertyNames_newList() unsafe.Pointer {
	return C.ApproveSignDataCtx687eda___dynamicPropertyNames_newList(ptr.Pointer())
}

func (ptr *ApproveSignDataCtx) __findChildren_atList(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQObjectFromPointer(C.ApproveSignDataCtx687eda___findChildren_atList(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *ApproveSignDataCtx) __findChildren_setList(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveSignDataCtx687eda___findChildren_setList(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *ApproveSignDataCtx) __findChildren_newList() unsafe.Pointer {
	return C.ApproveSignDataCtx687eda___findChildren_newList(ptr.Pointer())
}

func (ptr *ApproveSignDataCtx) __findChildren_atList3(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQObjectFromPointer(C.ApproveSignDataCtx687eda___findChildren_atList3(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *ApproveSignDataCtx) __findChildren_setList3(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveSignDataCtx687eda___findChildren_setList3(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *ApproveSignDataCtx) __findChildren_newList3() unsafe.Pointer {
	return C.ApproveSignDataCtx687eda___findChildren_newList3(ptr.Pointer())
}

func NewApproveSignDataCtx(parent std_core.QObject_ITF) *ApproveSignDataCtx {
	ApproveSignDataCtx_QRegisterMetaType()
	tmpValue := NewApproveSignDataCtxFromPointer(C.ApproveSignDataCtx687eda_NewApproveSignDataCtx(std_core.PointerFromQObject(parent)))
	if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
		tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
	}
	return tmpValue
}

//export callbackApproveSignDataCtx687eda_DestroyApproveSignDataCtx
func callbackApproveSignDataCtx687eda_DestroyApproveSignDataCtx(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "~ApproveSignDataCtx"); signal != nil {
		(*(*func())(signal))()
	} else {
		NewApproveSignDataCtxFromPointer(ptr).DestroyApproveSignDataCtxDefault()
	}
}

func (ptr *ApproveSignDataCtx) ConnectDestroyApproveSignDataCtx(f func()) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "~ApproveSignDataCtx"); signal != nil {
			f := func() {
				(*(*func())(signal))()
				f()
			}
			qt.ConnectSignal(ptr.Pointer(), "~ApproveSignDataCtx", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "~ApproveSignDataCtx", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveSignDataCtx) DisconnectDestroyApproveSignDataCtx() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "~ApproveSignDataCtx")
	}
}

func (ptr *ApproveSignDataCtx) DestroyApproveSignDataCtx() {
	if ptr.Pointer() != nil {

		qt.SetFinalizer(ptr, nil)
		C.ApproveSignDataCtx687eda_DestroyApproveSignDataCtx(ptr.Pointer())
		ptr.SetPointer(nil)
	}
}

func (ptr *ApproveSignDataCtx) DestroyApproveSignDataCtxDefault() {
	if ptr.Pointer() != nil {

		qt.SetFinalizer(ptr, nil)
		C.ApproveSignDataCtx687eda_DestroyApproveSignDataCtxDefault(ptr.Pointer())
		ptr.SetPointer(nil)
	}
}

//export callbackApproveSignDataCtx687eda_ChildEvent
func callbackApproveSignDataCtx687eda_ChildEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "childEvent"); signal != nil {
		(*(*func(*std_core.QChildEvent))(signal))(std_core.NewQChildEventFromPointer(event))
	} else {
		NewApproveSignDataCtxFromPointer(ptr).ChildEventDefault(std_core.NewQChildEventFromPointer(event))
	}
}

func (ptr *ApproveSignDataCtx) ChildEventDefault(event std_core.QChildEvent_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveSignDataCtx687eda_ChildEventDefault(ptr.Pointer(), std_core.PointerFromQChildEvent(event))
	}
}

//export callbackApproveSignDataCtx687eda_ConnectNotify
func callbackApproveSignDataCtx687eda_ConnectNotify(ptr unsafe.Pointer, sign unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "connectNotify"); signal != nil {
		(*(*func(*std_core.QMetaMethod))(signal))(std_core.NewQMetaMethodFromPointer(sign))
	} else {
		NewApproveSignDataCtxFromPointer(ptr).ConnectNotifyDefault(std_core.NewQMetaMethodFromPointer(sign))
	}
}

func (ptr *ApproveSignDataCtx) ConnectNotifyDefault(sign std_core.QMetaMethod_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveSignDataCtx687eda_ConnectNotifyDefault(ptr.Pointer(), std_core.PointerFromQMetaMethod(sign))
	}
}

//export callbackApproveSignDataCtx687eda_CustomEvent
func callbackApproveSignDataCtx687eda_CustomEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "customEvent"); signal != nil {
		(*(*func(*std_core.QEvent))(signal))(std_core.NewQEventFromPointer(event))
	} else {
		NewApproveSignDataCtxFromPointer(ptr).CustomEventDefault(std_core.NewQEventFromPointer(event))
	}
}

func (ptr *ApproveSignDataCtx) CustomEventDefault(event std_core.QEvent_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveSignDataCtx687eda_CustomEventDefault(ptr.Pointer(), std_core.PointerFromQEvent(event))
	}
}

//export callbackApproveSignDataCtx687eda_DeleteLater
func callbackApproveSignDataCtx687eda_DeleteLater(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "deleteLater"); signal != nil {
		(*(*func())(signal))()
	} else {
		NewApproveSignDataCtxFromPointer(ptr).DeleteLaterDefault()
	}
}

func (ptr *ApproveSignDataCtx) DeleteLaterDefault() {
	if ptr.Pointer() != nil {

		qt.SetFinalizer(ptr, nil)
		C.ApproveSignDataCtx687eda_DeleteLaterDefault(ptr.Pointer())
	}
}

//export callbackApproveSignDataCtx687eda_Destroyed
func callbackApproveSignDataCtx687eda_Destroyed(ptr unsafe.Pointer, obj unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "destroyed"); signal != nil {
		(*(*func(*std_core.QObject))(signal))(std_core.NewQObjectFromPointer(obj))
	}
	qt.Unregister(ptr)

}

//export callbackApproveSignDataCtx687eda_DisconnectNotify
func callbackApproveSignDataCtx687eda_DisconnectNotify(ptr unsafe.Pointer, sign unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "disconnectNotify"); signal != nil {
		(*(*func(*std_core.QMetaMethod))(signal))(std_core.NewQMetaMethodFromPointer(sign))
	} else {
		NewApproveSignDataCtxFromPointer(ptr).DisconnectNotifyDefault(std_core.NewQMetaMethodFromPointer(sign))
	}
}

func (ptr *ApproveSignDataCtx) DisconnectNotifyDefault(sign std_core.QMetaMethod_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveSignDataCtx687eda_DisconnectNotifyDefault(ptr.Pointer(), std_core.PointerFromQMetaMethod(sign))
	}
}

//export callbackApproveSignDataCtx687eda_Event
func callbackApproveSignDataCtx687eda_Event(ptr unsafe.Pointer, e unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "event"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QEvent) bool)(signal))(std_core.NewQEventFromPointer(e)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewApproveSignDataCtxFromPointer(ptr).EventDefault(std_core.NewQEventFromPointer(e)))))
}

func (ptr *ApproveSignDataCtx) EventDefault(e std_core.QEvent_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.ApproveSignDataCtx687eda_EventDefault(ptr.Pointer(), std_core.PointerFromQEvent(e))) != 0
	}
	return false
}

//export callbackApproveSignDataCtx687eda_EventFilter
func callbackApproveSignDataCtx687eda_EventFilter(ptr unsafe.Pointer, watched unsafe.Pointer, event unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "eventFilter"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QObject, *std_core.QEvent) bool)(signal))(std_core.NewQObjectFromPointer(watched), std_core.NewQEventFromPointer(event)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewApproveSignDataCtxFromPointer(ptr).EventFilterDefault(std_core.NewQObjectFromPointer(watched), std_core.NewQEventFromPointer(event)))))
}

func (ptr *ApproveSignDataCtx) EventFilterDefault(watched std_core.QObject_ITF, event std_core.QEvent_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.ApproveSignDataCtx687eda_EventFilterDefault(ptr.Pointer(), std_core.PointerFromQObject(watched), std_core.PointerFromQEvent(event))) != 0
	}
	return false
}

//export callbackApproveSignDataCtx687eda_ObjectNameChanged
func callbackApproveSignDataCtx687eda_ObjectNameChanged(ptr unsafe.Pointer, objectName C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "objectNameChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(objectName))
	}

}

//export callbackApproveSignDataCtx687eda_TimerEvent
func callbackApproveSignDataCtx687eda_TimerEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "timerEvent"); signal != nil {
		(*(*func(*std_core.QTimerEvent))(signal))(std_core.NewQTimerEventFromPointer(event))
	} else {
		NewApproveSignDataCtxFromPointer(ptr).TimerEventDefault(std_core.NewQTimerEventFromPointer(event))
	}
}

func (ptr *ApproveSignDataCtx) TimerEventDefault(event std_core.QTimerEvent_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveSignDataCtx687eda_TimerEventDefault(ptr.Pointer(), std_core.PointerFromQTimerEvent(event))
	}
}

type ApproveTxCtx_ITF interface {
	std_core.QObject_ITF
	ApproveTxCtx_PTR() *ApproveTxCtx
}

func (ptr *ApproveTxCtx) ApproveTxCtx_PTR() *ApproveTxCtx {
	return ptr
}

func (ptr *ApproveTxCtx) Pointer() unsafe.Pointer {
	if ptr != nil {
		return ptr.QObject_PTR().Pointer()
	}
	return nil
}

func (ptr *ApproveTxCtx) SetPointer(p unsafe.Pointer) {
	if ptr != nil {
		ptr.QObject_PTR().SetPointer(p)
	}
}

func PointerFromApproveTxCtx(ptr ApproveTxCtx_ITF) unsafe.Pointer {
	if ptr != nil {
		return ptr.ApproveTxCtx_PTR().Pointer()
	}
	return nil
}

func NewApproveTxCtxFromPointer(ptr unsafe.Pointer) (n *ApproveTxCtx) {
	if gPtr, ok := qt.Receive(ptr); !ok {
		n = new(ApproveTxCtx)
		n.SetPointer(ptr)
	} else {
		switch deduced := gPtr.(type) {
		case *ApproveTxCtx:
			n = deduced

		case *std_core.QObject:
			n = &ApproveTxCtx{QObject: *deduced}

		default:
			n = new(ApproveTxCtx)
			n.SetPointer(ptr)
		}
	}
	return
}
func (this *ApproveTxCtx) Init() { this.init() }

//export callbackApproveTxCtx687eda_Constructor
func callbackApproveTxCtx687eda_Constructor(ptr unsafe.Pointer) {
	this := NewApproveTxCtxFromPointer(ptr)
	qt.Register(ptr, this)
	this.ConnectApprove(this.approve)
	this.ConnectReject(this.reject)
	this.ConnectCheckTxDiff(this.checkTxDiff)
	this.ConnectBack(this.back)
	this.ConnectEdited(this.edited)
	this.ConnectChangeValueUnit(this.changeValueUnit)
	this.ConnectChangeGasPriceUnit(this.changeGasPriceUnit)
	this.init()
}

//export callbackApproveTxCtx687eda_Approve
func callbackApproveTxCtx687eda_Approve(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "approve"); signal != nil {
		(*(*func())(signal))()
	}

}

func (ptr *ApproveTxCtx) ConnectApprove(f func()) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "approve") {
			C.ApproveTxCtx687eda_ConnectApprove(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "approve")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "approve"); signal != nil {
			f := func() {
				(*(*func())(signal))()
				f()
			}
			qt.ConnectSignal(ptr.Pointer(), "approve", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "approve", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectApprove() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectApprove(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "approve")
	}
}

func (ptr *ApproveTxCtx) Approve() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_Approve(ptr.Pointer())
	}
}

//export callbackApproveTxCtx687eda_Reject
func callbackApproveTxCtx687eda_Reject(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "reject"); signal != nil {
		(*(*func())(signal))()
	}

}

func (ptr *ApproveTxCtx) ConnectReject(f func()) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "reject") {
			C.ApproveTxCtx687eda_ConnectReject(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "reject")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "reject"); signal != nil {
			f := func() {
				(*(*func())(signal))()
				f()
			}
			qt.ConnectSignal(ptr.Pointer(), "reject", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "reject", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectReject() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectReject(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "reject")
	}
}

func (ptr *ApproveTxCtx) Reject() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_Reject(ptr.Pointer())
	}
}

//export callbackApproveTxCtx687eda_CheckTxDiff
func callbackApproveTxCtx687eda_CheckTxDiff(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "checkTxDiff"); signal != nil {
		(*(*func())(signal))()
	}

}

func (ptr *ApproveTxCtx) ConnectCheckTxDiff(f func()) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "checkTxDiff") {
			C.ApproveTxCtx687eda_ConnectCheckTxDiff(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "checkTxDiff")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "checkTxDiff"); signal != nil {
			f := func() {
				(*(*func())(signal))()
				f()
			}
			qt.ConnectSignal(ptr.Pointer(), "checkTxDiff", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "checkTxDiff", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectCheckTxDiff() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectCheckTxDiff(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "checkTxDiff")
	}
}

func (ptr *ApproveTxCtx) CheckTxDiff() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_CheckTxDiff(ptr.Pointer())
	}
}

//export callbackApproveTxCtx687eda_Back
func callbackApproveTxCtx687eda_Back(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "back"); signal != nil {
		(*(*func())(signal))()
	}

}

func (ptr *ApproveTxCtx) ConnectBack(f func()) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "back") {
			C.ApproveTxCtx687eda_ConnectBack(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "back")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "back"); signal != nil {
			f := func() {
				(*(*func())(signal))()
				f()
			}
			qt.ConnectSignal(ptr.Pointer(), "back", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "back", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectBack() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectBack(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "back")
	}
}

func (ptr *ApproveTxCtx) Back() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_Back(ptr.Pointer())
	}
}

//export callbackApproveTxCtx687eda_Edited
func callbackApproveTxCtx687eda_Edited(ptr unsafe.Pointer, s C.struct_Moc_PackedString, v C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "edited"); signal != nil {
		(*(*func(string, string))(signal))(cGoUnpackString(s), cGoUnpackString(v))
	}

}

func (ptr *ApproveTxCtx) ConnectEdited(f func(s string, v string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "edited") {
			C.ApproveTxCtx687eda_ConnectEdited(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "edited")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "edited"); signal != nil {
			f := func(s string, v string) {
				(*(*func(string, string))(signal))(s, v)
				f(s, v)
			}
			qt.ConnectSignal(ptr.Pointer(), "edited", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "edited", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectEdited() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectEdited(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "edited")
	}
}

func (ptr *ApproveTxCtx) Edited(s string, v string) {
	if ptr.Pointer() != nil {
		var sC *C.char
		if s != "" {
			sC = C.CString(s)
			defer C.free(unsafe.Pointer(sC))
		}
		var vC *C.char
		if v != "" {
			vC = C.CString(v)
			defer C.free(unsafe.Pointer(vC))
		}
		C.ApproveTxCtx687eda_Edited(ptr.Pointer(), C.struct_Moc_PackedString{data: sC, len: C.longlong(len(s))}, C.struct_Moc_PackedString{data: vC, len: C.longlong(len(v))})
	}
}

//export callbackApproveTxCtx687eda_ChangeValueUnit
func callbackApproveTxCtx687eda_ChangeValueUnit(ptr unsafe.Pointer, v C.int) {
	if signal := qt.GetSignal(ptr, "changeValueUnit"); signal != nil {
		(*(*func(int))(signal))(int(int32(v)))
	}

}

func (ptr *ApproveTxCtx) ConnectChangeValueUnit(f func(v int)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "changeValueUnit") {
			C.ApproveTxCtx687eda_ConnectChangeValueUnit(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "changeValueUnit")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "changeValueUnit"); signal != nil {
			f := func(v int) {
				(*(*func(int))(signal))(v)
				f(v)
			}
			qt.ConnectSignal(ptr.Pointer(), "changeValueUnit", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "changeValueUnit", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectChangeValueUnit() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectChangeValueUnit(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "changeValueUnit")
	}
}

func (ptr *ApproveTxCtx) ChangeValueUnit(v int) {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_ChangeValueUnit(ptr.Pointer(), C.int(int32(v)))
	}
}

//export callbackApproveTxCtx687eda_ChangeGasPriceUnit
func callbackApproveTxCtx687eda_ChangeGasPriceUnit(ptr unsafe.Pointer, v C.int) {
	if signal := qt.GetSignal(ptr, "changeGasPriceUnit"); signal != nil {
		(*(*func(int))(signal))(int(int32(v)))
	}

}

func (ptr *ApproveTxCtx) ConnectChangeGasPriceUnit(f func(v int)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "changeGasPriceUnit") {
			C.ApproveTxCtx687eda_ConnectChangeGasPriceUnit(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "changeGasPriceUnit")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "changeGasPriceUnit"); signal != nil {
			f := func(v int) {
				(*(*func(int))(signal))(v)
				f(v)
			}
			qt.ConnectSignal(ptr.Pointer(), "changeGasPriceUnit", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "changeGasPriceUnit", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectChangeGasPriceUnit() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectChangeGasPriceUnit(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "changeGasPriceUnit")
	}
}

func (ptr *ApproveTxCtx) ChangeGasPriceUnit(v int) {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_ChangeGasPriceUnit(ptr.Pointer(), C.int(int32(v)))
	}
}

//export callbackApproveTxCtx687eda_ValueUnit
func callbackApproveTxCtx687eda_ValueUnit(ptr unsafe.Pointer) C.int {
	if signal := qt.GetSignal(ptr, "valueUnit"); signal != nil {
		return C.int(int32((*(*func() int)(signal))()))
	}

	return C.int(int32(NewApproveTxCtxFromPointer(ptr).ValueUnitDefault()))
}

func (ptr *ApproveTxCtx) ConnectValueUnit(f func() int) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "valueUnit"); signal != nil {
			f := func() int {
				(*(*func() int)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "valueUnit", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "valueUnit", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectValueUnit() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "valueUnit")
	}
}

func (ptr *ApproveTxCtx) ValueUnit() int {
	if ptr.Pointer() != nil {
		return int(int32(C.ApproveTxCtx687eda_ValueUnit(ptr.Pointer())))
	}
	return 0
}

func (ptr *ApproveTxCtx) ValueUnitDefault() int {
	if ptr.Pointer() != nil {
		return int(int32(C.ApproveTxCtx687eda_ValueUnitDefault(ptr.Pointer())))
	}
	return 0
}

//export callbackApproveTxCtx687eda_SetValueUnit
func callbackApproveTxCtx687eda_SetValueUnit(ptr unsafe.Pointer, valueUnit C.int) {
	if signal := qt.GetSignal(ptr, "setValueUnit"); signal != nil {
		(*(*func(int))(signal))(int(int32(valueUnit)))
	} else {
		NewApproveTxCtxFromPointer(ptr).SetValueUnitDefault(int(int32(valueUnit)))
	}
}

func (ptr *ApproveTxCtx) ConnectSetValueUnit(f func(valueUnit int)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setValueUnit"); signal != nil {
			f := func(valueUnit int) {
				(*(*func(int))(signal))(valueUnit)
				f(valueUnit)
			}
			qt.ConnectSignal(ptr.Pointer(), "setValueUnit", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setValueUnit", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectSetValueUnit() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setValueUnit")
	}
}

func (ptr *ApproveTxCtx) SetValueUnit(valueUnit int) {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_SetValueUnit(ptr.Pointer(), C.int(int32(valueUnit)))
	}
}

func (ptr *ApproveTxCtx) SetValueUnitDefault(valueUnit int) {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_SetValueUnitDefault(ptr.Pointer(), C.int(int32(valueUnit)))
	}
}

//export callbackApproveTxCtx687eda_ValueUnitChanged
func callbackApproveTxCtx687eda_ValueUnitChanged(ptr unsafe.Pointer, valueUnit C.int) {
	if signal := qt.GetSignal(ptr, "valueUnitChanged"); signal != nil {
		(*(*func(int))(signal))(int(int32(valueUnit)))
	}

}

func (ptr *ApproveTxCtx) ConnectValueUnitChanged(f func(valueUnit int)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "valueUnitChanged") {
			C.ApproveTxCtx687eda_ConnectValueUnitChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "valueUnitChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "valueUnitChanged"); signal != nil {
			f := func(valueUnit int) {
				(*(*func(int))(signal))(valueUnit)
				f(valueUnit)
			}
			qt.ConnectSignal(ptr.Pointer(), "valueUnitChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "valueUnitChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectValueUnitChanged() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectValueUnitChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "valueUnitChanged")
	}
}

func (ptr *ApproveTxCtx) ValueUnitChanged(valueUnit int) {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_ValueUnitChanged(ptr.Pointer(), C.int(int32(valueUnit)))
	}
}

//export callbackApproveTxCtx687eda_Remote
func callbackApproveTxCtx687eda_Remote(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "remote"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveTxCtxFromPointer(ptr).RemoteDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveTxCtx) ConnectRemote(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "remote"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "remote", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "remote", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectRemote() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "remote")
	}
}

func (ptr *ApproveTxCtx) Remote() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_Remote(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveTxCtx) RemoteDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_RemoteDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveTxCtx687eda_SetRemote
func callbackApproveTxCtx687eda_SetRemote(ptr unsafe.Pointer, remote C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setRemote"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(remote))
	} else {
		NewApproveTxCtxFromPointer(ptr).SetRemoteDefault(cGoUnpackString(remote))
	}
}

func (ptr *ApproveTxCtx) ConnectSetRemote(f func(remote string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setRemote"); signal != nil {
			f := func(remote string) {
				(*(*func(string))(signal))(remote)
				f(remote)
			}
			qt.ConnectSignal(ptr.Pointer(), "setRemote", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setRemote", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectSetRemote() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setRemote")
	}
}

func (ptr *ApproveTxCtx) SetRemote(remote string) {
	if ptr.Pointer() != nil {
		var remoteC *C.char
		if remote != "" {
			remoteC = C.CString(remote)
			defer C.free(unsafe.Pointer(remoteC))
		}
		C.ApproveTxCtx687eda_SetRemote(ptr.Pointer(), C.struct_Moc_PackedString{data: remoteC, len: C.longlong(len(remote))})
	}
}

func (ptr *ApproveTxCtx) SetRemoteDefault(remote string) {
	if ptr.Pointer() != nil {
		var remoteC *C.char
		if remote != "" {
			remoteC = C.CString(remote)
			defer C.free(unsafe.Pointer(remoteC))
		}
		C.ApproveTxCtx687eda_SetRemoteDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: remoteC, len: C.longlong(len(remote))})
	}
}

//export callbackApproveTxCtx687eda_RemoteChanged
func callbackApproveTxCtx687eda_RemoteChanged(ptr unsafe.Pointer, remote C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "remoteChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(remote))
	}

}

func (ptr *ApproveTxCtx) ConnectRemoteChanged(f func(remote string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "remoteChanged") {
			C.ApproveTxCtx687eda_ConnectRemoteChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "remoteChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "remoteChanged"); signal != nil {
			f := func(remote string) {
				(*(*func(string))(signal))(remote)
				f(remote)
			}
			qt.ConnectSignal(ptr.Pointer(), "remoteChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "remoteChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectRemoteChanged() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectRemoteChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "remoteChanged")
	}
}

func (ptr *ApproveTxCtx) RemoteChanged(remote string) {
	if ptr.Pointer() != nil {
		var remoteC *C.char
		if remote != "" {
			remoteC = C.CString(remote)
			defer C.free(unsafe.Pointer(remoteC))
		}
		C.ApproveTxCtx687eda_RemoteChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: remoteC, len: C.longlong(len(remote))})
	}
}

//export callbackApproveTxCtx687eda_Transport
func callbackApproveTxCtx687eda_Transport(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "transport"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveTxCtxFromPointer(ptr).TransportDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveTxCtx) ConnectTransport(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "transport"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "transport", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "transport", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectTransport() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "transport")
	}
}

func (ptr *ApproveTxCtx) Transport() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_Transport(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveTxCtx) TransportDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_TransportDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveTxCtx687eda_SetTransport
func callbackApproveTxCtx687eda_SetTransport(ptr unsafe.Pointer, transport C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setTransport"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(transport))
	} else {
		NewApproveTxCtxFromPointer(ptr).SetTransportDefault(cGoUnpackString(transport))
	}
}

func (ptr *ApproveTxCtx) ConnectSetTransport(f func(transport string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setTransport"); signal != nil {
			f := func(transport string) {
				(*(*func(string))(signal))(transport)
				f(transport)
			}
			qt.ConnectSignal(ptr.Pointer(), "setTransport", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setTransport", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectSetTransport() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setTransport")
	}
}

func (ptr *ApproveTxCtx) SetTransport(transport string) {
	if ptr.Pointer() != nil {
		var transportC *C.char
		if transport != "" {
			transportC = C.CString(transport)
			defer C.free(unsafe.Pointer(transportC))
		}
		C.ApproveTxCtx687eda_SetTransport(ptr.Pointer(), C.struct_Moc_PackedString{data: transportC, len: C.longlong(len(transport))})
	}
}

func (ptr *ApproveTxCtx) SetTransportDefault(transport string) {
	if ptr.Pointer() != nil {
		var transportC *C.char
		if transport != "" {
			transportC = C.CString(transport)
			defer C.free(unsafe.Pointer(transportC))
		}
		C.ApproveTxCtx687eda_SetTransportDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: transportC, len: C.longlong(len(transport))})
	}
}

//export callbackApproveTxCtx687eda_TransportChanged
func callbackApproveTxCtx687eda_TransportChanged(ptr unsafe.Pointer, transport C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "transportChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(transport))
	}

}

func (ptr *ApproveTxCtx) ConnectTransportChanged(f func(transport string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "transportChanged") {
			C.ApproveTxCtx687eda_ConnectTransportChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "transportChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "transportChanged"); signal != nil {
			f := func(transport string) {
				(*(*func(string))(signal))(transport)
				f(transport)
			}
			qt.ConnectSignal(ptr.Pointer(), "transportChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "transportChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectTransportChanged() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectTransportChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "transportChanged")
	}
}

func (ptr *ApproveTxCtx) TransportChanged(transport string) {
	if ptr.Pointer() != nil {
		var transportC *C.char
		if transport != "" {
			transportC = C.CString(transport)
			defer C.free(unsafe.Pointer(transportC))
		}
		C.ApproveTxCtx687eda_TransportChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: transportC, len: C.longlong(len(transport))})
	}
}

//export callbackApproveTxCtx687eda_Endpoint
func callbackApproveTxCtx687eda_Endpoint(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "endpoint"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveTxCtxFromPointer(ptr).EndpointDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveTxCtx) ConnectEndpoint(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "endpoint"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "endpoint", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "endpoint", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectEndpoint() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "endpoint")
	}
}

func (ptr *ApproveTxCtx) Endpoint() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_Endpoint(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveTxCtx) EndpointDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_EndpointDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveTxCtx687eda_SetEndpoint
func callbackApproveTxCtx687eda_SetEndpoint(ptr unsafe.Pointer, endpoint C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setEndpoint"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(endpoint))
	} else {
		NewApproveTxCtxFromPointer(ptr).SetEndpointDefault(cGoUnpackString(endpoint))
	}
}

func (ptr *ApproveTxCtx) ConnectSetEndpoint(f func(endpoint string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setEndpoint"); signal != nil {
			f := func(endpoint string) {
				(*(*func(string))(signal))(endpoint)
				f(endpoint)
			}
			qt.ConnectSignal(ptr.Pointer(), "setEndpoint", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setEndpoint", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectSetEndpoint() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setEndpoint")
	}
}

func (ptr *ApproveTxCtx) SetEndpoint(endpoint string) {
	if ptr.Pointer() != nil {
		var endpointC *C.char
		if endpoint != "" {
			endpointC = C.CString(endpoint)
			defer C.free(unsafe.Pointer(endpointC))
		}
		C.ApproveTxCtx687eda_SetEndpoint(ptr.Pointer(), C.struct_Moc_PackedString{data: endpointC, len: C.longlong(len(endpoint))})
	}
}

func (ptr *ApproveTxCtx) SetEndpointDefault(endpoint string) {
	if ptr.Pointer() != nil {
		var endpointC *C.char
		if endpoint != "" {
			endpointC = C.CString(endpoint)
			defer C.free(unsafe.Pointer(endpointC))
		}
		C.ApproveTxCtx687eda_SetEndpointDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: endpointC, len: C.longlong(len(endpoint))})
	}
}

//export callbackApproveTxCtx687eda_EndpointChanged
func callbackApproveTxCtx687eda_EndpointChanged(ptr unsafe.Pointer, endpoint C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "endpointChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(endpoint))
	}

}

func (ptr *ApproveTxCtx) ConnectEndpointChanged(f func(endpoint string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "endpointChanged") {
			C.ApproveTxCtx687eda_ConnectEndpointChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "endpointChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "endpointChanged"); signal != nil {
			f := func(endpoint string) {
				(*(*func(string))(signal))(endpoint)
				f(endpoint)
			}
			qt.ConnectSignal(ptr.Pointer(), "endpointChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "endpointChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectEndpointChanged() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectEndpointChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "endpointChanged")
	}
}

func (ptr *ApproveTxCtx) EndpointChanged(endpoint string) {
	if ptr.Pointer() != nil {
		var endpointC *C.char
		if endpoint != "" {
			endpointC = C.CString(endpoint)
			defer C.free(unsafe.Pointer(endpointC))
		}
		C.ApproveTxCtx687eda_EndpointChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: endpointC, len: C.longlong(len(endpoint))})
	}
}

//export callbackApproveTxCtx687eda_Data
func callbackApproveTxCtx687eda_Data(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "data"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveTxCtxFromPointer(ptr).DataDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveTxCtx) ConnectData(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "data"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "data", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "data", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectData() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "data")
	}
}

func (ptr *ApproveTxCtx) Data() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_Data(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveTxCtx) DataDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_DataDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveTxCtx687eda_SetData
func callbackApproveTxCtx687eda_SetData(ptr unsafe.Pointer, data C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setData"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(data))
	} else {
		NewApproveTxCtxFromPointer(ptr).SetDataDefault(cGoUnpackString(data))
	}
}

func (ptr *ApproveTxCtx) ConnectSetData(f func(data string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setData"); signal != nil {
			f := func(data string) {
				(*(*func(string))(signal))(data)
				f(data)
			}
			qt.ConnectSignal(ptr.Pointer(), "setData", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setData", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectSetData() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setData")
	}
}

func (ptr *ApproveTxCtx) SetData(data string) {
	if ptr.Pointer() != nil {
		var dataC *C.char
		if data != "" {
			dataC = C.CString(data)
			defer C.free(unsafe.Pointer(dataC))
		}
		C.ApproveTxCtx687eda_SetData(ptr.Pointer(), C.struct_Moc_PackedString{data: dataC, len: C.longlong(len(data))})
	}
}

func (ptr *ApproveTxCtx) SetDataDefault(data string) {
	if ptr.Pointer() != nil {
		var dataC *C.char
		if data != "" {
			dataC = C.CString(data)
			defer C.free(unsafe.Pointer(dataC))
		}
		C.ApproveTxCtx687eda_SetDataDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: dataC, len: C.longlong(len(data))})
	}
}

//export callbackApproveTxCtx687eda_DataChanged
func callbackApproveTxCtx687eda_DataChanged(ptr unsafe.Pointer, data C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "dataChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(data))
	}

}

func (ptr *ApproveTxCtx) ConnectDataChanged(f func(data string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "dataChanged") {
			C.ApproveTxCtx687eda_ConnectDataChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "dataChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "dataChanged"); signal != nil {
			f := func(data string) {
				(*(*func(string))(signal))(data)
				f(data)
			}
			qt.ConnectSignal(ptr.Pointer(), "dataChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "dataChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectDataChanged() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectDataChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "dataChanged")
	}
}

func (ptr *ApproveTxCtx) DataChanged(data string) {
	if ptr.Pointer() != nil {
		var dataC *C.char
		if data != "" {
			dataC = C.CString(data)
			defer C.free(unsafe.Pointer(dataC))
		}
		C.ApproveTxCtx687eda_DataChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: dataC, len: C.longlong(len(data))})
	}
}

//export callbackApproveTxCtx687eda_From
func callbackApproveTxCtx687eda_From(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "from"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveTxCtxFromPointer(ptr).FromDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveTxCtx) ConnectFrom(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "from"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "from", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "from", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectFrom() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "from")
	}
}

func (ptr *ApproveTxCtx) From() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_From(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveTxCtx) FromDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_FromDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveTxCtx687eda_SetFrom
func callbackApproveTxCtx687eda_SetFrom(ptr unsafe.Pointer, from C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setFrom"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(from))
	} else {
		NewApproveTxCtxFromPointer(ptr).SetFromDefault(cGoUnpackString(from))
	}
}

func (ptr *ApproveTxCtx) ConnectSetFrom(f func(from string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setFrom"); signal != nil {
			f := func(from string) {
				(*(*func(string))(signal))(from)
				f(from)
			}
			qt.ConnectSignal(ptr.Pointer(), "setFrom", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setFrom", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectSetFrom() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setFrom")
	}
}

func (ptr *ApproveTxCtx) SetFrom(from string) {
	if ptr.Pointer() != nil {
		var fromC *C.char
		if from != "" {
			fromC = C.CString(from)
			defer C.free(unsafe.Pointer(fromC))
		}
		C.ApproveTxCtx687eda_SetFrom(ptr.Pointer(), C.struct_Moc_PackedString{data: fromC, len: C.longlong(len(from))})
	}
}

func (ptr *ApproveTxCtx) SetFromDefault(from string) {
	if ptr.Pointer() != nil {
		var fromC *C.char
		if from != "" {
			fromC = C.CString(from)
			defer C.free(unsafe.Pointer(fromC))
		}
		C.ApproveTxCtx687eda_SetFromDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: fromC, len: C.longlong(len(from))})
	}
}

//export callbackApproveTxCtx687eda_FromChanged
func callbackApproveTxCtx687eda_FromChanged(ptr unsafe.Pointer, from C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "fromChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(from))
	}

}

func (ptr *ApproveTxCtx) ConnectFromChanged(f func(from string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "fromChanged") {
			C.ApproveTxCtx687eda_ConnectFromChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "fromChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "fromChanged"); signal != nil {
			f := func(from string) {
				(*(*func(string))(signal))(from)
				f(from)
			}
			qt.ConnectSignal(ptr.Pointer(), "fromChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "fromChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectFromChanged() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectFromChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "fromChanged")
	}
}

func (ptr *ApproveTxCtx) FromChanged(from string) {
	if ptr.Pointer() != nil {
		var fromC *C.char
		if from != "" {
			fromC = C.CString(from)
			defer C.free(unsafe.Pointer(fromC))
		}
		C.ApproveTxCtx687eda_FromChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: fromC, len: C.longlong(len(from))})
	}
}

//export callbackApproveTxCtx687eda_FromWarning
func callbackApproveTxCtx687eda_FromWarning(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "fromWarning"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveTxCtxFromPointer(ptr).FromWarningDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveTxCtx) ConnectFromWarning(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "fromWarning"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "fromWarning", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "fromWarning", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectFromWarning() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "fromWarning")
	}
}

func (ptr *ApproveTxCtx) FromWarning() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_FromWarning(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveTxCtx) FromWarningDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_FromWarningDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveTxCtx687eda_SetFromWarning
func callbackApproveTxCtx687eda_SetFromWarning(ptr unsafe.Pointer, fromWarning C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setFromWarning"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(fromWarning))
	} else {
		NewApproveTxCtxFromPointer(ptr).SetFromWarningDefault(cGoUnpackString(fromWarning))
	}
}

func (ptr *ApproveTxCtx) ConnectSetFromWarning(f func(fromWarning string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setFromWarning"); signal != nil {
			f := func(fromWarning string) {
				(*(*func(string))(signal))(fromWarning)
				f(fromWarning)
			}
			qt.ConnectSignal(ptr.Pointer(), "setFromWarning", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setFromWarning", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectSetFromWarning() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setFromWarning")
	}
}

func (ptr *ApproveTxCtx) SetFromWarning(fromWarning string) {
	if ptr.Pointer() != nil {
		var fromWarningC *C.char
		if fromWarning != "" {
			fromWarningC = C.CString(fromWarning)
			defer C.free(unsafe.Pointer(fromWarningC))
		}
		C.ApproveTxCtx687eda_SetFromWarning(ptr.Pointer(), C.struct_Moc_PackedString{data: fromWarningC, len: C.longlong(len(fromWarning))})
	}
}

func (ptr *ApproveTxCtx) SetFromWarningDefault(fromWarning string) {
	if ptr.Pointer() != nil {
		var fromWarningC *C.char
		if fromWarning != "" {
			fromWarningC = C.CString(fromWarning)
			defer C.free(unsafe.Pointer(fromWarningC))
		}
		C.ApproveTxCtx687eda_SetFromWarningDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: fromWarningC, len: C.longlong(len(fromWarning))})
	}
}

//export callbackApproveTxCtx687eda_FromWarningChanged
func callbackApproveTxCtx687eda_FromWarningChanged(ptr unsafe.Pointer, fromWarning C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "fromWarningChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(fromWarning))
	}

}

func (ptr *ApproveTxCtx) ConnectFromWarningChanged(f func(fromWarning string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "fromWarningChanged") {
			C.ApproveTxCtx687eda_ConnectFromWarningChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "fromWarningChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "fromWarningChanged"); signal != nil {
			f := func(fromWarning string) {
				(*(*func(string))(signal))(fromWarning)
				f(fromWarning)
			}
			qt.ConnectSignal(ptr.Pointer(), "fromWarningChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "fromWarningChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectFromWarningChanged() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectFromWarningChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "fromWarningChanged")
	}
}

func (ptr *ApproveTxCtx) FromWarningChanged(fromWarning string) {
	if ptr.Pointer() != nil {
		var fromWarningC *C.char
		if fromWarning != "" {
			fromWarningC = C.CString(fromWarning)
			defer C.free(unsafe.Pointer(fromWarningC))
		}
		C.ApproveTxCtx687eda_FromWarningChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: fromWarningC, len: C.longlong(len(fromWarning))})
	}
}

//export callbackApproveTxCtx687eda_IsFromVisible
func callbackApproveTxCtx687eda_IsFromVisible(ptr unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "isFromVisible"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func() bool)(signal))())))
	}

	return C.char(int8(qt.GoBoolToInt(NewApproveTxCtxFromPointer(ptr).IsFromVisibleDefault())))
}

func (ptr *ApproveTxCtx) ConnectIsFromVisible(f func() bool) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "isFromVisible"); signal != nil {
			f := func() bool {
				(*(*func() bool)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "isFromVisible", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "isFromVisible", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectIsFromVisible() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "isFromVisible")
	}
}

func (ptr *ApproveTxCtx) IsFromVisible() bool {
	if ptr.Pointer() != nil {
		return int8(C.ApproveTxCtx687eda_IsFromVisible(ptr.Pointer())) != 0
	}
	return false
}

func (ptr *ApproveTxCtx) IsFromVisibleDefault() bool {
	if ptr.Pointer() != nil {
		return int8(C.ApproveTxCtx687eda_IsFromVisibleDefault(ptr.Pointer())) != 0
	}
	return false
}

//export callbackApproveTxCtx687eda_SetFromVisible
func callbackApproveTxCtx687eda_SetFromVisible(ptr unsafe.Pointer, fromVisible C.char) {
	if signal := qt.GetSignal(ptr, "setFromVisible"); signal != nil {
		(*(*func(bool))(signal))(int8(fromVisible) != 0)
	} else {
		NewApproveTxCtxFromPointer(ptr).SetFromVisibleDefault(int8(fromVisible) != 0)
	}
}

func (ptr *ApproveTxCtx) ConnectSetFromVisible(f func(fromVisible bool)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setFromVisible"); signal != nil {
			f := func(fromVisible bool) {
				(*(*func(bool))(signal))(fromVisible)
				f(fromVisible)
			}
			qt.ConnectSignal(ptr.Pointer(), "setFromVisible", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setFromVisible", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectSetFromVisible() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setFromVisible")
	}
}

func (ptr *ApproveTxCtx) SetFromVisible(fromVisible bool) {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_SetFromVisible(ptr.Pointer(), C.char(int8(qt.GoBoolToInt(fromVisible))))
	}
}

func (ptr *ApproveTxCtx) SetFromVisibleDefault(fromVisible bool) {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_SetFromVisibleDefault(ptr.Pointer(), C.char(int8(qt.GoBoolToInt(fromVisible))))
	}
}

//export callbackApproveTxCtx687eda_FromVisibleChanged
func callbackApproveTxCtx687eda_FromVisibleChanged(ptr unsafe.Pointer, fromVisible C.char) {
	if signal := qt.GetSignal(ptr, "fromVisibleChanged"); signal != nil {
		(*(*func(bool))(signal))(int8(fromVisible) != 0)
	}

}

func (ptr *ApproveTxCtx) ConnectFromVisibleChanged(f func(fromVisible bool)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "fromVisibleChanged") {
			C.ApproveTxCtx687eda_ConnectFromVisibleChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "fromVisibleChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "fromVisibleChanged"); signal != nil {
			f := func(fromVisible bool) {
				(*(*func(bool))(signal))(fromVisible)
				f(fromVisible)
			}
			qt.ConnectSignal(ptr.Pointer(), "fromVisibleChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "fromVisibleChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectFromVisibleChanged() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectFromVisibleChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "fromVisibleChanged")
	}
}

func (ptr *ApproveTxCtx) FromVisibleChanged(fromVisible bool) {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_FromVisibleChanged(ptr.Pointer(), C.char(int8(qt.GoBoolToInt(fromVisible))))
	}
}

//export callbackApproveTxCtx687eda_To
func callbackApproveTxCtx687eda_To(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "to"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveTxCtxFromPointer(ptr).ToDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveTxCtx) ConnectTo(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "to"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "to", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "to", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectTo() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "to")
	}
}

func (ptr *ApproveTxCtx) To() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_To(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveTxCtx) ToDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_ToDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveTxCtx687eda_SetTo
func callbackApproveTxCtx687eda_SetTo(ptr unsafe.Pointer, to C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setTo"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(to))
	} else {
		NewApproveTxCtxFromPointer(ptr).SetToDefault(cGoUnpackString(to))
	}
}

func (ptr *ApproveTxCtx) ConnectSetTo(f func(to string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setTo"); signal != nil {
			f := func(to string) {
				(*(*func(string))(signal))(to)
				f(to)
			}
			qt.ConnectSignal(ptr.Pointer(), "setTo", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setTo", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectSetTo() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setTo")
	}
}

func (ptr *ApproveTxCtx) SetTo(to string) {
	if ptr.Pointer() != nil {
		var toC *C.char
		if to != "" {
			toC = C.CString(to)
			defer C.free(unsafe.Pointer(toC))
		}
		C.ApproveTxCtx687eda_SetTo(ptr.Pointer(), C.struct_Moc_PackedString{data: toC, len: C.longlong(len(to))})
	}
}

func (ptr *ApproveTxCtx) SetToDefault(to string) {
	if ptr.Pointer() != nil {
		var toC *C.char
		if to != "" {
			toC = C.CString(to)
			defer C.free(unsafe.Pointer(toC))
		}
		C.ApproveTxCtx687eda_SetToDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: toC, len: C.longlong(len(to))})
	}
}

//export callbackApproveTxCtx687eda_ToChanged
func callbackApproveTxCtx687eda_ToChanged(ptr unsafe.Pointer, to C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "toChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(to))
	}

}

func (ptr *ApproveTxCtx) ConnectToChanged(f func(to string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "toChanged") {
			C.ApproveTxCtx687eda_ConnectToChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "toChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "toChanged"); signal != nil {
			f := func(to string) {
				(*(*func(string))(signal))(to)
				f(to)
			}
			qt.ConnectSignal(ptr.Pointer(), "toChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "toChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectToChanged() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectToChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "toChanged")
	}
}

func (ptr *ApproveTxCtx) ToChanged(to string) {
	if ptr.Pointer() != nil {
		var toC *C.char
		if to != "" {
			toC = C.CString(to)
			defer C.free(unsafe.Pointer(toC))
		}
		C.ApproveTxCtx687eda_ToChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: toC, len: C.longlong(len(to))})
	}
}

//export callbackApproveTxCtx687eda_ToWarning
func callbackApproveTxCtx687eda_ToWarning(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "toWarning"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveTxCtxFromPointer(ptr).ToWarningDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveTxCtx) ConnectToWarning(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "toWarning"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "toWarning", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "toWarning", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectToWarning() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "toWarning")
	}
}

func (ptr *ApproveTxCtx) ToWarning() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_ToWarning(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveTxCtx) ToWarningDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_ToWarningDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveTxCtx687eda_SetToWarning
func callbackApproveTxCtx687eda_SetToWarning(ptr unsafe.Pointer, toWarning C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setToWarning"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(toWarning))
	} else {
		NewApproveTxCtxFromPointer(ptr).SetToWarningDefault(cGoUnpackString(toWarning))
	}
}

func (ptr *ApproveTxCtx) ConnectSetToWarning(f func(toWarning string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setToWarning"); signal != nil {
			f := func(toWarning string) {
				(*(*func(string))(signal))(toWarning)
				f(toWarning)
			}
			qt.ConnectSignal(ptr.Pointer(), "setToWarning", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setToWarning", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectSetToWarning() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setToWarning")
	}
}

func (ptr *ApproveTxCtx) SetToWarning(toWarning string) {
	if ptr.Pointer() != nil {
		var toWarningC *C.char
		if toWarning != "" {
			toWarningC = C.CString(toWarning)
			defer C.free(unsafe.Pointer(toWarningC))
		}
		C.ApproveTxCtx687eda_SetToWarning(ptr.Pointer(), C.struct_Moc_PackedString{data: toWarningC, len: C.longlong(len(toWarning))})
	}
}

func (ptr *ApproveTxCtx) SetToWarningDefault(toWarning string) {
	if ptr.Pointer() != nil {
		var toWarningC *C.char
		if toWarning != "" {
			toWarningC = C.CString(toWarning)
			defer C.free(unsafe.Pointer(toWarningC))
		}
		C.ApproveTxCtx687eda_SetToWarningDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: toWarningC, len: C.longlong(len(toWarning))})
	}
}

//export callbackApproveTxCtx687eda_ToWarningChanged
func callbackApproveTxCtx687eda_ToWarningChanged(ptr unsafe.Pointer, toWarning C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "toWarningChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(toWarning))
	}

}

func (ptr *ApproveTxCtx) ConnectToWarningChanged(f func(toWarning string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "toWarningChanged") {
			C.ApproveTxCtx687eda_ConnectToWarningChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "toWarningChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "toWarningChanged"); signal != nil {
			f := func(toWarning string) {
				(*(*func(string))(signal))(toWarning)
				f(toWarning)
			}
			qt.ConnectSignal(ptr.Pointer(), "toWarningChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "toWarningChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectToWarningChanged() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectToWarningChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "toWarningChanged")
	}
}

func (ptr *ApproveTxCtx) ToWarningChanged(toWarning string) {
	if ptr.Pointer() != nil {
		var toWarningC *C.char
		if toWarning != "" {
			toWarningC = C.CString(toWarning)
			defer C.free(unsafe.Pointer(toWarningC))
		}
		C.ApproveTxCtx687eda_ToWarningChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: toWarningC, len: C.longlong(len(toWarning))})
	}
}

//export callbackApproveTxCtx687eda_IsToVisible
func callbackApproveTxCtx687eda_IsToVisible(ptr unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "isToVisible"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func() bool)(signal))())))
	}

	return C.char(int8(qt.GoBoolToInt(NewApproveTxCtxFromPointer(ptr).IsToVisibleDefault())))
}

func (ptr *ApproveTxCtx) ConnectIsToVisible(f func() bool) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "isToVisible"); signal != nil {
			f := func() bool {
				(*(*func() bool)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "isToVisible", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "isToVisible", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectIsToVisible() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "isToVisible")
	}
}

func (ptr *ApproveTxCtx) IsToVisible() bool {
	if ptr.Pointer() != nil {
		return int8(C.ApproveTxCtx687eda_IsToVisible(ptr.Pointer())) != 0
	}
	return false
}

func (ptr *ApproveTxCtx) IsToVisibleDefault() bool {
	if ptr.Pointer() != nil {
		return int8(C.ApproveTxCtx687eda_IsToVisibleDefault(ptr.Pointer())) != 0
	}
	return false
}

//export callbackApproveTxCtx687eda_SetToVisible
func callbackApproveTxCtx687eda_SetToVisible(ptr unsafe.Pointer, toVisible C.char) {
	if signal := qt.GetSignal(ptr, "setToVisible"); signal != nil {
		(*(*func(bool))(signal))(int8(toVisible) != 0)
	} else {
		NewApproveTxCtxFromPointer(ptr).SetToVisibleDefault(int8(toVisible) != 0)
	}
}

func (ptr *ApproveTxCtx) ConnectSetToVisible(f func(toVisible bool)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setToVisible"); signal != nil {
			f := func(toVisible bool) {
				(*(*func(bool))(signal))(toVisible)
				f(toVisible)
			}
			qt.ConnectSignal(ptr.Pointer(), "setToVisible", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setToVisible", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectSetToVisible() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setToVisible")
	}
}

func (ptr *ApproveTxCtx) SetToVisible(toVisible bool) {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_SetToVisible(ptr.Pointer(), C.char(int8(qt.GoBoolToInt(toVisible))))
	}
}

func (ptr *ApproveTxCtx) SetToVisibleDefault(toVisible bool) {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_SetToVisibleDefault(ptr.Pointer(), C.char(int8(qt.GoBoolToInt(toVisible))))
	}
}

//export callbackApproveTxCtx687eda_ToVisibleChanged
func callbackApproveTxCtx687eda_ToVisibleChanged(ptr unsafe.Pointer, toVisible C.char) {
	if signal := qt.GetSignal(ptr, "toVisibleChanged"); signal != nil {
		(*(*func(bool))(signal))(int8(toVisible) != 0)
	}

}

func (ptr *ApproveTxCtx) ConnectToVisibleChanged(f func(toVisible bool)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "toVisibleChanged") {
			C.ApproveTxCtx687eda_ConnectToVisibleChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "toVisibleChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "toVisibleChanged"); signal != nil {
			f := func(toVisible bool) {
				(*(*func(bool))(signal))(toVisible)
				f(toVisible)
			}
			qt.ConnectSignal(ptr.Pointer(), "toVisibleChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "toVisibleChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectToVisibleChanged() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectToVisibleChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "toVisibleChanged")
	}
}

func (ptr *ApproveTxCtx) ToVisibleChanged(toVisible bool) {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_ToVisibleChanged(ptr.Pointer(), C.char(int8(qt.GoBoolToInt(toVisible))))
	}
}

//export callbackApproveTxCtx687eda_Gas
func callbackApproveTxCtx687eda_Gas(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "gas"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveTxCtxFromPointer(ptr).GasDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveTxCtx) ConnectGas(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "gas"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "gas", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "gas", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectGas() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "gas")
	}
}

func (ptr *ApproveTxCtx) Gas() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_Gas(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveTxCtx) GasDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_GasDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveTxCtx687eda_SetGas
func callbackApproveTxCtx687eda_SetGas(ptr unsafe.Pointer, gas C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setGas"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(gas))
	} else {
		NewApproveTxCtxFromPointer(ptr).SetGasDefault(cGoUnpackString(gas))
	}
}

func (ptr *ApproveTxCtx) ConnectSetGas(f func(gas string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setGas"); signal != nil {
			f := func(gas string) {
				(*(*func(string))(signal))(gas)
				f(gas)
			}
			qt.ConnectSignal(ptr.Pointer(), "setGas", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setGas", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectSetGas() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setGas")
	}
}

func (ptr *ApproveTxCtx) SetGas(gas string) {
	if ptr.Pointer() != nil {
		var gasC *C.char
		if gas != "" {
			gasC = C.CString(gas)
			defer C.free(unsafe.Pointer(gasC))
		}
		C.ApproveTxCtx687eda_SetGas(ptr.Pointer(), C.struct_Moc_PackedString{data: gasC, len: C.longlong(len(gas))})
	}
}

func (ptr *ApproveTxCtx) SetGasDefault(gas string) {
	if ptr.Pointer() != nil {
		var gasC *C.char
		if gas != "" {
			gasC = C.CString(gas)
			defer C.free(unsafe.Pointer(gasC))
		}
		C.ApproveTxCtx687eda_SetGasDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: gasC, len: C.longlong(len(gas))})
	}
}

//export callbackApproveTxCtx687eda_GasChanged
func callbackApproveTxCtx687eda_GasChanged(ptr unsafe.Pointer, gas C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "gasChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(gas))
	}

}

func (ptr *ApproveTxCtx) ConnectGasChanged(f func(gas string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "gasChanged") {
			C.ApproveTxCtx687eda_ConnectGasChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "gasChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "gasChanged"); signal != nil {
			f := func(gas string) {
				(*(*func(string))(signal))(gas)
				f(gas)
			}
			qt.ConnectSignal(ptr.Pointer(), "gasChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "gasChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectGasChanged() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectGasChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "gasChanged")
	}
}

func (ptr *ApproveTxCtx) GasChanged(gas string) {
	if ptr.Pointer() != nil {
		var gasC *C.char
		if gas != "" {
			gasC = C.CString(gas)
			defer C.free(unsafe.Pointer(gasC))
		}
		C.ApproveTxCtx687eda_GasChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: gasC, len: C.longlong(len(gas))})
	}
}

//export callbackApproveTxCtx687eda_GasPrice
func callbackApproveTxCtx687eda_GasPrice(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "gasPrice"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveTxCtxFromPointer(ptr).GasPriceDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveTxCtx) ConnectGasPrice(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "gasPrice"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "gasPrice", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "gasPrice", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectGasPrice() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "gasPrice")
	}
}

func (ptr *ApproveTxCtx) GasPrice() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_GasPrice(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveTxCtx) GasPriceDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_GasPriceDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveTxCtx687eda_SetGasPrice
func callbackApproveTxCtx687eda_SetGasPrice(ptr unsafe.Pointer, gasPrice C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setGasPrice"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(gasPrice))
	} else {
		NewApproveTxCtxFromPointer(ptr).SetGasPriceDefault(cGoUnpackString(gasPrice))
	}
}

func (ptr *ApproveTxCtx) ConnectSetGasPrice(f func(gasPrice string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setGasPrice"); signal != nil {
			f := func(gasPrice string) {
				(*(*func(string))(signal))(gasPrice)
				f(gasPrice)
			}
			qt.ConnectSignal(ptr.Pointer(), "setGasPrice", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setGasPrice", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectSetGasPrice() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setGasPrice")
	}
}

func (ptr *ApproveTxCtx) SetGasPrice(gasPrice string) {
	if ptr.Pointer() != nil {
		var gasPriceC *C.char
		if gasPrice != "" {
			gasPriceC = C.CString(gasPrice)
			defer C.free(unsafe.Pointer(gasPriceC))
		}
		C.ApproveTxCtx687eda_SetGasPrice(ptr.Pointer(), C.struct_Moc_PackedString{data: gasPriceC, len: C.longlong(len(gasPrice))})
	}
}

func (ptr *ApproveTxCtx) SetGasPriceDefault(gasPrice string) {
	if ptr.Pointer() != nil {
		var gasPriceC *C.char
		if gasPrice != "" {
			gasPriceC = C.CString(gasPrice)
			defer C.free(unsafe.Pointer(gasPriceC))
		}
		C.ApproveTxCtx687eda_SetGasPriceDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: gasPriceC, len: C.longlong(len(gasPrice))})
	}
}

//export callbackApproveTxCtx687eda_GasPriceChanged
func callbackApproveTxCtx687eda_GasPriceChanged(ptr unsafe.Pointer, gasPrice C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "gasPriceChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(gasPrice))
	}

}

func (ptr *ApproveTxCtx) ConnectGasPriceChanged(f func(gasPrice string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "gasPriceChanged") {
			C.ApproveTxCtx687eda_ConnectGasPriceChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "gasPriceChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "gasPriceChanged"); signal != nil {
			f := func(gasPrice string) {
				(*(*func(string))(signal))(gasPrice)
				f(gasPrice)
			}
			qt.ConnectSignal(ptr.Pointer(), "gasPriceChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "gasPriceChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectGasPriceChanged() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectGasPriceChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "gasPriceChanged")
	}
}

func (ptr *ApproveTxCtx) GasPriceChanged(gasPrice string) {
	if ptr.Pointer() != nil {
		var gasPriceC *C.char
		if gasPrice != "" {
			gasPriceC = C.CString(gasPrice)
			defer C.free(unsafe.Pointer(gasPriceC))
		}
		C.ApproveTxCtx687eda_GasPriceChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: gasPriceC, len: C.longlong(len(gasPrice))})
	}
}

//export callbackApproveTxCtx687eda_GasPriceUnit
func callbackApproveTxCtx687eda_GasPriceUnit(ptr unsafe.Pointer) C.int {
	if signal := qt.GetSignal(ptr, "gasPriceUnit"); signal != nil {
		return C.int(int32((*(*func() int)(signal))()))
	}

	return C.int(int32(NewApproveTxCtxFromPointer(ptr).GasPriceUnitDefault()))
}

func (ptr *ApproveTxCtx) ConnectGasPriceUnit(f func() int) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "gasPriceUnit"); signal != nil {
			f := func() int {
				(*(*func() int)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "gasPriceUnit", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "gasPriceUnit", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectGasPriceUnit() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "gasPriceUnit")
	}
}

func (ptr *ApproveTxCtx) GasPriceUnit() int {
	if ptr.Pointer() != nil {
		return int(int32(C.ApproveTxCtx687eda_GasPriceUnit(ptr.Pointer())))
	}
	return 0
}

func (ptr *ApproveTxCtx) GasPriceUnitDefault() int {
	if ptr.Pointer() != nil {
		return int(int32(C.ApproveTxCtx687eda_GasPriceUnitDefault(ptr.Pointer())))
	}
	return 0
}

//export callbackApproveTxCtx687eda_SetGasPriceUnit
func callbackApproveTxCtx687eda_SetGasPriceUnit(ptr unsafe.Pointer, gasPriceUnit C.int) {
	if signal := qt.GetSignal(ptr, "setGasPriceUnit"); signal != nil {
		(*(*func(int))(signal))(int(int32(gasPriceUnit)))
	} else {
		NewApproveTxCtxFromPointer(ptr).SetGasPriceUnitDefault(int(int32(gasPriceUnit)))
	}
}

func (ptr *ApproveTxCtx) ConnectSetGasPriceUnit(f func(gasPriceUnit int)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setGasPriceUnit"); signal != nil {
			f := func(gasPriceUnit int) {
				(*(*func(int))(signal))(gasPriceUnit)
				f(gasPriceUnit)
			}
			qt.ConnectSignal(ptr.Pointer(), "setGasPriceUnit", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setGasPriceUnit", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectSetGasPriceUnit() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setGasPriceUnit")
	}
}

func (ptr *ApproveTxCtx) SetGasPriceUnit(gasPriceUnit int) {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_SetGasPriceUnit(ptr.Pointer(), C.int(int32(gasPriceUnit)))
	}
}

func (ptr *ApproveTxCtx) SetGasPriceUnitDefault(gasPriceUnit int) {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_SetGasPriceUnitDefault(ptr.Pointer(), C.int(int32(gasPriceUnit)))
	}
}

//export callbackApproveTxCtx687eda_GasPriceUnitChanged
func callbackApproveTxCtx687eda_GasPriceUnitChanged(ptr unsafe.Pointer, gasPriceUnit C.int) {
	if signal := qt.GetSignal(ptr, "gasPriceUnitChanged"); signal != nil {
		(*(*func(int))(signal))(int(int32(gasPriceUnit)))
	}

}

func (ptr *ApproveTxCtx) ConnectGasPriceUnitChanged(f func(gasPriceUnit int)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "gasPriceUnitChanged") {
			C.ApproveTxCtx687eda_ConnectGasPriceUnitChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "gasPriceUnitChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "gasPriceUnitChanged"); signal != nil {
			f := func(gasPriceUnit int) {
				(*(*func(int))(signal))(gasPriceUnit)
				f(gasPriceUnit)
			}
			qt.ConnectSignal(ptr.Pointer(), "gasPriceUnitChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "gasPriceUnitChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectGasPriceUnitChanged() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectGasPriceUnitChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "gasPriceUnitChanged")
	}
}

func (ptr *ApproveTxCtx) GasPriceUnitChanged(gasPriceUnit int) {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_GasPriceUnitChanged(ptr.Pointer(), C.int(int32(gasPriceUnit)))
	}
}

//export callbackApproveTxCtx687eda_Nonce
func callbackApproveTxCtx687eda_Nonce(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "nonce"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveTxCtxFromPointer(ptr).NonceDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveTxCtx) ConnectNonce(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "nonce"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "nonce", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "nonce", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectNonce() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "nonce")
	}
}

func (ptr *ApproveTxCtx) Nonce() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_Nonce(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveTxCtx) NonceDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_NonceDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveTxCtx687eda_SetNonce
func callbackApproveTxCtx687eda_SetNonce(ptr unsafe.Pointer, nonce C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setNonce"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(nonce))
	} else {
		NewApproveTxCtxFromPointer(ptr).SetNonceDefault(cGoUnpackString(nonce))
	}
}

func (ptr *ApproveTxCtx) ConnectSetNonce(f func(nonce string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setNonce"); signal != nil {
			f := func(nonce string) {
				(*(*func(string))(signal))(nonce)
				f(nonce)
			}
			qt.ConnectSignal(ptr.Pointer(), "setNonce", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setNonce", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectSetNonce() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setNonce")
	}
}

func (ptr *ApproveTxCtx) SetNonce(nonce string) {
	if ptr.Pointer() != nil {
		var nonceC *C.char
		if nonce != "" {
			nonceC = C.CString(nonce)
			defer C.free(unsafe.Pointer(nonceC))
		}
		C.ApproveTxCtx687eda_SetNonce(ptr.Pointer(), C.struct_Moc_PackedString{data: nonceC, len: C.longlong(len(nonce))})
	}
}

func (ptr *ApproveTxCtx) SetNonceDefault(nonce string) {
	if ptr.Pointer() != nil {
		var nonceC *C.char
		if nonce != "" {
			nonceC = C.CString(nonce)
			defer C.free(unsafe.Pointer(nonceC))
		}
		C.ApproveTxCtx687eda_SetNonceDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: nonceC, len: C.longlong(len(nonce))})
	}
}

//export callbackApproveTxCtx687eda_NonceChanged
func callbackApproveTxCtx687eda_NonceChanged(ptr unsafe.Pointer, nonce C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "nonceChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(nonce))
	}

}

func (ptr *ApproveTxCtx) ConnectNonceChanged(f func(nonce string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "nonceChanged") {
			C.ApproveTxCtx687eda_ConnectNonceChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "nonceChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "nonceChanged"); signal != nil {
			f := func(nonce string) {
				(*(*func(string))(signal))(nonce)
				f(nonce)
			}
			qt.ConnectSignal(ptr.Pointer(), "nonceChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "nonceChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectNonceChanged() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectNonceChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "nonceChanged")
	}
}

func (ptr *ApproveTxCtx) NonceChanged(nonce string) {
	if ptr.Pointer() != nil {
		var nonceC *C.char
		if nonce != "" {
			nonceC = C.CString(nonce)
			defer C.free(unsafe.Pointer(nonceC))
		}
		C.ApproveTxCtx687eda_NonceChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: nonceC, len: C.longlong(len(nonce))})
	}
}

//export callbackApproveTxCtx687eda_Value
func callbackApproveTxCtx687eda_Value(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "value"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveTxCtxFromPointer(ptr).ValueDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveTxCtx) ConnectValue(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "value"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "value", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "value", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectValue() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "value")
	}
}

func (ptr *ApproveTxCtx) Value() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_Value(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveTxCtx) ValueDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_ValueDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveTxCtx687eda_SetValue
func callbackApproveTxCtx687eda_SetValue(ptr unsafe.Pointer, value C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setValue"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(value))
	} else {
		NewApproveTxCtxFromPointer(ptr).SetValueDefault(cGoUnpackString(value))
	}
}

func (ptr *ApproveTxCtx) ConnectSetValue(f func(value string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setValue"); signal != nil {
			f := func(value string) {
				(*(*func(string))(signal))(value)
				f(value)
			}
			qt.ConnectSignal(ptr.Pointer(), "setValue", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setValue", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectSetValue() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setValue")
	}
}

func (ptr *ApproveTxCtx) SetValue(value string) {
	if ptr.Pointer() != nil {
		var valueC *C.char
		if value != "" {
			valueC = C.CString(value)
			defer C.free(unsafe.Pointer(valueC))
		}
		C.ApproveTxCtx687eda_SetValue(ptr.Pointer(), C.struct_Moc_PackedString{data: valueC, len: C.longlong(len(value))})
	}
}

func (ptr *ApproveTxCtx) SetValueDefault(value string) {
	if ptr.Pointer() != nil {
		var valueC *C.char
		if value != "" {
			valueC = C.CString(value)
			defer C.free(unsafe.Pointer(valueC))
		}
		C.ApproveTxCtx687eda_SetValueDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: valueC, len: C.longlong(len(value))})
	}
}

//export callbackApproveTxCtx687eda_ValueChanged
func callbackApproveTxCtx687eda_ValueChanged(ptr unsafe.Pointer, value C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "valueChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(value))
	}

}

func (ptr *ApproveTxCtx) ConnectValueChanged(f func(value string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "valueChanged") {
			C.ApproveTxCtx687eda_ConnectValueChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "valueChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "valueChanged"); signal != nil {
			f := func(value string) {
				(*(*func(string))(signal))(value)
				f(value)
			}
			qt.ConnectSignal(ptr.Pointer(), "valueChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "valueChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectValueChanged() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectValueChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "valueChanged")
	}
}

func (ptr *ApproveTxCtx) ValueChanged(value string) {
	if ptr.Pointer() != nil {
		var valueC *C.char
		if value != "" {
			valueC = C.CString(value)
			defer C.free(unsafe.Pointer(valueC))
		}
		C.ApproveTxCtx687eda_ValueChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: valueC, len: C.longlong(len(value))})
	}
}

//export callbackApproveTxCtx687eda_Password
func callbackApproveTxCtx687eda_Password(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "password"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveTxCtxFromPointer(ptr).PasswordDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveTxCtx) ConnectPassword(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "password"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "password", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "password", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectPassword() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "password")
	}
}

func (ptr *ApproveTxCtx) Password() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_Password(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveTxCtx) PasswordDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_PasswordDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveTxCtx687eda_SetPassword
func callbackApproveTxCtx687eda_SetPassword(ptr unsafe.Pointer, password C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setPassword"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(password))
	} else {
		NewApproveTxCtxFromPointer(ptr).SetPasswordDefault(cGoUnpackString(password))
	}
}

func (ptr *ApproveTxCtx) ConnectSetPassword(f func(password string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setPassword"); signal != nil {
			f := func(password string) {
				(*(*func(string))(signal))(password)
				f(password)
			}
			qt.ConnectSignal(ptr.Pointer(), "setPassword", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setPassword", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectSetPassword() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setPassword")
	}
}

func (ptr *ApproveTxCtx) SetPassword(password string) {
	if ptr.Pointer() != nil {
		var passwordC *C.char
		if password != "" {
			passwordC = C.CString(password)
			defer C.free(unsafe.Pointer(passwordC))
		}
		C.ApproveTxCtx687eda_SetPassword(ptr.Pointer(), C.struct_Moc_PackedString{data: passwordC, len: C.longlong(len(password))})
	}
}

func (ptr *ApproveTxCtx) SetPasswordDefault(password string) {
	if ptr.Pointer() != nil {
		var passwordC *C.char
		if password != "" {
			passwordC = C.CString(password)
			defer C.free(unsafe.Pointer(passwordC))
		}
		C.ApproveTxCtx687eda_SetPasswordDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: passwordC, len: C.longlong(len(password))})
	}
}

//export callbackApproveTxCtx687eda_PasswordChanged
func callbackApproveTxCtx687eda_PasswordChanged(ptr unsafe.Pointer, password C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "passwordChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(password))
	}

}

func (ptr *ApproveTxCtx) ConnectPasswordChanged(f func(password string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "passwordChanged") {
			C.ApproveTxCtx687eda_ConnectPasswordChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "passwordChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "passwordChanged"); signal != nil {
			f := func(password string) {
				(*(*func(string))(signal))(password)
				f(password)
			}
			qt.ConnectSignal(ptr.Pointer(), "passwordChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "passwordChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectPasswordChanged() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectPasswordChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "passwordChanged")
	}
}

func (ptr *ApproveTxCtx) PasswordChanged(password string) {
	if ptr.Pointer() != nil {
		var passwordC *C.char
		if password != "" {
			passwordC = C.CString(password)
			defer C.free(unsafe.Pointer(passwordC))
		}
		C.ApproveTxCtx687eda_PasswordChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: passwordC, len: C.longlong(len(password))})
	}
}

//export callbackApproveTxCtx687eda_FromSrc
func callbackApproveTxCtx687eda_FromSrc(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "fromSrc"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveTxCtxFromPointer(ptr).FromSrcDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveTxCtx) ConnectFromSrc(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "fromSrc"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "fromSrc", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "fromSrc", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectFromSrc() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "fromSrc")
	}
}

func (ptr *ApproveTxCtx) FromSrc() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_FromSrc(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveTxCtx) FromSrcDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_FromSrcDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveTxCtx687eda_SetFromSrc
func callbackApproveTxCtx687eda_SetFromSrc(ptr unsafe.Pointer, fromSrc C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setFromSrc"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(fromSrc))
	} else {
		NewApproveTxCtxFromPointer(ptr).SetFromSrcDefault(cGoUnpackString(fromSrc))
	}
}

func (ptr *ApproveTxCtx) ConnectSetFromSrc(f func(fromSrc string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setFromSrc"); signal != nil {
			f := func(fromSrc string) {
				(*(*func(string))(signal))(fromSrc)
				f(fromSrc)
			}
			qt.ConnectSignal(ptr.Pointer(), "setFromSrc", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setFromSrc", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectSetFromSrc() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setFromSrc")
	}
}

func (ptr *ApproveTxCtx) SetFromSrc(fromSrc string) {
	if ptr.Pointer() != nil {
		var fromSrcC *C.char
		if fromSrc != "" {
			fromSrcC = C.CString(fromSrc)
			defer C.free(unsafe.Pointer(fromSrcC))
		}
		C.ApproveTxCtx687eda_SetFromSrc(ptr.Pointer(), C.struct_Moc_PackedString{data: fromSrcC, len: C.longlong(len(fromSrc))})
	}
}

func (ptr *ApproveTxCtx) SetFromSrcDefault(fromSrc string) {
	if ptr.Pointer() != nil {
		var fromSrcC *C.char
		if fromSrc != "" {
			fromSrcC = C.CString(fromSrc)
			defer C.free(unsafe.Pointer(fromSrcC))
		}
		C.ApproveTxCtx687eda_SetFromSrcDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: fromSrcC, len: C.longlong(len(fromSrc))})
	}
}

//export callbackApproveTxCtx687eda_FromSrcChanged
func callbackApproveTxCtx687eda_FromSrcChanged(ptr unsafe.Pointer, fromSrc C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "fromSrcChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(fromSrc))
	}

}

func (ptr *ApproveTxCtx) ConnectFromSrcChanged(f func(fromSrc string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "fromSrcChanged") {
			C.ApproveTxCtx687eda_ConnectFromSrcChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "fromSrcChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "fromSrcChanged"); signal != nil {
			f := func(fromSrc string) {
				(*(*func(string))(signal))(fromSrc)
				f(fromSrc)
			}
			qt.ConnectSignal(ptr.Pointer(), "fromSrcChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "fromSrcChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectFromSrcChanged() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectFromSrcChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "fromSrcChanged")
	}
}

func (ptr *ApproveTxCtx) FromSrcChanged(fromSrc string) {
	if ptr.Pointer() != nil {
		var fromSrcC *C.char
		if fromSrc != "" {
			fromSrcC = C.CString(fromSrc)
			defer C.free(unsafe.Pointer(fromSrcC))
		}
		C.ApproveTxCtx687eda_FromSrcChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: fromSrcC, len: C.longlong(len(fromSrc))})
	}
}

//export callbackApproveTxCtx687eda_ToSrc
func callbackApproveTxCtx687eda_ToSrc(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "toSrc"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveTxCtxFromPointer(ptr).ToSrcDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveTxCtx) ConnectToSrc(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "toSrc"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "toSrc", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "toSrc", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectToSrc() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "toSrc")
	}
}

func (ptr *ApproveTxCtx) ToSrc() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_ToSrc(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveTxCtx) ToSrcDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_ToSrcDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveTxCtx687eda_SetToSrc
func callbackApproveTxCtx687eda_SetToSrc(ptr unsafe.Pointer, toSrc C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setToSrc"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(toSrc))
	} else {
		NewApproveTxCtxFromPointer(ptr).SetToSrcDefault(cGoUnpackString(toSrc))
	}
}

func (ptr *ApproveTxCtx) ConnectSetToSrc(f func(toSrc string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setToSrc"); signal != nil {
			f := func(toSrc string) {
				(*(*func(string))(signal))(toSrc)
				f(toSrc)
			}
			qt.ConnectSignal(ptr.Pointer(), "setToSrc", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setToSrc", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectSetToSrc() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setToSrc")
	}
}

func (ptr *ApproveTxCtx) SetToSrc(toSrc string) {
	if ptr.Pointer() != nil {
		var toSrcC *C.char
		if toSrc != "" {
			toSrcC = C.CString(toSrc)
			defer C.free(unsafe.Pointer(toSrcC))
		}
		C.ApproveTxCtx687eda_SetToSrc(ptr.Pointer(), C.struct_Moc_PackedString{data: toSrcC, len: C.longlong(len(toSrc))})
	}
}

func (ptr *ApproveTxCtx) SetToSrcDefault(toSrc string) {
	if ptr.Pointer() != nil {
		var toSrcC *C.char
		if toSrc != "" {
			toSrcC = C.CString(toSrc)
			defer C.free(unsafe.Pointer(toSrcC))
		}
		C.ApproveTxCtx687eda_SetToSrcDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: toSrcC, len: C.longlong(len(toSrc))})
	}
}

//export callbackApproveTxCtx687eda_ToSrcChanged
func callbackApproveTxCtx687eda_ToSrcChanged(ptr unsafe.Pointer, toSrc C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "toSrcChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(toSrc))
	}

}

func (ptr *ApproveTxCtx) ConnectToSrcChanged(f func(toSrc string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "toSrcChanged") {
			C.ApproveTxCtx687eda_ConnectToSrcChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "toSrcChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "toSrcChanged"); signal != nil {
			f := func(toSrc string) {
				(*(*func(string))(signal))(toSrc)
				f(toSrc)
			}
			qt.ConnectSignal(ptr.Pointer(), "toSrcChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "toSrcChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectToSrcChanged() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectToSrcChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "toSrcChanged")
	}
}

func (ptr *ApproveTxCtx) ToSrcChanged(toSrc string) {
	if ptr.Pointer() != nil {
		var toSrcC *C.char
		if toSrc != "" {
			toSrcC = C.CString(toSrc)
			defer C.free(unsafe.Pointer(toSrcC))
		}
		C.ApproveTxCtx687eda_ToSrcChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: toSrcC, len: C.longlong(len(toSrc))})
	}
}

//export callbackApproveTxCtx687eda_Diff
func callbackApproveTxCtx687eda_Diff(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "diff"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveTxCtxFromPointer(ptr).DiffDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveTxCtx) ConnectDiff(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "diff"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "diff", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "diff", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectDiff() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "diff")
	}
}

func (ptr *ApproveTxCtx) Diff() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_Diff(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveTxCtx) DiffDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveTxCtx687eda_DiffDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveTxCtx687eda_SetDiff
func callbackApproveTxCtx687eda_SetDiff(ptr unsafe.Pointer, diff C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setDiff"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(diff))
	} else {
		NewApproveTxCtxFromPointer(ptr).SetDiffDefault(cGoUnpackString(diff))
	}
}

func (ptr *ApproveTxCtx) ConnectSetDiff(f func(diff string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setDiff"); signal != nil {
			f := func(diff string) {
				(*(*func(string))(signal))(diff)
				f(diff)
			}
			qt.ConnectSignal(ptr.Pointer(), "setDiff", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setDiff", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectSetDiff() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setDiff")
	}
}

func (ptr *ApproveTxCtx) SetDiff(diff string) {
	if ptr.Pointer() != nil {
		var diffC *C.char
		if diff != "" {
			diffC = C.CString(diff)
			defer C.free(unsafe.Pointer(diffC))
		}
		C.ApproveTxCtx687eda_SetDiff(ptr.Pointer(), C.struct_Moc_PackedString{data: diffC, len: C.longlong(len(diff))})
	}
}

func (ptr *ApproveTxCtx) SetDiffDefault(diff string) {
	if ptr.Pointer() != nil {
		var diffC *C.char
		if diff != "" {
			diffC = C.CString(diff)
			defer C.free(unsafe.Pointer(diffC))
		}
		C.ApproveTxCtx687eda_SetDiffDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: diffC, len: C.longlong(len(diff))})
	}
}

//export callbackApproveTxCtx687eda_DiffChanged
func callbackApproveTxCtx687eda_DiffChanged(ptr unsafe.Pointer, diff C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "diffChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(diff))
	}

}

func (ptr *ApproveTxCtx) ConnectDiffChanged(f func(diff string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "diffChanged") {
			C.ApproveTxCtx687eda_ConnectDiffChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "diffChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "diffChanged"); signal != nil {
			f := func(diff string) {
				(*(*func(string))(signal))(diff)
				f(diff)
			}
			qt.ConnectSignal(ptr.Pointer(), "diffChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "diffChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectDiffChanged() {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectDiffChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "diffChanged")
	}
}

func (ptr *ApproveTxCtx) DiffChanged(diff string) {
	if ptr.Pointer() != nil {
		var diffC *C.char
		if diff != "" {
			diffC = C.CString(diff)
			defer C.free(unsafe.Pointer(diffC))
		}
		C.ApproveTxCtx687eda_DiffChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: diffC, len: C.longlong(len(diff))})
	}
}

func ApproveTxCtx_QRegisterMetaType() int {
	return int(int32(C.ApproveTxCtx687eda_ApproveTxCtx687eda_QRegisterMetaType()))
}

func (ptr *ApproveTxCtx) QRegisterMetaType() int {
	return int(int32(C.ApproveTxCtx687eda_ApproveTxCtx687eda_QRegisterMetaType()))
}

func ApproveTxCtx_QRegisterMetaType2(typeName string) int {
	var typeNameC *C.char
	if typeName != "" {
		typeNameC = C.CString(typeName)
		defer C.free(unsafe.Pointer(typeNameC))
	}
	return int(int32(C.ApproveTxCtx687eda_ApproveTxCtx687eda_QRegisterMetaType2(typeNameC)))
}

func (ptr *ApproveTxCtx) QRegisterMetaType2(typeName string) int {
	var typeNameC *C.char
	if typeName != "" {
		typeNameC = C.CString(typeName)
		defer C.free(unsafe.Pointer(typeNameC))
	}
	return int(int32(C.ApproveTxCtx687eda_ApproveTxCtx687eda_QRegisterMetaType2(typeNameC)))
}

func ApproveTxCtx_QmlRegisterType() int {
	return int(int32(C.ApproveTxCtx687eda_ApproveTxCtx687eda_QmlRegisterType()))
}

func (ptr *ApproveTxCtx) QmlRegisterType() int {
	return int(int32(C.ApproveTxCtx687eda_ApproveTxCtx687eda_QmlRegisterType()))
}

func ApproveTxCtx_QmlRegisterType2(uri string, versionMajor int, versionMinor int, qmlName string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	return int(int32(C.ApproveTxCtx687eda_ApproveTxCtx687eda_QmlRegisterType2(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC)))
}

func (ptr *ApproveTxCtx) QmlRegisterType2(uri string, versionMajor int, versionMinor int, qmlName string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	return int(int32(C.ApproveTxCtx687eda_ApproveTxCtx687eda_QmlRegisterType2(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC)))
}

func ApproveTxCtx_QmlRegisterUncreatableType(uri string, versionMajor int, versionMinor int, qmlName string, message string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	var messageC *C.char
	if message != "" {
		messageC = C.CString(message)
		defer C.free(unsafe.Pointer(messageC))
	}
	return int(int32(C.ApproveTxCtx687eda_ApproveTxCtx687eda_QmlRegisterUncreatableType(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC, C.struct_Moc_PackedString{data: messageC, len: C.longlong(len(message))})))
}

func (ptr *ApproveTxCtx) QmlRegisterUncreatableType(uri string, versionMajor int, versionMinor int, qmlName string, message string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	var messageC *C.char
	if message != "" {
		messageC = C.CString(message)
		defer C.free(unsafe.Pointer(messageC))
	}
	return int(int32(C.ApproveTxCtx687eda_ApproveTxCtx687eda_QmlRegisterUncreatableType(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC, C.struct_Moc_PackedString{data: messageC, len: C.longlong(len(message))})))
}

func (ptr *ApproveTxCtx) __children_atList(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQObjectFromPointer(C.ApproveTxCtx687eda___children_atList(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *ApproveTxCtx) __children_setList(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda___children_setList(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *ApproveTxCtx) __children_newList() unsafe.Pointer {
	return C.ApproveTxCtx687eda___children_newList(ptr.Pointer())
}

func (ptr *ApproveTxCtx) __dynamicPropertyNames_atList(i int) *std_core.QByteArray {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQByteArrayFromPointer(C.ApproveTxCtx687eda___dynamicPropertyNames_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QByteArray).DestroyQByteArray)
		return tmpValue
	}
	return nil
}

func (ptr *ApproveTxCtx) __dynamicPropertyNames_setList(i std_core.QByteArray_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda___dynamicPropertyNames_setList(ptr.Pointer(), std_core.PointerFromQByteArray(i))
	}
}

func (ptr *ApproveTxCtx) __dynamicPropertyNames_newList() unsafe.Pointer {
	return C.ApproveTxCtx687eda___dynamicPropertyNames_newList(ptr.Pointer())
}

func (ptr *ApproveTxCtx) __findChildren_atList(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQObjectFromPointer(C.ApproveTxCtx687eda___findChildren_atList(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *ApproveTxCtx) __findChildren_setList(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda___findChildren_setList(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *ApproveTxCtx) __findChildren_newList() unsafe.Pointer {
	return C.ApproveTxCtx687eda___findChildren_newList(ptr.Pointer())
}

func (ptr *ApproveTxCtx) __findChildren_atList3(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQObjectFromPointer(C.ApproveTxCtx687eda___findChildren_atList3(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *ApproveTxCtx) __findChildren_setList3(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda___findChildren_setList3(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *ApproveTxCtx) __findChildren_newList3() unsafe.Pointer {
	return C.ApproveTxCtx687eda___findChildren_newList3(ptr.Pointer())
}

func NewApproveTxCtx(parent std_core.QObject_ITF) *ApproveTxCtx {
	ApproveTxCtx_QRegisterMetaType()
	tmpValue := NewApproveTxCtxFromPointer(C.ApproveTxCtx687eda_NewApproveTxCtx(std_core.PointerFromQObject(parent)))
	if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
		tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
	}
	return tmpValue
}

//export callbackApproveTxCtx687eda_DestroyApproveTxCtx
func callbackApproveTxCtx687eda_DestroyApproveTxCtx(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "~ApproveTxCtx"); signal != nil {
		(*(*func())(signal))()
	} else {
		NewApproveTxCtxFromPointer(ptr).DestroyApproveTxCtxDefault()
	}
}

func (ptr *ApproveTxCtx) ConnectDestroyApproveTxCtx(f func()) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "~ApproveTxCtx"); signal != nil {
			f := func() {
				(*(*func())(signal))()
				f()
			}
			qt.ConnectSignal(ptr.Pointer(), "~ApproveTxCtx", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "~ApproveTxCtx", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveTxCtx) DisconnectDestroyApproveTxCtx() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "~ApproveTxCtx")
	}
}

func (ptr *ApproveTxCtx) DestroyApproveTxCtx() {
	if ptr.Pointer() != nil {

		qt.SetFinalizer(ptr, nil)
		C.ApproveTxCtx687eda_DestroyApproveTxCtx(ptr.Pointer())
		ptr.SetPointer(nil)
	}
}

func (ptr *ApproveTxCtx) DestroyApproveTxCtxDefault() {
	if ptr.Pointer() != nil {

		qt.SetFinalizer(ptr, nil)
		C.ApproveTxCtx687eda_DestroyApproveTxCtxDefault(ptr.Pointer())
		ptr.SetPointer(nil)
	}
}

//export callbackApproveTxCtx687eda_ChildEvent
func callbackApproveTxCtx687eda_ChildEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "childEvent"); signal != nil {
		(*(*func(*std_core.QChildEvent))(signal))(std_core.NewQChildEventFromPointer(event))
	} else {
		NewApproveTxCtxFromPointer(ptr).ChildEventDefault(std_core.NewQChildEventFromPointer(event))
	}
}

func (ptr *ApproveTxCtx) ChildEventDefault(event std_core.QChildEvent_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_ChildEventDefault(ptr.Pointer(), std_core.PointerFromQChildEvent(event))
	}
}

//export callbackApproveTxCtx687eda_ConnectNotify
func callbackApproveTxCtx687eda_ConnectNotify(ptr unsafe.Pointer, sign unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "connectNotify"); signal != nil {
		(*(*func(*std_core.QMetaMethod))(signal))(std_core.NewQMetaMethodFromPointer(sign))
	} else {
		NewApproveTxCtxFromPointer(ptr).ConnectNotifyDefault(std_core.NewQMetaMethodFromPointer(sign))
	}
}

func (ptr *ApproveTxCtx) ConnectNotifyDefault(sign std_core.QMetaMethod_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_ConnectNotifyDefault(ptr.Pointer(), std_core.PointerFromQMetaMethod(sign))
	}
}

//export callbackApproveTxCtx687eda_CustomEvent
func callbackApproveTxCtx687eda_CustomEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "customEvent"); signal != nil {
		(*(*func(*std_core.QEvent))(signal))(std_core.NewQEventFromPointer(event))
	} else {
		NewApproveTxCtxFromPointer(ptr).CustomEventDefault(std_core.NewQEventFromPointer(event))
	}
}

func (ptr *ApproveTxCtx) CustomEventDefault(event std_core.QEvent_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_CustomEventDefault(ptr.Pointer(), std_core.PointerFromQEvent(event))
	}
}

//export callbackApproveTxCtx687eda_DeleteLater
func callbackApproveTxCtx687eda_DeleteLater(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "deleteLater"); signal != nil {
		(*(*func())(signal))()
	} else {
		NewApproveTxCtxFromPointer(ptr).DeleteLaterDefault()
	}
}

func (ptr *ApproveTxCtx) DeleteLaterDefault() {
	if ptr.Pointer() != nil {

		qt.SetFinalizer(ptr, nil)
		C.ApproveTxCtx687eda_DeleteLaterDefault(ptr.Pointer())
	}
}

//export callbackApproveTxCtx687eda_Destroyed
func callbackApproveTxCtx687eda_Destroyed(ptr unsafe.Pointer, obj unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "destroyed"); signal != nil {
		(*(*func(*std_core.QObject))(signal))(std_core.NewQObjectFromPointer(obj))
	}
	qt.Unregister(ptr)

}

//export callbackApproveTxCtx687eda_DisconnectNotify
func callbackApproveTxCtx687eda_DisconnectNotify(ptr unsafe.Pointer, sign unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "disconnectNotify"); signal != nil {
		(*(*func(*std_core.QMetaMethod))(signal))(std_core.NewQMetaMethodFromPointer(sign))
	} else {
		NewApproveTxCtxFromPointer(ptr).DisconnectNotifyDefault(std_core.NewQMetaMethodFromPointer(sign))
	}
}

func (ptr *ApproveTxCtx) DisconnectNotifyDefault(sign std_core.QMetaMethod_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_DisconnectNotifyDefault(ptr.Pointer(), std_core.PointerFromQMetaMethod(sign))
	}
}

//export callbackApproveTxCtx687eda_Event
func callbackApproveTxCtx687eda_Event(ptr unsafe.Pointer, e unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "event"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QEvent) bool)(signal))(std_core.NewQEventFromPointer(e)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewApproveTxCtxFromPointer(ptr).EventDefault(std_core.NewQEventFromPointer(e)))))
}

func (ptr *ApproveTxCtx) EventDefault(e std_core.QEvent_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.ApproveTxCtx687eda_EventDefault(ptr.Pointer(), std_core.PointerFromQEvent(e))) != 0
	}
	return false
}

//export callbackApproveTxCtx687eda_EventFilter
func callbackApproveTxCtx687eda_EventFilter(ptr unsafe.Pointer, watched unsafe.Pointer, event unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "eventFilter"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QObject, *std_core.QEvent) bool)(signal))(std_core.NewQObjectFromPointer(watched), std_core.NewQEventFromPointer(event)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewApproveTxCtxFromPointer(ptr).EventFilterDefault(std_core.NewQObjectFromPointer(watched), std_core.NewQEventFromPointer(event)))))
}

func (ptr *ApproveTxCtx) EventFilterDefault(watched std_core.QObject_ITF, event std_core.QEvent_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.ApproveTxCtx687eda_EventFilterDefault(ptr.Pointer(), std_core.PointerFromQObject(watched), std_core.PointerFromQEvent(event))) != 0
	}
	return false
}

//export callbackApproveTxCtx687eda_ObjectNameChanged
func callbackApproveTxCtx687eda_ObjectNameChanged(ptr unsafe.Pointer, objectName C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "objectNameChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(objectName))
	}

}

//export callbackApproveTxCtx687eda_TimerEvent
func callbackApproveTxCtx687eda_TimerEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "timerEvent"); signal != nil {
		(*(*func(*std_core.QTimerEvent))(signal))(std_core.NewQTimerEventFromPointer(event))
	} else {
		NewApproveTxCtxFromPointer(ptr).TimerEventDefault(std_core.NewQTimerEventFromPointer(event))
	}
}

func (ptr *ApproveTxCtx) TimerEventDefault(event std_core.QTimerEvent_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveTxCtx687eda_TimerEventDefault(ptr.Pointer(), std_core.PointerFromQTimerEvent(event))
	}
}

type LoginContext_ITF interface {
	std_core.QObject_ITF
	LoginContext_PTR() *LoginContext
}

func (ptr *LoginContext) LoginContext_PTR() *LoginContext {
	return ptr
}

func (ptr *LoginContext) Pointer() unsafe.Pointer {
	if ptr != nil {
		return ptr.QObject_PTR().Pointer()
	}
	return nil
}

func (ptr *LoginContext) SetPointer(p unsafe.Pointer) {
	if ptr != nil {
		ptr.QObject_PTR().SetPointer(p)
	}
}

func PointerFromLoginContext(ptr LoginContext_ITF) unsafe.Pointer {
	if ptr != nil {
		return ptr.LoginContext_PTR().Pointer()
	}
	return nil
}

func NewLoginContextFromPointer(ptr unsafe.Pointer) (n *LoginContext) {
	if gPtr, ok := qt.Receive(ptr); !ok {
		n = new(LoginContext)
		n.SetPointer(ptr)
	} else {
		switch deduced := gPtr.(type) {
		case *LoginContext:
			n = deduced

		case *std_core.QObject:
			n = &LoginContext{QObject: *deduced}

		default:
			n = new(LoginContext)
			n.SetPointer(ptr)
		}
	}
	return
}

//export callbackLoginContext687eda_Constructor
func callbackLoginContext687eda_Constructor(ptr unsafe.Pointer) {
	this := NewLoginContextFromPointer(ptr)
	qt.Register(ptr, this)
	this.ConnectStart(this.start)
	this.ConnectCheckPath(this.checkPath)
}

//export callbackLoginContext687eda_Start
func callbackLoginContext687eda_Start(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "start"); signal != nil {
		(*(*func())(signal))()
	}

}

func (ptr *LoginContext) ConnectStart(f func()) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "start") {
			C.LoginContext687eda_ConnectStart(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "start")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "start"); signal != nil {
			f := func() {
				(*(*func())(signal))()
				f()
			}
			qt.ConnectSignal(ptr.Pointer(), "start", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "start", unsafe.Pointer(&f))
		}
	}
}

func (ptr *LoginContext) DisconnectStart() {
	if ptr.Pointer() != nil {
		C.LoginContext687eda_DisconnectStart(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "start")
	}
}

func (ptr *LoginContext) Start() {
	if ptr.Pointer() != nil {
		C.LoginContext687eda_Start(ptr.Pointer())
	}
}

//export callbackLoginContext687eda_CheckPath
func callbackLoginContext687eda_CheckPath(ptr unsafe.Pointer, b C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "checkPath"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(b))
	}

}

func (ptr *LoginContext) ConnectCheckPath(f func(b string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "checkPath") {
			C.LoginContext687eda_ConnectCheckPath(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "checkPath")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "checkPath"); signal != nil {
			f := func(b string) {
				(*(*func(string))(signal))(b)
				f(b)
			}
			qt.ConnectSignal(ptr.Pointer(), "checkPath", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "checkPath", unsafe.Pointer(&f))
		}
	}
}

func (ptr *LoginContext) DisconnectCheckPath() {
	if ptr.Pointer() != nil {
		C.LoginContext687eda_DisconnectCheckPath(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "checkPath")
	}
}

func (ptr *LoginContext) CheckPath(b string) {
	if ptr.Pointer() != nil {
		var bC *C.char
		if b != "" {
			bC = C.CString(b)
			defer C.free(unsafe.Pointer(bC))
		}
		C.LoginContext687eda_CheckPath(ptr.Pointer(), C.struct_Moc_PackedString{data: bC, len: C.longlong(len(b))})
	}
}

//export callbackLoginContext687eda_ClefPath
func callbackLoginContext687eda_ClefPath(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "clefPath"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewLoginContextFromPointer(ptr).ClefPathDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *LoginContext) ConnectClefPath(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "clefPath"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "clefPath", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "clefPath", unsafe.Pointer(&f))
		}
	}
}

func (ptr *LoginContext) DisconnectClefPath() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "clefPath")
	}
}

func (ptr *LoginContext) ClefPath() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.LoginContext687eda_ClefPath(ptr.Pointer()))
	}
	return ""
}

func (ptr *LoginContext) ClefPathDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.LoginContext687eda_ClefPathDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackLoginContext687eda_SetClefPath
func callbackLoginContext687eda_SetClefPath(ptr unsafe.Pointer, clefPath C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setClefPath"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(clefPath))
	} else {
		NewLoginContextFromPointer(ptr).SetClefPathDefault(cGoUnpackString(clefPath))
	}
}

func (ptr *LoginContext) ConnectSetClefPath(f func(clefPath string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setClefPath"); signal != nil {
			f := func(clefPath string) {
				(*(*func(string))(signal))(clefPath)
				f(clefPath)
			}
			qt.ConnectSignal(ptr.Pointer(), "setClefPath", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setClefPath", unsafe.Pointer(&f))
		}
	}
}

func (ptr *LoginContext) DisconnectSetClefPath() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setClefPath")
	}
}

func (ptr *LoginContext) SetClefPath(clefPath string) {
	if ptr.Pointer() != nil {
		var clefPathC *C.char
		if clefPath != "" {
			clefPathC = C.CString(clefPath)
			defer C.free(unsafe.Pointer(clefPathC))
		}
		C.LoginContext687eda_SetClefPath(ptr.Pointer(), C.struct_Moc_PackedString{data: clefPathC, len: C.longlong(len(clefPath))})
	}
}

func (ptr *LoginContext) SetClefPathDefault(clefPath string) {
	if ptr.Pointer() != nil {
		var clefPathC *C.char
		if clefPath != "" {
			clefPathC = C.CString(clefPath)
			defer C.free(unsafe.Pointer(clefPathC))
		}
		C.LoginContext687eda_SetClefPathDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: clefPathC, len: C.longlong(len(clefPath))})
	}
}

//export callbackLoginContext687eda_ClefPathChanged
func callbackLoginContext687eda_ClefPathChanged(ptr unsafe.Pointer, clefPath C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "clefPathChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(clefPath))
	}

}

func (ptr *LoginContext) ConnectClefPathChanged(f func(clefPath string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "clefPathChanged") {
			C.LoginContext687eda_ConnectClefPathChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "clefPathChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "clefPathChanged"); signal != nil {
			f := func(clefPath string) {
				(*(*func(string))(signal))(clefPath)
				f(clefPath)
			}
			qt.ConnectSignal(ptr.Pointer(), "clefPathChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "clefPathChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *LoginContext) DisconnectClefPathChanged() {
	if ptr.Pointer() != nil {
		C.LoginContext687eda_DisconnectClefPathChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "clefPathChanged")
	}
}

func (ptr *LoginContext) ClefPathChanged(clefPath string) {
	if ptr.Pointer() != nil {
		var clefPathC *C.char
		if clefPath != "" {
			clefPathC = C.CString(clefPath)
			defer C.free(unsafe.Pointer(clefPathC))
		}
		C.LoginContext687eda_ClefPathChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: clefPathC, len: C.longlong(len(clefPath))})
	}
}

//export callbackLoginContext687eda_BinaryHash
func callbackLoginContext687eda_BinaryHash(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "binaryHash"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewLoginContextFromPointer(ptr).BinaryHashDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *LoginContext) ConnectBinaryHash(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "binaryHash"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "binaryHash", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "binaryHash", unsafe.Pointer(&f))
		}
	}
}

func (ptr *LoginContext) DisconnectBinaryHash() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "binaryHash")
	}
}

func (ptr *LoginContext) BinaryHash() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.LoginContext687eda_BinaryHash(ptr.Pointer()))
	}
	return ""
}

func (ptr *LoginContext) BinaryHashDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.LoginContext687eda_BinaryHashDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackLoginContext687eda_SetBinaryHash
func callbackLoginContext687eda_SetBinaryHash(ptr unsafe.Pointer, binaryHash C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setBinaryHash"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(binaryHash))
	} else {
		NewLoginContextFromPointer(ptr).SetBinaryHashDefault(cGoUnpackString(binaryHash))
	}
}

func (ptr *LoginContext) ConnectSetBinaryHash(f func(binaryHash string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setBinaryHash"); signal != nil {
			f := func(binaryHash string) {
				(*(*func(string))(signal))(binaryHash)
				f(binaryHash)
			}
			qt.ConnectSignal(ptr.Pointer(), "setBinaryHash", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setBinaryHash", unsafe.Pointer(&f))
		}
	}
}

func (ptr *LoginContext) DisconnectSetBinaryHash() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setBinaryHash")
	}
}

func (ptr *LoginContext) SetBinaryHash(binaryHash string) {
	if ptr.Pointer() != nil {
		var binaryHashC *C.char
		if binaryHash != "" {
			binaryHashC = C.CString(binaryHash)
			defer C.free(unsafe.Pointer(binaryHashC))
		}
		C.LoginContext687eda_SetBinaryHash(ptr.Pointer(), C.struct_Moc_PackedString{data: binaryHashC, len: C.longlong(len(binaryHash))})
	}
}

func (ptr *LoginContext) SetBinaryHashDefault(binaryHash string) {
	if ptr.Pointer() != nil {
		var binaryHashC *C.char
		if binaryHash != "" {
			binaryHashC = C.CString(binaryHash)
			defer C.free(unsafe.Pointer(binaryHashC))
		}
		C.LoginContext687eda_SetBinaryHashDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: binaryHashC, len: C.longlong(len(binaryHash))})
	}
}

//export callbackLoginContext687eda_BinaryHashChanged
func callbackLoginContext687eda_BinaryHashChanged(ptr unsafe.Pointer, binaryHash C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "binaryHashChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(binaryHash))
	}

}

func (ptr *LoginContext) ConnectBinaryHashChanged(f func(binaryHash string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "binaryHashChanged") {
			C.LoginContext687eda_ConnectBinaryHashChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "binaryHashChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "binaryHashChanged"); signal != nil {
			f := func(binaryHash string) {
				(*(*func(string))(signal))(binaryHash)
				f(binaryHash)
			}
			qt.ConnectSignal(ptr.Pointer(), "binaryHashChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "binaryHashChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *LoginContext) DisconnectBinaryHashChanged() {
	if ptr.Pointer() != nil {
		C.LoginContext687eda_DisconnectBinaryHashChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "binaryHashChanged")
	}
}

func (ptr *LoginContext) BinaryHashChanged(binaryHash string) {
	if ptr.Pointer() != nil {
		var binaryHashC *C.char
		if binaryHash != "" {
			binaryHashC = C.CString(binaryHash)
			defer C.free(unsafe.Pointer(binaryHashC))
		}
		C.LoginContext687eda_BinaryHashChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: binaryHashC, len: C.longlong(len(binaryHash))})
	}
}

//export callbackLoginContext687eda_Error
func callbackLoginContext687eda_Error(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "error"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewLoginContextFromPointer(ptr).ErrorDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *LoginContext) ConnectError(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "error"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "error", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "error", unsafe.Pointer(&f))
		}
	}
}

func (ptr *LoginContext) DisconnectError() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "error")
	}
}

func (ptr *LoginContext) Error() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.LoginContext687eda_Error(ptr.Pointer()))
	}
	return ""
}

func (ptr *LoginContext) ErrorDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.LoginContext687eda_ErrorDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackLoginContext687eda_SetError
func callbackLoginContext687eda_SetError(ptr unsafe.Pointer, error C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setError"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(error))
	} else {
		NewLoginContextFromPointer(ptr).SetErrorDefault(cGoUnpackString(error))
	}
}

func (ptr *LoginContext) ConnectSetError(f func(error string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setError"); signal != nil {
			f := func(error string) {
				(*(*func(string))(signal))(error)
				f(error)
			}
			qt.ConnectSignal(ptr.Pointer(), "setError", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setError", unsafe.Pointer(&f))
		}
	}
}

func (ptr *LoginContext) DisconnectSetError() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setError")
	}
}

func (ptr *LoginContext) SetError(error string) {
	if ptr.Pointer() != nil {
		var errorC *C.char
		if error != "" {
			errorC = C.CString(error)
			defer C.free(unsafe.Pointer(errorC))
		}
		C.LoginContext687eda_SetError(ptr.Pointer(), C.struct_Moc_PackedString{data: errorC, len: C.longlong(len(error))})
	}
}

func (ptr *LoginContext) SetErrorDefault(error string) {
	if ptr.Pointer() != nil {
		var errorC *C.char
		if error != "" {
			errorC = C.CString(error)
			defer C.free(unsafe.Pointer(errorC))
		}
		C.LoginContext687eda_SetErrorDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: errorC, len: C.longlong(len(error))})
	}
}

//export callbackLoginContext687eda_ErrorChanged
func callbackLoginContext687eda_ErrorChanged(ptr unsafe.Pointer, error C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "errorChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(error))
	}

}

func (ptr *LoginContext) ConnectErrorChanged(f func(error string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "errorChanged") {
			C.LoginContext687eda_ConnectErrorChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "errorChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "errorChanged"); signal != nil {
			f := func(error string) {
				(*(*func(string))(signal))(error)
				f(error)
			}
			qt.ConnectSignal(ptr.Pointer(), "errorChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "errorChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *LoginContext) DisconnectErrorChanged() {
	if ptr.Pointer() != nil {
		C.LoginContext687eda_DisconnectErrorChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "errorChanged")
	}
}

func (ptr *LoginContext) ErrorChanged(error string) {
	if ptr.Pointer() != nil {
		var errorC *C.char
		if error != "" {
			errorC = C.CString(error)
			defer C.free(unsafe.Pointer(errorC))
		}
		C.LoginContext687eda_ErrorChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: errorC, len: C.longlong(len(error))})
	}
}

func LoginContext_QRegisterMetaType() int {
	return int(int32(C.LoginContext687eda_LoginContext687eda_QRegisterMetaType()))
}

func (ptr *LoginContext) QRegisterMetaType() int {
	return int(int32(C.LoginContext687eda_LoginContext687eda_QRegisterMetaType()))
}

func LoginContext_QRegisterMetaType2(typeName string) int {
	var typeNameC *C.char
	if typeName != "" {
		typeNameC = C.CString(typeName)
		defer C.free(unsafe.Pointer(typeNameC))
	}
	return int(int32(C.LoginContext687eda_LoginContext687eda_QRegisterMetaType2(typeNameC)))
}

func (ptr *LoginContext) QRegisterMetaType2(typeName string) int {
	var typeNameC *C.char
	if typeName != "" {
		typeNameC = C.CString(typeName)
		defer C.free(unsafe.Pointer(typeNameC))
	}
	return int(int32(C.LoginContext687eda_LoginContext687eda_QRegisterMetaType2(typeNameC)))
}

func LoginContext_QmlRegisterType() int {
	return int(int32(C.LoginContext687eda_LoginContext687eda_QmlRegisterType()))
}

func (ptr *LoginContext) QmlRegisterType() int {
	return int(int32(C.LoginContext687eda_LoginContext687eda_QmlRegisterType()))
}

func LoginContext_QmlRegisterType2(uri string, versionMajor int, versionMinor int, qmlName string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	return int(int32(C.LoginContext687eda_LoginContext687eda_QmlRegisterType2(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC)))
}

func (ptr *LoginContext) QmlRegisterType2(uri string, versionMajor int, versionMinor int, qmlName string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	return int(int32(C.LoginContext687eda_LoginContext687eda_QmlRegisterType2(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC)))
}

func LoginContext_QmlRegisterUncreatableType(uri string, versionMajor int, versionMinor int, qmlName string, message string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	var messageC *C.char
	if message != "" {
		messageC = C.CString(message)
		defer C.free(unsafe.Pointer(messageC))
	}
	return int(int32(C.LoginContext687eda_LoginContext687eda_QmlRegisterUncreatableType(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC, C.struct_Moc_PackedString{data: messageC, len: C.longlong(len(message))})))
}

func (ptr *LoginContext) QmlRegisterUncreatableType(uri string, versionMajor int, versionMinor int, qmlName string, message string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	var messageC *C.char
	if message != "" {
		messageC = C.CString(message)
		defer C.free(unsafe.Pointer(messageC))
	}
	return int(int32(C.LoginContext687eda_LoginContext687eda_QmlRegisterUncreatableType(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC, C.struct_Moc_PackedString{data: messageC, len: C.longlong(len(message))})))
}

func (ptr *LoginContext) __children_atList(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQObjectFromPointer(C.LoginContext687eda___children_atList(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *LoginContext) __children_setList(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.LoginContext687eda___children_setList(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *LoginContext) __children_newList() unsafe.Pointer {
	return C.LoginContext687eda___children_newList(ptr.Pointer())
}

func (ptr *LoginContext) __dynamicPropertyNames_atList(i int) *std_core.QByteArray {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQByteArrayFromPointer(C.LoginContext687eda___dynamicPropertyNames_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QByteArray).DestroyQByteArray)
		return tmpValue
	}
	return nil
}

func (ptr *LoginContext) __dynamicPropertyNames_setList(i std_core.QByteArray_ITF) {
	if ptr.Pointer() != nil {
		C.LoginContext687eda___dynamicPropertyNames_setList(ptr.Pointer(), std_core.PointerFromQByteArray(i))
	}
}

func (ptr *LoginContext) __dynamicPropertyNames_newList() unsafe.Pointer {
	return C.LoginContext687eda___dynamicPropertyNames_newList(ptr.Pointer())
}

func (ptr *LoginContext) __findChildren_atList(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQObjectFromPointer(C.LoginContext687eda___findChildren_atList(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *LoginContext) __findChildren_setList(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.LoginContext687eda___findChildren_setList(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *LoginContext) __findChildren_newList() unsafe.Pointer {
	return C.LoginContext687eda___findChildren_newList(ptr.Pointer())
}

func (ptr *LoginContext) __findChildren_atList3(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQObjectFromPointer(C.LoginContext687eda___findChildren_atList3(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *LoginContext) __findChildren_setList3(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.LoginContext687eda___findChildren_setList3(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *LoginContext) __findChildren_newList3() unsafe.Pointer {
	return C.LoginContext687eda___findChildren_newList3(ptr.Pointer())
}

func NewLoginContext(parent std_core.QObject_ITF) *LoginContext {
	LoginContext_QRegisterMetaType()
	tmpValue := NewLoginContextFromPointer(C.LoginContext687eda_NewLoginContext(std_core.PointerFromQObject(parent)))
	if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
		tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
	}
	return tmpValue
}

//export callbackLoginContext687eda_DestroyLoginContext
func callbackLoginContext687eda_DestroyLoginContext(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "~LoginContext"); signal != nil {
		(*(*func())(signal))()
	} else {
		NewLoginContextFromPointer(ptr).DestroyLoginContextDefault()
	}
}

func (ptr *LoginContext) ConnectDestroyLoginContext(f func()) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "~LoginContext"); signal != nil {
			f := func() {
				(*(*func())(signal))()
				f()
			}
			qt.ConnectSignal(ptr.Pointer(), "~LoginContext", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "~LoginContext", unsafe.Pointer(&f))
		}
	}
}

func (ptr *LoginContext) DisconnectDestroyLoginContext() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "~LoginContext")
	}
}

func (ptr *LoginContext) DestroyLoginContext() {
	if ptr.Pointer() != nil {

		qt.SetFinalizer(ptr, nil)
		C.LoginContext687eda_DestroyLoginContext(ptr.Pointer())
		ptr.SetPointer(nil)
	}
}

func (ptr *LoginContext) DestroyLoginContextDefault() {
	if ptr.Pointer() != nil {

		qt.SetFinalizer(ptr, nil)
		C.LoginContext687eda_DestroyLoginContextDefault(ptr.Pointer())
		ptr.SetPointer(nil)
	}
}

//export callbackLoginContext687eda_ChildEvent
func callbackLoginContext687eda_ChildEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "childEvent"); signal != nil {
		(*(*func(*std_core.QChildEvent))(signal))(std_core.NewQChildEventFromPointer(event))
	} else {
		NewLoginContextFromPointer(ptr).ChildEventDefault(std_core.NewQChildEventFromPointer(event))
	}
}

func (ptr *LoginContext) ChildEventDefault(event std_core.QChildEvent_ITF) {
	if ptr.Pointer() != nil {
		C.LoginContext687eda_ChildEventDefault(ptr.Pointer(), std_core.PointerFromQChildEvent(event))
	}
}

//export callbackLoginContext687eda_ConnectNotify
func callbackLoginContext687eda_ConnectNotify(ptr unsafe.Pointer, sign unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "connectNotify"); signal != nil {
		(*(*func(*std_core.QMetaMethod))(signal))(std_core.NewQMetaMethodFromPointer(sign))
	} else {
		NewLoginContextFromPointer(ptr).ConnectNotifyDefault(std_core.NewQMetaMethodFromPointer(sign))
	}
}

func (ptr *LoginContext) ConnectNotifyDefault(sign std_core.QMetaMethod_ITF) {
	if ptr.Pointer() != nil {
		C.LoginContext687eda_ConnectNotifyDefault(ptr.Pointer(), std_core.PointerFromQMetaMethod(sign))
	}
}

//export callbackLoginContext687eda_CustomEvent
func callbackLoginContext687eda_CustomEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "customEvent"); signal != nil {
		(*(*func(*std_core.QEvent))(signal))(std_core.NewQEventFromPointer(event))
	} else {
		NewLoginContextFromPointer(ptr).CustomEventDefault(std_core.NewQEventFromPointer(event))
	}
}

func (ptr *LoginContext) CustomEventDefault(event std_core.QEvent_ITF) {
	if ptr.Pointer() != nil {
		C.LoginContext687eda_CustomEventDefault(ptr.Pointer(), std_core.PointerFromQEvent(event))
	}
}

//export callbackLoginContext687eda_DeleteLater
func callbackLoginContext687eda_DeleteLater(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "deleteLater"); signal != nil {
		(*(*func())(signal))()
	} else {
		NewLoginContextFromPointer(ptr).DeleteLaterDefault()
	}
}

func (ptr *LoginContext) DeleteLaterDefault() {
	if ptr.Pointer() != nil {

		qt.SetFinalizer(ptr, nil)
		C.LoginContext687eda_DeleteLaterDefault(ptr.Pointer())
	}
}

//export callbackLoginContext687eda_Destroyed
func callbackLoginContext687eda_Destroyed(ptr unsafe.Pointer, obj unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "destroyed"); signal != nil {
		(*(*func(*std_core.QObject))(signal))(std_core.NewQObjectFromPointer(obj))
	}
	qt.Unregister(ptr)

}

//export callbackLoginContext687eda_DisconnectNotify
func callbackLoginContext687eda_DisconnectNotify(ptr unsafe.Pointer, sign unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "disconnectNotify"); signal != nil {
		(*(*func(*std_core.QMetaMethod))(signal))(std_core.NewQMetaMethodFromPointer(sign))
	} else {
		NewLoginContextFromPointer(ptr).DisconnectNotifyDefault(std_core.NewQMetaMethodFromPointer(sign))
	}
}

func (ptr *LoginContext) DisconnectNotifyDefault(sign std_core.QMetaMethod_ITF) {
	if ptr.Pointer() != nil {
		C.LoginContext687eda_DisconnectNotifyDefault(ptr.Pointer(), std_core.PointerFromQMetaMethod(sign))
	}
}

//export callbackLoginContext687eda_Event
func callbackLoginContext687eda_Event(ptr unsafe.Pointer, e unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "event"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QEvent) bool)(signal))(std_core.NewQEventFromPointer(e)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewLoginContextFromPointer(ptr).EventDefault(std_core.NewQEventFromPointer(e)))))
}

func (ptr *LoginContext) EventDefault(e std_core.QEvent_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.LoginContext687eda_EventDefault(ptr.Pointer(), std_core.PointerFromQEvent(e))) != 0
	}
	return false
}

//export callbackLoginContext687eda_EventFilter
func callbackLoginContext687eda_EventFilter(ptr unsafe.Pointer, watched unsafe.Pointer, event unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "eventFilter"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QObject, *std_core.QEvent) bool)(signal))(std_core.NewQObjectFromPointer(watched), std_core.NewQEventFromPointer(event)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewLoginContextFromPointer(ptr).EventFilterDefault(std_core.NewQObjectFromPointer(watched), std_core.NewQEventFromPointer(event)))))
}

func (ptr *LoginContext) EventFilterDefault(watched std_core.QObject_ITF, event std_core.QEvent_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.LoginContext687eda_EventFilterDefault(ptr.Pointer(), std_core.PointerFromQObject(watched), std_core.PointerFromQEvent(event))) != 0
	}
	return false
}

//export callbackLoginContext687eda_ObjectNameChanged
func callbackLoginContext687eda_ObjectNameChanged(ptr unsafe.Pointer, objectName C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "objectNameChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(objectName))
	}

}

//export callbackLoginContext687eda_TimerEvent
func callbackLoginContext687eda_TimerEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "timerEvent"); signal != nil {
		(*(*func(*std_core.QTimerEvent))(signal))(std_core.NewQTimerEventFromPointer(event))
	} else {
		NewLoginContextFromPointer(ptr).TimerEventDefault(std_core.NewQTimerEventFromPointer(event))
	}
}

func (ptr *LoginContext) TimerEventDefault(event std_core.QTimerEvent_ITF) {
	if ptr.Pointer() != nil {
		C.LoginContext687eda_TimerEventDefault(ptr.Pointer(), std_core.PointerFromQTimerEvent(event))
	}
}

type TxListAccountsModel_ITF interface {
	std_core.QAbstractListModel_ITF
	TxListAccountsModel_PTR() *TxListAccountsModel
}

func (ptr *TxListAccountsModel) TxListAccountsModel_PTR() *TxListAccountsModel {
	return ptr
}

func (ptr *TxListAccountsModel) Pointer() unsafe.Pointer {
	if ptr != nil {
		return ptr.QAbstractListModel_PTR().Pointer()
	}
	return nil
}

func (ptr *TxListAccountsModel) SetPointer(p unsafe.Pointer) {
	if ptr != nil {
		ptr.QAbstractListModel_PTR().SetPointer(p)
	}
}

func PointerFromTxListAccountsModel(ptr TxListAccountsModel_ITF) unsafe.Pointer {
	if ptr != nil {
		return ptr.TxListAccountsModel_PTR().Pointer()
	}
	return nil
}

func NewTxListAccountsModelFromPointer(ptr unsafe.Pointer) (n *TxListAccountsModel) {
	if gPtr, ok := qt.Receive(ptr); !ok {
		n = new(TxListAccountsModel)
		n.SetPointer(ptr)
	} else {
		switch deduced := gPtr.(type) {
		case *TxListAccountsModel:
			n = deduced

		case *std_core.QAbstractListModel:
			n = &TxListAccountsModel{QAbstractListModel: *deduced}

		default:
			n = new(TxListAccountsModel)
			n.SetPointer(ptr)
		}
	}
	return
}
func (this *TxListAccountsModel) Init() { this.init() }

//export callbackTxListAccountsModel687eda_Constructor
func callbackTxListAccountsModel687eda_Constructor(ptr unsafe.Pointer) {
	this := NewTxListAccountsModelFromPointer(ptr)
	qt.Register(ptr, this)
	this.ConnectAdd(this.add)
	this.init()
}

//export callbackTxListAccountsModel687eda_Add
func callbackTxListAccountsModel687eda_Add(ptr unsafe.Pointer, tx C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "add"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(tx))
	}

}

func (ptr *TxListAccountsModel) ConnectAdd(f func(tx string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "add") {
			C.TxListAccountsModel687eda_ConnectAdd(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "add")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "add"); signal != nil {
			f := func(tx string) {
				(*(*func(string))(signal))(tx)
				f(tx)
			}
			qt.ConnectSignal(ptr.Pointer(), "add", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "add", unsafe.Pointer(&f))
		}
	}
}

func (ptr *TxListAccountsModel) DisconnectAdd() {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda_DisconnectAdd(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "add")
	}
}

func (ptr *TxListAccountsModel) Add(tx string) {
	if ptr.Pointer() != nil {
		var txC *C.char
		if tx != "" {
			txC = C.CString(tx)
			defer C.free(unsafe.Pointer(txC))
		}
		C.TxListAccountsModel687eda_Add(ptr.Pointer(), C.struct_Moc_PackedString{data: txC, len: C.longlong(len(tx))})
	}
}

func TxListAccountsModel_QRegisterMetaType() int {
	return int(int32(C.TxListAccountsModel687eda_TxListAccountsModel687eda_QRegisterMetaType()))
}

func (ptr *TxListAccountsModel) QRegisterMetaType() int {
	return int(int32(C.TxListAccountsModel687eda_TxListAccountsModel687eda_QRegisterMetaType()))
}

func TxListAccountsModel_QRegisterMetaType2(typeName string) int {
	var typeNameC *C.char
	if typeName != "" {
		typeNameC = C.CString(typeName)
		defer C.free(unsafe.Pointer(typeNameC))
	}
	return int(int32(C.TxListAccountsModel687eda_TxListAccountsModel687eda_QRegisterMetaType2(typeNameC)))
}

func (ptr *TxListAccountsModel) QRegisterMetaType2(typeName string) int {
	var typeNameC *C.char
	if typeName != "" {
		typeNameC = C.CString(typeName)
		defer C.free(unsafe.Pointer(typeNameC))
	}
	return int(int32(C.TxListAccountsModel687eda_TxListAccountsModel687eda_QRegisterMetaType2(typeNameC)))
}

func TxListAccountsModel_QmlRegisterType() int {
	return int(int32(C.TxListAccountsModel687eda_TxListAccountsModel687eda_QmlRegisterType()))
}

func (ptr *TxListAccountsModel) QmlRegisterType() int {
	return int(int32(C.TxListAccountsModel687eda_TxListAccountsModel687eda_QmlRegisterType()))
}

func TxListAccountsModel_QmlRegisterType2(uri string, versionMajor int, versionMinor int, qmlName string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	return int(int32(C.TxListAccountsModel687eda_TxListAccountsModel687eda_QmlRegisterType2(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC)))
}

func (ptr *TxListAccountsModel) QmlRegisterType2(uri string, versionMajor int, versionMinor int, qmlName string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	return int(int32(C.TxListAccountsModel687eda_TxListAccountsModel687eda_QmlRegisterType2(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC)))
}

func TxListAccountsModel_QmlRegisterUncreatableType(uri string, versionMajor int, versionMinor int, qmlName string, message string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	var messageC *C.char
	if message != "" {
		messageC = C.CString(message)
		defer C.free(unsafe.Pointer(messageC))
	}
	return int(int32(C.TxListAccountsModel687eda_TxListAccountsModel687eda_QmlRegisterUncreatableType(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC, C.struct_Moc_PackedString{data: messageC, len: C.longlong(len(message))})))
}

func (ptr *TxListAccountsModel) QmlRegisterUncreatableType(uri string, versionMajor int, versionMinor int, qmlName string, message string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	var messageC *C.char
	if message != "" {
		messageC = C.CString(message)
		defer C.free(unsafe.Pointer(messageC))
	}
	return int(int32(C.TxListAccountsModel687eda_TxListAccountsModel687eda_QmlRegisterUncreatableType(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC, C.struct_Moc_PackedString{data: messageC, len: C.longlong(len(message))})))
}

func (ptr *TxListAccountsModel) ____itemData_keyList_atList(i int) int {
	if ptr.Pointer() != nil {
		return int(int32(C.TxListAccountsModel687eda_____itemData_keyList_atList(ptr.Pointer(), C.int(int32(i)))))
	}
	return 0
}

func (ptr *TxListAccountsModel) ____itemData_keyList_setList(i int) {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda_____itemData_keyList_setList(ptr.Pointer(), C.int(int32(i)))
	}
}

func (ptr *TxListAccountsModel) ____itemData_keyList_newList() unsafe.Pointer {
	return C.TxListAccountsModel687eda_____itemData_keyList_newList(ptr.Pointer())
}

func (ptr *TxListAccountsModel) ____roleNames_keyList_atList(i int) int {
	if ptr.Pointer() != nil {
		return int(int32(C.TxListAccountsModel687eda_____roleNames_keyList_atList(ptr.Pointer(), C.int(int32(i)))))
	}
	return 0
}

func (ptr *TxListAccountsModel) ____roleNames_keyList_setList(i int) {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda_____roleNames_keyList_setList(ptr.Pointer(), C.int(int32(i)))
	}
}

func (ptr *TxListAccountsModel) ____roleNames_keyList_newList() unsafe.Pointer {
	return C.TxListAccountsModel687eda_____roleNames_keyList_newList(ptr.Pointer())
}

func (ptr *TxListAccountsModel) ____setItemData_roles_keyList_atList(i int) int {
	if ptr.Pointer() != nil {
		return int(int32(C.TxListAccountsModel687eda_____setItemData_roles_keyList_atList(ptr.Pointer(), C.int(int32(i)))))
	}
	return 0
}

func (ptr *TxListAccountsModel) ____setItemData_roles_keyList_setList(i int) {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda_____setItemData_roles_keyList_setList(ptr.Pointer(), C.int(int32(i)))
	}
}

func (ptr *TxListAccountsModel) ____setItemData_roles_keyList_newList() unsafe.Pointer {
	return C.TxListAccountsModel687eda_____setItemData_roles_keyList_newList(ptr.Pointer())
}

func (ptr *TxListAccountsModel) __changePersistentIndexList_from_atList(i int) *std_core.QModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQModelIndexFromPointer(C.TxListAccountsModel687eda___changePersistentIndexList_from_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QModelIndex).DestroyQModelIndex)
		return tmpValue
	}
	return nil
}

func (ptr *TxListAccountsModel) __changePersistentIndexList_from_setList(i std_core.QModelIndex_ITF) {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda___changePersistentIndexList_from_setList(ptr.Pointer(), std_core.PointerFromQModelIndex(i))
	}
}

func (ptr *TxListAccountsModel) __changePersistentIndexList_from_newList() unsafe.Pointer {
	return C.TxListAccountsModel687eda___changePersistentIndexList_from_newList(ptr.Pointer())
}

func (ptr *TxListAccountsModel) __changePersistentIndexList_to_atList(i int) *std_core.QModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQModelIndexFromPointer(C.TxListAccountsModel687eda___changePersistentIndexList_to_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QModelIndex).DestroyQModelIndex)
		return tmpValue
	}
	return nil
}

func (ptr *TxListAccountsModel) __changePersistentIndexList_to_setList(i std_core.QModelIndex_ITF) {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda___changePersistentIndexList_to_setList(ptr.Pointer(), std_core.PointerFromQModelIndex(i))
	}
}

func (ptr *TxListAccountsModel) __changePersistentIndexList_to_newList() unsafe.Pointer {
	return C.TxListAccountsModel687eda___changePersistentIndexList_to_newList(ptr.Pointer())
}

func (ptr *TxListAccountsModel) __dataChanged_roles_atList(i int) int {
	if ptr.Pointer() != nil {
		return int(int32(C.TxListAccountsModel687eda___dataChanged_roles_atList(ptr.Pointer(), C.int(int32(i)))))
	}
	return 0
}

func (ptr *TxListAccountsModel) __dataChanged_roles_setList(i int) {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda___dataChanged_roles_setList(ptr.Pointer(), C.int(int32(i)))
	}
}

func (ptr *TxListAccountsModel) __dataChanged_roles_newList() unsafe.Pointer {
	return C.TxListAccountsModel687eda___dataChanged_roles_newList(ptr.Pointer())
}

func (ptr *TxListAccountsModel) __itemData_atList(v int, i int) *std_core.QVariant {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQVariantFromPointer(C.TxListAccountsModel687eda___itemData_atList(ptr.Pointer(), C.int(int32(v)), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QVariant).DestroyQVariant)
		return tmpValue
	}
	return nil
}

func (ptr *TxListAccountsModel) __itemData_setList(key int, i std_core.QVariant_ITF) {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda___itemData_setList(ptr.Pointer(), C.int(int32(key)), std_core.PointerFromQVariant(i))
	}
}

func (ptr *TxListAccountsModel) __itemData_newList() unsafe.Pointer {
	return C.TxListAccountsModel687eda___itemData_newList(ptr.Pointer())
}

func (ptr *TxListAccountsModel) __itemData_keyList() []int {
	if ptr.Pointer() != nil {
		return func(l C.struct_Moc_PackedList) []int {
			out := make([]int, int(l.len))
			tmpList := NewTxListAccountsModelFromPointer(l.data)
			for i := 0; i < len(out); i++ {
				out[i] = tmpList.____itemData_keyList_atList(i)
			}
			return out
		}(C.TxListAccountsModel687eda___itemData_keyList(ptr.Pointer()))
	}
	return make([]int, 0)
}

func (ptr *TxListAccountsModel) __layoutAboutToBeChanged_parents_atList(i int) *std_core.QPersistentModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQPersistentModelIndexFromPointer(C.TxListAccountsModel687eda___layoutAboutToBeChanged_parents_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QPersistentModelIndex).DestroyQPersistentModelIndex)
		return tmpValue
	}
	return nil
}

func (ptr *TxListAccountsModel) __layoutAboutToBeChanged_parents_setList(i std_core.QPersistentModelIndex_ITF) {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda___layoutAboutToBeChanged_parents_setList(ptr.Pointer(), std_core.PointerFromQPersistentModelIndex(i))
	}
}

func (ptr *TxListAccountsModel) __layoutAboutToBeChanged_parents_newList() unsafe.Pointer {
	return C.TxListAccountsModel687eda___layoutAboutToBeChanged_parents_newList(ptr.Pointer())
}

func (ptr *TxListAccountsModel) __layoutChanged_parents_atList(i int) *std_core.QPersistentModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQPersistentModelIndexFromPointer(C.TxListAccountsModel687eda___layoutChanged_parents_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QPersistentModelIndex).DestroyQPersistentModelIndex)
		return tmpValue
	}
	return nil
}

func (ptr *TxListAccountsModel) __layoutChanged_parents_setList(i std_core.QPersistentModelIndex_ITF) {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda___layoutChanged_parents_setList(ptr.Pointer(), std_core.PointerFromQPersistentModelIndex(i))
	}
}

func (ptr *TxListAccountsModel) __layoutChanged_parents_newList() unsafe.Pointer {
	return C.TxListAccountsModel687eda___layoutChanged_parents_newList(ptr.Pointer())
}

func (ptr *TxListAccountsModel) __match_atList(i int) *std_core.QModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQModelIndexFromPointer(C.TxListAccountsModel687eda___match_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QModelIndex).DestroyQModelIndex)
		return tmpValue
	}
	return nil
}

func (ptr *TxListAccountsModel) __match_setList(i std_core.QModelIndex_ITF) {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda___match_setList(ptr.Pointer(), std_core.PointerFromQModelIndex(i))
	}
}

func (ptr *TxListAccountsModel) __match_newList() unsafe.Pointer {
	return C.TxListAccountsModel687eda___match_newList(ptr.Pointer())
}

func (ptr *TxListAccountsModel) __mimeData_indexes_atList(i int) *std_core.QModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQModelIndexFromPointer(C.TxListAccountsModel687eda___mimeData_indexes_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QModelIndex).DestroyQModelIndex)
		return tmpValue
	}
	return nil
}

func (ptr *TxListAccountsModel) __mimeData_indexes_setList(i std_core.QModelIndex_ITF) {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda___mimeData_indexes_setList(ptr.Pointer(), std_core.PointerFromQModelIndex(i))
	}
}

func (ptr *TxListAccountsModel) __mimeData_indexes_newList() unsafe.Pointer {
	return C.TxListAccountsModel687eda___mimeData_indexes_newList(ptr.Pointer())
}

func (ptr *TxListAccountsModel) __persistentIndexList_atList(i int) *std_core.QModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQModelIndexFromPointer(C.TxListAccountsModel687eda___persistentIndexList_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QModelIndex).DestroyQModelIndex)
		return tmpValue
	}
	return nil
}

func (ptr *TxListAccountsModel) __persistentIndexList_setList(i std_core.QModelIndex_ITF) {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda___persistentIndexList_setList(ptr.Pointer(), std_core.PointerFromQModelIndex(i))
	}
}

func (ptr *TxListAccountsModel) __persistentIndexList_newList() unsafe.Pointer {
	return C.TxListAccountsModel687eda___persistentIndexList_newList(ptr.Pointer())
}

func (ptr *TxListAccountsModel) __roleNames_atList(v int, i int) *std_core.QByteArray {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQByteArrayFromPointer(C.TxListAccountsModel687eda___roleNames_atList(ptr.Pointer(), C.int(int32(v)), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QByteArray).DestroyQByteArray)
		return tmpValue
	}
	return nil
}

func (ptr *TxListAccountsModel) __roleNames_setList(key int, i std_core.QByteArray_ITF) {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda___roleNames_setList(ptr.Pointer(), C.int(int32(key)), std_core.PointerFromQByteArray(i))
	}
}

func (ptr *TxListAccountsModel) __roleNames_newList() unsafe.Pointer {
	return C.TxListAccountsModel687eda___roleNames_newList(ptr.Pointer())
}

func (ptr *TxListAccountsModel) __roleNames_keyList() []int {
	if ptr.Pointer() != nil {
		return func(l C.struct_Moc_PackedList) []int {
			out := make([]int, int(l.len))
			tmpList := NewTxListAccountsModelFromPointer(l.data)
			for i := 0; i < len(out); i++ {
				out[i] = tmpList.____roleNames_keyList_atList(i)
			}
			return out
		}(C.TxListAccountsModel687eda___roleNames_keyList(ptr.Pointer()))
	}
	return make([]int, 0)
}

func (ptr *TxListAccountsModel) __setItemData_roles_atList(v int, i int) *std_core.QVariant {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQVariantFromPointer(C.TxListAccountsModel687eda___setItemData_roles_atList(ptr.Pointer(), C.int(int32(v)), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QVariant).DestroyQVariant)
		return tmpValue
	}
	return nil
}

func (ptr *TxListAccountsModel) __setItemData_roles_setList(key int, i std_core.QVariant_ITF) {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda___setItemData_roles_setList(ptr.Pointer(), C.int(int32(key)), std_core.PointerFromQVariant(i))
	}
}

func (ptr *TxListAccountsModel) __setItemData_roles_newList() unsafe.Pointer {
	return C.TxListAccountsModel687eda___setItemData_roles_newList(ptr.Pointer())
}

func (ptr *TxListAccountsModel) __setItemData_roles_keyList() []int {
	if ptr.Pointer() != nil {
		return func(l C.struct_Moc_PackedList) []int {
			out := make([]int, int(l.len))
			tmpList := NewTxListAccountsModelFromPointer(l.data)
			for i := 0; i < len(out); i++ {
				out[i] = tmpList.____setItemData_roles_keyList_atList(i)
			}
			return out
		}(C.TxListAccountsModel687eda___setItemData_roles_keyList(ptr.Pointer()))
	}
	return make([]int, 0)
}

func (ptr *TxListAccountsModel) ____doSetRoleNames_roleNames_keyList_atList(i int) int {
	if ptr.Pointer() != nil {
		return int(int32(C.TxListAccountsModel687eda_____doSetRoleNames_roleNames_keyList_atList(ptr.Pointer(), C.int(int32(i)))))
	}
	return 0
}

func (ptr *TxListAccountsModel) ____doSetRoleNames_roleNames_keyList_setList(i int) {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda_____doSetRoleNames_roleNames_keyList_setList(ptr.Pointer(), C.int(int32(i)))
	}
}

func (ptr *TxListAccountsModel) ____doSetRoleNames_roleNames_keyList_newList() unsafe.Pointer {
	return C.TxListAccountsModel687eda_____doSetRoleNames_roleNames_keyList_newList(ptr.Pointer())
}

func (ptr *TxListAccountsModel) ____setRoleNames_roleNames_keyList_atList(i int) int {
	if ptr.Pointer() != nil {
		return int(int32(C.TxListAccountsModel687eda_____setRoleNames_roleNames_keyList_atList(ptr.Pointer(), C.int(int32(i)))))
	}
	return 0
}

func (ptr *TxListAccountsModel) ____setRoleNames_roleNames_keyList_setList(i int) {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda_____setRoleNames_roleNames_keyList_setList(ptr.Pointer(), C.int(int32(i)))
	}
}

func (ptr *TxListAccountsModel) ____setRoleNames_roleNames_keyList_newList() unsafe.Pointer {
	return C.TxListAccountsModel687eda_____setRoleNames_roleNames_keyList_newList(ptr.Pointer())
}

func (ptr *TxListAccountsModel) __children_atList(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQObjectFromPointer(C.TxListAccountsModel687eda___children_atList(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *TxListAccountsModel) __children_setList(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda___children_setList(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *TxListAccountsModel) __children_newList() unsafe.Pointer {
	return C.TxListAccountsModel687eda___children_newList(ptr.Pointer())
}

func (ptr *TxListAccountsModel) __dynamicPropertyNames_atList(i int) *std_core.QByteArray {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQByteArrayFromPointer(C.TxListAccountsModel687eda___dynamicPropertyNames_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QByteArray).DestroyQByteArray)
		return tmpValue
	}
	return nil
}

func (ptr *TxListAccountsModel) __dynamicPropertyNames_setList(i std_core.QByteArray_ITF) {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda___dynamicPropertyNames_setList(ptr.Pointer(), std_core.PointerFromQByteArray(i))
	}
}

func (ptr *TxListAccountsModel) __dynamicPropertyNames_newList() unsafe.Pointer {
	return C.TxListAccountsModel687eda___dynamicPropertyNames_newList(ptr.Pointer())
}

func (ptr *TxListAccountsModel) __findChildren_atList(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQObjectFromPointer(C.TxListAccountsModel687eda___findChildren_atList(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *TxListAccountsModel) __findChildren_setList(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda___findChildren_setList(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *TxListAccountsModel) __findChildren_newList() unsafe.Pointer {
	return C.TxListAccountsModel687eda___findChildren_newList(ptr.Pointer())
}

func (ptr *TxListAccountsModel) __findChildren_atList3(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQObjectFromPointer(C.TxListAccountsModel687eda___findChildren_atList3(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *TxListAccountsModel) __findChildren_setList3(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda___findChildren_setList3(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *TxListAccountsModel) __findChildren_newList3() unsafe.Pointer {
	return C.TxListAccountsModel687eda___findChildren_newList3(ptr.Pointer())
}

func NewTxListAccountsModel(parent std_core.QObject_ITF) *TxListAccountsModel {
	TxListAccountsModel_QRegisterMetaType()
	tmpValue := NewTxListAccountsModelFromPointer(C.TxListAccountsModel687eda_NewTxListAccountsModel(std_core.PointerFromQObject(parent)))
	if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
		tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
	}
	return tmpValue
}

//export callbackTxListAccountsModel687eda_DestroyTxListAccountsModel
func callbackTxListAccountsModel687eda_DestroyTxListAccountsModel(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "~TxListAccountsModel"); signal != nil {
		(*(*func())(signal))()
	} else {
		NewTxListAccountsModelFromPointer(ptr).DestroyTxListAccountsModelDefault()
	}
}

func (ptr *TxListAccountsModel) ConnectDestroyTxListAccountsModel(f func()) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "~TxListAccountsModel"); signal != nil {
			f := func() {
				(*(*func())(signal))()
				f()
			}
			qt.ConnectSignal(ptr.Pointer(), "~TxListAccountsModel", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "~TxListAccountsModel", unsafe.Pointer(&f))
		}
	}
}

func (ptr *TxListAccountsModel) DisconnectDestroyTxListAccountsModel() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "~TxListAccountsModel")
	}
}

func (ptr *TxListAccountsModel) DestroyTxListAccountsModel() {
	if ptr.Pointer() != nil {

		qt.SetFinalizer(ptr, nil)
		C.TxListAccountsModel687eda_DestroyTxListAccountsModel(ptr.Pointer())
		ptr.SetPointer(nil)
	}
}

func (ptr *TxListAccountsModel) DestroyTxListAccountsModelDefault() {
	if ptr.Pointer() != nil {

		qt.SetFinalizer(ptr, nil)
		C.TxListAccountsModel687eda_DestroyTxListAccountsModelDefault(ptr.Pointer())
		ptr.SetPointer(nil)
	}
}

//export callbackTxListAccountsModel687eda_DropMimeData
func callbackTxListAccountsModel687eda_DropMimeData(ptr unsafe.Pointer, data unsafe.Pointer, action C.longlong, row C.int, column C.int, parent unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "dropMimeData"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QMimeData, std_core.Qt__DropAction, int, int, *std_core.QModelIndex) bool)(signal))(std_core.NewQMimeDataFromPointer(data), std_core.Qt__DropAction(action), int(int32(row)), int(int32(column)), std_core.NewQModelIndexFromPointer(parent)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListAccountsModelFromPointer(ptr).DropMimeDataDefault(std_core.NewQMimeDataFromPointer(data), std_core.Qt__DropAction(action), int(int32(row)), int(int32(column)), std_core.NewQModelIndexFromPointer(parent)))))
}

func (ptr *TxListAccountsModel) DropMimeDataDefault(data std_core.QMimeData_ITF, action std_core.Qt__DropAction, row int, column int, parent std_core.QModelIndex_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListAccountsModel687eda_DropMimeDataDefault(ptr.Pointer(), std_core.PointerFromQMimeData(data), C.longlong(action), C.int(int32(row)), C.int(int32(column)), std_core.PointerFromQModelIndex(parent))) != 0
	}
	return false
}

//export callbackTxListAccountsModel687eda_Flags
func callbackTxListAccountsModel687eda_Flags(ptr unsafe.Pointer, index unsafe.Pointer) C.longlong {
	if signal := qt.GetSignal(ptr, "flags"); signal != nil {
		return C.longlong((*(*func(*std_core.QModelIndex) std_core.Qt__ItemFlag)(signal))(std_core.NewQModelIndexFromPointer(index)))
	}

	return C.longlong(NewTxListAccountsModelFromPointer(ptr).FlagsDefault(std_core.NewQModelIndexFromPointer(index)))
}

func (ptr *TxListAccountsModel) FlagsDefault(index std_core.QModelIndex_ITF) std_core.Qt__ItemFlag {
	if ptr.Pointer() != nil {
		return std_core.Qt__ItemFlag(C.TxListAccountsModel687eda_FlagsDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(index)))
	}
	return 0
}

//export callbackTxListAccountsModel687eda_Index
func callbackTxListAccountsModel687eda_Index(ptr unsafe.Pointer, row C.int, column C.int, parent unsafe.Pointer) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "index"); signal != nil {
		return std_core.PointerFromQModelIndex((*(*func(int, int, *std_core.QModelIndex) *std_core.QModelIndex)(signal))(int(int32(row)), int(int32(column)), std_core.NewQModelIndexFromPointer(parent)))
	}

	return std_core.PointerFromQModelIndex(NewTxListAccountsModelFromPointer(ptr).IndexDefault(int(int32(row)), int(int32(column)), std_core.NewQModelIndexFromPointer(parent)))
}

func (ptr *TxListAccountsModel) IndexDefault(row int, column int, parent std_core.QModelIndex_ITF) *std_core.QModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQModelIndexFromPointer(C.TxListAccountsModel687eda_IndexDefault(ptr.Pointer(), C.int(int32(row)), C.int(int32(column)), std_core.PointerFromQModelIndex(parent)))
		qt.SetFinalizer(tmpValue, (*std_core.QModelIndex).DestroyQModelIndex)
		return tmpValue
	}
	return nil
}

//export callbackTxListAccountsModel687eda_Sibling
func callbackTxListAccountsModel687eda_Sibling(ptr unsafe.Pointer, row C.int, column C.int, idx unsafe.Pointer) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "sibling"); signal != nil {
		return std_core.PointerFromQModelIndex((*(*func(int, int, *std_core.QModelIndex) *std_core.QModelIndex)(signal))(int(int32(row)), int(int32(column)), std_core.NewQModelIndexFromPointer(idx)))
	}

	return std_core.PointerFromQModelIndex(NewTxListAccountsModelFromPointer(ptr).SiblingDefault(int(int32(row)), int(int32(column)), std_core.NewQModelIndexFromPointer(idx)))
}

func (ptr *TxListAccountsModel) SiblingDefault(row int, column int, idx std_core.QModelIndex_ITF) *std_core.QModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQModelIndexFromPointer(C.TxListAccountsModel687eda_SiblingDefault(ptr.Pointer(), C.int(int32(row)), C.int(int32(column)), std_core.PointerFromQModelIndex(idx)))
		qt.SetFinalizer(tmpValue, (*std_core.QModelIndex).DestroyQModelIndex)
		return tmpValue
	}
	return nil
}

//export callbackTxListAccountsModel687eda_Buddy
func callbackTxListAccountsModel687eda_Buddy(ptr unsafe.Pointer, index unsafe.Pointer) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "buddy"); signal != nil {
		return std_core.PointerFromQModelIndex((*(*func(*std_core.QModelIndex) *std_core.QModelIndex)(signal))(std_core.NewQModelIndexFromPointer(index)))
	}

	return std_core.PointerFromQModelIndex(NewTxListAccountsModelFromPointer(ptr).BuddyDefault(std_core.NewQModelIndexFromPointer(index)))
}

func (ptr *TxListAccountsModel) BuddyDefault(index std_core.QModelIndex_ITF) *std_core.QModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQModelIndexFromPointer(C.TxListAccountsModel687eda_BuddyDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(index)))
		qt.SetFinalizer(tmpValue, (*std_core.QModelIndex).DestroyQModelIndex)
		return tmpValue
	}
	return nil
}

//export callbackTxListAccountsModel687eda_CanDropMimeData
func callbackTxListAccountsModel687eda_CanDropMimeData(ptr unsafe.Pointer, data unsafe.Pointer, action C.longlong, row C.int, column C.int, parent unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "canDropMimeData"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QMimeData, std_core.Qt__DropAction, int, int, *std_core.QModelIndex) bool)(signal))(std_core.NewQMimeDataFromPointer(data), std_core.Qt__DropAction(action), int(int32(row)), int(int32(column)), std_core.NewQModelIndexFromPointer(parent)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListAccountsModelFromPointer(ptr).CanDropMimeDataDefault(std_core.NewQMimeDataFromPointer(data), std_core.Qt__DropAction(action), int(int32(row)), int(int32(column)), std_core.NewQModelIndexFromPointer(parent)))))
}

func (ptr *TxListAccountsModel) CanDropMimeDataDefault(data std_core.QMimeData_ITF, action std_core.Qt__DropAction, row int, column int, parent std_core.QModelIndex_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListAccountsModel687eda_CanDropMimeDataDefault(ptr.Pointer(), std_core.PointerFromQMimeData(data), C.longlong(action), C.int(int32(row)), C.int(int32(column)), std_core.PointerFromQModelIndex(parent))) != 0
	}
	return false
}

//export callbackTxListAccountsModel687eda_CanFetchMore
func callbackTxListAccountsModel687eda_CanFetchMore(ptr unsafe.Pointer, parent unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "canFetchMore"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QModelIndex) bool)(signal))(std_core.NewQModelIndexFromPointer(parent)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListAccountsModelFromPointer(ptr).CanFetchMoreDefault(std_core.NewQModelIndexFromPointer(parent)))))
}

func (ptr *TxListAccountsModel) CanFetchMoreDefault(parent std_core.QModelIndex_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListAccountsModel687eda_CanFetchMoreDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(parent))) != 0
	}
	return false
}

//export callbackTxListAccountsModel687eda_ColumnCount
func callbackTxListAccountsModel687eda_ColumnCount(ptr unsafe.Pointer, parent unsafe.Pointer) C.int {
	if signal := qt.GetSignal(ptr, "columnCount"); signal != nil {
		return C.int(int32((*(*func(*std_core.QModelIndex) int)(signal))(std_core.NewQModelIndexFromPointer(parent))))
	}

	return C.int(int32(NewTxListAccountsModelFromPointer(ptr).ColumnCountDefault(std_core.NewQModelIndexFromPointer(parent))))
}

func (ptr *TxListAccountsModel) ColumnCountDefault(parent std_core.QModelIndex_ITF) int {
	if ptr.Pointer() != nil {
		return int(int32(C.TxListAccountsModel687eda_ColumnCountDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(parent))))
	}
	return 0
}

//export callbackTxListAccountsModel687eda_ColumnsAboutToBeInserted
func callbackTxListAccountsModel687eda_ColumnsAboutToBeInserted(ptr unsafe.Pointer, parent unsafe.Pointer, first C.int, last C.int) {
	if signal := qt.GetSignal(ptr, "columnsAboutToBeInserted"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(first)), int(int32(last)))
	}

}

//export callbackTxListAccountsModel687eda_ColumnsAboutToBeMoved
func callbackTxListAccountsModel687eda_ColumnsAboutToBeMoved(ptr unsafe.Pointer, sourceParent unsafe.Pointer, sourceStart C.int, sourceEnd C.int, destinationParent unsafe.Pointer, destinationColumn C.int) {
	if signal := qt.GetSignal(ptr, "columnsAboutToBeMoved"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int, *std_core.QModelIndex, int))(signal))(std_core.NewQModelIndexFromPointer(sourceParent), int(int32(sourceStart)), int(int32(sourceEnd)), std_core.NewQModelIndexFromPointer(destinationParent), int(int32(destinationColumn)))
	}

}

//export callbackTxListAccountsModel687eda_ColumnsAboutToBeRemoved
func callbackTxListAccountsModel687eda_ColumnsAboutToBeRemoved(ptr unsafe.Pointer, parent unsafe.Pointer, first C.int, last C.int) {
	if signal := qt.GetSignal(ptr, "columnsAboutToBeRemoved"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(first)), int(int32(last)))
	}

}

//export callbackTxListAccountsModel687eda_ColumnsInserted
func callbackTxListAccountsModel687eda_ColumnsInserted(ptr unsafe.Pointer, parent unsafe.Pointer, first C.int, last C.int) {
	if signal := qt.GetSignal(ptr, "columnsInserted"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(first)), int(int32(last)))
	}

}

//export callbackTxListAccountsModel687eda_ColumnsMoved
func callbackTxListAccountsModel687eda_ColumnsMoved(ptr unsafe.Pointer, parent unsafe.Pointer, start C.int, end C.int, destination unsafe.Pointer, column C.int) {
	if signal := qt.GetSignal(ptr, "columnsMoved"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int, *std_core.QModelIndex, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(start)), int(int32(end)), std_core.NewQModelIndexFromPointer(destination), int(int32(column)))
	}

}

//export callbackTxListAccountsModel687eda_ColumnsRemoved
func callbackTxListAccountsModel687eda_ColumnsRemoved(ptr unsafe.Pointer, parent unsafe.Pointer, first C.int, last C.int) {
	if signal := qt.GetSignal(ptr, "columnsRemoved"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(first)), int(int32(last)))
	}

}

//export callbackTxListAccountsModel687eda_Data
func callbackTxListAccountsModel687eda_Data(ptr unsafe.Pointer, index unsafe.Pointer, role C.int) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "data"); signal != nil {
		return std_core.PointerFromQVariant((*(*func(*std_core.QModelIndex, int) *std_core.QVariant)(signal))(std_core.NewQModelIndexFromPointer(index), int(int32(role))))
	}

	return std_core.PointerFromQVariant(NewTxListAccountsModelFromPointer(ptr).DataDefault(std_core.NewQModelIndexFromPointer(index), int(int32(role))))
}

func (ptr *TxListAccountsModel) DataDefault(index std_core.QModelIndex_ITF, role int) *std_core.QVariant {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQVariantFromPointer(C.TxListAccountsModel687eda_DataDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(index), C.int(int32(role))))
		qt.SetFinalizer(tmpValue, (*std_core.QVariant).DestroyQVariant)
		return tmpValue
	}
	return nil
}

//export callbackTxListAccountsModel687eda_DataChanged
func callbackTxListAccountsModel687eda_DataChanged(ptr unsafe.Pointer, topLeft unsafe.Pointer, bottomRight unsafe.Pointer, roles C.struct_Moc_PackedList) {
	if signal := qt.GetSignal(ptr, "dataChanged"); signal != nil {
		(*(*func(*std_core.QModelIndex, *std_core.QModelIndex, []int))(signal))(std_core.NewQModelIndexFromPointer(topLeft), std_core.NewQModelIndexFromPointer(bottomRight), func(l C.struct_Moc_PackedList) []int {
			out := make([]int, int(l.len))
			tmpList := NewTxListAccountsModelFromPointer(l.data)
			for i := 0; i < len(out); i++ {
				out[i] = tmpList.__dataChanged_roles_atList(i)
			}
			return out
		}(roles))
	}

}

//export callbackTxListAccountsModel687eda_FetchMore
func callbackTxListAccountsModel687eda_FetchMore(ptr unsafe.Pointer, parent unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "fetchMore"); signal != nil {
		(*(*func(*std_core.QModelIndex))(signal))(std_core.NewQModelIndexFromPointer(parent))
	} else {
		NewTxListAccountsModelFromPointer(ptr).FetchMoreDefault(std_core.NewQModelIndexFromPointer(parent))
	}
}

func (ptr *TxListAccountsModel) FetchMoreDefault(parent std_core.QModelIndex_ITF) {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda_FetchMoreDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(parent))
	}
}

//export callbackTxListAccountsModel687eda_HasChildren
func callbackTxListAccountsModel687eda_HasChildren(ptr unsafe.Pointer, parent unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "hasChildren"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QModelIndex) bool)(signal))(std_core.NewQModelIndexFromPointer(parent)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListAccountsModelFromPointer(ptr).HasChildrenDefault(std_core.NewQModelIndexFromPointer(parent)))))
}

func (ptr *TxListAccountsModel) HasChildrenDefault(parent std_core.QModelIndex_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListAccountsModel687eda_HasChildrenDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(parent))) != 0
	}
	return false
}

//export callbackTxListAccountsModel687eda_HeaderData
func callbackTxListAccountsModel687eda_HeaderData(ptr unsafe.Pointer, section C.int, orientation C.longlong, role C.int) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "headerData"); signal != nil {
		return std_core.PointerFromQVariant((*(*func(int, std_core.Qt__Orientation, int) *std_core.QVariant)(signal))(int(int32(section)), std_core.Qt__Orientation(orientation), int(int32(role))))
	}

	return std_core.PointerFromQVariant(NewTxListAccountsModelFromPointer(ptr).HeaderDataDefault(int(int32(section)), std_core.Qt__Orientation(orientation), int(int32(role))))
}

func (ptr *TxListAccountsModel) HeaderDataDefault(section int, orientation std_core.Qt__Orientation, role int) *std_core.QVariant {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQVariantFromPointer(C.TxListAccountsModel687eda_HeaderDataDefault(ptr.Pointer(), C.int(int32(section)), C.longlong(orientation), C.int(int32(role))))
		qt.SetFinalizer(tmpValue, (*std_core.QVariant).DestroyQVariant)
		return tmpValue
	}
	return nil
}

//export callbackTxListAccountsModel687eda_HeaderDataChanged
func callbackTxListAccountsModel687eda_HeaderDataChanged(ptr unsafe.Pointer, orientation C.longlong, first C.int, last C.int) {
	if signal := qt.GetSignal(ptr, "headerDataChanged"); signal != nil {
		(*(*func(std_core.Qt__Orientation, int, int))(signal))(std_core.Qt__Orientation(orientation), int(int32(first)), int(int32(last)))
	}

}

//export callbackTxListAccountsModel687eda_InsertColumns
func callbackTxListAccountsModel687eda_InsertColumns(ptr unsafe.Pointer, column C.int, count C.int, parent unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "insertColumns"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(int, int, *std_core.QModelIndex) bool)(signal))(int(int32(column)), int(int32(count)), std_core.NewQModelIndexFromPointer(parent)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListAccountsModelFromPointer(ptr).InsertColumnsDefault(int(int32(column)), int(int32(count)), std_core.NewQModelIndexFromPointer(parent)))))
}

func (ptr *TxListAccountsModel) InsertColumnsDefault(column int, count int, parent std_core.QModelIndex_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListAccountsModel687eda_InsertColumnsDefault(ptr.Pointer(), C.int(int32(column)), C.int(int32(count)), std_core.PointerFromQModelIndex(parent))) != 0
	}
	return false
}

//export callbackTxListAccountsModel687eda_InsertRows
func callbackTxListAccountsModel687eda_InsertRows(ptr unsafe.Pointer, row C.int, count C.int, parent unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "insertRows"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(int, int, *std_core.QModelIndex) bool)(signal))(int(int32(row)), int(int32(count)), std_core.NewQModelIndexFromPointer(parent)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListAccountsModelFromPointer(ptr).InsertRowsDefault(int(int32(row)), int(int32(count)), std_core.NewQModelIndexFromPointer(parent)))))
}

func (ptr *TxListAccountsModel) InsertRowsDefault(row int, count int, parent std_core.QModelIndex_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListAccountsModel687eda_InsertRowsDefault(ptr.Pointer(), C.int(int32(row)), C.int(int32(count)), std_core.PointerFromQModelIndex(parent))) != 0
	}
	return false
}

//export callbackTxListAccountsModel687eda_ItemData
func callbackTxListAccountsModel687eda_ItemData(ptr unsafe.Pointer, index unsafe.Pointer) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "itemData"); signal != nil {
		return func() unsafe.Pointer {
			tmpList := NewTxListAccountsModelFromPointer(NewTxListAccountsModelFromPointer(nil).__itemData_newList())
			for k, v := range (*(*func(*std_core.QModelIndex) map[int]*std_core.QVariant)(signal))(std_core.NewQModelIndexFromPointer(index)) {
				tmpList.__itemData_setList(k, v)
			}
			return tmpList.Pointer()
		}()
	}

	return func() unsafe.Pointer {
		tmpList := NewTxListAccountsModelFromPointer(NewTxListAccountsModelFromPointer(nil).__itemData_newList())
		for k, v := range NewTxListAccountsModelFromPointer(ptr).ItemDataDefault(std_core.NewQModelIndexFromPointer(index)) {
			tmpList.__itemData_setList(k, v)
		}
		return tmpList.Pointer()
	}()
}

func (ptr *TxListAccountsModel) ItemDataDefault(index std_core.QModelIndex_ITF) map[int]*std_core.QVariant {
	if ptr.Pointer() != nil {
		return func(l C.struct_Moc_PackedList) map[int]*std_core.QVariant {
			out := make(map[int]*std_core.QVariant, int(l.len))
			tmpList := NewTxListAccountsModelFromPointer(l.data)
			for i, v := range tmpList.__itemData_keyList() {
				out[v] = tmpList.__itemData_atList(v, i)
			}
			return out
		}(C.TxListAccountsModel687eda_ItemDataDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(index)))
	}
	return make(map[int]*std_core.QVariant, 0)
}

//export callbackTxListAccountsModel687eda_LayoutAboutToBeChanged
func callbackTxListAccountsModel687eda_LayoutAboutToBeChanged(ptr unsafe.Pointer, parents C.struct_Moc_PackedList, hint C.longlong) {
	if signal := qt.GetSignal(ptr, "layoutAboutToBeChanged"); signal != nil {
		(*(*func([]*std_core.QPersistentModelIndex, std_core.QAbstractItemModel__LayoutChangeHint))(signal))(func(l C.struct_Moc_PackedList) []*std_core.QPersistentModelIndex {
			out := make([]*std_core.QPersistentModelIndex, int(l.len))
			tmpList := NewTxListAccountsModelFromPointer(l.data)
			for i := 0; i < len(out); i++ {
				out[i] = tmpList.__layoutAboutToBeChanged_parents_atList(i)
			}
			return out
		}(parents), std_core.QAbstractItemModel__LayoutChangeHint(hint))
	}

}

//export callbackTxListAccountsModel687eda_LayoutChanged
func callbackTxListAccountsModel687eda_LayoutChanged(ptr unsafe.Pointer, parents C.struct_Moc_PackedList, hint C.longlong) {
	if signal := qt.GetSignal(ptr, "layoutChanged"); signal != nil {
		(*(*func([]*std_core.QPersistentModelIndex, std_core.QAbstractItemModel__LayoutChangeHint))(signal))(func(l C.struct_Moc_PackedList) []*std_core.QPersistentModelIndex {
			out := make([]*std_core.QPersistentModelIndex, int(l.len))
			tmpList := NewTxListAccountsModelFromPointer(l.data)
			for i := 0; i < len(out); i++ {
				out[i] = tmpList.__layoutChanged_parents_atList(i)
			}
			return out
		}(parents), std_core.QAbstractItemModel__LayoutChangeHint(hint))
	}

}

//export callbackTxListAccountsModel687eda_Match
func callbackTxListAccountsModel687eda_Match(ptr unsafe.Pointer, start unsafe.Pointer, role C.int, value unsafe.Pointer, hits C.int, flags C.longlong) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "match"); signal != nil {
		return func() unsafe.Pointer {
			tmpList := NewTxListAccountsModelFromPointer(NewTxListAccountsModelFromPointer(nil).__match_newList())
			for _, v := range (*(*func(*std_core.QModelIndex, int, *std_core.QVariant, int, std_core.Qt__MatchFlag) []*std_core.QModelIndex)(signal))(std_core.NewQModelIndexFromPointer(start), int(int32(role)), std_core.NewQVariantFromPointer(value), int(int32(hits)), std_core.Qt__MatchFlag(flags)) {
				tmpList.__match_setList(v)
			}
			return tmpList.Pointer()
		}()
	}

	return func() unsafe.Pointer {
		tmpList := NewTxListAccountsModelFromPointer(NewTxListAccountsModelFromPointer(nil).__match_newList())
		for _, v := range NewTxListAccountsModelFromPointer(ptr).MatchDefault(std_core.NewQModelIndexFromPointer(start), int(int32(role)), std_core.NewQVariantFromPointer(value), int(int32(hits)), std_core.Qt__MatchFlag(flags)) {
			tmpList.__match_setList(v)
		}
		return tmpList.Pointer()
	}()
}

func (ptr *TxListAccountsModel) MatchDefault(start std_core.QModelIndex_ITF, role int, value std_core.QVariant_ITF, hits int, flags std_core.Qt__MatchFlag) []*std_core.QModelIndex {
	if ptr.Pointer() != nil {
		return func(l C.struct_Moc_PackedList) []*std_core.QModelIndex {
			out := make([]*std_core.QModelIndex, int(l.len))
			tmpList := NewTxListAccountsModelFromPointer(l.data)
			for i := 0; i < len(out); i++ {
				out[i] = tmpList.__match_atList(i)
			}
			return out
		}(C.TxListAccountsModel687eda_MatchDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(start), C.int(int32(role)), std_core.PointerFromQVariant(value), C.int(int32(hits)), C.longlong(flags)))
	}
	return make([]*std_core.QModelIndex, 0)
}

//export callbackTxListAccountsModel687eda_MimeData
func callbackTxListAccountsModel687eda_MimeData(ptr unsafe.Pointer, indexes C.struct_Moc_PackedList) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "mimeData"); signal != nil {
		return std_core.PointerFromQMimeData((*(*func([]*std_core.QModelIndex) *std_core.QMimeData)(signal))(func(l C.struct_Moc_PackedList) []*std_core.QModelIndex {
			out := make([]*std_core.QModelIndex, int(l.len))
			tmpList := NewTxListAccountsModelFromPointer(l.data)
			for i := 0; i < len(out); i++ {
				out[i] = tmpList.__mimeData_indexes_atList(i)
			}
			return out
		}(indexes)))
	}

	return std_core.PointerFromQMimeData(NewTxListAccountsModelFromPointer(ptr).MimeDataDefault(func(l C.struct_Moc_PackedList) []*std_core.QModelIndex {
		out := make([]*std_core.QModelIndex, int(l.len))
		tmpList := NewTxListAccountsModelFromPointer(l.data)
		for i := 0; i < len(out); i++ {
			out[i] = tmpList.__mimeData_indexes_atList(i)
		}
		return out
	}(indexes)))
}

func (ptr *TxListAccountsModel) MimeDataDefault(indexes []*std_core.QModelIndex) *std_core.QMimeData {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQMimeDataFromPointer(C.TxListAccountsModel687eda_MimeDataDefault(ptr.Pointer(), func() unsafe.Pointer {
			tmpList := NewTxListAccountsModelFromPointer(NewTxListAccountsModelFromPointer(nil).__mimeData_indexes_newList())
			for _, v := range indexes {
				tmpList.__mimeData_indexes_setList(v)
			}
			return tmpList.Pointer()
		}()))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

//export callbackTxListAccountsModel687eda_MimeTypes
func callbackTxListAccountsModel687eda_MimeTypes(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "mimeTypes"); signal != nil {
		tempVal := (*(*func() []string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(strings.Join(tempVal, "¡¦!")), len: C.longlong(len(strings.Join(tempVal, "¡¦!")))}
	}
	tempVal := NewTxListAccountsModelFromPointer(ptr).MimeTypesDefault()
	return C.struct_Moc_PackedString{data: C.CString(strings.Join(tempVal, "¡¦!")), len: C.longlong(len(strings.Join(tempVal, "¡¦!")))}
}

func (ptr *TxListAccountsModel) MimeTypesDefault() []string {
	if ptr.Pointer() != nil {
		return unpackStringList(cGoUnpackString(C.TxListAccountsModel687eda_MimeTypesDefault(ptr.Pointer())))
	}
	return make([]string, 0)
}

//export callbackTxListAccountsModel687eda_ModelAboutToBeReset
func callbackTxListAccountsModel687eda_ModelAboutToBeReset(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "modelAboutToBeReset"); signal != nil {
		(*(*func())(signal))()
	}

}

//export callbackTxListAccountsModel687eda_ModelReset
func callbackTxListAccountsModel687eda_ModelReset(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "modelReset"); signal != nil {
		(*(*func())(signal))()
	}

}

//export callbackTxListAccountsModel687eda_MoveColumns
func callbackTxListAccountsModel687eda_MoveColumns(ptr unsafe.Pointer, sourceParent unsafe.Pointer, sourceColumn C.int, count C.int, destinationParent unsafe.Pointer, destinationChild C.int) C.char {
	if signal := qt.GetSignal(ptr, "moveColumns"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QModelIndex, int, int, *std_core.QModelIndex, int) bool)(signal))(std_core.NewQModelIndexFromPointer(sourceParent), int(int32(sourceColumn)), int(int32(count)), std_core.NewQModelIndexFromPointer(destinationParent), int(int32(destinationChild))))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListAccountsModelFromPointer(ptr).MoveColumnsDefault(std_core.NewQModelIndexFromPointer(sourceParent), int(int32(sourceColumn)), int(int32(count)), std_core.NewQModelIndexFromPointer(destinationParent), int(int32(destinationChild))))))
}

func (ptr *TxListAccountsModel) MoveColumnsDefault(sourceParent std_core.QModelIndex_ITF, sourceColumn int, count int, destinationParent std_core.QModelIndex_ITF, destinationChild int) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListAccountsModel687eda_MoveColumnsDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(sourceParent), C.int(int32(sourceColumn)), C.int(int32(count)), std_core.PointerFromQModelIndex(destinationParent), C.int(int32(destinationChild)))) != 0
	}
	return false
}

//export callbackTxListAccountsModel687eda_MoveRows
func callbackTxListAccountsModel687eda_MoveRows(ptr unsafe.Pointer, sourceParent unsafe.Pointer, sourceRow C.int, count C.int, destinationParent unsafe.Pointer, destinationChild C.int) C.char {
	if signal := qt.GetSignal(ptr, "moveRows"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QModelIndex, int, int, *std_core.QModelIndex, int) bool)(signal))(std_core.NewQModelIndexFromPointer(sourceParent), int(int32(sourceRow)), int(int32(count)), std_core.NewQModelIndexFromPointer(destinationParent), int(int32(destinationChild))))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListAccountsModelFromPointer(ptr).MoveRowsDefault(std_core.NewQModelIndexFromPointer(sourceParent), int(int32(sourceRow)), int(int32(count)), std_core.NewQModelIndexFromPointer(destinationParent), int(int32(destinationChild))))))
}

func (ptr *TxListAccountsModel) MoveRowsDefault(sourceParent std_core.QModelIndex_ITF, sourceRow int, count int, destinationParent std_core.QModelIndex_ITF, destinationChild int) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListAccountsModel687eda_MoveRowsDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(sourceParent), C.int(int32(sourceRow)), C.int(int32(count)), std_core.PointerFromQModelIndex(destinationParent), C.int(int32(destinationChild)))) != 0
	}
	return false
}

//export callbackTxListAccountsModel687eda_Parent
func callbackTxListAccountsModel687eda_Parent(ptr unsafe.Pointer, index unsafe.Pointer) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "parent"); signal != nil {
		return std_core.PointerFromQModelIndex((*(*func(*std_core.QModelIndex) *std_core.QModelIndex)(signal))(std_core.NewQModelIndexFromPointer(index)))
	}

	return std_core.PointerFromQModelIndex(NewTxListAccountsModelFromPointer(ptr).ParentDefault(std_core.NewQModelIndexFromPointer(index)))
}

func (ptr *TxListAccountsModel) ParentDefault(index std_core.QModelIndex_ITF) *std_core.QModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQModelIndexFromPointer(C.TxListAccountsModel687eda_ParentDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(index)))
		qt.SetFinalizer(tmpValue, (*std_core.QModelIndex).DestroyQModelIndex)
		return tmpValue
	}
	return nil
}

//export callbackTxListAccountsModel687eda_RemoveColumns
func callbackTxListAccountsModel687eda_RemoveColumns(ptr unsafe.Pointer, column C.int, count C.int, parent unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "removeColumns"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(int, int, *std_core.QModelIndex) bool)(signal))(int(int32(column)), int(int32(count)), std_core.NewQModelIndexFromPointer(parent)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListAccountsModelFromPointer(ptr).RemoveColumnsDefault(int(int32(column)), int(int32(count)), std_core.NewQModelIndexFromPointer(parent)))))
}

func (ptr *TxListAccountsModel) RemoveColumnsDefault(column int, count int, parent std_core.QModelIndex_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListAccountsModel687eda_RemoveColumnsDefault(ptr.Pointer(), C.int(int32(column)), C.int(int32(count)), std_core.PointerFromQModelIndex(parent))) != 0
	}
	return false
}

//export callbackTxListAccountsModel687eda_RemoveRows
func callbackTxListAccountsModel687eda_RemoveRows(ptr unsafe.Pointer, row C.int, count C.int, parent unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "removeRows"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(int, int, *std_core.QModelIndex) bool)(signal))(int(int32(row)), int(int32(count)), std_core.NewQModelIndexFromPointer(parent)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListAccountsModelFromPointer(ptr).RemoveRowsDefault(int(int32(row)), int(int32(count)), std_core.NewQModelIndexFromPointer(parent)))))
}

func (ptr *TxListAccountsModel) RemoveRowsDefault(row int, count int, parent std_core.QModelIndex_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListAccountsModel687eda_RemoveRowsDefault(ptr.Pointer(), C.int(int32(row)), C.int(int32(count)), std_core.PointerFromQModelIndex(parent))) != 0
	}
	return false
}

//export callbackTxListAccountsModel687eda_ResetInternalData
func callbackTxListAccountsModel687eda_ResetInternalData(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "resetInternalData"); signal != nil {
		(*(*func())(signal))()
	} else {
		NewTxListAccountsModelFromPointer(ptr).ResetInternalDataDefault()
	}
}

func (ptr *TxListAccountsModel) ResetInternalDataDefault() {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda_ResetInternalDataDefault(ptr.Pointer())
	}
}

//export callbackTxListAccountsModel687eda_Revert
func callbackTxListAccountsModel687eda_Revert(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "revert"); signal != nil {
		(*(*func())(signal))()
	} else {
		NewTxListAccountsModelFromPointer(ptr).RevertDefault()
	}
}

func (ptr *TxListAccountsModel) RevertDefault() {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda_RevertDefault(ptr.Pointer())
	}
}

//export callbackTxListAccountsModel687eda_RoleNames
func callbackTxListAccountsModel687eda_RoleNames(ptr unsafe.Pointer) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "roleNames"); signal != nil {
		return func() unsafe.Pointer {
			tmpList := NewTxListAccountsModelFromPointer(NewTxListAccountsModelFromPointer(nil).__roleNames_newList())
			for k, v := range (*(*func() map[int]*std_core.QByteArray)(signal))() {
				tmpList.__roleNames_setList(k, v)
			}
			return tmpList.Pointer()
		}()
	}

	return func() unsafe.Pointer {
		tmpList := NewTxListAccountsModelFromPointer(NewTxListAccountsModelFromPointer(nil).__roleNames_newList())
		for k, v := range NewTxListAccountsModelFromPointer(ptr).RoleNamesDefault() {
			tmpList.__roleNames_setList(k, v)
		}
		return tmpList.Pointer()
	}()
}

func (ptr *TxListAccountsModel) RoleNamesDefault() map[int]*std_core.QByteArray {
	if ptr.Pointer() != nil {
		return func(l C.struct_Moc_PackedList) map[int]*std_core.QByteArray {
			out := make(map[int]*std_core.QByteArray, int(l.len))
			tmpList := NewTxListAccountsModelFromPointer(l.data)
			for i, v := range tmpList.__roleNames_keyList() {
				out[v] = tmpList.__roleNames_atList(v, i)
			}
			return out
		}(C.TxListAccountsModel687eda_RoleNamesDefault(ptr.Pointer()))
	}
	return make(map[int]*std_core.QByteArray, 0)
}

//export callbackTxListAccountsModel687eda_RowCount
func callbackTxListAccountsModel687eda_RowCount(ptr unsafe.Pointer, parent unsafe.Pointer) C.int {
	if signal := qt.GetSignal(ptr, "rowCount"); signal != nil {
		return C.int(int32((*(*func(*std_core.QModelIndex) int)(signal))(std_core.NewQModelIndexFromPointer(parent))))
	}

	return C.int(int32(NewTxListAccountsModelFromPointer(ptr).RowCountDefault(std_core.NewQModelIndexFromPointer(parent))))
}

func (ptr *TxListAccountsModel) RowCountDefault(parent std_core.QModelIndex_ITF) int {
	if ptr.Pointer() != nil {
		return int(int32(C.TxListAccountsModel687eda_RowCountDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(parent))))
	}
	return 0
}

//export callbackTxListAccountsModel687eda_RowsAboutToBeInserted
func callbackTxListAccountsModel687eda_RowsAboutToBeInserted(ptr unsafe.Pointer, parent unsafe.Pointer, start C.int, end C.int) {
	if signal := qt.GetSignal(ptr, "rowsAboutToBeInserted"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(start)), int(int32(end)))
	}

}

//export callbackTxListAccountsModel687eda_RowsAboutToBeMoved
func callbackTxListAccountsModel687eda_RowsAboutToBeMoved(ptr unsafe.Pointer, sourceParent unsafe.Pointer, sourceStart C.int, sourceEnd C.int, destinationParent unsafe.Pointer, destinationRow C.int) {
	if signal := qt.GetSignal(ptr, "rowsAboutToBeMoved"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int, *std_core.QModelIndex, int))(signal))(std_core.NewQModelIndexFromPointer(sourceParent), int(int32(sourceStart)), int(int32(sourceEnd)), std_core.NewQModelIndexFromPointer(destinationParent), int(int32(destinationRow)))
	}

}

//export callbackTxListAccountsModel687eda_RowsAboutToBeRemoved
func callbackTxListAccountsModel687eda_RowsAboutToBeRemoved(ptr unsafe.Pointer, parent unsafe.Pointer, first C.int, last C.int) {
	if signal := qt.GetSignal(ptr, "rowsAboutToBeRemoved"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(first)), int(int32(last)))
	}

}

//export callbackTxListAccountsModel687eda_RowsInserted
func callbackTxListAccountsModel687eda_RowsInserted(ptr unsafe.Pointer, parent unsafe.Pointer, first C.int, last C.int) {
	if signal := qt.GetSignal(ptr, "rowsInserted"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(first)), int(int32(last)))
	}

}

//export callbackTxListAccountsModel687eda_RowsMoved
func callbackTxListAccountsModel687eda_RowsMoved(ptr unsafe.Pointer, parent unsafe.Pointer, start C.int, end C.int, destination unsafe.Pointer, row C.int) {
	if signal := qt.GetSignal(ptr, "rowsMoved"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int, *std_core.QModelIndex, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(start)), int(int32(end)), std_core.NewQModelIndexFromPointer(destination), int(int32(row)))
	}

}

//export callbackTxListAccountsModel687eda_RowsRemoved
func callbackTxListAccountsModel687eda_RowsRemoved(ptr unsafe.Pointer, parent unsafe.Pointer, first C.int, last C.int) {
	if signal := qt.GetSignal(ptr, "rowsRemoved"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(first)), int(int32(last)))
	}

}

//export callbackTxListAccountsModel687eda_SetData
func callbackTxListAccountsModel687eda_SetData(ptr unsafe.Pointer, index unsafe.Pointer, value unsafe.Pointer, role C.int) C.char {
	if signal := qt.GetSignal(ptr, "setData"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QModelIndex, *std_core.QVariant, int) bool)(signal))(std_core.NewQModelIndexFromPointer(index), std_core.NewQVariantFromPointer(value), int(int32(role))))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListAccountsModelFromPointer(ptr).SetDataDefault(std_core.NewQModelIndexFromPointer(index), std_core.NewQVariantFromPointer(value), int(int32(role))))))
}

func (ptr *TxListAccountsModel) SetDataDefault(index std_core.QModelIndex_ITF, value std_core.QVariant_ITF, role int) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListAccountsModel687eda_SetDataDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(index), std_core.PointerFromQVariant(value), C.int(int32(role)))) != 0
	}
	return false
}

//export callbackTxListAccountsModel687eda_SetHeaderData
func callbackTxListAccountsModel687eda_SetHeaderData(ptr unsafe.Pointer, section C.int, orientation C.longlong, value unsafe.Pointer, role C.int) C.char {
	if signal := qt.GetSignal(ptr, "setHeaderData"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(int, std_core.Qt__Orientation, *std_core.QVariant, int) bool)(signal))(int(int32(section)), std_core.Qt__Orientation(orientation), std_core.NewQVariantFromPointer(value), int(int32(role))))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListAccountsModelFromPointer(ptr).SetHeaderDataDefault(int(int32(section)), std_core.Qt__Orientation(orientation), std_core.NewQVariantFromPointer(value), int(int32(role))))))
}

func (ptr *TxListAccountsModel) SetHeaderDataDefault(section int, orientation std_core.Qt__Orientation, value std_core.QVariant_ITF, role int) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListAccountsModel687eda_SetHeaderDataDefault(ptr.Pointer(), C.int(int32(section)), C.longlong(orientation), std_core.PointerFromQVariant(value), C.int(int32(role)))) != 0
	}
	return false
}

//export callbackTxListAccountsModel687eda_SetItemData
func callbackTxListAccountsModel687eda_SetItemData(ptr unsafe.Pointer, index unsafe.Pointer, roles C.struct_Moc_PackedList) C.char {
	if signal := qt.GetSignal(ptr, "setItemData"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QModelIndex, map[int]*std_core.QVariant) bool)(signal))(std_core.NewQModelIndexFromPointer(index), func(l C.struct_Moc_PackedList) map[int]*std_core.QVariant {
			out := make(map[int]*std_core.QVariant, int(l.len))
			tmpList := NewTxListAccountsModelFromPointer(l.data)
			for i, v := range tmpList.__setItemData_roles_keyList() {
				out[v] = tmpList.__setItemData_roles_atList(v, i)
			}
			return out
		}(roles)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListAccountsModelFromPointer(ptr).SetItemDataDefault(std_core.NewQModelIndexFromPointer(index), func(l C.struct_Moc_PackedList) map[int]*std_core.QVariant {
		out := make(map[int]*std_core.QVariant, int(l.len))
		tmpList := NewTxListAccountsModelFromPointer(l.data)
		for i, v := range tmpList.__setItemData_roles_keyList() {
			out[v] = tmpList.__setItemData_roles_atList(v, i)
		}
		return out
	}(roles)))))
}

func (ptr *TxListAccountsModel) SetItemDataDefault(index std_core.QModelIndex_ITF, roles map[int]*std_core.QVariant) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListAccountsModel687eda_SetItemDataDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(index), func() unsafe.Pointer {
			tmpList := NewTxListAccountsModelFromPointer(NewTxListAccountsModelFromPointer(nil).__setItemData_roles_newList())
			for k, v := range roles {
				tmpList.__setItemData_roles_setList(k, v)
			}
			return tmpList.Pointer()
		}())) != 0
	}
	return false
}

//export callbackTxListAccountsModel687eda_Sort
func callbackTxListAccountsModel687eda_Sort(ptr unsafe.Pointer, column C.int, order C.longlong) {
	if signal := qt.GetSignal(ptr, "sort"); signal != nil {
		(*(*func(int, std_core.Qt__SortOrder))(signal))(int(int32(column)), std_core.Qt__SortOrder(order))
	} else {
		NewTxListAccountsModelFromPointer(ptr).SortDefault(int(int32(column)), std_core.Qt__SortOrder(order))
	}
}

func (ptr *TxListAccountsModel) SortDefault(column int, order std_core.Qt__SortOrder) {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda_SortDefault(ptr.Pointer(), C.int(int32(column)), C.longlong(order))
	}
}

//export callbackTxListAccountsModel687eda_Span
func callbackTxListAccountsModel687eda_Span(ptr unsafe.Pointer, index unsafe.Pointer) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "span"); signal != nil {
		return std_core.PointerFromQSize((*(*func(*std_core.QModelIndex) *std_core.QSize)(signal))(std_core.NewQModelIndexFromPointer(index)))
	}

	return std_core.PointerFromQSize(NewTxListAccountsModelFromPointer(ptr).SpanDefault(std_core.NewQModelIndexFromPointer(index)))
}

func (ptr *TxListAccountsModel) SpanDefault(index std_core.QModelIndex_ITF) *std_core.QSize {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQSizeFromPointer(C.TxListAccountsModel687eda_SpanDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(index)))
		qt.SetFinalizer(tmpValue, (*std_core.QSize).DestroyQSize)
		return tmpValue
	}
	return nil
}

//export callbackTxListAccountsModel687eda_Submit
func callbackTxListAccountsModel687eda_Submit(ptr unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "submit"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func() bool)(signal))())))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListAccountsModelFromPointer(ptr).SubmitDefault())))
}

func (ptr *TxListAccountsModel) SubmitDefault() bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListAccountsModel687eda_SubmitDefault(ptr.Pointer())) != 0
	}
	return false
}

//export callbackTxListAccountsModel687eda_SupportedDragActions
func callbackTxListAccountsModel687eda_SupportedDragActions(ptr unsafe.Pointer) C.longlong {
	if signal := qt.GetSignal(ptr, "supportedDragActions"); signal != nil {
		return C.longlong((*(*func() std_core.Qt__DropAction)(signal))())
	}

	return C.longlong(NewTxListAccountsModelFromPointer(ptr).SupportedDragActionsDefault())
}

func (ptr *TxListAccountsModel) SupportedDragActionsDefault() std_core.Qt__DropAction {
	if ptr.Pointer() != nil {
		return std_core.Qt__DropAction(C.TxListAccountsModel687eda_SupportedDragActionsDefault(ptr.Pointer()))
	}
	return 0
}

//export callbackTxListAccountsModel687eda_SupportedDropActions
func callbackTxListAccountsModel687eda_SupportedDropActions(ptr unsafe.Pointer) C.longlong {
	if signal := qt.GetSignal(ptr, "supportedDropActions"); signal != nil {
		return C.longlong((*(*func() std_core.Qt__DropAction)(signal))())
	}

	return C.longlong(NewTxListAccountsModelFromPointer(ptr).SupportedDropActionsDefault())
}

func (ptr *TxListAccountsModel) SupportedDropActionsDefault() std_core.Qt__DropAction {
	if ptr.Pointer() != nil {
		return std_core.Qt__DropAction(C.TxListAccountsModel687eda_SupportedDropActionsDefault(ptr.Pointer()))
	}
	return 0
}

//export callbackTxListAccountsModel687eda_ChildEvent
func callbackTxListAccountsModel687eda_ChildEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "childEvent"); signal != nil {
		(*(*func(*std_core.QChildEvent))(signal))(std_core.NewQChildEventFromPointer(event))
	} else {
		NewTxListAccountsModelFromPointer(ptr).ChildEventDefault(std_core.NewQChildEventFromPointer(event))
	}
}

func (ptr *TxListAccountsModel) ChildEventDefault(event std_core.QChildEvent_ITF) {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda_ChildEventDefault(ptr.Pointer(), std_core.PointerFromQChildEvent(event))
	}
}

//export callbackTxListAccountsModel687eda_ConnectNotify
func callbackTxListAccountsModel687eda_ConnectNotify(ptr unsafe.Pointer, sign unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "connectNotify"); signal != nil {
		(*(*func(*std_core.QMetaMethod))(signal))(std_core.NewQMetaMethodFromPointer(sign))
	} else {
		NewTxListAccountsModelFromPointer(ptr).ConnectNotifyDefault(std_core.NewQMetaMethodFromPointer(sign))
	}
}

func (ptr *TxListAccountsModel) ConnectNotifyDefault(sign std_core.QMetaMethod_ITF) {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda_ConnectNotifyDefault(ptr.Pointer(), std_core.PointerFromQMetaMethod(sign))
	}
}

//export callbackTxListAccountsModel687eda_CustomEvent
func callbackTxListAccountsModel687eda_CustomEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "customEvent"); signal != nil {
		(*(*func(*std_core.QEvent))(signal))(std_core.NewQEventFromPointer(event))
	} else {
		NewTxListAccountsModelFromPointer(ptr).CustomEventDefault(std_core.NewQEventFromPointer(event))
	}
}

func (ptr *TxListAccountsModel) CustomEventDefault(event std_core.QEvent_ITF) {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda_CustomEventDefault(ptr.Pointer(), std_core.PointerFromQEvent(event))
	}
}

//export callbackTxListAccountsModel687eda_DeleteLater
func callbackTxListAccountsModel687eda_DeleteLater(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "deleteLater"); signal != nil {
		(*(*func())(signal))()
	} else {
		NewTxListAccountsModelFromPointer(ptr).DeleteLaterDefault()
	}
}

func (ptr *TxListAccountsModel) DeleteLaterDefault() {
	if ptr.Pointer() != nil {

		qt.SetFinalizer(ptr, nil)
		C.TxListAccountsModel687eda_DeleteLaterDefault(ptr.Pointer())
	}
}

//export callbackTxListAccountsModel687eda_Destroyed
func callbackTxListAccountsModel687eda_Destroyed(ptr unsafe.Pointer, obj unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "destroyed"); signal != nil {
		(*(*func(*std_core.QObject))(signal))(std_core.NewQObjectFromPointer(obj))
	}
	qt.Unregister(ptr)

}

//export callbackTxListAccountsModel687eda_DisconnectNotify
func callbackTxListAccountsModel687eda_DisconnectNotify(ptr unsafe.Pointer, sign unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "disconnectNotify"); signal != nil {
		(*(*func(*std_core.QMetaMethod))(signal))(std_core.NewQMetaMethodFromPointer(sign))
	} else {
		NewTxListAccountsModelFromPointer(ptr).DisconnectNotifyDefault(std_core.NewQMetaMethodFromPointer(sign))
	}
}

func (ptr *TxListAccountsModel) DisconnectNotifyDefault(sign std_core.QMetaMethod_ITF) {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda_DisconnectNotifyDefault(ptr.Pointer(), std_core.PointerFromQMetaMethod(sign))
	}
}

//export callbackTxListAccountsModel687eda_Event
func callbackTxListAccountsModel687eda_Event(ptr unsafe.Pointer, e unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "event"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QEvent) bool)(signal))(std_core.NewQEventFromPointer(e)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListAccountsModelFromPointer(ptr).EventDefault(std_core.NewQEventFromPointer(e)))))
}

func (ptr *TxListAccountsModel) EventDefault(e std_core.QEvent_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListAccountsModel687eda_EventDefault(ptr.Pointer(), std_core.PointerFromQEvent(e))) != 0
	}
	return false
}

//export callbackTxListAccountsModel687eda_EventFilter
func callbackTxListAccountsModel687eda_EventFilter(ptr unsafe.Pointer, watched unsafe.Pointer, event unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "eventFilter"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QObject, *std_core.QEvent) bool)(signal))(std_core.NewQObjectFromPointer(watched), std_core.NewQEventFromPointer(event)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListAccountsModelFromPointer(ptr).EventFilterDefault(std_core.NewQObjectFromPointer(watched), std_core.NewQEventFromPointer(event)))))
}

func (ptr *TxListAccountsModel) EventFilterDefault(watched std_core.QObject_ITF, event std_core.QEvent_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListAccountsModel687eda_EventFilterDefault(ptr.Pointer(), std_core.PointerFromQObject(watched), std_core.PointerFromQEvent(event))) != 0
	}
	return false
}

//export callbackTxListAccountsModel687eda_ObjectNameChanged
func callbackTxListAccountsModel687eda_ObjectNameChanged(ptr unsafe.Pointer, objectName C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "objectNameChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(objectName))
	}

}

//export callbackTxListAccountsModel687eda_TimerEvent
func callbackTxListAccountsModel687eda_TimerEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "timerEvent"); signal != nil {
		(*(*func(*std_core.QTimerEvent))(signal))(std_core.NewQTimerEventFromPointer(event))
	} else {
		NewTxListAccountsModelFromPointer(ptr).TimerEventDefault(std_core.NewQTimerEventFromPointer(event))
	}
}

func (ptr *TxListAccountsModel) TimerEventDefault(event std_core.QTimerEvent_ITF) {
	if ptr.Pointer() != nil {
		C.TxListAccountsModel687eda_TimerEventDefault(ptr.Pointer(), std_core.PointerFromQTimerEvent(event))
	}
}

type TxListCtx_ITF interface {
	std_core.QObject_ITF
	TxListCtx_PTR() *TxListCtx
}

func (ptr *TxListCtx) TxListCtx_PTR() *TxListCtx {
	return ptr
}

func (ptr *TxListCtx) Pointer() unsafe.Pointer {
	if ptr != nil {
		return ptr.QObject_PTR().Pointer()
	}
	return nil
}

func (ptr *TxListCtx) SetPointer(p unsafe.Pointer) {
	if ptr != nil {
		ptr.QObject_PTR().SetPointer(p)
	}
}

func PointerFromTxListCtx(ptr TxListCtx_ITF) unsafe.Pointer {
	if ptr != nil {
		return ptr.TxListCtx_PTR().Pointer()
	}
	return nil
}

func NewTxListCtxFromPointer(ptr unsafe.Pointer) (n *TxListCtx) {
	if gPtr, ok := qt.Receive(ptr); !ok {
		n = new(TxListCtx)
		n.SetPointer(ptr)
	} else {
		switch deduced := gPtr.(type) {
		case *TxListCtx:
			n = deduced

		case *std_core.QObject:
			n = &TxListCtx{QObject: *deduced}

		default:
			n = new(TxListCtx)
			n.SetPointer(ptr)
		}
	}
	return
}
func (this *TxListCtx) Init() { this.init() }

//export callbackTxListCtx687eda_Constructor
func callbackTxListCtx687eda_Constructor(ptr unsafe.Pointer) {
	this := NewTxListCtxFromPointer(ptr)
	qt.Register(ptr, this)
	this.ConnectClicked(this.clicked)
	this.init()
}

//export callbackTxListCtx687eda_Clicked
func callbackTxListCtx687eda_Clicked(ptr unsafe.Pointer, b C.int) {
	if signal := qt.GetSignal(ptr, "clicked"); signal != nil {
		(*(*func(int))(signal))(int(int32(b)))
	}

}

func (ptr *TxListCtx) ConnectClicked(f func(b int)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "clicked") {
			C.TxListCtx687eda_ConnectClicked(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "clicked")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "clicked"); signal != nil {
			f := func(b int) {
				(*(*func(int))(signal))(b)
				f(b)
			}
			qt.ConnectSignal(ptr.Pointer(), "clicked", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "clicked", unsafe.Pointer(&f))
		}
	}
}

func (ptr *TxListCtx) DisconnectClicked() {
	if ptr.Pointer() != nil {
		C.TxListCtx687eda_DisconnectClicked(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "clicked")
	}
}

func (ptr *TxListCtx) Clicked(b int) {
	if ptr.Pointer() != nil {
		C.TxListCtx687eda_Clicked(ptr.Pointer(), C.int(int32(b)))
	}
}

//export callbackTxListCtx687eda_ShortenAddress
func callbackTxListCtx687eda_ShortenAddress(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "shortenAddress"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewTxListCtxFromPointer(ptr).ShortenAddressDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *TxListCtx) ConnectShortenAddress(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "shortenAddress"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "shortenAddress", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "shortenAddress", unsafe.Pointer(&f))
		}
	}
}

func (ptr *TxListCtx) DisconnectShortenAddress() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "shortenAddress")
	}
}

func (ptr *TxListCtx) ShortenAddress() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.TxListCtx687eda_ShortenAddress(ptr.Pointer()))
	}
	return ""
}

func (ptr *TxListCtx) ShortenAddressDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.TxListCtx687eda_ShortenAddressDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackTxListCtx687eda_SetShortenAddress
func callbackTxListCtx687eda_SetShortenAddress(ptr unsafe.Pointer, shortenAddress C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setShortenAddress"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(shortenAddress))
	} else {
		NewTxListCtxFromPointer(ptr).SetShortenAddressDefault(cGoUnpackString(shortenAddress))
	}
}

func (ptr *TxListCtx) ConnectSetShortenAddress(f func(shortenAddress string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setShortenAddress"); signal != nil {
			f := func(shortenAddress string) {
				(*(*func(string))(signal))(shortenAddress)
				f(shortenAddress)
			}
			qt.ConnectSignal(ptr.Pointer(), "setShortenAddress", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setShortenAddress", unsafe.Pointer(&f))
		}
	}
}

func (ptr *TxListCtx) DisconnectSetShortenAddress() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setShortenAddress")
	}
}

func (ptr *TxListCtx) SetShortenAddress(shortenAddress string) {
	if ptr.Pointer() != nil {
		var shortenAddressC *C.char
		if shortenAddress != "" {
			shortenAddressC = C.CString(shortenAddress)
			defer C.free(unsafe.Pointer(shortenAddressC))
		}
		C.TxListCtx687eda_SetShortenAddress(ptr.Pointer(), C.struct_Moc_PackedString{data: shortenAddressC, len: C.longlong(len(shortenAddress))})
	}
}

func (ptr *TxListCtx) SetShortenAddressDefault(shortenAddress string) {
	if ptr.Pointer() != nil {
		var shortenAddressC *C.char
		if shortenAddress != "" {
			shortenAddressC = C.CString(shortenAddress)
			defer C.free(unsafe.Pointer(shortenAddressC))
		}
		C.TxListCtx687eda_SetShortenAddressDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: shortenAddressC, len: C.longlong(len(shortenAddress))})
	}
}

//export callbackTxListCtx687eda_ShortenAddressChanged
func callbackTxListCtx687eda_ShortenAddressChanged(ptr unsafe.Pointer, shortenAddress C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "shortenAddressChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(shortenAddress))
	}

}

func (ptr *TxListCtx) ConnectShortenAddressChanged(f func(shortenAddress string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "shortenAddressChanged") {
			C.TxListCtx687eda_ConnectShortenAddressChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "shortenAddressChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "shortenAddressChanged"); signal != nil {
			f := func(shortenAddress string) {
				(*(*func(string))(signal))(shortenAddress)
				f(shortenAddress)
			}
			qt.ConnectSignal(ptr.Pointer(), "shortenAddressChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "shortenAddressChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *TxListCtx) DisconnectShortenAddressChanged() {
	if ptr.Pointer() != nil {
		C.TxListCtx687eda_DisconnectShortenAddressChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "shortenAddressChanged")
	}
}

func (ptr *TxListCtx) ShortenAddressChanged(shortenAddress string) {
	if ptr.Pointer() != nil {
		var shortenAddressC *C.char
		if shortenAddress != "" {
			shortenAddressC = C.CString(shortenAddress)
			defer C.free(unsafe.Pointer(shortenAddressC))
		}
		C.TxListCtx687eda_ShortenAddressChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: shortenAddressC, len: C.longlong(len(shortenAddress))})
	}
}

//export callbackTxListCtx687eda_SelectedSrc
func callbackTxListCtx687eda_SelectedSrc(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "selectedSrc"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewTxListCtxFromPointer(ptr).SelectedSrcDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *TxListCtx) ConnectSelectedSrc(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "selectedSrc"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "selectedSrc", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "selectedSrc", unsafe.Pointer(&f))
		}
	}
}

func (ptr *TxListCtx) DisconnectSelectedSrc() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "selectedSrc")
	}
}

func (ptr *TxListCtx) SelectedSrc() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.TxListCtx687eda_SelectedSrc(ptr.Pointer()))
	}
	return ""
}

func (ptr *TxListCtx) SelectedSrcDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.TxListCtx687eda_SelectedSrcDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackTxListCtx687eda_SetSelectedSrc
func callbackTxListCtx687eda_SetSelectedSrc(ptr unsafe.Pointer, selectedSrc C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setSelectedSrc"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(selectedSrc))
	} else {
		NewTxListCtxFromPointer(ptr).SetSelectedSrcDefault(cGoUnpackString(selectedSrc))
	}
}

func (ptr *TxListCtx) ConnectSetSelectedSrc(f func(selectedSrc string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setSelectedSrc"); signal != nil {
			f := func(selectedSrc string) {
				(*(*func(string))(signal))(selectedSrc)
				f(selectedSrc)
			}
			qt.ConnectSignal(ptr.Pointer(), "setSelectedSrc", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setSelectedSrc", unsafe.Pointer(&f))
		}
	}
}

func (ptr *TxListCtx) DisconnectSetSelectedSrc() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setSelectedSrc")
	}
}

func (ptr *TxListCtx) SetSelectedSrc(selectedSrc string) {
	if ptr.Pointer() != nil {
		var selectedSrcC *C.char
		if selectedSrc != "" {
			selectedSrcC = C.CString(selectedSrc)
			defer C.free(unsafe.Pointer(selectedSrcC))
		}
		C.TxListCtx687eda_SetSelectedSrc(ptr.Pointer(), C.struct_Moc_PackedString{data: selectedSrcC, len: C.longlong(len(selectedSrc))})
	}
}

func (ptr *TxListCtx) SetSelectedSrcDefault(selectedSrc string) {
	if ptr.Pointer() != nil {
		var selectedSrcC *C.char
		if selectedSrc != "" {
			selectedSrcC = C.CString(selectedSrc)
			defer C.free(unsafe.Pointer(selectedSrcC))
		}
		C.TxListCtx687eda_SetSelectedSrcDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: selectedSrcC, len: C.longlong(len(selectedSrc))})
	}
}

//export callbackTxListCtx687eda_SelectedSrcChanged
func callbackTxListCtx687eda_SelectedSrcChanged(ptr unsafe.Pointer, selectedSrc C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "selectedSrcChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(selectedSrc))
	}

}

func (ptr *TxListCtx) ConnectSelectedSrcChanged(f func(selectedSrc string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "selectedSrcChanged") {
			C.TxListCtx687eda_ConnectSelectedSrcChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "selectedSrcChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "selectedSrcChanged"); signal != nil {
			f := func(selectedSrc string) {
				(*(*func(string))(signal))(selectedSrc)
				f(selectedSrc)
			}
			qt.ConnectSignal(ptr.Pointer(), "selectedSrcChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "selectedSrcChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *TxListCtx) DisconnectSelectedSrcChanged() {
	if ptr.Pointer() != nil {
		C.TxListCtx687eda_DisconnectSelectedSrcChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "selectedSrcChanged")
	}
}

func (ptr *TxListCtx) SelectedSrcChanged(selectedSrc string) {
	if ptr.Pointer() != nil {
		var selectedSrcC *C.char
		if selectedSrc != "" {
			selectedSrcC = C.CString(selectedSrc)
			defer C.free(unsafe.Pointer(selectedSrcC))
		}
		C.TxListCtx687eda_SelectedSrcChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: selectedSrcC, len: C.longlong(len(selectedSrc))})
	}
}

func TxListCtx_QRegisterMetaType() int {
	return int(int32(C.TxListCtx687eda_TxListCtx687eda_QRegisterMetaType()))
}

func (ptr *TxListCtx) QRegisterMetaType() int {
	return int(int32(C.TxListCtx687eda_TxListCtx687eda_QRegisterMetaType()))
}

func TxListCtx_QRegisterMetaType2(typeName string) int {
	var typeNameC *C.char
	if typeName != "" {
		typeNameC = C.CString(typeName)
		defer C.free(unsafe.Pointer(typeNameC))
	}
	return int(int32(C.TxListCtx687eda_TxListCtx687eda_QRegisterMetaType2(typeNameC)))
}

func (ptr *TxListCtx) QRegisterMetaType2(typeName string) int {
	var typeNameC *C.char
	if typeName != "" {
		typeNameC = C.CString(typeName)
		defer C.free(unsafe.Pointer(typeNameC))
	}
	return int(int32(C.TxListCtx687eda_TxListCtx687eda_QRegisterMetaType2(typeNameC)))
}

func TxListCtx_QmlRegisterType() int {
	return int(int32(C.TxListCtx687eda_TxListCtx687eda_QmlRegisterType()))
}

func (ptr *TxListCtx) QmlRegisterType() int {
	return int(int32(C.TxListCtx687eda_TxListCtx687eda_QmlRegisterType()))
}

func TxListCtx_QmlRegisterType2(uri string, versionMajor int, versionMinor int, qmlName string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	return int(int32(C.TxListCtx687eda_TxListCtx687eda_QmlRegisterType2(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC)))
}

func (ptr *TxListCtx) QmlRegisterType2(uri string, versionMajor int, versionMinor int, qmlName string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	return int(int32(C.TxListCtx687eda_TxListCtx687eda_QmlRegisterType2(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC)))
}

func TxListCtx_QmlRegisterUncreatableType(uri string, versionMajor int, versionMinor int, qmlName string, message string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	var messageC *C.char
	if message != "" {
		messageC = C.CString(message)
		defer C.free(unsafe.Pointer(messageC))
	}
	return int(int32(C.TxListCtx687eda_TxListCtx687eda_QmlRegisterUncreatableType(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC, C.struct_Moc_PackedString{data: messageC, len: C.longlong(len(message))})))
}

func (ptr *TxListCtx) QmlRegisterUncreatableType(uri string, versionMajor int, versionMinor int, qmlName string, message string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	var messageC *C.char
	if message != "" {
		messageC = C.CString(message)
		defer C.free(unsafe.Pointer(messageC))
	}
	return int(int32(C.TxListCtx687eda_TxListCtx687eda_QmlRegisterUncreatableType(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC, C.struct_Moc_PackedString{data: messageC, len: C.longlong(len(message))})))
}

func (ptr *TxListCtx) __children_atList(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQObjectFromPointer(C.TxListCtx687eda___children_atList(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *TxListCtx) __children_setList(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.TxListCtx687eda___children_setList(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *TxListCtx) __children_newList() unsafe.Pointer {
	return C.TxListCtx687eda___children_newList(ptr.Pointer())
}

func (ptr *TxListCtx) __dynamicPropertyNames_atList(i int) *std_core.QByteArray {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQByteArrayFromPointer(C.TxListCtx687eda___dynamicPropertyNames_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QByteArray).DestroyQByteArray)
		return tmpValue
	}
	return nil
}

func (ptr *TxListCtx) __dynamicPropertyNames_setList(i std_core.QByteArray_ITF) {
	if ptr.Pointer() != nil {
		C.TxListCtx687eda___dynamicPropertyNames_setList(ptr.Pointer(), std_core.PointerFromQByteArray(i))
	}
}

func (ptr *TxListCtx) __dynamicPropertyNames_newList() unsafe.Pointer {
	return C.TxListCtx687eda___dynamicPropertyNames_newList(ptr.Pointer())
}

func (ptr *TxListCtx) __findChildren_atList(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQObjectFromPointer(C.TxListCtx687eda___findChildren_atList(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *TxListCtx) __findChildren_setList(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.TxListCtx687eda___findChildren_setList(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *TxListCtx) __findChildren_newList() unsafe.Pointer {
	return C.TxListCtx687eda___findChildren_newList(ptr.Pointer())
}

func (ptr *TxListCtx) __findChildren_atList3(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQObjectFromPointer(C.TxListCtx687eda___findChildren_atList3(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *TxListCtx) __findChildren_setList3(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.TxListCtx687eda___findChildren_setList3(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *TxListCtx) __findChildren_newList3() unsafe.Pointer {
	return C.TxListCtx687eda___findChildren_newList3(ptr.Pointer())
}

func NewTxListCtx(parent std_core.QObject_ITF) *TxListCtx {
	TxListCtx_QRegisterMetaType()
	tmpValue := NewTxListCtxFromPointer(C.TxListCtx687eda_NewTxListCtx(std_core.PointerFromQObject(parent)))
	if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
		tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
	}
	return tmpValue
}

//export callbackTxListCtx687eda_DestroyTxListCtx
func callbackTxListCtx687eda_DestroyTxListCtx(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "~TxListCtx"); signal != nil {
		(*(*func())(signal))()
	} else {
		NewTxListCtxFromPointer(ptr).DestroyTxListCtxDefault()
	}
}

func (ptr *TxListCtx) ConnectDestroyTxListCtx(f func()) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "~TxListCtx"); signal != nil {
			f := func() {
				(*(*func())(signal))()
				f()
			}
			qt.ConnectSignal(ptr.Pointer(), "~TxListCtx", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "~TxListCtx", unsafe.Pointer(&f))
		}
	}
}

func (ptr *TxListCtx) DisconnectDestroyTxListCtx() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "~TxListCtx")
	}
}

func (ptr *TxListCtx) DestroyTxListCtx() {
	if ptr.Pointer() != nil {

		qt.SetFinalizer(ptr, nil)
		C.TxListCtx687eda_DestroyTxListCtx(ptr.Pointer())
		ptr.SetPointer(nil)
	}
}

func (ptr *TxListCtx) DestroyTxListCtxDefault() {
	if ptr.Pointer() != nil {

		qt.SetFinalizer(ptr, nil)
		C.TxListCtx687eda_DestroyTxListCtxDefault(ptr.Pointer())
		ptr.SetPointer(nil)
	}
}

//export callbackTxListCtx687eda_ChildEvent
func callbackTxListCtx687eda_ChildEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "childEvent"); signal != nil {
		(*(*func(*std_core.QChildEvent))(signal))(std_core.NewQChildEventFromPointer(event))
	} else {
		NewTxListCtxFromPointer(ptr).ChildEventDefault(std_core.NewQChildEventFromPointer(event))
	}
}

func (ptr *TxListCtx) ChildEventDefault(event std_core.QChildEvent_ITF) {
	if ptr.Pointer() != nil {
		C.TxListCtx687eda_ChildEventDefault(ptr.Pointer(), std_core.PointerFromQChildEvent(event))
	}
}

//export callbackTxListCtx687eda_ConnectNotify
func callbackTxListCtx687eda_ConnectNotify(ptr unsafe.Pointer, sign unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "connectNotify"); signal != nil {
		(*(*func(*std_core.QMetaMethod))(signal))(std_core.NewQMetaMethodFromPointer(sign))
	} else {
		NewTxListCtxFromPointer(ptr).ConnectNotifyDefault(std_core.NewQMetaMethodFromPointer(sign))
	}
}

func (ptr *TxListCtx) ConnectNotifyDefault(sign std_core.QMetaMethod_ITF) {
	if ptr.Pointer() != nil {
		C.TxListCtx687eda_ConnectNotifyDefault(ptr.Pointer(), std_core.PointerFromQMetaMethod(sign))
	}
}

//export callbackTxListCtx687eda_CustomEvent
func callbackTxListCtx687eda_CustomEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "customEvent"); signal != nil {
		(*(*func(*std_core.QEvent))(signal))(std_core.NewQEventFromPointer(event))
	} else {
		NewTxListCtxFromPointer(ptr).CustomEventDefault(std_core.NewQEventFromPointer(event))
	}
}

func (ptr *TxListCtx) CustomEventDefault(event std_core.QEvent_ITF) {
	if ptr.Pointer() != nil {
		C.TxListCtx687eda_CustomEventDefault(ptr.Pointer(), std_core.PointerFromQEvent(event))
	}
}

//export callbackTxListCtx687eda_DeleteLater
func callbackTxListCtx687eda_DeleteLater(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "deleteLater"); signal != nil {
		(*(*func())(signal))()
	} else {
		NewTxListCtxFromPointer(ptr).DeleteLaterDefault()
	}
}

func (ptr *TxListCtx) DeleteLaterDefault() {
	if ptr.Pointer() != nil {

		qt.SetFinalizer(ptr, nil)
		C.TxListCtx687eda_DeleteLaterDefault(ptr.Pointer())
	}
}

//export callbackTxListCtx687eda_Destroyed
func callbackTxListCtx687eda_Destroyed(ptr unsafe.Pointer, obj unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "destroyed"); signal != nil {
		(*(*func(*std_core.QObject))(signal))(std_core.NewQObjectFromPointer(obj))
	}
	qt.Unregister(ptr)

}

//export callbackTxListCtx687eda_DisconnectNotify
func callbackTxListCtx687eda_DisconnectNotify(ptr unsafe.Pointer, sign unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "disconnectNotify"); signal != nil {
		(*(*func(*std_core.QMetaMethod))(signal))(std_core.NewQMetaMethodFromPointer(sign))
	} else {
		NewTxListCtxFromPointer(ptr).DisconnectNotifyDefault(std_core.NewQMetaMethodFromPointer(sign))
	}
}

func (ptr *TxListCtx) DisconnectNotifyDefault(sign std_core.QMetaMethod_ITF) {
	if ptr.Pointer() != nil {
		C.TxListCtx687eda_DisconnectNotifyDefault(ptr.Pointer(), std_core.PointerFromQMetaMethod(sign))
	}
}

//export callbackTxListCtx687eda_Event
func callbackTxListCtx687eda_Event(ptr unsafe.Pointer, e unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "event"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QEvent) bool)(signal))(std_core.NewQEventFromPointer(e)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListCtxFromPointer(ptr).EventDefault(std_core.NewQEventFromPointer(e)))))
}

func (ptr *TxListCtx) EventDefault(e std_core.QEvent_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListCtx687eda_EventDefault(ptr.Pointer(), std_core.PointerFromQEvent(e))) != 0
	}
	return false
}

//export callbackTxListCtx687eda_EventFilter
func callbackTxListCtx687eda_EventFilter(ptr unsafe.Pointer, watched unsafe.Pointer, event unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "eventFilter"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QObject, *std_core.QEvent) bool)(signal))(std_core.NewQObjectFromPointer(watched), std_core.NewQEventFromPointer(event)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListCtxFromPointer(ptr).EventFilterDefault(std_core.NewQObjectFromPointer(watched), std_core.NewQEventFromPointer(event)))))
}

func (ptr *TxListCtx) EventFilterDefault(watched std_core.QObject_ITF, event std_core.QEvent_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListCtx687eda_EventFilterDefault(ptr.Pointer(), std_core.PointerFromQObject(watched), std_core.PointerFromQEvent(event))) != 0
	}
	return false
}

//export callbackTxListCtx687eda_ObjectNameChanged
func callbackTxListCtx687eda_ObjectNameChanged(ptr unsafe.Pointer, objectName C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "objectNameChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(objectName))
	}

}

//export callbackTxListCtx687eda_TimerEvent
func callbackTxListCtx687eda_TimerEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "timerEvent"); signal != nil {
		(*(*func(*std_core.QTimerEvent))(signal))(std_core.NewQTimerEventFromPointer(event))
	} else {
		NewTxListCtxFromPointer(ptr).TimerEventDefault(std_core.NewQTimerEventFromPointer(event))
	}
}

func (ptr *TxListCtx) TimerEventDefault(event std_core.QTimerEvent_ITF) {
	if ptr.Pointer() != nil {
		C.TxListCtx687eda_TimerEventDefault(ptr.Pointer(), std_core.PointerFromQTimerEvent(event))
	}
}

type TxListModel_ITF interface {
	std_core.QAbstractListModel_ITF
	TxListModel_PTR() *TxListModel
}

func (ptr *TxListModel) TxListModel_PTR() *TxListModel {
	return ptr
}

func (ptr *TxListModel) Pointer() unsafe.Pointer {
	if ptr != nil {
		return ptr.QAbstractListModel_PTR().Pointer()
	}
	return nil
}

func (ptr *TxListModel) SetPointer(p unsafe.Pointer) {
	if ptr != nil {
		ptr.QAbstractListModel_PTR().SetPointer(p)
	}
}

func PointerFromTxListModel(ptr TxListModel_ITF) unsafe.Pointer {
	if ptr != nil {
		return ptr.TxListModel_PTR().Pointer()
	}
	return nil
}

func NewTxListModelFromPointer(ptr unsafe.Pointer) (n *TxListModel) {
	if gPtr, ok := qt.Receive(ptr); !ok {
		n = new(TxListModel)
		n.SetPointer(ptr)
	} else {
		switch deduced := gPtr.(type) {
		case *TxListModel:
			n = deduced

		case *std_core.QAbstractListModel:
			n = &TxListModel{QAbstractListModel: *deduced}

		default:
			n = new(TxListModel)
			n.SetPointer(ptr)
		}
	}
	return
}
func (this *TxListModel) Init() { this.init() }

//export callbackTxListModel687eda_Constructor
func callbackTxListModel687eda_Constructor(ptr unsafe.Pointer) {
	this := NewTxListModelFromPointer(ptr)
	qt.Register(ptr, this)
	this.ConnectClear(this.clear)
	this.ConnectAdd(this.add)
	this.ConnectRemove(this.remove)
	this.init()
}

//export callbackTxListModel687eda_Clear
func callbackTxListModel687eda_Clear(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "clear"); signal != nil {
		(*(*func())(signal))()
	}

}

func (ptr *TxListModel) ConnectClear(f func()) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "clear") {
			C.TxListModel687eda_ConnectClear(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "clear")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "clear"); signal != nil {
			f := func() {
				(*(*func())(signal))()
				f()
			}
			qt.ConnectSignal(ptr.Pointer(), "clear", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "clear", unsafe.Pointer(&f))
		}
	}
}

func (ptr *TxListModel) DisconnectClear() {
	if ptr.Pointer() != nil {
		C.TxListModel687eda_DisconnectClear(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "clear")
	}
}

func (ptr *TxListModel) Clear() {
	if ptr.Pointer() != nil {
		C.TxListModel687eda_Clear(ptr.Pointer())
	}
}

//export callbackTxListModel687eda_Add
func callbackTxListModel687eda_Add(ptr unsafe.Pointer, tx C.uintptr_t) {
	var txD IncomingRequestItem
	if txI, ok := qt.ReceiveTemp(unsafe.Pointer(uintptr(tx))); ok {
		qt.UnregisterTemp(unsafe.Pointer(uintptr(tx)))
		txD = (*(*IncomingRequestItem)(txI))
	}
	if signal := qt.GetSignal(ptr, "add"); signal != nil {
		(*(*func(IncomingRequestItem))(signal))(txD)
	}

}

func (ptr *TxListModel) ConnectAdd(f func(tx IncomingRequestItem)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "add") {
			C.TxListModel687eda_ConnectAdd(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "add")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "add"); signal != nil {
			f := func(tx IncomingRequestItem) {
				(*(*func(IncomingRequestItem))(signal))(tx)
				f(tx)
			}
			qt.ConnectSignal(ptr.Pointer(), "add", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "add", unsafe.Pointer(&f))
		}
	}
}

func (ptr *TxListModel) DisconnectAdd() {
	if ptr.Pointer() != nil {
		C.TxListModel687eda_DisconnectAdd(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "add")
	}
}

func (ptr *TxListModel) Add(tx IncomingRequestItem) {
	if ptr.Pointer() != nil {
		txTID := time.Now().UnixNano() + int64(uintptr(unsafe.Pointer(&tx)))
		qt.RegisterTemp(unsafe.Pointer(uintptr(txTID)), unsafe.Pointer(&tx))
		C.TxListModel687eda_Add(ptr.Pointer(), C.uintptr_t(txTID))
	}
}

//export callbackTxListModel687eda_Remove
func callbackTxListModel687eda_Remove(ptr unsafe.Pointer, i C.int) {
	if signal := qt.GetSignal(ptr, "remove"); signal != nil {
		(*(*func(int))(signal))(int(int32(i)))
	}

}

func (ptr *TxListModel) ConnectRemove(f func(i int)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "remove") {
			C.TxListModel687eda_ConnectRemove(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "remove")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "remove"); signal != nil {
			f := func(i int) {
				(*(*func(int))(signal))(i)
				f(i)
			}
			qt.ConnectSignal(ptr.Pointer(), "remove", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "remove", unsafe.Pointer(&f))
		}
	}
}

func (ptr *TxListModel) DisconnectRemove() {
	if ptr.Pointer() != nil {
		C.TxListModel687eda_DisconnectRemove(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "remove")
	}
}

func (ptr *TxListModel) Remove(i int) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda_Remove(ptr.Pointer(), C.int(int32(i)))
	}
}

//export callbackTxListModel687eda_IsEmpty
func callbackTxListModel687eda_IsEmpty(ptr unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "isEmpty"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func() bool)(signal))())))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListModelFromPointer(ptr).IsEmptyDefault())))
}

func (ptr *TxListModel) ConnectIsEmpty(f func() bool) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "isEmpty"); signal != nil {
			f := func() bool {
				(*(*func() bool)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "isEmpty", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "isEmpty", unsafe.Pointer(&f))
		}
	}
}

func (ptr *TxListModel) DisconnectIsEmpty() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "isEmpty")
	}
}

func (ptr *TxListModel) IsEmpty() bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListModel687eda_IsEmpty(ptr.Pointer())) != 0
	}
	return false
}

func (ptr *TxListModel) IsEmptyDefault() bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListModel687eda_IsEmptyDefault(ptr.Pointer())) != 0
	}
	return false
}

//export callbackTxListModel687eda_SetIsEmpty
func callbackTxListModel687eda_SetIsEmpty(ptr unsafe.Pointer, isEmpty C.char) {
	if signal := qt.GetSignal(ptr, "setIsEmpty"); signal != nil {
		(*(*func(bool))(signal))(int8(isEmpty) != 0)
	} else {
		NewTxListModelFromPointer(ptr).SetIsEmptyDefault(int8(isEmpty) != 0)
	}
}

func (ptr *TxListModel) ConnectSetIsEmpty(f func(isEmpty bool)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setIsEmpty"); signal != nil {
			f := func(isEmpty bool) {
				(*(*func(bool))(signal))(isEmpty)
				f(isEmpty)
			}
			qt.ConnectSignal(ptr.Pointer(), "setIsEmpty", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setIsEmpty", unsafe.Pointer(&f))
		}
	}
}

func (ptr *TxListModel) DisconnectSetIsEmpty() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setIsEmpty")
	}
}

func (ptr *TxListModel) SetIsEmpty(isEmpty bool) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda_SetIsEmpty(ptr.Pointer(), C.char(int8(qt.GoBoolToInt(isEmpty))))
	}
}

func (ptr *TxListModel) SetIsEmptyDefault(isEmpty bool) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda_SetIsEmptyDefault(ptr.Pointer(), C.char(int8(qt.GoBoolToInt(isEmpty))))
	}
}

//export callbackTxListModel687eda_IsEmptyChanged
func callbackTxListModel687eda_IsEmptyChanged(ptr unsafe.Pointer, isEmpty C.char) {
	if signal := qt.GetSignal(ptr, "isEmptyChanged"); signal != nil {
		(*(*func(bool))(signal))(int8(isEmpty) != 0)
	}

}

func (ptr *TxListModel) ConnectIsEmptyChanged(f func(isEmpty bool)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "isEmptyChanged") {
			C.TxListModel687eda_ConnectIsEmptyChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "isEmptyChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "isEmptyChanged"); signal != nil {
			f := func(isEmpty bool) {
				(*(*func(bool))(signal))(isEmpty)
				f(isEmpty)
			}
			qt.ConnectSignal(ptr.Pointer(), "isEmptyChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "isEmptyChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *TxListModel) DisconnectIsEmptyChanged() {
	if ptr.Pointer() != nil {
		C.TxListModel687eda_DisconnectIsEmptyChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "isEmptyChanged")
	}
}

func (ptr *TxListModel) IsEmptyChanged(isEmpty bool) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda_IsEmptyChanged(ptr.Pointer(), C.char(int8(qt.GoBoolToInt(isEmpty))))
	}
}

func TxListModel_QRegisterMetaType() int {
	return int(int32(C.TxListModel687eda_TxListModel687eda_QRegisterMetaType()))
}

func (ptr *TxListModel) QRegisterMetaType() int {
	return int(int32(C.TxListModel687eda_TxListModel687eda_QRegisterMetaType()))
}

func TxListModel_QRegisterMetaType2(typeName string) int {
	var typeNameC *C.char
	if typeName != "" {
		typeNameC = C.CString(typeName)
		defer C.free(unsafe.Pointer(typeNameC))
	}
	return int(int32(C.TxListModel687eda_TxListModel687eda_QRegisterMetaType2(typeNameC)))
}

func (ptr *TxListModel) QRegisterMetaType2(typeName string) int {
	var typeNameC *C.char
	if typeName != "" {
		typeNameC = C.CString(typeName)
		defer C.free(unsafe.Pointer(typeNameC))
	}
	return int(int32(C.TxListModel687eda_TxListModel687eda_QRegisterMetaType2(typeNameC)))
}

func TxListModel_QmlRegisterType() int {
	return int(int32(C.TxListModel687eda_TxListModel687eda_QmlRegisterType()))
}

func (ptr *TxListModel) QmlRegisterType() int {
	return int(int32(C.TxListModel687eda_TxListModel687eda_QmlRegisterType()))
}

func TxListModel_QmlRegisterType2(uri string, versionMajor int, versionMinor int, qmlName string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	return int(int32(C.TxListModel687eda_TxListModel687eda_QmlRegisterType2(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC)))
}

func (ptr *TxListModel) QmlRegisterType2(uri string, versionMajor int, versionMinor int, qmlName string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	return int(int32(C.TxListModel687eda_TxListModel687eda_QmlRegisterType2(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC)))
}

func TxListModel_QmlRegisterUncreatableType(uri string, versionMajor int, versionMinor int, qmlName string, message string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	var messageC *C.char
	if message != "" {
		messageC = C.CString(message)
		defer C.free(unsafe.Pointer(messageC))
	}
	return int(int32(C.TxListModel687eda_TxListModel687eda_QmlRegisterUncreatableType(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC, C.struct_Moc_PackedString{data: messageC, len: C.longlong(len(message))})))
}

func (ptr *TxListModel) QmlRegisterUncreatableType(uri string, versionMajor int, versionMinor int, qmlName string, message string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	var messageC *C.char
	if message != "" {
		messageC = C.CString(message)
		defer C.free(unsafe.Pointer(messageC))
	}
	return int(int32(C.TxListModel687eda_TxListModel687eda_QmlRegisterUncreatableType(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC, C.struct_Moc_PackedString{data: messageC, len: C.longlong(len(message))})))
}

func (ptr *TxListModel) ____itemData_keyList_atList(i int) int {
	if ptr.Pointer() != nil {
		return int(int32(C.TxListModel687eda_____itemData_keyList_atList(ptr.Pointer(), C.int(int32(i)))))
	}
	return 0
}

func (ptr *TxListModel) ____itemData_keyList_setList(i int) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda_____itemData_keyList_setList(ptr.Pointer(), C.int(int32(i)))
	}
}

func (ptr *TxListModel) ____itemData_keyList_newList() unsafe.Pointer {
	return C.TxListModel687eda_____itemData_keyList_newList(ptr.Pointer())
}

func (ptr *TxListModel) ____roleNames_keyList_atList(i int) int {
	if ptr.Pointer() != nil {
		return int(int32(C.TxListModel687eda_____roleNames_keyList_atList(ptr.Pointer(), C.int(int32(i)))))
	}
	return 0
}

func (ptr *TxListModel) ____roleNames_keyList_setList(i int) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda_____roleNames_keyList_setList(ptr.Pointer(), C.int(int32(i)))
	}
}

func (ptr *TxListModel) ____roleNames_keyList_newList() unsafe.Pointer {
	return C.TxListModel687eda_____roleNames_keyList_newList(ptr.Pointer())
}

func (ptr *TxListModel) ____setItemData_roles_keyList_atList(i int) int {
	if ptr.Pointer() != nil {
		return int(int32(C.TxListModel687eda_____setItemData_roles_keyList_atList(ptr.Pointer(), C.int(int32(i)))))
	}
	return 0
}

func (ptr *TxListModel) ____setItemData_roles_keyList_setList(i int) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda_____setItemData_roles_keyList_setList(ptr.Pointer(), C.int(int32(i)))
	}
}

func (ptr *TxListModel) ____setItemData_roles_keyList_newList() unsafe.Pointer {
	return C.TxListModel687eda_____setItemData_roles_keyList_newList(ptr.Pointer())
}

func (ptr *TxListModel) __changePersistentIndexList_from_atList(i int) *std_core.QModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQModelIndexFromPointer(C.TxListModel687eda___changePersistentIndexList_from_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QModelIndex).DestroyQModelIndex)
		return tmpValue
	}
	return nil
}

func (ptr *TxListModel) __changePersistentIndexList_from_setList(i std_core.QModelIndex_ITF) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda___changePersistentIndexList_from_setList(ptr.Pointer(), std_core.PointerFromQModelIndex(i))
	}
}

func (ptr *TxListModel) __changePersistentIndexList_from_newList() unsafe.Pointer {
	return C.TxListModel687eda___changePersistentIndexList_from_newList(ptr.Pointer())
}

func (ptr *TxListModel) __changePersistentIndexList_to_atList(i int) *std_core.QModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQModelIndexFromPointer(C.TxListModel687eda___changePersistentIndexList_to_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QModelIndex).DestroyQModelIndex)
		return tmpValue
	}
	return nil
}

func (ptr *TxListModel) __changePersistentIndexList_to_setList(i std_core.QModelIndex_ITF) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda___changePersistentIndexList_to_setList(ptr.Pointer(), std_core.PointerFromQModelIndex(i))
	}
}

func (ptr *TxListModel) __changePersistentIndexList_to_newList() unsafe.Pointer {
	return C.TxListModel687eda___changePersistentIndexList_to_newList(ptr.Pointer())
}

func (ptr *TxListModel) __dataChanged_roles_atList(i int) int {
	if ptr.Pointer() != nil {
		return int(int32(C.TxListModel687eda___dataChanged_roles_atList(ptr.Pointer(), C.int(int32(i)))))
	}
	return 0
}

func (ptr *TxListModel) __dataChanged_roles_setList(i int) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda___dataChanged_roles_setList(ptr.Pointer(), C.int(int32(i)))
	}
}

func (ptr *TxListModel) __dataChanged_roles_newList() unsafe.Pointer {
	return C.TxListModel687eda___dataChanged_roles_newList(ptr.Pointer())
}

func (ptr *TxListModel) __itemData_atList(v int, i int) *std_core.QVariant {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQVariantFromPointer(C.TxListModel687eda___itemData_atList(ptr.Pointer(), C.int(int32(v)), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QVariant).DestroyQVariant)
		return tmpValue
	}
	return nil
}

func (ptr *TxListModel) __itemData_setList(key int, i std_core.QVariant_ITF) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda___itemData_setList(ptr.Pointer(), C.int(int32(key)), std_core.PointerFromQVariant(i))
	}
}

func (ptr *TxListModel) __itemData_newList() unsafe.Pointer {
	return C.TxListModel687eda___itemData_newList(ptr.Pointer())
}

func (ptr *TxListModel) __itemData_keyList() []int {
	if ptr.Pointer() != nil {
		return func(l C.struct_Moc_PackedList) []int {
			out := make([]int, int(l.len))
			tmpList := NewTxListModelFromPointer(l.data)
			for i := 0; i < len(out); i++ {
				out[i] = tmpList.____itemData_keyList_atList(i)
			}
			return out
		}(C.TxListModel687eda___itemData_keyList(ptr.Pointer()))
	}
	return make([]int, 0)
}

func (ptr *TxListModel) __layoutAboutToBeChanged_parents_atList(i int) *std_core.QPersistentModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQPersistentModelIndexFromPointer(C.TxListModel687eda___layoutAboutToBeChanged_parents_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QPersistentModelIndex).DestroyQPersistentModelIndex)
		return tmpValue
	}
	return nil
}

func (ptr *TxListModel) __layoutAboutToBeChanged_parents_setList(i std_core.QPersistentModelIndex_ITF) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda___layoutAboutToBeChanged_parents_setList(ptr.Pointer(), std_core.PointerFromQPersistentModelIndex(i))
	}
}

func (ptr *TxListModel) __layoutAboutToBeChanged_parents_newList() unsafe.Pointer {
	return C.TxListModel687eda___layoutAboutToBeChanged_parents_newList(ptr.Pointer())
}

func (ptr *TxListModel) __layoutChanged_parents_atList(i int) *std_core.QPersistentModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQPersistentModelIndexFromPointer(C.TxListModel687eda___layoutChanged_parents_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QPersistentModelIndex).DestroyQPersistentModelIndex)
		return tmpValue
	}
	return nil
}

func (ptr *TxListModel) __layoutChanged_parents_setList(i std_core.QPersistentModelIndex_ITF) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda___layoutChanged_parents_setList(ptr.Pointer(), std_core.PointerFromQPersistentModelIndex(i))
	}
}

func (ptr *TxListModel) __layoutChanged_parents_newList() unsafe.Pointer {
	return C.TxListModel687eda___layoutChanged_parents_newList(ptr.Pointer())
}

func (ptr *TxListModel) __match_atList(i int) *std_core.QModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQModelIndexFromPointer(C.TxListModel687eda___match_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QModelIndex).DestroyQModelIndex)
		return tmpValue
	}
	return nil
}

func (ptr *TxListModel) __match_setList(i std_core.QModelIndex_ITF) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda___match_setList(ptr.Pointer(), std_core.PointerFromQModelIndex(i))
	}
}

func (ptr *TxListModel) __match_newList() unsafe.Pointer {
	return C.TxListModel687eda___match_newList(ptr.Pointer())
}

func (ptr *TxListModel) __mimeData_indexes_atList(i int) *std_core.QModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQModelIndexFromPointer(C.TxListModel687eda___mimeData_indexes_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QModelIndex).DestroyQModelIndex)
		return tmpValue
	}
	return nil
}

func (ptr *TxListModel) __mimeData_indexes_setList(i std_core.QModelIndex_ITF) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda___mimeData_indexes_setList(ptr.Pointer(), std_core.PointerFromQModelIndex(i))
	}
}

func (ptr *TxListModel) __mimeData_indexes_newList() unsafe.Pointer {
	return C.TxListModel687eda___mimeData_indexes_newList(ptr.Pointer())
}

func (ptr *TxListModel) __persistentIndexList_atList(i int) *std_core.QModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQModelIndexFromPointer(C.TxListModel687eda___persistentIndexList_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QModelIndex).DestroyQModelIndex)
		return tmpValue
	}
	return nil
}

func (ptr *TxListModel) __persistentIndexList_setList(i std_core.QModelIndex_ITF) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda___persistentIndexList_setList(ptr.Pointer(), std_core.PointerFromQModelIndex(i))
	}
}

func (ptr *TxListModel) __persistentIndexList_newList() unsafe.Pointer {
	return C.TxListModel687eda___persistentIndexList_newList(ptr.Pointer())
}

func (ptr *TxListModel) __roleNames_atList(v int, i int) *std_core.QByteArray {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQByteArrayFromPointer(C.TxListModel687eda___roleNames_atList(ptr.Pointer(), C.int(int32(v)), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QByteArray).DestroyQByteArray)
		return tmpValue
	}
	return nil
}

func (ptr *TxListModel) __roleNames_setList(key int, i std_core.QByteArray_ITF) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda___roleNames_setList(ptr.Pointer(), C.int(int32(key)), std_core.PointerFromQByteArray(i))
	}
}

func (ptr *TxListModel) __roleNames_newList() unsafe.Pointer {
	return C.TxListModel687eda___roleNames_newList(ptr.Pointer())
}

func (ptr *TxListModel) __roleNames_keyList() []int {
	if ptr.Pointer() != nil {
		return func(l C.struct_Moc_PackedList) []int {
			out := make([]int, int(l.len))
			tmpList := NewTxListModelFromPointer(l.data)
			for i := 0; i < len(out); i++ {
				out[i] = tmpList.____roleNames_keyList_atList(i)
			}
			return out
		}(C.TxListModel687eda___roleNames_keyList(ptr.Pointer()))
	}
	return make([]int, 0)
}

func (ptr *TxListModel) __setItemData_roles_atList(v int, i int) *std_core.QVariant {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQVariantFromPointer(C.TxListModel687eda___setItemData_roles_atList(ptr.Pointer(), C.int(int32(v)), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QVariant).DestroyQVariant)
		return tmpValue
	}
	return nil
}

func (ptr *TxListModel) __setItemData_roles_setList(key int, i std_core.QVariant_ITF) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda___setItemData_roles_setList(ptr.Pointer(), C.int(int32(key)), std_core.PointerFromQVariant(i))
	}
}

func (ptr *TxListModel) __setItemData_roles_newList() unsafe.Pointer {
	return C.TxListModel687eda___setItemData_roles_newList(ptr.Pointer())
}

func (ptr *TxListModel) __setItemData_roles_keyList() []int {
	if ptr.Pointer() != nil {
		return func(l C.struct_Moc_PackedList) []int {
			out := make([]int, int(l.len))
			tmpList := NewTxListModelFromPointer(l.data)
			for i := 0; i < len(out); i++ {
				out[i] = tmpList.____setItemData_roles_keyList_atList(i)
			}
			return out
		}(C.TxListModel687eda___setItemData_roles_keyList(ptr.Pointer()))
	}
	return make([]int, 0)
}

func (ptr *TxListModel) ____doSetRoleNames_roleNames_keyList_atList(i int) int {
	if ptr.Pointer() != nil {
		return int(int32(C.TxListModel687eda_____doSetRoleNames_roleNames_keyList_atList(ptr.Pointer(), C.int(int32(i)))))
	}
	return 0
}

func (ptr *TxListModel) ____doSetRoleNames_roleNames_keyList_setList(i int) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda_____doSetRoleNames_roleNames_keyList_setList(ptr.Pointer(), C.int(int32(i)))
	}
}

func (ptr *TxListModel) ____doSetRoleNames_roleNames_keyList_newList() unsafe.Pointer {
	return C.TxListModel687eda_____doSetRoleNames_roleNames_keyList_newList(ptr.Pointer())
}

func (ptr *TxListModel) ____setRoleNames_roleNames_keyList_atList(i int) int {
	if ptr.Pointer() != nil {
		return int(int32(C.TxListModel687eda_____setRoleNames_roleNames_keyList_atList(ptr.Pointer(), C.int(int32(i)))))
	}
	return 0
}

func (ptr *TxListModel) ____setRoleNames_roleNames_keyList_setList(i int) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda_____setRoleNames_roleNames_keyList_setList(ptr.Pointer(), C.int(int32(i)))
	}
}

func (ptr *TxListModel) ____setRoleNames_roleNames_keyList_newList() unsafe.Pointer {
	return C.TxListModel687eda_____setRoleNames_roleNames_keyList_newList(ptr.Pointer())
}

func (ptr *TxListModel) __children_atList(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQObjectFromPointer(C.TxListModel687eda___children_atList(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *TxListModel) __children_setList(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda___children_setList(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *TxListModel) __children_newList() unsafe.Pointer {
	return C.TxListModel687eda___children_newList(ptr.Pointer())
}

func (ptr *TxListModel) __dynamicPropertyNames_atList(i int) *std_core.QByteArray {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQByteArrayFromPointer(C.TxListModel687eda___dynamicPropertyNames_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QByteArray).DestroyQByteArray)
		return tmpValue
	}
	return nil
}

func (ptr *TxListModel) __dynamicPropertyNames_setList(i std_core.QByteArray_ITF) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda___dynamicPropertyNames_setList(ptr.Pointer(), std_core.PointerFromQByteArray(i))
	}
}

func (ptr *TxListModel) __dynamicPropertyNames_newList() unsafe.Pointer {
	return C.TxListModel687eda___dynamicPropertyNames_newList(ptr.Pointer())
}

func (ptr *TxListModel) __findChildren_atList(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQObjectFromPointer(C.TxListModel687eda___findChildren_atList(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *TxListModel) __findChildren_setList(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda___findChildren_setList(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *TxListModel) __findChildren_newList() unsafe.Pointer {
	return C.TxListModel687eda___findChildren_newList(ptr.Pointer())
}

func (ptr *TxListModel) __findChildren_atList3(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQObjectFromPointer(C.TxListModel687eda___findChildren_atList3(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *TxListModel) __findChildren_setList3(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda___findChildren_setList3(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *TxListModel) __findChildren_newList3() unsafe.Pointer {
	return C.TxListModel687eda___findChildren_newList3(ptr.Pointer())
}

func NewTxListModel(parent std_core.QObject_ITF) *TxListModel {
	TxListModel_QRegisterMetaType()
	tmpValue := NewTxListModelFromPointer(C.TxListModel687eda_NewTxListModel(std_core.PointerFromQObject(parent)))
	if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
		tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
	}
	return tmpValue
}

//export callbackTxListModel687eda_DestroyTxListModel
func callbackTxListModel687eda_DestroyTxListModel(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "~TxListModel"); signal != nil {
		(*(*func())(signal))()
	} else {
		NewTxListModelFromPointer(ptr).DestroyTxListModelDefault()
	}
}

func (ptr *TxListModel) ConnectDestroyTxListModel(f func()) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "~TxListModel"); signal != nil {
			f := func() {
				(*(*func())(signal))()
				f()
			}
			qt.ConnectSignal(ptr.Pointer(), "~TxListModel", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "~TxListModel", unsafe.Pointer(&f))
		}
	}
}

func (ptr *TxListModel) DisconnectDestroyTxListModel() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "~TxListModel")
	}
}

func (ptr *TxListModel) DestroyTxListModel() {
	if ptr.Pointer() != nil {

		qt.SetFinalizer(ptr, nil)
		C.TxListModel687eda_DestroyTxListModel(ptr.Pointer())
		ptr.SetPointer(nil)
	}
}

func (ptr *TxListModel) DestroyTxListModelDefault() {
	if ptr.Pointer() != nil {

		qt.SetFinalizer(ptr, nil)
		C.TxListModel687eda_DestroyTxListModelDefault(ptr.Pointer())
		ptr.SetPointer(nil)
	}
}

//export callbackTxListModel687eda_DropMimeData
func callbackTxListModel687eda_DropMimeData(ptr unsafe.Pointer, data unsafe.Pointer, action C.longlong, row C.int, column C.int, parent unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "dropMimeData"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QMimeData, std_core.Qt__DropAction, int, int, *std_core.QModelIndex) bool)(signal))(std_core.NewQMimeDataFromPointer(data), std_core.Qt__DropAction(action), int(int32(row)), int(int32(column)), std_core.NewQModelIndexFromPointer(parent)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListModelFromPointer(ptr).DropMimeDataDefault(std_core.NewQMimeDataFromPointer(data), std_core.Qt__DropAction(action), int(int32(row)), int(int32(column)), std_core.NewQModelIndexFromPointer(parent)))))
}

func (ptr *TxListModel) DropMimeDataDefault(data std_core.QMimeData_ITF, action std_core.Qt__DropAction, row int, column int, parent std_core.QModelIndex_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListModel687eda_DropMimeDataDefault(ptr.Pointer(), std_core.PointerFromQMimeData(data), C.longlong(action), C.int(int32(row)), C.int(int32(column)), std_core.PointerFromQModelIndex(parent))) != 0
	}
	return false
}

//export callbackTxListModel687eda_Flags
func callbackTxListModel687eda_Flags(ptr unsafe.Pointer, index unsafe.Pointer) C.longlong {
	if signal := qt.GetSignal(ptr, "flags"); signal != nil {
		return C.longlong((*(*func(*std_core.QModelIndex) std_core.Qt__ItemFlag)(signal))(std_core.NewQModelIndexFromPointer(index)))
	}

	return C.longlong(NewTxListModelFromPointer(ptr).FlagsDefault(std_core.NewQModelIndexFromPointer(index)))
}

func (ptr *TxListModel) FlagsDefault(index std_core.QModelIndex_ITF) std_core.Qt__ItemFlag {
	if ptr.Pointer() != nil {
		return std_core.Qt__ItemFlag(C.TxListModel687eda_FlagsDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(index)))
	}
	return 0
}

//export callbackTxListModel687eda_Index
func callbackTxListModel687eda_Index(ptr unsafe.Pointer, row C.int, column C.int, parent unsafe.Pointer) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "index"); signal != nil {
		return std_core.PointerFromQModelIndex((*(*func(int, int, *std_core.QModelIndex) *std_core.QModelIndex)(signal))(int(int32(row)), int(int32(column)), std_core.NewQModelIndexFromPointer(parent)))
	}

	return std_core.PointerFromQModelIndex(NewTxListModelFromPointer(ptr).IndexDefault(int(int32(row)), int(int32(column)), std_core.NewQModelIndexFromPointer(parent)))
}

func (ptr *TxListModel) IndexDefault(row int, column int, parent std_core.QModelIndex_ITF) *std_core.QModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQModelIndexFromPointer(C.TxListModel687eda_IndexDefault(ptr.Pointer(), C.int(int32(row)), C.int(int32(column)), std_core.PointerFromQModelIndex(parent)))
		qt.SetFinalizer(tmpValue, (*std_core.QModelIndex).DestroyQModelIndex)
		return tmpValue
	}
	return nil
}

//export callbackTxListModel687eda_Sibling
func callbackTxListModel687eda_Sibling(ptr unsafe.Pointer, row C.int, column C.int, idx unsafe.Pointer) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "sibling"); signal != nil {
		return std_core.PointerFromQModelIndex((*(*func(int, int, *std_core.QModelIndex) *std_core.QModelIndex)(signal))(int(int32(row)), int(int32(column)), std_core.NewQModelIndexFromPointer(idx)))
	}

	return std_core.PointerFromQModelIndex(NewTxListModelFromPointer(ptr).SiblingDefault(int(int32(row)), int(int32(column)), std_core.NewQModelIndexFromPointer(idx)))
}

func (ptr *TxListModel) SiblingDefault(row int, column int, idx std_core.QModelIndex_ITF) *std_core.QModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQModelIndexFromPointer(C.TxListModel687eda_SiblingDefault(ptr.Pointer(), C.int(int32(row)), C.int(int32(column)), std_core.PointerFromQModelIndex(idx)))
		qt.SetFinalizer(tmpValue, (*std_core.QModelIndex).DestroyQModelIndex)
		return tmpValue
	}
	return nil
}

//export callbackTxListModel687eda_Buddy
func callbackTxListModel687eda_Buddy(ptr unsafe.Pointer, index unsafe.Pointer) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "buddy"); signal != nil {
		return std_core.PointerFromQModelIndex((*(*func(*std_core.QModelIndex) *std_core.QModelIndex)(signal))(std_core.NewQModelIndexFromPointer(index)))
	}

	return std_core.PointerFromQModelIndex(NewTxListModelFromPointer(ptr).BuddyDefault(std_core.NewQModelIndexFromPointer(index)))
}

func (ptr *TxListModel) BuddyDefault(index std_core.QModelIndex_ITF) *std_core.QModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQModelIndexFromPointer(C.TxListModel687eda_BuddyDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(index)))
		qt.SetFinalizer(tmpValue, (*std_core.QModelIndex).DestroyQModelIndex)
		return tmpValue
	}
	return nil
}

//export callbackTxListModel687eda_CanDropMimeData
func callbackTxListModel687eda_CanDropMimeData(ptr unsafe.Pointer, data unsafe.Pointer, action C.longlong, row C.int, column C.int, parent unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "canDropMimeData"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QMimeData, std_core.Qt__DropAction, int, int, *std_core.QModelIndex) bool)(signal))(std_core.NewQMimeDataFromPointer(data), std_core.Qt__DropAction(action), int(int32(row)), int(int32(column)), std_core.NewQModelIndexFromPointer(parent)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListModelFromPointer(ptr).CanDropMimeDataDefault(std_core.NewQMimeDataFromPointer(data), std_core.Qt__DropAction(action), int(int32(row)), int(int32(column)), std_core.NewQModelIndexFromPointer(parent)))))
}

func (ptr *TxListModel) CanDropMimeDataDefault(data std_core.QMimeData_ITF, action std_core.Qt__DropAction, row int, column int, parent std_core.QModelIndex_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListModel687eda_CanDropMimeDataDefault(ptr.Pointer(), std_core.PointerFromQMimeData(data), C.longlong(action), C.int(int32(row)), C.int(int32(column)), std_core.PointerFromQModelIndex(parent))) != 0
	}
	return false
}

//export callbackTxListModel687eda_CanFetchMore
func callbackTxListModel687eda_CanFetchMore(ptr unsafe.Pointer, parent unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "canFetchMore"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QModelIndex) bool)(signal))(std_core.NewQModelIndexFromPointer(parent)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListModelFromPointer(ptr).CanFetchMoreDefault(std_core.NewQModelIndexFromPointer(parent)))))
}

func (ptr *TxListModel) CanFetchMoreDefault(parent std_core.QModelIndex_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListModel687eda_CanFetchMoreDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(parent))) != 0
	}
	return false
}

//export callbackTxListModel687eda_ColumnCount
func callbackTxListModel687eda_ColumnCount(ptr unsafe.Pointer, parent unsafe.Pointer) C.int {
	if signal := qt.GetSignal(ptr, "columnCount"); signal != nil {
		return C.int(int32((*(*func(*std_core.QModelIndex) int)(signal))(std_core.NewQModelIndexFromPointer(parent))))
	}

	return C.int(int32(NewTxListModelFromPointer(ptr).ColumnCountDefault(std_core.NewQModelIndexFromPointer(parent))))
}

func (ptr *TxListModel) ColumnCountDefault(parent std_core.QModelIndex_ITF) int {
	if ptr.Pointer() != nil {
		return int(int32(C.TxListModel687eda_ColumnCountDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(parent))))
	}
	return 0
}

//export callbackTxListModel687eda_ColumnsAboutToBeInserted
func callbackTxListModel687eda_ColumnsAboutToBeInserted(ptr unsafe.Pointer, parent unsafe.Pointer, first C.int, last C.int) {
	if signal := qt.GetSignal(ptr, "columnsAboutToBeInserted"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(first)), int(int32(last)))
	}

}

//export callbackTxListModel687eda_ColumnsAboutToBeMoved
func callbackTxListModel687eda_ColumnsAboutToBeMoved(ptr unsafe.Pointer, sourceParent unsafe.Pointer, sourceStart C.int, sourceEnd C.int, destinationParent unsafe.Pointer, destinationColumn C.int) {
	if signal := qt.GetSignal(ptr, "columnsAboutToBeMoved"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int, *std_core.QModelIndex, int))(signal))(std_core.NewQModelIndexFromPointer(sourceParent), int(int32(sourceStart)), int(int32(sourceEnd)), std_core.NewQModelIndexFromPointer(destinationParent), int(int32(destinationColumn)))
	}

}

//export callbackTxListModel687eda_ColumnsAboutToBeRemoved
func callbackTxListModel687eda_ColumnsAboutToBeRemoved(ptr unsafe.Pointer, parent unsafe.Pointer, first C.int, last C.int) {
	if signal := qt.GetSignal(ptr, "columnsAboutToBeRemoved"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(first)), int(int32(last)))
	}

}

//export callbackTxListModel687eda_ColumnsInserted
func callbackTxListModel687eda_ColumnsInserted(ptr unsafe.Pointer, parent unsafe.Pointer, first C.int, last C.int) {
	if signal := qt.GetSignal(ptr, "columnsInserted"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(first)), int(int32(last)))
	}

}

//export callbackTxListModel687eda_ColumnsMoved
func callbackTxListModel687eda_ColumnsMoved(ptr unsafe.Pointer, parent unsafe.Pointer, start C.int, end C.int, destination unsafe.Pointer, column C.int) {
	if signal := qt.GetSignal(ptr, "columnsMoved"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int, *std_core.QModelIndex, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(start)), int(int32(end)), std_core.NewQModelIndexFromPointer(destination), int(int32(column)))
	}

}

//export callbackTxListModel687eda_ColumnsRemoved
func callbackTxListModel687eda_ColumnsRemoved(ptr unsafe.Pointer, parent unsafe.Pointer, first C.int, last C.int) {
	if signal := qt.GetSignal(ptr, "columnsRemoved"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(first)), int(int32(last)))
	}

}

//export callbackTxListModel687eda_Data
func callbackTxListModel687eda_Data(ptr unsafe.Pointer, index unsafe.Pointer, role C.int) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "data"); signal != nil {
		return std_core.PointerFromQVariant((*(*func(*std_core.QModelIndex, int) *std_core.QVariant)(signal))(std_core.NewQModelIndexFromPointer(index), int(int32(role))))
	}

	return std_core.PointerFromQVariant(NewTxListModelFromPointer(ptr).DataDefault(std_core.NewQModelIndexFromPointer(index), int(int32(role))))
}

func (ptr *TxListModel) DataDefault(index std_core.QModelIndex_ITF, role int) *std_core.QVariant {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQVariantFromPointer(C.TxListModel687eda_DataDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(index), C.int(int32(role))))
		qt.SetFinalizer(tmpValue, (*std_core.QVariant).DestroyQVariant)
		return tmpValue
	}
	return nil
}

//export callbackTxListModel687eda_DataChanged
func callbackTxListModel687eda_DataChanged(ptr unsafe.Pointer, topLeft unsafe.Pointer, bottomRight unsafe.Pointer, roles C.struct_Moc_PackedList) {
	if signal := qt.GetSignal(ptr, "dataChanged"); signal != nil {
		(*(*func(*std_core.QModelIndex, *std_core.QModelIndex, []int))(signal))(std_core.NewQModelIndexFromPointer(topLeft), std_core.NewQModelIndexFromPointer(bottomRight), func(l C.struct_Moc_PackedList) []int {
			out := make([]int, int(l.len))
			tmpList := NewTxListModelFromPointer(l.data)
			for i := 0; i < len(out); i++ {
				out[i] = tmpList.__dataChanged_roles_atList(i)
			}
			return out
		}(roles))
	}

}

//export callbackTxListModel687eda_FetchMore
func callbackTxListModel687eda_FetchMore(ptr unsafe.Pointer, parent unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "fetchMore"); signal != nil {
		(*(*func(*std_core.QModelIndex))(signal))(std_core.NewQModelIndexFromPointer(parent))
	} else {
		NewTxListModelFromPointer(ptr).FetchMoreDefault(std_core.NewQModelIndexFromPointer(parent))
	}
}

func (ptr *TxListModel) FetchMoreDefault(parent std_core.QModelIndex_ITF) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda_FetchMoreDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(parent))
	}
}

//export callbackTxListModel687eda_HasChildren
func callbackTxListModel687eda_HasChildren(ptr unsafe.Pointer, parent unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "hasChildren"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QModelIndex) bool)(signal))(std_core.NewQModelIndexFromPointer(parent)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListModelFromPointer(ptr).HasChildrenDefault(std_core.NewQModelIndexFromPointer(parent)))))
}

func (ptr *TxListModel) HasChildrenDefault(parent std_core.QModelIndex_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListModel687eda_HasChildrenDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(parent))) != 0
	}
	return false
}

//export callbackTxListModel687eda_HeaderData
func callbackTxListModel687eda_HeaderData(ptr unsafe.Pointer, section C.int, orientation C.longlong, role C.int) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "headerData"); signal != nil {
		return std_core.PointerFromQVariant((*(*func(int, std_core.Qt__Orientation, int) *std_core.QVariant)(signal))(int(int32(section)), std_core.Qt__Orientation(orientation), int(int32(role))))
	}

	return std_core.PointerFromQVariant(NewTxListModelFromPointer(ptr).HeaderDataDefault(int(int32(section)), std_core.Qt__Orientation(orientation), int(int32(role))))
}

func (ptr *TxListModel) HeaderDataDefault(section int, orientation std_core.Qt__Orientation, role int) *std_core.QVariant {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQVariantFromPointer(C.TxListModel687eda_HeaderDataDefault(ptr.Pointer(), C.int(int32(section)), C.longlong(orientation), C.int(int32(role))))
		qt.SetFinalizer(tmpValue, (*std_core.QVariant).DestroyQVariant)
		return tmpValue
	}
	return nil
}

//export callbackTxListModel687eda_HeaderDataChanged
func callbackTxListModel687eda_HeaderDataChanged(ptr unsafe.Pointer, orientation C.longlong, first C.int, last C.int) {
	if signal := qt.GetSignal(ptr, "headerDataChanged"); signal != nil {
		(*(*func(std_core.Qt__Orientation, int, int))(signal))(std_core.Qt__Orientation(orientation), int(int32(first)), int(int32(last)))
	}

}

//export callbackTxListModel687eda_InsertColumns
func callbackTxListModel687eda_InsertColumns(ptr unsafe.Pointer, column C.int, count C.int, parent unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "insertColumns"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(int, int, *std_core.QModelIndex) bool)(signal))(int(int32(column)), int(int32(count)), std_core.NewQModelIndexFromPointer(parent)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListModelFromPointer(ptr).InsertColumnsDefault(int(int32(column)), int(int32(count)), std_core.NewQModelIndexFromPointer(parent)))))
}

func (ptr *TxListModel) InsertColumnsDefault(column int, count int, parent std_core.QModelIndex_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListModel687eda_InsertColumnsDefault(ptr.Pointer(), C.int(int32(column)), C.int(int32(count)), std_core.PointerFromQModelIndex(parent))) != 0
	}
	return false
}

//export callbackTxListModel687eda_InsertRows
func callbackTxListModel687eda_InsertRows(ptr unsafe.Pointer, row C.int, count C.int, parent unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "insertRows"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(int, int, *std_core.QModelIndex) bool)(signal))(int(int32(row)), int(int32(count)), std_core.NewQModelIndexFromPointer(parent)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListModelFromPointer(ptr).InsertRowsDefault(int(int32(row)), int(int32(count)), std_core.NewQModelIndexFromPointer(parent)))))
}

func (ptr *TxListModel) InsertRowsDefault(row int, count int, parent std_core.QModelIndex_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListModel687eda_InsertRowsDefault(ptr.Pointer(), C.int(int32(row)), C.int(int32(count)), std_core.PointerFromQModelIndex(parent))) != 0
	}
	return false
}

//export callbackTxListModel687eda_ItemData
func callbackTxListModel687eda_ItemData(ptr unsafe.Pointer, index unsafe.Pointer) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "itemData"); signal != nil {
		return func() unsafe.Pointer {
			tmpList := NewTxListModelFromPointer(NewTxListModelFromPointer(nil).__itemData_newList())
			for k, v := range (*(*func(*std_core.QModelIndex) map[int]*std_core.QVariant)(signal))(std_core.NewQModelIndexFromPointer(index)) {
				tmpList.__itemData_setList(k, v)
			}
			return tmpList.Pointer()
		}()
	}

	return func() unsafe.Pointer {
		tmpList := NewTxListModelFromPointer(NewTxListModelFromPointer(nil).__itemData_newList())
		for k, v := range NewTxListModelFromPointer(ptr).ItemDataDefault(std_core.NewQModelIndexFromPointer(index)) {
			tmpList.__itemData_setList(k, v)
		}
		return tmpList.Pointer()
	}()
}

func (ptr *TxListModel) ItemDataDefault(index std_core.QModelIndex_ITF) map[int]*std_core.QVariant {
	if ptr.Pointer() != nil {
		return func(l C.struct_Moc_PackedList) map[int]*std_core.QVariant {
			out := make(map[int]*std_core.QVariant, int(l.len))
			tmpList := NewTxListModelFromPointer(l.data)
			for i, v := range tmpList.__itemData_keyList() {
				out[v] = tmpList.__itemData_atList(v, i)
			}
			return out
		}(C.TxListModel687eda_ItemDataDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(index)))
	}
	return make(map[int]*std_core.QVariant, 0)
}

//export callbackTxListModel687eda_LayoutAboutToBeChanged
func callbackTxListModel687eda_LayoutAboutToBeChanged(ptr unsafe.Pointer, parents C.struct_Moc_PackedList, hint C.longlong) {
	if signal := qt.GetSignal(ptr, "layoutAboutToBeChanged"); signal != nil {
		(*(*func([]*std_core.QPersistentModelIndex, std_core.QAbstractItemModel__LayoutChangeHint))(signal))(func(l C.struct_Moc_PackedList) []*std_core.QPersistentModelIndex {
			out := make([]*std_core.QPersistentModelIndex, int(l.len))
			tmpList := NewTxListModelFromPointer(l.data)
			for i := 0; i < len(out); i++ {
				out[i] = tmpList.__layoutAboutToBeChanged_parents_atList(i)
			}
			return out
		}(parents), std_core.QAbstractItemModel__LayoutChangeHint(hint))
	}

}

//export callbackTxListModel687eda_LayoutChanged
func callbackTxListModel687eda_LayoutChanged(ptr unsafe.Pointer, parents C.struct_Moc_PackedList, hint C.longlong) {
	if signal := qt.GetSignal(ptr, "layoutChanged"); signal != nil {
		(*(*func([]*std_core.QPersistentModelIndex, std_core.QAbstractItemModel__LayoutChangeHint))(signal))(func(l C.struct_Moc_PackedList) []*std_core.QPersistentModelIndex {
			out := make([]*std_core.QPersistentModelIndex, int(l.len))
			tmpList := NewTxListModelFromPointer(l.data)
			for i := 0; i < len(out); i++ {
				out[i] = tmpList.__layoutChanged_parents_atList(i)
			}
			return out
		}(parents), std_core.QAbstractItemModel__LayoutChangeHint(hint))
	}

}

//export callbackTxListModel687eda_Match
func callbackTxListModel687eda_Match(ptr unsafe.Pointer, start unsafe.Pointer, role C.int, value unsafe.Pointer, hits C.int, flags C.longlong) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "match"); signal != nil {
		return func() unsafe.Pointer {
			tmpList := NewTxListModelFromPointer(NewTxListModelFromPointer(nil).__match_newList())
			for _, v := range (*(*func(*std_core.QModelIndex, int, *std_core.QVariant, int, std_core.Qt__MatchFlag) []*std_core.QModelIndex)(signal))(std_core.NewQModelIndexFromPointer(start), int(int32(role)), std_core.NewQVariantFromPointer(value), int(int32(hits)), std_core.Qt__MatchFlag(flags)) {
				tmpList.__match_setList(v)
			}
			return tmpList.Pointer()
		}()
	}

	return func() unsafe.Pointer {
		tmpList := NewTxListModelFromPointer(NewTxListModelFromPointer(nil).__match_newList())
		for _, v := range NewTxListModelFromPointer(ptr).MatchDefault(std_core.NewQModelIndexFromPointer(start), int(int32(role)), std_core.NewQVariantFromPointer(value), int(int32(hits)), std_core.Qt__MatchFlag(flags)) {
			tmpList.__match_setList(v)
		}
		return tmpList.Pointer()
	}()
}

func (ptr *TxListModel) MatchDefault(start std_core.QModelIndex_ITF, role int, value std_core.QVariant_ITF, hits int, flags std_core.Qt__MatchFlag) []*std_core.QModelIndex {
	if ptr.Pointer() != nil {
		return func(l C.struct_Moc_PackedList) []*std_core.QModelIndex {
			out := make([]*std_core.QModelIndex, int(l.len))
			tmpList := NewTxListModelFromPointer(l.data)
			for i := 0; i < len(out); i++ {
				out[i] = tmpList.__match_atList(i)
			}
			return out
		}(C.TxListModel687eda_MatchDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(start), C.int(int32(role)), std_core.PointerFromQVariant(value), C.int(int32(hits)), C.longlong(flags)))
	}
	return make([]*std_core.QModelIndex, 0)
}

//export callbackTxListModel687eda_MimeData
func callbackTxListModel687eda_MimeData(ptr unsafe.Pointer, indexes C.struct_Moc_PackedList) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "mimeData"); signal != nil {
		return std_core.PointerFromQMimeData((*(*func([]*std_core.QModelIndex) *std_core.QMimeData)(signal))(func(l C.struct_Moc_PackedList) []*std_core.QModelIndex {
			out := make([]*std_core.QModelIndex, int(l.len))
			tmpList := NewTxListModelFromPointer(l.data)
			for i := 0; i < len(out); i++ {
				out[i] = tmpList.__mimeData_indexes_atList(i)
			}
			return out
		}(indexes)))
	}

	return std_core.PointerFromQMimeData(NewTxListModelFromPointer(ptr).MimeDataDefault(func(l C.struct_Moc_PackedList) []*std_core.QModelIndex {
		out := make([]*std_core.QModelIndex, int(l.len))
		tmpList := NewTxListModelFromPointer(l.data)
		for i := 0; i < len(out); i++ {
			out[i] = tmpList.__mimeData_indexes_atList(i)
		}
		return out
	}(indexes)))
}

func (ptr *TxListModel) MimeDataDefault(indexes []*std_core.QModelIndex) *std_core.QMimeData {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQMimeDataFromPointer(C.TxListModel687eda_MimeDataDefault(ptr.Pointer(), func() unsafe.Pointer {
			tmpList := NewTxListModelFromPointer(NewTxListModelFromPointer(nil).__mimeData_indexes_newList())
			for _, v := range indexes {
				tmpList.__mimeData_indexes_setList(v)
			}
			return tmpList.Pointer()
		}()))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

//export callbackTxListModel687eda_MimeTypes
func callbackTxListModel687eda_MimeTypes(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "mimeTypes"); signal != nil {
		tempVal := (*(*func() []string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(strings.Join(tempVal, "¡¦!")), len: C.longlong(len(strings.Join(tempVal, "¡¦!")))}
	}
	tempVal := NewTxListModelFromPointer(ptr).MimeTypesDefault()
	return C.struct_Moc_PackedString{data: C.CString(strings.Join(tempVal, "¡¦!")), len: C.longlong(len(strings.Join(tempVal, "¡¦!")))}
}

func (ptr *TxListModel) MimeTypesDefault() []string {
	if ptr.Pointer() != nil {
		return unpackStringList(cGoUnpackString(C.TxListModel687eda_MimeTypesDefault(ptr.Pointer())))
	}
	return make([]string, 0)
}

//export callbackTxListModel687eda_ModelAboutToBeReset
func callbackTxListModel687eda_ModelAboutToBeReset(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "modelAboutToBeReset"); signal != nil {
		(*(*func())(signal))()
	}

}

//export callbackTxListModel687eda_ModelReset
func callbackTxListModel687eda_ModelReset(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "modelReset"); signal != nil {
		(*(*func())(signal))()
	}

}

//export callbackTxListModel687eda_MoveColumns
func callbackTxListModel687eda_MoveColumns(ptr unsafe.Pointer, sourceParent unsafe.Pointer, sourceColumn C.int, count C.int, destinationParent unsafe.Pointer, destinationChild C.int) C.char {
	if signal := qt.GetSignal(ptr, "moveColumns"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QModelIndex, int, int, *std_core.QModelIndex, int) bool)(signal))(std_core.NewQModelIndexFromPointer(sourceParent), int(int32(sourceColumn)), int(int32(count)), std_core.NewQModelIndexFromPointer(destinationParent), int(int32(destinationChild))))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListModelFromPointer(ptr).MoveColumnsDefault(std_core.NewQModelIndexFromPointer(sourceParent), int(int32(sourceColumn)), int(int32(count)), std_core.NewQModelIndexFromPointer(destinationParent), int(int32(destinationChild))))))
}

func (ptr *TxListModel) MoveColumnsDefault(sourceParent std_core.QModelIndex_ITF, sourceColumn int, count int, destinationParent std_core.QModelIndex_ITF, destinationChild int) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListModel687eda_MoveColumnsDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(sourceParent), C.int(int32(sourceColumn)), C.int(int32(count)), std_core.PointerFromQModelIndex(destinationParent), C.int(int32(destinationChild)))) != 0
	}
	return false
}

//export callbackTxListModel687eda_MoveRows
func callbackTxListModel687eda_MoveRows(ptr unsafe.Pointer, sourceParent unsafe.Pointer, sourceRow C.int, count C.int, destinationParent unsafe.Pointer, destinationChild C.int) C.char {
	if signal := qt.GetSignal(ptr, "moveRows"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QModelIndex, int, int, *std_core.QModelIndex, int) bool)(signal))(std_core.NewQModelIndexFromPointer(sourceParent), int(int32(sourceRow)), int(int32(count)), std_core.NewQModelIndexFromPointer(destinationParent), int(int32(destinationChild))))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListModelFromPointer(ptr).MoveRowsDefault(std_core.NewQModelIndexFromPointer(sourceParent), int(int32(sourceRow)), int(int32(count)), std_core.NewQModelIndexFromPointer(destinationParent), int(int32(destinationChild))))))
}

func (ptr *TxListModel) MoveRowsDefault(sourceParent std_core.QModelIndex_ITF, sourceRow int, count int, destinationParent std_core.QModelIndex_ITF, destinationChild int) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListModel687eda_MoveRowsDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(sourceParent), C.int(int32(sourceRow)), C.int(int32(count)), std_core.PointerFromQModelIndex(destinationParent), C.int(int32(destinationChild)))) != 0
	}
	return false
}

//export callbackTxListModel687eda_Parent
func callbackTxListModel687eda_Parent(ptr unsafe.Pointer, index unsafe.Pointer) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "parent"); signal != nil {
		return std_core.PointerFromQModelIndex((*(*func(*std_core.QModelIndex) *std_core.QModelIndex)(signal))(std_core.NewQModelIndexFromPointer(index)))
	}

	return std_core.PointerFromQModelIndex(NewTxListModelFromPointer(ptr).ParentDefault(std_core.NewQModelIndexFromPointer(index)))
}

func (ptr *TxListModel) ParentDefault(index std_core.QModelIndex_ITF) *std_core.QModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQModelIndexFromPointer(C.TxListModel687eda_ParentDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(index)))
		qt.SetFinalizer(tmpValue, (*std_core.QModelIndex).DestroyQModelIndex)
		return tmpValue
	}
	return nil
}

//export callbackTxListModel687eda_RemoveColumns
func callbackTxListModel687eda_RemoveColumns(ptr unsafe.Pointer, column C.int, count C.int, parent unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "removeColumns"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(int, int, *std_core.QModelIndex) bool)(signal))(int(int32(column)), int(int32(count)), std_core.NewQModelIndexFromPointer(parent)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListModelFromPointer(ptr).RemoveColumnsDefault(int(int32(column)), int(int32(count)), std_core.NewQModelIndexFromPointer(parent)))))
}

func (ptr *TxListModel) RemoveColumnsDefault(column int, count int, parent std_core.QModelIndex_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListModel687eda_RemoveColumnsDefault(ptr.Pointer(), C.int(int32(column)), C.int(int32(count)), std_core.PointerFromQModelIndex(parent))) != 0
	}
	return false
}

//export callbackTxListModel687eda_RemoveRows
func callbackTxListModel687eda_RemoveRows(ptr unsafe.Pointer, row C.int, count C.int, parent unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "removeRows"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(int, int, *std_core.QModelIndex) bool)(signal))(int(int32(row)), int(int32(count)), std_core.NewQModelIndexFromPointer(parent)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListModelFromPointer(ptr).RemoveRowsDefault(int(int32(row)), int(int32(count)), std_core.NewQModelIndexFromPointer(parent)))))
}

func (ptr *TxListModel) RemoveRowsDefault(row int, count int, parent std_core.QModelIndex_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListModel687eda_RemoveRowsDefault(ptr.Pointer(), C.int(int32(row)), C.int(int32(count)), std_core.PointerFromQModelIndex(parent))) != 0
	}
	return false
}

//export callbackTxListModel687eda_ResetInternalData
func callbackTxListModel687eda_ResetInternalData(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "resetInternalData"); signal != nil {
		(*(*func())(signal))()
	} else {
		NewTxListModelFromPointer(ptr).ResetInternalDataDefault()
	}
}

func (ptr *TxListModel) ResetInternalDataDefault() {
	if ptr.Pointer() != nil {
		C.TxListModel687eda_ResetInternalDataDefault(ptr.Pointer())
	}
}

//export callbackTxListModel687eda_Revert
func callbackTxListModel687eda_Revert(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "revert"); signal != nil {
		(*(*func())(signal))()
	} else {
		NewTxListModelFromPointer(ptr).RevertDefault()
	}
}

func (ptr *TxListModel) RevertDefault() {
	if ptr.Pointer() != nil {
		C.TxListModel687eda_RevertDefault(ptr.Pointer())
	}
}

//export callbackTxListModel687eda_RoleNames
func callbackTxListModel687eda_RoleNames(ptr unsafe.Pointer) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "roleNames"); signal != nil {
		return func() unsafe.Pointer {
			tmpList := NewTxListModelFromPointer(NewTxListModelFromPointer(nil).__roleNames_newList())
			for k, v := range (*(*func() map[int]*std_core.QByteArray)(signal))() {
				tmpList.__roleNames_setList(k, v)
			}
			return tmpList.Pointer()
		}()
	}

	return func() unsafe.Pointer {
		tmpList := NewTxListModelFromPointer(NewTxListModelFromPointer(nil).__roleNames_newList())
		for k, v := range NewTxListModelFromPointer(ptr).RoleNamesDefault() {
			tmpList.__roleNames_setList(k, v)
		}
		return tmpList.Pointer()
	}()
}

func (ptr *TxListModel) RoleNamesDefault() map[int]*std_core.QByteArray {
	if ptr.Pointer() != nil {
		return func(l C.struct_Moc_PackedList) map[int]*std_core.QByteArray {
			out := make(map[int]*std_core.QByteArray, int(l.len))
			tmpList := NewTxListModelFromPointer(l.data)
			for i, v := range tmpList.__roleNames_keyList() {
				out[v] = tmpList.__roleNames_atList(v, i)
			}
			return out
		}(C.TxListModel687eda_RoleNamesDefault(ptr.Pointer()))
	}
	return make(map[int]*std_core.QByteArray, 0)
}

//export callbackTxListModel687eda_RowCount
func callbackTxListModel687eda_RowCount(ptr unsafe.Pointer, parent unsafe.Pointer) C.int {
	if signal := qt.GetSignal(ptr, "rowCount"); signal != nil {
		return C.int(int32((*(*func(*std_core.QModelIndex) int)(signal))(std_core.NewQModelIndexFromPointer(parent))))
	}

	return C.int(int32(NewTxListModelFromPointer(ptr).RowCountDefault(std_core.NewQModelIndexFromPointer(parent))))
}

func (ptr *TxListModel) RowCountDefault(parent std_core.QModelIndex_ITF) int {
	if ptr.Pointer() != nil {
		return int(int32(C.TxListModel687eda_RowCountDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(parent))))
	}
	return 0
}

//export callbackTxListModel687eda_RowsAboutToBeInserted
func callbackTxListModel687eda_RowsAboutToBeInserted(ptr unsafe.Pointer, parent unsafe.Pointer, start C.int, end C.int) {
	if signal := qt.GetSignal(ptr, "rowsAboutToBeInserted"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(start)), int(int32(end)))
	}

}

//export callbackTxListModel687eda_RowsAboutToBeMoved
func callbackTxListModel687eda_RowsAboutToBeMoved(ptr unsafe.Pointer, sourceParent unsafe.Pointer, sourceStart C.int, sourceEnd C.int, destinationParent unsafe.Pointer, destinationRow C.int) {
	if signal := qt.GetSignal(ptr, "rowsAboutToBeMoved"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int, *std_core.QModelIndex, int))(signal))(std_core.NewQModelIndexFromPointer(sourceParent), int(int32(sourceStart)), int(int32(sourceEnd)), std_core.NewQModelIndexFromPointer(destinationParent), int(int32(destinationRow)))
	}

}

//export callbackTxListModel687eda_RowsAboutToBeRemoved
func callbackTxListModel687eda_RowsAboutToBeRemoved(ptr unsafe.Pointer, parent unsafe.Pointer, first C.int, last C.int) {
	if signal := qt.GetSignal(ptr, "rowsAboutToBeRemoved"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(first)), int(int32(last)))
	}

}

//export callbackTxListModel687eda_RowsInserted
func callbackTxListModel687eda_RowsInserted(ptr unsafe.Pointer, parent unsafe.Pointer, first C.int, last C.int) {
	if signal := qt.GetSignal(ptr, "rowsInserted"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(first)), int(int32(last)))
	}

}

//export callbackTxListModel687eda_RowsMoved
func callbackTxListModel687eda_RowsMoved(ptr unsafe.Pointer, parent unsafe.Pointer, start C.int, end C.int, destination unsafe.Pointer, row C.int) {
	if signal := qt.GetSignal(ptr, "rowsMoved"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int, *std_core.QModelIndex, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(start)), int(int32(end)), std_core.NewQModelIndexFromPointer(destination), int(int32(row)))
	}

}

//export callbackTxListModel687eda_RowsRemoved
func callbackTxListModel687eda_RowsRemoved(ptr unsafe.Pointer, parent unsafe.Pointer, first C.int, last C.int) {
	if signal := qt.GetSignal(ptr, "rowsRemoved"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(first)), int(int32(last)))
	}

}

//export callbackTxListModel687eda_SetData
func callbackTxListModel687eda_SetData(ptr unsafe.Pointer, index unsafe.Pointer, value unsafe.Pointer, role C.int) C.char {
	if signal := qt.GetSignal(ptr, "setData"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QModelIndex, *std_core.QVariant, int) bool)(signal))(std_core.NewQModelIndexFromPointer(index), std_core.NewQVariantFromPointer(value), int(int32(role))))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListModelFromPointer(ptr).SetDataDefault(std_core.NewQModelIndexFromPointer(index), std_core.NewQVariantFromPointer(value), int(int32(role))))))
}

func (ptr *TxListModel) SetDataDefault(index std_core.QModelIndex_ITF, value std_core.QVariant_ITF, role int) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListModel687eda_SetDataDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(index), std_core.PointerFromQVariant(value), C.int(int32(role)))) != 0
	}
	return false
}

//export callbackTxListModel687eda_SetHeaderData
func callbackTxListModel687eda_SetHeaderData(ptr unsafe.Pointer, section C.int, orientation C.longlong, value unsafe.Pointer, role C.int) C.char {
	if signal := qt.GetSignal(ptr, "setHeaderData"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(int, std_core.Qt__Orientation, *std_core.QVariant, int) bool)(signal))(int(int32(section)), std_core.Qt__Orientation(orientation), std_core.NewQVariantFromPointer(value), int(int32(role))))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListModelFromPointer(ptr).SetHeaderDataDefault(int(int32(section)), std_core.Qt__Orientation(orientation), std_core.NewQVariantFromPointer(value), int(int32(role))))))
}

func (ptr *TxListModel) SetHeaderDataDefault(section int, orientation std_core.Qt__Orientation, value std_core.QVariant_ITF, role int) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListModel687eda_SetHeaderDataDefault(ptr.Pointer(), C.int(int32(section)), C.longlong(orientation), std_core.PointerFromQVariant(value), C.int(int32(role)))) != 0
	}
	return false
}

//export callbackTxListModel687eda_SetItemData
func callbackTxListModel687eda_SetItemData(ptr unsafe.Pointer, index unsafe.Pointer, roles C.struct_Moc_PackedList) C.char {
	if signal := qt.GetSignal(ptr, "setItemData"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QModelIndex, map[int]*std_core.QVariant) bool)(signal))(std_core.NewQModelIndexFromPointer(index), func(l C.struct_Moc_PackedList) map[int]*std_core.QVariant {
			out := make(map[int]*std_core.QVariant, int(l.len))
			tmpList := NewTxListModelFromPointer(l.data)
			for i, v := range tmpList.__setItemData_roles_keyList() {
				out[v] = tmpList.__setItemData_roles_atList(v, i)
			}
			return out
		}(roles)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListModelFromPointer(ptr).SetItemDataDefault(std_core.NewQModelIndexFromPointer(index), func(l C.struct_Moc_PackedList) map[int]*std_core.QVariant {
		out := make(map[int]*std_core.QVariant, int(l.len))
		tmpList := NewTxListModelFromPointer(l.data)
		for i, v := range tmpList.__setItemData_roles_keyList() {
			out[v] = tmpList.__setItemData_roles_atList(v, i)
		}
		return out
	}(roles)))))
}

func (ptr *TxListModel) SetItemDataDefault(index std_core.QModelIndex_ITF, roles map[int]*std_core.QVariant) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListModel687eda_SetItemDataDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(index), func() unsafe.Pointer {
			tmpList := NewTxListModelFromPointer(NewTxListModelFromPointer(nil).__setItemData_roles_newList())
			for k, v := range roles {
				tmpList.__setItemData_roles_setList(k, v)
			}
			return tmpList.Pointer()
		}())) != 0
	}
	return false
}

//export callbackTxListModel687eda_Sort
func callbackTxListModel687eda_Sort(ptr unsafe.Pointer, column C.int, order C.longlong) {
	if signal := qt.GetSignal(ptr, "sort"); signal != nil {
		(*(*func(int, std_core.Qt__SortOrder))(signal))(int(int32(column)), std_core.Qt__SortOrder(order))
	} else {
		NewTxListModelFromPointer(ptr).SortDefault(int(int32(column)), std_core.Qt__SortOrder(order))
	}
}

func (ptr *TxListModel) SortDefault(column int, order std_core.Qt__SortOrder) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda_SortDefault(ptr.Pointer(), C.int(int32(column)), C.longlong(order))
	}
}

//export callbackTxListModel687eda_Span
func callbackTxListModel687eda_Span(ptr unsafe.Pointer, index unsafe.Pointer) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "span"); signal != nil {
		return std_core.PointerFromQSize((*(*func(*std_core.QModelIndex) *std_core.QSize)(signal))(std_core.NewQModelIndexFromPointer(index)))
	}

	return std_core.PointerFromQSize(NewTxListModelFromPointer(ptr).SpanDefault(std_core.NewQModelIndexFromPointer(index)))
}

func (ptr *TxListModel) SpanDefault(index std_core.QModelIndex_ITF) *std_core.QSize {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQSizeFromPointer(C.TxListModel687eda_SpanDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(index)))
		qt.SetFinalizer(tmpValue, (*std_core.QSize).DestroyQSize)
		return tmpValue
	}
	return nil
}

//export callbackTxListModel687eda_Submit
func callbackTxListModel687eda_Submit(ptr unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "submit"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func() bool)(signal))())))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListModelFromPointer(ptr).SubmitDefault())))
}

func (ptr *TxListModel) SubmitDefault() bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListModel687eda_SubmitDefault(ptr.Pointer())) != 0
	}
	return false
}

//export callbackTxListModel687eda_SupportedDragActions
func callbackTxListModel687eda_SupportedDragActions(ptr unsafe.Pointer) C.longlong {
	if signal := qt.GetSignal(ptr, "supportedDragActions"); signal != nil {
		return C.longlong((*(*func() std_core.Qt__DropAction)(signal))())
	}

	return C.longlong(NewTxListModelFromPointer(ptr).SupportedDragActionsDefault())
}

func (ptr *TxListModel) SupportedDragActionsDefault() std_core.Qt__DropAction {
	if ptr.Pointer() != nil {
		return std_core.Qt__DropAction(C.TxListModel687eda_SupportedDragActionsDefault(ptr.Pointer()))
	}
	return 0
}

//export callbackTxListModel687eda_SupportedDropActions
func callbackTxListModel687eda_SupportedDropActions(ptr unsafe.Pointer) C.longlong {
	if signal := qt.GetSignal(ptr, "supportedDropActions"); signal != nil {
		return C.longlong((*(*func() std_core.Qt__DropAction)(signal))())
	}

	return C.longlong(NewTxListModelFromPointer(ptr).SupportedDropActionsDefault())
}

func (ptr *TxListModel) SupportedDropActionsDefault() std_core.Qt__DropAction {
	if ptr.Pointer() != nil {
		return std_core.Qt__DropAction(C.TxListModel687eda_SupportedDropActionsDefault(ptr.Pointer()))
	}
	return 0
}

//export callbackTxListModel687eda_ChildEvent
func callbackTxListModel687eda_ChildEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "childEvent"); signal != nil {
		(*(*func(*std_core.QChildEvent))(signal))(std_core.NewQChildEventFromPointer(event))
	} else {
		NewTxListModelFromPointer(ptr).ChildEventDefault(std_core.NewQChildEventFromPointer(event))
	}
}

func (ptr *TxListModel) ChildEventDefault(event std_core.QChildEvent_ITF) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda_ChildEventDefault(ptr.Pointer(), std_core.PointerFromQChildEvent(event))
	}
}

//export callbackTxListModel687eda_ConnectNotify
func callbackTxListModel687eda_ConnectNotify(ptr unsafe.Pointer, sign unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "connectNotify"); signal != nil {
		(*(*func(*std_core.QMetaMethod))(signal))(std_core.NewQMetaMethodFromPointer(sign))
	} else {
		NewTxListModelFromPointer(ptr).ConnectNotifyDefault(std_core.NewQMetaMethodFromPointer(sign))
	}
}

func (ptr *TxListModel) ConnectNotifyDefault(sign std_core.QMetaMethod_ITF) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda_ConnectNotifyDefault(ptr.Pointer(), std_core.PointerFromQMetaMethod(sign))
	}
}

//export callbackTxListModel687eda_CustomEvent
func callbackTxListModel687eda_CustomEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "customEvent"); signal != nil {
		(*(*func(*std_core.QEvent))(signal))(std_core.NewQEventFromPointer(event))
	} else {
		NewTxListModelFromPointer(ptr).CustomEventDefault(std_core.NewQEventFromPointer(event))
	}
}

func (ptr *TxListModel) CustomEventDefault(event std_core.QEvent_ITF) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda_CustomEventDefault(ptr.Pointer(), std_core.PointerFromQEvent(event))
	}
}

//export callbackTxListModel687eda_DeleteLater
func callbackTxListModel687eda_DeleteLater(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "deleteLater"); signal != nil {
		(*(*func())(signal))()
	} else {
		NewTxListModelFromPointer(ptr).DeleteLaterDefault()
	}
}

func (ptr *TxListModel) DeleteLaterDefault() {
	if ptr.Pointer() != nil {

		qt.SetFinalizer(ptr, nil)
		C.TxListModel687eda_DeleteLaterDefault(ptr.Pointer())
	}
}

//export callbackTxListModel687eda_Destroyed
func callbackTxListModel687eda_Destroyed(ptr unsafe.Pointer, obj unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "destroyed"); signal != nil {
		(*(*func(*std_core.QObject))(signal))(std_core.NewQObjectFromPointer(obj))
	}
	qt.Unregister(ptr)

}

//export callbackTxListModel687eda_DisconnectNotify
func callbackTxListModel687eda_DisconnectNotify(ptr unsafe.Pointer, sign unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "disconnectNotify"); signal != nil {
		(*(*func(*std_core.QMetaMethod))(signal))(std_core.NewQMetaMethodFromPointer(sign))
	} else {
		NewTxListModelFromPointer(ptr).DisconnectNotifyDefault(std_core.NewQMetaMethodFromPointer(sign))
	}
}

func (ptr *TxListModel) DisconnectNotifyDefault(sign std_core.QMetaMethod_ITF) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda_DisconnectNotifyDefault(ptr.Pointer(), std_core.PointerFromQMetaMethod(sign))
	}
}

//export callbackTxListModel687eda_Event
func callbackTxListModel687eda_Event(ptr unsafe.Pointer, e unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "event"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QEvent) bool)(signal))(std_core.NewQEventFromPointer(e)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListModelFromPointer(ptr).EventDefault(std_core.NewQEventFromPointer(e)))))
}

func (ptr *TxListModel) EventDefault(e std_core.QEvent_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListModel687eda_EventDefault(ptr.Pointer(), std_core.PointerFromQEvent(e))) != 0
	}
	return false
}

//export callbackTxListModel687eda_EventFilter
func callbackTxListModel687eda_EventFilter(ptr unsafe.Pointer, watched unsafe.Pointer, event unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "eventFilter"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QObject, *std_core.QEvent) bool)(signal))(std_core.NewQObjectFromPointer(watched), std_core.NewQEventFromPointer(event)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewTxListModelFromPointer(ptr).EventFilterDefault(std_core.NewQObjectFromPointer(watched), std_core.NewQEventFromPointer(event)))))
}

func (ptr *TxListModel) EventFilterDefault(watched std_core.QObject_ITF, event std_core.QEvent_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.TxListModel687eda_EventFilterDefault(ptr.Pointer(), std_core.PointerFromQObject(watched), std_core.PointerFromQEvent(event))) != 0
	}
	return false
}

//export callbackTxListModel687eda_ObjectNameChanged
func callbackTxListModel687eda_ObjectNameChanged(ptr unsafe.Pointer, objectName C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "objectNameChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(objectName))
	}

}

//export callbackTxListModel687eda_TimerEvent
func callbackTxListModel687eda_TimerEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "timerEvent"); signal != nil {
		(*(*func(*std_core.QTimerEvent))(signal))(std_core.NewQTimerEventFromPointer(event))
	} else {
		NewTxListModelFromPointer(ptr).TimerEventDefault(std_core.NewQTimerEventFromPointer(event))
	}
}

func (ptr *TxListModel) TimerEventDefault(event std_core.QTimerEvent_ITF) {
	if ptr.Pointer() != nil {
		C.TxListModel687eda_TimerEventDefault(ptr.Pointer(), std_core.PointerFromQTimerEvent(event))
	}
}

type ApproveListingCtx_ITF interface {
	std_core.QObject_ITF
	ApproveListingCtx_PTR() *ApproveListingCtx
}

func (ptr *ApproveListingCtx) ApproveListingCtx_PTR() *ApproveListingCtx {
	return ptr
}

func (ptr *ApproveListingCtx) Pointer() unsafe.Pointer {
	if ptr != nil {
		return ptr.QObject_PTR().Pointer()
	}
	return nil
}

func (ptr *ApproveListingCtx) SetPointer(p unsafe.Pointer) {
	if ptr != nil {
		ptr.QObject_PTR().SetPointer(p)
	}
}

func PointerFromApproveListingCtx(ptr ApproveListingCtx_ITF) unsafe.Pointer {
	if ptr != nil {
		return ptr.ApproveListingCtx_PTR().Pointer()
	}
	return nil
}

func NewApproveListingCtxFromPointer(ptr unsafe.Pointer) (n *ApproveListingCtx) {
	if gPtr, ok := qt.Receive(ptr); !ok {
		n = new(ApproveListingCtx)
		n.SetPointer(ptr)
	} else {
		switch deduced := gPtr.(type) {
		case *ApproveListingCtx:
			n = deduced

		case *std_core.QObject:
			n = &ApproveListingCtx{QObject: *deduced}

		default:
			n = new(ApproveListingCtx)
			n.SetPointer(ptr)
		}
	}
	return
}
func (this *ApproveListingCtx) Init() { this.init() }

//export callbackApproveListingCtx687eda_Constructor
func callbackApproveListingCtx687eda_Constructor(ptr unsafe.Pointer) {
	this := NewApproveListingCtxFromPointer(ptr)
	qt.Register(ptr, this)
	this.ConnectBack(this.back)
	this.ConnectApprove(this.approve)
	this.ConnectReject(this.reject)
	this.ConnectOnCheckStateChanged(this.onCheckStateChanged)
	this.init()
}

//export callbackApproveListingCtx687eda_Back
func callbackApproveListingCtx687eda_Back(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "back"); signal != nil {
		(*(*func())(signal))()
	}

}

func (ptr *ApproveListingCtx) ConnectBack(f func()) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "back") {
			C.ApproveListingCtx687eda_ConnectBack(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "back")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "back"); signal != nil {
			f := func() {
				(*(*func())(signal))()
				f()
			}
			qt.ConnectSignal(ptr.Pointer(), "back", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "back", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveListingCtx) DisconnectBack() {
	if ptr.Pointer() != nil {
		C.ApproveListingCtx687eda_DisconnectBack(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "back")
	}
}

func (ptr *ApproveListingCtx) Back() {
	if ptr.Pointer() != nil {
		C.ApproveListingCtx687eda_Back(ptr.Pointer())
	}
}

//export callbackApproveListingCtx687eda_Approve
func callbackApproveListingCtx687eda_Approve(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "approve"); signal != nil {
		(*(*func())(signal))()
	}

}

func (ptr *ApproveListingCtx) ConnectApprove(f func()) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "approve") {
			C.ApproveListingCtx687eda_ConnectApprove(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "approve")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "approve"); signal != nil {
			f := func() {
				(*(*func())(signal))()
				f()
			}
			qt.ConnectSignal(ptr.Pointer(), "approve", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "approve", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveListingCtx) DisconnectApprove() {
	if ptr.Pointer() != nil {
		C.ApproveListingCtx687eda_DisconnectApprove(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "approve")
	}
}

func (ptr *ApproveListingCtx) Approve() {
	if ptr.Pointer() != nil {
		C.ApproveListingCtx687eda_Approve(ptr.Pointer())
	}
}

//export callbackApproveListingCtx687eda_Reject
func callbackApproveListingCtx687eda_Reject(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "reject"); signal != nil {
		(*(*func())(signal))()
	}

}

func (ptr *ApproveListingCtx) ConnectReject(f func()) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "reject") {
			C.ApproveListingCtx687eda_ConnectReject(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "reject")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "reject"); signal != nil {
			f := func() {
				(*(*func())(signal))()
				f()
			}
			qt.ConnectSignal(ptr.Pointer(), "reject", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "reject", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveListingCtx) DisconnectReject() {
	if ptr.Pointer() != nil {
		C.ApproveListingCtx687eda_DisconnectReject(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "reject")
	}
}

func (ptr *ApproveListingCtx) Reject() {
	if ptr.Pointer() != nil {
		C.ApproveListingCtx687eda_Reject(ptr.Pointer())
	}
}

//export callbackApproveListingCtx687eda_OnCheckStateChanged
func callbackApproveListingCtx687eda_OnCheckStateChanged(ptr unsafe.Pointer, i C.int, checked C.char) {
	if signal := qt.GetSignal(ptr, "onCheckStateChanged"); signal != nil {
		(*(*func(int, bool))(signal))(int(int32(i)), int8(checked) != 0)
	}

}

func (ptr *ApproveListingCtx) ConnectOnCheckStateChanged(f func(i int, checked bool)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "onCheckStateChanged") {
			C.ApproveListingCtx687eda_ConnectOnCheckStateChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "onCheckStateChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "onCheckStateChanged"); signal != nil {
			f := func(i int, checked bool) {
				(*(*func(int, bool))(signal))(i, checked)
				f(i, checked)
			}
			qt.ConnectSignal(ptr.Pointer(), "onCheckStateChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "onCheckStateChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveListingCtx) DisconnectOnCheckStateChanged() {
	if ptr.Pointer() != nil {
		C.ApproveListingCtx687eda_DisconnectOnCheckStateChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "onCheckStateChanged")
	}
}

func (ptr *ApproveListingCtx) OnCheckStateChanged(i int, checked bool) {
	if ptr.Pointer() != nil {
		C.ApproveListingCtx687eda_OnCheckStateChanged(ptr.Pointer(), C.int(int32(i)), C.char(int8(qt.GoBoolToInt(checked))))
	}
}

//export callbackApproveListingCtx687eda_TriggerUpdate
func callbackApproveListingCtx687eda_TriggerUpdate(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "triggerUpdate"); signal != nil {
		(*(*func())(signal))()
	}

}

func (ptr *ApproveListingCtx) ConnectTriggerUpdate(f func()) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "triggerUpdate") {
			C.ApproveListingCtx687eda_ConnectTriggerUpdate(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "triggerUpdate")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "triggerUpdate"); signal != nil {
			f := func() {
				(*(*func())(signal))()
				f()
			}
			qt.ConnectSignal(ptr.Pointer(), "triggerUpdate", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "triggerUpdate", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveListingCtx) DisconnectTriggerUpdate() {
	if ptr.Pointer() != nil {
		C.ApproveListingCtx687eda_DisconnectTriggerUpdate(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "triggerUpdate")
	}
}

func (ptr *ApproveListingCtx) TriggerUpdate() {
	if ptr.Pointer() != nil {
		C.ApproveListingCtx687eda_TriggerUpdate(ptr.Pointer())
	}
}

//export callbackApproveListingCtx687eda_Remote
func callbackApproveListingCtx687eda_Remote(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "remote"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveListingCtxFromPointer(ptr).RemoteDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveListingCtx) ConnectRemote(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "remote"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "remote", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "remote", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveListingCtx) DisconnectRemote() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "remote")
	}
}

func (ptr *ApproveListingCtx) Remote() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveListingCtx687eda_Remote(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveListingCtx) RemoteDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveListingCtx687eda_RemoteDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveListingCtx687eda_SetRemote
func callbackApproveListingCtx687eda_SetRemote(ptr unsafe.Pointer, remote C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setRemote"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(remote))
	} else {
		NewApproveListingCtxFromPointer(ptr).SetRemoteDefault(cGoUnpackString(remote))
	}
}

func (ptr *ApproveListingCtx) ConnectSetRemote(f func(remote string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setRemote"); signal != nil {
			f := func(remote string) {
				(*(*func(string))(signal))(remote)
				f(remote)
			}
			qt.ConnectSignal(ptr.Pointer(), "setRemote", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setRemote", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveListingCtx) DisconnectSetRemote() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setRemote")
	}
}

func (ptr *ApproveListingCtx) SetRemote(remote string) {
	if ptr.Pointer() != nil {
		var remoteC *C.char
		if remote != "" {
			remoteC = C.CString(remote)
			defer C.free(unsafe.Pointer(remoteC))
		}
		C.ApproveListingCtx687eda_SetRemote(ptr.Pointer(), C.struct_Moc_PackedString{data: remoteC, len: C.longlong(len(remote))})
	}
}

func (ptr *ApproveListingCtx) SetRemoteDefault(remote string) {
	if ptr.Pointer() != nil {
		var remoteC *C.char
		if remote != "" {
			remoteC = C.CString(remote)
			defer C.free(unsafe.Pointer(remoteC))
		}
		C.ApproveListingCtx687eda_SetRemoteDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: remoteC, len: C.longlong(len(remote))})
	}
}

//export callbackApproveListingCtx687eda_RemoteChanged
func callbackApproveListingCtx687eda_RemoteChanged(ptr unsafe.Pointer, remote C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "remoteChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(remote))
	}

}

func (ptr *ApproveListingCtx) ConnectRemoteChanged(f func(remote string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "remoteChanged") {
			C.ApproveListingCtx687eda_ConnectRemoteChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "remoteChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "remoteChanged"); signal != nil {
			f := func(remote string) {
				(*(*func(string))(signal))(remote)
				f(remote)
			}
			qt.ConnectSignal(ptr.Pointer(), "remoteChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "remoteChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveListingCtx) DisconnectRemoteChanged() {
	if ptr.Pointer() != nil {
		C.ApproveListingCtx687eda_DisconnectRemoteChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "remoteChanged")
	}
}

func (ptr *ApproveListingCtx) RemoteChanged(remote string) {
	if ptr.Pointer() != nil {
		var remoteC *C.char
		if remote != "" {
			remoteC = C.CString(remote)
			defer C.free(unsafe.Pointer(remoteC))
		}
		C.ApproveListingCtx687eda_RemoteChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: remoteC, len: C.longlong(len(remote))})
	}
}

//export callbackApproveListingCtx687eda_Transport
func callbackApproveListingCtx687eda_Transport(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "transport"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveListingCtxFromPointer(ptr).TransportDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveListingCtx) ConnectTransport(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "transport"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "transport", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "transport", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveListingCtx) DisconnectTransport() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "transport")
	}
}

func (ptr *ApproveListingCtx) Transport() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveListingCtx687eda_Transport(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveListingCtx) TransportDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveListingCtx687eda_TransportDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveListingCtx687eda_SetTransport
func callbackApproveListingCtx687eda_SetTransport(ptr unsafe.Pointer, transport C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setTransport"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(transport))
	} else {
		NewApproveListingCtxFromPointer(ptr).SetTransportDefault(cGoUnpackString(transport))
	}
}

func (ptr *ApproveListingCtx) ConnectSetTransport(f func(transport string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setTransport"); signal != nil {
			f := func(transport string) {
				(*(*func(string))(signal))(transport)
				f(transport)
			}
			qt.ConnectSignal(ptr.Pointer(), "setTransport", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setTransport", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveListingCtx) DisconnectSetTransport() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setTransport")
	}
}

func (ptr *ApproveListingCtx) SetTransport(transport string) {
	if ptr.Pointer() != nil {
		var transportC *C.char
		if transport != "" {
			transportC = C.CString(transport)
			defer C.free(unsafe.Pointer(transportC))
		}
		C.ApproveListingCtx687eda_SetTransport(ptr.Pointer(), C.struct_Moc_PackedString{data: transportC, len: C.longlong(len(transport))})
	}
}

func (ptr *ApproveListingCtx) SetTransportDefault(transport string) {
	if ptr.Pointer() != nil {
		var transportC *C.char
		if transport != "" {
			transportC = C.CString(transport)
			defer C.free(unsafe.Pointer(transportC))
		}
		C.ApproveListingCtx687eda_SetTransportDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: transportC, len: C.longlong(len(transport))})
	}
}

//export callbackApproveListingCtx687eda_TransportChanged
func callbackApproveListingCtx687eda_TransportChanged(ptr unsafe.Pointer, transport C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "transportChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(transport))
	}

}

func (ptr *ApproveListingCtx) ConnectTransportChanged(f func(transport string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "transportChanged") {
			C.ApproveListingCtx687eda_ConnectTransportChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "transportChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "transportChanged"); signal != nil {
			f := func(transport string) {
				(*(*func(string))(signal))(transport)
				f(transport)
			}
			qt.ConnectSignal(ptr.Pointer(), "transportChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "transportChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveListingCtx) DisconnectTransportChanged() {
	if ptr.Pointer() != nil {
		C.ApproveListingCtx687eda_DisconnectTransportChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "transportChanged")
	}
}

func (ptr *ApproveListingCtx) TransportChanged(transport string) {
	if ptr.Pointer() != nil {
		var transportC *C.char
		if transport != "" {
			transportC = C.CString(transport)
			defer C.free(unsafe.Pointer(transportC))
		}
		C.ApproveListingCtx687eda_TransportChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: transportC, len: C.longlong(len(transport))})
	}
}

//export callbackApproveListingCtx687eda_Endpoint
func callbackApproveListingCtx687eda_Endpoint(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "endpoint"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveListingCtxFromPointer(ptr).EndpointDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveListingCtx) ConnectEndpoint(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "endpoint"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "endpoint", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "endpoint", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveListingCtx) DisconnectEndpoint() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "endpoint")
	}
}

func (ptr *ApproveListingCtx) Endpoint() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveListingCtx687eda_Endpoint(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveListingCtx) EndpointDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveListingCtx687eda_EndpointDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveListingCtx687eda_SetEndpoint
func callbackApproveListingCtx687eda_SetEndpoint(ptr unsafe.Pointer, endpoint C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setEndpoint"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(endpoint))
	} else {
		NewApproveListingCtxFromPointer(ptr).SetEndpointDefault(cGoUnpackString(endpoint))
	}
}

func (ptr *ApproveListingCtx) ConnectSetEndpoint(f func(endpoint string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setEndpoint"); signal != nil {
			f := func(endpoint string) {
				(*(*func(string))(signal))(endpoint)
				f(endpoint)
			}
			qt.ConnectSignal(ptr.Pointer(), "setEndpoint", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setEndpoint", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveListingCtx) DisconnectSetEndpoint() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setEndpoint")
	}
}

func (ptr *ApproveListingCtx) SetEndpoint(endpoint string) {
	if ptr.Pointer() != nil {
		var endpointC *C.char
		if endpoint != "" {
			endpointC = C.CString(endpoint)
			defer C.free(unsafe.Pointer(endpointC))
		}
		C.ApproveListingCtx687eda_SetEndpoint(ptr.Pointer(), C.struct_Moc_PackedString{data: endpointC, len: C.longlong(len(endpoint))})
	}
}

func (ptr *ApproveListingCtx) SetEndpointDefault(endpoint string) {
	if ptr.Pointer() != nil {
		var endpointC *C.char
		if endpoint != "" {
			endpointC = C.CString(endpoint)
			defer C.free(unsafe.Pointer(endpointC))
		}
		C.ApproveListingCtx687eda_SetEndpointDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: endpointC, len: C.longlong(len(endpoint))})
	}
}

//export callbackApproveListingCtx687eda_EndpointChanged
func callbackApproveListingCtx687eda_EndpointChanged(ptr unsafe.Pointer, endpoint C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "endpointChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(endpoint))
	}

}

func (ptr *ApproveListingCtx) ConnectEndpointChanged(f func(endpoint string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "endpointChanged") {
			C.ApproveListingCtx687eda_ConnectEndpointChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "endpointChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "endpointChanged"); signal != nil {
			f := func(endpoint string) {
				(*(*func(string))(signal))(endpoint)
				f(endpoint)
			}
			qt.ConnectSignal(ptr.Pointer(), "endpointChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "endpointChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveListingCtx) DisconnectEndpointChanged() {
	if ptr.Pointer() != nil {
		C.ApproveListingCtx687eda_DisconnectEndpointChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "endpointChanged")
	}
}

func (ptr *ApproveListingCtx) EndpointChanged(endpoint string) {
	if ptr.Pointer() != nil {
		var endpointC *C.char
		if endpoint != "" {
			endpointC = C.CString(endpoint)
			defer C.free(unsafe.Pointer(endpointC))
		}
		C.ApproveListingCtx687eda_EndpointChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: endpointC, len: C.longlong(len(endpoint))})
	}
}

//export callbackApproveListingCtx687eda_From
func callbackApproveListingCtx687eda_From(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "from"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveListingCtxFromPointer(ptr).FromDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveListingCtx) ConnectFrom(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "from"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "from", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "from", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveListingCtx) DisconnectFrom() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "from")
	}
}

func (ptr *ApproveListingCtx) From() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveListingCtx687eda_From(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveListingCtx) FromDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveListingCtx687eda_FromDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveListingCtx687eda_SetFrom
func callbackApproveListingCtx687eda_SetFrom(ptr unsafe.Pointer, from C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setFrom"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(from))
	} else {
		NewApproveListingCtxFromPointer(ptr).SetFromDefault(cGoUnpackString(from))
	}
}

func (ptr *ApproveListingCtx) ConnectSetFrom(f func(from string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setFrom"); signal != nil {
			f := func(from string) {
				(*(*func(string))(signal))(from)
				f(from)
			}
			qt.ConnectSignal(ptr.Pointer(), "setFrom", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setFrom", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveListingCtx) DisconnectSetFrom() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setFrom")
	}
}

func (ptr *ApproveListingCtx) SetFrom(from string) {
	if ptr.Pointer() != nil {
		var fromC *C.char
		if from != "" {
			fromC = C.CString(from)
			defer C.free(unsafe.Pointer(fromC))
		}
		C.ApproveListingCtx687eda_SetFrom(ptr.Pointer(), C.struct_Moc_PackedString{data: fromC, len: C.longlong(len(from))})
	}
}

func (ptr *ApproveListingCtx) SetFromDefault(from string) {
	if ptr.Pointer() != nil {
		var fromC *C.char
		if from != "" {
			fromC = C.CString(from)
			defer C.free(unsafe.Pointer(fromC))
		}
		C.ApproveListingCtx687eda_SetFromDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: fromC, len: C.longlong(len(from))})
	}
}

//export callbackApproveListingCtx687eda_FromChanged
func callbackApproveListingCtx687eda_FromChanged(ptr unsafe.Pointer, from C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "fromChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(from))
	}

}

func (ptr *ApproveListingCtx) ConnectFromChanged(f func(from string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "fromChanged") {
			C.ApproveListingCtx687eda_ConnectFromChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "fromChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "fromChanged"); signal != nil {
			f := func(from string) {
				(*(*func(string))(signal))(from)
				f(from)
			}
			qt.ConnectSignal(ptr.Pointer(), "fromChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "fromChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveListingCtx) DisconnectFromChanged() {
	if ptr.Pointer() != nil {
		C.ApproveListingCtx687eda_DisconnectFromChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "fromChanged")
	}
}

func (ptr *ApproveListingCtx) FromChanged(from string) {
	if ptr.Pointer() != nil {
		var fromC *C.char
		if from != "" {
			fromC = C.CString(from)
			defer C.free(unsafe.Pointer(fromC))
		}
		C.ApproveListingCtx687eda_FromChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: fromC, len: C.longlong(len(from))})
	}
}

//export callbackApproveListingCtx687eda_Message
func callbackApproveListingCtx687eda_Message(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "message"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveListingCtxFromPointer(ptr).MessageDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveListingCtx) ConnectMessage(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "message"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "message", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "message", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveListingCtx) DisconnectMessage() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "message")
	}
}

func (ptr *ApproveListingCtx) Message() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveListingCtx687eda_Message(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveListingCtx) MessageDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveListingCtx687eda_MessageDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveListingCtx687eda_SetMessage
func callbackApproveListingCtx687eda_SetMessage(ptr unsafe.Pointer, message C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setMessage"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(message))
	} else {
		NewApproveListingCtxFromPointer(ptr).SetMessageDefault(cGoUnpackString(message))
	}
}

func (ptr *ApproveListingCtx) ConnectSetMessage(f func(message string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setMessage"); signal != nil {
			f := func(message string) {
				(*(*func(string))(signal))(message)
				f(message)
			}
			qt.ConnectSignal(ptr.Pointer(), "setMessage", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setMessage", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveListingCtx) DisconnectSetMessage() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setMessage")
	}
}

func (ptr *ApproveListingCtx) SetMessage(message string) {
	if ptr.Pointer() != nil {
		var messageC *C.char
		if message != "" {
			messageC = C.CString(message)
			defer C.free(unsafe.Pointer(messageC))
		}
		C.ApproveListingCtx687eda_SetMessage(ptr.Pointer(), C.struct_Moc_PackedString{data: messageC, len: C.longlong(len(message))})
	}
}

func (ptr *ApproveListingCtx) SetMessageDefault(message string) {
	if ptr.Pointer() != nil {
		var messageC *C.char
		if message != "" {
			messageC = C.CString(message)
			defer C.free(unsafe.Pointer(messageC))
		}
		C.ApproveListingCtx687eda_SetMessageDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: messageC, len: C.longlong(len(message))})
	}
}

//export callbackApproveListingCtx687eda_MessageChanged
func callbackApproveListingCtx687eda_MessageChanged(ptr unsafe.Pointer, message C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "messageChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(message))
	}

}

func (ptr *ApproveListingCtx) ConnectMessageChanged(f func(message string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "messageChanged") {
			C.ApproveListingCtx687eda_ConnectMessageChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "messageChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "messageChanged"); signal != nil {
			f := func(message string) {
				(*(*func(string))(signal))(message)
				f(message)
			}
			qt.ConnectSignal(ptr.Pointer(), "messageChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "messageChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveListingCtx) DisconnectMessageChanged() {
	if ptr.Pointer() != nil {
		C.ApproveListingCtx687eda_DisconnectMessageChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "messageChanged")
	}
}

func (ptr *ApproveListingCtx) MessageChanged(message string) {
	if ptr.Pointer() != nil {
		var messageC *C.char
		if message != "" {
			messageC = C.CString(message)
			defer C.free(unsafe.Pointer(messageC))
		}
		C.ApproveListingCtx687eda_MessageChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: messageC, len: C.longlong(len(message))})
	}
}

//export callbackApproveListingCtx687eda_RawData
func callbackApproveListingCtx687eda_RawData(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "rawData"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveListingCtxFromPointer(ptr).RawDataDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveListingCtx) ConnectRawData(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "rawData"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "rawData", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "rawData", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveListingCtx) DisconnectRawData() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "rawData")
	}
}

func (ptr *ApproveListingCtx) RawData() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveListingCtx687eda_RawData(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveListingCtx) RawDataDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveListingCtx687eda_RawDataDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveListingCtx687eda_SetRawData
func callbackApproveListingCtx687eda_SetRawData(ptr unsafe.Pointer, rawData C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setRawData"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(rawData))
	} else {
		NewApproveListingCtxFromPointer(ptr).SetRawDataDefault(cGoUnpackString(rawData))
	}
}

func (ptr *ApproveListingCtx) ConnectSetRawData(f func(rawData string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setRawData"); signal != nil {
			f := func(rawData string) {
				(*(*func(string))(signal))(rawData)
				f(rawData)
			}
			qt.ConnectSignal(ptr.Pointer(), "setRawData", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setRawData", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveListingCtx) DisconnectSetRawData() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setRawData")
	}
}

func (ptr *ApproveListingCtx) SetRawData(rawData string) {
	if ptr.Pointer() != nil {
		var rawDataC *C.char
		if rawData != "" {
			rawDataC = C.CString(rawData)
			defer C.free(unsafe.Pointer(rawDataC))
		}
		C.ApproveListingCtx687eda_SetRawData(ptr.Pointer(), C.struct_Moc_PackedString{data: rawDataC, len: C.longlong(len(rawData))})
	}
}

func (ptr *ApproveListingCtx) SetRawDataDefault(rawData string) {
	if ptr.Pointer() != nil {
		var rawDataC *C.char
		if rawData != "" {
			rawDataC = C.CString(rawData)
			defer C.free(unsafe.Pointer(rawDataC))
		}
		C.ApproveListingCtx687eda_SetRawDataDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: rawDataC, len: C.longlong(len(rawData))})
	}
}

//export callbackApproveListingCtx687eda_RawDataChanged
func callbackApproveListingCtx687eda_RawDataChanged(ptr unsafe.Pointer, rawData C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "rawDataChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(rawData))
	}

}

func (ptr *ApproveListingCtx) ConnectRawDataChanged(f func(rawData string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "rawDataChanged") {
			C.ApproveListingCtx687eda_ConnectRawDataChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "rawDataChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "rawDataChanged"); signal != nil {
			f := func(rawData string) {
				(*(*func(string))(signal))(rawData)
				f(rawData)
			}
			qt.ConnectSignal(ptr.Pointer(), "rawDataChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "rawDataChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveListingCtx) DisconnectRawDataChanged() {
	if ptr.Pointer() != nil {
		C.ApproveListingCtx687eda_DisconnectRawDataChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "rawDataChanged")
	}
}

func (ptr *ApproveListingCtx) RawDataChanged(rawData string) {
	if ptr.Pointer() != nil {
		var rawDataC *C.char
		if rawData != "" {
			rawDataC = C.CString(rawData)
			defer C.free(unsafe.Pointer(rawDataC))
		}
		C.ApproveListingCtx687eda_RawDataChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: rawDataC, len: C.longlong(len(rawData))})
	}
}

//export callbackApproveListingCtx687eda_Hash
func callbackApproveListingCtx687eda_Hash(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "hash"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveListingCtxFromPointer(ptr).HashDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveListingCtx) ConnectHash(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "hash"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "hash", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "hash", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveListingCtx) DisconnectHash() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "hash")
	}
}

func (ptr *ApproveListingCtx) Hash() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveListingCtx687eda_Hash(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveListingCtx) HashDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveListingCtx687eda_HashDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveListingCtx687eda_SetHash
func callbackApproveListingCtx687eda_SetHash(ptr unsafe.Pointer, hash C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setHash"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(hash))
	} else {
		NewApproveListingCtxFromPointer(ptr).SetHashDefault(cGoUnpackString(hash))
	}
}

func (ptr *ApproveListingCtx) ConnectSetHash(f func(hash string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setHash"); signal != nil {
			f := func(hash string) {
				(*(*func(string))(signal))(hash)
				f(hash)
			}
			qt.ConnectSignal(ptr.Pointer(), "setHash", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setHash", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveListingCtx) DisconnectSetHash() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setHash")
	}
}

func (ptr *ApproveListingCtx) SetHash(hash string) {
	if ptr.Pointer() != nil {
		var hashC *C.char
		if hash != "" {
			hashC = C.CString(hash)
			defer C.free(unsafe.Pointer(hashC))
		}
		C.ApproveListingCtx687eda_SetHash(ptr.Pointer(), C.struct_Moc_PackedString{data: hashC, len: C.longlong(len(hash))})
	}
}

func (ptr *ApproveListingCtx) SetHashDefault(hash string) {
	if ptr.Pointer() != nil {
		var hashC *C.char
		if hash != "" {
			hashC = C.CString(hash)
			defer C.free(unsafe.Pointer(hashC))
		}
		C.ApproveListingCtx687eda_SetHashDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: hashC, len: C.longlong(len(hash))})
	}
}

//export callbackApproveListingCtx687eda_HashChanged
func callbackApproveListingCtx687eda_HashChanged(ptr unsafe.Pointer, hash C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "hashChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(hash))
	}

}

func (ptr *ApproveListingCtx) ConnectHashChanged(f func(hash string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "hashChanged") {
			C.ApproveListingCtx687eda_ConnectHashChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "hashChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "hashChanged"); signal != nil {
			f := func(hash string) {
				(*(*func(string))(signal))(hash)
				f(hash)
			}
			qt.ConnectSignal(ptr.Pointer(), "hashChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "hashChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveListingCtx) DisconnectHashChanged() {
	if ptr.Pointer() != nil {
		C.ApproveListingCtx687eda_DisconnectHashChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "hashChanged")
	}
}

func (ptr *ApproveListingCtx) HashChanged(hash string) {
	if ptr.Pointer() != nil {
		var hashC *C.char
		if hash != "" {
			hashC = C.CString(hash)
			defer C.free(unsafe.Pointer(hashC))
		}
		C.ApproveListingCtx687eda_HashChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: hashC, len: C.longlong(len(hash))})
	}
}

func ApproveListingCtx_QRegisterMetaType() int {
	return int(int32(C.ApproveListingCtx687eda_ApproveListingCtx687eda_QRegisterMetaType()))
}

func (ptr *ApproveListingCtx) QRegisterMetaType() int {
	return int(int32(C.ApproveListingCtx687eda_ApproveListingCtx687eda_QRegisterMetaType()))
}

func ApproveListingCtx_QRegisterMetaType2(typeName string) int {
	var typeNameC *C.char
	if typeName != "" {
		typeNameC = C.CString(typeName)
		defer C.free(unsafe.Pointer(typeNameC))
	}
	return int(int32(C.ApproveListingCtx687eda_ApproveListingCtx687eda_QRegisterMetaType2(typeNameC)))
}

func (ptr *ApproveListingCtx) QRegisterMetaType2(typeName string) int {
	var typeNameC *C.char
	if typeName != "" {
		typeNameC = C.CString(typeName)
		defer C.free(unsafe.Pointer(typeNameC))
	}
	return int(int32(C.ApproveListingCtx687eda_ApproveListingCtx687eda_QRegisterMetaType2(typeNameC)))
}

func ApproveListingCtx_QmlRegisterType() int {
	return int(int32(C.ApproveListingCtx687eda_ApproveListingCtx687eda_QmlRegisterType()))
}

func (ptr *ApproveListingCtx) QmlRegisterType() int {
	return int(int32(C.ApproveListingCtx687eda_ApproveListingCtx687eda_QmlRegisterType()))
}

func ApproveListingCtx_QmlRegisterType2(uri string, versionMajor int, versionMinor int, qmlName string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	return int(int32(C.ApproveListingCtx687eda_ApproveListingCtx687eda_QmlRegisterType2(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC)))
}

func (ptr *ApproveListingCtx) QmlRegisterType2(uri string, versionMajor int, versionMinor int, qmlName string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	return int(int32(C.ApproveListingCtx687eda_ApproveListingCtx687eda_QmlRegisterType2(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC)))
}

func ApproveListingCtx_QmlRegisterUncreatableType(uri string, versionMajor int, versionMinor int, qmlName string, message string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	var messageC *C.char
	if message != "" {
		messageC = C.CString(message)
		defer C.free(unsafe.Pointer(messageC))
	}
	return int(int32(C.ApproveListingCtx687eda_ApproveListingCtx687eda_QmlRegisterUncreatableType(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC, C.struct_Moc_PackedString{data: messageC, len: C.longlong(len(message))})))
}

func (ptr *ApproveListingCtx) QmlRegisterUncreatableType(uri string, versionMajor int, versionMinor int, qmlName string, message string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	var messageC *C.char
	if message != "" {
		messageC = C.CString(message)
		defer C.free(unsafe.Pointer(messageC))
	}
	return int(int32(C.ApproveListingCtx687eda_ApproveListingCtx687eda_QmlRegisterUncreatableType(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC, C.struct_Moc_PackedString{data: messageC, len: C.longlong(len(message))})))
}

func (ptr *ApproveListingCtx) __children_atList(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQObjectFromPointer(C.ApproveListingCtx687eda___children_atList(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *ApproveListingCtx) __children_setList(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveListingCtx687eda___children_setList(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *ApproveListingCtx) __children_newList() unsafe.Pointer {
	return C.ApproveListingCtx687eda___children_newList(ptr.Pointer())
}

func (ptr *ApproveListingCtx) __dynamicPropertyNames_atList(i int) *std_core.QByteArray {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQByteArrayFromPointer(C.ApproveListingCtx687eda___dynamicPropertyNames_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QByteArray).DestroyQByteArray)
		return tmpValue
	}
	return nil
}

func (ptr *ApproveListingCtx) __dynamicPropertyNames_setList(i std_core.QByteArray_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveListingCtx687eda___dynamicPropertyNames_setList(ptr.Pointer(), std_core.PointerFromQByteArray(i))
	}
}

func (ptr *ApproveListingCtx) __dynamicPropertyNames_newList() unsafe.Pointer {
	return C.ApproveListingCtx687eda___dynamicPropertyNames_newList(ptr.Pointer())
}

func (ptr *ApproveListingCtx) __findChildren_atList(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQObjectFromPointer(C.ApproveListingCtx687eda___findChildren_atList(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *ApproveListingCtx) __findChildren_setList(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveListingCtx687eda___findChildren_setList(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *ApproveListingCtx) __findChildren_newList() unsafe.Pointer {
	return C.ApproveListingCtx687eda___findChildren_newList(ptr.Pointer())
}

func (ptr *ApproveListingCtx) __findChildren_atList3(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQObjectFromPointer(C.ApproveListingCtx687eda___findChildren_atList3(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *ApproveListingCtx) __findChildren_setList3(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveListingCtx687eda___findChildren_setList3(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *ApproveListingCtx) __findChildren_newList3() unsafe.Pointer {
	return C.ApproveListingCtx687eda___findChildren_newList3(ptr.Pointer())
}

func NewApproveListingCtx(parent std_core.QObject_ITF) *ApproveListingCtx {
	ApproveListingCtx_QRegisterMetaType()
	tmpValue := NewApproveListingCtxFromPointer(C.ApproveListingCtx687eda_NewApproveListingCtx(std_core.PointerFromQObject(parent)))
	if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
		tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
	}
	return tmpValue
}

//export callbackApproveListingCtx687eda_DestroyApproveListingCtx
func callbackApproveListingCtx687eda_DestroyApproveListingCtx(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "~ApproveListingCtx"); signal != nil {
		(*(*func())(signal))()
	} else {
		NewApproveListingCtxFromPointer(ptr).DestroyApproveListingCtxDefault()
	}
}

func (ptr *ApproveListingCtx) ConnectDestroyApproveListingCtx(f func()) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "~ApproveListingCtx"); signal != nil {
			f := func() {
				(*(*func())(signal))()
				f()
			}
			qt.ConnectSignal(ptr.Pointer(), "~ApproveListingCtx", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "~ApproveListingCtx", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveListingCtx) DisconnectDestroyApproveListingCtx() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "~ApproveListingCtx")
	}
}

func (ptr *ApproveListingCtx) DestroyApproveListingCtx() {
	if ptr.Pointer() != nil {

		qt.SetFinalizer(ptr, nil)
		C.ApproveListingCtx687eda_DestroyApproveListingCtx(ptr.Pointer())
		ptr.SetPointer(nil)
	}
}

func (ptr *ApproveListingCtx) DestroyApproveListingCtxDefault() {
	if ptr.Pointer() != nil {

		qt.SetFinalizer(ptr, nil)
		C.ApproveListingCtx687eda_DestroyApproveListingCtxDefault(ptr.Pointer())
		ptr.SetPointer(nil)
	}
}

//export callbackApproveListingCtx687eda_ChildEvent
func callbackApproveListingCtx687eda_ChildEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "childEvent"); signal != nil {
		(*(*func(*std_core.QChildEvent))(signal))(std_core.NewQChildEventFromPointer(event))
	} else {
		NewApproveListingCtxFromPointer(ptr).ChildEventDefault(std_core.NewQChildEventFromPointer(event))
	}
}

func (ptr *ApproveListingCtx) ChildEventDefault(event std_core.QChildEvent_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveListingCtx687eda_ChildEventDefault(ptr.Pointer(), std_core.PointerFromQChildEvent(event))
	}
}

//export callbackApproveListingCtx687eda_ConnectNotify
func callbackApproveListingCtx687eda_ConnectNotify(ptr unsafe.Pointer, sign unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "connectNotify"); signal != nil {
		(*(*func(*std_core.QMetaMethod))(signal))(std_core.NewQMetaMethodFromPointer(sign))
	} else {
		NewApproveListingCtxFromPointer(ptr).ConnectNotifyDefault(std_core.NewQMetaMethodFromPointer(sign))
	}
}

func (ptr *ApproveListingCtx) ConnectNotifyDefault(sign std_core.QMetaMethod_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveListingCtx687eda_ConnectNotifyDefault(ptr.Pointer(), std_core.PointerFromQMetaMethod(sign))
	}
}

//export callbackApproveListingCtx687eda_CustomEvent
func callbackApproveListingCtx687eda_CustomEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "customEvent"); signal != nil {
		(*(*func(*std_core.QEvent))(signal))(std_core.NewQEventFromPointer(event))
	} else {
		NewApproveListingCtxFromPointer(ptr).CustomEventDefault(std_core.NewQEventFromPointer(event))
	}
}

func (ptr *ApproveListingCtx) CustomEventDefault(event std_core.QEvent_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveListingCtx687eda_CustomEventDefault(ptr.Pointer(), std_core.PointerFromQEvent(event))
	}
}

//export callbackApproveListingCtx687eda_DeleteLater
func callbackApproveListingCtx687eda_DeleteLater(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "deleteLater"); signal != nil {
		(*(*func())(signal))()
	} else {
		NewApproveListingCtxFromPointer(ptr).DeleteLaterDefault()
	}
}

func (ptr *ApproveListingCtx) DeleteLaterDefault() {
	if ptr.Pointer() != nil {

		qt.SetFinalizer(ptr, nil)
		C.ApproveListingCtx687eda_DeleteLaterDefault(ptr.Pointer())
	}
}

//export callbackApproveListingCtx687eda_Destroyed
func callbackApproveListingCtx687eda_Destroyed(ptr unsafe.Pointer, obj unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "destroyed"); signal != nil {
		(*(*func(*std_core.QObject))(signal))(std_core.NewQObjectFromPointer(obj))
	}
	qt.Unregister(ptr)

}

//export callbackApproveListingCtx687eda_DisconnectNotify
func callbackApproveListingCtx687eda_DisconnectNotify(ptr unsafe.Pointer, sign unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "disconnectNotify"); signal != nil {
		(*(*func(*std_core.QMetaMethod))(signal))(std_core.NewQMetaMethodFromPointer(sign))
	} else {
		NewApproveListingCtxFromPointer(ptr).DisconnectNotifyDefault(std_core.NewQMetaMethodFromPointer(sign))
	}
}

func (ptr *ApproveListingCtx) DisconnectNotifyDefault(sign std_core.QMetaMethod_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveListingCtx687eda_DisconnectNotifyDefault(ptr.Pointer(), std_core.PointerFromQMetaMethod(sign))
	}
}

//export callbackApproveListingCtx687eda_Event
func callbackApproveListingCtx687eda_Event(ptr unsafe.Pointer, e unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "event"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QEvent) bool)(signal))(std_core.NewQEventFromPointer(e)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewApproveListingCtxFromPointer(ptr).EventDefault(std_core.NewQEventFromPointer(e)))))
}

func (ptr *ApproveListingCtx) EventDefault(e std_core.QEvent_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.ApproveListingCtx687eda_EventDefault(ptr.Pointer(), std_core.PointerFromQEvent(e))) != 0
	}
	return false
}

//export callbackApproveListingCtx687eda_EventFilter
func callbackApproveListingCtx687eda_EventFilter(ptr unsafe.Pointer, watched unsafe.Pointer, event unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "eventFilter"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QObject, *std_core.QEvent) bool)(signal))(std_core.NewQObjectFromPointer(watched), std_core.NewQEventFromPointer(event)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewApproveListingCtxFromPointer(ptr).EventFilterDefault(std_core.NewQObjectFromPointer(watched), std_core.NewQEventFromPointer(event)))))
}

func (ptr *ApproveListingCtx) EventFilterDefault(watched std_core.QObject_ITF, event std_core.QEvent_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.ApproveListingCtx687eda_EventFilterDefault(ptr.Pointer(), std_core.PointerFromQObject(watched), std_core.PointerFromQEvent(event))) != 0
	}
	return false
}

//export callbackApproveListingCtx687eda_ObjectNameChanged
func callbackApproveListingCtx687eda_ObjectNameChanged(ptr unsafe.Pointer, objectName C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "objectNameChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(objectName))
	}

}

//export callbackApproveListingCtx687eda_TimerEvent
func callbackApproveListingCtx687eda_TimerEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "timerEvent"); signal != nil {
		(*(*func(*std_core.QTimerEvent))(signal))(std_core.NewQTimerEventFromPointer(event))
	} else {
		NewApproveListingCtxFromPointer(ptr).TimerEventDefault(std_core.NewQTimerEventFromPointer(event))
	}
}

func (ptr *ApproveListingCtx) TimerEventDefault(event std_core.QTimerEvent_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveListingCtx687eda_TimerEventDefault(ptr.Pointer(), std_core.PointerFromQTimerEvent(event))
	}
}

type ApproveNewAccountCtx_ITF interface {
	std_core.QObject_ITF
	ApproveNewAccountCtx_PTR() *ApproveNewAccountCtx
}

func (ptr *ApproveNewAccountCtx) ApproveNewAccountCtx_PTR() *ApproveNewAccountCtx {
	return ptr
}

func (ptr *ApproveNewAccountCtx) Pointer() unsafe.Pointer {
	if ptr != nil {
		return ptr.QObject_PTR().Pointer()
	}
	return nil
}

func (ptr *ApproveNewAccountCtx) SetPointer(p unsafe.Pointer) {
	if ptr != nil {
		ptr.QObject_PTR().SetPointer(p)
	}
}

func PointerFromApproveNewAccountCtx(ptr ApproveNewAccountCtx_ITF) unsafe.Pointer {
	if ptr != nil {
		return ptr.ApproveNewAccountCtx_PTR().Pointer()
	}
	return nil
}

func NewApproveNewAccountCtxFromPointer(ptr unsafe.Pointer) (n *ApproveNewAccountCtx) {
	if gPtr, ok := qt.Receive(ptr); !ok {
		n = new(ApproveNewAccountCtx)
		n.SetPointer(ptr)
	} else {
		switch deduced := gPtr.(type) {
		case *ApproveNewAccountCtx:
			n = deduced

		case *std_core.QObject:
			n = &ApproveNewAccountCtx{QObject: *deduced}

		default:
			n = new(ApproveNewAccountCtx)
			n.SetPointer(ptr)
		}
	}
	return
}
func (this *ApproveNewAccountCtx) Init() { this.init() }

//export callbackApproveNewAccountCtx687eda_Constructor
func callbackApproveNewAccountCtx687eda_Constructor(ptr unsafe.Pointer) {
	this := NewApproveNewAccountCtxFromPointer(ptr)
	qt.Register(ptr, this)
	this.ConnectClicked(this.clicked)
	this.ConnectBack(this.back)
	this.ConnectPasswordEdited(this.passwordEdited)
	this.ConnectConfirmPasswordEdited(this.confirmPasswordEdited)
	this.init()
}

//export callbackApproveNewAccountCtx687eda_Clicked
func callbackApproveNewAccountCtx687eda_Clicked(ptr unsafe.Pointer, b C.int) {
	if signal := qt.GetSignal(ptr, "clicked"); signal != nil {
		(*(*func(int))(signal))(int(int32(b)))
	}

}

func (ptr *ApproveNewAccountCtx) ConnectClicked(f func(b int)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "clicked") {
			C.ApproveNewAccountCtx687eda_ConnectClicked(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "clicked")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "clicked"); signal != nil {
			f := func(b int) {
				(*(*func(int))(signal))(b)
				f(b)
			}
			qt.ConnectSignal(ptr.Pointer(), "clicked", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "clicked", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveNewAccountCtx) DisconnectClicked() {
	if ptr.Pointer() != nil {
		C.ApproveNewAccountCtx687eda_DisconnectClicked(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "clicked")
	}
}

func (ptr *ApproveNewAccountCtx) Clicked(b int) {
	if ptr.Pointer() != nil {
		C.ApproveNewAccountCtx687eda_Clicked(ptr.Pointer(), C.int(int32(b)))
	}
}

//export callbackApproveNewAccountCtx687eda_Back
func callbackApproveNewAccountCtx687eda_Back(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "back"); signal != nil {
		(*(*func())(signal))()
	}

}

func (ptr *ApproveNewAccountCtx) ConnectBack(f func()) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "back") {
			C.ApproveNewAccountCtx687eda_ConnectBack(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "back")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "back"); signal != nil {
			f := func() {
				(*(*func())(signal))()
				f()
			}
			qt.ConnectSignal(ptr.Pointer(), "back", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "back", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveNewAccountCtx) DisconnectBack() {
	if ptr.Pointer() != nil {
		C.ApproveNewAccountCtx687eda_DisconnectBack(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "back")
	}
}

func (ptr *ApproveNewAccountCtx) Back() {
	if ptr.Pointer() != nil {
		C.ApproveNewAccountCtx687eda_Back(ptr.Pointer())
	}
}

//export callbackApproveNewAccountCtx687eda_PasswordEdited
func callbackApproveNewAccountCtx687eda_PasswordEdited(ptr unsafe.Pointer, b C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "passwordEdited"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(b))
	}

}

func (ptr *ApproveNewAccountCtx) ConnectPasswordEdited(f func(b string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "passwordEdited") {
			C.ApproveNewAccountCtx687eda_ConnectPasswordEdited(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "passwordEdited")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "passwordEdited"); signal != nil {
			f := func(b string) {
				(*(*func(string))(signal))(b)
				f(b)
			}
			qt.ConnectSignal(ptr.Pointer(), "passwordEdited", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "passwordEdited", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveNewAccountCtx) DisconnectPasswordEdited() {
	if ptr.Pointer() != nil {
		C.ApproveNewAccountCtx687eda_DisconnectPasswordEdited(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "passwordEdited")
	}
}

func (ptr *ApproveNewAccountCtx) PasswordEdited(b string) {
	if ptr.Pointer() != nil {
		var bC *C.char
		if b != "" {
			bC = C.CString(b)
			defer C.free(unsafe.Pointer(bC))
		}
		C.ApproveNewAccountCtx687eda_PasswordEdited(ptr.Pointer(), C.struct_Moc_PackedString{data: bC, len: C.longlong(len(b))})
	}
}

//export callbackApproveNewAccountCtx687eda_ConfirmPasswordEdited
func callbackApproveNewAccountCtx687eda_ConfirmPasswordEdited(ptr unsafe.Pointer, b C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "confirmPasswordEdited"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(b))
	}

}

func (ptr *ApproveNewAccountCtx) ConnectConfirmPasswordEdited(f func(b string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "confirmPasswordEdited") {
			C.ApproveNewAccountCtx687eda_ConnectConfirmPasswordEdited(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "confirmPasswordEdited")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "confirmPasswordEdited"); signal != nil {
			f := func(b string) {
				(*(*func(string))(signal))(b)
				f(b)
			}
			qt.ConnectSignal(ptr.Pointer(), "confirmPasswordEdited", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "confirmPasswordEdited", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveNewAccountCtx) DisconnectConfirmPasswordEdited() {
	if ptr.Pointer() != nil {
		C.ApproveNewAccountCtx687eda_DisconnectConfirmPasswordEdited(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "confirmPasswordEdited")
	}
}

func (ptr *ApproveNewAccountCtx) ConfirmPasswordEdited(b string) {
	if ptr.Pointer() != nil {
		var bC *C.char
		if b != "" {
			bC = C.CString(b)
			defer C.free(unsafe.Pointer(bC))
		}
		C.ApproveNewAccountCtx687eda_ConfirmPasswordEdited(ptr.Pointer(), C.struct_Moc_PackedString{data: bC, len: C.longlong(len(b))})
	}
}

//export callbackApproveNewAccountCtx687eda_Remote
func callbackApproveNewAccountCtx687eda_Remote(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "remote"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveNewAccountCtxFromPointer(ptr).RemoteDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveNewAccountCtx) ConnectRemote(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "remote"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "remote", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "remote", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveNewAccountCtx) DisconnectRemote() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "remote")
	}
}

func (ptr *ApproveNewAccountCtx) Remote() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveNewAccountCtx687eda_Remote(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveNewAccountCtx) RemoteDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveNewAccountCtx687eda_RemoteDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveNewAccountCtx687eda_SetRemote
func callbackApproveNewAccountCtx687eda_SetRemote(ptr unsafe.Pointer, remote C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setRemote"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(remote))
	} else {
		NewApproveNewAccountCtxFromPointer(ptr).SetRemoteDefault(cGoUnpackString(remote))
	}
}

func (ptr *ApproveNewAccountCtx) ConnectSetRemote(f func(remote string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setRemote"); signal != nil {
			f := func(remote string) {
				(*(*func(string))(signal))(remote)
				f(remote)
			}
			qt.ConnectSignal(ptr.Pointer(), "setRemote", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setRemote", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveNewAccountCtx) DisconnectSetRemote() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setRemote")
	}
}

func (ptr *ApproveNewAccountCtx) SetRemote(remote string) {
	if ptr.Pointer() != nil {
		var remoteC *C.char
		if remote != "" {
			remoteC = C.CString(remote)
			defer C.free(unsafe.Pointer(remoteC))
		}
		C.ApproveNewAccountCtx687eda_SetRemote(ptr.Pointer(), C.struct_Moc_PackedString{data: remoteC, len: C.longlong(len(remote))})
	}
}

func (ptr *ApproveNewAccountCtx) SetRemoteDefault(remote string) {
	if ptr.Pointer() != nil {
		var remoteC *C.char
		if remote != "" {
			remoteC = C.CString(remote)
			defer C.free(unsafe.Pointer(remoteC))
		}
		C.ApproveNewAccountCtx687eda_SetRemoteDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: remoteC, len: C.longlong(len(remote))})
	}
}

//export callbackApproveNewAccountCtx687eda_RemoteChanged
func callbackApproveNewAccountCtx687eda_RemoteChanged(ptr unsafe.Pointer, remote C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "remoteChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(remote))
	}

}

func (ptr *ApproveNewAccountCtx) ConnectRemoteChanged(f func(remote string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "remoteChanged") {
			C.ApproveNewAccountCtx687eda_ConnectRemoteChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "remoteChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "remoteChanged"); signal != nil {
			f := func(remote string) {
				(*(*func(string))(signal))(remote)
				f(remote)
			}
			qt.ConnectSignal(ptr.Pointer(), "remoteChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "remoteChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveNewAccountCtx) DisconnectRemoteChanged() {
	if ptr.Pointer() != nil {
		C.ApproveNewAccountCtx687eda_DisconnectRemoteChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "remoteChanged")
	}
}

func (ptr *ApproveNewAccountCtx) RemoteChanged(remote string) {
	if ptr.Pointer() != nil {
		var remoteC *C.char
		if remote != "" {
			remoteC = C.CString(remote)
			defer C.free(unsafe.Pointer(remoteC))
		}
		C.ApproveNewAccountCtx687eda_RemoteChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: remoteC, len: C.longlong(len(remote))})
	}
}

//export callbackApproveNewAccountCtx687eda_Transport
func callbackApproveNewAccountCtx687eda_Transport(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "transport"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveNewAccountCtxFromPointer(ptr).TransportDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveNewAccountCtx) ConnectTransport(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "transport"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "transport", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "transport", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveNewAccountCtx) DisconnectTransport() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "transport")
	}
}

func (ptr *ApproveNewAccountCtx) Transport() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveNewAccountCtx687eda_Transport(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveNewAccountCtx) TransportDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveNewAccountCtx687eda_TransportDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveNewAccountCtx687eda_SetTransport
func callbackApproveNewAccountCtx687eda_SetTransport(ptr unsafe.Pointer, transport C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setTransport"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(transport))
	} else {
		NewApproveNewAccountCtxFromPointer(ptr).SetTransportDefault(cGoUnpackString(transport))
	}
}

func (ptr *ApproveNewAccountCtx) ConnectSetTransport(f func(transport string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setTransport"); signal != nil {
			f := func(transport string) {
				(*(*func(string))(signal))(transport)
				f(transport)
			}
			qt.ConnectSignal(ptr.Pointer(), "setTransport", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setTransport", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveNewAccountCtx) DisconnectSetTransport() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setTransport")
	}
}

func (ptr *ApproveNewAccountCtx) SetTransport(transport string) {
	if ptr.Pointer() != nil {
		var transportC *C.char
		if transport != "" {
			transportC = C.CString(transport)
			defer C.free(unsafe.Pointer(transportC))
		}
		C.ApproveNewAccountCtx687eda_SetTransport(ptr.Pointer(), C.struct_Moc_PackedString{data: transportC, len: C.longlong(len(transport))})
	}
}

func (ptr *ApproveNewAccountCtx) SetTransportDefault(transport string) {
	if ptr.Pointer() != nil {
		var transportC *C.char
		if transport != "" {
			transportC = C.CString(transport)
			defer C.free(unsafe.Pointer(transportC))
		}
		C.ApproveNewAccountCtx687eda_SetTransportDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: transportC, len: C.longlong(len(transport))})
	}
}

//export callbackApproveNewAccountCtx687eda_TransportChanged
func callbackApproveNewAccountCtx687eda_TransportChanged(ptr unsafe.Pointer, transport C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "transportChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(transport))
	}

}

func (ptr *ApproveNewAccountCtx) ConnectTransportChanged(f func(transport string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "transportChanged") {
			C.ApproveNewAccountCtx687eda_ConnectTransportChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "transportChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "transportChanged"); signal != nil {
			f := func(transport string) {
				(*(*func(string))(signal))(transport)
				f(transport)
			}
			qt.ConnectSignal(ptr.Pointer(), "transportChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "transportChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveNewAccountCtx) DisconnectTransportChanged() {
	if ptr.Pointer() != nil {
		C.ApproveNewAccountCtx687eda_DisconnectTransportChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "transportChanged")
	}
}

func (ptr *ApproveNewAccountCtx) TransportChanged(transport string) {
	if ptr.Pointer() != nil {
		var transportC *C.char
		if transport != "" {
			transportC = C.CString(transport)
			defer C.free(unsafe.Pointer(transportC))
		}
		C.ApproveNewAccountCtx687eda_TransportChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: transportC, len: C.longlong(len(transport))})
	}
}

//export callbackApproveNewAccountCtx687eda_Endpoint
func callbackApproveNewAccountCtx687eda_Endpoint(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "endpoint"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveNewAccountCtxFromPointer(ptr).EndpointDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveNewAccountCtx) ConnectEndpoint(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "endpoint"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "endpoint", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "endpoint", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveNewAccountCtx) DisconnectEndpoint() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "endpoint")
	}
}

func (ptr *ApproveNewAccountCtx) Endpoint() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveNewAccountCtx687eda_Endpoint(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveNewAccountCtx) EndpointDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveNewAccountCtx687eda_EndpointDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveNewAccountCtx687eda_SetEndpoint
func callbackApproveNewAccountCtx687eda_SetEndpoint(ptr unsafe.Pointer, endpoint C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setEndpoint"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(endpoint))
	} else {
		NewApproveNewAccountCtxFromPointer(ptr).SetEndpointDefault(cGoUnpackString(endpoint))
	}
}

func (ptr *ApproveNewAccountCtx) ConnectSetEndpoint(f func(endpoint string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setEndpoint"); signal != nil {
			f := func(endpoint string) {
				(*(*func(string))(signal))(endpoint)
				f(endpoint)
			}
			qt.ConnectSignal(ptr.Pointer(), "setEndpoint", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setEndpoint", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveNewAccountCtx) DisconnectSetEndpoint() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setEndpoint")
	}
}

func (ptr *ApproveNewAccountCtx) SetEndpoint(endpoint string) {
	if ptr.Pointer() != nil {
		var endpointC *C.char
		if endpoint != "" {
			endpointC = C.CString(endpoint)
			defer C.free(unsafe.Pointer(endpointC))
		}
		C.ApproveNewAccountCtx687eda_SetEndpoint(ptr.Pointer(), C.struct_Moc_PackedString{data: endpointC, len: C.longlong(len(endpoint))})
	}
}

func (ptr *ApproveNewAccountCtx) SetEndpointDefault(endpoint string) {
	if ptr.Pointer() != nil {
		var endpointC *C.char
		if endpoint != "" {
			endpointC = C.CString(endpoint)
			defer C.free(unsafe.Pointer(endpointC))
		}
		C.ApproveNewAccountCtx687eda_SetEndpointDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: endpointC, len: C.longlong(len(endpoint))})
	}
}

//export callbackApproveNewAccountCtx687eda_EndpointChanged
func callbackApproveNewAccountCtx687eda_EndpointChanged(ptr unsafe.Pointer, endpoint C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "endpointChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(endpoint))
	}

}

func (ptr *ApproveNewAccountCtx) ConnectEndpointChanged(f func(endpoint string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "endpointChanged") {
			C.ApproveNewAccountCtx687eda_ConnectEndpointChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "endpointChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "endpointChanged"); signal != nil {
			f := func(endpoint string) {
				(*(*func(string))(signal))(endpoint)
				f(endpoint)
			}
			qt.ConnectSignal(ptr.Pointer(), "endpointChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "endpointChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveNewAccountCtx) DisconnectEndpointChanged() {
	if ptr.Pointer() != nil {
		C.ApproveNewAccountCtx687eda_DisconnectEndpointChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "endpointChanged")
	}
}

func (ptr *ApproveNewAccountCtx) EndpointChanged(endpoint string) {
	if ptr.Pointer() != nil {
		var endpointC *C.char
		if endpoint != "" {
			endpointC = C.CString(endpoint)
			defer C.free(unsafe.Pointer(endpointC))
		}
		C.ApproveNewAccountCtx687eda_EndpointChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: endpointC, len: C.longlong(len(endpoint))})
	}
}

//export callbackApproveNewAccountCtx687eda_Password
func callbackApproveNewAccountCtx687eda_Password(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "password"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveNewAccountCtxFromPointer(ptr).PasswordDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveNewAccountCtx) ConnectPassword(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "password"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "password", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "password", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveNewAccountCtx) DisconnectPassword() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "password")
	}
}

func (ptr *ApproveNewAccountCtx) Password() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveNewAccountCtx687eda_Password(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveNewAccountCtx) PasswordDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveNewAccountCtx687eda_PasswordDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveNewAccountCtx687eda_SetPassword
func callbackApproveNewAccountCtx687eda_SetPassword(ptr unsafe.Pointer, password C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setPassword"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(password))
	} else {
		NewApproveNewAccountCtxFromPointer(ptr).SetPasswordDefault(cGoUnpackString(password))
	}
}

func (ptr *ApproveNewAccountCtx) ConnectSetPassword(f func(password string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setPassword"); signal != nil {
			f := func(password string) {
				(*(*func(string))(signal))(password)
				f(password)
			}
			qt.ConnectSignal(ptr.Pointer(), "setPassword", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setPassword", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveNewAccountCtx) DisconnectSetPassword() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setPassword")
	}
}

func (ptr *ApproveNewAccountCtx) SetPassword(password string) {
	if ptr.Pointer() != nil {
		var passwordC *C.char
		if password != "" {
			passwordC = C.CString(password)
			defer C.free(unsafe.Pointer(passwordC))
		}
		C.ApproveNewAccountCtx687eda_SetPassword(ptr.Pointer(), C.struct_Moc_PackedString{data: passwordC, len: C.longlong(len(password))})
	}
}

func (ptr *ApproveNewAccountCtx) SetPasswordDefault(password string) {
	if ptr.Pointer() != nil {
		var passwordC *C.char
		if password != "" {
			passwordC = C.CString(password)
			defer C.free(unsafe.Pointer(passwordC))
		}
		C.ApproveNewAccountCtx687eda_SetPasswordDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: passwordC, len: C.longlong(len(password))})
	}
}

//export callbackApproveNewAccountCtx687eda_PasswordChanged
func callbackApproveNewAccountCtx687eda_PasswordChanged(ptr unsafe.Pointer, password C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "passwordChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(password))
	}

}

func (ptr *ApproveNewAccountCtx) ConnectPasswordChanged(f func(password string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "passwordChanged") {
			C.ApproveNewAccountCtx687eda_ConnectPasswordChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "passwordChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "passwordChanged"); signal != nil {
			f := func(password string) {
				(*(*func(string))(signal))(password)
				f(password)
			}
			qt.ConnectSignal(ptr.Pointer(), "passwordChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "passwordChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveNewAccountCtx) DisconnectPasswordChanged() {
	if ptr.Pointer() != nil {
		C.ApproveNewAccountCtx687eda_DisconnectPasswordChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "passwordChanged")
	}
}

func (ptr *ApproveNewAccountCtx) PasswordChanged(password string) {
	if ptr.Pointer() != nil {
		var passwordC *C.char
		if password != "" {
			passwordC = C.CString(password)
			defer C.free(unsafe.Pointer(passwordC))
		}
		C.ApproveNewAccountCtx687eda_PasswordChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: passwordC, len: C.longlong(len(password))})
	}
}

//export callbackApproveNewAccountCtx687eda_ConfirmPassword
func callbackApproveNewAccountCtx687eda_ConfirmPassword(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "confirmPassword"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewApproveNewAccountCtxFromPointer(ptr).ConfirmPasswordDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *ApproveNewAccountCtx) ConnectConfirmPassword(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "confirmPassword"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "confirmPassword", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "confirmPassword", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveNewAccountCtx) DisconnectConfirmPassword() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "confirmPassword")
	}
}

func (ptr *ApproveNewAccountCtx) ConfirmPassword() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveNewAccountCtx687eda_ConfirmPassword(ptr.Pointer()))
	}
	return ""
}

func (ptr *ApproveNewAccountCtx) ConfirmPasswordDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.ApproveNewAccountCtx687eda_ConfirmPasswordDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackApproveNewAccountCtx687eda_SetConfirmPassword
func callbackApproveNewAccountCtx687eda_SetConfirmPassword(ptr unsafe.Pointer, confirmPassword C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setConfirmPassword"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(confirmPassword))
	} else {
		NewApproveNewAccountCtxFromPointer(ptr).SetConfirmPasswordDefault(cGoUnpackString(confirmPassword))
	}
}

func (ptr *ApproveNewAccountCtx) ConnectSetConfirmPassword(f func(confirmPassword string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setConfirmPassword"); signal != nil {
			f := func(confirmPassword string) {
				(*(*func(string))(signal))(confirmPassword)
				f(confirmPassword)
			}
			qt.ConnectSignal(ptr.Pointer(), "setConfirmPassword", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setConfirmPassword", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveNewAccountCtx) DisconnectSetConfirmPassword() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setConfirmPassword")
	}
}

func (ptr *ApproveNewAccountCtx) SetConfirmPassword(confirmPassword string) {
	if ptr.Pointer() != nil {
		var confirmPasswordC *C.char
		if confirmPassword != "" {
			confirmPasswordC = C.CString(confirmPassword)
			defer C.free(unsafe.Pointer(confirmPasswordC))
		}
		C.ApproveNewAccountCtx687eda_SetConfirmPassword(ptr.Pointer(), C.struct_Moc_PackedString{data: confirmPasswordC, len: C.longlong(len(confirmPassword))})
	}
}

func (ptr *ApproveNewAccountCtx) SetConfirmPasswordDefault(confirmPassword string) {
	if ptr.Pointer() != nil {
		var confirmPasswordC *C.char
		if confirmPassword != "" {
			confirmPasswordC = C.CString(confirmPassword)
			defer C.free(unsafe.Pointer(confirmPasswordC))
		}
		C.ApproveNewAccountCtx687eda_SetConfirmPasswordDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: confirmPasswordC, len: C.longlong(len(confirmPassword))})
	}
}

//export callbackApproveNewAccountCtx687eda_ConfirmPasswordChanged
func callbackApproveNewAccountCtx687eda_ConfirmPasswordChanged(ptr unsafe.Pointer, confirmPassword C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "confirmPasswordChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(confirmPassword))
	}

}

func (ptr *ApproveNewAccountCtx) ConnectConfirmPasswordChanged(f func(confirmPassword string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "confirmPasswordChanged") {
			C.ApproveNewAccountCtx687eda_ConnectConfirmPasswordChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "confirmPasswordChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "confirmPasswordChanged"); signal != nil {
			f := func(confirmPassword string) {
				(*(*func(string))(signal))(confirmPassword)
				f(confirmPassword)
			}
			qt.ConnectSignal(ptr.Pointer(), "confirmPasswordChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "confirmPasswordChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveNewAccountCtx) DisconnectConfirmPasswordChanged() {
	if ptr.Pointer() != nil {
		C.ApproveNewAccountCtx687eda_DisconnectConfirmPasswordChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "confirmPasswordChanged")
	}
}

func (ptr *ApproveNewAccountCtx) ConfirmPasswordChanged(confirmPassword string) {
	if ptr.Pointer() != nil {
		var confirmPasswordC *C.char
		if confirmPassword != "" {
			confirmPasswordC = C.CString(confirmPassword)
			defer C.free(unsafe.Pointer(confirmPasswordC))
		}
		C.ApproveNewAccountCtx687eda_ConfirmPasswordChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: confirmPasswordC, len: C.longlong(len(confirmPassword))})
	}
}

func ApproveNewAccountCtx_QRegisterMetaType() int {
	return int(int32(C.ApproveNewAccountCtx687eda_ApproveNewAccountCtx687eda_QRegisterMetaType()))
}

func (ptr *ApproveNewAccountCtx) QRegisterMetaType() int {
	return int(int32(C.ApproveNewAccountCtx687eda_ApproveNewAccountCtx687eda_QRegisterMetaType()))
}

func ApproveNewAccountCtx_QRegisterMetaType2(typeName string) int {
	var typeNameC *C.char
	if typeName != "" {
		typeNameC = C.CString(typeName)
		defer C.free(unsafe.Pointer(typeNameC))
	}
	return int(int32(C.ApproveNewAccountCtx687eda_ApproveNewAccountCtx687eda_QRegisterMetaType2(typeNameC)))
}

func (ptr *ApproveNewAccountCtx) QRegisterMetaType2(typeName string) int {
	var typeNameC *C.char
	if typeName != "" {
		typeNameC = C.CString(typeName)
		defer C.free(unsafe.Pointer(typeNameC))
	}
	return int(int32(C.ApproveNewAccountCtx687eda_ApproveNewAccountCtx687eda_QRegisterMetaType2(typeNameC)))
}

func ApproveNewAccountCtx_QmlRegisterType() int {
	return int(int32(C.ApproveNewAccountCtx687eda_ApproveNewAccountCtx687eda_QmlRegisterType()))
}

func (ptr *ApproveNewAccountCtx) QmlRegisterType() int {
	return int(int32(C.ApproveNewAccountCtx687eda_ApproveNewAccountCtx687eda_QmlRegisterType()))
}

func ApproveNewAccountCtx_QmlRegisterType2(uri string, versionMajor int, versionMinor int, qmlName string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	return int(int32(C.ApproveNewAccountCtx687eda_ApproveNewAccountCtx687eda_QmlRegisterType2(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC)))
}

func (ptr *ApproveNewAccountCtx) QmlRegisterType2(uri string, versionMajor int, versionMinor int, qmlName string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	return int(int32(C.ApproveNewAccountCtx687eda_ApproveNewAccountCtx687eda_QmlRegisterType2(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC)))
}

func ApproveNewAccountCtx_QmlRegisterUncreatableType(uri string, versionMajor int, versionMinor int, qmlName string, message string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	var messageC *C.char
	if message != "" {
		messageC = C.CString(message)
		defer C.free(unsafe.Pointer(messageC))
	}
	return int(int32(C.ApproveNewAccountCtx687eda_ApproveNewAccountCtx687eda_QmlRegisterUncreatableType(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC, C.struct_Moc_PackedString{data: messageC, len: C.longlong(len(message))})))
}

func (ptr *ApproveNewAccountCtx) QmlRegisterUncreatableType(uri string, versionMajor int, versionMinor int, qmlName string, message string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	var messageC *C.char
	if message != "" {
		messageC = C.CString(message)
		defer C.free(unsafe.Pointer(messageC))
	}
	return int(int32(C.ApproveNewAccountCtx687eda_ApproveNewAccountCtx687eda_QmlRegisterUncreatableType(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC, C.struct_Moc_PackedString{data: messageC, len: C.longlong(len(message))})))
}

func (ptr *ApproveNewAccountCtx) __children_atList(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQObjectFromPointer(C.ApproveNewAccountCtx687eda___children_atList(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *ApproveNewAccountCtx) __children_setList(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveNewAccountCtx687eda___children_setList(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *ApproveNewAccountCtx) __children_newList() unsafe.Pointer {
	return C.ApproveNewAccountCtx687eda___children_newList(ptr.Pointer())
}

func (ptr *ApproveNewAccountCtx) __dynamicPropertyNames_atList(i int) *std_core.QByteArray {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQByteArrayFromPointer(C.ApproveNewAccountCtx687eda___dynamicPropertyNames_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QByteArray).DestroyQByteArray)
		return tmpValue
	}
	return nil
}

func (ptr *ApproveNewAccountCtx) __dynamicPropertyNames_setList(i std_core.QByteArray_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveNewAccountCtx687eda___dynamicPropertyNames_setList(ptr.Pointer(), std_core.PointerFromQByteArray(i))
	}
}

func (ptr *ApproveNewAccountCtx) __dynamicPropertyNames_newList() unsafe.Pointer {
	return C.ApproveNewAccountCtx687eda___dynamicPropertyNames_newList(ptr.Pointer())
}

func (ptr *ApproveNewAccountCtx) __findChildren_atList(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQObjectFromPointer(C.ApproveNewAccountCtx687eda___findChildren_atList(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *ApproveNewAccountCtx) __findChildren_setList(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveNewAccountCtx687eda___findChildren_setList(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *ApproveNewAccountCtx) __findChildren_newList() unsafe.Pointer {
	return C.ApproveNewAccountCtx687eda___findChildren_newList(ptr.Pointer())
}

func (ptr *ApproveNewAccountCtx) __findChildren_atList3(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQObjectFromPointer(C.ApproveNewAccountCtx687eda___findChildren_atList3(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *ApproveNewAccountCtx) __findChildren_setList3(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveNewAccountCtx687eda___findChildren_setList3(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *ApproveNewAccountCtx) __findChildren_newList3() unsafe.Pointer {
	return C.ApproveNewAccountCtx687eda___findChildren_newList3(ptr.Pointer())
}

func NewApproveNewAccountCtx(parent std_core.QObject_ITF) *ApproveNewAccountCtx {
	ApproveNewAccountCtx_QRegisterMetaType()
	tmpValue := NewApproveNewAccountCtxFromPointer(C.ApproveNewAccountCtx687eda_NewApproveNewAccountCtx(std_core.PointerFromQObject(parent)))
	if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
		tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
	}
	return tmpValue
}

//export callbackApproveNewAccountCtx687eda_DestroyApproveNewAccountCtx
func callbackApproveNewAccountCtx687eda_DestroyApproveNewAccountCtx(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "~ApproveNewAccountCtx"); signal != nil {
		(*(*func())(signal))()
	} else {
		NewApproveNewAccountCtxFromPointer(ptr).DestroyApproveNewAccountCtxDefault()
	}
}

func (ptr *ApproveNewAccountCtx) ConnectDestroyApproveNewAccountCtx(f func()) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "~ApproveNewAccountCtx"); signal != nil {
			f := func() {
				(*(*func())(signal))()
				f()
			}
			qt.ConnectSignal(ptr.Pointer(), "~ApproveNewAccountCtx", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "~ApproveNewAccountCtx", unsafe.Pointer(&f))
		}
	}
}

func (ptr *ApproveNewAccountCtx) DisconnectDestroyApproveNewAccountCtx() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "~ApproveNewAccountCtx")
	}
}

func (ptr *ApproveNewAccountCtx) DestroyApproveNewAccountCtx() {
	if ptr.Pointer() != nil {

		qt.SetFinalizer(ptr, nil)
		C.ApproveNewAccountCtx687eda_DestroyApproveNewAccountCtx(ptr.Pointer())
		ptr.SetPointer(nil)
	}
}

func (ptr *ApproveNewAccountCtx) DestroyApproveNewAccountCtxDefault() {
	if ptr.Pointer() != nil {

		qt.SetFinalizer(ptr, nil)
		C.ApproveNewAccountCtx687eda_DestroyApproveNewAccountCtxDefault(ptr.Pointer())
		ptr.SetPointer(nil)
	}
}

//export callbackApproveNewAccountCtx687eda_ChildEvent
func callbackApproveNewAccountCtx687eda_ChildEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "childEvent"); signal != nil {
		(*(*func(*std_core.QChildEvent))(signal))(std_core.NewQChildEventFromPointer(event))
	} else {
		NewApproveNewAccountCtxFromPointer(ptr).ChildEventDefault(std_core.NewQChildEventFromPointer(event))
	}
}

func (ptr *ApproveNewAccountCtx) ChildEventDefault(event std_core.QChildEvent_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveNewAccountCtx687eda_ChildEventDefault(ptr.Pointer(), std_core.PointerFromQChildEvent(event))
	}
}

//export callbackApproveNewAccountCtx687eda_ConnectNotify
func callbackApproveNewAccountCtx687eda_ConnectNotify(ptr unsafe.Pointer, sign unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "connectNotify"); signal != nil {
		(*(*func(*std_core.QMetaMethod))(signal))(std_core.NewQMetaMethodFromPointer(sign))
	} else {
		NewApproveNewAccountCtxFromPointer(ptr).ConnectNotifyDefault(std_core.NewQMetaMethodFromPointer(sign))
	}
}

func (ptr *ApproveNewAccountCtx) ConnectNotifyDefault(sign std_core.QMetaMethod_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveNewAccountCtx687eda_ConnectNotifyDefault(ptr.Pointer(), std_core.PointerFromQMetaMethod(sign))
	}
}

//export callbackApproveNewAccountCtx687eda_CustomEvent
func callbackApproveNewAccountCtx687eda_CustomEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "customEvent"); signal != nil {
		(*(*func(*std_core.QEvent))(signal))(std_core.NewQEventFromPointer(event))
	} else {
		NewApproveNewAccountCtxFromPointer(ptr).CustomEventDefault(std_core.NewQEventFromPointer(event))
	}
}

func (ptr *ApproveNewAccountCtx) CustomEventDefault(event std_core.QEvent_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveNewAccountCtx687eda_CustomEventDefault(ptr.Pointer(), std_core.PointerFromQEvent(event))
	}
}

//export callbackApproveNewAccountCtx687eda_DeleteLater
func callbackApproveNewAccountCtx687eda_DeleteLater(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "deleteLater"); signal != nil {
		(*(*func())(signal))()
	} else {
		NewApproveNewAccountCtxFromPointer(ptr).DeleteLaterDefault()
	}
}

func (ptr *ApproveNewAccountCtx) DeleteLaterDefault() {
	if ptr.Pointer() != nil {

		qt.SetFinalizer(ptr, nil)
		C.ApproveNewAccountCtx687eda_DeleteLaterDefault(ptr.Pointer())
	}
}

//export callbackApproveNewAccountCtx687eda_Destroyed
func callbackApproveNewAccountCtx687eda_Destroyed(ptr unsafe.Pointer, obj unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "destroyed"); signal != nil {
		(*(*func(*std_core.QObject))(signal))(std_core.NewQObjectFromPointer(obj))
	}
	qt.Unregister(ptr)

}

//export callbackApproveNewAccountCtx687eda_DisconnectNotify
func callbackApproveNewAccountCtx687eda_DisconnectNotify(ptr unsafe.Pointer, sign unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "disconnectNotify"); signal != nil {
		(*(*func(*std_core.QMetaMethod))(signal))(std_core.NewQMetaMethodFromPointer(sign))
	} else {
		NewApproveNewAccountCtxFromPointer(ptr).DisconnectNotifyDefault(std_core.NewQMetaMethodFromPointer(sign))
	}
}

func (ptr *ApproveNewAccountCtx) DisconnectNotifyDefault(sign std_core.QMetaMethod_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveNewAccountCtx687eda_DisconnectNotifyDefault(ptr.Pointer(), std_core.PointerFromQMetaMethod(sign))
	}
}

//export callbackApproveNewAccountCtx687eda_Event
func callbackApproveNewAccountCtx687eda_Event(ptr unsafe.Pointer, e unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "event"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QEvent) bool)(signal))(std_core.NewQEventFromPointer(e)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewApproveNewAccountCtxFromPointer(ptr).EventDefault(std_core.NewQEventFromPointer(e)))))
}

func (ptr *ApproveNewAccountCtx) EventDefault(e std_core.QEvent_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.ApproveNewAccountCtx687eda_EventDefault(ptr.Pointer(), std_core.PointerFromQEvent(e))) != 0
	}
	return false
}

//export callbackApproveNewAccountCtx687eda_EventFilter
func callbackApproveNewAccountCtx687eda_EventFilter(ptr unsafe.Pointer, watched unsafe.Pointer, event unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "eventFilter"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QObject, *std_core.QEvent) bool)(signal))(std_core.NewQObjectFromPointer(watched), std_core.NewQEventFromPointer(event)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewApproveNewAccountCtxFromPointer(ptr).EventFilterDefault(std_core.NewQObjectFromPointer(watched), std_core.NewQEventFromPointer(event)))))
}

func (ptr *ApproveNewAccountCtx) EventFilterDefault(watched std_core.QObject_ITF, event std_core.QEvent_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.ApproveNewAccountCtx687eda_EventFilterDefault(ptr.Pointer(), std_core.PointerFromQObject(watched), std_core.PointerFromQEvent(event))) != 0
	}
	return false
}

//export callbackApproveNewAccountCtx687eda_ObjectNameChanged
func callbackApproveNewAccountCtx687eda_ObjectNameChanged(ptr unsafe.Pointer, objectName C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "objectNameChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(objectName))
	}

}

//export callbackApproveNewAccountCtx687eda_TimerEvent
func callbackApproveNewAccountCtx687eda_TimerEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "timerEvent"); signal != nil {
		(*(*func(*std_core.QTimerEvent))(signal))(std_core.NewQTimerEventFromPointer(event))
	} else {
		NewApproveNewAccountCtxFromPointer(ptr).TimerEventDefault(std_core.NewQTimerEventFromPointer(event))
	}
}

func (ptr *ApproveNewAccountCtx) TimerEventDefault(event std_core.QTimerEvent_ITF) {
	if ptr.Pointer() != nil {
		C.ApproveNewAccountCtx687eda_TimerEventDefault(ptr.Pointer(), std_core.PointerFromQTimerEvent(event))
	}
}

type CustomListModel_ITF interface {
	std_core.QAbstractListModel_ITF
	CustomListModel_PTR() *CustomListModel
}

func (ptr *CustomListModel) CustomListModel_PTR() *CustomListModel {
	return ptr
}

func (ptr *CustomListModel) Pointer() unsafe.Pointer {
	if ptr != nil {
		return ptr.QAbstractListModel_PTR().Pointer()
	}
	return nil
}

func (ptr *CustomListModel) SetPointer(p unsafe.Pointer) {
	if ptr != nil {
		ptr.QAbstractListModel_PTR().SetPointer(p)
	}
}

func PointerFromCustomListModel(ptr CustomListModel_ITF) unsafe.Pointer {
	if ptr != nil {
		return ptr.CustomListModel_PTR().Pointer()
	}
	return nil
}

func NewCustomListModelFromPointer(ptr unsafe.Pointer) (n *CustomListModel) {
	if gPtr, ok := qt.Receive(ptr); !ok {
		n = new(CustomListModel)
		n.SetPointer(ptr)
	} else {
		switch deduced := gPtr.(type) {
		case *CustomListModel:
			n = deduced

		case *std_core.QAbstractListModel:
			n = &CustomListModel{QAbstractListModel: *deduced}

		default:
			n = new(CustomListModel)
			n.SetPointer(ptr)
		}
	}
	return
}
func (this *CustomListModel) Init() { this.init() }

//export callbackCustomListModel687eda_Constructor
func callbackCustomListModel687eda_Constructor(ptr unsafe.Pointer) {
	this := NewCustomListModelFromPointer(ptr)
	qt.Register(ptr, this)
	this.ConnectClear(this.clear)
	this.ConnectAdd(this.add)
	this.init()
}

//export callbackCustomListModel687eda_Clear
func callbackCustomListModel687eda_Clear(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "clear"); signal != nil {
		(*(*func())(signal))()
	}

}

func (ptr *CustomListModel) ConnectClear(f func()) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "clear") {
			C.CustomListModel687eda_ConnectClear(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "clear")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "clear"); signal != nil {
			f := func() {
				(*(*func())(signal))()
				f()
			}
			qt.ConnectSignal(ptr.Pointer(), "clear", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "clear", unsafe.Pointer(&f))
		}
	}
}

func (ptr *CustomListModel) DisconnectClear() {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda_DisconnectClear(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "clear")
	}
}

func (ptr *CustomListModel) Clear() {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda_Clear(ptr.Pointer())
	}
}

//export callbackCustomListModel687eda_Add
func callbackCustomListModel687eda_Add(ptr unsafe.Pointer, account C.uintptr_t) {
	var accountD custom_accounts_902b5fm.Account
	if accountI, ok := qt.ReceiveTemp(unsafe.Pointer(uintptr(account))); ok {
		qt.UnregisterTemp(unsafe.Pointer(uintptr(account)))
		accountD = (*(*custom_accounts_902b5fm.Account)(accountI))
	}
	if signal := qt.GetSignal(ptr, "add"); signal != nil {
		(*(*func(custom_accounts_902b5fm.Account))(signal))(accountD)
	}

}

func (ptr *CustomListModel) ConnectAdd(f func(account custom_accounts_902b5fm.Account)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "add") {
			C.CustomListModel687eda_ConnectAdd(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "add")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "add"); signal != nil {
			f := func(account custom_accounts_902b5fm.Account) {
				(*(*func(custom_accounts_902b5fm.Account))(signal))(account)
				f(account)
			}
			qt.ConnectSignal(ptr.Pointer(), "add", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "add", unsafe.Pointer(&f))
		}
	}
}

func (ptr *CustomListModel) DisconnectAdd() {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda_DisconnectAdd(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "add")
	}
}

func (ptr *CustomListModel) Add(account custom_accounts_902b5fm.Account) {
	if ptr.Pointer() != nil {
		accountTID := time.Now().UnixNano() + int64(uintptr(unsafe.Pointer(&account)))
		qt.RegisterTemp(unsafe.Pointer(uintptr(accountTID)), unsafe.Pointer(&account))
		C.CustomListModel687eda_Add(ptr.Pointer(), C.uintptr_t(accountTID))
	}
}

//export callbackCustomListModel687eda_CallbackFromQml
func callbackCustomListModel687eda_CallbackFromQml(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "callbackFromQml"); signal != nil {
		(*(*func())(signal))()
	}

}

func (ptr *CustomListModel) ConnectCallbackFromQml(f func()) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "callbackFromQml"); signal != nil {
			f := func() {
				(*(*func())(signal))()
				f()
			}
			qt.ConnectSignal(ptr.Pointer(), "callbackFromQml", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "callbackFromQml", unsafe.Pointer(&f))
		}
	}
}

func (ptr *CustomListModel) DisconnectCallbackFromQml() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "callbackFromQml")
	}
}

func (ptr *CustomListModel) CallbackFromQml() {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda_CallbackFromQml(ptr.Pointer())
	}
}

//export callbackCustomListModel687eda_UpdateRequired
func callbackCustomListModel687eda_UpdateRequired(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "updateRequired"); signal != nil {
		tempVal := (*(*func() string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
	}
	tempVal := NewCustomListModelFromPointer(ptr).UpdateRequiredDefault()
	return C.struct_Moc_PackedString{data: C.CString(tempVal), len: C.longlong(len(tempVal))}
}

func (ptr *CustomListModel) ConnectUpdateRequired(f func() string) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "updateRequired"); signal != nil {
			f := func() string {
				(*(*func() string)(signal))()
				return f()
			}
			qt.ConnectSignal(ptr.Pointer(), "updateRequired", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "updateRequired", unsafe.Pointer(&f))
		}
	}
}

func (ptr *CustomListModel) DisconnectUpdateRequired() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "updateRequired")
	}
}

func (ptr *CustomListModel) UpdateRequired() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.CustomListModel687eda_UpdateRequired(ptr.Pointer()))
	}
	return ""
}

func (ptr *CustomListModel) UpdateRequiredDefault() string {
	if ptr.Pointer() != nil {
		return cGoUnpackString(C.CustomListModel687eda_UpdateRequiredDefault(ptr.Pointer()))
	}
	return ""
}

//export callbackCustomListModel687eda_SetUpdateRequired
func callbackCustomListModel687eda_SetUpdateRequired(ptr unsafe.Pointer, updateRequired C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "setUpdateRequired"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(updateRequired))
	} else {
		NewCustomListModelFromPointer(ptr).SetUpdateRequiredDefault(cGoUnpackString(updateRequired))
	}
}

func (ptr *CustomListModel) ConnectSetUpdateRequired(f func(updateRequired string)) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "setUpdateRequired"); signal != nil {
			f := func(updateRequired string) {
				(*(*func(string))(signal))(updateRequired)
				f(updateRequired)
			}
			qt.ConnectSignal(ptr.Pointer(), "setUpdateRequired", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "setUpdateRequired", unsafe.Pointer(&f))
		}
	}
}

func (ptr *CustomListModel) DisconnectSetUpdateRequired() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "setUpdateRequired")
	}
}

func (ptr *CustomListModel) SetUpdateRequired(updateRequired string) {
	if ptr.Pointer() != nil {
		var updateRequiredC *C.char
		if updateRequired != "" {
			updateRequiredC = C.CString(updateRequired)
			defer C.free(unsafe.Pointer(updateRequiredC))
		}
		C.CustomListModel687eda_SetUpdateRequired(ptr.Pointer(), C.struct_Moc_PackedString{data: updateRequiredC, len: C.longlong(len(updateRequired))})
	}
}

func (ptr *CustomListModel) SetUpdateRequiredDefault(updateRequired string) {
	if ptr.Pointer() != nil {
		var updateRequiredC *C.char
		if updateRequired != "" {
			updateRequiredC = C.CString(updateRequired)
			defer C.free(unsafe.Pointer(updateRequiredC))
		}
		C.CustomListModel687eda_SetUpdateRequiredDefault(ptr.Pointer(), C.struct_Moc_PackedString{data: updateRequiredC, len: C.longlong(len(updateRequired))})
	}
}

//export callbackCustomListModel687eda_UpdateRequiredChanged
func callbackCustomListModel687eda_UpdateRequiredChanged(ptr unsafe.Pointer, updateRequired C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "updateRequiredChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(updateRequired))
	}

}

func (ptr *CustomListModel) ConnectUpdateRequiredChanged(f func(updateRequired string)) {
	if ptr.Pointer() != nil {

		if !qt.ExistsSignal(ptr.Pointer(), "updateRequiredChanged") {
			C.CustomListModel687eda_ConnectUpdateRequiredChanged(ptr.Pointer(), C.longlong(qt.ConnectionType(ptr.Pointer(), "updateRequiredChanged")))
		}

		if signal := qt.LendSignal(ptr.Pointer(), "updateRequiredChanged"); signal != nil {
			f := func(updateRequired string) {
				(*(*func(string))(signal))(updateRequired)
				f(updateRequired)
			}
			qt.ConnectSignal(ptr.Pointer(), "updateRequiredChanged", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "updateRequiredChanged", unsafe.Pointer(&f))
		}
	}
}

func (ptr *CustomListModel) DisconnectUpdateRequiredChanged() {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda_DisconnectUpdateRequiredChanged(ptr.Pointer())
		qt.DisconnectSignal(ptr.Pointer(), "updateRequiredChanged")
	}
}

func (ptr *CustomListModel) UpdateRequiredChanged(updateRequired string) {
	if ptr.Pointer() != nil {
		var updateRequiredC *C.char
		if updateRequired != "" {
			updateRequiredC = C.CString(updateRequired)
			defer C.free(unsafe.Pointer(updateRequiredC))
		}
		C.CustomListModel687eda_UpdateRequiredChanged(ptr.Pointer(), C.struct_Moc_PackedString{data: updateRequiredC, len: C.longlong(len(updateRequired))})
	}
}

func CustomListModel_QRegisterMetaType() int {
	return int(int32(C.CustomListModel687eda_CustomListModel687eda_QRegisterMetaType()))
}

func (ptr *CustomListModel) QRegisterMetaType() int {
	return int(int32(C.CustomListModel687eda_CustomListModel687eda_QRegisterMetaType()))
}

func CustomListModel_QRegisterMetaType2(typeName string) int {
	var typeNameC *C.char
	if typeName != "" {
		typeNameC = C.CString(typeName)
		defer C.free(unsafe.Pointer(typeNameC))
	}
	return int(int32(C.CustomListModel687eda_CustomListModel687eda_QRegisterMetaType2(typeNameC)))
}

func (ptr *CustomListModel) QRegisterMetaType2(typeName string) int {
	var typeNameC *C.char
	if typeName != "" {
		typeNameC = C.CString(typeName)
		defer C.free(unsafe.Pointer(typeNameC))
	}
	return int(int32(C.CustomListModel687eda_CustomListModel687eda_QRegisterMetaType2(typeNameC)))
}

func CustomListModel_QmlRegisterType() int {
	return int(int32(C.CustomListModel687eda_CustomListModel687eda_QmlRegisterType()))
}

func (ptr *CustomListModel) QmlRegisterType() int {
	return int(int32(C.CustomListModel687eda_CustomListModel687eda_QmlRegisterType()))
}

func CustomListModel_QmlRegisterType2(uri string, versionMajor int, versionMinor int, qmlName string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	return int(int32(C.CustomListModel687eda_CustomListModel687eda_QmlRegisterType2(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC)))
}

func (ptr *CustomListModel) QmlRegisterType2(uri string, versionMajor int, versionMinor int, qmlName string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	return int(int32(C.CustomListModel687eda_CustomListModel687eda_QmlRegisterType2(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC)))
}

func CustomListModel_QmlRegisterUncreatableType(uri string, versionMajor int, versionMinor int, qmlName string, message string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	var messageC *C.char
	if message != "" {
		messageC = C.CString(message)
		defer C.free(unsafe.Pointer(messageC))
	}
	return int(int32(C.CustomListModel687eda_CustomListModel687eda_QmlRegisterUncreatableType(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC, C.struct_Moc_PackedString{data: messageC, len: C.longlong(len(message))})))
}

func (ptr *CustomListModel) QmlRegisterUncreatableType(uri string, versionMajor int, versionMinor int, qmlName string, message string) int {
	var uriC *C.char
	if uri != "" {
		uriC = C.CString(uri)
		defer C.free(unsafe.Pointer(uriC))
	}
	var qmlNameC *C.char
	if qmlName != "" {
		qmlNameC = C.CString(qmlName)
		defer C.free(unsafe.Pointer(qmlNameC))
	}
	var messageC *C.char
	if message != "" {
		messageC = C.CString(message)
		defer C.free(unsafe.Pointer(messageC))
	}
	return int(int32(C.CustomListModel687eda_CustomListModel687eda_QmlRegisterUncreatableType(uriC, C.int(int32(versionMajor)), C.int(int32(versionMinor)), qmlNameC, C.struct_Moc_PackedString{data: messageC, len: C.longlong(len(message))})))
}

func (ptr *CustomListModel) ____itemData_keyList_atList(i int) int {
	if ptr.Pointer() != nil {
		return int(int32(C.CustomListModel687eda_____itemData_keyList_atList(ptr.Pointer(), C.int(int32(i)))))
	}
	return 0
}

func (ptr *CustomListModel) ____itemData_keyList_setList(i int) {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda_____itemData_keyList_setList(ptr.Pointer(), C.int(int32(i)))
	}
}

func (ptr *CustomListModel) ____itemData_keyList_newList() unsafe.Pointer {
	return C.CustomListModel687eda_____itemData_keyList_newList(ptr.Pointer())
}

func (ptr *CustomListModel) ____roleNames_keyList_atList(i int) int {
	if ptr.Pointer() != nil {
		return int(int32(C.CustomListModel687eda_____roleNames_keyList_atList(ptr.Pointer(), C.int(int32(i)))))
	}
	return 0
}

func (ptr *CustomListModel) ____roleNames_keyList_setList(i int) {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda_____roleNames_keyList_setList(ptr.Pointer(), C.int(int32(i)))
	}
}

func (ptr *CustomListModel) ____roleNames_keyList_newList() unsafe.Pointer {
	return C.CustomListModel687eda_____roleNames_keyList_newList(ptr.Pointer())
}

func (ptr *CustomListModel) ____setItemData_roles_keyList_atList(i int) int {
	if ptr.Pointer() != nil {
		return int(int32(C.CustomListModel687eda_____setItemData_roles_keyList_atList(ptr.Pointer(), C.int(int32(i)))))
	}
	return 0
}

func (ptr *CustomListModel) ____setItemData_roles_keyList_setList(i int) {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda_____setItemData_roles_keyList_setList(ptr.Pointer(), C.int(int32(i)))
	}
}

func (ptr *CustomListModel) ____setItemData_roles_keyList_newList() unsafe.Pointer {
	return C.CustomListModel687eda_____setItemData_roles_keyList_newList(ptr.Pointer())
}

func (ptr *CustomListModel) __changePersistentIndexList_from_atList(i int) *std_core.QModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQModelIndexFromPointer(C.CustomListModel687eda___changePersistentIndexList_from_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QModelIndex).DestroyQModelIndex)
		return tmpValue
	}
	return nil
}

func (ptr *CustomListModel) __changePersistentIndexList_from_setList(i std_core.QModelIndex_ITF) {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda___changePersistentIndexList_from_setList(ptr.Pointer(), std_core.PointerFromQModelIndex(i))
	}
}

func (ptr *CustomListModel) __changePersistentIndexList_from_newList() unsafe.Pointer {
	return C.CustomListModel687eda___changePersistentIndexList_from_newList(ptr.Pointer())
}

func (ptr *CustomListModel) __changePersistentIndexList_to_atList(i int) *std_core.QModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQModelIndexFromPointer(C.CustomListModel687eda___changePersistentIndexList_to_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QModelIndex).DestroyQModelIndex)
		return tmpValue
	}
	return nil
}

func (ptr *CustomListModel) __changePersistentIndexList_to_setList(i std_core.QModelIndex_ITF) {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda___changePersistentIndexList_to_setList(ptr.Pointer(), std_core.PointerFromQModelIndex(i))
	}
}

func (ptr *CustomListModel) __changePersistentIndexList_to_newList() unsafe.Pointer {
	return C.CustomListModel687eda___changePersistentIndexList_to_newList(ptr.Pointer())
}

func (ptr *CustomListModel) __dataChanged_roles_atList(i int) int {
	if ptr.Pointer() != nil {
		return int(int32(C.CustomListModel687eda___dataChanged_roles_atList(ptr.Pointer(), C.int(int32(i)))))
	}
	return 0
}

func (ptr *CustomListModel) __dataChanged_roles_setList(i int) {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda___dataChanged_roles_setList(ptr.Pointer(), C.int(int32(i)))
	}
}

func (ptr *CustomListModel) __dataChanged_roles_newList() unsafe.Pointer {
	return C.CustomListModel687eda___dataChanged_roles_newList(ptr.Pointer())
}

func (ptr *CustomListModel) __itemData_atList(v int, i int) *std_core.QVariant {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQVariantFromPointer(C.CustomListModel687eda___itemData_atList(ptr.Pointer(), C.int(int32(v)), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QVariant).DestroyQVariant)
		return tmpValue
	}
	return nil
}

func (ptr *CustomListModel) __itemData_setList(key int, i std_core.QVariant_ITF) {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda___itemData_setList(ptr.Pointer(), C.int(int32(key)), std_core.PointerFromQVariant(i))
	}
}

func (ptr *CustomListModel) __itemData_newList() unsafe.Pointer {
	return C.CustomListModel687eda___itemData_newList(ptr.Pointer())
}

func (ptr *CustomListModel) __itemData_keyList() []int {
	if ptr.Pointer() != nil {
		return func(l C.struct_Moc_PackedList) []int {
			out := make([]int, int(l.len))
			tmpList := NewCustomListModelFromPointer(l.data)
			for i := 0; i < len(out); i++ {
				out[i] = tmpList.____itemData_keyList_atList(i)
			}
			return out
		}(C.CustomListModel687eda___itemData_keyList(ptr.Pointer()))
	}
	return make([]int, 0)
}

func (ptr *CustomListModel) __layoutAboutToBeChanged_parents_atList(i int) *std_core.QPersistentModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQPersistentModelIndexFromPointer(C.CustomListModel687eda___layoutAboutToBeChanged_parents_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QPersistentModelIndex).DestroyQPersistentModelIndex)
		return tmpValue
	}
	return nil
}

func (ptr *CustomListModel) __layoutAboutToBeChanged_parents_setList(i std_core.QPersistentModelIndex_ITF) {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda___layoutAboutToBeChanged_parents_setList(ptr.Pointer(), std_core.PointerFromQPersistentModelIndex(i))
	}
}

func (ptr *CustomListModel) __layoutAboutToBeChanged_parents_newList() unsafe.Pointer {
	return C.CustomListModel687eda___layoutAboutToBeChanged_parents_newList(ptr.Pointer())
}

func (ptr *CustomListModel) __layoutChanged_parents_atList(i int) *std_core.QPersistentModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQPersistentModelIndexFromPointer(C.CustomListModel687eda___layoutChanged_parents_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QPersistentModelIndex).DestroyQPersistentModelIndex)
		return tmpValue
	}
	return nil
}

func (ptr *CustomListModel) __layoutChanged_parents_setList(i std_core.QPersistentModelIndex_ITF) {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda___layoutChanged_parents_setList(ptr.Pointer(), std_core.PointerFromQPersistentModelIndex(i))
	}
}

func (ptr *CustomListModel) __layoutChanged_parents_newList() unsafe.Pointer {
	return C.CustomListModel687eda___layoutChanged_parents_newList(ptr.Pointer())
}

func (ptr *CustomListModel) __match_atList(i int) *std_core.QModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQModelIndexFromPointer(C.CustomListModel687eda___match_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QModelIndex).DestroyQModelIndex)
		return tmpValue
	}
	return nil
}

func (ptr *CustomListModel) __match_setList(i std_core.QModelIndex_ITF) {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda___match_setList(ptr.Pointer(), std_core.PointerFromQModelIndex(i))
	}
}

func (ptr *CustomListModel) __match_newList() unsafe.Pointer {
	return C.CustomListModel687eda___match_newList(ptr.Pointer())
}

func (ptr *CustomListModel) __mimeData_indexes_atList(i int) *std_core.QModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQModelIndexFromPointer(C.CustomListModel687eda___mimeData_indexes_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QModelIndex).DestroyQModelIndex)
		return tmpValue
	}
	return nil
}

func (ptr *CustomListModel) __mimeData_indexes_setList(i std_core.QModelIndex_ITF) {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda___mimeData_indexes_setList(ptr.Pointer(), std_core.PointerFromQModelIndex(i))
	}
}

func (ptr *CustomListModel) __mimeData_indexes_newList() unsafe.Pointer {
	return C.CustomListModel687eda___mimeData_indexes_newList(ptr.Pointer())
}

func (ptr *CustomListModel) __persistentIndexList_atList(i int) *std_core.QModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQModelIndexFromPointer(C.CustomListModel687eda___persistentIndexList_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QModelIndex).DestroyQModelIndex)
		return tmpValue
	}
	return nil
}

func (ptr *CustomListModel) __persistentIndexList_setList(i std_core.QModelIndex_ITF) {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda___persistentIndexList_setList(ptr.Pointer(), std_core.PointerFromQModelIndex(i))
	}
}

func (ptr *CustomListModel) __persistentIndexList_newList() unsafe.Pointer {
	return C.CustomListModel687eda___persistentIndexList_newList(ptr.Pointer())
}

func (ptr *CustomListModel) __roleNames_atList(v int, i int) *std_core.QByteArray {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQByteArrayFromPointer(C.CustomListModel687eda___roleNames_atList(ptr.Pointer(), C.int(int32(v)), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QByteArray).DestroyQByteArray)
		return tmpValue
	}
	return nil
}

func (ptr *CustomListModel) __roleNames_setList(key int, i std_core.QByteArray_ITF) {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda___roleNames_setList(ptr.Pointer(), C.int(int32(key)), std_core.PointerFromQByteArray(i))
	}
}

func (ptr *CustomListModel) __roleNames_newList() unsafe.Pointer {
	return C.CustomListModel687eda___roleNames_newList(ptr.Pointer())
}

func (ptr *CustomListModel) __roleNames_keyList() []int {
	if ptr.Pointer() != nil {
		return func(l C.struct_Moc_PackedList) []int {
			out := make([]int, int(l.len))
			tmpList := NewCustomListModelFromPointer(l.data)
			for i := 0; i < len(out); i++ {
				out[i] = tmpList.____roleNames_keyList_atList(i)
			}
			return out
		}(C.CustomListModel687eda___roleNames_keyList(ptr.Pointer()))
	}
	return make([]int, 0)
}

func (ptr *CustomListModel) __setItemData_roles_atList(v int, i int) *std_core.QVariant {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQVariantFromPointer(C.CustomListModel687eda___setItemData_roles_atList(ptr.Pointer(), C.int(int32(v)), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QVariant).DestroyQVariant)
		return tmpValue
	}
	return nil
}

func (ptr *CustomListModel) __setItemData_roles_setList(key int, i std_core.QVariant_ITF) {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda___setItemData_roles_setList(ptr.Pointer(), C.int(int32(key)), std_core.PointerFromQVariant(i))
	}
}

func (ptr *CustomListModel) __setItemData_roles_newList() unsafe.Pointer {
	return C.CustomListModel687eda___setItemData_roles_newList(ptr.Pointer())
}

func (ptr *CustomListModel) __setItemData_roles_keyList() []int {
	if ptr.Pointer() != nil {
		return func(l C.struct_Moc_PackedList) []int {
			out := make([]int, int(l.len))
			tmpList := NewCustomListModelFromPointer(l.data)
			for i := 0; i < len(out); i++ {
				out[i] = tmpList.____setItemData_roles_keyList_atList(i)
			}
			return out
		}(C.CustomListModel687eda___setItemData_roles_keyList(ptr.Pointer()))
	}
	return make([]int, 0)
}

func (ptr *CustomListModel) ____doSetRoleNames_roleNames_keyList_atList(i int) int {
	if ptr.Pointer() != nil {
		return int(int32(C.CustomListModel687eda_____doSetRoleNames_roleNames_keyList_atList(ptr.Pointer(), C.int(int32(i)))))
	}
	return 0
}

func (ptr *CustomListModel) ____doSetRoleNames_roleNames_keyList_setList(i int) {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda_____doSetRoleNames_roleNames_keyList_setList(ptr.Pointer(), C.int(int32(i)))
	}
}

func (ptr *CustomListModel) ____doSetRoleNames_roleNames_keyList_newList() unsafe.Pointer {
	return C.CustomListModel687eda_____doSetRoleNames_roleNames_keyList_newList(ptr.Pointer())
}

func (ptr *CustomListModel) ____setRoleNames_roleNames_keyList_atList(i int) int {
	if ptr.Pointer() != nil {
		return int(int32(C.CustomListModel687eda_____setRoleNames_roleNames_keyList_atList(ptr.Pointer(), C.int(int32(i)))))
	}
	return 0
}

func (ptr *CustomListModel) ____setRoleNames_roleNames_keyList_setList(i int) {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda_____setRoleNames_roleNames_keyList_setList(ptr.Pointer(), C.int(int32(i)))
	}
}

func (ptr *CustomListModel) ____setRoleNames_roleNames_keyList_newList() unsafe.Pointer {
	return C.CustomListModel687eda_____setRoleNames_roleNames_keyList_newList(ptr.Pointer())
}

func (ptr *CustomListModel) __children_atList(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQObjectFromPointer(C.CustomListModel687eda___children_atList(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *CustomListModel) __children_setList(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda___children_setList(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *CustomListModel) __children_newList() unsafe.Pointer {
	return C.CustomListModel687eda___children_newList(ptr.Pointer())
}

func (ptr *CustomListModel) __dynamicPropertyNames_atList(i int) *std_core.QByteArray {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQByteArrayFromPointer(C.CustomListModel687eda___dynamicPropertyNames_atList(ptr.Pointer(), C.int(int32(i))))
		qt.SetFinalizer(tmpValue, (*std_core.QByteArray).DestroyQByteArray)
		return tmpValue
	}
	return nil
}

func (ptr *CustomListModel) __dynamicPropertyNames_setList(i std_core.QByteArray_ITF) {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda___dynamicPropertyNames_setList(ptr.Pointer(), std_core.PointerFromQByteArray(i))
	}
}

func (ptr *CustomListModel) __dynamicPropertyNames_newList() unsafe.Pointer {
	return C.CustomListModel687eda___dynamicPropertyNames_newList(ptr.Pointer())
}

func (ptr *CustomListModel) __findChildren_atList(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQObjectFromPointer(C.CustomListModel687eda___findChildren_atList(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *CustomListModel) __findChildren_setList(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda___findChildren_setList(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *CustomListModel) __findChildren_newList() unsafe.Pointer {
	return C.CustomListModel687eda___findChildren_newList(ptr.Pointer())
}

func (ptr *CustomListModel) __findChildren_atList3(i int) *std_core.QObject {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQObjectFromPointer(C.CustomListModel687eda___findChildren_atList3(ptr.Pointer(), C.int(int32(i))))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

func (ptr *CustomListModel) __findChildren_setList3(i std_core.QObject_ITF) {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda___findChildren_setList3(ptr.Pointer(), std_core.PointerFromQObject(i))
	}
}

func (ptr *CustomListModel) __findChildren_newList3() unsafe.Pointer {
	return C.CustomListModel687eda___findChildren_newList3(ptr.Pointer())
}

func NewCustomListModel(parent std_core.QObject_ITF) *CustomListModel {
	CustomListModel_QRegisterMetaType()
	tmpValue := NewCustomListModelFromPointer(C.CustomListModel687eda_NewCustomListModel(std_core.PointerFromQObject(parent)))
	if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
		tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
	}
	return tmpValue
}

//export callbackCustomListModel687eda_DestroyCustomListModel
func callbackCustomListModel687eda_DestroyCustomListModel(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "~CustomListModel"); signal != nil {
		(*(*func())(signal))()
	} else {
		NewCustomListModelFromPointer(ptr).DestroyCustomListModelDefault()
	}
}

func (ptr *CustomListModel) ConnectDestroyCustomListModel(f func()) {
	if ptr.Pointer() != nil {

		if signal := qt.LendSignal(ptr.Pointer(), "~CustomListModel"); signal != nil {
			f := func() {
				(*(*func())(signal))()
				f()
			}
			qt.ConnectSignal(ptr.Pointer(), "~CustomListModel", unsafe.Pointer(&f))
		} else {
			qt.ConnectSignal(ptr.Pointer(), "~CustomListModel", unsafe.Pointer(&f))
		}
	}
}

func (ptr *CustomListModel) DisconnectDestroyCustomListModel() {
	if ptr.Pointer() != nil {

		qt.DisconnectSignal(ptr.Pointer(), "~CustomListModel")
	}
}

func (ptr *CustomListModel) DestroyCustomListModel() {
	if ptr.Pointer() != nil {

		qt.SetFinalizer(ptr, nil)
		C.CustomListModel687eda_DestroyCustomListModel(ptr.Pointer())
		ptr.SetPointer(nil)
	}
}

func (ptr *CustomListModel) DestroyCustomListModelDefault() {
	if ptr.Pointer() != nil {

		qt.SetFinalizer(ptr, nil)
		C.CustomListModel687eda_DestroyCustomListModelDefault(ptr.Pointer())
		ptr.SetPointer(nil)
	}
}

//export callbackCustomListModel687eda_DropMimeData
func callbackCustomListModel687eda_DropMimeData(ptr unsafe.Pointer, data unsafe.Pointer, action C.longlong, row C.int, column C.int, parent unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "dropMimeData"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QMimeData, std_core.Qt__DropAction, int, int, *std_core.QModelIndex) bool)(signal))(std_core.NewQMimeDataFromPointer(data), std_core.Qt__DropAction(action), int(int32(row)), int(int32(column)), std_core.NewQModelIndexFromPointer(parent)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewCustomListModelFromPointer(ptr).DropMimeDataDefault(std_core.NewQMimeDataFromPointer(data), std_core.Qt__DropAction(action), int(int32(row)), int(int32(column)), std_core.NewQModelIndexFromPointer(parent)))))
}

func (ptr *CustomListModel) DropMimeDataDefault(data std_core.QMimeData_ITF, action std_core.Qt__DropAction, row int, column int, parent std_core.QModelIndex_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.CustomListModel687eda_DropMimeDataDefault(ptr.Pointer(), std_core.PointerFromQMimeData(data), C.longlong(action), C.int(int32(row)), C.int(int32(column)), std_core.PointerFromQModelIndex(parent))) != 0
	}
	return false
}

//export callbackCustomListModel687eda_Flags
func callbackCustomListModel687eda_Flags(ptr unsafe.Pointer, index unsafe.Pointer) C.longlong {
	if signal := qt.GetSignal(ptr, "flags"); signal != nil {
		return C.longlong((*(*func(*std_core.QModelIndex) std_core.Qt__ItemFlag)(signal))(std_core.NewQModelIndexFromPointer(index)))
	}

	return C.longlong(NewCustomListModelFromPointer(ptr).FlagsDefault(std_core.NewQModelIndexFromPointer(index)))
}

func (ptr *CustomListModel) FlagsDefault(index std_core.QModelIndex_ITF) std_core.Qt__ItemFlag {
	if ptr.Pointer() != nil {
		return std_core.Qt__ItemFlag(C.CustomListModel687eda_FlagsDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(index)))
	}
	return 0
}

//export callbackCustomListModel687eda_Index
func callbackCustomListModel687eda_Index(ptr unsafe.Pointer, row C.int, column C.int, parent unsafe.Pointer) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "index"); signal != nil {
		return std_core.PointerFromQModelIndex((*(*func(int, int, *std_core.QModelIndex) *std_core.QModelIndex)(signal))(int(int32(row)), int(int32(column)), std_core.NewQModelIndexFromPointer(parent)))
	}

	return std_core.PointerFromQModelIndex(NewCustomListModelFromPointer(ptr).IndexDefault(int(int32(row)), int(int32(column)), std_core.NewQModelIndexFromPointer(parent)))
}

func (ptr *CustomListModel) IndexDefault(row int, column int, parent std_core.QModelIndex_ITF) *std_core.QModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQModelIndexFromPointer(C.CustomListModel687eda_IndexDefault(ptr.Pointer(), C.int(int32(row)), C.int(int32(column)), std_core.PointerFromQModelIndex(parent)))
		qt.SetFinalizer(tmpValue, (*std_core.QModelIndex).DestroyQModelIndex)
		return tmpValue
	}
	return nil
}

//export callbackCustomListModel687eda_Sibling
func callbackCustomListModel687eda_Sibling(ptr unsafe.Pointer, row C.int, column C.int, idx unsafe.Pointer) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "sibling"); signal != nil {
		return std_core.PointerFromQModelIndex((*(*func(int, int, *std_core.QModelIndex) *std_core.QModelIndex)(signal))(int(int32(row)), int(int32(column)), std_core.NewQModelIndexFromPointer(idx)))
	}

	return std_core.PointerFromQModelIndex(NewCustomListModelFromPointer(ptr).SiblingDefault(int(int32(row)), int(int32(column)), std_core.NewQModelIndexFromPointer(idx)))
}

func (ptr *CustomListModel) SiblingDefault(row int, column int, idx std_core.QModelIndex_ITF) *std_core.QModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQModelIndexFromPointer(C.CustomListModel687eda_SiblingDefault(ptr.Pointer(), C.int(int32(row)), C.int(int32(column)), std_core.PointerFromQModelIndex(idx)))
		qt.SetFinalizer(tmpValue, (*std_core.QModelIndex).DestroyQModelIndex)
		return tmpValue
	}
	return nil
}

//export callbackCustomListModel687eda_Buddy
func callbackCustomListModel687eda_Buddy(ptr unsafe.Pointer, index unsafe.Pointer) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "buddy"); signal != nil {
		return std_core.PointerFromQModelIndex((*(*func(*std_core.QModelIndex) *std_core.QModelIndex)(signal))(std_core.NewQModelIndexFromPointer(index)))
	}

	return std_core.PointerFromQModelIndex(NewCustomListModelFromPointer(ptr).BuddyDefault(std_core.NewQModelIndexFromPointer(index)))
}

func (ptr *CustomListModel) BuddyDefault(index std_core.QModelIndex_ITF) *std_core.QModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQModelIndexFromPointer(C.CustomListModel687eda_BuddyDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(index)))
		qt.SetFinalizer(tmpValue, (*std_core.QModelIndex).DestroyQModelIndex)
		return tmpValue
	}
	return nil
}

//export callbackCustomListModel687eda_CanDropMimeData
func callbackCustomListModel687eda_CanDropMimeData(ptr unsafe.Pointer, data unsafe.Pointer, action C.longlong, row C.int, column C.int, parent unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "canDropMimeData"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QMimeData, std_core.Qt__DropAction, int, int, *std_core.QModelIndex) bool)(signal))(std_core.NewQMimeDataFromPointer(data), std_core.Qt__DropAction(action), int(int32(row)), int(int32(column)), std_core.NewQModelIndexFromPointer(parent)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewCustomListModelFromPointer(ptr).CanDropMimeDataDefault(std_core.NewQMimeDataFromPointer(data), std_core.Qt__DropAction(action), int(int32(row)), int(int32(column)), std_core.NewQModelIndexFromPointer(parent)))))
}

func (ptr *CustomListModel) CanDropMimeDataDefault(data std_core.QMimeData_ITF, action std_core.Qt__DropAction, row int, column int, parent std_core.QModelIndex_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.CustomListModel687eda_CanDropMimeDataDefault(ptr.Pointer(), std_core.PointerFromQMimeData(data), C.longlong(action), C.int(int32(row)), C.int(int32(column)), std_core.PointerFromQModelIndex(parent))) != 0
	}
	return false
}

//export callbackCustomListModel687eda_CanFetchMore
func callbackCustomListModel687eda_CanFetchMore(ptr unsafe.Pointer, parent unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "canFetchMore"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QModelIndex) bool)(signal))(std_core.NewQModelIndexFromPointer(parent)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewCustomListModelFromPointer(ptr).CanFetchMoreDefault(std_core.NewQModelIndexFromPointer(parent)))))
}

func (ptr *CustomListModel) CanFetchMoreDefault(parent std_core.QModelIndex_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.CustomListModel687eda_CanFetchMoreDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(parent))) != 0
	}
	return false
}

//export callbackCustomListModel687eda_ColumnCount
func callbackCustomListModel687eda_ColumnCount(ptr unsafe.Pointer, parent unsafe.Pointer) C.int {
	if signal := qt.GetSignal(ptr, "columnCount"); signal != nil {
		return C.int(int32((*(*func(*std_core.QModelIndex) int)(signal))(std_core.NewQModelIndexFromPointer(parent))))
	}

	return C.int(int32(NewCustomListModelFromPointer(ptr).ColumnCountDefault(std_core.NewQModelIndexFromPointer(parent))))
}

func (ptr *CustomListModel) ColumnCountDefault(parent std_core.QModelIndex_ITF) int {
	if ptr.Pointer() != nil {
		return int(int32(C.CustomListModel687eda_ColumnCountDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(parent))))
	}
	return 0
}

//export callbackCustomListModel687eda_ColumnsAboutToBeInserted
func callbackCustomListModel687eda_ColumnsAboutToBeInserted(ptr unsafe.Pointer, parent unsafe.Pointer, first C.int, last C.int) {
	if signal := qt.GetSignal(ptr, "columnsAboutToBeInserted"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(first)), int(int32(last)))
	}

}

//export callbackCustomListModel687eda_ColumnsAboutToBeMoved
func callbackCustomListModel687eda_ColumnsAboutToBeMoved(ptr unsafe.Pointer, sourceParent unsafe.Pointer, sourceStart C.int, sourceEnd C.int, destinationParent unsafe.Pointer, destinationColumn C.int) {
	if signal := qt.GetSignal(ptr, "columnsAboutToBeMoved"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int, *std_core.QModelIndex, int))(signal))(std_core.NewQModelIndexFromPointer(sourceParent), int(int32(sourceStart)), int(int32(sourceEnd)), std_core.NewQModelIndexFromPointer(destinationParent), int(int32(destinationColumn)))
	}

}

//export callbackCustomListModel687eda_ColumnsAboutToBeRemoved
func callbackCustomListModel687eda_ColumnsAboutToBeRemoved(ptr unsafe.Pointer, parent unsafe.Pointer, first C.int, last C.int) {
	if signal := qt.GetSignal(ptr, "columnsAboutToBeRemoved"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(first)), int(int32(last)))
	}

}

//export callbackCustomListModel687eda_ColumnsInserted
func callbackCustomListModel687eda_ColumnsInserted(ptr unsafe.Pointer, parent unsafe.Pointer, first C.int, last C.int) {
	if signal := qt.GetSignal(ptr, "columnsInserted"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(first)), int(int32(last)))
	}

}

//export callbackCustomListModel687eda_ColumnsMoved
func callbackCustomListModel687eda_ColumnsMoved(ptr unsafe.Pointer, parent unsafe.Pointer, start C.int, end C.int, destination unsafe.Pointer, column C.int) {
	if signal := qt.GetSignal(ptr, "columnsMoved"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int, *std_core.QModelIndex, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(start)), int(int32(end)), std_core.NewQModelIndexFromPointer(destination), int(int32(column)))
	}

}

//export callbackCustomListModel687eda_ColumnsRemoved
func callbackCustomListModel687eda_ColumnsRemoved(ptr unsafe.Pointer, parent unsafe.Pointer, first C.int, last C.int) {
	if signal := qt.GetSignal(ptr, "columnsRemoved"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(first)), int(int32(last)))
	}

}

//export callbackCustomListModel687eda_Data
func callbackCustomListModel687eda_Data(ptr unsafe.Pointer, index unsafe.Pointer, role C.int) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "data"); signal != nil {
		return std_core.PointerFromQVariant((*(*func(*std_core.QModelIndex, int) *std_core.QVariant)(signal))(std_core.NewQModelIndexFromPointer(index), int(int32(role))))
	}

	return std_core.PointerFromQVariant(NewCustomListModelFromPointer(ptr).DataDefault(std_core.NewQModelIndexFromPointer(index), int(int32(role))))
}

func (ptr *CustomListModel) DataDefault(index std_core.QModelIndex_ITF, role int) *std_core.QVariant {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQVariantFromPointer(C.CustomListModel687eda_DataDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(index), C.int(int32(role))))
		qt.SetFinalizer(tmpValue, (*std_core.QVariant).DestroyQVariant)
		return tmpValue
	}
	return nil
}

//export callbackCustomListModel687eda_DataChanged
func callbackCustomListModel687eda_DataChanged(ptr unsafe.Pointer, topLeft unsafe.Pointer, bottomRight unsafe.Pointer, roles C.struct_Moc_PackedList) {
	if signal := qt.GetSignal(ptr, "dataChanged"); signal != nil {
		(*(*func(*std_core.QModelIndex, *std_core.QModelIndex, []int))(signal))(std_core.NewQModelIndexFromPointer(topLeft), std_core.NewQModelIndexFromPointer(bottomRight), func(l C.struct_Moc_PackedList) []int {
			out := make([]int, int(l.len))
			tmpList := NewCustomListModelFromPointer(l.data)
			for i := 0; i < len(out); i++ {
				out[i] = tmpList.__dataChanged_roles_atList(i)
			}
			return out
		}(roles))
	}

}

//export callbackCustomListModel687eda_FetchMore
func callbackCustomListModel687eda_FetchMore(ptr unsafe.Pointer, parent unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "fetchMore"); signal != nil {
		(*(*func(*std_core.QModelIndex))(signal))(std_core.NewQModelIndexFromPointer(parent))
	} else {
		NewCustomListModelFromPointer(ptr).FetchMoreDefault(std_core.NewQModelIndexFromPointer(parent))
	}
}

func (ptr *CustomListModel) FetchMoreDefault(parent std_core.QModelIndex_ITF) {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda_FetchMoreDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(parent))
	}
}

//export callbackCustomListModel687eda_HasChildren
func callbackCustomListModel687eda_HasChildren(ptr unsafe.Pointer, parent unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "hasChildren"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QModelIndex) bool)(signal))(std_core.NewQModelIndexFromPointer(parent)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewCustomListModelFromPointer(ptr).HasChildrenDefault(std_core.NewQModelIndexFromPointer(parent)))))
}

func (ptr *CustomListModel) HasChildrenDefault(parent std_core.QModelIndex_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.CustomListModel687eda_HasChildrenDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(parent))) != 0
	}
	return false
}

//export callbackCustomListModel687eda_HeaderData
func callbackCustomListModel687eda_HeaderData(ptr unsafe.Pointer, section C.int, orientation C.longlong, role C.int) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "headerData"); signal != nil {
		return std_core.PointerFromQVariant((*(*func(int, std_core.Qt__Orientation, int) *std_core.QVariant)(signal))(int(int32(section)), std_core.Qt__Orientation(orientation), int(int32(role))))
	}

	return std_core.PointerFromQVariant(NewCustomListModelFromPointer(ptr).HeaderDataDefault(int(int32(section)), std_core.Qt__Orientation(orientation), int(int32(role))))
}

func (ptr *CustomListModel) HeaderDataDefault(section int, orientation std_core.Qt__Orientation, role int) *std_core.QVariant {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQVariantFromPointer(C.CustomListModel687eda_HeaderDataDefault(ptr.Pointer(), C.int(int32(section)), C.longlong(orientation), C.int(int32(role))))
		qt.SetFinalizer(tmpValue, (*std_core.QVariant).DestroyQVariant)
		return tmpValue
	}
	return nil
}

//export callbackCustomListModel687eda_HeaderDataChanged
func callbackCustomListModel687eda_HeaderDataChanged(ptr unsafe.Pointer, orientation C.longlong, first C.int, last C.int) {
	if signal := qt.GetSignal(ptr, "headerDataChanged"); signal != nil {
		(*(*func(std_core.Qt__Orientation, int, int))(signal))(std_core.Qt__Orientation(orientation), int(int32(first)), int(int32(last)))
	}

}

//export callbackCustomListModel687eda_InsertColumns
func callbackCustomListModel687eda_InsertColumns(ptr unsafe.Pointer, column C.int, count C.int, parent unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "insertColumns"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(int, int, *std_core.QModelIndex) bool)(signal))(int(int32(column)), int(int32(count)), std_core.NewQModelIndexFromPointer(parent)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewCustomListModelFromPointer(ptr).InsertColumnsDefault(int(int32(column)), int(int32(count)), std_core.NewQModelIndexFromPointer(parent)))))
}

func (ptr *CustomListModel) InsertColumnsDefault(column int, count int, parent std_core.QModelIndex_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.CustomListModel687eda_InsertColumnsDefault(ptr.Pointer(), C.int(int32(column)), C.int(int32(count)), std_core.PointerFromQModelIndex(parent))) != 0
	}
	return false
}

//export callbackCustomListModel687eda_InsertRows
func callbackCustomListModel687eda_InsertRows(ptr unsafe.Pointer, row C.int, count C.int, parent unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "insertRows"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(int, int, *std_core.QModelIndex) bool)(signal))(int(int32(row)), int(int32(count)), std_core.NewQModelIndexFromPointer(parent)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewCustomListModelFromPointer(ptr).InsertRowsDefault(int(int32(row)), int(int32(count)), std_core.NewQModelIndexFromPointer(parent)))))
}

func (ptr *CustomListModel) InsertRowsDefault(row int, count int, parent std_core.QModelIndex_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.CustomListModel687eda_InsertRowsDefault(ptr.Pointer(), C.int(int32(row)), C.int(int32(count)), std_core.PointerFromQModelIndex(parent))) != 0
	}
	return false
}

//export callbackCustomListModel687eda_ItemData
func callbackCustomListModel687eda_ItemData(ptr unsafe.Pointer, index unsafe.Pointer) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "itemData"); signal != nil {
		return func() unsafe.Pointer {
			tmpList := NewCustomListModelFromPointer(NewCustomListModelFromPointer(nil).__itemData_newList())
			for k, v := range (*(*func(*std_core.QModelIndex) map[int]*std_core.QVariant)(signal))(std_core.NewQModelIndexFromPointer(index)) {
				tmpList.__itemData_setList(k, v)
			}
			return tmpList.Pointer()
		}()
	}

	return func() unsafe.Pointer {
		tmpList := NewCustomListModelFromPointer(NewCustomListModelFromPointer(nil).__itemData_newList())
		for k, v := range NewCustomListModelFromPointer(ptr).ItemDataDefault(std_core.NewQModelIndexFromPointer(index)) {
			tmpList.__itemData_setList(k, v)
		}
		return tmpList.Pointer()
	}()
}

func (ptr *CustomListModel) ItemDataDefault(index std_core.QModelIndex_ITF) map[int]*std_core.QVariant {
	if ptr.Pointer() != nil {
		return func(l C.struct_Moc_PackedList) map[int]*std_core.QVariant {
			out := make(map[int]*std_core.QVariant, int(l.len))
			tmpList := NewCustomListModelFromPointer(l.data)
			for i, v := range tmpList.__itemData_keyList() {
				out[v] = tmpList.__itemData_atList(v, i)
			}
			return out
		}(C.CustomListModel687eda_ItemDataDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(index)))
	}
	return make(map[int]*std_core.QVariant, 0)
}

//export callbackCustomListModel687eda_LayoutAboutToBeChanged
func callbackCustomListModel687eda_LayoutAboutToBeChanged(ptr unsafe.Pointer, parents C.struct_Moc_PackedList, hint C.longlong) {
	if signal := qt.GetSignal(ptr, "layoutAboutToBeChanged"); signal != nil {
		(*(*func([]*std_core.QPersistentModelIndex, std_core.QAbstractItemModel__LayoutChangeHint))(signal))(func(l C.struct_Moc_PackedList) []*std_core.QPersistentModelIndex {
			out := make([]*std_core.QPersistentModelIndex, int(l.len))
			tmpList := NewCustomListModelFromPointer(l.data)
			for i := 0; i < len(out); i++ {
				out[i] = tmpList.__layoutAboutToBeChanged_parents_atList(i)
			}
			return out
		}(parents), std_core.QAbstractItemModel__LayoutChangeHint(hint))
	}

}

//export callbackCustomListModel687eda_LayoutChanged
func callbackCustomListModel687eda_LayoutChanged(ptr unsafe.Pointer, parents C.struct_Moc_PackedList, hint C.longlong) {
	if signal := qt.GetSignal(ptr, "layoutChanged"); signal != nil {
		(*(*func([]*std_core.QPersistentModelIndex, std_core.QAbstractItemModel__LayoutChangeHint))(signal))(func(l C.struct_Moc_PackedList) []*std_core.QPersistentModelIndex {
			out := make([]*std_core.QPersistentModelIndex, int(l.len))
			tmpList := NewCustomListModelFromPointer(l.data)
			for i := 0; i < len(out); i++ {
				out[i] = tmpList.__layoutChanged_parents_atList(i)
			}
			return out
		}(parents), std_core.QAbstractItemModel__LayoutChangeHint(hint))
	}

}

//export callbackCustomListModel687eda_Match
func callbackCustomListModel687eda_Match(ptr unsafe.Pointer, start unsafe.Pointer, role C.int, value unsafe.Pointer, hits C.int, flags C.longlong) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "match"); signal != nil {
		return func() unsafe.Pointer {
			tmpList := NewCustomListModelFromPointer(NewCustomListModelFromPointer(nil).__match_newList())
			for _, v := range (*(*func(*std_core.QModelIndex, int, *std_core.QVariant, int, std_core.Qt__MatchFlag) []*std_core.QModelIndex)(signal))(std_core.NewQModelIndexFromPointer(start), int(int32(role)), std_core.NewQVariantFromPointer(value), int(int32(hits)), std_core.Qt__MatchFlag(flags)) {
				tmpList.__match_setList(v)
			}
			return tmpList.Pointer()
		}()
	}

	return func() unsafe.Pointer {
		tmpList := NewCustomListModelFromPointer(NewCustomListModelFromPointer(nil).__match_newList())
		for _, v := range NewCustomListModelFromPointer(ptr).MatchDefault(std_core.NewQModelIndexFromPointer(start), int(int32(role)), std_core.NewQVariantFromPointer(value), int(int32(hits)), std_core.Qt__MatchFlag(flags)) {
			tmpList.__match_setList(v)
		}
		return tmpList.Pointer()
	}()
}

func (ptr *CustomListModel) MatchDefault(start std_core.QModelIndex_ITF, role int, value std_core.QVariant_ITF, hits int, flags std_core.Qt__MatchFlag) []*std_core.QModelIndex {
	if ptr.Pointer() != nil {
		return func(l C.struct_Moc_PackedList) []*std_core.QModelIndex {
			out := make([]*std_core.QModelIndex, int(l.len))
			tmpList := NewCustomListModelFromPointer(l.data)
			for i := 0; i < len(out); i++ {
				out[i] = tmpList.__match_atList(i)
			}
			return out
		}(C.CustomListModel687eda_MatchDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(start), C.int(int32(role)), std_core.PointerFromQVariant(value), C.int(int32(hits)), C.longlong(flags)))
	}
	return make([]*std_core.QModelIndex, 0)
}

//export callbackCustomListModel687eda_MimeData
func callbackCustomListModel687eda_MimeData(ptr unsafe.Pointer, indexes C.struct_Moc_PackedList) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "mimeData"); signal != nil {
		return std_core.PointerFromQMimeData((*(*func([]*std_core.QModelIndex) *std_core.QMimeData)(signal))(func(l C.struct_Moc_PackedList) []*std_core.QModelIndex {
			out := make([]*std_core.QModelIndex, int(l.len))
			tmpList := NewCustomListModelFromPointer(l.data)
			for i := 0; i < len(out); i++ {
				out[i] = tmpList.__mimeData_indexes_atList(i)
			}
			return out
		}(indexes)))
	}

	return std_core.PointerFromQMimeData(NewCustomListModelFromPointer(ptr).MimeDataDefault(func(l C.struct_Moc_PackedList) []*std_core.QModelIndex {
		out := make([]*std_core.QModelIndex, int(l.len))
		tmpList := NewCustomListModelFromPointer(l.data)
		for i := 0; i < len(out); i++ {
			out[i] = tmpList.__mimeData_indexes_atList(i)
		}
		return out
	}(indexes)))
}

func (ptr *CustomListModel) MimeDataDefault(indexes []*std_core.QModelIndex) *std_core.QMimeData {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQMimeDataFromPointer(C.CustomListModel687eda_MimeDataDefault(ptr.Pointer(), func() unsafe.Pointer {
			tmpList := NewCustomListModelFromPointer(NewCustomListModelFromPointer(nil).__mimeData_indexes_newList())
			for _, v := range indexes {
				tmpList.__mimeData_indexes_setList(v)
			}
			return tmpList.Pointer()
		}()))
		if !qt.ExistsSignal(tmpValue.Pointer(), "destroyed") {
			tmpValue.ConnectDestroyed(func(*std_core.QObject) { tmpValue.SetPointer(nil) })
		}
		return tmpValue
	}
	return nil
}

//export callbackCustomListModel687eda_MimeTypes
func callbackCustomListModel687eda_MimeTypes(ptr unsafe.Pointer) C.struct_Moc_PackedString {
	if signal := qt.GetSignal(ptr, "mimeTypes"); signal != nil {
		tempVal := (*(*func() []string)(signal))()
		return C.struct_Moc_PackedString{data: C.CString(strings.Join(tempVal, "¡¦!")), len: C.longlong(len(strings.Join(tempVal, "¡¦!")))}
	}
	tempVal := NewCustomListModelFromPointer(ptr).MimeTypesDefault()
	return C.struct_Moc_PackedString{data: C.CString(strings.Join(tempVal, "¡¦!")), len: C.longlong(len(strings.Join(tempVal, "¡¦!")))}
}

func (ptr *CustomListModel) MimeTypesDefault() []string {
	if ptr.Pointer() != nil {
		return unpackStringList(cGoUnpackString(C.CustomListModel687eda_MimeTypesDefault(ptr.Pointer())))
	}
	return make([]string, 0)
}

//export callbackCustomListModel687eda_ModelAboutToBeReset
func callbackCustomListModel687eda_ModelAboutToBeReset(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "modelAboutToBeReset"); signal != nil {
		(*(*func())(signal))()
	}

}

//export callbackCustomListModel687eda_ModelReset
func callbackCustomListModel687eda_ModelReset(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "modelReset"); signal != nil {
		(*(*func())(signal))()
	}

}

//export callbackCustomListModel687eda_MoveColumns
func callbackCustomListModel687eda_MoveColumns(ptr unsafe.Pointer, sourceParent unsafe.Pointer, sourceColumn C.int, count C.int, destinationParent unsafe.Pointer, destinationChild C.int) C.char {
	if signal := qt.GetSignal(ptr, "moveColumns"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QModelIndex, int, int, *std_core.QModelIndex, int) bool)(signal))(std_core.NewQModelIndexFromPointer(sourceParent), int(int32(sourceColumn)), int(int32(count)), std_core.NewQModelIndexFromPointer(destinationParent), int(int32(destinationChild))))))
	}

	return C.char(int8(qt.GoBoolToInt(NewCustomListModelFromPointer(ptr).MoveColumnsDefault(std_core.NewQModelIndexFromPointer(sourceParent), int(int32(sourceColumn)), int(int32(count)), std_core.NewQModelIndexFromPointer(destinationParent), int(int32(destinationChild))))))
}

func (ptr *CustomListModel) MoveColumnsDefault(sourceParent std_core.QModelIndex_ITF, sourceColumn int, count int, destinationParent std_core.QModelIndex_ITF, destinationChild int) bool {
	if ptr.Pointer() != nil {
		return int8(C.CustomListModel687eda_MoveColumnsDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(sourceParent), C.int(int32(sourceColumn)), C.int(int32(count)), std_core.PointerFromQModelIndex(destinationParent), C.int(int32(destinationChild)))) != 0
	}
	return false
}

//export callbackCustomListModel687eda_MoveRows
func callbackCustomListModel687eda_MoveRows(ptr unsafe.Pointer, sourceParent unsafe.Pointer, sourceRow C.int, count C.int, destinationParent unsafe.Pointer, destinationChild C.int) C.char {
	if signal := qt.GetSignal(ptr, "moveRows"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QModelIndex, int, int, *std_core.QModelIndex, int) bool)(signal))(std_core.NewQModelIndexFromPointer(sourceParent), int(int32(sourceRow)), int(int32(count)), std_core.NewQModelIndexFromPointer(destinationParent), int(int32(destinationChild))))))
	}

	return C.char(int8(qt.GoBoolToInt(NewCustomListModelFromPointer(ptr).MoveRowsDefault(std_core.NewQModelIndexFromPointer(sourceParent), int(int32(sourceRow)), int(int32(count)), std_core.NewQModelIndexFromPointer(destinationParent), int(int32(destinationChild))))))
}

func (ptr *CustomListModel) MoveRowsDefault(sourceParent std_core.QModelIndex_ITF, sourceRow int, count int, destinationParent std_core.QModelIndex_ITF, destinationChild int) bool {
	if ptr.Pointer() != nil {
		return int8(C.CustomListModel687eda_MoveRowsDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(sourceParent), C.int(int32(sourceRow)), C.int(int32(count)), std_core.PointerFromQModelIndex(destinationParent), C.int(int32(destinationChild)))) != 0
	}
	return false
}

//export callbackCustomListModel687eda_Parent
func callbackCustomListModel687eda_Parent(ptr unsafe.Pointer, index unsafe.Pointer) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "parent"); signal != nil {
		return std_core.PointerFromQModelIndex((*(*func(*std_core.QModelIndex) *std_core.QModelIndex)(signal))(std_core.NewQModelIndexFromPointer(index)))
	}

	return std_core.PointerFromQModelIndex(NewCustomListModelFromPointer(ptr).ParentDefault(std_core.NewQModelIndexFromPointer(index)))
}

func (ptr *CustomListModel) ParentDefault(index std_core.QModelIndex_ITF) *std_core.QModelIndex {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQModelIndexFromPointer(C.CustomListModel687eda_ParentDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(index)))
		qt.SetFinalizer(tmpValue, (*std_core.QModelIndex).DestroyQModelIndex)
		return tmpValue
	}
	return nil
}

//export callbackCustomListModel687eda_RemoveColumns
func callbackCustomListModel687eda_RemoveColumns(ptr unsafe.Pointer, column C.int, count C.int, parent unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "removeColumns"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(int, int, *std_core.QModelIndex) bool)(signal))(int(int32(column)), int(int32(count)), std_core.NewQModelIndexFromPointer(parent)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewCustomListModelFromPointer(ptr).RemoveColumnsDefault(int(int32(column)), int(int32(count)), std_core.NewQModelIndexFromPointer(parent)))))
}

func (ptr *CustomListModel) RemoveColumnsDefault(column int, count int, parent std_core.QModelIndex_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.CustomListModel687eda_RemoveColumnsDefault(ptr.Pointer(), C.int(int32(column)), C.int(int32(count)), std_core.PointerFromQModelIndex(parent))) != 0
	}
	return false
}

//export callbackCustomListModel687eda_RemoveRows
func callbackCustomListModel687eda_RemoveRows(ptr unsafe.Pointer, row C.int, count C.int, parent unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "removeRows"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(int, int, *std_core.QModelIndex) bool)(signal))(int(int32(row)), int(int32(count)), std_core.NewQModelIndexFromPointer(parent)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewCustomListModelFromPointer(ptr).RemoveRowsDefault(int(int32(row)), int(int32(count)), std_core.NewQModelIndexFromPointer(parent)))))
}

func (ptr *CustomListModel) RemoveRowsDefault(row int, count int, parent std_core.QModelIndex_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.CustomListModel687eda_RemoveRowsDefault(ptr.Pointer(), C.int(int32(row)), C.int(int32(count)), std_core.PointerFromQModelIndex(parent))) != 0
	}
	return false
}

//export callbackCustomListModel687eda_ResetInternalData
func callbackCustomListModel687eda_ResetInternalData(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "resetInternalData"); signal != nil {
		(*(*func())(signal))()
	} else {
		NewCustomListModelFromPointer(ptr).ResetInternalDataDefault()
	}
}

func (ptr *CustomListModel) ResetInternalDataDefault() {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda_ResetInternalDataDefault(ptr.Pointer())
	}
}

//export callbackCustomListModel687eda_Revert
func callbackCustomListModel687eda_Revert(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "revert"); signal != nil {
		(*(*func())(signal))()
	} else {
		NewCustomListModelFromPointer(ptr).RevertDefault()
	}
}

func (ptr *CustomListModel) RevertDefault() {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda_RevertDefault(ptr.Pointer())
	}
}

//export callbackCustomListModel687eda_RoleNames
func callbackCustomListModel687eda_RoleNames(ptr unsafe.Pointer) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "roleNames"); signal != nil {
		return func() unsafe.Pointer {
			tmpList := NewCustomListModelFromPointer(NewCustomListModelFromPointer(nil).__roleNames_newList())
			for k, v := range (*(*func() map[int]*std_core.QByteArray)(signal))() {
				tmpList.__roleNames_setList(k, v)
			}
			return tmpList.Pointer()
		}()
	}

	return func() unsafe.Pointer {
		tmpList := NewCustomListModelFromPointer(NewCustomListModelFromPointer(nil).__roleNames_newList())
		for k, v := range NewCustomListModelFromPointer(ptr).RoleNamesDefault() {
			tmpList.__roleNames_setList(k, v)
		}
		return tmpList.Pointer()
	}()
}

func (ptr *CustomListModel) RoleNamesDefault() map[int]*std_core.QByteArray {
	if ptr.Pointer() != nil {
		return func(l C.struct_Moc_PackedList) map[int]*std_core.QByteArray {
			out := make(map[int]*std_core.QByteArray, int(l.len))
			tmpList := NewCustomListModelFromPointer(l.data)
			for i, v := range tmpList.__roleNames_keyList() {
				out[v] = tmpList.__roleNames_atList(v, i)
			}
			return out
		}(C.CustomListModel687eda_RoleNamesDefault(ptr.Pointer()))
	}
	return make(map[int]*std_core.QByteArray, 0)
}

//export callbackCustomListModel687eda_RowCount
func callbackCustomListModel687eda_RowCount(ptr unsafe.Pointer, parent unsafe.Pointer) C.int {
	if signal := qt.GetSignal(ptr, "rowCount"); signal != nil {
		return C.int(int32((*(*func(*std_core.QModelIndex) int)(signal))(std_core.NewQModelIndexFromPointer(parent))))
	}

	return C.int(int32(NewCustomListModelFromPointer(ptr).RowCountDefault(std_core.NewQModelIndexFromPointer(parent))))
}

func (ptr *CustomListModel) RowCountDefault(parent std_core.QModelIndex_ITF) int {
	if ptr.Pointer() != nil {
		return int(int32(C.CustomListModel687eda_RowCountDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(parent))))
	}
	return 0
}

//export callbackCustomListModel687eda_RowsAboutToBeInserted
func callbackCustomListModel687eda_RowsAboutToBeInserted(ptr unsafe.Pointer, parent unsafe.Pointer, start C.int, end C.int) {
	if signal := qt.GetSignal(ptr, "rowsAboutToBeInserted"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(start)), int(int32(end)))
	}

}

//export callbackCustomListModel687eda_RowsAboutToBeMoved
func callbackCustomListModel687eda_RowsAboutToBeMoved(ptr unsafe.Pointer, sourceParent unsafe.Pointer, sourceStart C.int, sourceEnd C.int, destinationParent unsafe.Pointer, destinationRow C.int) {
	if signal := qt.GetSignal(ptr, "rowsAboutToBeMoved"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int, *std_core.QModelIndex, int))(signal))(std_core.NewQModelIndexFromPointer(sourceParent), int(int32(sourceStart)), int(int32(sourceEnd)), std_core.NewQModelIndexFromPointer(destinationParent), int(int32(destinationRow)))
	}

}

//export callbackCustomListModel687eda_RowsAboutToBeRemoved
func callbackCustomListModel687eda_RowsAboutToBeRemoved(ptr unsafe.Pointer, parent unsafe.Pointer, first C.int, last C.int) {
	if signal := qt.GetSignal(ptr, "rowsAboutToBeRemoved"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(first)), int(int32(last)))
	}

}

//export callbackCustomListModel687eda_RowsInserted
func callbackCustomListModel687eda_RowsInserted(ptr unsafe.Pointer, parent unsafe.Pointer, first C.int, last C.int) {
	if signal := qt.GetSignal(ptr, "rowsInserted"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(first)), int(int32(last)))
	}

}

//export callbackCustomListModel687eda_RowsMoved
func callbackCustomListModel687eda_RowsMoved(ptr unsafe.Pointer, parent unsafe.Pointer, start C.int, end C.int, destination unsafe.Pointer, row C.int) {
	if signal := qt.GetSignal(ptr, "rowsMoved"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int, *std_core.QModelIndex, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(start)), int(int32(end)), std_core.NewQModelIndexFromPointer(destination), int(int32(row)))
	}

}

//export callbackCustomListModel687eda_RowsRemoved
func callbackCustomListModel687eda_RowsRemoved(ptr unsafe.Pointer, parent unsafe.Pointer, first C.int, last C.int) {
	if signal := qt.GetSignal(ptr, "rowsRemoved"); signal != nil {
		(*(*func(*std_core.QModelIndex, int, int))(signal))(std_core.NewQModelIndexFromPointer(parent), int(int32(first)), int(int32(last)))
	}

}

//export callbackCustomListModel687eda_SetData
func callbackCustomListModel687eda_SetData(ptr unsafe.Pointer, index unsafe.Pointer, value unsafe.Pointer, role C.int) C.char {
	if signal := qt.GetSignal(ptr, "setData"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QModelIndex, *std_core.QVariant, int) bool)(signal))(std_core.NewQModelIndexFromPointer(index), std_core.NewQVariantFromPointer(value), int(int32(role))))))
	}

	return C.char(int8(qt.GoBoolToInt(NewCustomListModelFromPointer(ptr).SetDataDefault(std_core.NewQModelIndexFromPointer(index), std_core.NewQVariantFromPointer(value), int(int32(role))))))
}

func (ptr *CustomListModel) SetDataDefault(index std_core.QModelIndex_ITF, value std_core.QVariant_ITF, role int) bool {
	if ptr.Pointer() != nil {
		return int8(C.CustomListModel687eda_SetDataDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(index), std_core.PointerFromQVariant(value), C.int(int32(role)))) != 0
	}
	return false
}

//export callbackCustomListModel687eda_SetHeaderData
func callbackCustomListModel687eda_SetHeaderData(ptr unsafe.Pointer, section C.int, orientation C.longlong, value unsafe.Pointer, role C.int) C.char {
	if signal := qt.GetSignal(ptr, "setHeaderData"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(int, std_core.Qt__Orientation, *std_core.QVariant, int) bool)(signal))(int(int32(section)), std_core.Qt__Orientation(orientation), std_core.NewQVariantFromPointer(value), int(int32(role))))))
	}

	return C.char(int8(qt.GoBoolToInt(NewCustomListModelFromPointer(ptr).SetHeaderDataDefault(int(int32(section)), std_core.Qt__Orientation(orientation), std_core.NewQVariantFromPointer(value), int(int32(role))))))
}

func (ptr *CustomListModel) SetHeaderDataDefault(section int, orientation std_core.Qt__Orientation, value std_core.QVariant_ITF, role int) bool {
	if ptr.Pointer() != nil {
		return int8(C.CustomListModel687eda_SetHeaderDataDefault(ptr.Pointer(), C.int(int32(section)), C.longlong(orientation), std_core.PointerFromQVariant(value), C.int(int32(role)))) != 0
	}
	return false
}

//export callbackCustomListModel687eda_SetItemData
func callbackCustomListModel687eda_SetItemData(ptr unsafe.Pointer, index unsafe.Pointer, roles C.struct_Moc_PackedList) C.char {
	if signal := qt.GetSignal(ptr, "setItemData"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QModelIndex, map[int]*std_core.QVariant) bool)(signal))(std_core.NewQModelIndexFromPointer(index), func(l C.struct_Moc_PackedList) map[int]*std_core.QVariant {
			out := make(map[int]*std_core.QVariant, int(l.len))
			tmpList := NewCustomListModelFromPointer(l.data)
			for i, v := range tmpList.__setItemData_roles_keyList() {
				out[v] = tmpList.__setItemData_roles_atList(v, i)
			}
			return out
		}(roles)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewCustomListModelFromPointer(ptr).SetItemDataDefault(std_core.NewQModelIndexFromPointer(index), func(l C.struct_Moc_PackedList) map[int]*std_core.QVariant {
		out := make(map[int]*std_core.QVariant, int(l.len))
		tmpList := NewCustomListModelFromPointer(l.data)
		for i, v := range tmpList.__setItemData_roles_keyList() {
			out[v] = tmpList.__setItemData_roles_atList(v, i)
		}
		return out
	}(roles)))))
}

func (ptr *CustomListModel) SetItemDataDefault(index std_core.QModelIndex_ITF, roles map[int]*std_core.QVariant) bool {
	if ptr.Pointer() != nil {
		return int8(C.CustomListModel687eda_SetItemDataDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(index), func() unsafe.Pointer {
			tmpList := NewCustomListModelFromPointer(NewCustomListModelFromPointer(nil).__setItemData_roles_newList())
			for k, v := range roles {
				tmpList.__setItemData_roles_setList(k, v)
			}
			return tmpList.Pointer()
		}())) != 0
	}
	return false
}

//export callbackCustomListModel687eda_Sort
func callbackCustomListModel687eda_Sort(ptr unsafe.Pointer, column C.int, order C.longlong) {
	if signal := qt.GetSignal(ptr, "sort"); signal != nil {
		(*(*func(int, std_core.Qt__SortOrder))(signal))(int(int32(column)), std_core.Qt__SortOrder(order))
	} else {
		NewCustomListModelFromPointer(ptr).SortDefault(int(int32(column)), std_core.Qt__SortOrder(order))
	}
}

func (ptr *CustomListModel) SortDefault(column int, order std_core.Qt__SortOrder) {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda_SortDefault(ptr.Pointer(), C.int(int32(column)), C.longlong(order))
	}
}

//export callbackCustomListModel687eda_Span
func callbackCustomListModel687eda_Span(ptr unsafe.Pointer, index unsafe.Pointer) unsafe.Pointer {
	if signal := qt.GetSignal(ptr, "span"); signal != nil {
		return std_core.PointerFromQSize((*(*func(*std_core.QModelIndex) *std_core.QSize)(signal))(std_core.NewQModelIndexFromPointer(index)))
	}

	return std_core.PointerFromQSize(NewCustomListModelFromPointer(ptr).SpanDefault(std_core.NewQModelIndexFromPointer(index)))
}

func (ptr *CustomListModel) SpanDefault(index std_core.QModelIndex_ITF) *std_core.QSize {
	if ptr.Pointer() != nil {
		tmpValue := std_core.NewQSizeFromPointer(C.CustomListModel687eda_SpanDefault(ptr.Pointer(), std_core.PointerFromQModelIndex(index)))
		qt.SetFinalizer(tmpValue, (*std_core.QSize).DestroyQSize)
		return tmpValue
	}
	return nil
}

//export callbackCustomListModel687eda_Submit
func callbackCustomListModel687eda_Submit(ptr unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "submit"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func() bool)(signal))())))
	}

	return C.char(int8(qt.GoBoolToInt(NewCustomListModelFromPointer(ptr).SubmitDefault())))
}

func (ptr *CustomListModel) SubmitDefault() bool {
	if ptr.Pointer() != nil {
		return int8(C.CustomListModel687eda_SubmitDefault(ptr.Pointer())) != 0
	}
	return false
}

//export callbackCustomListModel687eda_SupportedDragActions
func callbackCustomListModel687eda_SupportedDragActions(ptr unsafe.Pointer) C.longlong {
	if signal := qt.GetSignal(ptr, "supportedDragActions"); signal != nil {
		return C.longlong((*(*func() std_core.Qt__DropAction)(signal))())
	}

	return C.longlong(NewCustomListModelFromPointer(ptr).SupportedDragActionsDefault())
}

func (ptr *CustomListModel) SupportedDragActionsDefault() std_core.Qt__DropAction {
	if ptr.Pointer() != nil {
		return std_core.Qt__DropAction(C.CustomListModel687eda_SupportedDragActionsDefault(ptr.Pointer()))
	}
	return 0
}

//export callbackCustomListModel687eda_SupportedDropActions
func callbackCustomListModel687eda_SupportedDropActions(ptr unsafe.Pointer) C.longlong {
	if signal := qt.GetSignal(ptr, "supportedDropActions"); signal != nil {
		return C.longlong((*(*func() std_core.Qt__DropAction)(signal))())
	}

	return C.longlong(NewCustomListModelFromPointer(ptr).SupportedDropActionsDefault())
}

func (ptr *CustomListModel) SupportedDropActionsDefault() std_core.Qt__DropAction {
	if ptr.Pointer() != nil {
		return std_core.Qt__DropAction(C.CustomListModel687eda_SupportedDropActionsDefault(ptr.Pointer()))
	}
	return 0
}

//export callbackCustomListModel687eda_ChildEvent
func callbackCustomListModel687eda_ChildEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "childEvent"); signal != nil {
		(*(*func(*std_core.QChildEvent))(signal))(std_core.NewQChildEventFromPointer(event))
	} else {
		NewCustomListModelFromPointer(ptr).ChildEventDefault(std_core.NewQChildEventFromPointer(event))
	}
}

func (ptr *CustomListModel) ChildEventDefault(event std_core.QChildEvent_ITF) {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda_ChildEventDefault(ptr.Pointer(), std_core.PointerFromQChildEvent(event))
	}
}

//export callbackCustomListModel687eda_ConnectNotify
func callbackCustomListModel687eda_ConnectNotify(ptr unsafe.Pointer, sign unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "connectNotify"); signal != nil {
		(*(*func(*std_core.QMetaMethod))(signal))(std_core.NewQMetaMethodFromPointer(sign))
	} else {
		NewCustomListModelFromPointer(ptr).ConnectNotifyDefault(std_core.NewQMetaMethodFromPointer(sign))
	}
}

func (ptr *CustomListModel) ConnectNotifyDefault(sign std_core.QMetaMethod_ITF) {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda_ConnectNotifyDefault(ptr.Pointer(), std_core.PointerFromQMetaMethod(sign))
	}
}

//export callbackCustomListModel687eda_CustomEvent
func callbackCustomListModel687eda_CustomEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "customEvent"); signal != nil {
		(*(*func(*std_core.QEvent))(signal))(std_core.NewQEventFromPointer(event))
	} else {
		NewCustomListModelFromPointer(ptr).CustomEventDefault(std_core.NewQEventFromPointer(event))
	}
}

func (ptr *CustomListModel) CustomEventDefault(event std_core.QEvent_ITF) {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda_CustomEventDefault(ptr.Pointer(), std_core.PointerFromQEvent(event))
	}
}

//export callbackCustomListModel687eda_DeleteLater
func callbackCustomListModel687eda_DeleteLater(ptr unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "deleteLater"); signal != nil {
		(*(*func())(signal))()
	} else {
		NewCustomListModelFromPointer(ptr).DeleteLaterDefault()
	}
}

func (ptr *CustomListModel) DeleteLaterDefault() {
	if ptr.Pointer() != nil {

		qt.SetFinalizer(ptr, nil)
		C.CustomListModel687eda_DeleteLaterDefault(ptr.Pointer())
	}
}

//export callbackCustomListModel687eda_Destroyed
func callbackCustomListModel687eda_Destroyed(ptr unsafe.Pointer, obj unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "destroyed"); signal != nil {
		(*(*func(*std_core.QObject))(signal))(std_core.NewQObjectFromPointer(obj))
	}
	qt.Unregister(ptr)

}

//export callbackCustomListModel687eda_DisconnectNotify
func callbackCustomListModel687eda_DisconnectNotify(ptr unsafe.Pointer, sign unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "disconnectNotify"); signal != nil {
		(*(*func(*std_core.QMetaMethod))(signal))(std_core.NewQMetaMethodFromPointer(sign))
	} else {
		NewCustomListModelFromPointer(ptr).DisconnectNotifyDefault(std_core.NewQMetaMethodFromPointer(sign))
	}
}

func (ptr *CustomListModel) DisconnectNotifyDefault(sign std_core.QMetaMethod_ITF) {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda_DisconnectNotifyDefault(ptr.Pointer(), std_core.PointerFromQMetaMethod(sign))
	}
}

//export callbackCustomListModel687eda_Event
func callbackCustomListModel687eda_Event(ptr unsafe.Pointer, e unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "event"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QEvent) bool)(signal))(std_core.NewQEventFromPointer(e)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewCustomListModelFromPointer(ptr).EventDefault(std_core.NewQEventFromPointer(e)))))
}

func (ptr *CustomListModel) EventDefault(e std_core.QEvent_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.CustomListModel687eda_EventDefault(ptr.Pointer(), std_core.PointerFromQEvent(e))) != 0
	}
	return false
}

//export callbackCustomListModel687eda_EventFilter
func callbackCustomListModel687eda_EventFilter(ptr unsafe.Pointer, watched unsafe.Pointer, event unsafe.Pointer) C.char {
	if signal := qt.GetSignal(ptr, "eventFilter"); signal != nil {
		return C.char(int8(qt.GoBoolToInt((*(*func(*std_core.QObject, *std_core.QEvent) bool)(signal))(std_core.NewQObjectFromPointer(watched), std_core.NewQEventFromPointer(event)))))
	}

	return C.char(int8(qt.GoBoolToInt(NewCustomListModelFromPointer(ptr).EventFilterDefault(std_core.NewQObjectFromPointer(watched), std_core.NewQEventFromPointer(event)))))
}

func (ptr *CustomListModel) EventFilterDefault(watched std_core.QObject_ITF, event std_core.QEvent_ITF) bool {
	if ptr.Pointer() != nil {
		return int8(C.CustomListModel687eda_EventFilterDefault(ptr.Pointer(), std_core.PointerFromQObject(watched), std_core.PointerFromQEvent(event))) != 0
	}
	return false
}

//export callbackCustomListModel687eda_ObjectNameChanged
func callbackCustomListModel687eda_ObjectNameChanged(ptr unsafe.Pointer, objectName C.struct_Moc_PackedString) {
	if signal := qt.GetSignal(ptr, "objectNameChanged"); signal != nil {
		(*(*func(string))(signal))(cGoUnpackString(objectName))
	}

}

//export callbackCustomListModel687eda_TimerEvent
func callbackCustomListModel687eda_TimerEvent(ptr unsafe.Pointer, event unsafe.Pointer) {
	if signal := qt.GetSignal(ptr, "timerEvent"); signal != nil {
		(*(*func(*std_core.QTimerEvent))(signal))(std_core.NewQTimerEventFromPointer(event))
	} else {
		NewCustomListModelFromPointer(ptr).TimerEventDefault(std_core.NewQTimerEventFromPointer(event))
	}
}

func (ptr *CustomListModel) TimerEventDefault(event std_core.QTimerEvent_ITF) {
	if ptr.Pointer() != nil {
		C.CustomListModel687eda_TimerEventDefault(ptr.Pointer(), std_core.PointerFromQTimerEvent(event))
	}
}

func init() {
	qt.ItfMap["main.ApproveListingCtx_ITF"] = ApproveListingCtx{}
	qt.FuncMap["main.NewApproveListingCtx"] = NewApproveListingCtx
	qt.ItfMap["main.ApproveNewAccountCtx_ITF"] = ApproveNewAccountCtx{}
	qt.FuncMap["main.NewApproveNewAccountCtx"] = NewApproveNewAccountCtx
	qt.ItfMap["main.ApproveSignDataCtx_ITF"] = ApproveSignDataCtx{}
	qt.FuncMap["main.NewApproveSignDataCtx"] = NewApproveSignDataCtx
	qt.ItfMap["main.ApproveTxCtx_ITF"] = ApproveTxCtx{}
	qt.FuncMap["main.NewApproveTxCtx"] = NewApproveTxCtx
	qt.ItfMap["main.CustomListModel_ITF"] = CustomListModel{}
	qt.FuncMap["main.NewCustomListModel"] = NewCustomListModel
	qt.ItfMap["main.TxListCtx_ITF"] = TxListCtx{}
	qt.FuncMap["main.NewTxListCtx"] = NewTxListCtx
	qt.ItfMap["main.LoginContext_ITF"] = LoginContext{}
	qt.FuncMap["main.NewLoginContext"] = NewLoginContext
	qt.ItfMap["main.TxListAccountsModel_ITF"] = TxListAccountsModel{}
	qt.FuncMap["main.NewTxListAccountsModel"] = NewTxListAccountsModel
	qt.ItfMap["main.TxListModel_ITF"] = TxListModel{}
	qt.FuncMap["main.NewTxListModel"] = NewTxListModel
}
