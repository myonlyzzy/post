# 记一次grpc panic的问题

最近一个grpc的项目中遇到这样一个问题, 在grpc的rpc 中通过defer捕获了异常,但是程序还是异常退出。但是同样的问题我在grpc的helloworld 抛出异常却不会引起程序退出。

## 1. 使用 go-grpc-middleware 拦截所有rpc方法的panic .



```
opts := []grpc_recovery.Option{
   grpc_recovery.WithRecoveryHandler(recoverHandleFunc),
}
rpcServer := grpc.NewServer(grpc_middleware.WithUnaryServerChain(
   grpc_recovery.UnaryServerInterceptor(opts...),
),
   grpc_middleware.WithStreamServerChain(
      grpc_recovery.StreamServerInterceptor(opts...),
   ), )
 
 func recoverHandleFunc(p interface{}) (err error) {
    if p != nil {
		glog.KeeperLog.Errorln(errors.WithStack(errors.New("panic")))
	} 
	return nil
}

   
```

但是仍然会panic,看来问题不是在这里

## 2. 对比panic的栈

### 2.1 不panic的grpc server端 

通过panic抛出异常,在defer中捕获

 #### 2.1.1 不退出的panic的调用栈

```
ERROR: 2018/07/26 19:12:17 grpc: server failed to encode response:  rpc error: code = Internal desc = grpc: error while marshaling: proto: Marshal called with nil
2018/07/26 19:12:17 goroutine 78 [running]:
runtime/debug.Stack(0xc42004f918, 0x1406220, 0x1519740)
	/usr/local/go/src/runtime/debug/stack.go:24 +0xa7
main.(*server).SayHello.func1()
	/Users/myonlyzzy/go/src/github.com/myonlyzzy/go-exmaple/examples/helloworld/greeter_server/main.go:45 +0x48
panic(0x1406220, 0x1519740)
	/usr/local/go/src/runtime/panic.go:502 +0x229
main.(*server).SayHello(0x1787d70, 0x1520b40, 0xc4201b0fc0, 0xc4201b0ff0, 0x0, 0x0, 0x0)
	/Users/myonlyzzy/go/src/github.com/myonlyzzy/go-exmaple/examples/helloworld/greeter_server/main.go:48 +0x66
google.golang.org/grpc/examples/helloworld/helloworld._Greeter_SayHello_Handler(0x1425760, 0x1787d70, 0x1520b40, 0xc4201b0fc0, 0xc4201aef50, 0x0, 0x0, 0x0, 0x0, 0x0)
	/Users/myonlyzzy/go/src/google.golang.org/grpc/examples/helloworld/helloworld/helloworld.pb.go:158 +0x241
google.golang.org/grpc.(*Server).processUnaryRPC(0xc420190000, 0x1522c20, 0xc4200fc200, 0xc42021ed00, 0xc4200a3da0, 0x1759fa0, 0x0, 0x0, 0x0)
	/Users/myonlyzzy/go/src/google.golang.org/grpc/server.go:1011 +0x4fc
google.golang.org/grpc.(*Server).handleStream(0xc420190000, 0x1522c20, 0xc4200fc200, 0xc42021ed00, 0x0)
	/Users/myonlyzzy/go/src/google.golang.org/grpc/server.go:1249 +0x1318
google.golang.org/grpc.(*Server).serveStreams.func1.1(0xc4200a8640, 0xc420190000, 0x1522c20, 0xc4200fc200, 0xc42021ed00)
	/Users/myonlyzzy/go/src/google.golang.org/grpc/server.go:680 +0x9f
created by google.golang.org/grpc.(*Server).serveStreams.func1
	/Users/myonlyzzy/go/src/google.golang.org/grpc/server.go:678 +0xa1

ERROR: 2018/07/26 19:12:17 grpc: server failed to encode response:  rpc error: code = Internal desc = grpc: error while marshaling: proto: Marshal called with nil
```

#### 2.1.2程序退出的调用栈

```
	/usr/local/go/src/runtime/debug/stack.go:24 +0xa7
git.2dfire-inc.com/platform/zerodb/zero-keeper/server.(*Server).SwitchDB.func1()
	/Users/myonlyzzy/go/src/git.2dfire-inc.com/platform/zerodb/zero-keeper/server/rpcservice.go:301 +0x48
panic(0x19b2a20, 0x1c35cb0)
	/usr/local/go/src/runtime/panic.go:502 +0x229
git.2dfire-inc.com/platform/zerodb/zero-keeper/server.(*Server).SwitchDB(0xc4203e4820, 0x1c481e0, 0xc42001a630, 0xc4200bf440, 0x0, 0x0, 0x0)
	/Users/myonlyzzy/go/src/git.2dfire-inc.com/platform/zerodb/zero-keeper/server/rpcservice.go:304 +0x66
git.2dfire-inc.com/platform/zerodb/zero-common/zeroproto/pkg/keeper._Keeper_SwitchDB_Handler(0x1b65420, 0xc4203e4820, 0x1c481e0, 0xc42001a630, 0xc4200e6f50, 0x0, 0x0, 0x0, 0x0, 0x0)
	/Users/myonlyzzy/go/src/git.2dfire-inc.com/platform/zerodb/zero-common/zeroproto/pkg/keeper/keeper.pb.go:574 +0x241
google.golang.org/grpc.(*Server).processUnaryRPC(0xc4200f2000, 0x1c4c9c0, 0xc420452800, 0xc420432500, 0xc420030420, 0x2225010, 0x0, 0x0, 0x0)
	/Users/myonlyzzy/go/src/google.golang.org/grpc/server.go:1011 +0x4fc
google.golang.org/grpc.(*Server).handleStream(0xc4200f2000, 0x1c4c9c0, 0xc420452800, 0xc420432500, 0x0)
	/Users/myonlyzzy/go/src/google.golang.org/grpc/server.go:1249 +0x1318
google.golang.org/grpc.(*Server).serveStreams.func1.1(0xc4203aeed0, 0xc4200f2000, 0x1c4c9c0, 0xc420452800, 0xc420432500)
	/Users/myonlyzzy/go/src/google.golang.org/grpc/server.go:680 +0x9f
created by google.golang.org/grpc.(*Server).serveStreams.func1
	/Users/myonlyzzy/go/src/google.golang.org/grpc/server.go:678 +0xa1

panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x1536902]

goroutine 26 [running]:
git.2dfire-inc.com/platform/zerodb/zero-common/zeroproto/pkg/keeper.(*SwitchDBResponse).Size(0x0, 0xc42050d7d8)
	/Users/myonlyzzy/go/src/git.2dfire-inc.com/platform/zerodb/zero-common/zeroproto/pkg/keeper/keeper.pb.go:1316 +0x22
git.2dfire-inc.com/platform/zerodb/zero-common/zeroproto/pkg/keeper.(*SwitchDBResponse).Marshal(0x0, 0x1ab0a40, 0x0, 0x33e5470, 0x0, 0x3c7b3f501)
	/Users/myonlyzzy/go/src/git.2dfire-inc.com/platform/zerodb/zero-common/zeroproto/pkg/keeper/keeper.pb.go:1105 +0x2f
google.golang.org/grpc/encoding/proto.codec.Marshal(0x1ab0a40, 0x0, 0x0, 0xc42050d908, 0x100f1cd, 0xc42001e000, 0x1a21b00)
	/Users/myonlyzzy/go/src/google.golang.org/grpc/encoding/proto/proto.go:70 +0x19e
google.golang.org/grpc.encode(0x320c360, 0x22523a0, 0x1ab0a40, 0x0, 0x22523a0, 0x22329e0, 0xc420285760, 0xc420285740, 0xc4204484a8)
	/Users/myonlyzzy/go/src/google.golang.org/grpc/rpc_util.go:487 +0x61
google.golang.org/grpc.(*Server).sendResponse(0xc4200f2000, 0x1c4c9c0, 0xc420452800, 0xc420432500, 0x1ab0a40, 0x0, 0x0, 0x0, 0xc4200bcb6c, 0x0, ...)
	/Users/myonlyzzy/go/src/google.golang.org/grpc/server.go:830 +0x89
google.golang.org/grpc.(*Server).processUnaryRPC(0xc4200f2000, 0x1c4c9c0, 0xc420452800, 0xc420432500, 0xc420030420, 0x2225010, 0x0, 0x0, 0x0)
	/Users/myonlyzzy/go/src/google.golang.org/grpc/server.go:1036 +0x5df
google.golang.org/grpc.(*Server).handleStream(0xc4200f2000, 0x1c4c9c0, 0xc420452800, 0xc420432500, 0x0)
	/Users/myonlyzzy/go/src/google.golang.org/grpc/server.go:1249 +0x1318
google.golang.org/grpc.(*Server).serveStreams.func1.1(0xc4203aeed0, 0xc4200f2000, 0x1c4c9c0, 0xc420452800, 0xc420432500)
	/Users/myonlyzzy/go/src/google.golang.org/grpc/server.go:680 +0x9f
created by google.golang.org/grpc.(*Server).serveStreams.func1
	/Users/myonlyzzy/go/src/google.golang.org/grpc/server.go:678 +0xa1
```

### 2.1.3 异常的调用栈

第二次调用到了pb.go 文件中的Size方法

```
func (m *SwitchDBResponse) Size() (n int) {
	var l int
	_ = l
	if m.BasicResp != nil {
		l = m.BasicResp.Size()
		n += 1 + l + sovKeeper(uint64(l))
	}
	return n
}
```

第4行m.BasicResp 处产生异常,但是为什么第一次的rpc中不会调用到这个方法呢,而且我发现也没有生成Size方法。于是就要google了一把protobuf 文件生成到go文件的机制。

继续看上面的异常退出的函数调用栈.server.go 的1249行

```
if md, ok := srv.md[method]; ok {
		s.processUnaryRPC(t, stream, srv, md, trInfo)
		return
	}
```

调用到server.go 的1036行的sendResponse 方法

```
 fmt.Println("start send rpc response!")
if err := s.sendResponse(t, stream, reply, cp, opts, comp); err != nil {
   if err == io.EOF {
      // The entire stream is done (for unary RPC only).
      
```

sendResponse 调用encode函数对返回的数据进行编码.

```
func encode(c baseCodec, msg interface{}) ([]byte, error) {
   if msg == nil { // NOTE: typed nils will not be caught by this check
      return nil, nil
   }
   b, err := c.Marshal(msg)
   if err != nil {
      return nil, status.Errorf(codes.Internal, "grpc: error while marshaling: %v", err.Error())
   }
   if uint(len(b)) > math.MaxUint32 {
      return nil, status.Errorf(codes.ResourceExhausted, "grpc: message too large (%d bytes)", len(b))
   }
   return b, nil
}
```

encode函数中调用了codec的Marshal方法进行真正的序列化逻辑

```
func (codec) Marshal(v interface{}) ([]byte, error) {
   if pm, ok := v.(proto.Marshaler); ok {
      // object can marshal itself, no need for buffer
      println("start marshal...")
      return pm.Marshal()
   }
   println("cache marshal")
   cb := protoBufferPool.Get().(*cachedProtoBuffer)
   out, err := marshal(v, cb)

   // put back buffer and lose the ref to the slice
   cb.SetBuf(nil)
   protoBufferPool.Put(cb)
   return out, err
}
```

Marshal方法其实也很简单,先对返回的数据类型进行断言,如果是proto.Marshaler 类型就调用自己的Marshal方法，如果不是就再Marshal函数中进行序列化。是不是Marshaler类型返回数据结构上是否实现了Marshaler接口。

```
// Marshaler is the interface representing objects that can marshal themselves.
type Marshaler interface {
   Marshal() ([]byte, error)
}
```

原因找到,因为后一个会panic的grpc服务时使用了gogoproto的protoc-gen-gogofast .默认是给所有的类型生成marshal size 等方法,当然也可以使用gogoproto的options来关闭marshal 和unmarshal,但是gogoproto之所以比google原生的proto性能好,就是因为gogoproto生成的marshal unmarshal等方法要比原生的使用反射性能好。如果关闭了使用gogoproto的意义也不大了。



### 2.1.4  测试etcd 中的grpc服务  

写到这里原因找到了,要解决问题就使用原生的proto好了活着使用gogoproto关闭marshal.但是我想到gogoproto的文档中例举使用了很多gogoproto的开源go项目,第一个就是etcd,如果是这样etcd中会不会也存着这这样的问题，如果不存在它有是如何解决的呢。

####  2.1.4.1 etcd中为 putResponse生成的 Marshal方法和size方法

```
func (m *PutResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}
func (m *PutResponse) Size() (n int) {
	var l int
	_ = l
	if m.Header != nil {
		l = m.Header.Size()
		n += 1 + l + sovRpc(uint64(l))
	}
	if m.PrevKv != nil {
		l = m.PrevKv.Size()
		n += 1 + l + sovRpc(uint64(l))
	}
	return n
}
```

可以看到在size方法中有个m.header,如果put方法中产生panic,那么putResponse 肯定是nil,这里m.header 肯定会panic,如果你只在put方法中捕获异常，这里的panic肯定会引起程序崩溃。

####  2.1.4.2  修改etcd中put方法,抛出异常。 

我们在put方法中调用panic函数产生一个异常，然后再通过etcdctl调用put方法,居然发现etcd的进程直接挂掉退出去了.

```
func (s *kvServer) Put(ctx context.Context, r *pb.PutRequest) (*pb.PutResponse, error) {
	defer func(){
		if r:=recover();r!=nil{
			fmt.Println(string(debug.Stack()))
		}
	}()
	if err := checkPutRequest(r); err != nil {
		return nil, err
	}
	panic("test")
	resp, err := s.kv.Put(ctx, r)
	if err != nil {
		return nil, togRPCError(err)
	}

	s.hdr.fill(resp.Header)
	return resp, nil
}
```

```
6:37:28 etcd1 | github.com/coreos/etcd/etcdserver/etcdserverpb.(*PutResponse).Size(0x0, 0x1c87c80)
16:37:28 etcd1 |        /Users/myonlyzzy/go/src/github.com/coreos/etcd/etcdserver/etcdserverpb/rpc.pb.go:8102 +0x26
16:37:28 etcd1 | github.com/coreos/etcd/etcdserver/etcdserverpb.(*PutResponse).Marshal(0x0, 
16:37:28 etcd1 | 0x1c8d840, 0x0, 0x26967c8, 0x0, 0x1c87c01)
16:37:28 etcd1 |        /Users/myonlyzzy/go/src/github.com/coreos/etcd/etcdserver/etcdserverpb/rpc.pb.go:5135 +0x2f
16:37:28 etcd1 | github.com/coreos/etcd/vendor/github.com/gogo/protobuf/proto.Marshal(0x1c8d840, 0x0, 0x0, 0x1c8d840, 0x0, 0x3, 0x3)
16:37:28 etcd1 |        /Users/myonlyzzy/go/src/github.com/coreos/etcd/vendor/github.com/gogo/protobuf/proto/encode.go:233 +0x115
16:37:28 etcd1 | github.com/coreos/etcd/etcdserver/api/v3rpc.(*codec).Marshal(0x2212b20, 0x1af9d40, 0x0, 0xc4201c3d10, 0xc400000000, 0x1b17b80, 0xc426116a00, 0x1af9d40)
16:37:28 etcd1 |        /Users/myonlyzzy/go/src/github.com/coreos/etcd/etcdserver/api/v3rpc/codec.go:22 +0x60
16:37:28 etcd1 | github.com/coreos/etcd/vendor/google.golang.org/grpc.encode(0x1c8bf40, 0x2212b20, 0x1af9d40, 0x0, 0x0, 0x0, 0x0, 0x0, 0x100000001af9d40, 0x0, ...)
16:37:28 etcd1 |        /Users/myonlyzzy/go/src/github.com/coreos/etcd/vendor/google.golang.org/grpc/rpc_util.go:306 +0x25a
16:37:28 etcd1 | github.com/coreos/etcd/vendor/google.golang.org/grpc.(*Server).sendResponse(0xc4231d43c0, 0x1c97fc0, 0xc424bde420, 0xc420496a00, 0x1af9d40, 0x0, 0x0, 0x0, 0xc425d08bd0, 0x0, ...)
16:37:28 etcd1 |        /Users/myonlyzzy/go/src/github.com/coreos/etcd/vendor/google.golang.org/grpc/server.go:717 +0xad
16:37:28 etcd1 | github.com/coreos/etcd/vendor/google.golang.org/grpc.(*Server).processUnaryRPC(0xc4231d43c0, 0x1c97fc0, 0xc424bde420, 0xc420496a00, 0xc426a53770, 0x21e3c98, 0x0, 0x0, 0x0)
16:37:28 etcd1 |        /Users/myonlyzzy/go/src/github.com/coreos/etcd/vendor/google.golang.org/grpc/server.go:867 +0xf8a
16:37:28 etcd1 | github.com/coreos/etcd/vendor/google.golang.org/grpc.(*Server).handleStream(0xc4231d43c0, 0x1c97fc0, 0xc424bde420, 0xc420496a00, 0x0)
16:37:28 etcd1 |        /Users/myonlyzzy/go/src/github.com/coreos/etcd/vendor/google.golang.org/grpc/server.go:1040 +0x1318
16:37:28 etcd1 | github.com/coreos/etcd/vendor/google.golang.org/grpc.(*Server).serveStreams.func1.1(0xc424bd6a60, 0xc4231d43c0, 0x1c97fc0, 0xc424bde420, 0xc420496a00)
16:37:28 etcd1 |        /Users/myonlyzzy/go/src/github.com/coreos/etcd/vendor/google.golang.org/grpc/server.go:589 +0x9f
16:37:28 etcd1 | created by github.com/coreos/etcd/vendor/google.golang.org/grpc.(*Server).serveStreams.func1
16:37:28 etcd1 |        /Users/myonlyzzy/go/src/github.com/coreos/etcd/vendor/google.golang.org/grpc/server.go:587 +0xa1
16:37:28 etcd3 | 2018-07-31 16:37:28.806698 W | rafthttp: lost the TCP streaming connection with peer 8211f1d0f64f3269 (stream Message reader)
16:37:28 etcd3 | 2018-07-31 16:37:28.806718 E | rafthttp: failed to read 8211f1d0f64f3269 on stream Message (unexpected

```

可见etcd中的Put 方法在里面捕获异常也没有用,因为生成的size方法还是会panic。

## 3. 小结

通过上面的比较我们大概可以得出一下的结论.

* 使用google的原生的proto,在rpc方法中捕获异常是可以的。
* 如果每个方法都捕获嫌太麻烦,可以使用go-grpc-middleware grpc_recovery包。
* 如果想使用gogoproto并且打开marshal 会在生成的size 有可能会在产生panic,所以你在rpc捕获了，有可能程序还是会退出
* 如果使用gogoproto 并且不生成marshal,那还不如直接用原生的.
* 像etcd等开源项目中对rpc 方法中产生的panic 直接不处理,当然它可能认为不可能产生panic。