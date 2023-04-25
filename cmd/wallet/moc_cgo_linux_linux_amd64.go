//go:build !ubports
// +build !ubports

package main

/*
#cgo CFLAGS: -pipe -O2 -Wall -W -D_REENTRANT -fPIC -DQT_NO_DEBUG -DQT_WIDGETS_LIB -DQT_QUICK_LIB -DQT_GUI_LIB -DQT_QML_LIB -DQT_NETWORK_LIB -DQT_CORE_LIB
#cgo CXXFLAGS: -pipe -O2 -Wall -W -D_REENTRANT -fPIC -DQT_NO_DEBUG -DQT_WIDGETS_LIB -DQT_QUICK_LIB -DQT_GUI_LIB -DQT_QML_LIB -DQT_NETWORK_LIB -DQT_CORE_LIB
#cgo CXXFLAGS: -I../../cmd -I. -I../../vendor/github.com/therecipe/env_linux_amd64_513/5.13.0/gcc_64/include -I../../vendor/github.com/therecipe/env_linux_amd64_513/5.13.0/gcc_64/include/QtWidgets -I../../vendor/github.com/therecipe/env_linux_amd64_513/5.13.0/gcc_64/include/QtQuick -I../../vendor/github.com/therecipe/env_linux_amd64_513/5.13.0/gcc_64/include/QtGui -I../../vendor/github.com/therecipe/env_linux_amd64_513/5.13.0/gcc_64/include/QtQml -I../../vendor/github.com/therecipe/env_linux_amd64_513/5.13.0/gcc_64/include/QtNetwork -I../../vendor/github.com/therecipe/env_linux_amd64_513/5.13.0/gcc_64/include/QtCore -I. -isystem /usr/include/libdrm -I../../vendor/github.com/therecipe/env_linux_amd64_513/5.13.0/gcc_64/mkspecs/linux-g++
#cgo LDFLAGS: -O1 -Wl,-rpath,/home/parallels/go/src/github.com/pavelkrolevets/MIR/vendor/github.com/therecipe/env_linux_amd64_513/5.13.0/gcc_64/lib
#cgo LDFLAGS:  /home/parallels/go/src/github.com/pavelkrolevets/MIR/vendor/github.com/therecipe/env_linux_amd64_513/5.13.0/gcc_64/lib/libQt5Widgets.so /home/parallels/go/src/github.com/pavelkrolevets/MIR/vendor/github.com/therecipe/env_linux_amd64_513/5.13.0/gcc_64/lib/libQt5Quick.so /home/parallels/go/src/github.com/pavelkrolevets/MIR/vendor/github.com/therecipe/env_linux_amd64_513/5.13.0/gcc_64/lib/libQt5Gui.so /home/parallels/go/src/github.com/pavelkrolevets/MIR/vendor/github.com/therecipe/env_linux_amd64_513/5.13.0/gcc_64/lib/libQt5Qml.so /home/parallels/go/src/github.com/pavelkrolevets/MIR/vendor/github.com/therecipe/env_linux_amd64_513/5.13.0/gcc_64/lib/libQt5Network.so /home/parallels/go/src/github.com/pavelkrolevets/MIR/vendor/github.com/therecipe/env_linux_amd64_513/5.13.0/gcc_64/lib/libQt5Core.so -lGL -lpthread
#cgo CFLAGS: -Wno-unused-parameter -Wno-unused-variable -Wno-return-type
#cgo CXXFLAGS: -Wno-unused-parameter -Wno-unused-variable -Wno-return-type
*/
import "C"
