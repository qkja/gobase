package matcher

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/nyaruka/phonenumbers"

	"github.com/qkja/gobase/constants"

	"github.com/qkja/gobase/logger"
)

type ModelMatch struct {
	BlackWhiteMatch

	isIdCard  bool
	modelName string
	check     func(value string) bool
}

var modelMap = map[string]func(value string) bool{}

func (modelMatch *ModelMatch) Match(_ map[string]interface{}, _ any, field reflect.StructField, fieldValue any) bool {
	if nil == fieldValue {
		return false
	}

	if field.Type.Kind() != reflect.String {
		return false
	}

	if modelMatch.check == nil {
		return false
	}

	if modelMatch.check(fmt.Sprintf("%v", fieldValue)) {
		if modelMatch.isIdCard {
			modelMatch.SetBlackMsg("属性 %v 的值 %v 符合身份证要求", field.Name, fieldValue)
		} else {
			modelMatch.SetBlackMsg("属性 %v 的值 %v 命中不允许的类型 [%v]", field.Name, fieldValue, modelMatch.modelName)
		}
		return true
	}

	if modelMatch.isIdCard {
		modelMatch.SetWhiteMsg("属性 %v 的值 %v 不符合身份证要求", field.Name, fieldValue)
	} else {
		modelMatch.SetWhiteMsg("属性 %v 的值 %v 没有命中只允许类型 [%v]", field.Name, fieldValue, modelMatch.modelName)
	}
	return false
}

func (modelMatch *ModelMatch) IsEmpty() bool {
	return modelMatch.check == nil
}

func BuildModelMatcher(objectTypeFullName string, fieldKind reflect.Kind, objectFieldName string, tagName string, subCondition string, errMsg string) {
	if constants.MATCH != tagName {
		return
	}

	if fieldKind == reflect.Slice {
		return
	}

	if !strings.Contains(subCondition, constants.Model) || !strings.Contains(subCondition, constants.EQUAL) {
		return
	}

	index := strings.Index(subCondition, "=")
	modelKey := strings.TrimSpace(subCondition[index+1:])

	// 手机号：支持 tag 内冒号指定国家/地区（google/libphonenumber）。
	//   model=phone       默认 CN（中国大陆 11 位）
	//   model=phone:US    按美国号段校验
	//   model=phone:ZZ    国际号码，号码须带国家码（如 +8613800138000）
	if modelKey == constants.Phone || strings.HasPrefix(modelKey, constants.Phone+":") {
		region := "CN"
		if strings.Contains(modelKey, ":") {
			region = strings.SplitN(modelKey, ":", 2)[1]
		}
		check := func(v string) bool {
			num, err := phonenumbers.Parse(v, region)
			if err != nil {
				return false
			}
			return phonenumbers.IsValidNumber(num)
		}
		addMatcher(objectTypeFullName, objectFieldName, &ModelMatch{check: check, modelName: constants.Phone}, errMsg, true)
		return
	}

	check, contain := modelMap[modelKey]
	if !contain {
		logger.Error("不包含模式%v", modelKey)
		return
	}

	if modelKey == constants.IdCard {
		addMatcher(objectTypeFullName, objectFieldName, &ModelMatch{check: check, isIdCard: true, modelName: modelKey}, errMsg, true)
	} else {
		addMatcher(objectTypeFullName, objectFieldName, &ModelMatch{check: check, modelName: modelKey}, errMsg, true)
	}
}

func init() {
	// 固定电话
	fixedPhoneReg, _ := regexp.Compile("^(([0+]\\d{2,3}-)?(0\\d{2,3})-)(\\d{7,8})(-(\\d{3,}))?$")
	modelMap[constants.FixedPhone] = func(v string) bool { return fixedPhoneReg.MatchString(v) }

	// 邮箱：govalidator 战备实现（含国际化域名）
	modelMap[constants.MAIL] = govalidator.IsEmail

	// IP地址
	ipReg, _ := regexp.Compile("^((25[0-5]|2[0-4]\\d|[01]?\\d\\d?)\\.){3}(25[0-5]|2[0-4]\\d|[01]?\\d\\d?)$")
	modelMap[constants.IpAddress] = func(v string) bool { return ipReg.MatchString(v) }

	// 身份证号：算法校验（校验码 + 出生日期）
	modelMap[constants.IdCard] = func(v string) bool { return idCardIsValidate(v) }
}
