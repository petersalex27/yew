package repl

import "gopkg.in/yaml.v3"

type Args struct {
	Import []string `yaml:"import"`
	// nil if not literate, otherwise empty string for no output and a string for output file
	Literate *string `yaml:"literate"`
	Output string `yaml:"output"`
	IgnoreComments bool `yaml:"ignore-comments"`
}

func makeDefault() Args {
	return Args{
		Import: []string{},
		Literate: nil,
		Output: "",
		IgnoreComments: true,
	}
}

func ParseArgs(args string) (Args, error) {
	out := makeDefault()
	err := yaml.Unmarshal([]byte(args), out)
	if err != nil {
		return makeDefault(), err
	}
	return out, nil
}