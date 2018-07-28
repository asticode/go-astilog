package astilog

import "github.com/sirupsen/logrus"

type appNameHook struct {
	appName string
}

func newAppNameHook(appName string) *appNameHook {
	return &appNameHook{appName: appName}
}

func (h *appNameHook) Fire(e *logrus.Entry) error {
	e.Data["app_name"] = h.appName
	return nil
}

func (h *appNameHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
