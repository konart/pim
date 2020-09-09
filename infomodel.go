package pim

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	InfomodelPath = "StructureGroup/StructureGroupAttribute"
)

var (
	//fields for requests
	InfomodelFields = []string{"StructureGroupAttributeLang.Name(Russian)", "StructureGroupAttribute.Datatype",
		"StructureGroupAttributeLang.DomainValue(Russian)", "StructureGroupAttribute.IsMandatory",
		"StructureGroupAttribute.MultiValue"}

	//errors
	TypeCastErr = fmt.Errorf("cant cast value to correct type")
)

type StructureGroup struct {
	Identifier  string
	StructureID int
	Features    map[string]Feature
}

type Feature struct {
	Name         string
	DataType     string
	PresetValues []string
	Mandatory    bool
	Multivalued  bool
}

type StructureGroupProvider struct {
	c *Client
}

func (i *StructureGroupProvider) GetInfomodelByIdentifier(identifier string, structureID int) (*StructureGroup, error) {
	url := i.c.baseUrl() + InfomodelPath + "/byItems?" +
		"items=" + "'" + identifier + "'@" + strconv.Itoa(structureID) +
		"&fields=" + strings.Join(InfomodelFields, ",") +
		"&pageSize=-1"
	res, err := i.c.get(url)
	if err != nil {
		return nil, err
	}
	fs := make(map[string]Feature)
	for _, row := range res.Rows {
		if len(row.Values) != len(InfomodelFields) {
			return nil, fmt.Errorf("cant parse infomodel, wrong num of values in a row")
		}
		name, ok := row.Values[0].(string)
		if !ok {
			return nil, TypeCastErr
		}
		dataType, ok := row.Values[1].(string)
		if !ok {
			return nil, TypeCastErr
		}
		pi, ok := row.Values[2].([]interface{})
		if !ok {
			return nil, TypeCastErr
		}
		presets := make([]string, 0)
		for _, vi := range pi {
			val, ok := vi.(string)
			if !ok {
				return nil, TypeCastErr
			}
			if val == "" {
				continue
			}
			presets = append(presets, val)
		}
		manda, ok := row.Values[3].(bool)
		if !ok {
			return nil, TypeCastErr
		}
		multi, ok := row.Values[4].(bool)
		if !ok {
			return nil, TypeCastErr
		}
		fs[name] = Feature{
			Name:         name,
			DataType:     dataType,
			PresetValues: presets,
			Mandatory:    manda,
			Multivalued:  multi,
		}
	}
	return &StructureGroup{
		Identifier:  identifier,
		StructureID: structureID,
		Features:    fs,
	}, nil
}