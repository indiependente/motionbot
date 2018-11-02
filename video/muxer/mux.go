package muxer

type Muxer interface {
	Mux(string) (string, error)
}
