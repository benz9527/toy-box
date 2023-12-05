# Why Timing Wheels

![时间轮动态图](https://pic2.zhimg.com/v2-3ecdb3b5a6be2d82d004bbd8371daafd_b.gif)

常规（传统）的任务队列调度，需要不断的轮询任务队列，然后进行调度。这种方式会存在两个问题：

- 轮询的开销
- 调度的精度

尽可能不损失调度精度的前提下，减少调度的开销。 比如使用全局唯一的时间轮来进行批量化调度，充分利用协程/多线程资源，避免存在多个不同的调度中心，导致资源浪费。
主要的表现是，时间轮可以大批量地管理各种延时任务和周期任务。

# Timing Wheels Implementation vs Traditional Implementation

## Traditional Implementation

- JDK Timer + DelayQueue（小根堆-优先级队列）
- Java ScheduledThreadPoolExecutor（DelayQueue + thread pool）

## Timing Wheels Implementation

hash table + linked list + ticker（心跳信号、节拍器、555 信号发生器）

- Kafka

# 单层级时间轮

这种实现可能会存在调度时间不精确的现象。

# 多层级时间轮

通过对时间进行层级划分的方式来提高精确度。

# 实现细节讨论

## slot/bucket 如何划分时间范围？

按照 kafka 的实现，timing wheel 的时间包括 startMs, tickMs, wheelSize, interval 和 currentTimeMs，到 slot 上就是 expirationMs。

但是 kafka 的时间是靠 DelayQueue 在 timeoutMs 时间内把所有过期的 buckets 取出来，再取出内部的 tasks 进行调度和 reinsert 操作。

在 golang 中，有 ticker 作为高精度的时间信号发生器来驱动时间轮的转动，所以不需要 DelayQueue 来进行过期的 bucket 的取出和 reinsert 操作。

ticker 的精度就是至关重要了，如果精度不够，那么就会导致调度时间不准确，比如跨度大，会导致在 slot 里面比较远离 expirationMs 的任务被延后调度。

例子：
假设 tickMs = 10ms，它相当于 kafka 中的 interval，那么 slot 的划分就是 [0, 10), [10, 20), [20, 30) 这样的。

如果 task 被设置的任务间隔是 2ms，那么它第一次应该是在 [0, 10) 的 slot 里面，但是触发调度的时间点却是在 10ms，这样就导致了调度时间不准确。

虽然可以通过限制提交进入时间轮的最小时间间隔来避免这个问题，但这样会导致时间轮可面向的场景变少，比如无法支持 1ms 的时间间隔。

## slot/bucket 为什么要使用多层级？
slots 的数量是有限制的，也就是时间间隔的数量是一定的，时间在向前流逝，slot 也会不断地指向下一个。

虽然使用（逻辑上/物理上）的环境数据结构，可以在 slot 移动过程中把已经执行的任务的 slot 变更为下一个周期可以使用的 slot（复用）。

但是如果要添加的任务时间间隔大于了当前时间轮的 slot 数量，那么就会导致没有 slot 可以承接这个任务，这样就会导致任务丢失。

![移动](https://pic1.zhimg.com/v2-480d1f3a2ddea9a5ccbc87235accc5d0_b.gif)

## slot/bucket 如何实现多层级？

任务的提交涉及到时间跨度大的问题，比如秒级，分钟级，小时级，天级，月级，年级，这样的时间跨度，如果都放在一个时间轮里面，那么就会导致时间轮的大小很大，
比如 365 * 24 * 60 * 60 * 1000 = 31536000000，这样的时间轮就会有 31536000000 个 slot，这样的时间轮是不现实的。

所以需要对时间轮进行分层，比如秒级的时间轮，分钟级的时间轮，小时级的时间轮，天级的时间轮，月级的时间轮，年级的时间轮。

但是这样又会导致时间轮的每一次层的 slot 数量不一致。

![动态时间轮升级和降级](https://pic4.zhimg.com/v2-32981874757256e8b1bff6841f60a2cf_b.gif)

## task 分类
- 一次性延迟执行任务，这种必然有一个距离当前时间的延迟间隔，假设为 delayMs，那么它的过期时间就是 currentTimeMs + delayMs
- 小周期性执行任务，这种任务有一个固定的周期，假设为 periodMs，那么它的过期时间就是 currentTimeMs + periodMs
- 大跨度的周期性执行任务，这种需要计算它的下一次过期时间 expiredMs，因为大跨度的执行，需要考虑年，月，日的变化（极端的时候
还需要考虑到秒级别时间的补偿之类）


## Others
- tickMs `u ms`
- wheelSize `n`
- startMs
- interval = tickMs * wheelSize
- slot (circle array，使用取余模拟循环效果)
- task list (linked list，每个 slot 下挂载着对应的任务)

令 u 为时间单元，一个大小为 n 的时间轮有 n 个桶，能够持有 n * u 个时间间隔的定时任务。
那么每个 slot 都持有进入相应时间范围的定时任务。

- 第一个 slot 持有 [0, u) 范围的任务
- 第二个 slot 持有 [u, 2u) 范围的任务
- 第 n 个 slot 持有 [u * (n - 1), u * n) 范围的任务
  随着时间单元 u 的数量持续增加，也就是 slot 不断地被移动，slot 移动之后其中所有的定时任务都会过期。
  由于任务已经过期，此时就不会有任务被添加大当前的 slot 中，调度器应该立刻运行过期的任务。
  因为空 slot 在下一轮是可用的，所以如果当前的 slot 对应时间 t，那么它会在 tick 后变成 [t + u * n, t + (n + 1) * u) 的 slot。
  简单地说，当前 slot 中的任务会移动到下一个间隔了 t ms 的 slot 中。

slot 切分的间隔需要小一点，比如 2ms，5ms 这样的，防止间隔跨度大，导致调度时间偏移（不准确）
