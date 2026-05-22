package Config

func (f *config) initHTTPSProto() {
	if Config.HTTPSProto == "" {
		Config.HTTPSProto = "{\"ConnectionFlow\":12517377,\"HeaderPriority\":{\"StreamDep\":13,\"Exclusive\":false,\"Weight\":41},\"Priorities\":[{\"PriorityParam\":{\"StreamDep\":0,\"Exclusive\":false,\"Weight\":200},\"StreamID\":3},{\"PriorityParam\":{\"StreamDep\":0,\"Exclusive\":false,\"Weight\":100},\"StreamID\":5},{\"PriorityParam\":{\"StreamDep\":0,\"Exclusive\":false,\"Weight\":0},\"StreamID\":7},{\"PriorityParam\":{\"StreamDep\":7,\"Exclusive\":false,\"Weight\":0},\"StreamID\":9},{\"PriorityParam\":{\"StreamDep\":3,\"Exclusive\":false,\"Weight\":0},\"StreamID\":11},{\"PriorityParam\":{\"StreamDep\":0,\"Exclusive\":false,\"Weight\":240},\"StreamID\":13}],\"PseudoHeaderOrder\":[\":method\",\":path\",\":authority\",\":scheme\"],\"Settings\":{\"1\":65536,\"4\":131072,\"5\":16384},\"SettingsOrder\":[1,4,5]}"
		Config.RandomJa3 = true
	}
}
