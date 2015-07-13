package main

import (
	"time"
)

const (
	WorkerIdBits      = 5
	WorkIdMax         = 1<<WorkerIdBits - 1
	WorkIdShift       = WorkerIdBits
	DataCenterIdBits  = 5
	DataCenterIdMax   = 1<<DataCenterIdBits - 1
	DataCenterIdShift = WorkerIdBits + DataCenterIdBits
	SequenceBits      = 12
	SequenceMask      = 1<<SequenceBits - 1
	TimestampShift    = WorkerIdBits + DataCenterIdBits + SequenceBits
	MillisEpoch       = 1436715587000
)

type Worker struct {
	WorkerId      uint32
	DataCenterId  uint32
	Sequence      uint32
	LastTimestamp int64
	Out           chan uint64
}

func NewWorker(workerId uint32, dataCenterId uint32, out chan uint64) (*Worker, error) {
	if workerId > WorkIdMax || dataCenterId > DataCenterIdMax {
		return nil, IllegalArgumentError
	}
	w := &Worker{WorkerId: workerId, DataCenterId: dataCenterId, Sequence: 0, Out: out}
	return w, nil
}

func currentMillis() int64 {
	return time.Now().UnixNano()/1000000 - MillisEpoch
}

func (w *Worker) waitNextMillis() int64 {
	timestamp := currentMillis()
	// Busy wait
	for timestamp < w.LastTimestamp {
		timestamp = currentMillis()
	}
	return timestamp
}

func (w *Worker) YieldId() {
	w.Out <- w.GenerateId()
}

func (w *Worker) GenerateId() uint64 {
	timestamp := currentMillis()
	if timestamp < w.LastTimestamp {
		panic("timestamp less than w.LastTimestamp.")
	}

	if timestamp == w.LastTimestamp {
		w.Sequence = (w.Sequence + 1) & SequenceMask
		if w.Sequence == 0 {
			w.waitNextMillis()
		}
	} else {
		w.Sequence = 0
	}

	w.LastTimestamp = timestamp
	return uint64(w.LastTimestamp<<TimestampShift) | uint64(w.DataCenterId<<DataCenterIdShift) | uint64(w.WorkerId<<WorkIdShift) | uint64(w.Sequence)
}
