package ini

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	GLOBAL_SECTION = "_GLOBAL_SECTION_"
)

type ConfSet struct {
	filePath string
	parsed   bool
	sections map[string]*Section
}

type Section struct {
	name string
	vals map[string]interface{}
}

func NewConf(filePath string) *ConfSet {
	return &ConfSet{filePath, false, make(map[string]*Section)}
}

func (c ConfSet) Int(sectionName, name string, values ...int) (*int, error) {
	var rst = new(int)
	section, ok := c.sections[sectionName]
	if !ok {
		er := fmt.Errorf("不存在名为【%s】的section", sectionName)
		return nil, er
	}
	vals, ok := section.vals[name]
	str := vals.(string)
	val, err := strconv.Atoi(str)
	if nil != err {
		return nil, err
	}
	if ok {
		*rst = val
	} else {
		if 0 == len(values) {
			er := fmt.Errorf("配置文件中没有【%s】的值，也没有设置默认值")
			return nil, er
		}
		*rst = values[0]
	}
	return rst, nil
}

func (c ConfSet) StringSlice(sectionName, name, separator string, values ...[]string) ([]string, error) {
	section, ok := c.sections[sectionName]
	if !ok {
		er := fmt.Errorf("不存在名为【%s】的section", sectionName)
		return nil, er
	}
	vals, ok := section.vals[name]
	var rst []string
	if ok {
		strs := strings.Split(vals.(string), separator)
		for _, v := range strs {
			str := strings.TrimSpace(v)
			rst = append(rst, str)
		}
	} else {
		if 0 == len(values) {
			er := fmt.Errorf("配置文件中没有【%s】的值，也没有设置默认值")
			return nil, er
		}
		rst = values[0]
	}
	return rst, nil
}

func (c ConfSet) String(sectionName, name string, values ...string) (*string, error) {
	var str = new(string)
	section, ok := c.sections[sectionName]
	if !ok {
		er := fmt.Errorf("不存在名为【%s】的section", sectionName)
		return nil, er
	}
	vals, ok := section.vals[name]
	if ok {
		*str = vals.(string)
	} else {
		if 0 == len(values) {
			er := fmt.Errorf("配置文件中没有【%s】的值，也没有设置默认值")
			return nil, er
		}
		*str = values[0]
	}
	return str, nil
}

func (c *ConfSet) parseOne(sectionName, line string) error {
	if nil == c.sections {
		er := fmt.Errorf("ConfSet未初始化，请先调用NewConf函数")
		return er
	}
	section, ok := c.sections[sectionName]
	// handle line
	strs := strings.SplitN(line, "=", 2)
	name, value := strings.TrimSpace(strs[0]), strings.TrimSpace(strs[1])
	if !ok {
		section = new(Section)
	}
	if nil == section.vals {
		section.vals = make(map[string]interface{})
	} else if _, ok := section.vals[name]; ok {
		er := fmt.Errorf("section 为【%s】的name【%s】已存在", sectionName, name)
		return er
	}
	section.vals[name] = value
	section.name = sectionName
	c.sections[sectionName] = section
	return nil
}

func (c *ConfSet) Parse() error {
	c.parsed = true
	currentSection := GLOBAL_SECTION

	file, err := os.Open(c.filePath)
	if nil != err {
		er := fmt.Errorf("打开文件出错，请检查文件目录是否正确")
		fmt.Println(er.Error(), err)
		return er
	}
	reader := bufio.NewReader(file)
	for {
		// ReadLine 是一个低水平的行读取原语，大多数情况下，应该使用
		// ReadBytes('\n') 或 ReadString('\n')，或者使用一个 Scanner
		line, err := reader.ReadBytes('\n')

		if io.EOF == err && 0 == len(line) {
			break
		} else if 0 == len(line) {
			continue
		}

		l := strings.TrimSpace(string(line))
		switch l[0] {
		case '[':
			l := strings.TrimSpace(l)
			if ']' == l[len(l)-1] {
				currentSection = l[1 : len(l)-1]
				continue
			}
		case '#':
			continue
		case '/':
			if '/' == l[1] {
				continue
			}
		}

		// parse item
		err = c.parseOne(currentSection, l)
		if nil != err {
			return err
		}
	}
	return nil
}
