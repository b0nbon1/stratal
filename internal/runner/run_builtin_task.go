package runner

import "fmt"

func RunBuiltinTask(name string, params map[string]string) error {
	fmt.Println(params)
	return nil
}
