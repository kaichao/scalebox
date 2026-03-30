# side-car

边车模式的用户定制脚本：
- 初始化脚本：ACTION_SETUP
- 退出脚本：ACTION_TEARDOWN
- 任务运行脚本：ACTION_RUN
- 准入检测脚本：ACTION_CHECK

## 测试1(缺省脚本)

```sh
echo 0 | scalebox run
```

## 测试2(定制脚本)
```sh
ACTION_RUN=/app/bin/run-custom.sh \
ACTION_SETUP=/app/bin/setup-custom.sh \
ACTION_TEARDOWN=/app/bin/teardown-custom.sh \
ACTION_CHECK=/app/bin/check-custom.sh \
scalebox run
```
