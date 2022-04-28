# cveonslack
Automated script for updated CVE pushed on slack.

To install

```
git clone github.com/th3r4id/cveonslack.git
cd cveonslack
```
Update slack tokens and channel-id in code\
```
go build cvesearch.go
```

This will create a binary file of cvesearch add this in crontab to run it frequently.

To run this in every 5 min create cron job using below command
```
crontab -e
5 * * * * * /cveonslack
```

