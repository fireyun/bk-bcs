/*
Tencent is pleased to support the open source community by making Basic Service Configuration Platform available.
Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
Licensed under the MIT License (the "License"); you may not use this file except
in compliance with the License. You may obtain a copy of the License at
http://opensource.org/licenses/MIT
Unless required by applicable law or agreed to in writing, software distributed under
the License is distributed on an "as IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
either express or implied. See the License for the specific language governing permissions and
limitations under the License.
*/

// Package table NOTES
package table

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"bscp.io/pkg/criteria/enumor"
)

// Columns defines the column's details prepared for
// ORM operation usage.
type Columns struct {
	descriptors []ColumnDescriptor
	// columns defines all the table's columns
	columns []string
	// columnType defines a table's column and it's type
	columnType map[string]enumor.ColumnType
	// columnExpr is the joined columns with comma
	columnExpr string
	// namedExpr is the joined 'named' columns with comma.
	// which column may have the 'prefix'.
	// such as App.Spec.Name named columns should be:
	// "spec.name"
	namedExpr string
	// ColonNameExpr is the joined 'named' columns with comma.
	// which column may have the 'prefix' and each column is
	// prefix with a colon.
	// such as App.Spec.Name named columns should be:
	// ":spec.name"
	colonNameExpr string
}

// Columns returns all the db columns
func (col Columns) Columns() []string {
	copied := make([]string, len(col.columns))
	for idx := range col.columns {
		copied[idx] = col.columns[idx]
	}
	return copied
}

// ColumnTypes returns each column and it's data type
func (col Columns) ColumnTypes() map[string]enumor.ColumnType {
	copied := make(map[string]enumor.ColumnType)
	for k, v := range col.columnType {
		copied[k] = v
	}

	return copied
}

// ColumnExpr returns the joined columns with comma
func (col Columns) ColumnExpr() string {
	return col.columnExpr
}

// NamedExpr returns the joined 'named' columns with comma
// like: "name as spec.name"
func (col Columns) NamedExpr() string {
	return col.namedExpr
}

// ColonNameExpr returns the joined 'named' columns with comma and
// prefixed with colon, like: ":spec.name"
func (col Columns) ColonNameExpr() string {
	return col.colonNameExpr
}

// WithoutColumn remove one or more columns from the 'origin' columns.
func (col Columns) WithoutColumn(column ...string) map[string]enumor.ColumnType {
	if len(column) == 0 {
		return col.ColumnTypes()
	}

	reminder := make(map[string]bool)
	for _, one := range column {
		reminder[one] = true
	}

	copied := make(map[string]enumor.ColumnType)
	for col, typ := range col.columnType {
		if reminder[col] {
			continue
		}
		copied[col] = typ
	}

	return copied
}

// mergeColumns merge table columns together.
func mergeColumns(all ...ColumnDescriptors) *Columns {
	tc := &Columns{
		descriptors: make([]ColumnDescriptor, 0),
		columns:     make([]string, 0),
		columnType:  make(map[string]enumor.ColumnType),
		columnExpr:  "",
		namedExpr:   "",
	}
	if len(all) == 0 {
		return tc
	}

	namedExpr := make([]string, 0)
	colorExpr := make([]string, 0)

	for _, nc := range all {
		for _, col := range nc {
			tc.descriptors = append(tc.descriptors, col)
			tc.columnType[col.Column] = col.Type
			tc.columns = append(tc.columns, col.Column)

			if col.Column == col.NamedC {
				namedExpr = append(namedExpr, col.Column)
				colorExpr = append(colorExpr, col.Column)
			} else {
				namedExpr = append(namedExpr, fmt.Sprintf("%s as '%s'", col.Column, col.NamedC))
				colorExpr = append(colorExpr, col.NamedC)
			}
		}
	}

	tc.columnExpr = strings.Join(tc.columns, ", ")
	tc.namedExpr = strings.Join(namedExpr, ", ")
	tc.colonNameExpr = ":" + strings.Join(colorExpr, ", :")
	return tc
}

// ColumnDescriptor defines a table's column related information.
type ColumnDescriptor struct {
	// Column is column's name
	Column string
	// NamedC is named column's name
	NamedC string
	// Type is this column's data type.
	Type enumor.ColumnType
	_    struct{}
}

// ColumnDescriptors is a collection of ColumnDescriptor
type ColumnDescriptors []ColumnDescriptor

// mergeColumnDescriptors merge column descriptors to one map.
func mergeColumnDescriptors(prefix string, namedC ...ColumnDescriptors) ColumnDescriptors {
	if len(namedC) == 0 {
		return make([]ColumnDescriptor, 0)
	}

	merged := make([]ColumnDescriptor, 0)
	if len(prefix) == 0 {
		for _, one := range namedC {
			merged = append(merged, one...)
		}
	} else {
		for _, one := range namedC {
			for _, col := range one {
				col.NamedC = prefix + "." + col.NamedC
				merged = append(merged, col)
			}
		}
	}

	return merged
}

// Tables defines all the database table
// related resources.
type Tables interface {
	TableName() Name
}

// Name is database table's name type
type Name string

// Validate whether the table name is valid or not.
func (n Name) Validate() error {
	switch n {
	case AppTable:
	case ArchivedAppTable:
	case ContentTable:
	case ConfigItemTable:
	case CommitsTable:
	case ReleaseTable:
	case ReleasedConfigItemTable:
	case StrategySetTable:
	case StrategyTable:
	case GroupCategoryTable:
	case GroupTable:
	case CurrentPublishedStrategyTable:
	case PublishedStrategyHistoryTable:
	case CurrentReleasedInstanceTable:
	case EventTable:
	case ShardingDBTable:
	case ShardingBizTable:
	case IDGeneratorTable:
	case AuditTable:
	case ResourceLockTable:
	default:
		return fmt.Errorf("unknown table name: %s", n)
	}

	return nil
}

const (
	// AppTable is app table's name
	AppTable Name = "application"
	// ArchivedAppTable is archived app table's name
	ArchivedAppTable Name = "archived_app"
	// ContentTable is content table's name
	ContentTable Name = "content"
	// ConfigItemTable is config item table's name
	ConfigItemTable Name = "config_item"
	// CommitsTable is commits table's name
	CommitsTable Name = "commits"
	// ReleaseTable is release table's name
	ReleaseTable Name = "releases"
	// ReleasedConfigItemTable is released config item table's name
	ReleasedConfigItemTable Name = "released_config_item"
	// GroupCategory is group category table's name
	GroupCategoryTable Name = "group_category"
	// GroupTable is group table's name
	GroupTable Name = "group_"
	// StrategySetTable is strategy set table's name
	StrategySetTable Name = "strategy_set"
	// StrategyTable is strategy table's name
	StrategyTable Name = "strategy"
	// CurrentPublishedStrategyTable is current published strategy table's name
	CurrentPublishedStrategyTable Name = "current_published_strategy"
	// PublishedStrategyHistoryTable is published strategy history table's name
	PublishedStrategyHistoryTable Name = "published_strategy_history"
	// CurrentReleasedInstanceTable is current released instance table's name
	CurrentReleasedInstanceTable Name = "current_released_instance"
	// EventTable is event table's name
	EventTable Name = "event"
	// ShardingDBTable is sharding db table's name
	ShardingDBTable Name = "sharding_db"
	// ShardingBizTable is sharding biz table's name
	ShardingBizTable Name = "sharding_biz"
	// IDGeneratorTable is id generator table's name
	IDGeneratorTable Name = "id_generator"
	// AuditTable is audit table's name
	AuditTable Name = "audit"
	// ResourceLockTable is lock table's name
	ResourceLockTable Name = "resource_lock"
)

// RevisionColumns defines all the Revision table's columns.
var RevisionColumns = mergeColumns(RevisionColumnDescriptor)

// RevisionColumnDescriptor is Revision's column descriptors.
var RevisionColumnDescriptor = ColumnDescriptors{
	{Column: "creator", NamedC: "creator", Type: enumor.String},
	{Column: "reviser", NamedC: "reviser", Type: enumor.String},
	{Column: "created_at", NamedC: "created_at", Type: enumor.Time},
	{Column: "updated_at", NamedC: "updated_at", Type: enumor.Time},
}

// Revision is a resource's status information
type Revision struct {
	Creator   string    `db:"creator" json:"creator"`
	Reviser   string    `db:"reviser" json:"reviser"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// IsEmpty test whether a revision is empty or not.
func (r Revision) IsEmpty() bool {
	if len(r.Creator) != 0 {
		return false
	}

	if len(r.Reviser) != 0 {
		return false
	}

	if !r.CreatedAt.IsZero() {
		return false
	}

	if !r.UpdatedAt.IsZero() {
		return false
	}

	return true
}

const lagSeconds = 5 * 60

// ValidateCreate validate revision when created
func (r Revision) ValidateCreate() error {

	if len(r.Creator) == 0 {
		return errors.New("creator can not be empty")
	}

	now := time.Now().Unix()
	if (r.CreatedAt.Unix() <= (now - lagSeconds)) || (r.CreatedAt.Unix() >= (now + lagSeconds)) {
		return errors.New("invalid create time")
	}

	return nil
}

// ValidateUpdate validate revision when updated
func (r Revision) ValidateUpdate() error {
	if len(r.Reviser) == 0 {
		return errors.New("reviser can not be empty")
	}

	if len(r.Creator) != 0 {
		return errors.New("creator can not be updated")
	}

	if !r.CreatedAt.IsZero() {
		return errors.New("create_time can not be updated")
	}

	now := time.Now().Unix()
	if (r.UpdatedAt.Unix() <= (now - lagSeconds)) || (r.UpdatedAt.Unix() >= (now + lagSeconds)) {
		return errors.New("invalid update time")
	}

	return nil
}

// CreatedRevisionColumns defines all the Revision table's columns.
var CreatedRevisionColumns = mergeColumns(CreatedRevisionColumnDescriptor)

// CreatedRevisionColumnDescriptor is CreatedRevision's column descriptors.
var CreatedRevisionColumnDescriptor = ColumnDescriptors{
	{Column: "creator", NamedC: "creator", Type: enumor.String},
	{Column: "created_at", NamedC: "created_at", Type: enumor.Time}}

// CreatedRevision is a resource's reversion information being created.
type CreatedRevision struct {
	Creator   string    `db:"creator" json:"creator"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Validate revision when created
func (r CreatedRevision) Validate() error {

	if len(r.Creator) == 0 {
		return errors.New("creator can not be empty")
	}

	now := time.Now().Unix()
	if (r.CreatedAt.Unix() <= (now - lagSeconds)) || (r.CreatedAt.Unix() >= (now + lagSeconds)) {
		return errors.New("invalid create time")
	}

	return nil
}
