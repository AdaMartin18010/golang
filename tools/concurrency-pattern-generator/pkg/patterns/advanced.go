// Package patterns - 高级模式
package patterns

import "fmt"

// GenerateActorModel 生成Actor模型
func GenerateActorModel(pkg string) string {
	return fmt.Sprintf("package %s\n\nimport \"context\"\n\n// Message Actor消息\ntype Message struct {\n\tType string\n\tData interface{}\n}\n\n// Actor Actor模型\ntype Actor struct {\n\tmailbox chan Message\n}\n\n// NewActor 创建Actor\nfunc NewActor() *Actor {\n\treturn &Actor{\n\t\tmailbox: make(chan Message, 10),\n\t}\n}\n\n// Send 发送消息\nfunc (a *Actor) Send(msg Message) {\n\ta.mailbox <- msg\n}\n\n// Start 启动Actor\nfunc (a *Actor) Start(ctx context.Context) {\n\tfor {\n\t\tselect {\n\t\tcase msg := <-a.mailbox:\n\t\t\t// Handle message\n\t\t\t_ = msg\n\t\tcase <-ctx.Done():\n\t\t\treturn\n\t\t}\n\t}\n}\n", pkg)
}

// GenerateFuturePromise 生成Future/Promise模式
func GenerateFuturePromise(pkg string) string {
	return fmt.Sprintf("package %s\n\n// Future 异步结果\ntype Future struct {\n\tresult chan interface{}\n\terr    chan error\n}\n\n// NewFuture 创建Future\nfunc NewFuture(fn func() (interface{}, error)) *Future {\n\tf := &Future{\n\t\tresult: make(chan interface{}, 1),\n\t\terr:    make(chan error, 1),\n\t}\n\t\n\tgo func() {\n\t\tres, err := fn()\n\t\tif err != nil {\n\t\t\tf.err <- err\n\t\t} else {\n\t\t\tf.result <- res\n\t\t}\n\t}()\n\t\n\treturn f\n}\n\n// Get 获取结果\nfunc (f *Future) Get() (interface{}, error) {\n\tselect {\n\tcase res := <-f.result:\n\t\treturn res, nil\n\tcase err := <-f.err:\n\t\treturn nil, err\n\t}\n}\n", pkg)
}

// GenerateMapReduce 生成MapReduce模式
func GenerateMapReduce(pkg string) string {
	return fmt.Sprintf("package %s\n\nimport \"sync\"\n\n// MapFunc map函数类型\ntype MapFunc func(interface{}) interface{}\n\n// ReduceFunc reduce函数类型\ntype ReduceFunc func(interface{}, interface{}) interface{}\n\n// MapReduce MapReduce模式\nfunc MapReduce(data []interface{}, mapper MapFunc, reducer ReduceFunc, initial interface{}) interface{} {\n\t// Map phase\n\tmapped := make([]interface{}, len(data))\n\tvar wg sync.WaitGroup\n\t\n\tfor i, item := range data {\n\t\twg.Add(1)\n\t\tgo func(idx int, val interface{}) {\n\t\t\tdefer wg.Done()\n\t\t\tmapped[idx] = mapper(val)\n\t\t}(i, item)\n\t}\n\twg.Wait()\n\t\n\t// Reduce phase\n\tresult := initial\n\tfor _, item := range mapped {\n\t\tresult = reducer(result, item)\n\t}\n\t\n\treturn result\n}\n", pkg)
}

// GeneratePubSub 生成发布订阅模式
func GeneratePubSub(pkg string) string {
	return fmt.Sprintf("package %s\n\nimport \"sync\"\n\n// PubSub 发布订阅\ntype PubSub struct {\n\tmu          sync.RWMutex\n\tsubscribers map[string][]chan interface{}\n}\n\n// NewPubSub 创建PubSub\nfunc NewPubSub() *PubSub {\n\treturn &PubSub{\n\t\tsubscribers: make(map[string][]chan interface{}),\n\t}\n}\n\n// Subscribe 订阅主题\nfunc (ps *PubSub) Subscribe(topic string) <-chan interface{} {\n\tps.mu.Lock()\n\tdefer ps.mu.Unlock()\n\t\n\tch := make(chan interface{}, 10)\n\tps.subscribers[topic] = append(ps.subscribers[topic], ch)\n\treturn ch\n}\n\n// Publish 发布消息\nfunc (ps *PubSub) Publish(topic string, msg interface{}) {\n\tps.mu.RLock()\n\tdefer ps.mu.RUnlock()\n\t\n\tfor _, ch := range ps.subscribers[topic] {\n\t\tselect {\n\t\tcase ch <- msg:\n\t\tdefault:\n\t\t}\n\t}\n}\n", pkg)
}

// GenerateSessionTypes 生成Session Types模式
func GenerateSessionTypes(pkg string) string {
	return fmt.Sprintf("package %s\n\n// Session 会话类型\ntype Session struct {\n\tsend    chan<- interface{}\n\treceive <-chan interface{}\n}\n\n// NewSession 创建会话\nfunc NewSession() (*Session, *Session) {\n\tch1 := make(chan interface{})\n\tch2 := make(chan interface{})\n\t\n\ts1 := &Session{send: ch1, receive: ch2}\n\ts2 := &Session{send: ch2, receive: ch1}\n\t\n\treturn s1, s2\n}\n\n// Send 发送数据\nfunc (s *Session) Send(data interface{}) {\n\ts.send <- data\n}\n\n// Receive 接收数据\nfunc (s *Session) Receive() interface{} {\n\treturn <-s.receive\n}\n", pkg)
}
