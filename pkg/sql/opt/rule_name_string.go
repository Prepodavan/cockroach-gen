// Code generated by "stringer -output=pkg/sql/opt/rule_name_string.go -type=RuleName pkg/sql/opt/rule_name.go pkg/sql/opt/rule_name.og.go"; DO NOT EDIT.

package opt

import "strconv"

const _RuleName_name = "InvalidRuleNameSimplifyRootOrderingPruneRootColsSimplifyZeroCardinalityGroupNumManualRuleNamesEliminateAggDistinctNormalizeNestedAndsSimplifyTrueAndSimplifyAndTrueSimplifyFalseAndSimplifyAndFalseSimplifyTrueOrSimplifyOrTrueSimplifyFalseOrSimplifyOrFalseFoldNullAndOrFoldNotTrueFoldNotFalseNegateComparisonEliminateNotNegateAndNegateOrExtractRedundantConjunctCommuteVarInequalityCommuteConstInequalityNormalizeCmpPlusConstNormalizeCmpMinusConstNormalizeCmpConstMinusNormalizeTupleEqualityFoldNullComparisonLeftFoldNullComparisonRightFoldIsNullFoldNonNullIsNullFoldIsNotNullFoldNonNullIsNotNullCommuteNullIsDecorrelateJoinDecorrelateProjectSetTryDecorrelateSelectTryDecorrelateProjectTryDecorrelateProjectSelectTryDecorrelateProjectInnerJoinTryDecorrelateInnerJoinTryDecorrelateInnerLeftJoinTryDecorrelateGroupByTryDecorrelateScalarGroupByTryDecorrelateSemiJoinTryDecorrelateLimitOneTryDecorrelateProjectSetHoistSelectExistsHoistSelectNotExistsHoistSelectSubqueryHoistProjectSubqueryHoistJoinSubqueryHoistValuesSubqueryHoistProjectSetSubqueryNormalizeSelectAnyFilterNormalizeJoinAnyFilterNormalizeSelectNotAnyFilterNormalizeJoinNotAnyFilterFoldArrayFoldBinaryFoldUnaryFoldComparisonConvertGroupByToDistinctEliminateDistinctEliminateGroupByProjectReduceGroupingColsEliminateAggDistinctForKeysPushSelectIntoInlinableProjectInlineProjectInProjectSimplifyJoinFiltersDetectJoinContradictionPushFilterIntoJoinLeftAndRightMapFilterIntoJoinLeftMapFilterIntoJoinRightPushFilterIntoJoinLeftPushFilterIntoJoinRightSimplifyLeftJoinWithoutFiltersSimplifyRightJoinWithoutFiltersSimplifyLeftJoinWithFiltersSimplifyRightJoinWithFiltersEliminateSemiJoinEliminateAntiJoinEliminateJoinNoColsLeftEliminateJoinNoColsRightHoistJoinProjectSimplifyJoinNotNullEqualityExtractJoinEqualitiesEliminateLimitPushLimitIntoProjectPushOffsetIntoProjectEliminateMax1RowFoldPlusZeroFoldZeroPlusFoldMinusZeroFoldMultOneFoldOneMultFoldDivOneInvertMinusEliminateUnaryMinusSimplifyLimitOrderingSimplifyOffsetOrderingSimplifyGroupByOrderingSimplifyRowNumberOrderingSimplifyExplainOrderingSimplifyMutationOrderingEliminateProjectMergeProjectsMergeProjectWithValuesPruneProjectColsPruneScanColsPruneSelectColsPruneLimitColsPruneOffsetColsPruneJoinLeftColsPruneJoinRightColsPruneAggColsPruneGroupByColsPruneValuesColsPruneRowNumberColsPruneExplainColsPruneProjectSetColsRejectNullsLeftJoinRejectNullsRightJoinRejectNullsGroupByCommuteVarCommuteConstEliminateCoalesceSimplifyCoalesceEliminateCastFoldNullCastFoldNullUnaryFoldNullBinaryLeftFoldNullBinaryRightFoldNullInNonEmptyFoldNullInEmptyFoldNullNotInEmptyNormalizeInConstFoldInNullUnifyComparisonTypesEliminateExistsProjectEliminateExistsGroupByNormalizeJSONFieldAccessNormalizeJSONContainsSimplifyCaseWhenConstValueSimplifyEqualsAnyTupleSimplifyAnyScalarArrayFoldCollateNormalizeArrayFlattenToAggSimplifySelectFiltersDetectSelectContradictionEliminateSelectMergeSelectsPushSelectIntoProjectMergeSelectInnerJoinPushSelectCondLeftIntoJoinLeftAndRightPushSelectCondRightIntoJoinLeftAndRightPushSelectIntoJoinLeftPushSelectIntoJoinRightPushSelectIntoGroupByRemoveNotNullConditionEliminateUnionAllLeftEliminateUnionAllRightstartExploreRuleReplaceMinWithLimitReplaceMaxWithLimitGenerateStreamingGroupByCommuteJoinCommuteLeftJoinCommuteRightJoinGenerateMergeJoinsGenerateLookupJoinsGenerateZigzagJoinsGenerateLookupJoinsWithFilterGenerateLimitedScansPushLimitIntoConstrainedScanPushLimitIntoIndexJoinGenerateIndexScansGenerateConstrainedScansGenerateInvertedIndexScansNumRuleNames"

var _RuleName_index = [...]uint16{0, 15, 35, 48, 76, 94, 114, 133, 148, 163, 179, 195, 209, 223, 238, 253, 266, 277, 289, 305, 317, 326, 334, 358, 378, 400, 421, 443, 465, 487, 509, 532, 542, 559, 572, 592, 605, 620, 641, 661, 682, 709, 739, 762, 789, 810, 837, 859, 881, 905, 922, 942, 961, 981, 998, 1017, 1040, 1064, 1086, 1113, 1138, 1147, 1157, 1166, 1180, 1204, 1221, 1244, 1262, 1289, 1319, 1341, 1360, 1383, 1413, 1434, 1456, 1478, 1501, 1531, 1562, 1589, 1617, 1634, 1651, 1674, 1698, 1714, 1741, 1762, 1776, 1796, 1817, 1833, 1845, 1857, 1870, 1881, 1892, 1902, 1913, 1932, 1953, 1975, 1998, 2023, 2046, 2070, 2086, 2099, 2121, 2137, 2150, 2165, 2179, 2194, 2211, 2229, 2241, 2257, 2272, 2290, 2306, 2325, 2344, 2364, 2382, 2392, 2404, 2421, 2437, 2450, 2462, 2475, 2493, 2512, 2530, 2545, 2563, 2579, 2589, 2609, 2631, 2653, 2677, 2698, 2724, 2746, 2768, 2779, 2805, 2826, 2851, 2866, 2878, 2899, 2919, 2957, 2996, 3018, 3041, 3062, 3084, 3105, 3127, 3143, 3162, 3181, 3205, 3216, 3231, 3247, 3265, 3284, 3303, 3332, 3352, 3380, 3402, 3420, 3444, 3470, 3482}

func (i RuleName) String() string {
	if i >= RuleName(len(_RuleName_index)-1) {
		return "RuleName(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _RuleName_name[_RuleName_index[i]:_RuleName_index[i+1]]
}
