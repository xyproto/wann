package wann

// Config is a struct that is used when initializing new Network structs.
// The idea is that referring to fields by name is more explicit, and that it can
// be re-used in connection with having a configuration file, in the future.
type Config struct {
	Inputs          int
	ConnectionRatio float64
	SharedWeight    float64
}
