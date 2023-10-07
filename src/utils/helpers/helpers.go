package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"moadian/src/utils/format"
	"moadian/src/utils/verhoeff"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)


func GetNormalizedData(data any) string {
	arrayOfValues := []string{}
	flattenData := FlattenData(data)
	sortedFlattedData := format.GetSortedMapKeys(flattenData)
	for _, value := range sortedFlattedData {
		dataValue := flattenData[value]
		switch dataValue.(type) {
		case nil:
			arrayOfValues = append(arrayOfValues, "#")
		default:
			newValue := fmt.Sprint(dataValue)
			if newValue == "" {
				arrayOfValues = append(arrayOfValues, "#")
			} else {
				arrayOfValues = append(arrayOfValues, strings.ReplaceAll(newValue, "#", "##"))
			}
		}
	}
	return strings.Join(arrayOfValues, "#")
}


func FlattenData(data any) (map[string]any) {
	falattenRes := make(map[string]any)
	falattenRes = Flatten(data, falattenRes, "")
	return falattenRes
}

func Flatten(data any, falattenRes map[string]any, keyName string) (map[string]any) {
	switch data.(type) {
	case map[string]any:
		nestedMap := data.(map[string]any)
		for key, value := range nestedMap {
			falattenRes = Flatten(value, falattenRes, keyName + key + ".")
		}
	case []interface{}:
		valueSliceIntf := data.([]interface{})
		for i, nestedSliceValue := range valueSliceIntf {
			falattenRes = Flatten(nestedSliceValue, falattenRes, keyName + fmt.Sprint(i) + ".")
		}
	case []map[string]any:
		valueSliceMap := data.([]map[string]any)
		for i, nestedSliceValue := range valueSliceMap {
			falattenRes = Flatten(nestedSliceValue, falattenRes, keyName + fmt.Sprint(i) + ".")
		}
	case []string:
		valueSliceString := data.([]string)
		for i, nestedSliceValue := range valueSliceString {
			falattenRes = Flatten(nestedSliceValue, falattenRes, keyName + fmt.Sprint(i) + ".")
		}
	case []int:
		valueSliceInt := data.([]string)
		for i, nestedSliceValue := range valueSliceInt {
			falattenRes = Flatten(nestedSliceValue, falattenRes, keyName + fmt.Sprint(i) + ".")
		}
	case nil:
		falattenRes[keyName[:len(keyName)-1]] = data
	case bool:
		boolValue := data.(bool)
		falattenRes[keyName[:len(keyName)-1]] = boolValue
	default:
		stringValue := fmt.Sprint(data)
		falattenRes[keyName[:len(keyName)-1]] = stringValue
	}
	return falattenRes
}

func RandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
	  return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func ConvertTaxFiscalIDToNumeric(value string) string {
	resValue := ""
	for _, val := range value {
		valString := string(val)
		if (unicode.IsDigit(val)) {
			resValue += valString
		} else {
			codePoint, _ := utf8.DecodeRuneInString(valString)
			resValue += fmt.Sprint(int(codePoint))
		}
	}
	return resValue
}

func GenerateTaxID(fiscalID string, datetime time.Time, factorID uint) string {
	days := int64(math.Floor(datetime.Sub(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).Seconds() / 86400))
	daysPadded := fmt.Sprintf("%06d", days)
	dateHex := strconv.FormatInt(days, 16)
	dateHex = fmt.Sprintf("%05s", strings.ToUpper(dateHex))
	factorIDPad := fmt.Sprintf("%012d", factorID)
	factorIDHex := fmt.Sprintf("%x", factorID)
	factorIDHexPad := fmt.Sprintf("%010s", factorIDHex)
	fiscalIDNumeric := ConvertTaxFiscalIDToNumeric(fiscalID)
	verhoeffNum := verhoeff.GenerateVerhoeff(fiscalIDNumeric + daysPadded + factorIDPad)
	verhoeffNumStr := fmt.Sprint(verhoeffNum)
	return strings.ToUpper(fiscalID) + dateHex + factorIDHexPad + verhoeffNumStr
}