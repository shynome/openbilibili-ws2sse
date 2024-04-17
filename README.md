## 简介

将 B 站直播开放平台的 WebSocket 转换为 EventStream, EventStream 不需要像 WebSocket 那样处理重连问题

## 使用

```sh
# 安装
go install github.com/shynome/openbilibili-ws2sse@v0.0.1
# 启动
# https://open-live.bilibili.com/open-manage 中个人资料tab获取参数值
openbilibili-ws2sse --key {{access_key_id}} --secret {{access_key_secret}} --appid {{appid}}
```

### 另起终端测试

```sh
curl 'http://127.0.0.1:7070/?IDCode={{IDCode}}'
```

## 碎碎念

不应该贪一时口快说要做的, 做这东西还挺花时间的
