// Code generated by "stringer -output=pkg/sql/opt/rule_name_string.go -type=RuleName pkg/sql/opt/rule_name.go pkg/sql/opt/rule_name.og.go"; DO NOT EDIT.

package opt

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[InvalidRuleName-0]
	_ = x[SimplifyRootOrdering-1]
	_ = x[PruneRootCols-2]
	_ = x[SimplifyZeroCardinalityGroup-3]
	_ = x[NumManualRuleNames-4]
	_ = x[startAutoRule-4]
	_ = x[EliminateAggDistinct-5]
	_ = x[NormalizeNestedAnds-6]
	_ = x[SimplifyTrueAnd-7]
	_ = x[SimplifyAndTrue-8]
	_ = x[SimplifyFalseAnd-9]
	_ = x[SimplifyAndFalse-10]
	_ = x[SimplifyTrueOr-11]
	_ = x[SimplifyOrTrue-12]
	_ = x[SimplifyFalseOr-13]
	_ = x[SimplifyOrFalse-14]
	_ = x[SimplifyRange-15]
	_ = x[FoldNullAndOr-16]
	_ = x[FoldNotTrue-17]
	_ = x[FoldNotFalse-18]
	_ = x[FoldNotNull-19]
	_ = x[NegateComparison-20]
	_ = x[EliminateNot-21]
	_ = x[NegateAnd-22]
	_ = x[NegateOr-23]
	_ = x[ExtractRedundantConjunct-24]
	_ = x[CommuteVarInequality-25]
	_ = x[CommuteConstInequality-26]
	_ = x[NormalizeCmpPlusConst-27]
	_ = x[NormalizeCmpMinusConst-28]
	_ = x[NormalizeCmpConstMinus-29]
	_ = x[NormalizeTupleEquality-30]
	_ = x[FoldNullComparisonLeft-31]
	_ = x[FoldNullComparisonRight-32]
	_ = x[FoldIsNull-33]
	_ = x[FoldNonNullIsNull-34]
	_ = x[FoldNullTupleIsTupleNull-35]
	_ = x[FoldNonNullTupleIsTupleNull-36]
	_ = x[FoldIsNotNull-37]
	_ = x[FoldNonNullIsNotNull-38]
	_ = x[FoldNonNullTupleIsTupleNotNull-39]
	_ = x[FoldNullTupleIsTupleNotNull-40]
	_ = x[CommuteNullIs-41]
	_ = x[NormalizeCmpTimeZoneFunction-42]
	_ = x[NormalizeCmpTimeZoneFunctionTZ-43]
	_ = x[DecorrelateJoin-44]
	_ = x[DecorrelateProjectSet-45]
	_ = x[TryDecorrelateSelect-46]
	_ = x[TryDecorrelateProject-47]
	_ = x[TryDecorrelateProjectSelect-48]
	_ = x[TryDecorrelateProjectInnerJoin-49]
	_ = x[TryDecorrelateInnerJoin-50]
	_ = x[TryDecorrelateInnerLeftJoin-51]
	_ = x[TryDecorrelateGroupBy-52]
	_ = x[TryDecorrelateScalarGroupBy-53]
	_ = x[TryDecorrelateSemiJoin-54]
	_ = x[TryDecorrelateLimitOne-55]
	_ = x[TryDecorrelateProjectSet-56]
	_ = x[TryDecorrelateWindow-57]
	_ = x[TryDecorrelateMax1Row-58]
	_ = x[HoistSelectExists-59]
	_ = x[HoistSelectNotExists-60]
	_ = x[HoistSelectSubquery-61]
	_ = x[HoistProjectSubquery-62]
	_ = x[HoistJoinSubquery-63]
	_ = x[HoistValuesSubquery-64]
	_ = x[HoistProjectSetSubquery-65]
	_ = x[NormalizeSelectAnyFilter-66]
	_ = x[NormalizeJoinAnyFilter-67]
	_ = x[NormalizeSelectNotAnyFilter-68]
	_ = x[NormalizeJoinNotAnyFilter-69]
	_ = x[FoldNullCast-70]
	_ = x[FoldNullUnary-71]
	_ = x[FoldNullBinaryLeft-72]
	_ = x[FoldNullBinaryRight-73]
	_ = x[FoldNullInNonEmpty-74]
	_ = x[FoldInEmpty-75]
	_ = x[FoldNotInEmpty-76]
	_ = x[FoldArray-77]
	_ = x[FoldBinary-78]
	_ = x[FoldUnary-79]
	_ = x[FoldComparison-80]
	_ = x[FoldCast-81]
	_ = x[FoldIndirection-82]
	_ = x[FoldColumnAccess-83]
	_ = x[FoldFunction-84]
	_ = x[FoldEqualsAnyNull-85]
	_ = x[ConvertGroupByToDistinct-86]
	_ = x[EliminateGroupByProject-87]
	_ = x[EliminateJoinUnderGroupByLeft-88]
	_ = x[EliminateJoinUnderGroupByRight-89]
	_ = x[EliminateDistinct-90]
	_ = x[ReduceGroupingCols-91]
	_ = x[ReduceNotNullGroupingCols-92]
	_ = x[EliminateAggDistinctForKeys-93]
	_ = x[EliminateAggFilteredDistinctForKeys-94]
	_ = x[EliminateDistinctNoColumns-95]
	_ = x[EliminateEnsureDistinctNoColumns-96]
	_ = x[EliminateDistinctOnValues-97]
	_ = x[PushAggDistinctIntoGroupBy-98]
	_ = x[PushAggFilterIntoScalarGroupBy-99]
	_ = x[ConvertCountToCountRows-100]
	_ = x[FoldGroupingOperators-101]
	_ = x[InlineConstVar-102]
	_ = x[InlineProjectConstants-103]
	_ = x[InlineSelectConstants-104]
	_ = x[InlineJoinConstantsLeft-105]
	_ = x[InlineJoinConstantsRight-106]
	_ = x[PushSelectIntoInlinableProject-107]
	_ = x[InlineProjectInProject-108]
	_ = x[CommuteRightJoin-109]
	_ = x[SimplifyJoinFilters-110]
	_ = x[DetectJoinContradiction-111]
	_ = x[PushFilterIntoJoinLeftAndRight-112]
	_ = x[MapFilterIntoJoinLeft-113]
	_ = x[MapFilterIntoJoinRight-114]
	_ = x[MapEqualityIntoJoinLeftAndRight-115]
	_ = x[PushFilterIntoJoinLeft-116]
	_ = x[PushFilterIntoJoinRight-117]
	_ = x[SimplifyLeftJoin-118]
	_ = x[SimplifyRightJoin-119]
	_ = x[EliminateSemiJoin-120]
	_ = x[SimplifyZeroCardinalitySemiJoin-121]
	_ = x[EliminateAntiJoin-122]
	_ = x[SimplifyZeroCardinalityAntiJoin-123]
	_ = x[EliminateJoinNoColsLeft-124]
	_ = x[EliminateJoinNoColsRight-125]
	_ = x[HoistJoinProjectRight-126]
	_ = x[HoistJoinProjectLeft-127]
	_ = x[SimplifyJoinNotNullEquality-128]
	_ = x[ExtractJoinEqualities-129]
	_ = x[SortFiltersInJoin-130]
	_ = x[LeftAssociateJoinsLeft-131]
	_ = x[LeftAssociateJoinsRight-132]
	_ = x[RightAssociateJoinsLeft-133]
	_ = x[RightAssociateJoinsRight-134]
	_ = x[EliminateLimit-135]
	_ = x[EliminateOffset-136]
	_ = x[PushLimitIntoProject-137]
	_ = x[PushOffsetIntoProject-138]
	_ = x[PushLimitIntoOffset-139]
	_ = x[PushLimitIntoOrdinality-140]
	_ = x[PushLimitIntoJoinLeft-141]
	_ = x[PushLimitIntoJoinRight-142]
	_ = x[FoldLimits-143]
	_ = x[AssociateLimitJoinsLeft-144]
	_ = x[AssociateLimitJoinsRight-145]
	_ = x[EliminateMax1Row-146]
	_ = x[FoldPlusZero-147]
	_ = x[FoldZeroPlus-148]
	_ = x[FoldMinusZero-149]
	_ = x[FoldMultOne-150]
	_ = x[FoldOneMult-151]
	_ = x[FoldDivOne-152]
	_ = x[InvertMinus-153]
	_ = x[EliminateUnaryMinus-154]
	_ = x[SimplifyLimitOrdering-155]
	_ = x[SimplifyOffsetOrdering-156]
	_ = x[SimplifyGroupByOrdering-157]
	_ = x[SimplifyOrdinalityOrdering-158]
	_ = x[SimplifyExplainOrdering-159]
	_ = x[EliminateJoinUnderProjectLeft-160]
	_ = x[EliminateJoinUnderProjectRight-161]
	_ = x[EliminateProject-162]
	_ = x[MergeProjects-163]
	_ = x[MergeProjectWithValues-164]
	_ = x[PushColumnRemappingIntoValues-165]
	_ = x[FoldTupleAccessIntoValues-166]
	_ = x[FoldJSONAccessIntoValues-167]
	_ = x[ConvertZipArraysToValues-168]
	_ = x[PruneProjectCols-169]
	_ = x[PruneScanCols-170]
	_ = x[PruneSelectCols-171]
	_ = x[PruneLimitCols-172]
	_ = x[PruneOffsetCols-173]
	_ = x[PruneJoinLeftCols-174]
	_ = x[PruneJoinRightCols-175]
	_ = x[PruneSemiAntiJoinRightCols-176]
	_ = x[PruneAggCols-177]
	_ = x[PruneGroupByCols-178]
	_ = x[PruneValuesCols-179]
	_ = x[PruneOrdinalityCols-180]
	_ = x[PruneExplainCols-181]
	_ = x[PruneProjectSetCols-182]
	_ = x[PruneWindowOutputCols-183]
	_ = x[PruneWindowInputCols-184]
	_ = x[PruneMutationFetchCols-185]
	_ = x[PruneMutationInputCols-186]
	_ = x[PruneMutationReturnCols-187]
	_ = x[PruneWithScanCols-188]
	_ = x[PruneWithCols-189]
	_ = x[PruneUnionAllCols-190]
	_ = x[RejectNullsLeftJoin-191]
	_ = x[RejectNullsRightJoin-192]
	_ = x[RejectNullsGroupBy-193]
	_ = x[CommuteVar-194]
	_ = x[CommuteConst-195]
	_ = x[EliminateCoalesce-196]
	_ = x[SimplifyCoalesce-197]
	_ = x[EliminateCast-198]
	_ = x[NormalizeInConst-199]
	_ = x[FoldInNull-200]
	_ = x[UnifyComparisonTypes-201]
	_ = x[EliminateExistsZeroRows-202]
	_ = x[EliminateExistsProject-203]
	_ = x[EliminateExistsGroupBy-204]
	_ = x[IntroduceExistsLimit-205]
	_ = x[EliminateExistsLimit-206]
	_ = x[NormalizeJSONFieldAccess-207]
	_ = x[NormalizeJSONContains-208]
	_ = x[SimplifyCaseWhenConstValue-209]
	_ = x[InlineAnyValuesSingleCol-210]
	_ = x[InlineAnyValuesMultiCol-211]
	_ = x[SimplifyEqualsAnyTuple-212]
	_ = x[SimplifyAnyScalarArray-213]
	_ = x[FoldCollate-214]
	_ = x[NormalizeArrayFlattenToAgg-215]
	_ = x[SimplifySameVarEqualities-216]
	_ = x[SimplifySameVarInequalities-217]
	_ = x[SimplifySelectFilters-218]
	_ = x[ConsolidateSelectFilters-219]
	_ = x[DetectSelectContradiction-220]
	_ = x[EliminateSelect-221]
	_ = x[MergeSelects-222]
	_ = x[PushSelectIntoProject-223]
	_ = x[MergeSelectInnerJoin-224]
	_ = x[PushSelectCondLeftIntoJoinLeftAndRight-225]
	_ = x[PushSelectIntoJoinLeft-226]
	_ = x[PushSelectIntoGroupBy-227]
	_ = x[RemoveNotNullCondition-228]
	_ = x[PushSelectIntoProjectSet-229]
	_ = x[PushFilterIntoSetOp-230]
	_ = x[EliminateUnionAllLeft-231]
	_ = x[EliminateUnionAllRight-232]
	_ = x[EliminateWindow-233]
	_ = x[ReduceWindowPartitionCols-234]
	_ = x[SimplifyWindowOrdering-235]
	_ = x[PushSelectIntoWindow-236]
	_ = x[PushLimitIntoWindow-237]
	_ = x[InlineWith-238]
	_ = x[startExploreRule-239]
	_ = x[ReplaceScalarMinMaxWithLimit-240]
	_ = x[ReplaceMinWithLimit-241]
	_ = x[ReplaceMaxWithLimit-242]
	_ = x[GenerateStreamingGroupBy-243]
	_ = x[CommuteJoin-244]
	_ = x[CommuteLeftJoin-245]
	_ = x[CommuteSemiJoin-246]
	_ = x[GenerateMergeJoins-247]
	_ = x[GenerateLookupJoins-248]
	_ = x[GenerateGeoLookupJoins-249]
	_ = x[GenerateZigzagJoins-250]
	_ = x[GenerateInvertedIndexZigzagJoins-251]
	_ = x[GenerateLookupJoinsWithFilter-252]
	_ = x[AssociateJoin-253]
	_ = x[GenerateLimitedScans-254]
	_ = x[PushLimitIntoConstrainedScan-255]
	_ = x[PushLimitIntoIndexJoin-256]
	_ = x[SplitScanIntoUnionScans-257]
	_ = x[GenerateIndexScans-258]
	_ = x[GeneratePartialIndexScans-259]
	_ = x[GenerateConstrainedScans-260]
	_ = x[GenerateInvertedIndexScans-261]
	_ = x[SplitDisjunction-262]
	_ = x[SplitDisjunctionAddKey-263]
	_ = x[NumRuleNames-264]
}

const _RuleName_name = "InvalidRuleNameSimplifyRootOrderingPruneRootColsSimplifyZeroCardinalityGroupNumManualRuleNamesEliminateAggDistinctNormalizeNestedAndsSimplifyTrueAndSimplifyAndTrueSimplifyFalseAndSimplifyAndFalseSimplifyTrueOrSimplifyOrTrueSimplifyFalseOrSimplifyOrFalseSimplifyRangeFoldNullAndOrFoldNotTrueFoldNotFalseFoldNotNullNegateComparisonEliminateNotNegateAndNegateOrExtractRedundantConjunctCommuteVarInequalityCommuteConstInequalityNormalizeCmpPlusConstNormalizeCmpMinusConstNormalizeCmpConstMinusNormalizeTupleEqualityFoldNullComparisonLeftFoldNullComparisonRightFoldIsNullFoldNonNullIsNullFoldNullTupleIsTupleNullFoldNonNullTupleIsTupleNullFoldIsNotNullFoldNonNullIsNotNullFoldNonNullTupleIsTupleNotNullFoldNullTupleIsTupleNotNullCommuteNullIsNormalizeCmpTimeZoneFunctionNormalizeCmpTimeZoneFunctionTZDecorrelateJoinDecorrelateProjectSetTryDecorrelateSelectTryDecorrelateProjectTryDecorrelateProjectSelectTryDecorrelateProjectInnerJoinTryDecorrelateInnerJoinTryDecorrelateInnerLeftJoinTryDecorrelateGroupByTryDecorrelateScalarGroupByTryDecorrelateSemiJoinTryDecorrelateLimitOneTryDecorrelateProjectSetTryDecorrelateWindowTryDecorrelateMax1RowHoistSelectExistsHoistSelectNotExistsHoistSelectSubqueryHoistProjectSubqueryHoistJoinSubqueryHoistValuesSubqueryHoistProjectSetSubqueryNormalizeSelectAnyFilterNormalizeJoinAnyFilterNormalizeSelectNotAnyFilterNormalizeJoinNotAnyFilterFoldNullCastFoldNullUnaryFoldNullBinaryLeftFoldNullBinaryRightFoldNullInNonEmptyFoldInEmptyFoldNotInEmptyFoldArrayFoldBinaryFoldUnaryFoldComparisonFoldCastFoldIndirectionFoldColumnAccessFoldFunctionFoldEqualsAnyNullConvertGroupByToDistinctEliminateGroupByProjectEliminateJoinUnderGroupByLeftEliminateJoinUnderGroupByRightEliminateDistinctReduceGroupingColsReduceNotNullGroupingColsEliminateAggDistinctForKeysEliminateAggFilteredDistinctForKeysEliminateDistinctNoColumnsEliminateEnsureDistinctNoColumnsEliminateDistinctOnValuesPushAggDistinctIntoGroupByPushAggFilterIntoScalarGroupByConvertCountToCountRowsFoldGroupingOperatorsInlineConstVarInlineProjectConstantsInlineSelectConstantsInlineJoinConstantsLeftInlineJoinConstantsRightPushSelectIntoInlinableProjectInlineProjectInProjectCommuteRightJoinSimplifyJoinFiltersDetectJoinContradictionPushFilterIntoJoinLeftAndRightMapFilterIntoJoinLeftMapFilterIntoJoinRightMapEqualityIntoJoinLeftAndRightPushFilterIntoJoinLeftPushFilterIntoJoinRightSimplifyLeftJoinSimplifyRightJoinEliminateSemiJoinSimplifyZeroCardinalitySemiJoinEliminateAntiJoinSimplifyZeroCardinalityAntiJoinEliminateJoinNoColsLeftEliminateJoinNoColsRightHoistJoinProjectRightHoistJoinProjectLeftSimplifyJoinNotNullEqualityExtractJoinEqualitiesSortFiltersInJoinLeftAssociateJoinsLeftLeftAssociateJoinsRightRightAssociateJoinsLeftRightAssociateJoinsRightEliminateLimitEliminateOffsetPushLimitIntoProjectPushOffsetIntoProjectPushLimitIntoOffsetPushLimitIntoOrdinalityPushLimitIntoJoinLeftPushLimitIntoJoinRightFoldLimitsAssociateLimitJoinsLeftAssociateLimitJoinsRightEliminateMax1RowFoldPlusZeroFoldZeroPlusFoldMinusZeroFoldMultOneFoldOneMultFoldDivOneInvertMinusEliminateUnaryMinusSimplifyLimitOrderingSimplifyOffsetOrderingSimplifyGroupByOrderingSimplifyOrdinalityOrderingSimplifyExplainOrderingEliminateJoinUnderProjectLeftEliminateJoinUnderProjectRightEliminateProjectMergeProjectsMergeProjectWithValuesPushColumnRemappingIntoValuesFoldTupleAccessIntoValuesFoldJSONAccessIntoValuesConvertZipArraysToValuesPruneProjectColsPruneScanColsPruneSelectColsPruneLimitColsPruneOffsetColsPruneJoinLeftColsPruneJoinRightColsPruneSemiAntiJoinRightColsPruneAggColsPruneGroupByColsPruneValuesColsPruneOrdinalityColsPruneExplainColsPruneProjectSetColsPruneWindowOutputColsPruneWindowInputColsPruneMutationFetchColsPruneMutationInputColsPruneMutationReturnColsPruneWithScanColsPruneWithColsPruneUnionAllColsRejectNullsLeftJoinRejectNullsRightJoinRejectNullsGroupByCommuteVarCommuteConstEliminateCoalesceSimplifyCoalesceEliminateCastNormalizeInConstFoldInNullUnifyComparisonTypesEliminateExistsZeroRowsEliminateExistsProjectEliminateExistsGroupByIntroduceExistsLimitEliminateExistsLimitNormalizeJSONFieldAccessNormalizeJSONContainsSimplifyCaseWhenConstValueInlineAnyValuesSingleColInlineAnyValuesMultiColSimplifyEqualsAnyTupleSimplifyAnyScalarArrayFoldCollateNormalizeArrayFlattenToAggSimplifySameVarEqualitiesSimplifySameVarInequalitiesSimplifySelectFiltersConsolidateSelectFiltersDetectSelectContradictionEliminateSelectMergeSelectsPushSelectIntoProjectMergeSelectInnerJoinPushSelectCondLeftIntoJoinLeftAndRightPushSelectIntoJoinLeftPushSelectIntoGroupByRemoveNotNullConditionPushSelectIntoProjectSetPushFilterIntoSetOpEliminateUnionAllLeftEliminateUnionAllRightEliminateWindowReduceWindowPartitionColsSimplifyWindowOrderingPushSelectIntoWindowPushLimitIntoWindowInlineWithstartExploreRuleReplaceScalarMinMaxWithLimitReplaceMinWithLimitReplaceMaxWithLimitGenerateStreamingGroupByCommuteJoinCommuteLeftJoinCommuteSemiJoinGenerateMergeJoinsGenerateLookupJoinsGenerateGeoLookupJoinsGenerateZigzagJoinsGenerateInvertedIndexZigzagJoinsGenerateLookupJoinsWithFilterAssociateJoinGenerateLimitedScansPushLimitIntoConstrainedScanPushLimitIntoIndexJoinSplitScanIntoUnionScansGenerateIndexScansGeneratePartialIndexScansGenerateConstrainedScansGenerateInvertedIndexScansSplitDisjunctionSplitDisjunctionAddKeyNumRuleNames"

var _RuleName_index = [...]uint16{0, 15, 35, 48, 76, 94, 114, 133, 148, 163, 179, 195, 209, 223, 238, 253, 266, 279, 290, 302, 313, 329, 341, 350, 358, 382, 402, 424, 445, 467, 489, 511, 533, 556, 566, 583, 607, 634, 647, 667, 697, 724, 737, 765, 795, 810, 831, 851, 872, 899, 929, 952, 979, 1000, 1027, 1049, 1071, 1095, 1115, 1136, 1153, 1173, 1192, 1212, 1229, 1248, 1271, 1295, 1317, 1344, 1369, 1381, 1394, 1412, 1431, 1449, 1460, 1474, 1483, 1493, 1502, 1516, 1524, 1539, 1555, 1567, 1584, 1608, 1631, 1660, 1690, 1707, 1725, 1750, 1777, 1812, 1838, 1870, 1895, 1921, 1951, 1974, 1995, 2009, 2031, 2052, 2075, 2099, 2129, 2151, 2167, 2186, 2209, 2239, 2260, 2282, 2313, 2335, 2358, 2374, 2391, 2408, 2439, 2456, 2487, 2510, 2534, 2555, 2575, 2602, 2623, 2640, 2662, 2685, 2708, 2732, 2746, 2761, 2781, 2802, 2821, 2844, 2865, 2887, 2897, 2920, 2944, 2960, 2972, 2984, 2997, 3008, 3019, 3029, 3040, 3059, 3080, 3102, 3125, 3151, 3174, 3203, 3233, 3249, 3262, 3284, 3313, 3338, 3362, 3386, 3402, 3415, 3430, 3444, 3459, 3476, 3494, 3520, 3532, 3548, 3563, 3582, 3598, 3617, 3638, 3658, 3680, 3702, 3725, 3742, 3755, 3772, 3791, 3811, 3829, 3839, 3851, 3868, 3884, 3897, 3913, 3923, 3943, 3966, 3988, 4010, 4030, 4050, 4074, 4095, 4121, 4145, 4168, 4190, 4212, 4223, 4249, 4274, 4301, 4322, 4346, 4371, 4386, 4398, 4419, 4439, 4477, 4499, 4520, 4542, 4566, 4585, 4606, 4628, 4643, 4668, 4690, 4710, 4729, 4739, 4755, 4783, 4802, 4821, 4845, 4856, 4871, 4886, 4904, 4923, 4945, 4964, 4996, 5025, 5038, 5058, 5086, 5108, 5131, 5149, 5174, 5198, 5224, 5240, 5262, 5274}

func (i RuleName) String() string {
	if i >= RuleName(len(_RuleName_index)-1) {
		return "RuleName(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _RuleName_name[_RuleName_index[i]:_RuleName_index[i+1]]
}
