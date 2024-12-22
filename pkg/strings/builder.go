package strings

import "strings"

type Builder struct {
	inner strings.Builder
}

func (b *Builder) WriteString(data string) *Builder {
	b.inner.WriteString(data)
	return b
}

func (b *Builder) WriteByte(data byte) *Builder {
	b.inner.WriteByte(data)
	return b
}

func (b *Builder) WriteEnter() *Builder {
	return b.WriteByte('\n')
}

func (b *Builder) WriteStringWithEnter(data string) *Builder {
	return b.WriteString(data + "\n")
}

func (b *Builder) String() string {
	return b.inner.String()
}

func (b *Builder) Bytes() []byte {
	return StringToBytes(b.inner.String())
}
