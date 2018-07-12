// Code generated by optgen; DO NOT EDIT.

package opt

const (
	startAutoRule RuleName = iota + NumManualRuleNames

	// ------------------------------------------------------------
	// Normalize Rule Names
	// ------------------------------------------------------------
	EliminateEmptyAnd
	EliminateEmptyOr
	EliminateSingletonAndOr
	SimplifyAnd
	SimplifyOr
	SimplifyFilters
	FoldNullAndOr
	NegateComparison
	EliminateNot
	NegateAnd
	NegateOr
	ExtractRedundantClause
	ExtractRedundantSubclause
	CommuteVarInequality
	CommuteConstInequality
	NormalizeCmpPlusConst
	NormalizeCmpMinusConst
	NormalizeCmpConstMinus
	NormalizeTupleEquality
	FoldNullComparisonLeft
	FoldNullComparisonRight
	FoldIsNull
	FoldNonNullIsNull
	FoldIsNotNull
	FoldNonNullIsNotNull
	CommuteNullIs
	DecorrelateJoin
	TryDecorrelateSelect
	TryDecorrelateProject
	TryDecorrelateProjectSelect
	TryDecorrelateInnerJoin
	TryDecorrelateGroupBy
	TryDecorrelateScalarGroupBy
	HoistSelectExists
	HoistSelectNotExists
	HoistSelectSubquery
	HoistProjectSubquery
	HoistJoinSubquery
	HoistValuesSubquery
	NormalizeAnyFilter
	NormalizeNotAnyFilter
	EliminateDistinct
	EliminateGroupByProject
	PushSelectIntoInlinableProject
	InlineProjectInProject
	EnsureJoinFiltersAnd
	EnsureJoinFilters
	PushFilterIntoJoinLeftAndRight
	MapFilterIntoJoinLeft
	MapFilterIntoJoinRight
	PushFilterIntoJoinLeft
	PushFilterIntoJoinRight
	SimplifyLeftJoin
	SimplifyRightJoin
	EliminateSemiJoin
	EliminateAntiJoin
	EliminateJoinNoColsLeft
	EliminateJoinNoColsRight
	EliminateLimit
	PushLimitIntoProject
	PushOffsetIntoProject
	EliminateMax1Row
	FoldPlusZero
	FoldZeroPlus
	FoldMinusZero
	FoldMultOne
	FoldOneMult
	FoldDivOne
	InvertMinus
	EliminateUnaryMinus
	FoldUnaryMinus
	SimplifyLimitOrdering
	SimplifyOffsetOrdering
	SimplifyGroupByOrdering
	SimplifyRowNumberOrdering
	SimplifyExplainOrdering
	EliminateProject
	EliminateProjectProject
	PruneProjectCols
	PruneScanCols
	PruneSelectCols
	PruneLimitCols
	PruneOffsetCols
	PruneJoinLeftCols
	PruneJoinRightCols
	PruneAggCols
	PruneGroupByCols
	PruneValuesCols
	PruneRowNumberCols
	PruneExplainCols
	CommuteVar
	CommuteConst
	EliminateCoalesce
	SimplifyCoalesce
	EliminateCast
	FoldNullCast
	FoldNullUnary
	FoldNullBinaryLeft
	FoldNullBinaryRight
	FoldNullInNonEmpty
	FoldNullInEmpty
	FoldNullNotInEmpty
	NormalizeInConst
	FoldInNull
	EliminateExistsProject
	EliminateExistsGroupBy
	EliminateSelect
	EnsureSelectFiltersAnd
	EnsureSelectFilters
	MergeSelects
	PushSelectIntoProject
	SimplifySelectLeftJoin
	SimplifySelectRightJoin
	MergeSelectInnerJoin
	PushSelectIntoJoinLeft
	PushSelectIntoJoinRight
	PushSelectIntoGroupBy
	RemoveNotNullCondition

	// startExploreRule tracks the number of normalization rules;
	// all rules greater than this value are exploration rules.
	startExploreRule

	// ------------------------------------------------------------
	// Explore Rule Names
	// ------------------------------------------------------------
	ReplaceMinWithLimit
	ReplaceMaxWithLimit
	CommuteJoin
	CommuteLeftJoin
	CommuteRightJoin
	GenerateMergeJoins
	GenerateLookupJoin
	GenerateLookupJoinWithFilter
	PushLimitIntoScan
	PushLimitIntoIndexJoin
	GenerateIndexScans
	ConstrainScan
	PushFilterIntoIndexJoinNoRemainder
	PushFilterIntoIndexJoin
	ConstrainIndexJoinScan

	// NumRuleNames tracks the total count of rule names.
	NumRuleNames
)
