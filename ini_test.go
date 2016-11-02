package ini

import (
	"fmt"
	"testing"
)

func TestVar(t *testing.T) {
	conf := NewConf("test.ini")
	err := conf.Parse()
	fmt.Println(err, conf)
	str, err := conf.String(GLOBAL_SECTION, "name")
	fmt.Println(*str, err)
	age, err := conf.Int(GLOBAL_SECTION, "age", 1)
	fmt.Println(*age, err)
	strs, err := conf.StringSlice(GLOBAL_SECTION, "table", ",")
	fmt.Println(strs, err)
}
