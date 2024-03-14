package http2

import "fmt"

type (
	SettingsType uint8
)

const (
	SettingsHeaderTableSize      = SettingsType(0x1)
	SettingsEnablePush           = SettingsType(0x2)
	SettingsMaxConcurrentStreams = SettingsType(0x3)
	SettingsInitialWindowSize    = SettingsType(0x4)
	SettingsMaxFrameSize         = SettingsType(0x5)
	SettingsMaxHeaderListSize    = SettingsType(0x6)
)

func (s SettingsType) String() string {
	switch s {
	case SettingsHeaderTableSize:
		return `SETTINGS_HEADER_TABLE_SIZE`
	case SettingsEnablePush:
		return `SETTINGS_ENABLE_PUSH`
	case SettingsMaxConcurrentStreams:
		return `SETTINGS_MAX_CONCURRENT_STREAMS`
	case SettingsInitialWindowSize:
		return `SETTINGS_INITIAL_WINDOW_SIZE`
	case SettingsMaxFrameSize:
		return `SETTINGS_MAX_FRAME_SIZE`
	case SettingsMaxHeaderListSize:
		return `SETTINGS_MAX_HEADER_LIST_SIZE`
	default:
		return fmt.Sprintf(`SettingsType(%d)`, s)
	}
}
