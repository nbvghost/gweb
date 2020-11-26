/*
DATE: 2020/11/14
AUTHOR: wushunqing

hhmmss 工具类

*/
package hhmmss

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type _HhMmSs [3]uint

func NewHhMmSsWithSecond(second int) _HhMmSs {
	hh := second / 60 / 60 % 24
	mm := second / 60 % 60
	ss := second % 60
	return _HhMmSs{uint(hh), uint(mm), uint(ss)}
}

//hhmmss 格式为02:12:12或201212
func NewHhMmSs(hms string) (_HhMmSs, error) {
	hhmmss := [3]int{}

	if strings.Contains(hms, ":") {
		timeArray := strings.Split(string(hms), ":")
		if len(timeArray) != 3 {
			return _HhMmSs{}, errors.New("时间格式不正确")
		}
		var err error
		hhmmss[0], err = strconv.Atoi(timeArray[0])
		if err != nil {
			return _HhMmSs{}, err
		}
		hhmmss[1], err = strconv.Atoi(timeArray[1])
		if err != nil {
			return _HhMmSs{}, err
		}
		hhmmss[2], err = strconv.Atoi(timeArray[2])
		if err != nil {
			return _HhMmSs{}, err
		}

		if hhmmss[0] >= 24 {
			return _HhMmSs{}, errors.New("小时格式不正确")
		}
		if hhmmss[1] >= 60 {
			return _HhMmSs{}, errors.New("分钟格式不正确")
		}
		if hhmmss[2] >= 60 {
			return _HhMmSs{}, errors.New("秒钟格式不正确")
		}
		return _HhMmSs{uint(hhmmss[0]), uint(hhmmss[1]), uint(hhmmss[2])}, nil
	} else {
		var err error
		if len(hms) != 3*2 {
			return _HhMmSs{}, errors.New("时间格式不正确,必须是6位")
		}

		hhmmss[0], err = strconv.Atoi(hms[0:2])
		if err != nil {
			return [3]uint{}, err
		}
		hhmmss[1], err = strconv.Atoi(hms[2:4])
		if err != nil {
			return [3]uint{}, err
		}
		hhmmss[2], err = strconv.Atoi(hms[4:6])
		if err != nil {
			return [3]uint{}, err
		}

		if hhmmss[0] >= 24 {
			return _HhMmSs{}, errors.New("小时格式不正确")
		}
		if hhmmss[1] >= 60 {
			return _HhMmSs{}, errors.New("分钟格式不正确")
		}
		if hhmmss[2] >= 60 {
			return _HhMmSs{}, errors.New("秒钟格式不正确")
		}

		return _HhMmSs{uint(hhmmss[0]), uint(hhmmss[1]), uint(hhmmss[2])}, nil
	}
}

func (hms _HhMmSs) ToHhMmSs() (hh, mm, ss int) {

	return int(hms[0]), int(hms[1]), int(hms[2])
}

//返回总秒数
func (hms _HhMmSs) Second() int {
	hh := hms[0]
	mm := hms[1]
	ss := hms[2]
	return int(ss + (mm * 60) + (hh * 60 * 60))
}

//返回总秒数
func (hms _HhMmSs) Add(second int) _HhMmSs {
	s := hms.Second() + second
	if s > (24 * 60 * 60) {
		s = s - (24 * 60 * 60)
	}
	return NewHhMmSsWithSecond(s)
}

//将时间转成今天的时间
func (hms _HhMmSs) ToDay() time.Time {

	hh, mm, ss := hms.ToHhMmSs()
	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day(), hh, mm, ss, 0, now.Location())
	return t
}
