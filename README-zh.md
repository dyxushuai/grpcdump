# gRPC dump tool

## 示例

```sh
» sudo go run cmd/grpcdump.go -d lo0 -port 8085 -ip 127.0.0.1 -proto ./grpc_example/helloworld/helloworld/helloworld.proto
2018/06/06 18:20:04 Starting capture on device "lo0"
2018/06/06 18:20:04 reading in packets

REQUEST > 2018-06-06T18:20:06.731710593+08:00: 127.0.0.1:53913 ---> 127.0.0.1:8085
  HEADERS:
    :method = "POST"
    :scheme = "http"
    :path = "/helloworld.Greeter/SayHello"
    :authority = "127.0.0.1:8085"
    content-type = "application/grpc"
    user-agent = "grpc-go/1.13.0-dev"
    te = "trailers"
  BODY:
    {"name":"world","i18":"zh"}

RESPONSE > 2018-06-06T18:20:06.73193703+08:00: 127.0.0.1:53913 <--- 127.0.0.1:8085
  HEADERS:
    :status = "200"
    content-type = "application/grpc"
  BODY:
    {"message":"Hello, 你好"}

RESPONSE > 2018-06-06T18:20:06.731994357+08:00: 127.0.0.1:53913 <--- 127.0.0.1:8085
  HEADERS:
    grpc-status = "0"
    grpc-message = ""
```

## 需求

设计一个类似[tcpdump][]的工具，用于[gRPC][]请求的数据包抓包。

## [tcpdump][]工具简介

> TcpDump可以将网络中传送的数据包完全截获下来提供分析。它支持针对网络层、协议、主机、网络或端口的过滤，并提供and、or、not等逻辑语句来帮助你去掉无用的信息。 - 百度百科

1. tcpdump工作于网络协议的IP层，使用[AF_PACKET][]套接字抓取所有进出的IP包。
2. 通过解析IP包头过滤主机信息。
3. 通过解析TCP包头过滤端口信息。

## `grpcdump`设计目标

1. [x] 使用[AF_PACKET][]原始套接字抓取并过滤IP包。
2. [ ] 使用`keylog`解密TLS/SSL密文。
3. [x] 使用HTTP2协议栈解码。
4. [x] 动态加载[protobuf][]定义文件(.proto)，通过HTTP2 [:path][]头匹配对应Request和Response Message，解码body。
5. [x] 带缓存的格式化打印(colorful?)。

## `grpcdump`设计细节

### [AF_PACKET][]的golang实现[gopacket][]

[代码](pcap.go)

使用[gopacket][]监听网卡，设置[BPF][]规则，为了能够串联请求和返回包，必须要求提供gPRC服务的IP和端口。并在头部打印流信息，如下：

```sh
127.0.0.1:50422 ---> 127.0.0.1:8085
127.0.0.1:50422 <--- 127.0.0.1:8085
```

**http2包需要过滤掉[PRI][]/[代码](skip.go)** [issue](https://github.com/golang/go/issues/14141)

### TLS/SSL解密(暂未实现)

[gRPC][]客户端或者服务端需要使用[tls.Config.KeyLogWriter](tls.Config)记录[key_log][], grpcdump读取[key_log][]信息解密。

**只推荐在debug模式下记录[key_log][]**


### 动态加载[protobuf][]

[代码](protobuf.go)

指定[protobuf][]的定义文件,来实现运行时解析[protobuf][]编码的包。并打印序列化json，如下：

```sh
{"name":"world","i18":"zh"}
{"message":"Hello, 你好"}
```

**http2.DataFrame读出的body需要去空格！**
　
### 带缓存的格式化打印(readable)

输出可读性比较好的格式, 如下：

```sh
REQUEST > 2018-06-06T18:20:06.731710593+08:00: 127.0.0.1:53913 ---> 127.0.0.1:8085
  HEADERS:
    :method = "POST"
    :scheme = "http"
    :path = "/helloworld.Greeter/SayHello"
    :authority = "127.0.0.1:8085"
    content-type = "application/grpc"
    user-agent = "grpc-go/1.13.0-dev"
    te = "trailers"
  BODY:
    {"name":"world","i18":"zh"}

RESPONSE > 2018-06-06T18:20:06.73193703+08:00: 127.0.0.1:53913 <--- 127.0.0.1:8085
  HEADERS:
    :status = "200"
    content-type = "application/grpc"
  BODY:
    {"message":"Hello, 你好"}
```


[tcpdump]: http://www.tcpdump.org/
[gRPC]: https://grpc.io/
[AF_PACKET]: http://man7.org/linux/man-pages/man7/packet.7.html
[protobuf]: https://developers.google.com/protocol-buffers/
[gopacket]: https://github.com/google/gopacket
[:path]: https://tools.ietf.org/html/rfc3986#section-3.3
[BPF]: https://zh.wikipedia.org/wiki/BPF
[PRI]: https://http2.github.io/http2-spec/#rfc.section.11.6
[tls.Config]: https://godoc.org/crypto/tls#Config
[key_log]: https://developer.mozilla.org/en-US/docs/Mozilla/Projects/NSS/Key_Log_Format