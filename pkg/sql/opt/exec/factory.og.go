// Code generated by optgen; DO NOT EDIT.

package exec

import (
	"github.com/cockroachdb/cockroach/pkg/sql/opt"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/cat"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/constraint"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/invertedexpr"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlbase"
)

// Factory defines the interface for building an execution plan, which consists
// of a tree of execution nodes (currently a sql.planNode tree).
//
// The tree is always built bottom-up. The Construct methods either construct
// leaf nodes, or they take other nodes previously constructed by this same
// factory as children.
//
// The TypedExprs passed to these functions refer to columns of the input node
// via IndexedVars.
type Factory interface {
	// ConstructPlan creates a plan enclosing the given plan and (optionally)
	// subqueries, cascades, and checks.
	//
	// Subqueries are executed before the root tree, which can refer to subquery
	// results using tree.Subquery nodes.
	//
	// Cascades are executed after the root tree. They can return more cascades
	// and checks which should also be executed.
	//
	// Checks are executed after all cascades have been executed. They don't
	// return results but can generate errors (e.g. foreign key check failures).
	ConstructPlan(
		root Node, subqueries []Subquery, cascades []Cascade, checks []Node,
	) (Plan, error)

	// ConstructScan creates a node for a Scan operation.
	//
	// Scan runs a scan of a specified index of a table, possibly with an index
	// constraint and/or a hard limit.
	ConstructScan(
		table cat.Table,
		index cat.Index,
		params ScanParams,
		reqOrdering OutputOrdering,
	) (Node, error)

	// ConstructValues creates a node for a Values operation.
	ConstructValues(
		rows [][]tree.TypedExpr,
		columns sqlbase.ResultColumns,
	) (Node, error)

	// ConstructFilter creates a node for a Filter operation.
	//
	// Filter applies a filter on the results of the given input node.
	ConstructFilter(
		input Node,
		filter tree.TypedExpr,
		reqOrdering OutputOrdering,
	) (Node, error)

	// ConstructInvertedFilter creates a node for a InvertedFilter operation.
	//
	// InvertedFilter applies a span expression on the results of the given input
	// node.
	ConstructInvertedFilter(
		input Node,
		invFilter *invertedexpr.SpanExpression,
		invColumn NodeColumnOrdinal,
	) (Node, error)

	// ConstructSimpleProject creates a node for a SimpleProject operation.
	//
	// SimpleProject applies a "simple" projection on the results of the given input
	// node. A simple projection is one that does not involve new expressions; it's
	// just a reshuffling of columns. This is a more efficient version of
	// ConstructRender.  The colNames argument is optional; if it is nil, the names
	// of the corresponding input columns are kept.
	ConstructSimpleProject(
		input Node,
		cols []NodeColumnOrdinal,
		colNames []string,
		reqOrdering OutputOrdering,
	) (Node, error)

	// ConstructRender creates a node for a Render operation.
	//
	// Render applies a projection on the results of the given input node. The
	// projection can contain new expressions. The input expression slice will be
	// modified.
	ConstructRender(
		input Node,
		columns sqlbase.ResultColumns,
		exprs tree.TypedExprs,
		reqOrdering OutputOrdering,
	) (Node, error)

	// ConstructApplyJoin creates a node for a ApplyJoin operation.
	//
	// ApplyJoin runs an apply join between an input node (the left side of the join)
	// and a RelExpr that has outer columns (the right side of the join) by replacing
	// the outer columns of the right side RelExpr with data from each row of the
	// left side of the join according to the data in leftBoundColMap. The apply join
	// can be any kind of join except for right outer and full outer.
	//
	// To plan the right-hand side, planRightSideFn must be called for each left
	// row. This function generates a plan (using the same factory) that produces
	// the rightColumns (in order).
	//
	// onCond is the join condition.
	ConstructApplyJoin(
		joinType sqlbase.JoinType,
		left Node,
		rightColumns sqlbase.ResultColumns,
		onCond tree.TypedExpr,
		planRightSideFn ApplyJoinPlanRightSideFn,
	) (Node, error)

	// ConstructHashJoin creates a node for a HashJoin operation.
	//
	// HashJoin runs a hash-join between the results of two input nodes.
	//
	// The leftEqColsAreKey/rightEqColsAreKey flags, if set, indicate that the
	// equality columns form a key in the left/right input.
	//
	// The extraOnCond expression can refer to columns from both inputs using
	// IndexedVars (first the left columns, then the right columns).
	ConstructHashJoin(
		joinType sqlbase.JoinType,
		left Node,
		right Node,
		leftEqCols []NodeColumnOrdinal,
		rightEqCols []NodeColumnOrdinal,
		leftEqColsAreKey bool,
		rightEqColsAreKey bool,
		extraOnCond tree.TypedExpr,
	) (Node, error)

	// ConstructMergeJoin creates a node for a MergeJoin operation.
	//
	// The ON expression can refer to columns from both inputs using IndexedVars
	// (first the left columns, then the right columns). In addition, the i-th
	// column in leftOrdering is constrained to equal the i-th column in
	// rightOrdering. The directions must match between the two orderings.
	ConstructMergeJoin(
		joinType sqlbase.JoinType,
		left Node,
		right Node,
		onCond tree.TypedExpr,
		leftOrdering sqlbase.ColumnOrdering,
		rightOrdering sqlbase.ColumnOrdering,
		reqOrdering OutputOrdering,
		leftEqColsAreKey bool,
		rightEqColsAreKey bool,
	) (Node, error)

	// ConstructInterleavedJoin creates a node for a InterleavedJoin operation.
	//
	// InterleavedJoin runs a join between two interleaved tables. One table is the
	// ancestor of the other in the interleaving hierarchy (as per leftIsAncestor).
	//
	// Semantically, the join is identical to a merge-join between these two
	// tables, where the equality columns are all the index column of the ancestor
	// index.
	//
	// The two scans are guaranteed to have the same direction, and to not have
	// any hard limits.
	//
	// Since the interleaved joiner does a single scan for both tables, only the
	// Locking clause for the ancestor is used.
	//
	ConstructInterleavedJoin(
		joinType sqlbase.JoinType,
		leftTable cat.Table,
		leftIndex cat.Index,
		leftParams ScanParams,
		leftFilter tree.TypedExpr,
		rightTable cat.Table,
		rightIndex cat.Index,
		rightParams ScanParams,
		rightFilter tree.TypedExpr,
		leftIsAncestor bool,
		onCond tree.TypedExpr,
		reqOrdering OutputOrdering,
	) (Node, error)

	// ConstructGroupBy creates a node for a GroupBy operation.
	//
	// GroupBy runs an aggregation. A set of aggregations is performed for each group
	// of values on the groupCols.
	//
	// If the input is guaranteed to have an ordering on grouping columns, a
	// "streaming" aggregation is performed (i.e. aggregation happens separately
	// for each distinct set of values on the set of columns in the ordering).
	ConstructGroupBy(
		input Node,
		groupCols []NodeColumnOrdinal,
		groupColOrdering sqlbase.ColumnOrdering,
		aggregations []AggInfo,
		reqOrdering OutputOrdering,
	) (Node, error)

	// ConstructScalarGroupBy creates a node for a ScalarGroupBy operation.
	//
	// ScalarGroupBy runs a scalar aggregation, i.e.  one which performs a set of
	// aggregations on all the input rows (as a single group) and has exactly one
	// result row (even when there are no input rows).
	ConstructScalarGroupBy(
		input Node,
		aggregations []AggInfo,
	) (Node, error)

	// ConstructDistinct creates a node for a Distinct operation.
	//
	// Distinct filters out rows such that only the first row is kept for each set of
	// values along the distinct columns.  The orderedCols are a subset of
	// distinctCols; the input is required to be ordered along these columns (i.e.
	// all rows with the same values on these columns are a contiguous part of the
	// input).
	ConstructDistinct(
		input Node,
		distinctCols NodeColumnOrdinalSet,
		orderedCols NodeColumnOrdinalSet,
		reqOrdering OutputOrdering,
		nullsAreDistinct bool,
		errorOnDup string,
	) (Node, error)

	// ConstructSetOp creates a node for a SetOp operation.
	//
	// SetOp performs a UNION / INTERSECT / EXCEPT operation (either the ALL or the
	// DISTINCT version). The left and right nodes must have the same number of
	// columns.
	ConstructSetOp(
		typ tree.UnionType,
		all bool,
		left Node,
		right Node,
	) (Node, error)

	// ConstructSort creates a node for a Sort operation.
	//
	// Sort performs a resorting of the rows produced by the input node.
	//
	// When the input is partially sorted we can execute a "segmented" sort. In
	// this case alreadyOrderedPrefix is non-zero and the input is ordered by
	// ordering[:alreadyOrderedPrefix].
	ConstructSort(
		input Node,
		ordering sqlbase.ColumnOrdering,
		alreadyOrderedPrefix int,
	) (Node, error)

	// ConstructOrdinality creates a node for a Ordinality operation.
	//
	// Ordinality appends an ordinality column to each row in the input node.
	ConstructOrdinality(
		input Node,
		colName string,
	) (Node, error)

	// ConstructIndexJoin creates a node for a IndexJoin operation.
	//
	// IndexJoin performs an index join. The input contains the primary key (on the
	// columns identified as keyCols).
	//
	// The index join produces the given table columns (in ordinal order).
	ConstructIndexJoin(
		input Node,
		table cat.Table,
		keyCols []NodeColumnOrdinal,
		tableCols TableColumnOrdinalSet,
		reqOrdering OutputOrdering,
	) (Node, error)

	// ConstructLookupJoin creates a node for a LookupJoin operation.
	//
	// LookupJoin performs a lookup join.
	//
	// The eqCols are columns from the input used
	// as keys for the columns of the index (or a prefix of them); lookupCols are
	// ordinals for the table columns we are retrieving.
	//
	// The node produces the columns in the input and (unless join type is
	// LeftSemiJoin or LeftAntiJoin) the lookupCols, ordered by ordinal. The ON
	// condition can refer to these using IndexedVars.
	ConstructLookupJoin(
		joinType sqlbase.JoinType,
		input Node,
		table cat.Table,
		index cat.Index,
		eqCols []NodeColumnOrdinal,
		eqColsAreKey bool,
		lookupCols TableColumnOrdinalSet,
		onCond tree.TypedExpr,
		reqOrdering OutputOrdering,
	) (Node, error)

	// ConstructInvertedJoin creates a node for a InvertedJoin operation.
	//
	// InvertedJoin performs a lookup join into an inverted index.
	//
	// invertedExpr is used along with inputCol (a column from the input) to
	// find the keys to look up in the index; lookupCols are ordinals for the
	// table columns we are retrieving.
	//
	// The node produces the columns in the input and (unless join type is
	// LeftSemiJoin or LeftAntiJoin) the lookupCols, ordered by ordinal. The ON
	// condition can refer to these using IndexedVars. Note that lookupCols
	// includes the inverted column.
	ConstructInvertedJoin(
		joinType sqlbase.JoinType,
		invertedExpr tree.TypedExpr,
		input Node,
		table cat.Table,
		index cat.Index,
		inputCol NodeColumnOrdinal,
		lookupCols TableColumnOrdinalSet,
		onCond tree.TypedExpr,
		reqOrdering OutputOrdering,
	) (Node, error)

	// ConstructZigzagJoin creates a node for a ZigzagJoin operation.
	//
	// ZigzagJoin performs a zigzag join.
	//
	// Each side of the join has two kinds of columns that form a prefix
	// of the specified index: fixed columns (with values specified in
	// fixedVals), and equal columns (with column ordinals specified in
	// {left,right}EqCols). The lengths of leftEqCols and rightEqCols
	// must match.
	ConstructZigzagJoin(
		leftTable cat.Table,
		leftIndex cat.Index,
		rightTable cat.Table,
		rightIndex cat.Index,
		leftEqCols []NodeColumnOrdinal,
		rightEqCols []NodeColumnOrdinal,
		leftCols NodeColumnOrdinalSet,
		rightCols NodeColumnOrdinalSet,
		onCond tree.TypedExpr,
		fixedVals []Node,
		reqOrdering OutputOrdering,
	) (Node, error)

	// ConstructLimit creates a node for a Limit operation.
	//
	// Limit implements LIMIT and/or OFFSET on the results of the given node. If one
	// or the other is not needed, then it is set to nil.
	ConstructLimit(
		input Node,
		limit tree.TypedExpr,
		offset tree.TypedExpr,
	) (Node, error)

	// ConstructMax1Row creates a node for a Max1Row operation.
	//
	// Max1Row permits at most one row from the given input node, causing an error
	// with the given text at runtime if the node tries to return more than one row.
	ConstructMax1Row(
		input Node,
		errorText string,
	) (Node, error)

	// ConstructProjectSet creates a node for a ProjectSet operation.
	//
	// ProjectSet performs a lateral cross join between the output of the given node
	// and the functional zip of the given expressions.
	ConstructProjectSet(
		input Node,
		exprs tree.TypedExprs,
		zipCols sqlbase.ResultColumns,
		numColsPerGen []int,
	) (Node, error)

	// ConstructWindow creates a node for a Window operation.
	//
	// Window executes a window function over the given node.
	ConstructWindow(
		input Node,
		window WindowInfo,
	) (Node, error)

	// ConstructRenameColumns creates a node for a RenameColumns operation.
	//
	// RenameColumns modifies the column names of a node.
	ConstructRenameColumns(
		input Node,
		colNames []string,
	) (Node, error)

	// ConstructExplainOpt creates a node for a ExplainOpt operation.
	//
	// Explain implements EXPLAIN (OPT), showing information about the given plan.
	ConstructExplainOpt(
		plan string,
		envOpts ExplainEnvData,
	) (Node, error)

	// ConstructExplain creates a node for a Explain operation.
	//
	// Explain implements EXPLAIN, showing information about the given plan.
	ConstructExplain(
		options *tree.ExplainOptions,
		stmtType tree.StatementType,
		plan Plan,
	) (Node, error)

	// ConstructShowTrace creates a node for a ShowTrace operation.
	//
	// ShowTrace implements a SHOW TRACE FOR SESSION statement.
	ConstructShowTrace(
		typ tree.ShowTraceType,
		compact bool,
	) (Node, error)

	// ConstructInsert creates a node for a Insert operation.
	//
	// Insert implements an INSERT statement (including ON CONFLICT DO NOTHING, but
	// not other ON CONFLICT clauses).
	//
	// The input columns are inserted into a subset of columns in the table, in the
	// same order they're defined. The insertCols set contains the ordinal positions
	// of columns in the table into which values are inserted. All columns are
	// expected to be present except delete-only mutation columns, since those do not
	// need to participate in an insert operation.
	//
	// If allowAutoCommit is set, the operator is allowed to commit the
	// transaction (if appropriate, i.e. if it is in an implicit transaction).
	// This is false if there are multiple mutations in a statement, or the output
	// of the mutation is processed through side-effecting expressions.
	ConstructInsert(
		input Node,
		table cat.Table,
		insertCols TableColumnOrdinalSet,
		returnCols TableColumnOrdinalSet,
		checkCols CheckOrdinalSet,
		allowAutoCommit bool,
	) (Node, error)

	// ConstructInsertFastPath creates a node for a InsertFastPath operation.
	//
	// InsertFastPath implements a special (but very common) case of insert,
	// satisfying the following conditions:
	//  - the input is Values with at most InsertFastPathMaxRows, and there are no
	//    subqueries;
	//  - there are no other mutations in the statement, and the output of the
	//    insert is not processed through side-effecting expressions (see
	//    allowAutoCommit flag for ConstructInsert);
	//  - there are no self-referencing foreign keys;
	//  - all FK checks can be performed using direct lookups into unique indexes.
	//
	// In this case, the foreign-key checks can run before (or even concurrently
	// with) the insert. If they are run before, the insert is allowed to
	// auto-commit.
	ConstructInsertFastPath(
		rows [][]tree.TypedExpr,
		table cat.Table,
		insertCols TableColumnOrdinalSet,
		returnCols TableColumnOrdinalSet,
		checkCols CheckOrdinalSet,
		fkChecks []InsertFastPathFKCheck,
	) (Node, error)

	// ConstructUpdate creates a node for a Update operation.
	//
	// Update implements an UPDATE statement. The input contains columns that were
	// fetched from the target table, and that provide existing values that can be
	// used to formulate the new encoded value that will be written back to the table
	// (updating any column in a family requires having the values of all other
	// columns). The input also contains computed columns that provide new values for
	// any updated columns.
	//
	// The fetchCols and updateCols sets contain the ordinal positions of the
	// fetch and update columns in the target table. The input must contain those
	// columns in the same order as they appear in the table schema, with the
	// fetch columns first and the update columns second.
	//
	// The passthrough parameter contains all the result columns that are part of
	// the input node that the update node needs to return (passing through from
	// the input). The pass through columns are used to return any column from the
	// FROM tables that are referenced in the RETURNING clause.
	//
	// If allowAutoCommit is set, the operator is allowed to commit the
	// transaction (if appropriate, i.e. if it is in an implicit transaction).
	// This is false if there are multiple mutations in a statement, or the output
	// of the mutation is processed through side-effecting expressions.
	ConstructUpdate(
		input Node,
		table cat.Table,
		fetchCols TableColumnOrdinalSet,
		updateCols TableColumnOrdinalSet,
		returnCols TableColumnOrdinalSet,
		checks CheckOrdinalSet,
		passthrough sqlbase.ResultColumns,
		allowAutoCommit bool,
	) (Node, error)

	// ConstructUpsert creates a node for a Upsert operation.
	//
	// Upsert implements an INSERT..ON CONFLICT DO UPDATE or UPSERT statement.
	//
	// For each input row, Upsert will test the canaryCol. If it is null, then it
	// will insert a new row. If not-null, then Upsert will update an existing row.
	// The input is expected to contain the columns to be inserted, followed by the
	// columns containing existing values, and finally the columns containing new
	// values.
	//
	// The length of each group of input columns can be up to the number of
	// columns in the given table. The insertCols, fetchCols, and updateCols sets
	// contain the ordinal positions of the table columns that are involved in
	// the Upsert. For example:
	//
	//   CREATE TABLE abc (a INT PRIMARY KEY, b INT, c INT)
	//   INSERT INTO abc VALUES (10, 20, 30) ON CONFLICT (a) DO UPDATE SET b=25
	//
	//   insertCols = {0, 1, 2}
	//   fetchCols  = {0, 1, 2}
	//   updateCols = {1}
	//
	// The input is expected to first have 3 columns that will be inserted into
	// columns {0, 1, 2} of the table. The next 3 columns contain the existing
	// values of columns {0, 1, 2} of the table. The last column contains the
	// new value for column {1} of the table.
	//
	// If allowAutoCommit is set, the operator is allowed to commit the
	// transaction (if appropriate, i.e. if it is in an implicit transaction).
	// This is false if there are multiple mutations in a statement, or the output
	// of the mutation is processed through side-effecting expressions.
	ConstructUpsert(
		input Node,
		table cat.Table,
		canaryCol NodeColumnOrdinal,
		insertCols TableColumnOrdinalSet,
		fetchCols TableColumnOrdinalSet,
		updateCols TableColumnOrdinalSet,
		returnCols TableColumnOrdinalSet,
		checks CheckOrdinalSet,
		allowAutoCommit bool,
	) (Node, error)

	// ConstructDelete creates a node for a Delete operation.
	//
	// Delete implements a DELETE statement. The input contains columns that were
	// fetched from the target table, and that will be deleted.
	//
	// The fetchCols set contains the ordinal positions of the fetch columns in
	// the target table. The input must contain those columns in the same order
	// as they appear in the table schema.
	//
	// If allowAutoCommit is set, the operator is allowed to commit the
	// transaction (if appropriate, i.e. if it is in an implicit transaction).
	// This is false if there are multiple mutations in a statement, or the output
	// of the mutation is processed through side-effecting expressions.
	ConstructDelete(
		input Node,
		table cat.Table,
		fetchCols TableColumnOrdinalSet,
		returnCols TableColumnOrdinalSet,
		allowAutoCommit bool,
	) (Node, error)

	// ConstructDeleteRange creates a node for a DeleteRange operation.
	//
	// DeleteRange efficiently deletes contiguous rows stored in the given table's
	// primary index. This fast path is only possible when certain conditions hold
	// true:
	//  - there are no secondary indexes;
	//  - the input to the delete is a scan (without limits);
	//  - the table is not involved in interleaving, or it is at the root of an
	//    interleaving hierarchy with cascading FKs such that a delete of a row
	//    cascades and deletes all interleaved rows corresponding to that row;
	//  - there are no inbound FKs to the table (other than within the
	//    interleaving as described above).
	//
	// See the comment for ConstructScan for descriptions of the needed and
	// indexConstraint parameters, since DeleteRange combines Delete + Scan into a
	// single operator.
	//
	// If any interleavedTables are passed, they are all the descendant tables in
	// an interleaving hierarchy we are deleting from.
	ConstructDeleteRange(
		table cat.Table,
		needed TableColumnOrdinalSet,
		indexConstraint *constraint.Constraint,
		interleavedTables []cat.Table,
		maxReturnedKeys int,
		allowAutoCommit bool,
	) (Node, error)

	// ConstructCreateTable creates a node for a CreateTable operation.
	//
	// CreateTable implements a CREATE TABLE statement.
	ConstructCreateTable(
		input Node,
		schema cat.Schema,
		ct *tree.CreateTable,
	) (Node, error)

	// ConstructCreateView creates a node for a CreateView operation.
	//
	// CreateView implements a CREATE VIEW statement.
	ConstructCreateView(
		schema cat.Schema,
		viewName *cat.DataSourceName,
		ifNotExists bool,
		replace bool,
		temporary bool,
		viewQuery string,
		columns sqlbase.ResultColumns,
		deps opt.ViewDeps,
	) (Node, error)

	// ConstructSequenceSelect creates a node for a SequenceSelect operation.
	//
	// SequenceSelect implements a scan of a sequence as a data source.
	ConstructSequenceSelect(
		sequence cat.Sequence,
	) (Node, error)

	// ConstructSaveTable creates a node for a SaveTable operation.
	//
	// SaveTable passes through all the input rows unchanged, but also creates a
	// table and inserts all the rows into it.
	ConstructSaveTable(
		input Node,
		table *cat.DataSourceName,
		colNames []string,
	) (Node, error)

	// ConstructErrorIfRows creates a node for a ErrorIfRows operation.
	//
	// ErrorIfRows returns no results, but causes an execution error if the input
	// returns any rows.
	ConstructErrorIfRows(
		input Node,
		// mkErr is used to create the error; it is passed an input row.
		mkErr MkErrFn,
	) (Node, error)

	// ConstructOpaque creates a node for a Opaque operation.
	//
	// Opaque implements operators that have no relational inputs and which require
	// no specific treatment by the optimizer.
	ConstructOpaque(
		metadata opt.OpaqueMetadata,
	) (Node, error)

	// ConstructAlterTableSplit creates a node for a AlterTableSplit operation.
	//
	// AlterTableSplit implements ALTER TABLE/INDEX SPLIT AT.
	ConstructAlterTableSplit(
		index cat.Index,
		input Node,
		expiration tree.TypedExpr,
	) (Node, error)

	// ConstructAlterTableUnsplit creates a node for a AlterTableUnsplit operation.
	//
	// AlterTableUnsplit implements ALTER TABLE/INDEX UNSPLIT AT.
	ConstructAlterTableUnsplit(
		index cat.Index,
		input Node,
	) (Node, error)

	// ConstructAlterTableUnsplitAll creates a node for a AlterTableUnsplitAll operation.
	//
	// AlterTableUnsplitAll implements ALTER TABLE/INDEX UNSPLIT ALL.
	ConstructAlterTableUnsplitAll(
		index cat.Index,
	) (Node, error)

	// ConstructAlterTableRelocate creates a node for a AlterTableRelocate operation.
	//
	// AlterTableRelocate implements ALTER TABLE/INDEX UNSPLIT AT.
	ConstructAlterTableRelocate(
		index cat.Index,
		input Node,
		relocateLease bool,
	) (Node, error)

	// ConstructBuffer creates a node for a Buffer operation.
	//
	// Buffer passes through the input rows but also saves them in a buffer, which
	// can be referenced from elsewhere in the query (using ScanBuffer).
	ConstructBuffer(
		input Node,
		label string,
	) (Node, error)

	// ConstructScanBuffer creates a node for a ScanBuffer operation.
	//
	// ScanBuffer refers to a node constructed by Buffer or passed to
	// RecursiveCTEIterationFn.
	ConstructScanBuffer(
		ref Node,
		label string,
	) (Node, error)

	// ConstructRecursiveCTE creates a node for a RecursiveCTE operation.
	//
	// RecursiveCTE executes a recursive CTE:
	//   * the initial plan is run first; the results are emitted and also saved
	//     in a buffer.
	//   * so long as the last buffer is not empty:
	//     - the RecursiveCTEIterationFn is used to create a plan for the
	//       recursive side; a reference to the last buffer is passed to this
	//       function. The returned plan uses this reference with a
	//       ConstructScanBuffer call.
	//     - the plan is executed; the results are emitted and also saved in a new
	//       buffer for the next iteration.
	ConstructRecursiveCTE(
		initial Node,
		fn RecursiveCTEIterationFn,
		label string,
	) (Node, error)

	// ConstructControlJobs creates a node for a ControlJobs operation.
	//
	// ControlJobs implements PAUSE/CANCEL/RESUME JOBS.
	ConstructControlJobs(
		command tree.JobCommand,
		input Node,
	) (Node, error)

	// ConstructControlSchedules creates a node for a ControlSchedules operation.
	//
	// ControlSchedules implements PAUSE/CANCEL/DROP SCHEDULES.
	ConstructControlSchedules(
		command tree.ScheduleCommand,
		input Node,
	) (Node, error)

	// ConstructCancelQueries creates a node for a CancelQueries operation.
	//
	// CancelQueries implements CANCEL QUERIES.
	ConstructCancelQueries(
		input Node,
		ifExists bool,
	) (Node, error)

	// ConstructCancelSessions creates a node for a CancelSessions operation.
	//
	// CancelSessions implements CANCEL SESSIONS.
	ConstructCancelSessions(
		input Node,
		ifExists bool,
	) (Node, error)

	// ConstructExport creates a node for a Export operation.
	//
	// Export implements EXPORT.
	ConstructExport(
		input Node,
		fileName tree.TypedExpr,
		fileFormat string,
		options []KVOption,
	) (Node, error)
}

// StubFactory is a do-nothing implementation of Factory, used for testing.
type StubFactory struct{}

var _ Factory = StubFactory{}

func (StubFactory) ConstructPlan(
	root Node, subqueries []Subquery, cascades []Cascade, checks []Node,
) (Plan, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructScan(
	table cat.Table,
	index cat.Index,
	params ScanParams,
	reqOrdering OutputOrdering,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructValues(
	rows [][]tree.TypedExpr,
	columns sqlbase.ResultColumns,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructFilter(
	input Node,
	filter tree.TypedExpr,
	reqOrdering OutputOrdering,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructInvertedFilter(
	input Node,
	invFilter *invertedexpr.SpanExpression,
	invColumn NodeColumnOrdinal,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructSimpleProject(
	input Node,
	cols []NodeColumnOrdinal,
	colNames []string,
	reqOrdering OutputOrdering,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructRender(
	input Node,
	columns sqlbase.ResultColumns,
	exprs tree.TypedExprs,
	reqOrdering OutputOrdering,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructApplyJoin(
	joinType sqlbase.JoinType,
	left Node,
	rightColumns sqlbase.ResultColumns,
	onCond tree.TypedExpr,
	planRightSideFn ApplyJoinPlanRightSideFn,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructHashJoin(
	joinType sqlbase.JoinType,
	left Node,
	right Node,
	leftEqCols []NodeColumnOrdinal,
	rightEqCols []NodeColumnOrdinal,
	leftEqColsAreKey bool,
	rightEqColsAreKey bool,
	extraOnCond tree.TypedExpr,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructMergeJoin(
	joinType sqlbase.JoinType,
	left Node,
	right Node,
	onCond tree.TypedExpr,
	leftOrdering sqlbase.ColumnOrdering,
	rightOrdering sqlbase.ColumnOrdering,
	reqOrdering OutputOrdering,
	leftEqColsAreKey bool,
	rightEqColsAreKey bool,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructInterleavedJoin(
	joinType sqlbase.JoinType,
	leftTable cat.Table,
	leftIndex cat.Index,
	leftParams ScanParams,
	leftFilter tree.TypedExpr,
	rightTable cat.Table,
	rightIndex cat.Index,
	rightParams ScanParams,
	rightFilter tree.TypedExpr,
	leftIsAncestor bool,
	onCond tree.TypedExpr,
	reqOrdering OutputOrdering,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructGroupBy(
	input Node,
	groupCols []NodeColumnOrdinal,
	groupColOrdering sqlbase.ColumnOrdering,
	aggregations []AggInfo,
	reqOrdering OutputOrdering,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructScalarGroupBy(
	input Node,
	aggregations []AggInfo,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructDistinct(
	input Node,
	distinctCols NodeColumnOrdinalSet,
	orderedCols NodeColumnOrdinalSet,
	reqOrdering OutputOrdering,
	nullsAreDistinct bool,
	errorOnDup string,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructSetOp(
	typ tree.UnionType,
	all bool,
	left Node,
	right Node,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructSort(
	input Node,
	ordering sqlbase.ColumnOrdering,
	alreadyOrderedPrefix int,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructOrdinality(
	input Node,
	colName string,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructIndexJoin(
	input Node,
	table cat.Table,
	keyCols []NodeColumnOrdinal,
	tableCols TableColumnOrdinalSet,
	reqOrdering OutputOrdering,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructLookupJoin(
	joinType sqlbase.JoinType,
	input Node,
	table cat.Table,
	index cat.Index,
	eqCols []NodeColumnOrdinal,
	eqColsAreKey bool,
	lookupCols TableColumnOrdinalSet,
	onCond tree.TypedExpr,
	reqOrdering OutputOrdering,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructInvertedJoin(
	joinType sqlbase.JoinType,
	invertedExpr tree.TypedExpr,
	input Node,
	table cat.Table,
	index cat.Index,
	inputCol NodeColumnOrdinal,
	lookupCols TableColumnOrdinalSet,
	onCond tree.TypedExpr,
	reqOrdering OutputOrdering,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructZigzagJoin(
	leftTable cat.Table,
	leftIndex cat.Index,
	rightTable cat.Table,
	rightIndex cat.Index,
	leftEqCols []NodeColumnOrdinal,
	rightEqCols []NodeColumnOrdinal,
	leftCols NodeColumnOrdinalSet,
	rightCols NodeColumnOrdinalSet,
	onCond tree.TypedExpr,
	fixedVals []Node,
	reqOrdering OutputOrdering,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructLimit(
	input Node,
	limit tree.TypedExpr,
	offset tree.TypedExpr,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructMax1Row(
	input Node,
	errorText string,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructProjectSet(
	input Node,
	exprs tree.TypedExprs,
	zipCols sqlbase.ResultColumns,
	numColsPerGen []int,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructWindow(
	input Node,
	window WindowInfo,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructRenameColumns(
	input Node,
	colNames []string,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructExplainOpt(
	plan string,
	envOpts ExplainEnvData,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructExplain(
	options *tree.ExplainOptions,
	stmtType tree.StatementType,
	plan Plan,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructShowTrace(
	typ tree.ShowTraceType,
	compact bool,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructInsert(
	input Node,
	table cat.Table,
	insertCols TableColumnOrdinalSet,
	returnCols TableColumnOrdinalSet,
	checkCols CheckOrdinalSet,
	allowAutoCommit bool,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructInsertFastPath(
	rows [][]tree.TypedExpr,
	table cat.Table,
	insertCols TableColumnOrdinalSet,
	returnCols TableColumnOrdinalSet,
	checkCols CheckOrdinalSet,
	fkChecks []InsertFastPathFKCheck,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructUpdate(
	input Node,
	table cat.Table,
	fetchCols TableColumnOrdinalSet,
	updateCols TableColumnOrdinalSet,
	returnCols TableColumnOrdinalSet,
	checks CheckOrdinalSet,
	passthrough sqlbase.ResultColumns,
	allowAutoCommit bool,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructUpsert(
	input Node,
	table cat.Table,
	canaryCol NodeColumnOrdinal,
	insertCols TableColumnOrdinalSet,
	fetchCols TableColumnOrdinalSet,
	updateCols TableColumnOrdinalSet,
	returnCols TableColumnOrdinalSet,
	checks CheckOrdinalSet,
	allowAutoCommit bool,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructDelete(
	input Node,
	table cat.Table,
	fetchCols TableColumnOrdinalSet,
	returnCols TableColumnOrdinalSet,
	allowAutoCommit bool,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructDeleteRange(
	table cat.Table,
	needed TableColumnOrdinalSet,
	indexConstraint *constraint.Constraint,
	interleavedTables []cat.Table,
	maxReturnedKeys int,
	allowAutoCommit bool,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructCreateTable(
	input Node,
	schema cat.Schema,
	ct *tree.CreateTable,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructCreateView(
	schema cat.Schema,
	viewName *cat.DataSourceName,
	ifNotExists bool,
	replace bool,
	temporary bool,
	viewQuery string,
	columns sqlbase.ResultColumns,
	deps opt.ViewDeps,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructSequenceSelect(
	sequence cat.Sequence,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructSaveTable(
	input Node,
	table *cat.DataSourceName,
	colNames []string,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructErrorIfRows(
	input Node,
	mkErr MkErrFn,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructOpaque(
	metadata opt.OpaqueMetadata,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructAlterTableSplit(
	index cat.Index,
	input Node,
	expiration tree.TypedExpr,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructAlterTableUnsplit(
	index cat.Index,
	input Node,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructAlterTableUnsplitAll(
	index cat.Index,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructAlterTableRelocate(
	index cat.Index,
	input Node,
	relocateLease bool,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructBuffer(
	input Node,
	label string,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructScanBuffer(
	ref Node,
	label string,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructRecursiveCTE(
	initial Node,
	fn RecursiveCTEIterationFn,
	label string,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructControlJobs(
	command tree.JobCommand,
	input Node,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructControlSchedules(
	command tree.ScheduleCommand,
	input Node,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructCancelQueries(
	input Node,
	ifExists bool,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructCancelSessions(
	input Node,
	ifExists bool,
) (Node, error) {
	return struct{}{}, nil
}

func (StubFactory) ConstructExport(
	input Node,
	fileName tree.TypedExpr,
	fileFormat string,
	options []KVOption,
) (Node, error) {
	return struct{}{}, nil
}