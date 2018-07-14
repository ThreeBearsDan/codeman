### FAQ:
* 在创建consumer时，需要传入`go-nsq.NewConfig()`这样一个配置，而且只能通过这种方式传入，不可以自己定一个`Config`结构体常量。配置中有一个`max-in-flight`参数，这个参数什么意思？

    这个参数其实是`consumer`端做流量控制，当然也就相当于客户端的一个并发开关。`consumer`通过`nsqlookupd`查询到自己订阅的`topic/channel`的`nsqd`之后与`nsqd`建立`TCP`连接，在发送给`nsqd`的请求中包含`RDY n`这样的命令，这条命令的意思是告诉`nsqd`一次发送过来`n`条消息，我可以全部处理。当`consumer`处理的消息还剩下接近25%*`n`时，这时`consumer`会再次发送`RDY n`这个请求，告诉服务端，再来`n`条。所以这里的`n`可以认为是`max-in-flight`，当然是针对一个`nsqd`节点。如果一个`consumer`消费多个`nsqd`节点上的消息，则发送给每个节点的`RDY n`，`n`的值将变为`max-in-flight`消息的一半，如果是奇数可能不会均分。所以代码中如果尝试将`max-in-flight`设置为`0`，便会发现`nsqd`不会向`consumer`推送消息。