package api

import "testing"

type breadIF interface {
	Name() string
}

type XiaoLongBao struct{}

func (XiaoLongBao) Name() string {
	return "小笼包"
}

func BenchmarkBaseInterface(b *testing.B) {
	// IF convert is slower than directly and IF
	b.Run("XiaoLongBao IF convert", func(b *testing.B) {
		var bread breadIF = XiaoLongBao{}
		for i := 0; i < b.N; i++ {
			bread.(XiaoLongBao).Name()
		}
	})
	// No great difference between IF and directly
	b.Run("XiaoLongBao IF", func(b *testing.B) {
		var bread breadIF = XiaoLongBao{}
		for i := 0; i < b.N; i++ {
			bread.Name()
		}
	})
	b.Run("XiaoLongBao directly", func(b *testing.B) {
		var bread = XiaoLongBao{}
		for i := 0; i < b.N; i++ {
			bread.Name()
		}
	})
}
