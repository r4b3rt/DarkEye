package common

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"sort"
	"strconv"
)

// Convert JSON to CSV
func Convert(r io.Reader, w io.Writer) error {
	dec := json.NewDecoder(r)

	var data interface{}

	if err := dec.Decode(&data); err != nil {
		return err
	}

	var rows []map[string]string

	switch value := data.(type) {
	case []interface{}:
		for i := range value {
			rows = append(rows, topLevelObject(value[i]))
		}
	default:
		rows = append(rows, topLevelObject(value))
	}

	columns := make(map[string]string)
	for i := range rows {
		for col := range rows[i] {
			columns[col] = col
		}
	}
	var colRecord []string
	for c := range columns {
		colRecord = append(colRecord, c)
	}
	sort.Strings(colRecord)

	cw := csv.NewWriter(w)

	if err := cw.Write(colRecord); err != nil {
		return err
	}

	for i := range rows {
		record := make([]string, 0, len(columns))
		for _, col := range colRecord {
			record = append(record, rows[i][col])
		}
		if err := cw.Write(record); err != nil {
			return err
		}

	}

	cw.Flush()

	return nil
}

func topLevelObject(object interface{}) map[string]string {
	values := make(map[string]string)

	switch value := object.(type) {
	case string:
		values["text"] = value
	case map[string]interface{}:
		flattenObject("", values, value)
	case float64:
		values["number"] = strconv.FormatFloat(value, 'f', -1, 64)
	case []interface{}:
		addValue("", values, value)
	}

	return values
}

func flattenObject(path string, values map[string]string, obj map[string]interface{}) {
	for k, v := range obj {
		p := k
		if path != "" {
			p = path + "." + k
		}
		addValue(p, values, v)
	}
}

func addValue(path string, values map[string]string, v interface{}) {
	switch value := v.(type) {
	case string:
		values[path] = value
	case map[string]interface{}:
		flattenObject(path, values, value)
	case float64:
		values[path] = strconv.FormatFloat(value, 'f', -1, 64)
	case []interface{}:
		for i := range value {
			p := strconv.Itoa(i)
			if path != "" {
				p = path + "." + p
			}
			addValue(p, values, value[i])
		}
	}
}
