package util

import (
	"fmt"
	"io"
	"os"

	"github.com/petersalex27/yew/api"
)

type (
	fileSource struct {
		path string
		content string
	}

	emptySource struct {}

	stringSource string

	// content is hidden to ensure immutability
	bytesSource struct {
		content []byte
	}

	freeSource struct {
		path string
		content string
	}
)

func ExposeSource(source api.Source) string {
	return fmt.Sprintf("Source{path: %q, content: %q}", source.Path(), source.String())
}

func (fs fileSource) String() string { return fs.content }

func (fs fileSource) Path() string { return fs.path }

func (es emptySource) String() string { return "" }

func (es emptySource) Path() string { return "<empty>" }

func (ss stringSource) String() string { return string(ss) }

func (ss stringSource) Path() string { return "<string>" }

func (bs bytesSource) String() string { return string(bs.content) }

func (bs bytesSource) Path() string { return "<bytes>" }

func (fs freeSource) String() string { return fs.content }

func (fs freeSource) Path() string { return fs.path }

func FileSource(path string) (api.Source, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	var content []byte

	content, err = io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return fileSource{path: path, content: string(content)}, nil
}

func EmptySource() api.Source { return emptySource{} }

func StringSource(content string) api.Source { return stringSource(content) }

func BytesSource(content []byte) api.Source { return bytesSource{content: content} }

func FreeSource(path, content string) api.Source { return freeSource{path: path, content: content} }