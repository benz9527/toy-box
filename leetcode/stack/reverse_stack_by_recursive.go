package stack

type myStack struct {
	arr []int
}

func (s *myStack) Pop() int {
	val := s.arr[len(s.arr)-1]
	s.arr = s.arr[0 : len(s.arr)-1]
	return val
}

func (s *myStack) Push(item int) {
	s.arr = append(s.arr, item)
}

func (s *myStack) IsEmpty() bool {
	return len(s.arr) == 0
}

func (s *myStack) GetArray() []int {
	return append([]int{}, s.arr...)
}

func GetAndRemoveLast(s *myStack) int {
	result := s.Pop()
	if s.IsEmpty() {
		return result
	}
	last := GetAndRemoveLast(s)
	s.Push(result)
	return last
}

func ReverseStack(s *myStack) {
	if s.IsEmpty() {
		return
	}
	last := GetAndRemoveLast(s)
	ReverseStack(s)
	s.Push(last)
}
