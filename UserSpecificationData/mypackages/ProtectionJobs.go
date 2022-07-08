package mypackages

import (
	"encoding/json"
	"strconv"
)

var ProtectionJobsList []string 
var ProtectionJobKeys = make(map[string]string)
var ProtectiongroupEnvironment = make(map[string]string)

var virtual_machine_env = [7]string{"kAcropolis","kAWS","kAzure","kGCP","kHyperV","kKVM","kVMware"}
var data_base_env = [6]string{"kCassandra","kCouchbase","kMongoDB","kOracle","kSQL","kUDA"}
var NAS_env = [6]string{"kElastifile","kFlashBlade","kGenericNas","kGPFS","kIsilon","kNetapp"}

type metrics struct {
	Runs      []run `json:"runs"`
	TotalRuns int   `json:"totalRuns"`
}
type run struct {
	Id              string          `json:"id"`
	Environment 	string 			`json:"environment"`
	BackupInfo      jobinfo         `json:"localBackupInfo"`
	ReplicationInfo replicationInfo `json:"replicationInfo"`
	ArchivalInfo    archivalInfo    `json:"archivalInfo"`
}

type jobinfo struct {
	StartTimeUsecs  int64  `json:"startTimeUsecs"`
	EndTimeUsecs    int64  `json:"endTimeUsecs"`
	QueuedTimeUsecs int64  `json:"queuedTimeUsecs"`
	Status          string `json:"status"`
	IsSlaViolated   bool   `json:"isSlaViolated"`
}

type replicationInfo struct {
	TargetResults []jobinfo `json:"replicationTargetResults"`
}
type archivalInfo struct {
	TargetResults []jobinfo `json:"archivalTargetResults"`
}
type jobs struct {
	ProtectionJobs []protectionjob `json:"protectionGroups"`
}

type protectionjob struct {
	JobName string `json:"name"`
	Id      string `json:"id"`
}

func ProtectionJobInfo(jobname string) (job_environment string,backup_count int,sum_of_jobback_times int64,average_jobbackup_time string,average_jobbackup_time_Usecs int64 ,
	 replication_count int,sum_of_jobreplication_times int64,average_jobreplication_time string, average_jobreplication_time_Usecs int64,
	 archival_count int,sum_of_jobarchival_times int64,average_jobarchival_time string, average_jobarchival_time_Usecs int64 ,sla_times int, runs int) {
	
	jobid := ProtectionJobKeys[jobname]
	timePeriod := "day"
	starttime, endtime := TimeStampGneerator(timePeriod, "U")

	newUrl := GenerateNewURLforProtectionJobs(starttime, endtime, jobid)

	response := PostRequestForAccessToken()

	data := GetRequestForJsonData(response, newUrl)

	var responsedata metrics

	json.Unmarshal(data, &responsedata)

	job_environment = JobEnvironment(responsedata)

	backup_count, sum_of_jobback_times, average_jobbackup_time, average_jobbackup_time_Usecs = JobBackup(responsedata)
	replication_count, sum_of_jobreplication_times, average_jobreplication_time, average_jobreplication_time_Usecs = JobReplication(responsedata)
	archival_count, sum_of_jobarchival_times, average_jobarchival_time, average_jobarchival_time_Usecs = JobArchival(responsedata)
	sla_times = SlaTimes(responsedata)
	runs = responsedata.TotalRuns

	return
}

func FillProtectionJobKeys() {

	ProtectionJobsList = nil

	access_token := PostRequestForAccessToken()

	url := "https://10.14.19.226/v2/data-protect/protection-groups?useCachedData=false&pruneSourceIds=true&isDeleted=false&includeTenants=true&includeLastRunInfo=true"

	jsondata := GetRequestForJsonData(access_token, url)

	var responsejsondata jobs

	json.Unmarshal(jsondata, &responsejsondata)

	l := len(responsejsondata.ProtectionJobs)

	for i := 0; i < l; i++ {
		ProtectionJobsList = append(ProtectionJobsList, responsejsondata.ProtectionJobs[i].JobName)
		id := responsejsondata.ProtectionJobs[i].Id
		key := GenerateProtectionJobKeys(id)
		ProtectionJobKeys[ProtectionJobsList[i]] = key
	}
	
	for i:=0;i<7;i++{
		ProtectiongroupEnvironment[virtual_machine_env[i]] = "VirtualMachines"
	}        
	
	for i:=0;i<6;i++{
		ProtectiongroupEnvironment[data_base_env[i]] = "Databases"
	}        

	for i:=0;i<6;i++{
		ProtectiongroupEnvironment[NAS_env[i]] = "NAS"
	}      

	ProtectiongroupEnvironment["kO365"] = "Microsoft365"
	ProtectiongroupEnvironment["kPhysical"] = "PhysicalServers"
	ProtectiongroupEnvironment["kAD"] = "Applications"
	ProtectiongroupEnvironment["kExchange"] = "Applications"
	ProtectiongroupEnvironment["kPure"] = "SAN"
	ProtectiongroupEnvironment["kView"] = "CohesityViews"
	
}

func GenerateProtectionJobKeys(id string) (result string) {

	last_digit_in_id := len(id) - 1
	var key string

	for id[last_digit_in_id] != 58 {
		key = key + string(id[last_digit_in_id])
		last_digit_in_id -= 1
	}

	for _, v := range key {
		result = string(v) + result
	}
	return
}

func JobEnvironment(responsedata metrics)(job_environment string){
	if 0 < len(responsedata.Runs){
	job_environment = responsedata.Runs[0].Environment
	}
	return
}

func JobBackup(responsedata metrics) (backup_count int,sum_time int64,average_jobbackup_time string, average_jobbackup_time_Usecs int64) {
	var start_time_vector []int64
	var end_time_vector []int64

	for i := 0; i < len(responsedata.Runs); i++ {

		status := responsedata.Runs[i].BackupInfo.Status

		if status == "Succeeded" || status == "SucceededWithWarning" || status == "Failed" {

			start_time_vector = append(start_time_vector, responsedata.Runs[i].BackupInfo.StartTimeUsecs)
			end_time_vector = append(end_time_vector, responsedata.Runs[i].BackupInfo.EndTimeUsecs)

		} else if status == "Running" {
	}
}
size := len(start_time_vector)
sum_time, average_time := SumAndAverageTime(size, start_time_vector, end_time_vector)
converted_time := ConvertUnixTime(average_time)
return size, sum_time,converted_time , average_time
}

func JobReplication(responsedata metrics) (replication_count int,sum_time int64,average_job_replication_time string , average_job_replication_time_Usecs int64) {
	var start_time_vector []int64
	var end_time_vector []int64

	for i := 0; i < len(responsedata.Runs); i++ {

		if 0 < len(responsedata.Runs[i].ReplicationInfo.TargetResults) {
			for j := 0; j < (len(responsedata.Runs[i].ReplicationInfo.TargetResults)); j++ {

				status := responsedata.Runs[i].ReplicationInfo.TargetResults[j].Status

				if status == "Succeeded" || status == "SucceededWithWarning" || status == "Failed" {

					start_time_vector = append(start_time_vector, responsedata.Runs[i].ReplicationInfo.TargetResults[j].StartTimeUsecs)
					end_time_vector = append(end_time_vector, responsedata.Runs[i].ReplicationInfo.TargetResults[j].EndTimeUsecs)
				
					} else if status == "Running" {}
			}
		} else {}
	}
	size := len(start_time_vector)
	sum_time, average_time := SumAndAverageTime(size, start_time_vector, end_time_vector)
	converted_time := ConvertUnixTime(average_time)
	return size,sum_time,converted_time , average_time
}

func JobArchival(responsedata metrics) (archival_count int,sum_time int64, average_jobarchival_time string, average_jobarchival_time_Usecs int64) {

	var start_time_vector []int64
	var end_time_vector []int64

	for i := 0; i < len(responsedata.Runs); i++ {

		if 0 < len(responsedata.Runs[i].ArchivalInfo.TargetResults) {

			for j := 0; j < (len(responsedata.Runs[i].ArchivalInfo.TargetResults)); j++ {

				status := responsedata.Runs[i].ArchivalInfo.TargetResults[j].Status
	
				if status == "Succeeded" || status == "SucceededWithWarning" {

					start_time_vector = append(start_time_vector, responsedata.Runs[i].ArchivalInfo.TargetResults[j].StartTimeUsecs)
					end_time_vector = append(end_time_vector, responsedata.Runs[i].ArchivalInfo.TargetResults[j].EndTimeUsecs)

				} else if status == "Running" {

				} else if status == "Failed" {

				}
			}
		} else {}
	}

	size := len(start_time_vector)
	sum_time, average_time := SumAndAverageTime(size, start_time_vector, end_time_vector)
	converted_time := ConvertUnixTime(average_time)
	return size, sum_time,converted_time , average_time
}

func SlaTimes(responsedata metrics) (slacount int) {

	slacount = 0

	for i := 0; i < len(responsedata.Runs); i++ {

		slaViolation := responsedata.Runs[i].BackupInfo.IsSlaViolated

		if !slaViolation {
			slacount++
		}
	}
	return slacount
}

func SumAndAverageTime(size_of_vector int, start_time_vector []int64, end_time_vector []int64) (sum int64, avg int64) {

	for i := 0; i < size_of_vector; i++ {
		sum += (end_time_vector[i] - start_time_vector[i])
	}
	if size_of_vector > 0{
	avg = sum / int64(size_of_vector)
	}else{
		avg = 0
	}

	return
}

func ConvertUnixTime(initialtime int64) (finaltime string) {

	initialtime /= 1000000
	hrs := initialtime / 3600
	initialtime = initialtime - hrs*3600
	minutes := initialtime / 60
	seconds := initialtime - minutes*60

	hrsstr := strconv.FormatInt(hrs, 10)
	minstr := strconv.FormatInt(minutes, 10)
	secstr := strconv.FormatInt(seconds, 10)

	finaltime = hrsstr + "hr " + minstr + "min " + secstr + "sec"
	return
}
