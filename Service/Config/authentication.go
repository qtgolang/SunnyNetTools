package Config

func (f *config) initAuthentication() {
	if f.Authentication == nil {
		f.Authentication = make(map[int]*AuthenticationInfo)
	}
}
