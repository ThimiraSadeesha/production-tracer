package schema

import (
	"time"

	"gorm.io/gorm"
)

type OperationProcess struct {
	ID                int            `json:"id" gorm:"primaryKey;column:id"`
	Sequence          int            `json:"sequence" gorm:"column:sequence;default:0"`
	PlanDateTime      LocalDateTime  `json:"planDateTime" gorm:"column:plan_date_time;type:datetime"`
	StartedDateTime   time.Time      `json:"startedDateTime" gorm:"column:started_date_time"`
	CompletedDateTime time.Time      `json:"completedDateTime" gorm:"column:completed_date_time"`
	PlannedValue      float64        `json:"plannedValue" gorm:"column:planned_value;type:decimal(15,2);default:0"`
	ActualValue       float64        `json:"actualValue" gorm:"column:actual_value;type:decimal(15,2);default:0"`
	Status            string         `json:"status" gorm:"type:varchar(100);column:status;index:idx_op_woi_status,priority:2;index:idx_op_wo_scope_status,priority:3"`
	EstimateTime      float64        `json:"estimateTime" gorm:"column:estimate_time;default:0"`
	StartedRoasterID  *int           `json:"startedRoasterId,omitempty" gorm:"column:started_roaster_id"`
	EndRoasterID      *int           `json:"endRoasterId,omitempty" gorm:"column:end_roaster_id"`
	MachineID         *int           `json:"machineId" gorm:"column:machine_id;index:idx_op_order_dup,priority:4;index:idx_op_item_dup,priority:4"`
	WorkOrderItemID   *int           `json:"workOrderItemId,omitempty" gorm:"column:work_order_item_id;index:idx_op_woi_status,priority:1;index:idx_op_item_dup,priority:1"`
	WorkOrderID       *int           `json:"workOrderId,omitempty" gorm:"column:work_order_id;index:idx_op_wo_scope_status,priority:1;index:idx_op_order_dup,priority:1"`
	OperationScope    OperationScope `json:"operationScope" gorm:"type:enum('ITEMS','ORDER','PARTIAL');column:operation_scope;default:ITEMS;index:idx_op_wo_scope_status,priority:2;index:idx_op_order_dup,priority:3;index:idx_op_item_dup,priority:3"`
	ProcessID         int            `json:"processId" gorm:"column:process_id;index:idx_op_order_dup,priority:2;index:idx_op_item_dup,priority:2"`
	IsRedo            bool           `json:"isRedo" gorm:"column:is_redo;default:false"`
	Remark            *string        `json:"remark" gorm:"type:varchar(255);column:remark"`
	CreatedBy         *string        `json:"createdBy" gorm:"type:varchar(255);column:created_by"`
	UpdatedBy         *string        `json:"updatedBy" gorm:"type:varchar(255);column:updated_by"`
	DeletedBy         *string        `json:"deletedBy" gorm:"type:varchar(255);column:deleted_by"`
	DeletedAt         gorm.DeletedAt `json:"deletedAt" gorm:"index;column:deleted_at"`

	// Relationships
	StartedRoaster *ShiftRoster                `json:"startedRoaster,omitempty" gorm:"foreignKey:StartedRoasterID;references:ID"`
	EndRoaster     *ShiftRoster                `json:"endRoaster,omitempty" gorm:"foreignKey:EndRoasterID;references:ID"`
	Machine        *Machine                    `json:"machine,omitempty" gorm:"foreignKey:MachineID"`
	WorkOrder      *WorkOrder                  `json:"workOrder,omitempty" gorm:"foreignKey:WorkOrderID"`
	WorkOrderItem  *WorkOrderItem              `json:"workOrderItem,omitempty" gorm:"foreignKey:WorkOrderItemID"`
	Process        *Process                    `json:"process,omitempty" gorm:"foreignKey:ProcessID"`
	OperationLogs  []OperationLog              `json:"operationLogs,omitempty" gorm:"foreignKey:JobOperationID;references:ID"`
	Attributes     []OperationProcessAttribute `json:"attributes,omitempty" gorm:"foreignKey:OperationProcessID"`
	MakeReadies    []MakeReady                 `json:"makeReadies,omitempty" gorm:"foreignKey:OperationProcessID"`
}

func (OperationProcess) TableName() string {
	return "tbl_operation_process"
}

type OperationScope string

const (
	OperationScopeItems   OperationScope = "ITEMS"
	OperationScopeOrder   OperationScope = "ORDER"
	OperationScopePartial OperationScope = "PARTIAL"
)
