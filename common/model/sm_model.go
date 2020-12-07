package model

import (
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/project-flogo/core/data"
)

type StateMachines struct {
	StateMachines []StateMachine `json:"state-machines"`
	smMap         map[string]*StateMachine
}
type StateMachine struct {
	Descriptor   TupleDescriptor `json:"state-machine"`
	InitialState string          `json:"initial-state"`
	States       []SmState       `json:"states"`
	EndState     string          `json:"end-state"`

	//derived
	stateMap    map[string]*SmState
	ParentSm    string
	ParentState string
}
type SmState struct {
	State        string         `json:"state"`
	EntryAction  string         `json:"entry-action"`
	ExitAction   string         `json:"exit-action"`
	Timeout      int            `json:"timeout"`
	TimeoutState string         `json:"timeout-state"`
	Transitions  []SmTransition `json:"transitions"`
}

type SmTransition struct {
	ToState   string `json:"to-state"`
	Condition string `json:"condition"`
	//ChildSm   string `json:"child-sm"`
	StartSm string `json:"start-sm"`
	//ExitSm    string `json:"exit-sm"`
}

func (sms *StateMachines) UnmarshalJSON(d []byte) error {

	smsSer := &struct {
		SmArr []StateMachine `json:"state-machines"`
	}{}
	if err := yaml.Unmarshal(d, &smsSer); err != nil {
		return err
	}
	sms.StateMachines = smsSer.SmArr

	//setup parent SM links
	sms.smMap = map[string]*StateMachine{}
	for i := range sms.StateMachines {
		sm := &sms.StateMachines[i]
		smName := sm.Descriptor.Name
		_, found := sms.smMap[smName]
		if found {
			return fmt.Errorf("state machine already defined [%s]", smName)
		}
		sms.smMap[smName] = sm
	}

	//todo: detect cycles
	//

	for i := range sms.StateMachines {
		sm := &sms.StateMachines[i]

		for j := range sm.States {
			state := &sm.States[j]
			for k := range state.Transitions {
				transition := &state.Transitions[k]
				if transition.StartSm != "" {
					childSm := sms.smMap[transition.StartSm]
					if childSm == nil {
						return fmt.Errorf("child state machine not found [%s]", transition.StartSm)
					}
					if childSm.ParentSm != "" {
						return fmt.Errorf("in state machine [%s] / state [%s], child state machine [%s] already associated with a "+
							"different parent [%s]",
							sm.Descriptor.Name,
							state.State, transition.StartSm, childSm.ParentSm)
					}

					childSm.ParentSm = sm.Descriptor.Name
					childSm.ParentState = state.State
				}
				sms.addStateEntryForState(sm, transition.ToState)
			}
			sms.addStateEntryForState(sm, state.TimeoutState)
		}
		sms.addStateEntryForState(sm, sm.EndState)

	}

	return nil
}

func (sms *StateMachines) addStateEntryForState(sm *StateMachine, stateName string) {
	stateEntry := sm.stateMap[stateName]
	if stateEntry == nil {
		//add a default
		stateEntry1 := SmState{
			State:        stateName,
			EntryAction:  "",
			ExitAction:   "",
			Timeout:      -1,
			TimeoutState: "",
			Transitions:  []SmTransition{},
		}
		sm.States = append(sm.States, stateEntry1)
		sm.stateMap[stateName] = &stateEntry1
	}
}

func (sms *StateMachines) GetSm(smName string) *StateMachine {
	return sms.smMap[smName]
}

func (s *StateMachine) UnmarshalJSON(d []byte) error {
	ser := &struct {
		Descriptor   TupleDescriptor `json:"state-machine"`
		InitialState string          `json:"initial-state"`
		States       []SmState       `json:"states"`
		EndState     string          `json:"end-state"`
	}{}

	if err := yaml.Unmarshal(d, ser); err != nil {
		return err
	}

	s.InitialState = ser.InitialState
	s.EndState = ser.EndState
	s.States = ser.States

	s.Descriptor = TupleDescriptor{
		Name:         ser.Descriptor.Name,
		TTLInSeconds: ser.Descriptor.TTLInSeconds,
		Props:        ser.Descriptor.Props,
		//do not set the key props here!!
		//keyProps:     ser.Descriptor.GetKeyProps(),
	}
	//add the primary key
	keyProp := TuplePropertyDescriptor{
		Name:     "sm_key",
		PropType: data.TypeString,
		KeyIndex: 0,
	}
	s.Descriptor.Props = append(s.Descriptor.Props, keyProp)

	//add the state key
	stateProp := TuplePropertyDescriptor{
		Name:     "sm_state",
		PropType: data.TypeString,
		KeyIndex: -1,
	}
	s.Descriptor.Props = append(s.Descriptor.Props, stateProp)

	s.stateMap = map[string]*SmState{}

	for i := range s.States {
		state := &s.States[i]
		_, found := s.stateMap[state.State]
		if found {
			return fmt.Errorf("duplicate state entry found [%s]", state.State)
		}
		s.stateMap[state.State] = state
	}

	return nil
}

func (s *StateMachine) GetSmForState(state string) *SmState {
	return s.stateMap[state]
}

func RegisterSmTypes(sms []StateMachine) error {
	for i := range sms {
		smTd := sms[i].Descriptor
		err := RegisterTupleDescriptorsFromTds([]TupleDescriptor{smTd})
		if err != nil {
			return err
		}
	}
	err := registerTimerType()
	if err != nil {
		return err
	}
	return nil
}

func registerTimerType() error {
	td := TupleDescriptor{
		Name:         "timer",
		TTLInSeconds: 0,
		Props: []TuplePropertyDescriptor{
			{
				Name:     "ctime",
				PropType: data.TypeInt,
				KeyIndex: 0,
			},
			{
				Name:     "ctx",
				PropType: data.TypeString,
				KeyIndex: -1,
			},
			{
				Name:     "ruleName",
				PropType: data.TypeString,
				KeyIndex: -1,
			},
		},
	}
	return RegisterTupleDescriptorsFromTds([]TupleDescriptor{td})
}
