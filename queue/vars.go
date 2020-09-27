package queue

var Params = struct {
	PoolSize              int //每个pool缓冲大小
	MaxPoolNum            int //可创建最大的缓冲区的数量，0为不限制
	PoolTimeOut           int //针对空闲的pool等待多久移出缓冲列表(毫秒)
	MaxProcessMessageNum  int //每个协程处理消息的数据
	MaxWaitCollectMessage int //最大等待收集消息的时间（毫秒），超过这个时间返回空的消息处理列表
	NextWaitReadMessage   int //下一次协程处理消息的等待时间（毫秒）
}{
	PoolSize:              512,
	MaxPoolNum:            1024,
	PoolTimeOut:           60 * 1000,
	MaxProcessMessageNum:  1000,
	MaxWaitCollectMessage: 1000,
	NextWaitReadMessage:   800,
}
