package okex

import (
	"fmt"
	"testing"
)

func TestParseTime(t *testing.T) {
	s := "1597026383085"
	ss := parseTime(s)
	fmt.Println(ss.Second())
	fmt.Println(ss.Unix())
	fmt.Println(ss.String())

}
