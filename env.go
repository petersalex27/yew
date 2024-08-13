package main

import (
	"os"

	"github.com/petersalex27/yew/errors"
)

const pathErrorMessage string = "cannot create yew directory because the home directory cannot be located"

// TODO: finish
func writeLibrary(path string) (e error) {
	panic("TODO: implement")
	return nil
}

// TODO: finish
func createYewDir() (path string, e error) {
	// find home dir
	path, e = os.UserHomeDir()
	if e != nil {
		// error: can't locate home dir
		return path, errors.MakeError("OS", pathErrorMessage)
	}

	// append yew to path
	if len(path) > 0 && path[len(path)-1] != '/' {
		path = path + "/.yew"
	} else {
		path = path + ".yew"
	}

	// create ~/.yew dir
	if e = os.Mkdir(path, os.ModeDir); e != nil {
		if pe, ok := e.(*os.PathError); ok {
			// error: some kind of path error
			return path, errors.MakeError("OS", pe.Err.Error())
		}
		// error: some other kind of error
		return path, errors.MakeError("OS", e.Error())
	}
	return path, nil
}