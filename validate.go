package util

import (
	"regexp"
	"strconv"
)

// NewPassword 生成密码
func NewPassword(pwd string, code string) string {
	return NewMD5(code + Base64(pwd))
}

// VerifyPassword 验证密码，一致返回true，否则返回false
func VerifyPassword(n string, o string, code string) bool {
	n = NewPassword(n, code)
	if n == o {
		return true
	}
	return false
}

// CheckMobile 校验手机号(只校验大陆手机号)
func CheckMobile(mobile string, countryCode int) bool {
	if countryCode == 86 {
		reg := regexp.MustCompile(`^1[23456789]\d{9}$`)
		return reg.MatchString(mobile)
	}
	return true
}

// CheckPassword 校验密码
func CheckPassword(password string) bool {
	reg := regexp.MustCompile(`^\w{6,16}$`)
	return reg.MatchString(password)
}

// CheckIDCard 校验身份证
func CheckIDCard(idCard string) bool {
	if len(idCard) != 18 {
		return false
	}

	var idCardArr [18]byte // 'X' == byte(88)， 'X'在byte中表示为88
	var idCardArrCopy [17]byte

	// 将字符串，转换成[]byte,arrIdCard数组当中
	for k, v := range []byte(idCard) {
		idCardArr[k] = byte(v)
	}

	//arrIdCard[18]前17位元素到arrIdCardCopy数组当中
	for j := 0; j < 17; j++ {
		idCardArrCopy[j] = idCardArr[j]
	}

	checkID := func(id [17]byte) int {
		arr := make([]int, 17)
		for index, value := range id {
			arr[index], _ = strconv.Atoi(string(value))
		}

		wi := [...]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
		var res int
		for i := 0; i < 17; i++ {
			res += arr[i] * wi[i]
		}
		return res % 11
	}

	byte2int := func(x byte) byte {
		if x == 88 {
			return 'X'
		}
		return x - 48 // 'X' - 48 = 40;
	}

	verify := checkID(idCardArrCopy)
	last := byte2int(idCardArr[17])
	var temp byte
	var i int
	a18 := [11]byte{1, 0, 'X', 9, 8, 7, 6, 5, 4, 3, 2}

	for i = 0; i < 11; i++ {
		if i == verify {
			temp = a18[i]
			break
		}
	}

	if temp == last {
		return true
	}
	return false

}
