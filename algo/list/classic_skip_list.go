package list

import (
	"math/rand"
	"sync/atomic"
	"time"
)

var (
	_ SkipListNodeElement[struct{}] = (*classicSkipListNodeElement[struct{}])(nil) // Type check assertion
	_ SkipListLevel[struct{}]       = (*classicSkipListLevel[struct{}])(nil)       // Type check assertion
	_ SkipList[struct{}]            = (*classicSkipList[struct{}])(nil)            // Type check assertion
)

type classicSkipListNodeElement[E comparable] struct {
	object E
	// 指向垂直方向上的下一个结点，通常 level 为 0 的结点，没有 vBackward；但是作为 levels 的node，一定有 vBackward
	// 但是也只有最底层的才需要设置 vBackward，因为 levels 是和 vBackward 一起的
	vBackward SkipListNodeElement[E]
	// 当前结点作为非索引部分时（单纯存放数据），levels 为空
	// 当前结点包含索引数据时，levels[0] 就已经是索引了，和哨兵的 levels[0] 是不一样的
	levels []SkipListLevel[E]
}

func newClassicSkipListNodeElement[E comparable](level int, obj E) SkipListNodeElement[E] {
	e := &classicSkipListNodeElement[E]{
		object: obj,
		levels: make([]SkipListLevel[E], level), // 一开始就把每一层的间距都设置为 0，也就是一开始就做好了数据分配
	}
	for i := 0; i < level; i++ {
		e.levels[i] = newClassicSkipListLevel[E](0, nil)
	}
	return e
}

func (e *classicSkipListNodeElement[E]) GetObject() E {
	return e.object
}

func (e *classicSkipListNodeElement[E]) GetVerticalBackward() SkipListNodeElement[E] {
	return e.vBackward
}

func (e *classicSkipListNodeElement[E]) SetVerticalBackward(backward SkipListNodeElement[E]) {
	e.vBackward = backward
}

func (e *classicSkipListNodeElement[E]) GetLevels() []SkipListLevel[E] {
	return e.levels
}

func (e *classicSkipListNodeElement[E]) Free() {
	e.object = *new(E)
	e.vBackward = nil
	e.levels = nil
}

type classicSkipListLevel[E comparable] struct {
	span     int64                  // 间距数量
	hForward SkipListNodeElement[E] // 指向水平方向的下一个结点
}

func newClassicSkipListLevel[E comparable](span int64, forward SkipListNodeElement[E]) SkipListLevel[E] {
	return &classicSkipListLevel[E]{
		span:     span,
		hForward: forward,
	}
}

func (lvl *classicSkipListLevel[E]) GetSpan() int64 {
	return atomic.LoadInt64(&lvl.span)
}

func (lvl *classicSkipListLevel[E]) SetSpan(span int64) {
	atomic.StoreInt64(&lvl.span, span)
}

func (lvl *classicSkipListLevel[E]) GetHorizontalForward() SkipListNodeElement[E] {
	return lvl.hForward
}

func (lvl *classicSkipListLevel[E]) SetHorizontalForward(forward SkipListNodeElement[E]) {
	lvl.hForward = forward
}

type classicSkipList[E comparable] struct {
	// 当前 skip list 实际使用了的最大层数，最大不会超过 ClassicSkipListMaxLevel
	// 这个类本身是不包含节点的，这就是个哨兵，用来指向分散在堆中的结点
	// 所以 head.levels[0].hForward 指向的是第一个结点，而这一层级的所有结点
	// 就是完整的单向链表，从 1 开始的层级才是索引
	level           int
	len             int64
	head            SkipListNodeElement[E] // 哨兵结点
	tail            SkipListNodeElement[E] // 哨兵结点
	localCompareTo  compareTo[E]
	randomGenerator *rand.Rand
}

func NewClassicSkipList[E comparable](compareTo compareTo[E]) SkipList[E] {
	sl := &classicSkipList[E]{
		level:           1,
		len:             0,
		localCompareTo:  compareTo,
		randomGenerator: rand.New(rand.NewSource(time.Now().Unix())),
	}
	sl.head = newClassicSkipListNodeElement[E](ClassicSkipListMaxLevel, *new(E))
	// 哨兵头结点的层数一定是 ClassicSkipListMaxLevel
	// 并且做好了初始化，每一层的间距都是 0
	// 防止后续插入的时候，出现空指针异常
	for i := 0; i < ClassicSkipListMaxLevel; i++ {
		sl.head.GetLevels()[i].SetSpan(0)
		sl.head.GetLevels()[i].SetHorizontalForward(nil)
	}
	sl.head.SetVerticalBackward(nil)
	sl.tail = nil
	return sl
}

func (sl *classicSkipList[E]) randomLevel() int {
	level := 1
	for float64(sl.randomGenerator.Int63()&0xFFFF) < ClassicSkipListProbability*0xFFFF {
		level += 1
	}
	if level < ClassicSkipListMaxLevel {
		return level
	}
	return ClassicSkipListMaxLevel
}

func (sl *classicSkipList[E]) GetLevel() int {
	return sl.level
}

func (sl *classicSkipList[E]) Len() int64 {
	return atomic.LoadInt64(&sl.len)
}

func (sl *classicSkipList[E]) Insert(obj E) SkipListNodeElement[E] {
	var (
		update     [ClassicSkipListMaxLevel]SkipListNodeElement[E]
		x          SkipListNodeElement[E]
		levelSpans [ClassicSkipListMaxLevel]int64 // levelSpans[i] 表示哨兵 levels 第 i 层的间距，第 0 层是数据层
		levelIdx   int
		level      int
	)
	// 临时对象，获取哨兵结点，从头开始遍历
	x = sl.head
	// 从最高层开始查找当前新元素的插入位置
	for levelIdx = sl.level - 1; levelIdx >= 0; levelIdx-- {
		if levelIdx == sl.level-1 {
			// 第一次遍历，必然是会进入到这里，当前最高层的间距为 0
			// 因为这里不计算第 0 层的间距，需要高度退一层
			levelSpans[levelIdx] = 0
		} else {
			// 从第二次遍历开始，当前层的间距为上一层索引结点的间距
			levelSpans[levelIdx] = levelSpans[levelIdx+1]
		}

		// 1. 第一次遍历且跳表元素为空，如果当前索引没有下一个索引结点，排名（间距）不变
		// 2. 第N次遍历，判断当前索引节点是否有下一个索引结点，如果有，则比较当前新元素和下一个索引节点的指向元素大小
		//    如果当前新元素大于下一个索引节点的指向元素，则当前新元素的排名为当前索引节点的排名
		//    所谓排名，就是当前索引节点到下一个索引节点的距离
		for x.GetLevels()[levelIdx].GetHorizontalForward() != nil &&
			sl.localCompareTo(x.GetLevels()[levelIdx].GetHorizontalForward().GetObject(), obj) < 0 {
			// 更新当前索引节点到下一个索引节点的距离，即排名
			levelSpans[levelIdx] += x.GetLevels()[levelIdx].GetSpan()
			// 当前新元素的值比较大，更新临时结点为当前比较的结点，继续向后遍历，直到找到下一个比
			// 当前新元素大的索引节点，作为右边界
			x = x.GetLevels()[levelIdx].GetHorizontalForward()
		}
		// 找到了右边界，暂存右边界
		// 没找到右边界，也需要暂存当前临时结点，因为当前临时结点是当前新元素的左边界
		// 继续往下一层遍历
		update[levelIdx] = x // 这里就有可能让 update[0] 指向哨兵头结点，因为第一次跳表为空或者是没有下一个元素（末尾了）
	}
	// update 的最后一层必然是数据部分，而不是索引部分

	// 到这里就相当于找到了当前新元素的左边界和右边界
	// 如果当前新元素的值和右边界的值相等，则不需要插入
	if x.GetLevels()[0].GetHorizontalForward() != nil &&
		sl.localCompareTo(x.GetLevels()[0].GetHorizontalForward().GetObject(), obj) == 0 {
		return nil
	}

	// 元素不存在，需要插入
	// 要生成随机层数
	level = sl.randomLevel()
	if level > sl.level {
		// 如果随机层数大于当前 skip list 的最大层数，则需要更新当前 skip list 的最大层数
		for lvl := sl.level; lvl < level; lvl++ {
			// 从当前 skip list 的最大层数开始，更新每一层的间距
			levelSpans[lvl] = 0   // 当前层的间距，相当于之后的运算中不需要减去这一层的数量
			update[lvl] = sl.head // 新增的层，指向头哨兵结点
			// 只有哨兵头结点的层数是 ClassicSkipListMaxLevel，才能进行遍历
			update[lvl].GetLevels()[lvl].SetSpan(sl.len)
		}
		// 更新当前 skip list 实际使用了的最大层数
		sl.level = level
	}

	// 生成新的结点，准备插入。一个需要插入的元素，一定是会从最底层开始
	x = newClassicSkipListNodeElement[E](level, obj)
	// 从最低层开始，更新每一层的间距
	for levelIdx = 0; levelIdx < level; levelIdx++ {
		// 当前结点的当前层指向
		// update 保留了之前每一层遍历右边界结果，相当于是路径记录
		// 因为新增了索引，需要调整索引之间的指向，类似于双向链表的指针调整
		// 直接是指针复制，需要调整指针指向
		// 如果是第一次插入且跳表为空，下一个指向肯定是 nil
		x.GetLevels()[levelIdx].SetHorizontalForward(update[levelIdx].GetLevels()[levelIdx].GetHorizontalForward())
		update[levelIdx].GetLevels()[levelIdx].SetHorizontalForward(x)

		// 插入新的索引之后，需要重新调整各个层的间距数量（旧索引1 ----> 新索引1 ----> 旧索引2）
		// levelSpans[0] 必须大于 levelSpans[levelIdx]，才能算出差值
		x.GetLevels()[levelIdx].SetSpan(update[levelIdx].GetLevels()[levelIdx].GetSpan() - (levelSpans[0] - levelSpans[levelIdx]))
		update[levelIdx].GetLevels()[levelIdx].SetSpan(levelSpans[0] - levelSpans[levelIdx] + 1)
	}

	// 逐层+1更新更高级的索引层的间距，因为 update 在这就是一个路径，而且路径会指向不同的索引和不同的层
	// 前面更新了是底部的索引层的信息，如果当前插入的结点的层数小于经过路径的层数
	// 就会存在没有被更新到的索引层
	for levelIdx = level; levelIdx < sl.level; levelIdx++ {
		update[levelIdx].GetLevels()[levelIdx].SetSpan(update[levelIdx].GetLevels()[levelIdx].GetSpan() + 1)
	}

	if update[0] == sl.head {
		// 如果最底层的索引结点指向了头结点
		// 1. 空跳表
		// 2. 最后一个元素
		x.SetVerticalBackward(nil) // 当前结点的下一个结点设置为空
	} else {
		x.SetVerticalBackward(update[0]) // 设置为下一个元素
	}

	if x.GetLevels()[0].GetHorizontalForward() != nil {
		// 通常来说，插入了新元素，而且有索引，一定是走这里
		// 因为前面是指针的复制，原本 vBackward 的指向还是旧索引1
		// 只有最底层才需要调整 vBackward，因为 levels 是和 vBackward 一起的
		x.GetLevels()[0].GetHorizontalForward().SetVerticalBackward(x)
	} else {
		// 如果是最后一个元素，需要更改哨兵尾结点的指向
		sl.tail = x
	}
	sl.len++
	return x
}

func (sl *classicSkipList[E]) Remove(obj E) SkipListNodeElement[E] {
	var (
		update [ClassicSkipListMaxLevel]SkipListNodeElement[E]
		x      SkipListNodeElement[E]
		idx    int
	)
	// 临时对象，获取哨兵结点，从头开始遍历
	x = sl.head
	for idx = sl.level - 1; idx >= 0; idx-- {
		// 从最高层开始查找当前待删除元素的位置
		for x.GetLevels()[idx].GetHorizontalForward() != nil &&
			sl.localCompareTo(x.GetLevels()[idx].GetHorizontalForward().GetObject(), obj) < 0 {
			// 如果当前索引节点的下一个索引节点的值小于待删除元素的值
			// 则当前索引节点的下一个索引节点就是待删除元素的右边界
			x = x.GetLevels()[idx].GetHorizontalForward()
		}
		// 不管找没找到右边界，都需要记录当前索引节点
		update[idx] = x
		// 转入下一次的遍历
	}

	// 到这里就相当于找到了当前待删除元素的左边界和右边界
	x = x.GetLevels()[0].GetHorizontalForward()
	if x != nil && sl.localCompareTo(x.GetObject(), obj) == 0 {
		// 找到了待删除元素
		sl.deleteNode(x, update)
		return x
	}
	// 没有找到待删除元素
	return nil
}

func (sl *classicSkipList[E]) deleteNode(x SkipListNodeElement[E], update [32]SkipListNodeElement[E]) {
	var idx int
	// 从底层开始，逐级向上调整索引指向和间距
	for idx = 0; idx < sl.level; idx++ {
		if update[idx].GetLevels()[idx].GetHorizontalForward() == x {
			// 如果当前索引节点的下一个索引节点就是待删除元素
			// 调整当前索引节点的下一个索引节点为待删除元素的下一个索引节点
			// 调整当前索引节点的间距为待删除元素的间距
			update[idx].GetLevels()[idx].SetSpan(update[idx].GetLevels()[idx].GetSpan() + x.GetLevels()[idx].GetSpan() - 1)
			update[idx].GetLevels()[idx].SetHorizontalForward(x.GetLevels()[idx].GetHorizontalForward())
		} else {
			// 如果当前索引节点的下一个索引节点不是待删除元素
			// 调整当前索引节点的间距为待删除元素的间距
			update[idx].GetLevels()[idx].SetSpan(update[idx].GetLevels()[idx].GetSpan() - 1)
		}
	}
	if x.GetLevels()[0].GetHorizontalForward() != nil {
		// 如果待删除元素的下一个索引节点不为空，也就是不是最后一个元素
		// 调整待删除元素的下一个索引节点的 vBackward 为待删除元素的 vBackward
		x.GetLevels()[0].GetHorizontalForward().SetVerticalBackward(x.GetVerticalBackward())
	} else {
		// 如果待删除元素的下一个索引节点为空，也就是是最后一个元素
		// 调整哨兵尾结点的指向为待删除元素的 vBackward
		sl.tail = x.GetVerticalBackward()
	}
	for sl.level > 1 && sl.head.GetLevels()[sl.level-1].GetHorizontalForward() == nil {
		sl.level--
	}
	sl.len--
}

func (sl *classicSkipList[E]) Find(obj E) SkipListNodeElement[E] {
	var (
		x   SkipListNodeElement[E]
		idx int
	)
	// 临时对象，获取哨兵结点，从头开始遍历
	x = sl.head
	// 从最高层开始查找当前匹配的元素的位置
	for idx = sl.level - 1; idx >= 0; idx-- {
		for x.GetLevels()[idx].GetHorizontalForward() != nil &&
			sl.localCompareTo(x.GetLevels()[idx].GetHorizontalForward().GetObject(), obj) < 0 {
			// 还是和插入一样，找到右边界
			x = x.GetLevels()[idx].GetHorizontalForward()
		}
		// 转入下一次的遍历
	}
	// 到这里就相当于找到了当前待删除元素的左边界和右边界
	x = x.GetLevels()[0].GetHorizontalForward()
	// 如果找到了右边界，且右边界的值和待查找元素的值相等，则找到了
	if x != nil && sl.localCompareTo(x.GetObject(), obj) == 0 {
		return x
	}
	// 没有找到
	return nil
}

func (sl *classicSkipList[E]) PopHead() (obj E) {
	x := sl.head
	x = x.GetLevels()[0].GetHorizontalForward()
	if x == nil {
		return obj
	}
	obj = x.GetObject()
	sl.Remove(obj)
	return
}

func (sl *classicSkipList[E]) PopTail() (obj E) {
	x := sl.tail
	if x == nil {
		return *new(E)
	}
	obj = x.GetObject()
	sl.Remove(obj)
	return
}

func (sl *classicSkipList[E]) Free() {
	var (
		x, next SkipListNodeElement[E]
		idx     int
	)
	x = sl.head.GetLevels()[0].GetHorizontalForward()
	for x != nil {
		next = x.GetLevels()[0].GetHorizontalForward()
		x.Free()
		x = nil
		x = next
	}
	for idx = 0; idx < ClassicSkipListMaxLevel; idx++ {
		sl.head.GetLevels()[idx].SetHorizontalForward(nil)
		sl.head.GetLevels()[idx].SetSpan(0)
	}
	sl.tail = nil
	sl.level = 0
	sl.len = 0
}

func (sl *classicSkipList[E]) ForEach(fn func(idx int64, v E)) {
	var (
		x   SkipListNodeElement[E]
		idx int64
	)
	x = sl.head.GetLevels()[0].GetHorizontalForward()
	for x != nil {
		next := x.GetLevels()[0].GetHorizontalForward()
		fn(idx, x.GetObject())
		idx++
		x = next
	}
}
