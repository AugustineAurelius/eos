package generics

type NumericIterator[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64] Iterator[T]

func NewNumericIterator[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](slice []T) NumericIterator[T] {
	return NumericIterator[T](NewIterator(slice))
}

func (i NumericIterator[T]) Sum() T {
	var res T
	for elem := range i {
		res += elem
	}
	return res
}

func (i NumericIterator[T]) Min() T {
	var min T
	first := true
	for elem := range i {
		if first || elem < min {
			min = elem
			first = false
		}
	}
	return min
}

func (i NumericIterator[T]) Max() T {
	var max T
	first := true
	for elem := range i {
		if first || elem > max {
			max = elem
			first = false
		}
	}
	return max
}

func (i NumericIterator[T]) Average() T {
	var sum, count T
	for elem := range i {
		sum += elem
		count++
	}
	if count == 0 {
		return 0
	}
	return sum / count
}
