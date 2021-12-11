# TEMPO LOGGER

The "Last Chance" Tool. Allows you to automatically register your working time for some tasks from the template, if you
were too busy and forgot about it

## Installation

1. Download repo and open tempo-logger root. Build it:

```cmd
go build
```

2. Rename .env.default -> .env
3. Prepare .env:
    1. JIRA_URL - jira root url. <span style="color: green">For example https://jira.organization.com </span>
    2. LOGIN, PASSWORD - jira account for time logging
    3. WORKLOG - an array of jira worklogs. <span style="color: green">For example: RFS-123,5,code review; INC-2,3,bug
       fixing</span>
        1. <jira_ticket> - Required field. Ticket from jira.
        2. <max_hours_spent> - Required field. Hours limit (int val) for logging. The Tool will generate a random
           int-value for logging [1..N], where N = <max_hours_spent>. Note: if it is last task in the array and <
           max_hours_spent> < remainder time, the Tool will spend the entire remainder on that task.
        3. <worklog_comment> - Required field. Comment for logged time

## Usage

Run it manually or create scheduler to run automatically. For example, replace ........ and import xml-config for
Windows Task Scheduler:

``` xml
<?xml version="1.0" encoding="UTF-16"?>
<Task version="1.2" xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
  <RegistrationInfo>
    <Date>2021-12-10T09:50:22.5002972</Date>
    <Author>........MY_DOMAIN\USER_LOGIN</Author>
    <URI>\Tempo Logger</URI>
  </RegistrationInfo>
  <Triggers>
    <CalendarTrigger>
      <StartBoundary>........2021-12-10T10:30:00</StartBoundary>
      <Enabled>true</Enabled>
      <ScheduleByDay>
        <DaysInterval>1</DaysInterval>
      </ScheduleByDay>
    </CalendarTrigger>
  </Triggers>
  <Settings>
    <MultipleInstancesPolicy>IgnoreNew</MultipleInstancesPolicy>
    <DisallowStartIfOnBatteries>false</DisallowStartIfOnBatteries>
    <StopIfGoingOnBatteries>true</StopIfGoingOnBatteries>
    <AllowHardTerminate>true</AllowHardTerminate>
    <StartWhenAvailable>false</StartWhenAvailable>
    <RunOnlyIfNetworkAvailable>false</RunOnlyIfNetworkAvailable>
    <IdleSettings>
      <Duration>PT10M</Duration>
      <WaitTimeout>PT1H</WaitTimeout>
      <StopOnIdleEnd>true</StopOnIdleEnd>
      <RestartOnIdle>false</RestartOnIdle>
    </IdleSettings>
    <AllowStartOnDemand>true</AllowStartOnDemand>
    <Enabled>true</Enabled>
    <Hidden>false</Hidden>
    <RunOnlyIfIdle>false</RunOnlyIfIdle>
    <WakeToRun>true</WakeToRun>
    <ExecutionTimeLimit>PT72H</ExecutionTimeLimit>
    <Priority>7</Priority>
  </Settings>
  <Actions Context="Author">
    <Exec>
      <Command>........C:\tempo-logger\logger.exe</Command>
      <WorkingDirectory>........C:\tempo-logger</WorkingDirectory>
    </Exec>
  </Actions>
</Task>
```

Enjoy it
