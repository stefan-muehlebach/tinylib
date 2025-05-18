package tinylib

import (
	"time"
)

//----------------------------------------------------------------------------

type Mappable interface {
	~int | ~int16 | ~uint16 | ~float64 | time.Duration
}

func Map[IT, OT Mappable](val, srcMin, srcMax IT, dstMin, dstMax OT) OT {
	return dstMin + OT(val-srcMin)*(dstMax-dstMin)/OT(srcMax-srcMin)
}

//----------------------------------------------------------------------------

func millis() int64 {
	return time.Now().UnixMilli()
}
