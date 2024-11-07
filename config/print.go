package config

import (
	"bytes"
	"io"

	"gopkg.in/yaml.v3"
)

var (
	_Prefix = "[INFO] Config: \n\t"
)

// Print self. Unsafe function, use only on the startup.
func (c Config) Print(to io.Writer) {
	ybytes, err := yaml.Marshal(c)
	if err != nil {
		return
	}
	buf := bytes.Buffer{}

	_, _ = buf.WriteString(_Prefix)
	for _, b := range ybytes {
		if b == '\n' {
			_, _ = buf.Write([]byte{'\n', '\t'})

			continue
		}
		_ = buf.WriteByte(b)
	}
	_, _ = to.Write(buf.Bytes())
}
