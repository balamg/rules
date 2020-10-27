package model

import (
	"context"
)

// RuleContext associated with every rule
type RuleContext interface{}

//Rule ... a Rule interface
type Rule interface {
	GetName() string
	GetIdentifiers() []TupleType
	GetConditions() []Condition
	GetActionFn() ActionFunction
	String() string
	GetPriority() int
	GetDeps() map[TupleType]map[string]bool
	GetContext() RuleContext
}

//MutableRule interface has methods to add conditions and actions
type MutableRule interface {
	Rule
	AddCondition(conditionName string, idrs []string, cFn ConditionEvaluator, ctx RuleContext) (err error)
	SetAction(actionFn ActionFunction)
	SetPriority(priority int)
	SetContext(ctx RuleContext)
	AddExprCondition(conditionName string, cExpr string, ctx RuleContext) error
	AddIdrsToRule(idrs []TupleType)
	//will try to coerce value to the named property's type
	//SetValue(ctx context.Context, name string, value interface{}) error
	//SetValues(ctx context.Context, values map[string]interface{}) error
	SetFlowBasedAction(ctx RuleContext) error
}

//Condition interface to maintain/get various condition properties
type Condition interface {
	GetName() string
	GetRule() Rule
	GetIdentifiers() []TupleType
	GetContext() RuleContext
	String() string
	Evaluate(string, string, map[TupleType]Tuple, RuleContext) (bool, error)
}

// RuleSession to maintain rules and assert tuples against those rules
type RuleSession interface {
	GetName() string

	AddRule(rule Rule) (err error)
	DeleteRule(ruleName string)
	GetRules() []Rule

	Assert(ctx context.Context, tuple Tuple) (err error)
	Retract(ctx context.Context, tuple Tuple)

	ScheduleAssert(ctx context.Context, delayInMillis uint64, key interface{}, tuple Tuple)
	CancelScheduledAssert(ctx context.Context, key interface{})

	Unregister()

	//Optional, called before asserting a tuple but after adding all rules
	SetStartupFunction(startupFn StartupRSFunction)

	GetStartupFunction() (startupFn StartupRSFunction)

	//To be called when the rule session is ready to start accepting tuples
	//This will invoke the StartupFunction
	Start(startupCtx map[string]interface{}) (err error)

	//return the asserted tuple, nil if not found
	GetAssertedTuple(key TupleKey) Tuple

	//Retract, and remove
	Delete(ctx context.Context, tuple Tuple)

	//RtcTransactionHandler
	RegisterRtcTransactionHandler(txnHandler RtcTransactionHandler, handlerCtx interface{})

	//replay existing tuples into a rule
	ReplayTuplesForRule(ruleName string) (err error)
	GetActionDataChannel() chan ActionData
}

//ConditionEvaluator is a function pointer for handling condition evaluations on the server side
//i.e, part of the server side API
type ConditionEvaluator func(string, string, map[TupleType]Tuple, RuleContext) bool

//ActionFunction is a function pointer for handling action callbacks on the server side
//i.e part of the server side API
type ActionFunction func(context.Context, RuleSession, string, map[TupleType]Tuple, RuleContext)

//StartupRSFunction is called once after creation of a RuleSession
type StartupRSFunction func(ctx context.Context, rs RuleSession, sessionCtx map[string]interface{}) (err error)

// ValueChangeListener to pickup and process tuple value changes
type ValueChangeListener interface {
	OnValueChange(tuple Tuple, prop string)
}

type RtcTxn interface {
	//map of type and map of key/tuple
	GetRtcAdded() map[string]map[string]Tuple
	GetRtcModified() map[string]map[string]RtcModified
	GetRtcDeleted() map[string]map[string]Tuple
}

type RtcModified interface {
	GetTuple() Tuple
	GetModifiedProps() map[string]bool
}

type RtcTransactionHandler func(ctx context.Context, rs RuleSession, txn RtcTxn, txnContext interface{})

type ActionData struct {
	RuleSession RuleSession
	Context     context.Context
	RuleName    string
	Tuples      map[TupleType]Tuple
	Done        chan bool
}

type StateMachine interface {
	MutableTuple
	CancelTimer()
	SetTimer()
	Start()
}
