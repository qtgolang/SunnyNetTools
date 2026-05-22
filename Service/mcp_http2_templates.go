package Service

import (
	"strings"
)

// http2FingerprintTemplates 与前端 HTTPSProto.vue options 一致。
var http2FingerprintTemplates = []map[string]any{
	{"name": "Firefox", "label": "Firefox", "config": "{\"ConnectionFlow\":12517377,\"HeaderPriority\":{\"StreamDep\":13,\"Exclusive\":false,\"Weight\":41},\"Priorities\":[{\"PriorityParam\":{\"StreamDep\":0,\"Exclusive\":false,\"Weight\":200},\"StreamID\":3},{\"PriorityParam\":{\"StreamDep\":0,\"Exclusive\":false,\"Weight\":100},\"StreamID\":5},{\"PriorityParam\":{\"StreamDep\":0,\"Exclusive\":false,\"Weight\":0},\"StreamID\":7},{\"PriorityParam\":{\"StreamDep\":7,\"Exclusive\":false,\"Weight\":0},\"StreamID\":9},{\"PriorityParam\":{\"StreamDep\":3,\"Exclusive\":false,\"Weight\":0},\"StreamID\":11},{\"PriorityParam\":{\"StreamDep\":0,\"Exclusive\":false,\"Weight\":240},\"StreamID\":13}],\"PseudoHeaderOrder\":[\":method\",\":path\",\":authority\",\":scheme\"],\"Settings\":{\"1\":65536,\"4\":131072,\"5\":16384},\"SettingsOrder\":[1,4,5]}"},
	{"name": "Opera", "label": "Opera", "config": "{\"ConnectionFlow\":15663105,\"HeaderPriority\":null,\"Priorities\":null,\"PseudoHeaderOrder\":[\":method\",\":authority\",\":scheme\",\":path\"],\"Settings\":{\"1\":65536,\"3\":1000,\"4\":6291456,\"6\":262144},\"SettingsOrder\":[1,3,4,6]}"},
	{"name": "Safari_IOS_17_0", "label": "Safari_IOS_17_0", "config": "{\"ConnectionFlow\":10485760,\"HeaderPriority\":null,\"Priorities\":null,\"PseudoHeaderOrder\":[\":method\",\":scheme\",\":path\",\":authority\"],\"Settings\":{\"2\":0,\"3\":100,\"4\":2097152},\"SettingsOrder\":[2,4,3]}"},
	{"name": "Safari_IOS_16_0", "label": "Safari_IOS_16_0", "config": "{\"ConnectionFlow\":10485760,\"HeaderPriority\":null,\"Priorities\":null,\"PseudoHeaderOrder\":[\":method\",\":scheme\",\":path\",\":authority\"],\"Settings\":{\"3\":100,\"4\":2097152},\"SettingsOrder\":[4,3]}"},
	{"name": "Chrome_117_120_124", "label": "Chrome_117_120_124", "config": "{\"ConnectionFlow\":10485760,\"HeaderPriority\":null,\"Priorities\":null,\"PseudoHeaderOrder\":[\":method\",\":scheme\",\":path\",\":authority\"],\"Settings\":{\"3\":100,\"4\":4194304},\"SettingsOrder\":[4,3]}"},
	{"name": "Chrome_106_116", "label": "Chrome_106_116", "config": "{\"ConnectionFlow\":15663105,\"HeaderPriority\":null,\"Priorities\":null,\"PseudoHeaderOrder\":[\":method\",\":authority\",\":scheme\",\":path\"],\"Settings\":{\"1\":65536,\"2\":0,\"4\":6291456,\"6\":262144},\"SettingsOrder\":[1,2,4,6]}"},
	{"name": "Chrome_103_105", "label": "Chrome_103_105", "config": "{\"ConnectionFlow\":15663105,\"HeaderPriority\":null,\"Priorities\":null,\"PseudoHeaderOrder\":[\":method\",\":authority\",\":scheme\",\":path\"],\"Settings\":{\"1\":65536,\"3\":1000,\"4\":6291456,\"6\":262144},\"SettingsOrder\":[1,3,4,6]}"},
}

func findHTTP2Template(name string) (map[string]any, bool) {
	name = strings.TrimSpace(name)
	for _, t := range http2FingerprintTemplates {
		if strings.TrimSpace(t["name"].(string)) == name {
			return t, true
		}
	}
	return nil, false
}
