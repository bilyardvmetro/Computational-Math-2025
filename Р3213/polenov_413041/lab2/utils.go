package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func Abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func MaxAbs(x float64, y float64) float64 {
	if Abs(x) > Abs(y) {
		return Abs(x)
	}
	return Abs(y)
}

// FastPow Быстрое возведение в степень
func FastPow(x float64, k int) float64 {
	if k == 0 {
		return 1
	} else if k == 1 {
		return x
	} else if k%2 == 1 {
		return FastPow(x, k-1) * x
	}
	tmp := FastPow(x, k/2)
	return tmp * tmp
}

// GetOut Функция выхода из программы с ошибкой
func GetOut(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func CheckIsFloat(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func CheckIsInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func StrToFloat(s string) float64 {
	number, _ := strconv.ParseFloat(s, 64)
	return number
}

func StrToInt(s string) int {
	number, _ := strconv.Atoi(s)
	return number
}

func ReadInt(in *bufio.Reader, tmp *int, prefix bool) {
	var tmpS string
	fmt.Fscan(in, &tmpS)
	if !CheckIsInt(tmpS) {
		GetOut(OptionError{})
	}
	*tmp = StrToInt(tmpS)
	if prefix {
		in.ReadLine()
	}
}

func ReadFloat(in *bufio.Reader, tmp *float64, prefix bool, textError string) {
	var x string
	fmt.Fscan(in, &x)
	if !CheckIsFloat(x) {
		GetOut(ParseError{textError})
	}
	*tmp = StrToFloat(x)
	if prefix {
		in.ReadLine()
	}
}
