package main

import (
	"testing"

	"github.com/qkja/gobase/validate"
)

type ValueModelIdCardEntity struct {
	Data string `match:"model=id_card"`
}

type ValueModelPhone struct {
	Data string `match:"model=phone"`
}

type ValueModelFixedPhoneEntity struct {
	Data string `match:"model=fixed_phone"`
}

type ValueModelEmailEntity struct {
	Data string `match:"model=mail"`
}

type ValueModelIpAddressEntity struct {
	Data string `match:"model=ip"`
}

// 身份证号
func TestModelIdCard(t *testing.T) {
	var value ValueModelIdCardEntity
	var result bool
	var err string

	// 测试 异常情况
	value = ValueModelIdCardEntity{Data: "4109281002226311"}
	result, err = validate.Check(value)
	Equal(t, err, "属性 Data 的值 4109281002226311 不符合身份证要求", result, false)

	// 测试 异常情况
	value = ValueModelIdCardEntity{Data: "28712381"}
	result, err = validate.Check(value)
	Equal(t, err, "属性 Data 的值 28712381 不符合身份证要求", result, false)
}

// 手机号
func TestModelPhone(t *testing.T) {
	var value ValueModelPhone
	var result bool
	var err string

	// 测试 正常情况
	value = ValueModelPhone{Data: "15700092345"}
	result, err = validate.Check(value)
	TrueErr(t, result, err)

	// 测试 正常情况
	value = ValueModelPhone{Data: "15200092345"}
	result, err = validate.Check(value)
	TrueErr(t, result, err)

	// 测试 异常情况
	value = ValueModelPhone{Data: "28712381887"}
	result, err = validate.Check(value)
	Equal(t, err, "属性 Data 的值 28712381887 没有命中只允许类型 [phone]", result, false)
}

// 手机号
func TestModelPhone2(t *testing.T) {
	var value ValueModelPhone
	var result bool
	var err string

	// 测试 正常情况
	value = ValueModelPhone{Data: "14500092345"}
	result, err = validate.Check(value)
	TrueErr(t, result, err)

	// 测试 异常情况
	value = ValueModelPhone{Data: "28712381"}
	result, err = validate.Check(value)
	Equal(t, err, "属性 Data 的值 28712381 没有命中只允许类型 [phone]", result, false)
}

// 固定电话
func TestModelFixedPhone(t *testing.T) {
	var value ValueModelFixedPhoneEntity
	var result bool
	var err string

	// 测试 正常情况
	value = ValueModelFixedPhoneEntity{Data: "0393-3879765"}
	result, err = validate.Check(value)
	TrueErr(t, result, err)

	// 测试 异常情况
	value = ValueModelFixedPhoneEntity{Data: "1387772"}
	result, err = validate.Check(value)
	Equal(t, err, "属性 Data 的值 1387772 没有命中只允许类型 [fixed_phone]", result, false)
}

// 邮箱
func TestModelMail(t *testing.T) {
	var value ValueModelEmailEntity
	var result bool
	var err string

	// 测试 正常情况
	value = ValueModelEmailEntity{Data: "123lan@163.com"}
	result, err = validate.Check(value)
	TrueErr(t, result, err)

	// 测试 异常情况
	value = ValueModelEmailEntity{Data: "123@"}
	result, err = validate.Check(value)
	Equal(t, err, "属性 Data 的值 123@ 没有命中只允许类型 [mail]", result, false)
}

// ip地址
func TestModelIpAddress(t *testing.T) {
	var value ValueModelIpAddressEntity
	var result bool
	var err string

	// 测试 正常情况
	value = ValueModelIpAddressEntity{Data: "192.123.231.222"}
	result, err = validate.Check(value)
	TrueErr(t, result, err)

	// 测试 异常情况
	value = ValueModelIpAddressEntity{Data: "123.231.222"}
	result, err = validate.Check(value)
	Equal(t, err, "属性 Data 的值 123.231.222 没有命中只允许类型 [ip]", result, false)

	// 测试 异常情况
	value = ValueModelIpAddressEntity{Data: "192.123.231.adf"}
	result, err = validate.Check(value)
	Equal(t, err, "属性 Data 的值 192.123.231.adf 没有命中只允许类型 [ip]", result, false)
}

// 国际手机号：按国家/地区校验（google/libphonenumber）
type ValueModelPhoneUS struct {
	Data string `match:"model=phone:US"`
}

type ValueModelPhoneZZ struct {
	Data string `match:"model=phone:ZZ"`
}

// 手机号：指定国家/地区（model=phone:US）
func TestModelPhoneLocale(t *testing.T) {
	var value ValueModelPhoneUS
	var result bool
	var err string

	// 美国号 → 通过
	value = ValueModelPhoneUS{Data: "2015550123"}
	result, err = validate.Check(value)
	TrueErr(t, result, err)

	// 大陆号在 US 下 → 拒绝
	value = ValueModelPhoneUS{Data: "13800138000"}
	result, err = validate.Check(value)
	Equal(t, err, "属性 Data 的值 13800138000 没有命中只允许类型 [phone]", result, false)
}

// 手机号：国际号码（model=phone:ZZ，号码须带国家码）
func TestModelPhoneInternational(t *testing.T) {
	var value ValueModelPhoneZZ
	var result bool
	var err string

	// 带国家码 → 通过
	value = ValueModelPhoneZZ{Data: "+8613800138000"}
	result, err = validate.Check(value)
	TrueErr(t, result, err)

	// 不带国家码 → 拒绝
	value = ValueModelPhoneZZ{Data: "13800138000"}
	result, err = validate.Check(value)
	Equal(t, err, "属性 Data 的值 13800138000 没有命中只允许类型 [phone]", result, false)
}

// 邮箱边界（govalidator.IsEmail）
func TestModelMailBoundary(t *testing.T) {
	var value ValueModelEmailEntity
	var result bool
	var err string

	// 合法邮箱 → 通过
	value = ValueModelEmailEntity{Data: "abc@example.com"}
	result, err = validate.Check(value)
	TrueErr(t, result, err)

	// 缺少域名 → 拒绝
	value = ValueModelEmailEntity{Data: "abc"}
	result, err = validate.Check(value)
	Equal(t, err, "属性 Data 的值 abc 没有命中只允许类型 [mail]", result, false)
}
