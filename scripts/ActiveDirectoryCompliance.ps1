#
# ActiveDirectory.ps1
#
Param(
    $Name
)

Import-Module Compliance


$result = New-ComplianceResult
$checkOkMessage = "Passed check."

function Check-CurrentAccountForKerberos
{
	$currentUser = [System.Security.Principal.WindowsIdentity]::GetCurrent();

	if ($currentUser.AuthenticationType -ne "Kerberos")
	{
		$messageText = $currentUser.Name + " is not authenticated via Kerberos."
		$msg = New-ComplianceMessage -Title $MyInvocation.MyCommand -Message $messageText -Level CRITICAL
	}
	else
	{
		$msg = New-ComplianceMessage -Title $MyInvocation.MyCommand -Message $checkOkMessage -Level OK
	}

	$result.AddComplianceMessage($msg)
}

function Check-ADDomainConnectionAndTimeSkew
{
	# timeskew in minutes
	$TIMESKEW = 5.0
	try
	{
		$domainInfo = [System.DirectoryServices.ActiveDirectory.Domain]::GetComputerDomain()
		$localDateTime = Get-Date

		$timespan = New-TimeSpan -Start $domainInfo.DomainControllers[0].CurrentTime -End $localDateTime.ToUniversalTime()

		if ($timespan.TotalMinutes -gt $TIMESKEW)
		{
			$messageText = "Time skew exceeds " + $TIMESKEW  + " min: " + $timespan.TotalMinutes
			$msg = New-ComplianceMessage -Title $MyInvocation.MyCommand -Message $messageText -Level WARNING
		}
		else
		{
			$msg = New-ComplianceMessage -Title $MyInvocation.MyCommand -Message $checkOkMessage -Level OK
		}
	}
	catch [System.DirectoryServices.ActiveDirectory.ActiveDirectoryObjectNotFoundException]
	{
		$msg = New-ComplianceMessage -Title $MyInvocation.MyCommand -Message $_.Exception.Message -Level CRITICAL
	}

	$result.AddComplianceMessage($msg)
}

Check-CurrentAccountForKerberos
Check-ADDomainConnectionAndTimeSkew

$result.SetOverallStatusLevel()
$result.PrintJsonToConsole()
