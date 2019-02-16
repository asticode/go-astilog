package astilog

import "github.com/sirupsen/logrus"

type withFieldHook struct {
	k string
	v interface{}
}

func newWithFieldHook(k string, v interface{}) *withFieldHook {
	return &withFieldHook{
		k: k,
		v: v,
	}
}

func (h *withFieldHook) Fire(e *logrus.Entry) error {
	if h.v != nil {
		e.Data[h.k] = h.v
	}
	return nil
}

func (h *withFieldHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
