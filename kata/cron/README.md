## cron service

This program is a simple cron service that runs a set of commands at specified intervals.

### Configuration

A cron expression represents a set of times, using 5 space-separated fields.

Field name   | Mandatory? | Allowed values  | Allowed special characters
-------------|------------|-----------------|----------------------------
Seconds      | Yes        | 0-59            | * / , -
Minutes      | Yes        | 0-59            | * / , -
Hours        | Yes        | 0-23            | * / , -
Day of month | Yes        | 1-31            | * / , - ?
Month        | Yes        | 1-12 or JAN-DEC | * / , -
Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?


```
tasks:
  - name: "Task 1"
    schedule: "*/5 * * * * *"  # Run every 5 seconds
    command: "echo 'Hello, Task 1!'"
  - name: "Task 2"
    deadline: "2025-04-10T14:30:00+08:00"
    command: "echo 'Hello, Task 2!'"
```

### Install Dependencies

You need to install the required Go packages:

```bash
go mod init cron-service
go get github.com/robfig/cron/v3
go get gopkg.in/yaml.v2
```

### Build the Program

Build the Go program:

```bash
go build -o cron-service
```

### Run the Program

Run the compiled binary:

```bash
./cron-service
```

### Run as a macOS Service

To run this as a macOS service, you can create a `plist` file and use `launchctl` to manage it.

1. Create a `com.fanyamin.cron-service.plist` file in `~/Library/LaunchAgents/`:

```xml

<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
  <dict>
    <key>Label</key>
    <string>com.fanyamin.cron-service</string>
    <key>ProgramArguments</key>
    <array>
      <string>~/cron-service</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardErrorPath</key>
    <string>/tmp/cron-service.err</string>
    <key>StandardOutPath</key>
    <string>/tmp/cron-service.out</string>
  </dict>
</plist>


```

2. Load the service:

```bash
launchctl load ~/Library/LaunchAgents/com.fanyamin.cron-service.plist
```

3. Start the service:

```bash
launchctl start com.fanyamin.cron-service
```

### Verify the Service

You can check the logs to verify that the service is running and executing the jobs as expected:

```bash
tail -f /var/log/system.log
```
