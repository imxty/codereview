package errors

import (
	"regexp"
	"strconv"

	"github.com/jinmukeji/plat-pkg/v4/micro/errors/codes"
)

var (
	// 定义正则表达式
	re = regexp.MustCompile(`errcode:\s*(\d+)`)
)

// GetErrCode
func GetErrCode(err error) codes.Code {
	matches := re.FindStringSubmatch(err.Error())
	if len(matches) > 1 {
		cd, err := strconv.Atoi(matches[1])
		if err != nil {
			return codes.Unknown
		}
		return codes.Code(cd)
	}
	return codes.Unknown
}
