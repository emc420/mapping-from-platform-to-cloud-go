package ontology

import "errors"

type JmesPathPointType string

const (
	STRING__Type JmesPathPointType = "string"
	INT64__Type  JmesPathPointType = "int64"
	DOUBLE_Type  JmesPathPointType = "double"
	OBIX_Type    JmesPathPointType = "obix"
	XML_Type     JmesPathPointType = "xml"
	BOOLEAN_Type JmesPathPointType = "boolean"
	OBJECT_Type  JmesPathPointType = "object"
)

func (lt *JmesPathPointType) FromValue(text string) (JmesPathPointType, error) {
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

func (lt *JmesPathPointType) GetValue(pointType JmesPathPointType) (string, error) {
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
