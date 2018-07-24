package rete

import "github.com/TIBCOSoftware/bego/common/model"

type ruleImpl struct {
	id          int
	name        string
	identifiers []model.TupleTypeAlias
	conditions  []model.Condition
	actionFn    model.ActionFunction
	priority    int
}

//NewRule ... Create a new rule
func NewRule(name string) model.MutableRule {
	rule := ruleImpl{}
	rule.initRuleImpl(name)
	return &rule
}
func (rule *ruleImpl) initRuleImpl(name string) {
	currentNodeID++
	rule.id = currentNodeID
	rule.name = name
}

func (rule *ruleImpl) GetName() string {
	return rule.name
}

func (rule *ruleImpl) GetID() int {
	return rule.id
}

//GetConditions ... get the rule's condition set
func (rule *ruleImpl) GetConditions() []model.Condition {
	return rule.conditions
}

func (rule *ruleImpl) GetIdentifiers() []model.TupleTypeAlias {
	return rule.identifiers
}

func (rule *ruleImpl) AddCondition(conditionName string, idrs []model.TupleTypeAlias, cfn model.ConditionEvaluator) {
	//strIds := make([]string, len(idrs))
	//for i := 0; i < len(idrs); i++ {
	//	strIds[i] = idrs[i].GetName()
	//}
	c := newCondition(conditionName, rule, idrs, cfn)
	rule.conditions = append(rule.conditions, c)
	for _, idr := range idrs {
		// rule.addIdentifier(idr[0], idr[1])
		rule.addIdentifier(idr)
	}
}

func (rule *ruleImpl) addIdentifier(identifierName model.TupleTypeAlias) model.TupleTypeAlias {

	//idrNew := model.TupleTypeAlias(identifierName)
	//TODO: Optimize this, perhaps using a map (need it to be ordered)
	//search for the idr. if exists, skip, else add it

	for _, idr := range rule.identifiers {
		if idr == identifierName {
			return idr
		}
	}

	rule.identifiers = append(rule.identifiers, identifierName)
	return identifierName
}

func (rule *ruleImpl) GetPriority() int {
	return rule.priority
}

func (rule *ruleImpl) SetPriority(priority int) {
	rule.priority = priority
}

func (rule *ruleImpl) SetAction(actionFn model.ActionFunction) {
	rule.actionFn = actionFn
}

func (rule *ruleImpl) GetActionFn() model.ActionFunction {
	return rule.actionFn
}

func (rule *ruleImpl) String() string {
	str := ""
	str += "[Rule: (" + ") " + rule.name + "\n"
	//str += "[Rule: (" + strconv.Itoa(rule.id) + ") " + rule.name + "\n"

	str += "\t[Conditions:\n"
	for _, cond := range rule.conditions {
		str += "\t\t" + cond.String() + "\n"
	}
	// idrs := ""
	// for i := 0; i < len(rule.identifiers); i++ {
	// 	idrs += rule.identifiers[i].String() + ", "
	// }
	str += "\t[Idrs:" + model.IdentifiersToString(rule.identifiers) + "]\n"
	return str
	// return str + idrs + "]\n"
}
