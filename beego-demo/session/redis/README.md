单独使用了`beego`的`session`模块实现`session`，底层使用`redis`作为存储。数据在`redis`以`session id`作为一个`string`类型的键，每次请求之后都会从这个键中取出值然后decode到一个`map`中，这个`map`和每个
`session`唯一绑定，所以对于无状态的`http`协议，在服务器端通过`session`的方式将客户端的信息存储在服务端，使得用户会有更好的上网体验。

