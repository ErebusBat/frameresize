package frameresize

type StringSlice []string

func (s StringSlice) Contains(test string) bool {
	for _, a := range s {
		if a == test {
			return true
		}
	}
	return false
}
