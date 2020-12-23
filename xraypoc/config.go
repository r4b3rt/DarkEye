package xraypoc

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

//Poc add comment
type Poc struct {
	Name   string            `yaml:"name"`
	Set    map[string]string `yaml:"set"`
	Rules  []Rules           `yaml:"rules"`
	Detail Detail            `yaml:"detail"`
}

//Rules add comment
type Rules struct {
	Method          string            `yaml:"method"`
	Path            string            `yaml:"path"`
	Headers         map[string]string `yaml:"headers"`
	Body            string            `yaml:"body"`
	Search          string            `yaml:"search"`
	FollowRedirects bool              `yaml:"follow_redirects"`
	Expression      string            `yaml:"expression"`
}

//Detail add comment
type Detail struct {
	Author      string   `yaml:"author"`
	Links       []string `yaml:"links"`
	Description string   `yaml:"description"`
	Version     string   `yaml:"version"`
}

//XArYPoc add comment
type XArYPoc struct {
	ReverseUrl      string
	ReverseUrlCheck string
}

//LoadPoc add comment
func LoadPoc(fileName string) (*Poc, error) {
	p := &Poc{}
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, p)
	if err != nil {
		return nil, err
	}
	return p, err
}

//LoadPocByData add comment
func LoadPocByData(data []byte) (*Poc, error) {
	if data == nil {
		return nil, fmt.Errorf("NoData")
	}
	p := &Poc{}
	err := yaml.Unmarshal(data, p)
	if err != nil {
		return nil, err
	}
	return p, err
}
