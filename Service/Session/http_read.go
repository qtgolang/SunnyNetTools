package Session

// ReadLocked 在断点 Wait 期间回调仍可能占用 Mutex；MCP/UI 只读时若拿不到锁且处于 IsWait，则无锁读取快照字段。
func (s *HttpSession) ReadLocked(fn func()) {
	if s == nil || fn == nil {
		return
	}
	if s.TryLock() {
		defer s.Unlock()
		fn()
		return
	}
	if s.IsWait() {
		fn()
		return
	}
	s.Lock()
	defer s.Unlock()
	fn()
}
