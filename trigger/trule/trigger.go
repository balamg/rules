package trule

import (
	"fmt"

	"github.com/project-flogo/rules/ruleapi"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/utils"

	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
	"github.com/project-flogo/rules/common"
)

var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{}, &Reply{})

func init() {
	_ = trigger.Register(&Trigger{}, &Factory{})
}

type Factory struct {
}

// Metadata implements trigger.Factory.Metadata
func (*Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

// New implements trigger.Factory.New
func (*Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	s := &Settings{}
	err := metadata.MapToStruct(config.Settings, s, true)
	if err != nil {
		return nil, err
	}

	return &Trigger{id: config.Id, settings: s}, nil
}

type Trigger struct {
	settings *Settings
	id       string
	logger   log.Logger
	//	serverInstanceID string
	handlerMap     map[string]trigger.Handler
	actionDataChan <-chan model.ActionData
}

func (t *Trigger) Initialize(ctx trigger.InitContext) error {

	t.logger = ctx.Logger()
	rs := ruleapi.GetRuleSession(t.settings.Rs)
	if rs == nil {
		return fmt.Errorf("rule session not set [%s]", t.settings.Rs)
	} else {
		fmt.Printf("FOUND RS [%s]\n", t.settings.Rs)
	}
	t.actionDataChan = rs.GetActionDataChannel()

	t.handlerMap = map[string]trigger.Handler{}
	// Init handlers
	for _, handler := range ctx.GetHandlers() {

		s := &HandlerSettings{}
		err := metadata.MapToStruct(handler.Settings(), s, true)
		if err != nil {
			return err
		}
		name := s.Name

		rule := addRulesFromSettings(s.Name, s.Condition, s.Priority)
		err = rs.AddRule(rule)
		if err != nil {
			return fmt.Errorf("ERROR during setting up rules: %s", err)
		}

		t.logger.Infof("Added rule [%s]", name)

		t.handlerMap[name] = handler
	}

	t.logger.Debugf("Configured on rulesession %d", t.settings.Rs)

	return nil
}

func addRulesFromSettings(name string, condition string, priority int) model.Rule {
	rule := ruleapi.NewRule(name)
	rule.SetContext("This is a test of context")
	rule.SetFlowBasedAction(nil)
	rule.SetPriority(priority)
	rule.AddExprCondition(condition, condition, nil)

	/*
		state-machine=SM1
		initial-state: S1

		state S1, condition: C1 next: S2, timeout: T1, timeout-state:S3

	*/

	//rule.AddExprCondition("x", "$.sm1.state == S1 && C1")
	//rule.SetAction(setNextState)

	//	//now add explicit rule identifiers if any
	//	if ruleCfg.Identifiers != nil {
	//		idrs := []model.TupleType{}
	//		for _, idr := range ruleCfg.Identifiers {
	//			idrs = append(idrs, model.TupleType(idr))
	//		}
	//		rule.AddIdrsToRule(idrs)
	//	}
	//
	//	rs.AddRule(rule)
	//}
	return rule
}

//
//func setNextState(ctx context.Context, session model.RuleSession, s string, m map[model.TupleType]model.Tuple, context model.RuleContext) {
//	//get the name of the sm-type from ruleCtx
//	var t model.StateMachine
//	t, ok := m["sm-name"].(model.StateMachine)
//	if ok {
//		t.CancelTimer()
//	}
//	t.SetState(nextState)
//	t.SetTimer(forNextState)
//}

func (t *Trigger) Start() error {
	fmt.Printf("STARTED\n")
	go func() {
		for action := range t.actionDataChan {

			handler := t.handlerMap[action.RuleName]
			if handler != nil {
				rs_uid, _ := common.GetUniqueId()
				utils.SetVar(rs_uid, action.RuleSession)
				ctx_uid, _ := common.GetUniqueId()
				utils.SetVar(ctx_uid, action.Context)
				out := &Output{
					Rs:       rs_uid,
					Ctx:      ctx_uid,
					Rulename: action.RuleName,
					Tuples:   model.TuplesToMap(action.Tuples),
				}

				result, err := handler.Handle(action.Context, out)
				if err != nil {
					fmt.Printf("error while invoking handler [%s]", err)
				} else {
					fmt.Printf("rule invocation result: [%s][%v]", out.Rulename, result)
				}

				utils.RemoveVar(rs_uid)
				utils.RemoveVar(ctx_uid)

			} else {
				fmt.Printf("No handler defined for rule [%s]", action.RuleName)
			}
			action.Done <- true

		}
	}()
	fmt.Printf("STARTED complete\n")
	return nil
}

// Stop implements util.Managed.Stop
func (t *Trigger) Stop() error {
	return nil
}

//
//func newActionHandler(rt *Trigger, method string, handler trigger.Handler) httprouter.Handle {
//
//	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//
//		logger := rt.logger
//
//		logger.Debugf("Received request for id '%s'", rt.id)
//
//		c := cors.New(CorsPrefix, logger)
//		c.WriteCorsActualRequestHeaders(w)
//
//		//add server instance id to response
//		//w.Header().Add("X-Server-Instance-Id", rt.serverInstanceID)
//
//		out := &Output{}
//		out.Method = method
//
//		out.PathParams = make(map[string]string)
//		for _, param := range ps {
//			out.PathParams[param.Key] = param.Value
//		}
//
//		queryValues := r.URL.Query()
//		out.QueryParams = make(map[string]string, len(queryValues))
//		out.Headers = make(map[string]string, len(r.Header))
//
//		for key, value := range r.Header {
//			out.Headers[key] = strings.Join(value, ",")
//		}
//
//		for key, value := range queryValues {
//			out.QueryParams[key] = strings.Join(value, ",")
//		}
//
//		// Check the HTTP Header Content-Type
//		contentType := r.Header.Get("Content-Type")
//		switch contentType {
//		case "application/x-www-form-urlencoded":
//			buf := new(bytes.Buffer)
//			_, err := buf.ReadFrom(r.Body)
//			if err != nil {
//				logger.Debugf("Error reading body: %s", err.Error())
//				http.Error(w, err.Error(), http.StatusBadRequest)
//				return
//			}
//
//			s := buf.String()
//			m, err := url.ParseQuery(s)
//			if err != nil {
//				logger.Debugf("Error parsing query string: %s", err.Error())
//				http.Error(w, err.Error(), http.StatusBadRequest)
//				return
//			}
//
//			content := make(map[string]interface{}, 0)
//			for key, val := range m {
//				if len(val) == 1 {
//					content[key] = val[0]
//				} else {
//					content[key] = val[0]
//				}
//			}
//
//			out.Content = content
//		case "application/json":
//			var content interface{}
//			err := json.NewDecoder(r.Body).Decode(&content)
//			if err != nil {
//				switch {
//				case err == io.EOF:
//					// empty body
//					//todo what should handler say if content is expected?
//				default:
//					logger.Debugf("Error parsing json body: %s", err.Error())
//					http.Error(w, err.Error(), http.StatusBadRequest)
//					return
//				}
//			}
//			out.Content = content
//		default:
//			if strings.Contains(contentType, "multipart/form-data") {
//				// need to still extract the body, only handling the multipart data for now...
//
//				if err := r.ParseMultipartForm(32); err != nil {
//					logger.Debugf("Error parsing multipart form: %s", err.Error())
//					http.Error(w, err.Error(), http.StatusBadRequest)
//					return
//				}
//
//				var files []map[string]interface{}
//
//				for key, fh := range r.MultipartForm.File {
//					for _, header := range fh {
//
//						fileDetails, err := getFileDetails(key, header)
//						if err != nil {
//							logger.Debugf("Error getting attached file details: %s", err.Error())
//							http.Error(w, err.Error(), http.StatusBadRequest)
//							return
//						}
//
//						files = append(files, fileDetails)
//					}
//				}
//
//				// The content output from the trigger
//				content := map[string]interface{}{
//					"body":  nil,
//					"files": files,
//				}
//				out.Content = content
//			} else {
//				b, err := ioutil.ReadAll(r.Body)
//				if err != nil {
//					logger.Debugf("Error reading body: %s", err.Error())
//					http.Error(w, err.Error(), http.StatusBadRequest)
//					return
//				}
//
//				out.Content = string(b)
//			}
//		}
//
//		results, err := handler.Handle(context.Background(), out)
//		if err != nil {
//			logger.Debugf("Error handling request: %s", err.Error())
//			http.Error(w, err.Error(), http.StatusBadRequest)
//			return
//		}
//
//		if logger.TraceEnabled() {
//			logger.Tracef("Action Results: %#v", results)
//		}
//
//		reply := &Reply{}
//		err = reply.FromMap(results)
//		if err != nil {
//			logger.Debugf("Error mapping results: %s", err.Error())
//			http.Error(w, err.Error(), http.StatusBadRequest)
//			return
//		}
//
//		// add response headers
//		if len(reply.Headers) > 0 {
//			if logger.TraceEnabled() {
//				logger.Tracef("Adding Headers")
//			}
//
//			for key, value := range reply.Headers {
//				w.Header().Set(key, value)
//			}
//		}
//
//		if len(reply.Cookies) > 0 {
//			if logger.TraceEnabled() {
//				logger.Tracef("Adding Cookies")
//			}
//
//			err := addCookies(w, reply.Cookies)
//			if err != nil {
//				logger.Debugf("Error handling request: %s", err.Error())
//				http.Error(w, err.Error(), http.StatusInternalServerError)
//				return
//			}
//		}
//
//		if reply.Code == 0 {
//			reply.Code = http.StatusOK
//		}
//
//		if reply.Data != nil {
//
//			if logger.DebugEnabled() {
//				logger.Debugf("The http reply code is: %d", reply.Code)
//				logger.Debugf("The http reply data is: %#v", reply.Data)
//			}
//
//			switch t := reply.Data.(type) {
//			case string:
//				var v interface{}
//				err := json.Unmarshal([]byte(t), &v)
//				if err != nil {
//					//Not a json
//					w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
//				} else {
//					//Json
//					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//				}
//
//				w.WriteHeader(reply.Code)
//				_, err = w.Write([]byte(t))
//				if err != nil {
//					logger.Debugf("Error writing body: %s", err.Error())
//				}
//				return
//			default:
//				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//				w.WriteHeader(reply.Code)
//				if err := json.NewEncoder(w).Encode(reply.Data); err != nil {
//					logger.Debugf("Error encoding json reply: %s", err.Error())
//				}
//				return
//			}
//		}
//
//		logger.Debugf("The reply http code is: %d", reply.Code)
//		w.WriteHeader(reply.Code)
//	}
//}
