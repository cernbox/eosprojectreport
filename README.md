### Download

You can grab the binary for linux amd64 from the Releases page.

### Requirements
The machine where the tool runs needs to have EOS root access to the MGM.
So either you run it on the mgm node or on a trusted gateway like the cbox-webng-01.

### Usage

```
[labkode@labradorbox eosprojectreport]$ eosprojectreport -h
Usage of eosprojectreport:
  -filled-less int
        returns projects with usage less that <n>%
  -filled-more int
        returns projects with usage bigger that <n>%
  -groupby string
        aggreate by time dimension (day,month,year,one)
  -human
        output human readable values
  -mgm string
        mgm url where projects live (default "root://eosuser-slave.cern.ch")
  -newer int
        returns projects newer than <n> days
  -older int
        returns projects older than <n> days
  -only-empty
        show only empty projects
  -s    remove header from output
  -sep string
        separator to use in output (default " ")
```

List project space information:

```
$ eosprojectreport
#UID PATH CTIME TOTAL USED USAGE
89737 /eos/project/a/abpdata/ 1452524194 60000000000000 36572094146486 60.95%
99113 /eos/project/a/abtua9/ 1480697750 2500000000000 1049321964 0.04%
```

List project space information with human readable output, (-s for silent):
```
$ eosprojectreport --human -s
abpdata /eos/project/a/abpdata/ 2016/01/11 60TB 37TB 60.95%
abtua9 /eos/project/a/abtua9/ 2016/12/02 2.5TB 1.0GB 0.04%
```

Find projects created in the last 30 days:

```
[labkode@labradorbox eosprojectreport]$ eosprojectreport --human --newer 30
#UID PATH CTIME TOTAL USED USAGE
rmunzers /eos/project/a/alice-tpc-upgrade/ 2018/11/29 2.0TB 0B 0.00%
atlsuep /eos/project/a/atlas-suep/ 2018/10/31 1.0TB 74GB 7.44%
```

Find projects older than 1 year:

```
$ eosprojectreport --human --older 365 -s
abpdata /eos/project/a/abpdata/ 2016/01/11 60TB 37TB 60.95%
abtua9 /eos/project/a/abtua9/ 2016/12/02 2.5TB 1.0GB 0.04%
```


Find empty projects:

```
$ eosprojectreport --human -s -only-empty
ad0 /eos/project/a/ad-alice/ 2017/10/11 1.0TB 0B 0.00%
aliglan /eos/project/a/alice-sams/ 2017/09/19 1.0TB 0B 0.00%
```

Find projects with usage bigger than 80%:

```
$ eosprojectreport --human -s --filled-more 80
aliceits /eos/project/a/alice-its/ 2016/09/06 80TB 65TB 81.48%
```

Find projects with usage less than 10%:

```
$ eosprojectreport --human -s --filled-less 10
abtua9 /eos/project/a/abtua9/ 2016/12/02 2.5TB 1.0GB 0.04%
adadt /eos/project/a/ad-adt/ 2018/04/26 1.0TB 396B 0.00%
```

Aggregate counters to one single value:

```
$ eosprojectreport --human --groupby one
#GROUPBY COUNT TOTAL USED USAGE
one 322 1.1PB 272TB 0.26%
```


Aggreate counters by year:

```
$ eosprojectreport --human --groupby year | sort -n
#GROUPBY COUNT TOTAL USED USAGE
2015 1 1.0TB 528GB 0.53%
2016 50 372TB 147TB 0.39%
2017 148 372TB 96TB 0.26%
2018 123 318TB 29TB 0.09%
```

Set output to CSV format:

```
$ eosprojectreport  -sep ,
#UID,PATH,CTIME,TOTAL,USED,USAGE
89737,/eos/project/a/abpdata/,1452524194,60000000000000,36572341004069,60.95%
99113,/eos/project/a/abtua9/,1480697750,2500000000000,1049321964,0.04%
99103,/eos/project/a/active_halo_collimation/,1480625039,2500000000000,799933255404,32.00%
```
