package transactions

import "fmt"

func itoa(v int) string { return fmt.Sprintf("%d", v) }
func itoa64(v int64) string { return fmt.Sprintf("%d", v) }
func toString(v any) string { return fmt.Sprintf("%v", v) }
