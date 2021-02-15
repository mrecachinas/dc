// +build dev

package ui

import "net/http"

var WebUI http.FileSystem = http.Dir("webapp/build")
