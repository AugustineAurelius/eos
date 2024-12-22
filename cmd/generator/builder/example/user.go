package example

//go:generate github.com/AugustineAurelius/eos generate builder
type User struct {
	Name    string
	Surname string
}
