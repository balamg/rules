package trule

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	Rs string `md:"rs,required"` // rule session
}

type HandlerSettings struct {
	Rulename string `md:"rulename,required"` // handler for this rule
}

type Output struct {
	Rs       string                 `md:"rs"`       // The path parameters (e.g., 'id' in http://.../pet/:id/name )
	Ctx      string                 `md:"ctx"`      // The query parameters (e.g., 'id' in http://.../pet?id=someValue )
	Rulename string                 `md:"rulename"` // The HTTP header parameters
	Tuples   map[string]interface{} `md:"tuples"`   // The content of the request
}

type Reply struct{}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"rs":       o.Rs,
		"ctx":      o.Ctx,
		"rulename": o.Rulename,
		"tuples":   o.Tuples,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	var err error
	o.Rs, err = coerce.ToString(values["rs"])
	if err != nil {
		return err
	}
	o.Ctx, err = coerce.ToString(values["ctx"])
	if err != nil {
		return err
	}
	o.Rulename, err = coerce.ToString(values["rulename"])
	if err != nil {
		return err
	}
	o.Tuples, err = coerce.ToObject(values["tuples"])
	if err != nil {
		return err
	}

	return nil
}

func (r *Reply) ToMap() map[string]interface{} {
	return map[string]interface{}{}
}

func (r *Reply) FromMap(values map[string]interface{}) error {

	return nil
}
