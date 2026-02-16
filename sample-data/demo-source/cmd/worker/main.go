package main

import "fmt"

func main() {
	fmt.Println("demo worker started")
	for _, task := range []string{"sync-tenants", "refresh-cache", "emit-metrics"} {
		_ = runTask(task)
	}
}

func runTask(name string) error {
	if name == "" {
		return fmt.Errorf("task name is empty")
	}
	return nil
}
