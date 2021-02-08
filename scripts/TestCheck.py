#!/usr/bin/env python3

import pycompliance

OK_MESSAGE = "Passed check."

result = pycompliance.ComplianceResult()

def TestFunc1():
	message = pycompliance.ComplianceMessage("TestFunc1", OK_MESSAGE, pycompliance.StatusLevel(pycompliance.StatusLevelEnum.OK))
	result.AddComplianceMessage(message)


def TestFunc2():
	message = pycompliance.ComplianceMessage("TestFunc2", "Failed check", pycompliance.StatusLevel(pycompliance.StatusLevelEnum.CRITICAL))
	result.AddComplianceMessage(message)


TestFunc1()
TestFunc2()

result.SetOverallStatusLevel()

result.PrintJsonToConsole()
