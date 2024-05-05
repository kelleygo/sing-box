package vmess

import (
	"bytes"
	"os"
	"strings"
	"text/template"
)

func TemplateText(tpl string, tplArgv map[string]interface{}) (string, error) {
	tpl1, err := template.New("x").Parse(tpl)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	_ = tpl1.Execute(&b, tplArgv)
	return b.String(), nil
}

func TemplateContent(tpl string, tplArgv map[string]interface{}) ([]byte, error) {
	tpl1, err := template.New("x").Parse(tpl)
	if err != nil {
		return []byte{}, err
	}

	var b bytes.Buffer
	_ = tpl1.Execute(&b, tplArgv)
	return b.Bytes(), nil
}

func CurrentPath() string {
	return getCurrentPath()
}

func getCurrentPath() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
