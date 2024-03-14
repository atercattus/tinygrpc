package http2

type (
	Settings struct { // page 40
		HeaderTableSize      uint32
		EnablePush           bool
		MaxConcurrentStreams uint32
		InitialWindowSize    uint32
		MaxFrameSize         uint32
		MaxHeaderListSize    uint32
	}
)

func GetDefaultSettings() Settings {
	return Settings{
		HeaderTableSize:      initialHeaderTableSize,
		EnablePush:           false,
		MaxConcurrentStreams: initialSettingsMaxConcurrentStreams,
		InitialWindowSize:    initialSettingsInitialWindowSize,
		MaxFrameSize:         initialSettingsMaxFrameSize,
		MaxHeaderListSize:    initialMaxHeaderListSize, //math.MaxInt32, // unlimited
	}
}

func (s *Settings) UpdateFromSettingsFrame(frame *SettingsFrame) {
	for _, param := range frame.Params {
		switch param.Id {
		case SettingsHeaderTableSize:
			s.HeaderTableSize = param.Value
		case SettingsEnablePush:
			s.EnablePush = param.Value != 0
		case SettingsMaxConcurrentStreams:
			s.MaxConcurrentStreams = param.Value
		case SettingsInitialWindowSize:
			s.InitialWindowSize = param.Value
		case SettingsMaxFrameSize:
			s.MaxFrameSize = param.Value
		case SettingsMaxHeaderListSize:
			s.MaxHeaderListSize = param.Value
		}
	}
}

func (s *Settings) ToSettingsFrame() *SettingsFrame {
	enablePush := uint32(0)
	if s.EnablePush {
		enablePush = 1
	}

	var sf SettingsFrame
	sf.Params = []SettingsFrameParam{
		{Id: SettingsMaxFrameSize, Value: s.MaxFrameSize},
		{Id: SettingsMaxConcurrentStreams, Value: s.MaxConcurrentStreams},
		{Id: SettingsMaxHeaderListSize, Value: s.MaxHeaderListSize},
		{Id: SettingsHeaderTableSize, Value: s.HeaderTableSize},
		{Id: SettingsInitialWindowSize, Value: s.InitialWindowSize},
		{Id: SettingsEnablePush, Value: enablePush},
	}
	return &sf
}
