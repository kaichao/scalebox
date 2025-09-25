# side-car

## 测试1(缺省脚本)

```sh
echo 0 | scalebox app run
```

## 测试2(定制脚本)
```sh
ACTION_RUN=/app/bin/run-custom.sh \
ACTION_SETUP=/app/bin/setup-custom.sh \
ACTION_TEARDOWN=/app/bin/teardown-custom.sh \
ACTION_CHECK=/app/bin/check-custom.sh \
scalebox app create
```
