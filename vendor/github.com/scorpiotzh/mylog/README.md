# mylog

* zap log > stdout

```go

log := NewLogger("test", LevelDebug)

log.Debug("aaaaa")
log.Debugf("aaa %s aaa", "bbb")

log.Info("aaaaa")
log.Infof("aaa %s aaa", "bbb")

log.Warn("aaaaa")
log.Warnf("aaa %s aaa", "bbb")

log.Error("aaaaa")
log.Errorf("aaa %s aaa", "bbb")

```

* log > file
```go

log := NewLoggerDefault("test", LevelDebug, nil)

log.Debug("aaaaa")
log.Debugf("aaa %s aaa", "bbb")

```