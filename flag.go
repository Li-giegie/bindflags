package bindflags

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"strings"
)

type IFlagTag interface {
	GetName() string
	GetValue() string
	GetUsage() string
}

type GetFlagTag interface {
	GetFlagTag() IFlagTag
}

func BindFlags(f *flag.FlagSet, a any, group ...string) error {
	rv := reflect.ValueOf(a)
	if rv.Kind() != reflect.Ptr {
		return errors.New("a must be a pointer")
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return errors.New("a must be a struct")
	}
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
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
			if err = BindFlags(f, fv.Addr().Interface(), group...); err != nil {
				return err
			}
		case reflect.String:
			f.StringVar((*string)(fv.Addr().UnsafePointer()), flagTag.Name, flagTag.Value, flagTag.Usage)
		case reflect.Int:
			f.IntVar((*int)(fv.Addr().UnsafePointer()), flagTag.Name, convertValue(flagTag.Name, flagTag.Value, "int").(int), flagTag.Usage)
		case reflect.Int64:
			f.Int64Var((*int64)(fv.Addr().UnsafePointer()), flagTag.Name, convertValue(flagTag.Name, flagTag.Value, "int64").(int64), flagTag.Usage)
		case reflect.Uint:
			f.UintVar((*uint)(fv.Addr().UnsafePointer()), flagTag.Name, convertValue(flagTag.Name, flagTag.Value, "uint").(uint), flagTag.Usage)
		case reflect.Uint64:
			f.Uint64Var((*uint64)(fv.Addr().UnsafePointer()), flagTag.Name, convertValue(flagTag.Name, flagTag.Value, "uint64").(uint64), flagTag.Usage)
		case reflect.Float64:
			f.Float64Var((*float64)(fv.Addr().UnsafePointer()), flagTag.Name, convertValue(flagTag.Name, flagTag.Value, "float64").(float64), flagTag.Usage)
		case reflect.Bool:
			f.BoolVar((*bool)(fv.Addr().UnsafePointer()), flagTag.Name, convertValue(flagTag.Name, flagTag.Value, "bool").(bool), flagTag.Usage)
		default:
			panic(fmt.Sprintf("unsupported type: %T", fv.Interface()))
		}
	}
	return nil
}

func MustBindFlags(f *flag.FlagSet, a any, group ...string) {
	err := BindFlags(f, a, group...)
	if err != nil {
		panic(err)
	}
}
