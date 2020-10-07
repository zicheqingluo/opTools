package util

/*
@Time : 2020/10/07
@Author : yangxiaokang
@File : retry.go
*/

import (
	zlog "common/zaplog"
	"time"
)

func Retry(attempts int, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if s, ok := err.(stop); ok {
			zlog.Warn("不需要重试:%s", s.error)
			return s.error
		}

		if attempts--; attempts > 0 {
			zlog.Warn("retry func error: %s. attemps #%d after %s.", err.Error(), attempts, sleep)
			time.Sleep(sleep)
			return Retry(attempts, 2*sleep, fn)
		}
		zlog.Warn("达到最大重试次数")
		return err
	}
	return nil
}

type stop struct {
	error
}

func NoRetryError(err error) stop {
	return stop{err}
}
