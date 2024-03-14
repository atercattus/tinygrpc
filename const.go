package http2

var (
	httpHeaderSep = []byte("\r\n")

	// httpHeaderConnection    = []byte("Connection:")
	// httpHeaderUpgrade       = []byte("Upgrade:")
	// httpHeaderHTTP2Settings = []byte("HTTP2-Settings:")

	connectionPreface          = []byte("PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n")
	httpSwitchingProtoResponse = []byte("HTTP/1.1 101 Switching Protocols\r\nConnection: Upgrade\r\nUpgrade: h2c\r\n\r\n")
)

const (
	RecvBufferSize = 16 * 1024

	// https://httpwg.org/specs/rfc7540.html#SettingValues
	initialSettingsMaxConcurrentStreams = 250
	initialHeaderTableSize              = 4096 // Go client sends 10485760
	initialMaxHeaderListSize            = 1024 * 1024
	initialSettingsMaxFrameSize         = 16 * 1024
	initialSettingsInitialWindowSize    = 1024 * 1024 // Go client sends 4194304

	httpSP = ' ' // SP = <US-ASCII SP, space (32)>
)
