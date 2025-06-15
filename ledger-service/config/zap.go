package config

import "go.uber.org/zap"

func ConfigureZap() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(logger)
	return nil
}
