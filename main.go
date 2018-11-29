package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"regexp"
	"strings"
	"time"

	"github.com/cernbox/reva/api/storage_eos/eosclient"
	"github.com/dustin/go-humanize"
)

var (
	newerFlag      int64
	olderFlag      int64
	humanFlag      bool
	silentFlag     bool
	sepFlag        string
	urlFlag        string
	groupByFlag    string
	emptyOnlyFlag  bool
	filledMoreFlag int64
	filledLessFlag int64

	matchProject = regexp.MustCompile(`/eos/project/[a-z]/`)
)

const (
	groupByDay   = "day"
	groupByMonth = "month"
	groupByYear  = "year"
	groupByOne   = "one"
)

func init() {
	flag.Int64Var(&newerFlag, "newer", 0, "returns projects newer than <n> days")
	flag.Int64Var(&olderFlag, "older", 0, "returns projects older than <n> days")
	flag.BoolVar(&humanFlag, "human", false, "output human readable values")
	flag.BoolVar(&silentFlag, "s", false, "remove header from output")
	flag.BoolVar(&emptyOnlyFlag, "only-empty", false, "show only empty projects")
	flag.StringVar(&sepFlag, "sep", " ", "separator to use in output")
	flag.StringVar(&urlFlag, "mgm", "root://eosuser-slave.cern.ch", "mgm url where projects live")
	flag.StringVar(&groupByFlag, "groupby", "", "aggreate by time dimension (day,month,year,one)")
	flag.Int64Var(&filledMoreFlag, "filled-more", 0, "returns projects with usage bigger that <n>%")
	flag.Int64Var(&filledLessFlag, "filled-less", 0, "returns projects with usage less that <n>%")
	flag.Parse()
}

type info struct {
	md                                *eosclient.FileInfo
	total, used                       int
	cTimeHuman, totalHuman, usedHuman string
	uidHuman                          string
	usage                             float64
}

func (i *info) addHuman() {
	i.cTimeHuman = time.Unix(int64(i.md.CTime), 0).Format("2006/01/02")
	i.totalHuman = strings.Replace(humanize.Bytes(uint64(i.total)), " ", "", -1)
	i.usedHuman = strings.Replace(humanize.Bytes(uint64(i.used)), " ", "", -1)
	u, err := user.LookupId(i.md.UID)
	if err == nil {
		i.uidHuman = u.Username
	}
}

func main() {

	ctx := context.Background()
	opts := &eosclient.Options{URL: urlFlag}
	client, err := eosclient.New(opts)
	if err != nil {
		log.Fatal(err)
	}

	mds, err := client.List(ctx, "root", "/eos/project")
	if err != nil {
		log.Fatal(err)
	}

	infos := []*info{}
	if groupByFlag == "" {
		printHeader()
	}

	for _, md := range mds {
		if matchProject.MatchString(md.File) {
			projmds, err := client.List(ctx, "root", md.File)
			if err != nil {
				log.Fatal(err)
			}

			for _, projmd := range projmds {
				total, used, err := client.GetQuota(ctx, projmd.UID, projmd.File)
				if err != nil {
					log.Fatalln(err)
				}

				info := &info{
					md:    projmd,
					total: total,
					used:  used,
					usage: (float64(used) / float64(total)) * 100,
				}
				infos = append(infos, info)
				if groupByFlag == "" {
					process(info)
				}
			}
		}
	}

	if groupByFlag == "" {
		os.Exit(0)
	}

	// groupBy is set
	buckets := map[string][]*info{}
	for _, i := range infos {
		var key string
		if groupByFlag == groupByDay {
			key = time.Unix(int64(i.md.CTime), 0).Format("2006/01/02")
		} else if groupByFlag == groupByMonth {
			key = time.Unix(int64(i.md.CTime), 0).Format("2006/01")
		} else if groupByFlag == groupByYear {
			key = time.Unix(int64(i.md.CTime), 0).Format("2006")
		} else if groupByFlag == groupByOne {
			key = groupByFlag
		} else {
			log.Fatal("unkown value for groupByFlag: ", groupByFlag)
		}

		val := buckets[key]
		if val == nil {
			buckets[key] = []*info{i}
		} else {
			buckets[key] = append(buckets[key], i)
		}
	}

	printGroupByHeader()
	for key, infos := range buckets {
		var total, used int
		for _, i := range infos {
			total += i.total
			used += i.used
		}

		usage := float64(used) / float64(total)

		fields := []string{}
		if humanFlag {
			totalHuman := strings.Replace(humanize.Bytes(uint64(total)), " ", "", -1)
			usedHuman := strings.Replace(humanize.Bytes(uint64(used)), " ", "", -1)
			fields = append(fields, key, fmt.Sprintf("%d", len(infos)), totalHuman, usedHuman, fmt.Sprintf("%.2f", usage)+"%")
		} else {
			fields = append(fields, key, fmt.Sprintf("%d", len(infos)), fmt.Sprintf("%d", total), fmt.Sprintf("%d", used), fmt.Sprintf("%.2f", usage)+"%")
		}
		fmt.Println(strings.Join(fields, sepFlag))
	}
}

func printInfo(i *info) {
	fields := []string{}
	if humanFlag {
		i.addHuman()
		fields = append(fields, i.uidHuman, i.md.File, i.cTimeHuman, i.totalHuman, i.usedHuman, fmt.Sprintf("%.2f", i.usage)+"%")
	} else {
		fields = append(fields, i.md.UID, i.md.File, fmt.Sprintf("%d", i.md.CTime), fmt.Sprintf("%d", i.total), fmt.Sprintf("%d", i.used), fmt.Sprintf("%.2f", i.usage)+"%")
	}
	fmt.Println(strings.Join(fields, sepFlag))
}

func process(i *info) {
	if emptyOnlyFlag {
		if i.used == 0 {
			printInfo(i)
		}
		return
	}

	if newerFlag > int64(0) {
		from := time.Now().Add(-(time.Minute * 24 * 60 * time.Duration(newerFlag))).Unix()
		if i.md.CTime > uint64(from) {
			printInfo(i)
		}
		return
	}

	if olderFlag > int64(0) {
		from := time.Now().Add(-(time.Minute * 24 * 60 * time.Duration(olderFlag))).Unix()
		if i.md.CTime < uint64(from) {
			printInfo(i)
		}
		return
	}

	if filledMoreFlag > 0 {
		if i.usage > float64(filledMoreFlag) {
			printInfo(i)
		}
		return
	}

	if filledLessFlag > 0 {
		if i.usage < float64(filledLessFlag) {
			printInfo(i)
		}
		return
	}

	printInfo(i)
}

func printHeader() {
	if !silentFlag {
		header := []string{"#UID", "PATH", "CTIME", "TOTAL", "USED", "USAGE"}
		fmt.Println(strings.Join(header, sepFlag))
	}
}

func printGroupByHeader() {
	if !silentFlag {
		header := []string{"#GROUPBY", "COUNT", "TOTAL", "USED", "USAGE"}
		fmt.Println(strings.Join(header, sepFlag))
	}
}
