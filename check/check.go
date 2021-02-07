package check

import (
	"encoding/json"
	"time"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/iwanteggroll/compliancebeat/config"
	"github.com/iwanteggroll/compliancebeat/models"
	"github.com/iwanteggroll/compliancebeat/scripting"
)

type ComplianceCheck struct {
	period   time.Duration
	name     string
	enabled  bool
	category string
	path     string
	program  string
	params   []string
}

func StatusText(status int) string {
	switch status {
	case 0:
		return "OK"
	case 1:
		return "WARNING"
	case 2:
		return "CRITICAL"
	case 3:
		return "UNKNOWN"
	}
	return "INVALID"
}

func (complianceCheck *ComplianceCheck) Setup(config *config.CheckConfig) scripting.Script {

	if config == nil {
		logp.Err("Invalid Compliance Check config.")
		return nil
	}

	if config.Name == "" {
		logp.Err("Must specify name.")
		return nil
	}

	complianceCheck.name = config.Name

	if config.Period == "" {
		logp.Err("Must specific period")
		return nil
	}

	period, err := time.ParseDuration(config.Period)

	if err != nil {
		logp.Err("Couldn't read period.")
		return nil
	}

	complianceCheck.period = period

	// if config.Enabled == nil {
	// 	return errors.New("Must specify whether the script is enabled or not.")
	// }
	complianceCheck.enabled = config.Enabled

	if config.Program == "" {
		logp.Err("Must specify powershell or python.")
		return nil
	}

	if config.Category == "" {
		logp.Err("Must specify category of check.")
		return nil
	}

	complianceCheck.category = config.Category

	complianceCheck.program = config.Program

	if config.Path == "" {
		logp.Err("Must specify absolute or relative path to script.")
		return nil
	}

	complianceCheck.params = config.Params

	complianceCheck.path = config.Path

	var sc scripting.Script
	if config.Program == "powershell" {
		pos := scripting.PowershellScript{ScriptName: complianceCheck.name, ScriptPath: complianceCheck.path, ScriptParams: complianceCheck.params}
		sc = &pos
		//fmt.Printf("%+v\n", sc) // remove for prod
		logp.Debug("examplebeat-Setup", "In powershell config if.")
	}

	if config.Program == "python" {
		pys := scripting.PythonScript{ScriptName: complianceCheck.name, ScriptPath: complianceCheck.path, ScriptParams: complianceCheck.params}
		sc = &pys
		logp.Debug("examplebeat-Setup", "In python config if.")
	}

	return sc
}

func (complianceCheck *ComplianceCheck) Run(publish func([]beat.Event), sc scripting.Script) {

	if !complianceCheck.enabled {
		logp.Info("Check %s not starting; disabled in config.", complianceCheck.name)
		return
	}

	logp.Info("Starting check %s with period of %s.", complianceCheck.name, complianceCheck.period.String())

	ticker := time.NewTicker(complianceCheck.period)

	defer ticker.Stop()
	for range ticker.C {
		events, err := complianceCheck.Check(sc)

		if err != nil {
			logp.Err("Check Error: %q: %v", complianceCheck.name, err)
		}

		publish(events)
	}
}

func (complianceCheck *ComplianceCheck) Check(sc scripting.Script) (events []beat.Event, err error) {

	// find powershell or python
	output, err := sc.Execute()

	if err != nil {
		logp.Err("Errors using %s to run %s", complianceCheck.program, complianceCheck.name)
		return
	}

	// execute script and get results
	//jsonString := `{"StatusMessages":[{"ComplianceCheckLevel":{"Status":0},"ComplianceCheckFunction":"Check-CurrentAccountForKerberos","ComplianceCheckMessageText":"Passed check.","ComplianceCheckTimestamp":"2020-06-19T14:45:52.0891727Z"},{"ComplianceCheckLevel":{"Status":0},"ComplianceCheckFunction":"Check-ADDomainConnectionAndTimeSkew","ComplianceCheckMessageText":"Passed check.","ComplianceCheckTimestamp":"2020-06-19T14:45:52.1047982Z"}],"ResultLevel":{"Status":0},"ExecutionId":"4ece732d-51ad-455a-b471-94f617d73997"}`

	var eventReceiver models.ComplianceEventReceiver

	err = json.Unmarshal(output, &eventReceiver)

	if err != nil {
		logp.Err("JSON error: %v", err)
	}

	ce, err := eventReceiver.ToComplianceEvents()

	if err != nil {
		logp.Err("Conversion to compliance events array failed: %v", err)
	}

	for _, complianceEvent := range ce {
		checkEvent := beat.Event{
			Timestamp: complianceEvent.EventTimestamp,
			Fields: common.MapStr{
				"compliance.category":          complianceCheck.category,
				"compliance.name":              complianceEvent.ComplianceCheckFunction,
				"compliance.status":            complianceEvent.EventStatus,
				"compliance.eventid":           complianceEvent.EventID,
				"compliance.resultlevel":       complianceEvent.ResultLevel,
				"compliance.resultlevelnumber": complianceEvent.IntResultLevel,
				"compliance.messagetext":       complianceEvent.ComplianceCheckResultMessage,
			},
		}

		events = append(events, checkEvent)
	}
	return

}
