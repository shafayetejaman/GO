package main

type formatter interface {
	format() string
}

type plainText struct {
	message string
}
type bold struct {
	message string
}
type code struct {
	message string
}

func (p plainText) format() string {
	return p.message
}
func (p bold) format() string {
	return "**" + p.message + "**"
}
func (p code) format() string {
	return "`" + p.message + "`"
}

// Don't Touch below this line

func sendMessage(format formatter) string {
	return format.format() // Adjusted to call Format without an argument
}
