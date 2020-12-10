// Code generated by execgen; DO NOT EDIT.
// Copyright 2018 The Cockroach Authors.
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package colexecagg

import (
	"time"
	"unsafe"

	"github.com/cockroachdb/apd/v2"
	"github.com/cockroachdb/cockroach/pkg/col/coldata"
	"github.com/cockroachdb/cockroach/pkg/col/coldataext"
	"github.com/cockroachdb/cockroach/pkg/col/typeconv"
	"github.com/cockroachdb/cockroach/pkg/sql/colmem"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
	"github.com/cockroachdb/cockroach/pkg/util/duration"
	"github.com/cockroachdb/errors"
)

func newAnyNotNullHashAggAlloc(
	allocator *colmem.Allocator, t *types.T, allocSize int64,
) (aggregateFuncAlloc, error) {
	allocBase := aggAllocBase{allocator: allocator, allocSize: allocSize}
	switch typeconv.TypeFamilyToCanonicalTypeFamily(t.Family()) {
	case types.BoolFamily:
		switch t.Width() {
		case -1:
		default:
			return &anyNotNullBoolHashAggAlloc{aggAllocBase: allocBase}, nil
		}
	case types.BytesFamily:
		switch t.Width() {
		case -1:
		default:
			return &anyNotNullBytesHashAggAlloc{aggAllocBase: allocBase}, nil
		}
	case types.DecimalFamily:
		switch t.Width() {
		case -1:
		default:
			return &anyNotNullDecimalHashAggAlloc{aggAllocBase: allocBase}, nil
		}
	case types.IntFamily:
		switch t.Width() {
		case 16:
			return &anyNotNullInt16HashAggAlloc{aggAllocBase: allocBase}, nil
		case 32:
			return &anyNotNullInt32HashAggAlloc{aggAllocBase: allocBase}, nil
		case -1:
		default:
			return &anyNotNullInt64HashAggAlloc{aggAllocBase: allocBase}, nil
		}
	case types.FloatFamily:
		switch t.Width() {
		case -1:
		default:
			return &anyNotNullFloat64HashAggAlloc{aggAllocBase: allocBase}, nil
		}
	case types.TimestampTZFamily:
		switch t.Width() {
		case -1:
		default:
			return &anyNotNullTimestampHashAggAlloc{aggAllocBase: allocBase}, nil
		}
	case types.IntervalFamily:
		switch t.Width() {
		case -1:
		default:
			return &anyNotNullIntervalHashAggAlloc{aggAllocBase: allocBase}, nil
		}
	case typeconv.DatumVecCanonicalTypeFamily:
		switch t.Width() {
		case -1:
		default:
			return &anyNotNullDatumHashAggAlloc{aggAllocBase: allocBase}, nil
		}
	}
	return nil, errors.Errorf("unsupported any not null agg type %s", t.Name())
}

// anyNotNullBoolHashAgg implements the ANY_NOT_NULL aggregate, returning the
// first non-null value in the input column.
type anyNotNullBoolHashAgg struct {
	hashAggregateFuncBase
	col                         coldata.Bools
	curAgg                      bool
	foundNonNullForCurrentGroup bool
}

var _ AggregateFunc = &anyNotNullBoolHashAgg{}

func (a *anyNotNullBoolHashAgg) SetOutput(vec coldata.Vec) {
	a.hashAggregateFuncBase.SetOutput(vec)
	a.col = vec.Bool()
}

func (a *anyNotNullBoolHashAgg) Compute(
	vecs []coldata.Vec, inputIdxs []uint32, inputLen int, sel []int,
) {
	if a.foundNonNullForCurrentGroup {
		// We have already seen non-null for the current group, and since there
		// is at most a single group when performing hash aggregation, we can
		// finish computing.
		return
	}

	var oldCurAggSize uintptr
	vec := vecs[inputIdxs[0]]
	col, nulls := vec.Bool(), vec.Nulls()
	a.allocator.PerformOperation([]coldata.Vec{a.vec}, func() {
		// Capture col to force bounds check to work. See
		// https://github.com/golang/go/issues/39756
		col := col
		_ = col.Get(inputLen - 1)
		{
			sel = sel[:inputLen]
			if nulls.MaybeHasNulls() {
				for _, i := range sel {

					var isNull bool
					isNull = nulls.NullAt(i)
					if !a.foundNonNullForCurrentGroup && !isNull {
						// If we haven't seen any non-nulls for the current group yet, and the
						// current value is non-null, then we can pick the current value to be
						// the output.
						val := col.Get(i)
						a.curAgg = val
						a.foundNonNullForCurrentGroup = true
						// We have already seen non-null for the current group, and since there
						// is at most a single group when performing hash aggregation, we can
						// finish computing.
						return
					}
				}
			} else {
				for _, i := range sel {

					var isNull bool
					isNull = false
					if !a.foundNonNullForCurrentGroup && !isNull {
						// If we haven't seen any non-nulls for the current group yet, and the
						// current value is non-null, then we can pick the current value to be
						// the output.
						val := col.Get(i)
						a.curAgg = val
						a.foundNonNullForCurrentGroup = true
						// We have already seen non-null for the current group, and since there
						// is at most a single group when performing hash aggregation, we can
						// finish computing.
						return
					}
				}
			}
		}
	},
	)
	var newCurAggSize uintptr
	if newCurAggSize != oldCurAggSize {
		a.allocator.AdjustMemoryUsage(int64(newCurAggSize - oldCurAggSize))
	}
}

func (a *anyNotNullBoolHashAgg) Flush(outputIdx int) {
	// If we haven't found any non-nulls for this group so far, the output for
	// this group should be null.
	if !a.foundNonNullForCurrentGroup {
		a.nulls.SetNull(outputIdx)
	} else {
		a.col[outputIdx] = a.curAgg
	}
}

type anyNotNullBoolHashAggAlloc struct {
	aggAllocBase
	aggFuncs []anyNotNullBoolHashAgg
}

var _ aggregateFuncAlloc = &anyNotNullBoolHashAggAlloc{}

const sizeOfAnyNotNullBoolHashAgg = int64(unsafe.Sizeof(anyNotNullBoolHashAgg{}))
const anyNotNullBoolHashAggSliceOverhead = int64(unsafe.Sizeof([]anyNotNullBoolHashAgg{}))

func (a *anyNotNullBoolHashAggAlloc) newAggFunc() AggregateFunc {
	if len(a.aggFuncs) == 0 {
		a.allocator.AdjustMemoryUsage(anyNotNullBoolHashAggSliceOverhead + sizeOfAnyNotNullBoolHashAgg*a.allocSize)
		a.aggFuncs = make([]anyNotNullBoolHashAgg, a.allocSize)
	}
	f := &a.aggFuncs[0]
	f.allocator = a.allocator
	a.aggFuncs = a.aggFuncs[1:]
	return f
}

// anyNotNullBytesHashAgg implements the ANY_NOT_NULL aggregate, returning the
// first non-null value in the input column.
type anyNotNullBytesHashAgg struct {
	hashAggregateFuncBase
	col                         *coldata.Bytes
	curAgg                      []byte
	foundNonNullForCurrentGroup bool
}

var _ AggregateFunc = &anyNotNullBytesHashAgg{}

func (a *anyNotNullBytesHashAgg) SetOutput(vec coldata.Vec) {
	a.hashAggregateFuncBase.SetOutput(vec)
	a.col = vec.Bytes()
}

func (a *anyNotNullBytesHashAgg) Compute(
	vecs []coldata.Vec, inputIdxs []uint32, inputLen int, sel []int,
) {
	if a.foundNonNullForCurrentGroup {
		// We have already seen non-null for the current group, and since there
		// is at most a single group when performing hash aggregation, we can
		// finish computing.
		return
	}

	oldCurAggSize := len(a.curAgg)
	vec := vecs[inputIdxs[0]]
	col, nulls := vec.Bytes(), vec.Nulls()
	a.allocator.PerformOperation([]coldata.Vec{a.vec}, func() {
		// Capture col to force bounds check to work. See
		// https://github.com/golang/go/issues/39756
		col := col
		_ = col.Get(inputLen - 1)
		{
			sel = sel[:inputLen]
			if nulls.MaybeHasNulls() {
				for _, i := range sel {

					var isNull bool
					isNull = nulls.NullAt(i)
					if !a.foundNonNullForCurrentGroup && !isNull {
						// If we haven't seen any non-nulls for the current group yet, and the
						// current value is non-null, then we can pick the current value to be
						// the output.
						val := col.Get(i)
						a.curAgg = append(a.curAgg[:0], val...)
						a.foundNonNullForCurrentGroup = true
						// We have already seen non-null for the current group, and since there
						// is at most a single group when performing hash aggregation, we can
						// finish computing.
						return
					}
				}
			} else {
				for _, i := range sel {

					var isNull bool
					isNull = false
					if !a.foundNonNullForCurrentGroup && !isNull {
						// If we haven't seen any non-nulls for the current group yet, and the
						// current value is non-null, then we can pick the current value to be
						// the output.
						val := col.Get(i)
						a.curAgg = append(a.curAgg[:0], val...)
						a.foundNonNullForCurrentGroup = true
						// We have already seen non-null for the current group, and since there
						// is at most a single group when performing hash aggregation, we can
						// finish computing.
						return
					}
				}
			}
		}
	},
	)
	newCurAggSize := len(a.curAgg)
	if newCurAggSize != oldCurAggSize {
		a.allocator.AdjustMemoryUsage(int64(newCurAggSize - oldCurAggSize))
	}
}

func (a *anyNotNullBytesHashAgg) Flush(outputIdx int) {
	// If we haven't found any non-nulls for this group so far, the output for
	// this group should be null.
	if !a.foundNonNullForCurrentGroup {
		a.nulls.SetNull(outputIdx)
	} else {
		a.col.Set(outputIdx, a.curAgg)
	}
	// Release the reference to curAgg eagerly.
	a.allocator.AdjustMemoryUsage(-int64(len(a.curAgg)))
	a.curAgg = nil
}

type anyNotNullBytesHashAggAlloc struct {
	aggAllocBase
	aggFuncs []anyNotNullBytesHashAgg
}

var _ aggregateFuncAlloc = &anyNotNullBytesHashAggAlloc{}

const sizeOfAnyNotNullBytesHashAgg = int64(unsafe.Sizeof(anyNotNullBytesHashAgg{}))
const anyNotNullBytesHashAggSliceOverhead = int64(unsafe.Sizeof([]anyNotNullBytesHashAgg{}))

func (a *anyNotNullBytesHashAggAlloc) newAggFunc() AggregateFunc {
	if len(a.aggFuncs) == 0 {
		a.allocator.AdjustMemoryUsage(anyNotNullBytesHashAggSliceOverhead + sizeOfAnyNotNullBytesHashAgg*a.allocSize)
		a.aggFuncs = make([]anyNotNullBytesHashAgg, a.allocSize)
	}
	f := &a.aggFuncs[0]
	f.allocator = a.allocator
	a.aggFuncs = a.aggFuncs[1:]
	return f
}

// anyNotNullDecimalHashAgg implements the ANY_NOT_NULL aggregate, returning the
// first non-null value in the input column.
type anyNotNullDecimalHashAgg struct {
	hashAggregateFuncBase
	col                         coldata.Decimals
	curAgg                      apd.Decimal
	foundNonNullForCurrentGroup bool
}

var _ AggregateFunc = &anyNotNullDecimalHashAgg{}

func (a *anyNotNullDecimalHashAgg) SetOutput(vec coldata.Vec) {
	a.hashAggregateFuncBase.SetOutput(vec)
	a.col = vec.Decimal()
}

func (a *anyNotNullDecimalHashAgg) Compute(
	vecs []coldata.Vec, inputIdxs []uint32, inputLen int, sel []int,
) {
	if a.foundNonNullForCurrentGroup {
		// We have already seen non-null for the current group, and since there
		// is at most a single group when performing hash aggregation, we can
		// finish computing.
		return
	}

	oldCurAggSize := tree.SizeOfDecimal(&a.curAgg)
	vec := vecs[inputIdxs[0]]
	col, nulls := vec.Decimal(), vec.Nulls()
	a.allocator.PerformOperation([]coldata.Vec{a.vec}, func() {
		// Capture col to force bounds check to work. See
		// https://github.com/golang/go/issues/39756
		col := col
		_ = col.Get(inputLen - 1)
		{
			sel = sel[:inputLen]
			if nulls.MaybeHasNulls() {
				for _, i := range sel {

					var isNull bool
					isNull = nulls.NullAt(i)
					if !a.foundNonNullForCurrentGroup && !isNull {
						// If we haven't seen any non-nulls for the current group yet, and the
						// current value is non-null, then we can pick the current value to be
						// the output.
						val := col.Get(i)
						a.curAgg.Set(&val)
						a.foundNonNullForCurrentGroup = true
						// We have already seen non-null for the current group, and since there
						// is at most a single group when performing hash aggregation, we can
						// finish computing.
						return
					}
				}
			} else {
				for _, i := range sel {

					var isNull bool
					isNull = false
					if !a.foundNonNullForCurrentGroup && !isNull {
						// If we haven't seen any non-nulls for the current group yet, and the
						// current value is non-null, then we can pick the current value to be
						// the output.
						val := col.Get(i)
						a.curAgg.Set(&val)
						a.foundNonNullForCurrentGroup = true
						// We have already seen non-null for the current group, and since there
						// is at most a single group when performing hash aggregation, we can
						// finish computing.
						return
					}
				}
			}
		}
	},
	)
	newCurAggSize := tree.SizeOfDecimal(&a.curAgg)
	if newCurAggSize != oldCurAggSize {
		a.allocator.AdjustMemoryUsage(int64(newCurAggSize - oldCurAggSize))
	}
}

func (a *anyNotNullDecimalHashAgg) Flush(outputIdx int) {
	// If we haven't found any non-nulls for this group so far, the output for
	// this group should be null.
	if !a.foundNonNullForCurrentGroup {
		a.nulls.SetNull(outputIdx)
	} else {
		a.col[outputIdx].Set(&a.curAgg)
	}
}

type anyNotNullDecimalHashAggAlloc struct {
	aggAllocBase
	aggFuncs []anyNotNullDecimalHashAgg
}

var _ aggregateFuncAlloc = &anyNotNullDecimalHashAggAlloc{}

const sizeOfAnyNotNullDecimalHashAgg = int64(unsafe.Sizeof(anyNotNullDecimalHashAgg{}))
const anyNotNullDecimalHashAggSliceOverhead = int64(unsafe.Sizeof([]anyNotNullDecimalHashAgg{}))

func (a *anyNotNullDecimalHashAggAlloc) newAggFunc() AggregateFunc {
	if len(a.aggFuncs) == 0 {
		a.allocator.AdjustMemoryUsage(anyNotNullDecimalHashAggSliceOverhead + sizeOfAnyNotNullDecimalHashAgg*a.allocSize)
		a.aggFuncs = make([]anyNotNullDecimalHashAgg, a.allocSize)
	}
	f := &a.aggFuncs[0]
	f.allocator = a.allocator
	a.aggFuncs = a.aggFuncs[1:]
	return f
}

// anyNotNullInt16HashAgg implements the ANY_NOT_NULL aggregate, returning the
// first non-null value in the input column.
type anyNotNullInt16HashAgg struct {
	hashAggregateFuncBase
	col                         coldata.Int16s
	curAgg                      int16
	foundNonNullForCurrentGroup bool
}

var _ AggregateFunc = &anyNotNullInt16HashAgg{}

func (a *anyNotNullInt16HashAgg) SetOutput(vec coldata.Vec) {
	a.hashAggregateFuncBase.SetOutput(vec)
	a.col = vec.Int16()
}

func (a *anyNotNullInt16HashAgg) Compute(
	vecs []coldata.Vec, inputIdxs []uint32, inputLen int, sel []int,
) {
	if a.foundNonNullForCurrentGroup {
		// We have already seen non-null for the current group, and since there
		// is at most a single group when performing hash aggregation, we can
		// finish computing.
		return
	}

	var oldCurAggSize uintptr
	vec := vecs[inputIdxs[0]]
	col, nulls := vec.Int16(), vec.Nulls()
	a.allocator.PerformOperation([]coldata.Vec{a.vec}, func() {
		// Capture col to force bounds check to work. See
		// https://github.com/golang/go/issues/39756
		col := col
		_ = col.Get(inputLen - 1)
		{
			sel = sel[:inputLen]
			if nulls.MaybeHasNulls() {
				for _, i := range sel {

					var isNull bool
					isNull = nulls.NullAt(i)
					if !a.foundNonNullForCurrentGroup && !isNull {
						// If we haven't seen any non-nulls for the current group yet, and the
						// current value is non-null, then we can pick the current value to be
						// the output.
						val := col.Get(i)
						a.curAgg = val
						a.foundNonNullForCurrentGroup = true
						// We have already seen non-null for the current group, and since there
						// is at most a single group when performing hash aggregation, we can
						// finish computing.
						return
					}
				}
			} else {
				for _, i := range sel {

					var isNull bool
					isNull = false
					if !a.foundNonNullForCurrentGroup && !isNull {
						// If we haven't seen any non-nulls for the current group yet, and the
						// current value is non-null, then we can pick the current value to be
						// the output.
						val := col.Get(i)
						a.curAgg = val
						a.foundNonNullForCurrentGroup = true
						// We have already seen non-null for the current group, and since there
						// is at most a single group when performing hash aggregation, we can
						// finish computing.
						return
					}
				}
			}
		}
	},
	)
	var newCurAggSize uintptr
	if newCurAggSize != oldCurAggSize {
		a.allocator.AdjustMemoryUsage(int64(newCurAggSize - oldCurAggSize))
	}
}

func (a *anyNotNullInt16HashAgg) Flush(outputIdx int) {
	// If we haven't found any non-nulls for this group so far, the output for
	// this group should be null.
	if !a.foundNonNullForCurrentGroup {
		a.nulls.SetNull(outputIdx)
	} else {
		a.col[outputIdx] = a.curAgg
	}
}

type anyNotNullInt16HashAggAlloc struct {
	aggAllocBase
	aggFuncs []anyNotNullInt16HashAgg
}

var _ aggregateFuncAlloc = &anyNotNullInt16HashAggAlloc{}

const sizeOfAnyNotNullInt16HashAgg = int64(unsafe.Sizeof(anyNotNullInt16HashAgg{}))
const anyNotNullInt16HashAggSliceOverhead = int64(unsafe.Sizeof([]anyNotNullInt16HashAgg{}))

func (a *anyNotNullInt16HashAggAlloc) newAggFunc() AggregateFunc {
	if len(a.aggFuncs) == 0 {
		a.allocator.AdjustMemoryUsage(anyNotNullInt16HashAggSliceOverhead + sizeOfAnyNotNullInt16HashAgg*a.allocSize)
		a.aggFuncs = make([]anyNotNullInt16HashAgg, a.allocSize)
	}
	f := &a.aggFuncs[0]
	f.allocator = a.allocator
	a.aggFuncs = a.aggFuncs[1:]
	return f
}

// anyNotNullInt32HashAgg implements the ANY_NOT_NULL aggregate, returning the
// first non-null value in the input column.
type anyNotNullInt32HashAgg struct {
	hashAggregateFuncBase
	col                         coldata.Int32s
	curAgg                      int32
	foundNonNullForCurrentGroup bool
}

var _ AggregateFunc = &anyNotNullInt32HashAgg{}

func (a *anyNotNullInt32HashAgg) SetOutput(vec coldata.Vec) {
	a.hashAggregateFuncBase.SetOutput(vec)
	a.col = vec.Int32()
}

func (a *anyNotNullInt32HashAgg) Compute(
	vecs []coldata.Vec, inputIdxs []uint32, inputLen int, sel []int,
) {
	if a.foundNonNullForCurrentGroup {
		// We have already seen non-null for the current group, and since there
		// is at most a single group when performing hash aggregation, we can
		// finish computing.
		return
	}

	var oldCurAggSize uintptr
	vec := vecs[inputIdxs[0]]
	col, nulls := vec.Int32(), vec.Nulls()
	a.allocator.PerformOperation([]coldata.Vec{a.vec}, func() {
		// Capture col to force bounds check to work. See
		// https://github.com/golang/go/issues/39756
		col := col
		_ = col.Get(inputLen - 1)
		{
			sel = sel[:inputLen]
			if nulls.MaybeHasNulls() {
				for _, i := range sel {

					var isNull bool
					isNull = nulls.NullAt(i)
					if !a.foundNonNullForCurrentGroup && !isNull {
						// If we haven't seen any non-nulls for the current group yet, and the
						// current value is non-null, then we can pick the current value to be
						// the output.
						val := col.Get(i)
						a.curAgg = val
						a.foundNonNullForCurrentGroup = true
						// We have already seen non-null for the current group, and since there
						// is at most a single group when performing hash aggregation, we can
						// finish computing.
						return
					}
				}
			} else {
				for _, i := range sel {

					var isNull bool
					isNull = false
					if !a.foundNonNullForCurrentGroup && !isNull {
						// If we haven't seen any non-nulls for the current group yet, and the
						// current value is non-null, then we can pick the current value to be
						// the output.
						val := col.Get(i)
						a.curAgg = val
						a.foundNonNullForCurrentGroup = true
						// We have already seen non-null for the current group, and since there
						// is at most a single group when performing hash aggregation, we can
						// finish computing.
						return
					}
				}
			}
		}
	},
	)
	var newCurAggSize uintptr
	if newCurAggSize != oldCurAggSize {
		a.allocator.AdjustMemoryUsage(int64(newCurAggSize - oldCurAggSize))
	}
}

func (a *anyNotNullInt32HashAgg) Flush(outputIdx int) {
	// If we haven't found any non-nulls for this group so far, the output for
	// this group should be null.
	if !a.foundNonNullForCurrentGroup {
		a.nulls.SetNull(outputIdx)
	} else {
		a.col[outputIdx] = a.curAgg
	}
}

type anyNotNullInt32HashAggAlloc struct {
	aggAllocBase
	aggFuncs []anyNotNullInt32HashAgg
}

var _ aggregateFuncAlloc = &anyNotNullInt32HashAggAlloc{}

const sizeOfAnyNotNullInt32HashAgg = int64(unsafe.Sizeof(anyNotNullInt32HashAgg{}))
const anyNotNullInt32HashAggSliceOverhead = int64(unsafe.Sizeof([]anyNotNullInt32HashAgg{}))

func (a *anyNotNullInt32HashAggAlloc) newAggFunc() AggregateFunc {
	if len(a.aggFuncs) == 0 {
		a.allocator.AdjustMemoryUsage(anyNotNullInt32HashAggSliceOverhead + sizeOfAnyNotNullInt32HashAgg*a.allocSize)
		a.aggFuncs = make([]anyNotNullInt32HashAgg, a.allocSize)
	}
	f := &a.aggFuncs[0]
	f.allocator = a.allocator
	a.aggFuncs = a.aggFuncs[1:]
	return f
}

// anyNotNullInt64HashAgg implements the ANY_NOT_NULL aggregate, returning the
// first non-null value in the input column.
type anyNotNullInt64HashAgg struct {
	hashAggregateFuncBase
	col                         coldata.Int64s
	curAgg                      int64
	foundNonNullForCurrentGroup bool
}

var _ AggregateFunc = &anyNotNullInt64HashAgg{}

func (a *anyNotNullInt64HashAgg) SetOutput(vec coldata.Vec) {
	a.hashAggregateFuncBase.SetOutput(vec)
	a.col = vec.Int64()
}

func (a *anyNotNullInt64HashAgg) Compute(
	vecs []coldata.Vec, inputIdxs []uint32, inputLen int, sel []int,
) {
	if a.foundNonNullForCurrentGroup {
		// We have already seen non-null for the current group, and since there
		// is at most a single group when performing hash aggregation, we can
		// finish computing.
		return
	}

	var oldCurAggSize uintptr
	vec := vecs[inputIdxs[0]]
	col, nulls := vec.Int64(), vec.Nulls()
	a.allocator.PerformOperation([]coldata.Vec{a.vec}, func() {
		// Capture col to force bounds check to work. See
		// https://github.com/golang/go/issues/39756
		col := col
		_ = col.Get(inputLen - 1)
		{
			sel = sel[:inputLen]
			if nulls.MaybeHasNulls() {
				for _, i := range sel {

					var isNull bool
					isNull = nulls.NullAt(i)
					if !a.foundNonNullForCurrentGroup && !isNull {
						// If we haven't seen any non-nulls for the current group yet, and the
						// current value is non-null, then we can pick the current value to be
						// the output.
						val := col.Get(i)
						a.curAgg = val
						a.foundNonNullForCurrentGroup = true
						// We have already seen non-null for the current group, and since there
						// is at most a single group when performing hash aggregation, we can
						// finish computing.
						return
					}
				}
			} else {
				for _, i := range sel {

					var isNull bool
					isNull = false
					if !a.foundNonNullForCurrentGroup && !isNull {
						// If we haven't seen any non-nulls for the current group yet, and the
						// current value is non-null, then we can pick the current value to be
						// the output.
						val := col.Get(i)
						a.curAgg = val
						a.foundNonNullForCurrentGroup = true
						// We have already seen non-null for the current group, and since there
						// is at most a single group when performing hash aggregation, we can
						// finish computing.
						return
					}
				}
			}
		}
	},
	)
	var newCurAggSize uintptr
	if newCurAggSize != oldCurAggSize {
		a.allocator.AdjustMemoryUsage(int64(newCurAggSize - oldCurAggSize))
	}
}

func (a *anyNotNullInt64HashAgg) Flush(outputIdx int) {
	// If we haven't found any non-nulls for this group so far, the output for
	// this group should be null.
	if !a.foundNonNullForCurrentGroup {
		a.nulls.SetNull(outputIdx)
	} else {
		a.col[outputIdx] = a.curAgg
	}
}

type anyNotNullInt64HashAggAlloc struct {
	aggAllocBase
	aggFuncs []anyNotNullInt64HashAgg
}

var _ aggregateFuncAlloc = &anyNotNullInt64HashAggAlloc{}

const sizeOfAnyNotNullInt64HashAgg = int64(unsafe.Sizeof(anyNotNullInt64HashAgg{}))
const anyNotNullInt64HashAggSliceOverhead = int64(unsafe.Sizeof([]anyNotNullInt64HashAgg{}))

func (a *anyNotNullInt64HashAggAlloc) newAggFunc() AggregateFunc {
	if len(a.aggFuncs) == 0 {
		a.allocator.AdjustMemoryUsage(anyNotNullInt64HashAggSliceOverhead + sizeOfAnyNotNullInt64HashAgg*a.allocSize)
		a.aggFuncs = make([]anyNotNullInt64HashAgg, a.allocSize)
	}
	f := &a.aggFuncs[0]
	f.allocator = a.allocator
	a.aggFuncs = a.aggFuncs[1:]
	return f
}

// anyNotNullFloat64HashAgg implements the ANY_NOT_NULL aggregate, returning the
// first non-null value in the input column.
type anyNotNullFloat64HashAgg struct {
	hashAggregateFuncBase
	col                         coldata.Float64s
	curAgg                      float64
	foundNonNullForCurrentGroup bool
}

var _ AggregateFunc = &anyNotNullFloat64HashAgg{}

func (a *anyNotNullFloat64HashAgg) SetOutput(vec coldata.Vec) {
	a.hashAggregateFuncBase.SetOutput(vec)
	a.col = vec.Float64()
}

func (a *anyNotNullFloat64HashAgg) Compute(
	vecs []coldata.Vec, inputIdxs []uint32, inputLen int, sel []int,
) {
	if a.foundNonNullForCurrentGroup {
		// We have already seen non-null for the current group, and since there
		// is at most a single group when performing hash aggregation, we can
		// finish computing.
		return
	}

	var oldCurAggSize uintptr
	vec := vecs[inputIdxs[0]]
	col, nulls := vec.Float64(), vec.Nulls()
	a.allocator.PerformOperation([]coldata.Vec{a.vec}, func() {
		// Capture col to force bounds check to work. See
		// https://github.com/golang/go/issues/39756
		col := col
		_ = col.Get(inputLen - 1)
		{
			sel = sel[:inputLen]
			if nulls.MaybeHasNulls() {
				for _, i := range sel {

					var isNull bool
					isNull = nulls.NullAt(i)
					if !a.foundNonNullForCurrentGroup && !isNull {
						// If we haven't seen any non-nulls for the current group yet, and the
						// current value is non-null, then we can pick the current value to be
						// the output.
						val := col.Get(i)
						a.curAgg = val
						a.foundNonNullForCurrentGroup = true
						// We have already seen non-null for the current group, and since there
						// is at most a single group when performing hash aggregation, we can
						// finish computing.
						return
					}
				}
			} else {
				for _, i := range sel {

					var isNull bool
					isNull = false
					if !a.foundNonNullForCurrentGroup && !isNull {
						// If we haven't seen any non-nulls for the current group yet, and the
						// current value is non-null, then we can pick the current value to be
						// the output.
						val := col.Get(i)
						a.curAgg = val
						a.foundNonNullForCurrentGroup = true
						// We have already seen non-null for the current group, and since there
						// is at most a single group when performing hash aggregation, we can
						// finish computing.
						return
					}
				}
			}
		}
	},
	)
	var newCurAggSize uintptr
	if newCurAggSize != oldCurAggSize {
		a.allocator.AdjustMemoryUsage(int64(newCurAggSize - oldCurAggSize))
	}
}

func (a *anyNotNullFloat64HashAgg) Flush(outputIdx int) {
	// If we haven't found any non-nulls for this group so far, the output for
	// this group should be null.
	if !a.foundNonNullForCurrentGroup {
		a.nulls.SetNull(outputIdx)
	} else {
		a.col[outputIdx] = a.curAgg
	}
}

type anyNotNullFloat64HashAggAlloc struct {
	aggAllocBase
	aggFuncs []anyNotNullFloat64HashAgg
}

var _ aggregateFuncAlloc = &anyNotNullFloat64HashAggAlloc{}

const sizeOfAnyNotNullFloat64HashAgg = int64(unsafe.Sizeof(anyNotNullFloat64HashAgg{}))
const anyNotNullFloat64HashAggSliceOverhead = int64(unsafe.Sizeof([]anyNotNullFloat64HashAgg{}))

func (a *anyNotNullFloat64HashAggAlloc) newAggFunc() AggregateFunc {
	if len(a.aggFuncs) == 0 {
		a.allocator.AdjustMemoryUsage(anyNotNullFloat64HashAggSliceOverhead + sizeOfAnyNotNullFloat64HashAgg*a.allocSize)
		a.aggFuncs = make([]anyNotNullFloat64HashAgg, a.allocSize)
	}
	f := &a.aggFuncs[0]
	f.allocator = a.allocator
	a.aggFuncs = a.aggFuncs[1:]
	return f
}

// anyNotNullTimestampHashAgg implements the ANY_NOT_NULL aggregate, returning the
// first non-null value in the input column.
type anyNotNullTimestampHashAgg struct {
	hashAggregateFuncBase
	col                         coldata.Times
	curAgg                      time.Time
	foundNonNullForCurrentGroup bool
}

var _ AggregateFunc = &anyNotNullTimestampHashAgg{}

func (a *anyNotNullTimestampHashAgg) SetOutput(vec coldata.Vec) {
	a.hashAggregateFuncBase.SetOutput(vec)
	a.col = vec.Timestamp()
}

func (a *anyNotNullTimestampHashAgg) Compute(
	vecs []coldata.Vec, inputIdxs []uint32, inputLen int, sel []int,
) {
	if a.foundNonNullForCurrentGroup {
		// We have already seen non-null for the current group, and since there
		// is at most a single group when performing hash aggregation, we can
		// finish computing.
		return
	}

	var oldCurAggSize uintptr
	vec := vecs[inputIdxs[0]]
	col, nulls := vec.Timestamp(), vec.Nulls()
	a.allocator.PerformOperation([]coldata.Vec{a.vec}, func() {
		// Capture col to force bounds check to work. See
		// https://github.com/golang/go/issues/39756
		col := col
		_ = col.Get(inputLen - 1)
		{
			sel = sel[:inputLen]
			if nulls.MaybeHasNulls() {
				for _, i := range sel {

					var isNull bool
					isNull = nulls.NullAt(i)
					if !a.foundNonNullForCurrentGroup && !isNull {
						// If we haven't seen any non-nulls for the current group yet, and the
						// current value is non-null, then we can pick the current value to be
						// the output.
						val := col.Get(i)
						a.curAgg = val
						a.foundNonNullForCurrentGroup = true
						// We have already seen non-null for the current group, and since there
						// is at most a single group when performing hash aggregation, we can
						// finish computing.
						return
					}
				}
			} else {
				for _, i := range sel {

					var isNull bool
					isNull = false
					if !a.foundNonNullForCurrentGroup && !isNull {
						// If we haven't seen any non-nulls for the current group yet, and the
						// current value is non-null, then we can pick the current value to be
						// the output.
						val := col.Get(i)
						a.curAgg = val
						a.foundNonNullForCurrentGroup = true
						// We have already seen non-null for the current group, and since there
						// is at most a single group when performing hash aggregation, we can
						// finish computing.
						return
					}
				}
			}
		}
	},
	)
	var newCurAggSize uintptr
	if newCurAggSize != oldCurAggSize {
		a.allocator.AdjustMemoryUsage(int64(newCurAggSize - oldCurAggSize))
	}
}

func (a *anyNotNullTimestampHashAgg) Flush(outputIdx int) {
	// If we haven't found any non-nulls for this group so far, the output for
	// this group should be null.
	if !a.foundNonNullForCurrentGroup {
		a.nulls.SetNull(outputIdx)
	} else {
		a.col[outputIdx] = a.curAgg
	}
}

type anyNotNullTimestampHashAggAlloc struct {
	aggAllocBase
	aggFuncs []anyNotNullTimestampHashAgg
}

var _ aggregateFuncAlloc = &anyNotNullTimestampHashAggAlloc{}

const sizeOfAnyNotNullTimestampHashAgg = int64(unsafe.Sizeof(anyNotNullTimestampHashAgg{}))
const anyNotNullTimestampHashAggSliceOverhead = int64(unsafe.Sizeof([]anyNotNullTimestampHashAgg{}))

func (a *anyNotNullTimestampHashAggAlloc) newAggFunc() AggregateFunc {
	if len(a.aggFuncs) == 0 {
		a.allocator.AdjustMemoryUsage(anyNotNullTimestampHashAggSliceOverhead + sizeOfAnyNotNullTimestampHashAgg*a.allocSize)
		a.aggFuncs = make([]anyNotNullTimestampHashAgg, a.allocSize)
	}
	f := &a.aggFuncs[0]
	f.allocator = a.allocator
	a.aggFuncs = a.aggFuncs[1:]
	return f
}

// anyNotNullIntervalHashAgg implements the ANY_NOT_NULL aggregate, returning the
// first non-null value in the input column.
type anyNotNullIntervalHashAgg struct {
	hashAggregateFuncBase
	col                         coldata.Durations
	curAgg                      duration.Duration
	foundNonNullForCurrentGroup bool
}

var _ AggregateFunc = &anyNotNullIntervalHashAgg{}

func (a *anyNotNullIntervalHashAgg) SetOutput(vec coldata.Vec) {
	a.hashAggregateFuncBase.SetOutput(vec)
	a.col = vec.Interval()
}

func (a *anyNotNullIntervalHashAgg) Compute(
	vecs []coldata.Vec, inputIdxs []uint32, inputLen int, sel []int,
) {
	if a.foundNonNullForCurrentGroup {
		// We have already seen non-null for the current group, and since there
		// is at most a single group when performing hash aggregation, we can
		// finish computing.
		return
	}

	var oldCurAggSize uintptr
	vec := vecs[inputIdxs[0]]
	col, nulls := vec.Interval(), vec.Nulls()
	a.allocator.PerformOperation([]coldata.Vec{a.vec}, func() {
		// Capture col to force bounds check to work. See
		// https://github.com/golang/go/issues/39756
		col := col
		_ = col.Get(inputLen - 1)
		{
			sel = sel[:inputLen]
			if nulls.MaybeHasNulls() {
				for _, i := range sel {

					var isNull bool
					isNull = nulls.NullAt(i)
					if !a.foundNonNullForCurrentGroup && !isNull {
						// If we haven't seen any non-nulls for the current group yet, and the
						// current value is non-null, then we can pick the current value to be
						// the output.
						val := col.Get(i)
						a.curAgg = val
						a.foundNonNullForCurrentGroup = true
						// We have already seen non-null for the current group, and since there
						// is at most a single group when performing hash aggregation, we can
						// finish computing.
						return
					}
				}
			} else {
				for _, i := range sel {

					var isNull bool
					isNull = false
					if !a.foundNonNullForCurrentGroup && !isNull {
						// If we haven't seen any non-nulls for the current group yet, and the
						// current value is non-null, then we can pick the current value to be
						// the output.
						val := col.Get(i)
						a.curAgg = val
						a.foundNonNullForCurrentGroup = true
						// We have already seen non-null for the current group, and since there
						// is at most a single group when performing hash aggregation, we can
						// finish computing.
						return
					}
				}
			}
		}
	},
	)
	var newCurAggSize uintptr
	if newCurAggSize != oldCurAggSize {
		a.allocator.AdjustMemoryUsage(int64(newCurAggSize - oldCurAggSize))
	}
}

func (a *anyNotNullIntervalHashAgg) Flush(outputIdx int) {
	// If we haven't found any non-nulls for this group so far, the output for
	// this group should be null.
	if !a.foundNonNullForCurrentGroup {
		a.nulls.SetNull(outputIdx)
	} else {
		a.col[outputIdx] = a.curAgg
	}
}

type anyNotNullIntervalHashAggAlloc struct {
	aggAllocBase
	aggFuncs []anyNotNullIntervalHashAgg
}

var _ aggregateFuncAlloc = &anyNotNullIntervalHashAggAlloc{}

const sizeOfAnyNotNullIntervalHashAgg = int64(unsafe.Sizeof(anyNotNullIntervalHashAgg{}))
const anyNotNullIntervalHashAggSliceOverhead = int64(unsafe.Sizeof([]anyNotNullIntervalHashAgg{}))

func (a *anyNotNullIntervalHashAggAlloc) newAggFunc() AggregateFunc {
	if len(a.aggFuncs) == 0 {
		a.allocator.AdjustMemoryUsage(anyNotNullIntervalHashAggSliceOverhead + sizeOfAnyNotNullIntervalHashAgg*a.allocSize)
		a.aggFuncs = make([]anyNotNullIntervalHashAgg, a.allocSize)
	}
	f := &a.aggFuncs[0]
	f.allocator = a.allocator
	a.aggFuncs = a.aggFuncs[1:]
	return f
}

// anyNotNullDatumHashAgg implements the ANY_NOT_NULL aggregate, returning the
// first non-null value in the input column.
type anyNotNullDatumHashAgg struct {
	hashAggregateFuncBase
	col                         coldata.DatumVec
	curAgg                      interface{}
	foundNonNullForCurrentGroup bool
}

var _ AggregateFunc = &anyNotNullDatumHashAgg{}

func (a *anyNotNullDatumHashAgg) SetOutput(vec coldata.Vec) {
	a.hashAggregateFuncBase.SetOutput(vec)
	a.col = vec.Datum()
}

func (a *anyNotNullDatumHashAgg) Compute(
	vecs []coldata.Vec, inputIdxs []uint32, inputLen int, sel []int,
) {
	if a.foundNonNullForCurrentGroup {
		// We have already seen non-null for the current group, and since there
		// is at most a single group when performing hash aggregation, we can
		// finish computing.
		return
	}

	var oldCurAggSize uintptr
	if a.curAgg != nil {
		oldCurAggSize = a.curAgg.(*coldataext.Datum).Size()
	}
	vec := vecs[inputIdxs[0]]
	col, nulls := vec.Datum(), vec.Nulls()
	a.allocator.PerformOperation([]coldata.Vec{a.vec}, func() {
		// Capture col to force bounds check to work. See
		// https://github.com/golang/go/issues/39756
		col := col
		_ = col.Get(inputLen - 1)
		{
			sel = sel[:inputLen]
			if nulls.MaybeHasNulls() {
				for _, i := range sel {

					var isNull bool
					isNull = nulls.NullAt(i)
					if !a.foundNonNullForCurrentGroup && !isNull {
						// If we haven't seen any non-nulls for the current group yet, and the
						// current value is non-null, then we can pick the current value to be
						// the output.
						val := col.Get(i)
						a.curAgg = val
						a.foundNonNullForCurrentGroup = true
						// We have already seen non-null for the current group, and since there
						// is at most a single group when performing hash aggregation, we can
						// finish computing.
						return
					}
				}
			} else {
				for _, i := range sel {

					var isNull bool
					isNull = false
					if !a.foundNonNullForCurrentGroup && !isNull {
						// If we haven't seen any non-nulls for the current group yet, and the
						// current value is non-null, then we can pick the current value to be
						// the output.
						val := col.Get(i)
						a.curAgg = val
						a.foundNonNullForCurrentGroup = true
						// We have already seen non-null for the current group, and since there
						// is at most a single group when performing hash aggregation, we can
						// finish computing.
						return
					}
				}
			}
		}
	},
	)

	var newCurAggSize uintptr
	if a.curAgg != nil {
		newCurAggSize = a.curAgg.(*coldataext.Datum).Size()
	}
	if newCurAggSize != oldCurAggSize {
		a.allocator.AdjustMemoryUsage(int64(newCurAggSize - oldCurAggSize))
	}
}

func (a *anyNotNullDatumHashAgg) Flush(outputIdx int) {
	// If we haven't found any non-nulls for this group so far, the output for
	// this group should be null.
	if !a.foundNonNullForCurrentGroup {
		a.nulls.SetNull(outputIdx)
	} else {
		a.col.Set(outputIdx, a.curAgg)
	}
	// Release the reference to curAgg eagerly.
	if d, ok := a.curAgg.(*coldataext.Datum); ok {
		a.allocator.AdjustMemoryUsage(-int64(d.Size()))
	}
	a.curAgg = nil
}

type anyNotNullDatumHashAggAlloc struct {
	aggAllocBase
	aggFuncs []anyNotNullDatumHashAgg
}

var _ aggregateFuncAlloc = &anyNotNullDatumHashAggAlloc{}

const sizeOfAnyNotNullDatumHashAgg = int64(unsafe.Sizeof(anyNotNullDatumHashAgg{}))
const anyNotNullDatumHashAggSliceOverhead = int64(unsafe.Sizeof([]anyNotNullDatumHashAgg{}))

func (a *anyNotNullDatumHashAggAlloc) newAggFunc() AggregateFunc {
	if len(a.aggFuncs) == 0 {
		a.allocator.AdjustMemoryUsage(anyNotNullDatumHashAggSliceOverhead + sizeOfAnyNotNullDatumHashAgg*a.allocSize)
		a.aggFuncs = make([]anyNotNullDatumHashAgg, a.allocSize)
	}
	f := &a.aggFuncs[0]
	f.allocator = a.allocator
	a.aggFuncs = a.aggFuncs[1:]
	return f
}
