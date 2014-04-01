package comparable

type Interface interface {
	Compare(i, j int) int
	Len() int
}

func Max(data Interface) (value int) {
	for i := 1; i < data.Len(); i++ {
		switch data.Compare(value, i) {
		case -1:
			value = i
		case 0:
			break
		case 1:
			break
		}
	}
	return value
}
