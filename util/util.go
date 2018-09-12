package util

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"time"
)

func Md5(s string) string {
	data := []byte(s)
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	return hex.EncodeToString(md5Ctx.Sum(nil))
}

func GraceFloat64P4(f float64) float64 {
	if f > 0 {
		return float64(int((f+0.00005)*10000)) / 10000.0
	} else {
		return float64(int((f-0.00005)*10000)) / 10000.0
	}
}

func GraceFloat64P3(f float64) float64 {
	if f > 0 {
		return float64(int((f+0.0005)*1000)) / 1000.0
	} else {
		return float64(int((f-0.0005)*1000)) / 1000.0
	}
}

// 每日24点（次日0点）的timer
func Start24Timer(f func()) {
	go func() {
		for {
			now := time.Now()
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			<-t.C
			f()
		}
	}()
}

func StartTimer(duration time.Duration, f func()) {
	go func() {
		for {
			now := time.Now()
			next := now.Add(duration)
			t := time.NewTimer(next.Sub(now))
			<-t.C
			f()
		}
	}()
}

func Float2String(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func String2Float(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
