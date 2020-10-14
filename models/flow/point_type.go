package flow

import "errors"

type PointType string

const (
	STRING__Type PointType = "string"
	INT64__Type  PointType = "int64"
	DOUBLE_Type  PointType = "double"
	OBIX_Type    PointType = "obix"
	XML_Type     PointType = "xml"
	BOOLEAN_Type PointType = "boolean"
	OBJECT_Type  PointType = "object"
)

func (lt *PointType) FromValue(text string) (PointType, error) {
	switch {
	case text == "string":
		return STRING__Type, nil
	case text == "int64":
		return INT64__Type, nil
	case text == "double":
		return DOUBLE_Type, nil
	case text == "obix":
		return OBIX_Type, nil
	case text == "xml":
		return XML_Type, nil
	case text == "boolean":
		return BOOLEAN_Type, nil
	case text == "object":
		return OBJECT_Type, nil
	default:
		return "", errors.New("no matching values found")
	}
}

func (lt *PointType) GetValue(pointType PointType) (string, error) {
	switch {
	case pointType == STRING__Type:
		return "string", nil
	case pointType == INT64__Type:
		return "int64", nil
	case pointType == DOUBLE_Type:
		return "double", nil
	case pointType == OBIX_Type:
		return "obix", nil
	case pointType == XML_Type:
		return "xml", nil
	case pointType == BOOLEAN_Type:
		return "boolean", nil
	case pointType == OBJECT_Type:
		return "object", nil
	default:
		return "", errors.New("no matching values found")
	}
}
