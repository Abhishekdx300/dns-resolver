package dns

import "fmt"

type Stack []string

func (s *Stack) Push(item string) {
	*s = append(*s, item)
}
func (s *Stack) Pop() (string, error) {
	if len(*s) == 0 {
		return "", fmt.Errorf("empty stack")
	}
	ind := len(*s) - 1
	element := (*s)[ind]
	*s = (*s)[:ind]
	return element, nil
}
