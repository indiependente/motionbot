package muxer

// Muxer muxes an input file into an output file.
type Muxer interface {
	Mux(string) (string, error)
}
