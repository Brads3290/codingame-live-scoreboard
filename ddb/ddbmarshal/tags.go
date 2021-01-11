package ddbmarshal

import "strings"

type ddbTag struct {
	FieldName string
	keyType   string
}

func newDdbTag(tagStr string) ddbTag {
	var dbt ddbTag

	// Split the tag string by commas
	s := strings.Split(tagStr, ",")

	// Trim all the strings
	for i, v := range s {
		s[i] = strings.Trim(v, " \t")
	}

	switch true {
	case len(s) >= 2:
		dbt.keyType = s[1]
		fallthrough
	case len(s) >= 1:
		dbt.FieldName = s[0]
	}

	return dbt
}

func (dbt ddbTag) IsIgnored() bool {
	return dbt.FieldName == "-"
}

func (dbt ddbTag) IsKey() bool {
	return dbt.keyType == "key"
}
