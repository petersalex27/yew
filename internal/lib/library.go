package lib

import "github.com/petersalex27/yew/internal/pkg"

// TODO: change this to the actual path
const YEW_PATH = "../libs"

type (
	Library struct {
		// path relative to YEW_PATH 
		path string
		packages map[string]pkg.Loader
	}
)