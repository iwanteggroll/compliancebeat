# Compliancebeat

## Description

Compliancebeat is a custom elastic beat designed to check services, configurations, connectivity, etc. and send the status to an elasticsearch server with the results viewable in kibana. Modeled after Nagios, Compliancebeat is similar but it's designed to run locally by calling python3 or powershell scripts to perform the checks and return a status level along with an error message. With the power of elasticsearch and kibana, the raw data can be transformed and visualized into various dashboards.

Compliancebeat was borne out of a need to verify and report Active Directory compliance during blue/red team events. During these events, frustrated red teams normally try to leverage credential reuse within an Active Directory enclave but discover an environment is not properly setup even though Active Directory is a requirement for blue teams to implement. Active Directory checks are the initial use case but compliancebeat can support many more checks. If the check can be written in python3 or powershell, then compliancebeat can provide a framework to format and send the data into a elasticsearch/kibana instance.

## Supported Operating Systems
* Linux w/ systemd and python3
* Microsoft Windows 8/2012 and higher

## Installing the Packages
* Windows
  * Open an administrative powershell prompt
  * Run `.\installer.ps1 -elasticsearch_ip IP -team STRING` (no spaces for the team parameter)
* Linux (TBD)

## Writing a Check

Under the `scripts` directory, are 2 examples of how to write a check, one in powershell (ActiveDirectoryCompliance.ps1) and the other in python (TestCheck.py). The general structure is to instantiante a result object globally, perform a check within a function, adding a message to the result within the function, calling the function to execute, tabulate the results, and print the json out so compliancebeat can read it and send it off to elasticsearch.

