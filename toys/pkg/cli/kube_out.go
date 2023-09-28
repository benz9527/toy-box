package cli

import (
	"bytes"
	"encoding/json"

	"github.com/jedib0t/go-pretty/v6/table"
	YamlV3 "gopkg.in/yaml.v3"
	MetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	KGeneric "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Frequently encode and decode is resource waste.

func outputFormatFactory(typ kubeResultOutFormatType) visitorFn {
	var v visitorFn
	switch typ {
	case kubeOutYaml:
		v = outAsYaml
	case kubeOutJson:
		v = outAsJson
	case kubeOutTable:
		fallthrough
	default:
		v = outAsTable
	}
	return noResultInterceptor(v)
}

func outAsTable(info *commandContext, err error) error {
	var tabObj = MetaV1.Table{}
	if err = json.Unmarshal(info.genericObj, &tabObj); err != nil {
		return err
	}

	w := table.NewWriter()
	th := make(table.Row, 0, 8)
	for _, h := range tabObj.ColumnDefinitions {
		th = append(th, h.Name)
	}
	w.AppendHeader(th)
	trList := make([]table.Row, 0, 8)
	for _, r := range tabObj.Rows {
		trList = append(trList, r.Cells)
	}
	w.AppendRows(trList)
	info.result = w.Render()
	return nil
}

func outAsJson(info *commandContext, err error) error {
	info.result = bytes.NewBuffer(info.genericObj).String()
	return nil
}

func outAsYaml(info *commandContext, err error) error {
	var y = KGeneric.Unstructured{}
	if err = json.Unmarshal(info.genericObj, &y); err != nil {
		return err
	}
	var data []byte
	if data, err = YamlV3.Marshal(&y.Object); err != nil {
		return err
	}
	info.result = bytes.NewBuffer(data).String()
	return nil
}

func noResultInterceptor(fn visitorFn) visitorFn {
	return func(info *commandContext, err error) error {
		if len(info.genericObj) == 0 {
			return errCliExecEmptyResult
		}
		return fn(info, err)
	}
}
