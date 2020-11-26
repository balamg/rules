package model

import (
	"github.com/ghodss/yaml"
	"github.com/project-flogo/core/data"
)

type StateMachineModel struct {
	Name         string          `json:"name"`
	Descriptor   TupleDescriptor `json:"descriptor"`
	InitialState string          `json:"initial-state"`
	States       []SmState       `json:"states"`
	EndState     string          `json:"end-state"`
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
	ToState      string `json:"to-state"`
	Condition    string `json:"condition"`
	EntryAction  string `json:"entry-action"`
	ExitAction   string `json:"exit-action"`
	Timeout      int    `json:"timeout"`
	TimeoutState string `json:"timeout-state"`
}

func (s *StateMachineModel) UnmarshalJSON(d []byte) error {
	ser := &struct {
		Name         string          `json:"name"`
		Descriptor   TupleDescriptor `json:"descriptor"`
		InitialState string          `json:"initial-state"`
		States       []SmState       `json:"states"`
		EndState     string          `json:"end-state"`
	}{}

	if err := yaml.Unmarshal(d, ser); err != nil {
		return err
	}

	s.Name = ser.Name
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
	return nil
}

func RegisterSmTypes(sms []StateMachineModel) error {
	for i := range sms {
		smTd := sms[i].Descriptor
		err := RegisterTupleDescriptorsFromTds([]TupleDescriptor{smTd})
		if err != nil {
			return err
		}
	}
	return nil
}
