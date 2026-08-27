package schema

var Migrations = []interface{}{
	// Roots (no foreign keys)
	&Department{},
	&Unit{},
	&HoldReason{},
	&Device{},

	// Depend on the roots
	&Operator{},
	&Process{},
	&ProcessAttribute{},
	&MachineType{},
	&Machine{},
	&MachineHoldReason{},
	&Shift{},
	&ShiftRoster{},

	// Work orders
	&WorkOrder{},
	&WorkOrderItem{},

	// Operations
	&OperationProcess{},
	&OperationLog{},
	&OperationProcessAttribute{},
	&MachineLog{},

	// Make-ready
	&MakeReady{},
	&MakeReadyAttributes{},
	&MakeReadyLog{},
}
