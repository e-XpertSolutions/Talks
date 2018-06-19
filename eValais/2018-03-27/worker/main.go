package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

// Command line arguments.
var (
	n = flag.Int("worker", 3, "number of concurrent worker")
)

// A Task represents an methematical operation to be exeucted by a worker.
type Task struct {
	A  int    // First operand
	Op string // Operator (+, -, * or /)
	B  int    // Second operan
}

func (t Task) Do() (int, error) {
	var res int
	switch t.Op {
	case "+":
		res = t.A + t.B
	case "-":
		res = t.A - t.B
	case "*":
		res = t.A * t.B
	case "/":
		res = t.A / t.B
	default:
		return 0, fmt.Errorf("invalid operator %q", t.Op)
	}
	return res, nil
}

func main() {
	flag.Parse()

	// Channel where the tasks will be sent to.
	taskc := make(chan Task, *n)

	// Start n workers. Each worker will read from the channel and execute the
	// task. Output is written into stdout.
	for i := 0; i < *n; i++ {
		go func(ii int) {
			for t := range taskc {
				res, err := t.Do()
				if err != nil {
					fmt.Printf("[worker %02d] %v", ii, err)
				}
				fmt.Printf("[worker %02d] %d %s %d = %d\n", ii, t.A, t.Op, t.B, res)
			}
		}(i)
	}

	// Generate random operations and send them to the tasks channel.
	ops := []string{"+", "-", "*", "/"}
	for i := 0; i < 10; i++ {
		a := rand.Intn(100) + 1
		b := rand.Intn(100) + 1
		op := ops[rand.Intn(len(ops))]
		taskc <- Task{A: a, Op: op, B: b}
	}

	close(taskc)

	// In order to keep this example simple, we just wait 1 second to be sure
	// all tasks have been executed.
	time.Sleep(1 * time.Second)
}
