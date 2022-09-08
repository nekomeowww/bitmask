package bitmask

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var (
	// ErrInvalidTagValue tag 值不合法
	ErrInvalidTagValue = errors.New("invalid tag value")
	// ErrNoNilPointerAllowed 不允许 nil 指针
	ErrNoNilPointerAllowed = errors.New("no nil pointer allowed")
	// ErrNotAStructPointer 不是结构体指针
	ErrNotAStructPointer = errors.New("dest must be a pointer to struct")
	// ErrNotStructPointerOrStruct 不是结构体指针或结构体
	ErrNotStructPointerOrStruct = errors.New("dest must be a pointer to struct or a struct")
	// ErrBitOutOfRange bit 超出范围
	ErrBitOutOfRange = errors.New("bit should within 1 to 64")
)

const (
	// TagName 标签名
	TagName = "bitmask"
)

// BitMask 最大支持 64 位
type BitMask int

// New 新建一个 BitMask
func New(val int) BitMask {
	b := BitMask(val)
	return b
}

// IsSet 判断第 i 位是否为 1，0 <= i <= 32
func (b BitMask) IsSet(i int) bool {
	return b&(1<<(i-1)) != 0
}

// Set 翻转第 i 位置为 1，0 <= i <= 32
func (b *BitMask) Set(i int) {
	*b |= (1 << (i - 1))
}

// Unset 翻转第 i 位置为 0，0 <= i <= 32
func (b *BitMask) Unset(i int) {
	*b &^= (1 << (i - 1))
}

// Marshal 将结构体转换为 BitMask
func Marshal(dest any) (BitMask, error) {
	var bitmask BitMask
	if dest == nil {
		return bitmask, nil
	}

	v := reflect.ValueOf(dest)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return bitmask, nil
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		return Marshal(v.Interface())
	}
	if v.Kind() != reflect.Struct {
		return bitmask, ErrNotStructPointerOrStruct
	}

	destType := v.Type()
	for i := 0; i < destType.NumField(); i++ {
		field := destType.Field(i)
		tagValue := field.Tag.Get(TagName)
		if tagValue == "" {
			continue
		}

		value, err := strconv.ParseInt(tagValue, 10, 64)
		if err != nil {
			return bitmask, fmt.Errorf("%w: %s", ErrInvalidTagValue, tagValue)
		}
		if value < 1 || value > 64 {
			return bitmask, fmt.Errorf("%w: %s", ErrBitOutOfRange, tagValue)
		}

		switch field.Type.Kind() {
		case reflect.Bool:
			if v.Field(i).Bool() {
				bitmask.Set(int(value))
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v.Field(i).Int() != 0 {
				bitmask.Set(int(value))
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if v.Field(i).Uint() != 0 {
				bitmask.Set(int(value))
			}
		}
	}

	return bitmask, nil
}

// Unmarshal 将 BitMask 转换到结构体
func Unmarshal(bitmask BitMask, dest any) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return ErrNotAStructPointer
	}

	v = v.Elem()
	if v.Kind() == reflect.Ptr {
		return Unmarshal(bitmask, v.Interface())
	}
	if v.Kind() != reflect.Struct {
		return ErrNotAStructPointer
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tagValue := f.Tag.Get(TagName)
		if tagValue == "" {
			continue
		}

		value, err := strconv.ParseInt(tagValue, 10, 64)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrInvalidTagValue, tagValue)
		}
		if value < 1 || value > 64 {
			return fmt.Errorf("%w: %s", ErrBitOutOfRange, tagValue)
		}

		switch f.Type.Kind() {
		case reflect.Bool:
			if bitmask.IsSet(int(value)) {
				v.Field(i).SetBool(true)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if bitmask.IsSet(int(value)) {
				v.Field(i).SetInt(1)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if bitmask.IsSet(int(value)) {
				v.Field(i).SetUint(1)
			}
		}

	}

	return nil
}
