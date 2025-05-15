package bindflags

import (
	"errors"
	"fmt"
	"github.com/spf13/pflag"
	"reflect"
	"strings"
	"time"
)

var TagName = "flag"

type GetFlagTag interface {
	GetFlagTag() IFlagTag
}

// BindFlags binds the struct member field to the cobra flag, and the struct input parameter must be a pointer type; Add a declaration to the tag of the field such as
// Key-value pair:Name string 'flag:"Name:name; shorthand:n; value:ss; usage:name of student"`
// No key: Name string 'flag:"name; n; ss; name of student"`
// Blend mode: "Name string 'flag:"Name:name; n; ss; name of student"`”
// Or define a custom type, and then implement the GetFlagTag interface for the type
// BindFlags 把结构体成员字段绑定到cobra FlagSet中，结构体入参必须是指针类型；在字段的tag加上声明如
// 键值对：Name  string `flag:"Name:name;shorthand:n;value:ss;usage:name of student"`
// 无键值：Name  string `flag:"name;n;ss;name of student"`
// 混合模式： “Name  string `flag:"Name:name;n;ss;name of student"`”
// 再或者 定义一个自定义类型，然后给类型实现 GetFlagTag 接口
func BindFlags(flag *pflag.FlagSet, a any, group ...string) error {
	rv := reflect.ValueOf(a)
	if rv.Kind() != reflect.Ptr {
		return errors.New("a must be a pointer")
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return errors.New("a must be a struct")
	}
	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		if !rt.Field(i).IsExported() {
			continue
		}
		ft := rt.Field(i)
		tag := ft.Tag.Get(TagName)
		if tag == "-" {
			continue
		}
		fv := rv.Field(i)
		if fv.Kind() == reflect.Ptr {
			if fv.IsNil() {
				fv.Set(reflect.New(fv.Type().Elem()))
			}
			fv = fv.Elem()
		}
		var flagTag = new(FlagTag)
		var err error
		if tag == "" {
			if ft.Type.Implements(reflect.TypeOf((*GetFlagTag)(nil)).Elem()) {
				result := rv.Field(i).Interface().(GetFlagTag).GetFlagTag()
				flagTag.Name = result.GetName()
				flagTag.Usage = result.GetUsage()
				flagTag.Value = result.GetValue()
				flagTag.Shorthand = result.GetShorthand()
			} else if ft.Type.Kind() != reflect.Struct {
				continue
			}
		} else {
			flagTag, err = scanFlagTag(tag)
			if err != nil {
				return err
			}
		}
		if flagTag.Name == "" && rv.Kind() != reflect.Struct {
			return fmt.Errorf("flag '%s' is required", flagTag.Name)
		}
		groupName := flagTag.Name
		if flagTag.Name != "" {
			flagTag.Name = strings.Join(append(group, flagTag.Name), ".")
		}
		switch fv.Kind() {
		case reflect.Struct:
			if groupName != "" {
				group = append(group, groupName)
			}
			if err = BindFlags(flag, fv.Addr().Interface(), group...); err != nil {
				return err
			}
		case reflect.Slice:
			switch fv.Type().Elem().Kind() {
			case reflect.String:
				if flagTag.Shorthand != "" && flagTag.Shorthand != "-" {
					flag.StringSliceVarP((*[]string)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.Shorthand, flagTag.convertValue("string", true).([]string), flagTag.Usage)
				} else {
					flag.StringSliceVar((*[]string)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.convertValue("string", true).([]string), flagTag.Usage)
				}
			case reflect.Int:
				if flagTag.Shorthand != "" && flagTag.Shorthand != "-" {
					flag.IntSliceVarP((*[]int)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.Shorthand, flagTag.convertValue("int", true).([]int), flagTag.Usage)
				} else {
					flag.IntSliceVar((*[]int)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.convertValue("int", true).([]int), flagTag.Usage)
				}
			case reflect.Int32:
				if flagTag.Shorthand != "" && flagTag.Shorthand != "-" {
					flag.Int32SliceVarP((*[]int32)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.Shorthand, flagTag.convertValue("int32", true).([]int32), flagTag.Usage)
				} else {
					flag.Int32SliceVar((*[]int32)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.convertValue("int32", true).([]int32), flagTag.Usage)
				}
			case reflect.Int64:
				if flagTag.Shorthand != "" && flagTag.Shorthand != "-" {
					flag.DurationSliceVarP((*[]time.Duration)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.Shorthand, flagTag.convertValue("duration", true).([]time.Duration), flagTag.Usage)
				} else {
					flag.DurationSliceVar((*[]time.Duration)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.convertValue("duration", true).([]time.Duration), flagTag.Usage)
				}
			case reflect.Uint:
				if flagTag.Shorthand != "" && flagTag.Shorthand != "-" {
					flag.UintSliceVarP((*[]uint)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.Shorthand, flagTag.convertValue("uint", true).([]uint), flagTag.Usage)
				} else {
					flag.UintSliceVar((*[]uint)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.convertValue("uint", true).([]uint), flagTag.Usage)
				}
			case reflect.Float32:
				if flagTag.Shorthand != "" && flagTag.Shorthand != "-" {
					flag.Float32SliceVarP((*[]float32)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.Shorthand, flagTag.convertValue("float32", true).([]float32), flagTag.Usage)
				} else {
					flag.Float32SliceVar((*[]float32)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.convertValue("float32", true).([]float32), flagTag.Usage)
				}
			case reflect.Float64:
				if flagTag.Shorthand != "" && flagTag.Shorthand != "-" {
					flag.Float64SliceVarP((*[]float64)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.Shorthand, flagTag.convertValue("float64", true).([]float64), flagTag.Usage)
				} else {
					flag.Float64SliceVar((*[]float64)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.convertValue("float64", true).([]float64), flagTag.Usage)
				}
			case reflect.Bool:
				if flagTag.Shorthand != "" && flagTag.Shorthand != "-" {
					flag.BoolSliceVarP((*[]bool)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.Shorthand, flagTag.convertValue("bool", true).([]bool), flagTag.Usage)
				} else {
					flag.BoolSliceVar((*[]bool)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.convertValue("bool", true).([]bool), flagTag.Usage)
				}
			default:
				panic(fmt.Sprintf("unsupported type: %T", fv.Interface()))
			}
		case reflect.String:
			if flagTag.Shorthand != "" && flagTag.Shorthand != "-" {
				flag.StringVarP((*string)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.Shorthand, flagTag.Value, flagTag.Usage)
			} else {
				flag.StringVar((*string)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.Value, flagTag.Usage)
			}
		case reflect.Int:
			if flagTag.Shorthand != "" && flagTag.Shorthand != "-" {
				flag.IntVarP((*int)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.Shorthand, flagTag.convertValue("int").(int), flagTag.Usage)
			} else {
				flag.IntVar((*int)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.convertValue("int").(int), flagTag.Usage)
			}
		case reflect.Int8:
			if flagTag.Shorthand != "" && flagTag.Shorthand != "-" {
				flag.Int8VarP((*int8)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.Shorthand, flagTag.convertValue("int8").(int8), flagTag.Usage)
			} else {
				flag.Int8Var((*int8)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.convertValue("int8").(int8), flagTag.Usage)
			}
		case reflect.Int16:
			if flagTag.Shorthand != "" && flagTag.Shorthand != "-" {
				flag.Int16VarP((*int16)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.Shorthand, flagTag.convertValue("int16").(int16), flagTag.Usage)
			} else {
				flag.Int16Var((*int16)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.convertValue("int16").(int16), flagTag.Usage)
			}
		case reflect.Int32:
			if flagTag.Shorthand != "" && flagTag.Shorthand != "-" {
				flag.Int32VarP((*int32)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.Shorthand, flagTag.convertValue("int32").(int32), flagTag.Usage)
			} else {
				flag.Int32Var((*int32)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.convertValue("int32").(int32), flagTag.Usage)
			}
		case reflect.Int64:
			if flagTag.Shorthand != "" && flagTag.Shorthand != "-" {
				flag.Int64VarP((*int64)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.Shorthand, flagTag.convertValue("int64").(int64), flagTag.Usage)
			} else {
				flag.Int64Var((*int64)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.convertValue("int64").(int64), flagTag.Usage)
			}
		case reflect.Uint:
			if flagTag.Shorthand != "" && flagTag.Shorthand != "-" {
				flag.UintVarP((*uint)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.Shorthand, flagTag.convertValue("uint").(uint), flagTag.Usage)
			} else {
				flag.UintVar((*uint)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.convertValue("uint").(uint), flagTag.Usage)
			}
		case reflect.Uint8:
			if flagTag.Shorthand != "" && flagTag.Shorthand != "-" {
				flag.Uint8VarP((*uint8)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.Shorthand, flagTag.convertValue("uint8").(uint8), flagTag.Usage)
			} else {
				flag.Uint8Var((*uint8)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.convertValue("uint8").(uint8), flagTag.Usage)
			}
		case reflect.Uint16:
			if flagTag.Shorthand != "" && flagTag.Shorthand != "-" {
				flag.Uint16VarP((*uint16)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.Shorthand, flagTag.convertValue("uint16").(uint16), flagTag.Usage)
			} else {
				flag.Uint16Var((*uint16)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.convertValue("uint16").(uint16), flagTag.Usage)
			}
		case reflect.Uint32:
			if flagTag.Shorthand != "" && flagTag.Shorthand != "-" {
				flag.Uint32VarP((*uint32)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.Shorthand, flagTag.convertValue("uint32").(uint32), flagTag.Usage)
			} else {
				flag.Uint32Var((*uint32)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.convertValue("uint32").(uint32), flagTag.Usage)
			}
		case reflect.Uint64:
			if flagTag.Shorthand != "" && flagTag.Shorthand != "-" {
				flag.Uint64VarP((*uint64)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.Shorthand, flagTag.convertValue("uint64").(uint64), flagTag.Usage)
			} else {
				flag.Uint64Var((*uint64)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.convertValue("uint64").(uint64), flagTag.Usage)
			}
		case reflect.Float32:
			if flagTag.Shorthand != "" && flagTag.Shorthand != "-" {
				flag.Float32VarP((*float32)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.Shorthand, flagTag.convertValue("float32").(float32), flagTag.Usage)
			} else {
				flag.Float32Var((*float32)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.convertValue("float32").(float32), flagTag.Usage)
			}
		case reflect.Float64:
			if flagTag.Shorthand != "" && flagTag.Shorthand != "-" {
				flag.Float64VarP((*float64)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.Shorthand, flagTag.convertValue("float64").(float64), flagTag.Usage)
			} else {
				flag.Float64Var((*float64)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.convertValue("float64").(float64), flagTag.Usage)
			}
		case reflect.Bool:
			if flagTag.Shorthand != "" && flagTag.Shorthand != "-" {
				flag.BoolVarP((*bool)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.Shorthand, flagTag.convertValue("bool").(bool), flagTag.Usage)
			} else {
				flag.BoolVar((*bool)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.convertValue("bool").(bool), flagTag.Usage)
			}
		default:
			panic(fmt.Sprintf("unsupported type: %T", fv.Interface()))
		}
	}
	return nil
}

func MustBindFlags(flag *pflag.FlagSet, a any) {
	err := BindFlags(flag, a)
	if err != nil {
		panic(err)
	}
}
