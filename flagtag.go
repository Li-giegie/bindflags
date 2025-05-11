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

type IFlagTag interface {
	GetName() string
	GetShorthand() string
	GetValue() string
	GetUsage() string
}

type FlagTag struct {
	Name      string
	Shorthand string
	Value     string
	Usage     string
}

func (f *FlagTag) GetName() string {
	return f.Name
}
func (f *FlagTag) GetShorthand() string {
	return f.Shorthand
}
func (f *FlagTag) GetValue() string {
	return f.Value
}
func (f *FlagTag) GetUsage() string {
	return f.Usage
}

func (f *FlagTag) convertValue(typ string, isSlice ...bool) interface{} {
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
			err = json.Unmarshal([]byte(f.Value), &v)
		} else {
			v = f.Value
		}
	case "int":
		if slice {
			v = []int{}
			err = json.Unmarshal([]byte(f.Value), &v)
		} else {
			v, err = strconv.Atoi(f.Value)
		}
	case "int8":
		if slice {
			v = []int8{}
			err = json.Unmarshal([]byte(f.Value), &v)
		} else {
			n, e := strconv.ParseInt(f.Value, 10, 8)
			v = int8(n)
			err = e
		}

	case "int16":
		if slice {
			v = []int16{}
			err = json.Unmarshal([]byte(f.Value), &v)
		} else {
			n, e := strconv.ParseInt(f.Value, 10, 16)
			v = int16(n)
			err = e
		}

	case "int32":
		if slice {
			v = []int32{}
			err = json.Unmarshal([]byte(f.Value), &v)
		} else {
			n, e := strconv.ParseInt(f.Value, 10, 32)
			v = int32(n)
			err = e
		}

	case "int64":
		if slice {
			v = []int64{}
			err = json.Unmarshal([]byte(f.Value), &v)
		} else {
			v, err = strconv.ParseInt(f.Value, 10, 64)
		}

	case "uint":
		if slice {
			v = []uint{}
			err = json.Unmarshal([]byte(f.Value), &v)
		} else {
			n, e := strconv.ParseUint(f.Value, 10, 32)
			v = uint(n)
			err = e
		}

	case "uint8":
		if slice {
			v = []uint8{}
			err = json.Unmarshal([]byte(f.Value), &v)
		} else {
			n, e := strconv.ParseUint(f.Value, 10, 8)
			v = uint8(n)
			err = e
		}

	case "uint16":
		if slice {
			v = []uint16{}
			err = json.Unmarshal([]byte(f.Value), &v)
		} else {
			n, e := strconv.ParseUint(f.Value, 10, 16)
			v = uint16(n)
			err = e
		}

	case "uint32":
		if slice {
			v = []uint32{}
			err = json.Unmarshal([]byte(f.Value), &v)
		} else {
			n, e := strconv.ParseUint(f.Value, 10, 32)
			v = uint32(n)
			err = e
		}

	case "uint64":
		if slice {
			v = []uint64{}
			err = json.Unmarshal([]byte(f.Value), &v)
		} else {
			v, err = strconv.ParseUint(f.Value, 10, 64)
		}

	case "float32":
		if slice {
			v = []float32{}
			err = json.Unmarshal([]byte(f.Value), &v)
		} else {
			n, e := strconv.ParseFloat(f.Value, 32)
			v = float32(n)
			err = e
		}

	case "float64":
		if slice {
			v = []float64{}
			err = json.Unmarshal([]byte(f.Value), &v)
		} else {
			v, err = strconv.ParseFloat(f.Value, 64)
		}

	case "bool":
		if slice {
			v = []bool{}
			err = json.Unmarshal([]byte(f.Value), &v)
		} else {
			v, err = strconv.ParseBool(f.Value)
		}

	case "duration", "time.Duration":
		if slice {
			v = []time.Duration{}
			err = json.Unmarshal([]byte(f.Value), &v)
		} else {
			v, err = time.ParseDuration(f.Value)
		}
	default:
		err = errors.New(typ + ": unsupported type")
	}
	if err != nil && f.Value != "" {
		panic(fmt.Sprintf("flag tag %#v value invalid: %v", f.Name, err))
	}
	return v
}

var flagNames = []string{"name", "shorthand", "value", "usage"}

func scanFlagTag(s string) (*FlagTag, error) {
	worlds, err := scanWorld(strings.NewReader(s), ';')
	result, err := scanKV(worlds)
	if err != nil {
		return nil, err
	}
	result = formatKV(result)
	return &FlagTag{
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

func scanKV(worlds []string) (map[string]string, error) {
	result := make(map[string]string)
	defaultValues := []string{}
	for _, word := range worlds {
		n := strings.IndexByte(word, ':')
		if n == -1 {
			defaultValues = append(defaultValues, word)
			continue
		}
		tempName := strings.ToLower(strings.TrimSpace(word[:n]))
		switch tempName {
		case "name", "shorthand", "value", "usage":
			result[tempName] = word[n+1:]
		default:
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
