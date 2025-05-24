package bindflags

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

var pFlagNames = []string{"name", "shorthand", "value", "usage"}
var flagNames = []string{"name", "value", "usage"}

func scanFlagTag(s string) (*FlagTag, error) {
	worlds, err := scanWorld(strings.NewReader(s), ';')
	result, err := scanKV(worlds, flagNames)
	if err != nil {
		return nil, err
	}
	result = formatKV(result)
	return &FlagTag{
		Name:  result["name"],
		Value: result["value"],
		Usage: result["usage"],
	}, nil
}

func scanPFlagTag(s string) (*PFlagTag, error) {
	worlds, err := scanWorld(strings.NewReader(s), ';')
	result, err := scanKV(worlds, pFlagNames)
	if err != nil {
		return nil, err
	}
	result = formatKV(result)
	return &PFlagTag{
		Name:      result["name"],
		Shorthand: result["shorthand"],
		Value:     result["value"],
		Usage:     result["usage"],
	}, nil
}

func scanWorld(r io.RuneReader, split rune) ([]string, error) {
	var builder strings.Builder
	var item = make([]string, 0, 4)
	var tag1, tag2 bool
	for {
		char, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				if tag1 || tag2 {
					c := "\""
					if tag2 {
						c = "'"
					}
					return nil, errors.New("Syntax error: Closing character could not be found " + c)
				}
				if builder.Len() > 0 {
					item = append(item, builder.String())
				}
				builder.Reset()
				break
			}
			return nil, err
		}
		switch char {
		case '"':
			if !tag2 {
				tag1 = !tag1
			}
			builder.WriteRune(char)
		case '\'':
			if !tag1 {
				tag2 = !tag2
			}
			builder.WriteRune(char)
		case split:
			if tag1 || tag2 {
				builder.WriteRune(char)
			} else {
				item = append(item, builder.String())
				builder.Reset()
			}
		default:
			builder.WriteRune(char)
		}
	}
	return item, nil
}

func scanKV(worlds []string, flagNames []string) (map[string]string, error) {
	result := make(map[string]string)
	defaultValues := []string{}
	var isScan bool
	for _, word := range worlds {
		n := strings.IndexByte(word, ':')
		if n == -1 {
			defaultValues = append(defaultValues, word)
			continue
		}
		tempName := strings.ToLower(strings.TrimSpace(word[:n]))
		isScan = false
		for _, fn := range flagNames {
			if fn == tempName {
				result[tempName] = word[n+1:]
				isScan = true
				break
			}
		}
		if !isScan {
			return nil, errors.New("Invalid flag name: " + word[:n])
		}
	}
	for _, name := range flagNames {
		if len(defaultValues) == 0 {
			break
		}
		if _, ok := result[name]; ok {
			continue
		}
		result[name] = defaultValues[0]
		defaultValues = defaultValues[1:]
	}
	return result, nil
}

func formatKV(kv map[string]string) map[string]string {
	for k, v := range kv {
		if strings.HasPrefix(v, "\"") && strings.HasSuffix(v, "\"") {
			kv[k] = v[1 : len(v)-1]
		} else if strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'") {
			kv[k] = v[1 : len(v)-1]
		} else {
			kv[k] = strings.TrimSpace(v)
		}
	}
	return kv
}

func convertValue(name, value, typ string, isSlice ...bool) interface{} {
	slice := false
	if len(isSlice) > 0 {
		slice = isSlice[0]
	}
	var v interface{}
	var err error
	switch typ {
	case "string":
		if slice {
			v = []string{}
			err = json.Unmarshal([]byte(value), &v)
		} else {
			v = value
		}
	case "int":
		if slice {
			v = []int{}
			err = json.Unmarshal([]byte(value), &v)
		} else {
			v, err = strconv.Atoi(value)
		}
	case "int8":
		if slice {
			v = []int8{}
			err = json.Unmarshal([]byte(value), &v)
		} else {
			n, e := strconv.ParseInt(value, 10, 8)
			v = int8(n)
			err = e
		}

	case "int16":
		if slice {
			v = []int16{}
			err = json.Unmarshal([]byte(value), &v)
		} else {
			n, e := strconv.ParseInt(value, 10, 16)
			v = int16(n)
			err = e
		}

	case "int32":
		if slice {
			v = []int32{}
			err = json.Unmarshal([]byte(value), &v)
		} else {
			n, e := strconv.ParseInt(value, 10, 32)
			v = int32(n)
			err = e
		}

	case "int64":
		if slice {
			v = []int64{}
			err = json.Unmarshal([]byte(value), &v)
		} else {
			v, err = strconv.ParseInt(value, 10, 64)
		}

	case "uint":
		if slice {
			v = []uint{}
			err = json.Unmarshal([]byte(value), &v)
		} else {
			n, e := strconv.ParseUint(value, 10, 32)
			v = uint(n)
			err = e
		}

	case "uint8":
		if slice {
			v = []uint8{}
			err = json.Unmarshal([]byte(value), &v)
		} else {
			n, e := strconv.ParseUint(value, 10, 8)
			v = uint8(n)
			err = e
		}

	case "uint16":
		if slice {
			v = []uint16{}
			err = json.Unmarshal([]byte(value), &v)
		} else {
			n, e := strconv.ParseUint(value, 10, 16)
			v = uint16(n)
			err = e
		}

	case "uint32":
		if slice {
			v = []uint32{}
			err = json.Unmarshal([]byte(value), &v)
		} else {
			n, e := strconv.ParseUint(value, 10, 32)
			v = uint32(n)
			err = e
		}

	case "uint64":
		if slice {
			v = []uint64{}
			err = json.Unmarshal([]byte(value), &v)
		} else {
			v, err = strconv.ParseUint(value, 10, 64)
		}

	case "float32":
		if slice {
			v = []float32{}
			err = json.Unmarshal([]byte(value), &v)
		} else {
			n, e := strconv.ParseFloat(value, 32)
			v = float32(n)
			err = e
		}

	case "float64":
		if slice {
			v = []float64{}
			err = json.Unmarshal([]byte(value), &v)
		} else {
			v, err = strconv.ParseFloat(value, 64)
		}

	case "bool":
		if slice {
			v = []bool{}
			err = json.Unmarshal([]byte(value), &v)
		} else {
			v, err = strconv.ParseBool(value)
		}

	case "duration", "time.Duration":
		if slice {
			v = []time.Duration{}
			err = json.Unmarshal([]byte(value), &v)
		} else {
			v, err = time.ParseDuration(value)
		}
	default:
		err = errors.New(typ + ": unsupported type")
	}
	if err != nil && value != "" {
		panic(fmt.Sprintf("flag tag %#v value %#v invalid: %v", name, value, err))
	}
	return v
}
