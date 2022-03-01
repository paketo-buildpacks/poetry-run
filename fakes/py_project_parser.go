package fakes

import "sync"

type PyProjectParser struct {
	ParseCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			String string
		}
		Returns struct {
			String string
			Error  error
		}
		Stub func(string) (string, error)
	}
}

func (f *PyProjectParser) Parse(param1 string) (string, error) {
	f.ParseCall.mutex.Lock()
	defer f.ParseCall.mutex.Unlock()
	f.ParseCall.CallCount++
	f.ParseCall.Receives.String = param1
	if f.ParseCall.Stub != nil {
		return f.ParseCall.Stub(param1)
	}
	return f.ParseCall.Returns.String, f.ParseCall.Returns.Error
}
