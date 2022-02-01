package wsnterminal

type wsnTerminal struct {
	port string
	// TODO: handlers to call on a certain type of message
}

type WsnTerminal interface {
}

func New(port string) WsnTerminal {
	wt := wsnTerminal{
		port: port,
	}
	return wt
}

func start() {
}

func read() {
}

func write() {
}

// USe go usb: https://pkg.go.dev/github.com/karalabe/usb#section-documentation
