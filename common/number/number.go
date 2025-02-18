package number

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~float32 | ~float64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func SelectBigger[T Number](num1 T, num2 T) T {
	if num1 >= num2 {
		return num1
	}
	return num2
}

func SelectSmaller[T Number](num1 T, num2 T) T {
	if num1 <= num2 {
		return num1
	}
	return num2
}

func SelectNotZero[T Number](num1 T, num2 T) T {
	if num1 != 0 {
		return num1
	}
	return num1
}
