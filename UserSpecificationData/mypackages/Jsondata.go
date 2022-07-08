package mypackages

import (
	"encoding/json"
	"fmt"
)

type jsondata struct {
	Id                              string            `json:"id"`
	ClusterSoftwareVersion 			string 			  `json:"clusterSoftwareVersion"`
	TimeStamp                       int64             `json:"timeStamp"`
	FileCreateRate                  int64             `json:"fileCreateRate"`
	FileCreateSum                   int64             `json:"fileCreateSum"`
	FileLatency                     float64           `json:"filerLatency"`
	FileLatencyUnits                string            `json:"fileLatencyUnits"`
	SystemUtilizationChangeRate     float64           `json:"systemUtilizationChangeRate"`
	SystemUtilizationChangeRateUnit string            `json:"systemUtilizationChangeRateUnit"`
	GarbageCollection               float64           `json:"garbageCollection"`
	GarbageCollectionUnit           string            `json:"garbageCollectionUnit"`
	ClusterSpaceUtilization         float64           `json:"clusterSpaceUtilization"`
	ClusterSpaceUtilizationUnits    string            `json:"clusterSpaceUtilizationUnits"`
	MetaDataUtilization             float64           `json:"metaDataUtilization"`
	MetaDataUtilizationUnits        string            `json:"metaDataUtilizationUnits"`
	MetaDataSpacePercentage         float64           `json:"meta_data_space_percentage"`
	Deduplication					float64			  `json:"deduplication"`
	IndexingBacklog                 int64             `json:"indexingBacklog"`
	ProtectionJobsInfo              []protectiongroup `json:"protectionGroups"`
	AverageJobbackupTimeForAllJobsUsecs int64		  `json:"averageJobbackupTimeForAllJobsUsecs"`
	AverageJobreplicationTimeForAllJobsUsecs int64    `json:"averageJobreplicationTimeForAllJobsUsecs"`
	AverageJobarchivalTimeForAllJobsUsecs int64		  `json:"averageJobarchivalTimeForAllJobsUsecs"`
	TotalSlaviolations				int				  `json:"totalSlaViolations"`	
	TotalRuns 						int 			  `json:"totalRuns"`
	Alerts                          []alert           `json:"alerts"`
}

type protectiongroup struct {
	ProtectionGroupName 		string 	 `json:"protectionGroupName"`
	ProtectionGroupType 		string 	 `json:"protectionGroupType"`
	JobBackuptime       		string	 `json:"jobBackupTime"`
	JobBackuptimeUsecs  		int64 	 `json:"jobBackupTimeUsecs"`
	JobReplicationTime  		string 	 `json:"jobReplicationTime"`
	JobReplicationTimeUsecs  	int64 	 `json:"jobReplicationTimeUsecs"`
	JobArchivalTime     		string	 `json:"jobArchivalTime"`
	JobArchivalTimeUsecs    	int64	 `json:"jobArchivalTimeUsecs"`
	TimesSlaNotMet     		 	int   	 `json:"timesSlaNotMet"`
	Runs                		int      `json:"runs"`
}

var Response_for_elastic jsondata
var Present_time_stamp int64

func GenerateJson(time_stamp int64) (data []byte){
	var total_backup_time, total_archival_time, total_replication_time int64
	var average_jobbackup_time_Usecs, average_jobarchival_time_Usecs, average_jobreplication_time_Usecs int64
	total_backup_time, total_archival_time, total_replication_time = 0,0,0 
	total_backup_count,total_replication_count,total_archival_count,total_sla_count, total_runs := 0,0,0,0,0
	 
	Present_time_stamp = time_stamp 
	fmt.Println("TimeStamp : ", Present_time_stamp)
	FillProtectionJobKeys()
	total_jobs := len(ProtectionJobsList)
	var myProtectionGroups []protectiongroup
	for i := 0; i < total_jobs; i++ {
		job_environment, backup_count,job_backup_time_sum,job_backup_time, job_backup_time_Usecs, replication_count,job_replication_time_sum, job_replication_time, job_replication_time_Usecs, archival_count,job_archival_time_sum, job_archival_time, job_archival_time_Usecs,times_sla_not_met, runs := ProtectionJobInfo(ProtectionJobsList[i])
		mygroup := protectiongroup{
			ProtectionGroupName: ProtectionJobsList[i],
			ProtectionGroupType: ProtectiongroupEnvironment[job_environment],
			JobBackuptime:       job_backup_time,
			JobBackuptimeUsecs: job_backup_time_Usecs,
			JobReplicationTime:  job_replication_time,
			JobReplicationTimeUsecs: job_replication_time_Usecs,
			JobArchivalTime:     job_archival_time,
			JobArchivalTimeUsecs: job_archival_time_Usecs,
			TimesSlaNotMet:      times_sla_not_met,
			Runs:                runs,
		}
		
		total_sla_count += times_sla_not_met
		total_backup_count += backup_count
		total_backup_time += job_backup_time_sum
		total_archival_count += archival_count
		total_archival_time += job_archival_time_sum
		total_replication_count += replication_count
		total_replication_time += job_replication_time_sum
		total_runs += runs

		myProtectionGroups = append(myProtectionGroups, mygroup)
	}

	if(total_backup_count > 0){
		average_jobbackup_time_Usecs = total_backup_time/int64(total_backup_count)
	}else{ average_jobbackup_time_Usecs = 0}

	if total_archival_count >0 {
		average_jobarchival_time_Usecs = total_archival_time/int64(total_archival_count)
	}else{ average_jobarchival_time_Usecs = 0 }

	if total_replication_count>0{
		average_jobreplication_time_Usecs = total_replication_time/int64(total_replication_count)
	}else{ average_jobreplication_time_Usecs = 0 }

	file_create_rate := Filecreaterate()
	file_create_sum := FileCreateSum()
	file_latency := FileLatency()
	system_utilization_change_rate := SystemUtilizationChangeRate()
	garbage_collection := GarbageCollection()
	cluster_space_utilization := Cluster_Space_Utilization
	cluster_version := ClusterVersion()
	indexing_backlog := IndexingBacklog()
	meta_data_utilization := Meta_Data_Utilization
	meta_data_space_percentage := Meta_Data_Space_Percentage
	deduplication := Deduplication()
	Alerts()
	myjsondata := jsondata{
		Id:                              "longivitycluster",
		ClusterSoftwareVersion: 		 cluster_version,
		TimeStamp:                       Present_time_stamp,
		FileCreateRate:                  file_create_rate,
		FileCreateSum:     		 	   	 file_create_sum,
		FileLatency:                     file_latency,
		FileLatencyUnits:                "ms",
		SystemUtilizationChangeRate:     system_utilization_change_rate,
		SystemUtilizationChangeRateUnit: "GiB/sec",
		GarbageCollection:               garbage_collection,
		GarbageCollectionUnit:           "GiB",
		ClusterSpaceUtilization:         cluster_space_utilization,
		ClusterSpaceUtilizationUnits:    "TiB",
		MetaDataUtilization:             meta_data_utilization,
		MetaDataUtilizationUnits:        "TiB",
		MetaDataSpacePercentage:         meta_data_space_percentage,
		Deduplication:  			     deduplication,				
		IndexingBacklog:                 indexing_backlog,
		ProtectionJobsInfo:              myProtectionGroups,
		AverageJobbackupTimeForAllJobsUsecs: average_jobbackup_time_Usecs,
		AverageJobreplicationTimeForAllJobsUsecs: average_jobreplication_time_Usecs,
		AverageJobarchivalTimeForAllJobsUsecs: average_jobarchival_time_Usecs,
		TotalRuns: 						 total_runs,		
		TotalSlaviolations: 			 total_sla_count,			
		Alerts:                          AlertsList,
	}

	data, _ = json.Marshal(myjsondata)
	return data
}
