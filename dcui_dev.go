// +build dev

package dcui

import "net/http"

var WebUI http.FileSystem = http.Dir("webapp/build")
