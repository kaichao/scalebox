# 4. 常见问答集

Frequently Asked Questions

## 1. local集群的hello-scalebox应用不能正常运行，如何检查？
- 检查本地节点的sshd的端口号是否为缺省的22，若不是，则在server/local/mycluster.yaml文件中设置该端口号，具体在parameters中加上`port=<port_number>`
- 用```hostname -i```检查本地ip地址是否正确设置？若不正确，则用```hostname -I```列出所有本地IP地址，再在defs.mk中设置变量```LOCAL_ADDR```或```LOCAL_IP_INDEX```

## 2. 

