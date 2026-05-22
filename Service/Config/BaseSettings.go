package Config

func (f *config) initBaseSettings() {
	if f.Port == 0 {
		f.Port = 2025
		f.DisableCache = true
	}
	if f.LimitRequestSize < 10240 {
		f.LimitRequestSize = 1024000
	}
}
