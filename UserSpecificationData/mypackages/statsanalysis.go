package mypackages

import (
	"encoding/json"
	"time"
	//	"mypackages/generator"

	"fmt"
	"math"
	"sort"
	"strconv"
)

type datastats struct {
	Id                               int64                  `json:"id"`
	ClusterSoftwareVersion           string                 `json:"clusterSoftwareVersion"`
	FileCreateRateStats              values                 `json:"fileCreateRateStats"`
	FileCreateSumStats               values                 `json:"fileCreateSumStats"`
	FileLatencyStats                 values                 `json:"fileLatencyStats"`
	SystemUtilizationChangeRateStats values                 `json:"systemUtilizationChangeRateStats"`
	GarbageCollectionStats           values                 `json:"garbageCollectionStats"`
	ClusterSpaceUtilizationStats     values                 `json:"clusterSpaceUtilizationStats"`
	MetaDataUtilizationStats         values                 `json:"metaDataUtilizationStats"`
	DeduplicationStats               values                 `json:"deduplicationStats"`
	JobBackUpTimeStats               values                 `json:"jobBackUpTimesStats"`
	JobArchivalTimeStats             values                 `json:"jobArchivalTimeStats"`
	JobReplicationTimeStats          values                 `json:"jobReplicationTimeStats"`
	Protectiongroupstats             []protectiongroupstats `json:"protectionGroupStats"`
}

type protectiongroupstats struct {
	Name                    string `json:"name"`
	JobBackUpTimeStats      values `json:"jobBackUpTime"`
	JobArchivalTimeStats    values `json:"jobArchivalTimeStats"`
	JobReplicationTimeStats values `json:"jobReplicationTimeStats"`
}
type values struct {
	Mean     float64   `json:"mean"`
	SD       float64   `json:"standardDeviation"`
	PastData []float64 `json:"pastData"`
}

var Stats_response datastats
var alr = true
var count = 0
var previous_sample_size float64
var IndividualProtectionJobsMap = make(map[string]protectiongroupstats)
var myProtectionGroups []protectiongroupstats

func mean(prev_mean float64, prev_sample_size float64, new_data float64) (new_mean float64) {

	val1 := (prev_sample_size) / (prev_sample_size + 1)
	val1 *= prev_mean

	val2 := new_data / (prev_sample_size + 1)

	new_mean = val1 + val2
	return
}

func standardDeviation(prev_sd, prev_mean, new_mean, prev_sample_size, new_data float64) (new_sd float64) {

	val1 := (prev_sample_size) / (prev_sample_size + 1)
	val1 *= prev_sd * prev_sd

	val2 := (new_data - new_mean)
	val2 /= (prev_sample_size + 1)
	val2 *= (new_data - prev_mean)
	val1 += val2

	new_sd = math.Sqrt(val1)
	return
}

func Examine(timestamp,num_of_stats int64) {
	RetriveStatsFronElasticDb(num_of_stats)
	CheckSoftwareUpdate()

	alr = true
	var job_backup_percent_deviated, job_replication_percent_deviated, job_archival_percent_deviated float64
	previous_sample_size = float64(num_of_stats)

	if MetricDataAnalysis(float64(Response_for_elastic.AverageJobbackupTimeForAllJobsUsecs), Stats_response.JobBackUpTimeStats.SD, Stats_response.JobBackUpTimeStats.PastData) {
		job_backup_percent_deviated = PercentageCalculator(GetMedian(Stats_response.JobBackUpTimeStats.PastData), float64(Response_for_elastic.AverageJobbackupTimeForAllJobsUsecs))
		alr = alr && false
	}

	if MetricDataAnalysis(float64(Response_for_elastic.AverageJobreplicationTimeForAllJobsUsecs), Stats_response.JobReplicationTimeStats.SD, Stats_response.JobReplicationTimeStats.PastData) {
		job_replication_percent_deviated = PercentageCalculator(GetMedian(Stats_response.JobReplicationTimeStats.PastData), float64(Response_for_elastic.AverageJobreplicationTimeForAllJobsUsecs))
		alr = alr && false
	}

	if MetricDataAnalysis(float64(Response_for_elastic.AverageJobarchivalTimeForAllJobsUsecs), Stats_response.JobArchivalTimeStats.SD, Stats_response.JobArchivalTimeStats.PastData) {
		job_archival_percent_deviated = PercentageCalculator(GetMedian(Stats_response.JobArchivalTimeStats.PastData), float64(Response_for_elastic.AverageJobarchivalTimeForAllJobsUsecs))
		alr = alr && false
	}

	if job_backup_percent_deviated > 0 || job_replication_percent_deviated > 0 || job_archival_percent_deviated > 0 {
		count++
		colorRed := "\033[31m"
		fmt.Println(string(colorRed))
		fmt.Println("						ALERT - ", count)
		fmt.Println()
		PrintTime(timestamp)
		fmt.Println()
		colorReset := "\033[0m"
		fmt.Println(string(colorReset))
		if job_backup_percent_deviated >= job_replication_percent_deviated {
			if job_backup_percent_deviated >= job_archival_percent_deviated {
				fmt.Println("- Average Backup time increased by ", math.Round(math.Abs(job_backup_percent_deviated)), "% . While the Backup was running these Replication/Archival jobs were running.")
			} else {
				fmt.Println("- Average Job Archival time increased by ", math.Round(math.Abs(job_archival_percent_deviated)), "% . While the Archival was running these Backup/Replication jobs were running.")
			}
		} else {
			if job_replication_percent_deviated >= job_archival_percent_deviated {
				fmt.Println("- Average Job Replication time increased by ", math.Round(math.Abs(job_replication_percent_deviated)), "%. While the Replication was running these Backup/Archival jobs were running.")
			} else {
				fmt.Println("- Average Job Archival time increased by ", math.Round(math.Abs(job_archival_percent_deviated)), "% . While the Archival was running these Backup/Replication jobs were running.")
			}
		}
	}
	//else{
	// 	fmt.Println("- No backup job deviated")
	// }

	// fmt.Println()
	// fmt.Println("- Mertics that got deviated from the medain are : ")
	// fmt.Println()
	// if MetricDataAnalysis(float64(Response_for_elastic.AverageJobbackupTimeForAllJobsUsecs), Stats_response.JobBackUpTimeStats.SD, Stats_response.JobBackUpTimeStats.PastData) {
	// 	percent_deviated = PercentageCalculator(GetMedian(Stats_response.JobBackUpTimeStats.PastData), float64(Response_for_elastic.AverageJobbackupTimeForAllJobsUsecs))
	// 	fmt.Println("	Job Back up time deviated by ", percent_deviated, "%")
	// }

	// if MetricDataAnalysis(float64(Response_for_elastic.AverageJobreplicationTimeForAllJobsUsecs), Stats_response.JobReplicationTimeStats.SD, Stats_response.JobReplicationTimeStats.PastData) {
	// 	percent_deviated = PercentageCalculator(GetMedian(Stats_response.JobReplicationTimeStats.PastData), float64(Response_for_elastic.AverageJobreplicationTimeForAllJobsUsecs))
	// 	fmt.Println("	Job Replication time deviated by ", percent_deviated, "%")
	// }

	// if MetricDataAnalysis(float64(Response_for_elastic.AverageJobarchivalTimeForAllJobsUsecs), Stats_response.JobArchivalTimeStats.SD, Stats_response.JobArchivalTimeStats.PastData) {
	// 	percent_deviated = PercentageCalculator(GetMedian(Stats_response.JobArchivalTimeStats.PastData), float64(Response_for_elastic.AverageJobarchivalTimeForAllJobsUsecs))
	// 	fmt.Println("	Job Archival time deviated by ", percent_deviated, "%")
	// }

	// if MetricDataAnalysis(float64(Response_for_elastic.FileLatency), Stats_response.FileLatencyStats.SD, Stats_response.FileLatencyStats.PastData) {
	// 	percent_deviated = PercentageCalculator(GetMedian(Stats_response.FileLatencyStats.PastData), float64(Response_for_elastic.FileLatency))
	// 	fmt.Println("	File latency deviated by ", percent_deviated, "%")
	// }

	// if MetricDataAnalysis(float64(Response_for_elastic.FileCreateSum), Stats_response.FileCreateSumStats.SD, Stats_response.FileCreateSumStats.PastData) {
	// 	percent_deviated = PercentageCalculator(GetMedian(Stats_response.FileCreateSumStats.PastData), float64(Response_for_elastic.FileCreateSum))
	// 	fmt.Println("	File create sum deviated by ", percent_deviated, "%")
	// }

	// if MetricDataAnalysis(float64(Response_for_elastic.SystemUtilizationChangeRate), Stats_response.SystemUtilizationChangeRateStats.SD, Stats_response.SystemUtilizationChangeRateStats.PastData) {
	// 	percent_deviated = PercentageCalculator(GetMedian(Stats_response.SystemUtilizationChangeRateStats.PastData), float64(Response_for_elastic.SystemUtilizationChangeRate))
	// 	fmt.Println("	System Utilization Change Rate deviated by ", percent_deviated, "%")
	// }

	// if MetricDataAnalysis(float64(Response_for_elastic.GarbageCollection), Stats_response.GarbageCollectionStats.SD, Stats_response.GarbageCollectionStats.PastData) {
	// 	percent_deviated = PercentageCalculator(GetMedian(Stats_response.GarbageCollectionStats.PastData), float64(Response_for_elastic.GarbageCollection))
	// 	fmt.Println("	Garbage Collection deviated by ", percent_deviated, "%")
	// }

	// if MetricDataAnalysis(float64(Response_for_elastic.ClusterSpaceUtilization), Stats_response.ClusterSpaceUtilizationStats.SD, Stats_response.ClusterSpaceUtilizationStats.PastData) {
	// 	percent_deviated = PercentageCalculator(GetMedian(Stats_response.ClusterSpaceUtilizationStats.PastData), float64(Response_for_elastic.ClusterSpaceUtilization))
	// 	fmt.Println("	Cluster Space Utilization deviated by ", percent_deviated, "%")
	// }

	// if MetricDataAnalysis(float64(Response_for_elastic.MetaDataUtilization), Stats_response.MetaDataUtilizationStats.SD, Stats_response.MetaDataUtilizationStats.PastData) {
	// 	percent_deviated = PercentageCalculator(GetMedian(Stats_response.MetaDataUtilizationStats.PastData), float64(Response_for_elastic.MetaDataUtilization))
	// 	fmt.Println("	Meta Data Utilization deviated by ", percent_deviated, "%")
	// }

	// if MetricDataAnalysis(float64(Response_for_elastic.Deduplication), Stats_response.DeduplicationStats.SD, Stats_response.DeduplicationStats.PastData) {
	// 	percent_deviated = PercentageCalculator(GetMedian(Stats_response.DeduplicationStats.PastData), float64(Response_for_elastic.Deduplication))
	// 	fmt.Println("	Deduplication factor deviated by ", percent_deviated, "%")
	// }

	if bool(alr) {
		
	} else {
		
		NewProtectionJobs()
		DisplayData()

	}

	//UpdateStats()
	//fmt.Println()

}

func CheckSoftwareUpdate() {
	ver1 := Stats_response.ClusterSoftwareVersion

	ver2 := Response_for_elastic.ClusterSoftwareVersion

	if ver1 != ver2 {
		fmt.Println("Cluster updated from version ", ver1, " to ", ver2)
	}
}

func PercentageCalculator(median, datavalue float64) (perdev float64) {
	perdev = (datavalue - median)
	perdev /= median
	perdev *= 100
	return
}

func MetricDataAnalysis(datavalue, sd float64, pastdata []float64) (alert bool) {

	median := GetMedian(pastdata)
	alert = CheckAlert(median, sd, datavalue)
	return
}

func GetMedian(data []float64) (median float64) {

	sort.Float64s(data)
	num_of_data := len(data)
	if num_of_data > 0 {
		if num_of_data%2 == 0 {
			median = data[(num_of_data-2)/2] + data[(num_of_data)/2]
			median /= 2
		} else {
			median = data[(num_of_data-1)/2]
		}
	}
	return
}

func CheckAlert(median, SD, value float64) (alert bool) {

	left_boundry := median - 2*SD
	right_boundry := median + 2*SD

	if value > right_boundry || value < left_boundry {
		alert = true
	} else {
		alert = false
	}
	return
}

func DisplayData() {

	//fmt.Println("Id = ", Response_for_elastic.Id)

	// `fmt.Println("- Behaviour of other metrics during this time :")`
	//fmt.Println()
	var percent_deviated float64
	colorReset := "\033[0m"
	fmt.Println(string(colorReset))

	// if !MetricDataAnalysis(float64(Response_for_elastic.AverageJobbackupTimeForAllJobsUsecs), Stats_response.JobBackUpTimeStats.SD, Stats_response.JobBackUpTimeStats.PastData) {
	// 	percent_deviated = PercentageCalculator(GetMedian(Stats_response.JobBackUpTimeStats.PastData), float64(Response_for_elastic.AverageJobbackupTimeForAllJobsUsecs))
	// 	if percent_deviated > 0{
	// 	fmt.Println("	Job Back up time increased by ",  math.Round(math.Abs(percent_deviated)), "%")
	// 	}else{
	// 		fmt.Println("	Job Back up time decreased by ",  math.Round(math.Abs(percent_deviated)), "%")
	// 	}
	// }

	// if !MetricDataAnalysis(float64(Response_for_elastic.AverageJobreplicationTimeForAllJobsUsecs), Stats_response.JobReplicationTimeStats.SD, Stats_response.JobReplicationTimeStats.PastData) {
	// 	percent_deviated = PercentageCalculator(GetMedian(Stats_response.JobReplicationTimeStats.PastData), float64(Response_for_elastic.AverageJobreplicationTimeForAllJobsUsecs))
	// 	if percent_deviated > 0{
	// 		fmt.Println("	Job Replication time increased by ",  math.Round(math.Abs(percent_deviated)), "%")
	// 		}else{
	// 			fmt.Println("	Job Replication time decreased by ",  math.Round(math.Abs(percent_deviated)), "%")
	// 		}
	// }

	// if !MetricDataAnalysis(float64(Response_for_elastic.AverageJobarchivalTimeForAllJobsUsecs), Stats_response.JobArchivalTimeStats.SD, Stats_response.JobArchivalTimeStats.PastData) {
	// 	percent_deviated = PercentageCalculator(GetMedian(Stats_response.JobArchivalTimeStats.PastData), float64(Response_for_elastic.AverageJobarchivalTimeForAllJobsUsecs))
	// 	if percent_deviated > 0{
	// 		fmt.Println("	Job Archival time increased by ",  math.Round(math.Abs(percent_deviated)), "%")
	// 		}else{
	// 			fmt.Println("	Job Archival time decreased by ", math.Round(math.Abs(percent_deviated)), "%")
	// 		}
	// }

	//if!MetricDataAnalysis(float64(Response_for_elastic.GarbageCollection), Stats_response.GarbageCollectionStats.SD, Stats_response.GarbageCollectionStats.PastData) {
	percent_deviated = PercentageCalculator(GetMedian(Stats_response.GarbageCollectionStats.PastData), float64(Response_for_elastic.GarbageCollection))

	if percent_deviated > 0 {
		fmt.Print("- Garbage Collection increased by ", math.Round(math.Abs(percent_deviated)), "%")
	} else {
		fmt.Print("- Garbage Collection decreased by ", math.Round(math.Abs(percent_deviated)), "%")
	}

	//if !MetricDataAnalysis(float64(Response_for_elastic.Deduplication), Stats_response.DeduplicationStats.SD, Stats_response.DeduplicationStats.PastData) {
	percent_deviated = PercentageCalculator(GetMedian(Stats_response.DeduplicationStats.PastData), float64(Response_for_elastic.Deduplication))

	if percent_deviated > 0 {
		fmt.Println(" during the period. Also, Deduplication factor increased by ", math.Round(math.Abs(percent_deviated)), "%")
	} else {
		fmt.Println(" during the period. Also, Deduplication factor decreased by ", math.Round(math.Abs(percent_deviated)), "%")
	}

	fmt.Println()

	//if !MetricDataAnalysis(float64(Response_for_elastic.FileLatency), Stats_response.FileLatencyStats.SD, Stats_response.FileLatencyStats.PastData) {
	percent_deviated = PercentageCalculator(GetMedian(Stats_response.FileLatencyStats.PastData), float64(Response_for_elastic.FileLatency))

	if percent_deviated > 0 {
		fmt.Print("- Filecreate latency increased by ", math.Round(math.Abs(percent_deviated)), "% [current value at ",roundFloat(Response_for_elastic.FileLatency,1) ," ms].")
	}else{
		fmt.Print("- Filecreate latency decreased by ", math.Round(math.Abs(percent_deviated)), "% [current value at ",roundFloat(Response_for_elastic.FileLatency,1) ," ms].")
	}

	//if !MetricDataAnalysis(float64(Response_for_elastic.SystemUtilizationChangeRate), Stats_response.SystemUtilizationChangeRateStats.SD, Stats_response.SystemUtilizationChangeRateStats.PastData) {
	percent_deviated = PercentageCalculator(GetMedian(Stats_response.SystemUtilizationChangeRateStats.PastData), float64(Response_for_elastic.SystemUtilizationChangeRate))

	if percent_deviated > 0 {
		fmt.Print(" Also, System Utilization(Space) change rate increased by ", math.Round(math.Abs(percent_deviated)), "% [current value at ",roundFloat(Response_for_elastic.SystemUtilizationChangeRate,1)," GiB/sec]" )
	} else {
		fmt.Print(" Also, System Utilization(Space) change rate decreased by ", math.Round(math.Abs(percent_deviated)), "% [current value at ",roundFloat(Response_for_elastic.SystemUtilizationChangeRate,1)," GiB/sec]" )
	}

	//if !MetricDataAnalysis(float64(Response_for_elastic.ClusterSpaceUtilization), Stats_response.ClusterSpaceUtilizationStats.SD, Stats_response.ClusterSpaceUtilizationStats.PastData) {
	percent_deviated = PercentageCalculator(GetMedian(Stats_response.ClusterSpaceUtilizationStats.PastData), float64(Response_for_elastic.ClusterSpaceUtilization))

	if percent_deviated > 0 {
		fmt.Println(" and Cluster Space Utilization increased by ", math.Round(math.Abs(percent_deviated)), "% [current value at ", roundFloat(Response_for_elastic.ClusterSpaceUtilization,1) ," TiB]"  )
	} else {
		fmt.Println(" and Cluster Space Utilization decreased by ", math.Round(math.Abs(percent_deviated)), "% [current value at ", roundFloat(Response_for_elastic.ClusterSpaceUtilization,1) ," TiB]")
	}
	fmt.Println()

	fmt.Println("- Meta Data Utilisation stands at ", roundFloat(Response_for_elastic.MetaDataSpacePercentage,1) ,"% with a current value of",roundFloat(Response_for_elastic.MetaDataUtilization,1) ,"TiB" )
	fmt.Println()
	// 	if !MetricDataAnalysis(float64(Response_for_elastic.FileCreateSum), Stats_response.FileCreateSumStats.SD, Stats_response.FileCreateSumStats.PastData) {
	// 		percent_deviated = PercentageCalculator(GetMedian(Stats_response.FileCreateSumStats.PastData), float64(Response_for_elastic.FileCreateSum))
	// 		if percent_deviated > 0{
	// 			fmt.Println("	File Create Sum increased by ", math.Round(math.Abs(percent_deviated)), "%")
	// 			}else{
	// 				fmt.Println("	File Create Sum decreased by ", math.Round(math.Abs(percent_deviated)), "%")
	// 			}
	// 	}

	// 	if !MetricDataAnalysis(float64(Response_for_elastic.MetaDataUtilization), Stats_response.MetaDataUtilizationStats.SD, Stats_response.MetaDataUtilizationStats.PastData) {
	// 		percent_deviated = PercentageCalculator(GetMedian(Stats_response.MetaDataUtilizationStats.PastData), float64(Response_for_elastic.MetaDataUtilization))
	// 		if percent_deviated > 0{
	// 			fmt.Println("	Meta Data Utilization increased by ", math.Round(math.Abs(percent_deviated)), "%")
	// 			}else{
	// 				fmt.Println("	Meta Data Utilization decreased by ", math.Round(math.Abs(percent_deviated)), "%")
	// 			}
	// 	}
}

func AssignValues(mean, sd float64, pastdata []float64) (jsonvalues values) {

	jsonvalues = values{
		Mean:     mean,
		SD:       sd,
		PastData: pastdata,
	}
	return
}

func UpdateValues(past_sd, past_mean, new_value float64, past_data []float64) (new_mean, new_sd float64, new_data []float64) {

	if new_value > 0 {
		past_smp_size := float64(len(past_data))
		past_data = append(past_data, new_value)
		new_data = past_data
		new_mean = mean(past_mean, past_smp_size, new_value)
		new_sd = standardDeviation(past_sd, past_mean, new_mean, past_smp_size, new_value)
	} else {
		new_mean = past_mean
		new_sd = past_sd
		new_data = past_data
	}

	return
}

func NewProtectionJobs() {

	number_of_jobs := len(Stats_response.Protectiongroupstats)

	for i := 0; i < number_of_jobs; i++ {
		IndividualProtectionJobsMap[Stats_response.Protectiongroupstats[i].Name] = Stats_response.Protectiongroupstats[i]
	}

	num_of_jobs := len(Response_for_elastic.ProtectionJobsInfo)
	c := 0
	for i := 0; i < num_of_jobs; i++ {
		for j := 0; j < Response_for_elastic.ProtectionJobsInfo[i].Runs; j++ {
			//	Addmetrics("Sanchit", "b", 1, "r", 2, "a", 3)
			fmt.Println()
			colorYellow := "\033[33m"
			if c == 0 {
				fmt.Println(string(colorYellow), "- Jobs that were completed during this time : ")
			}
			c++
			protectiongroup_jobname := Response_for_elastic.ProtectionJobsInfo[i].ProtectionGroupName
			temp := IndividualProtectionJobsMap[protectiongroup_jobname]
			temp.Name = protectiongroup_jobname

			fmt.Println()
			fmt.Println("	Job Name : ", temp.Name)
			fmt.Println()
			protectiongroup_jb_new_data := float64(Response_for_elastic.ProtectionJobsInfo[i].JobBackuptimeUsecs)
			
			var tempbackup string
			if protectiongroup_jb_new_data > 0 {
				fmt.Println("	JobBackuptime : ", Response_for_elastic.ProtectionJobsInfo[i].JobBackuptime)
				protectiongroup_jb_prev_sample_size := float64(len(Stats_response.JobBackUpTimeStats.PastData))
				temp.JobBackUpTimeStats.PastData = append(Stats_response.JobBackUpTimeStats.PastData, protectiongroup_jb_new_data)
				protectiongroup_jb_new_mean := mean(Stats_response.JobBackUpTimeStats.Mean, protectiongroup_jb_prev_sample_size, protectiongroup_jb_new_data)
				protectiongroup_jb_new_sd := standardDeviation(Stats_response.JobBackUpTimeStats.SD, Stats_response.JobBackUpTimeStats.Mean, protectiongroup_jb_new_mean, protectiongroup_jb_prev_sample_size, protectiongroup_jb_new_data)
				temp.JobBackUpTimeStats.Mean = protectiongroup_jb_new_mean
				temp.JobBackUpTimeStats.SD = protectiongroup_jb_new_sd

				percent_deviated := PercentageCalculator(GetMedian(temp.JobBackUpTimeStats.PastData), protectiongroup_jb_new_data)

				if percent_deviated > 0 {
					fmt.Println("	Job Back up time increased by ", math.Round(math.Abs(percent_deviated)), "%")
					s := strconv.FormatFloat(math.Round(math.Abs(percent_deviated)), 'f', 0, 64)
					tempbackup = "	Job Back up time increased by " + s + "%"
					//Addmetrics(temp.Name,Response_for_elastic.ProtectionJobsInfo[i].JobBackuptime,tempstr)
				} else {
					fmt.Println("	Job Back up time decreased by ", math.Round(math.Abs(percent_deviated)), "%")
					s := strconv.FormatFloat(math.Round(math.Abs(percent_deviated)), 'f', 0, 64)
					tempbackup = "	Job Back up time decreased by " + s + "%"
					//	Addmetrics(temp.Name,Response_for_elastic.ProtectionJobsInfo[i].JobBackuptime,)
				}

				fmt.Println()

			}

			protectiongroup_jr_new_data := float64(Response_for_elastic.ProtectionJobsInfo[i].JobReplicationTimeUsecs)
			
			var temprepl string
			if protectiongroup_jr_new_data > 0 {
				fmt.Println("	JobReplicationTime : ", Response_for_elastic.ProtectionJobsInfo[i].JobReplicationTime)
				protectiongroup_jr_prev_sample_size := float64(len(Stats_response.JobReplicationTimeStats.PastData))
				temp.JobReplicationTimeStats.PastData = append(Stats_response.JobReplicationTimeStats.PastData, protectiongroup_jr_new_data)
				protectiongroup_jr_new_mean := mean(Stats_response.JobReplicationTimeStats.Mean, protectiongroup_jr_prev_sample_size, protectiongroup_jr_new_data)
				protectiongroup_jr_new_sd := standardDeviation(Stats_response.JobReplicationTimeStats.SD, temp.JobReplicationTimeStats.Mean, protectiongroup_jr_new_mean, protectiongroup_jr_prev_sample_size, protectiongroup_jr_new_data)
				temp.JobReplicationTimeStats.Mean = protectiongroup_jr_new_mean
				temp.JobReplicationTimeStats.SD = protectiongroup_jr_new_sd

				percent_deviated := PercentageCalculator(GetMedian(temp.JobReplicationTimeStats.PastData), protectiongroup_jr_new_data)

				if percent_deviated > 0 {
					fmt.Println("	Job Replication time increased by ", math.Round(math.Abs(percent_deviated)), "%")
					s := strconv.FormatFloat(math.Round(math.Abs(percent_deviated)), 'f', 0, 64)
					temprepl = "	Job Replication time increased by  " + s + "%"
				} else {
					fmt.Println("	Job Replication time decreased by ", math.Round(math.Abs(percent_deviated)), "%")
					s := strconv.FormatFloat(math.Round(math.Abs(percent_deviated)), 'f', 0, 64)
					temprepl = "	Job Replication time decreased by  " + s + "%"
				}
				fmt.Println()

			}

			protectiongroup_ja_new_data := float64(Response_for_elastic.ProtectionJobsInfo[i].JobArchivalTimeUsecs)
			
			var temparch string
			if protectiongroup_ja_new_data > 0 {
				fmt.Println("	JobArchivalTime : ", Response_for_elastic.ProtectionJobsInfo[i].JobArchivalTime)
				protectiongroup_ja_prev_sample_size := float64(len(Stats_response.JobArchivalTimeStats.PastData))
				temp.JobArchivalTimeStats.PastData = append(Stats_response.JobArchivalTimeStats.PastData, protectiongroup_ja_new_data)
				protectiongroup_ja_new_mean := mean(Stats_response.JobArchivalTimeStats.Mean, protectiongroup_ja_prev_sample_size, protectiongroup_ja_new_data)
				protectiongroup_ja_new_sd := standardDeviation(Stats_response.JobArchivalTimeStats.SD, Stats_response.JobArchivalTimeStats.Mean, protectiongroup_ja_new_mean, protectiongroup_ja_prev_sample_size, protectiongroup_ja_new_data)
				temp.JobArchivalTimeStats.Mean = protectiongroup_ja_new_mean
				temp.JobArchivalTimeStats.SD = protectiongroup_ja_new_sd

				percent_deviated := PercentageCalculator(GetMedian(temp.JobArchivalTimeStats.PastData), protectiongroup_ja_new_data)

				if percent_deviated > 0 {
					fmt.Println("	Job Archival time increased by ", math.Round(math.Abs(percent_deviated)), "%")
					s := strconv.FormatFloat(math.Round(math.Abs(percent_deviated)), 'f', 0, 64)
					temparch = "	Job archival time increased by  " + s + "%"
				} else {
					fmt.Println("	Job Archival time decreased by ", math.Round(math.Abs(percent_deviated)), "%")
					s := strconv.FormatFloat(math.Round(math.Abs(percent_deviated)), 'f', 0, 64)
					temparch = "	Job archival time decreased by  " + s + "%"
				}
				// fmt.Println()

			}
			Addmetrics(temp.Name, "	Job BackUp Time : "+Response_for_elastic.ProtectionJobsInfo[i].JobBackuptime, tempbackup, "	Job Replication Time : "+Response_for_elastic.ProtectionJobsInfo[i].JobReplicationTime, temprepl, "	JobArchivalTime : "+Response_for_elastic.ProtectionJobsInfo[i].JobArchivalTime, temparch)

			IndividualProtectionJobsMap[protectiongroup_jobname] = temp
			//fmt.Println("IndividualProtectionJobsMap",protectiongroup_jobname, " : ",IndividualProtectionJobsMap[protectiongroup_jobname])
		}

		//generator()

		for _, v := range IndividualProtectionJobsMap {

			myProtectionGroups = append(myProtectionGroups, v)
		}
	}
}

func UpdateStats() {

	FileCreateRate_mean, FileCreateRate_SD, FileCreateRate_data := UpdateValues(Stats_response.FileCreateRateStats.SD, Stats_response.FileCreateRateStats.Mean, float64(Response_for_elastic.FileCreateRate), Stats_response.FileCreateRateStats.PastData)
	FileCreateSum_mean, FileCreateSum_SD, FileCreateSum_data := UpdateValues(Stats_response.FileCreateSumStats.SD, Stats_response.FileCreateSumStats.Mean, float64(Response_for_elastic.FileCreateSum), Stats_response.FileCreateSumStats.PastData)
	SystemUtilizationChangeRate_mean, SystemUtilizationChangeRate_SD, SystemUtilizationChangeRate_data := UpdateValues(Stats_response.SystemUtilizationChangeRateStats.SD, Stats_response.SystemUtilizationChangeRateStats.Mean, float64(Response_for_elastic.SystemUtilizationChangeRate), Stats_response.SystemUtilizationChangeRateStats.PastData)
	GarbageCollection_mean, GarbageCollection_SD, GarbageCollection_data := UpdateValues(Stats_response.GarbageCollectionStats.SD, Stats_response.GarbageCollectionStats.Mean, float64(Response_for_elastic.GarbageCollection), Stats_response.GarbageCollectionStats.PastData)
	ClusterSpaceUtilization_mean, ClusterSpaceUtilization_SD, ClusterSpaceUtilization_data := UpdateValues(Stats_response.ClusterSpaceUtilizationStats.SD, Stats_response.ClusterSpaceUtilizationStats.Mean, float64(Response_for_elastic.ClusterSpaceUtilization), Stats_response.ClusterSpaceUtilizationStats.PastData)
	MetaDataUtilization_mean, MetaDataUtilization_SD, MetaDataUtilization_data := UpdateValues(Stats_response.MetaDataUtilizationStats.SD, Stats_response.MetaDataUtilizationStats.Mean, float64(Response_for_elastic.MetaDataUtilization), Stats_response.MetaDataUtilizationStats.PastData)
	Deduplication_mean, Deduplication_SD, Deduplication_data := UpdateValues(Stats_response.DeduplicationStats.SD, Stats_response.DeduplicationStats.Mean, float64(Response_for_elastic.Deduplication), Stats_response.DeduplicationStats.PastData)
	FileLatency_mean, FileLatency_SD, FileLatency_data := UpdateValues(Stats_response.FileLatencyStats.SD, Stats_response.FileLatencyStats.Mean, float64(Response_for_elastic.FileLatency), Stats_response.FileLatencyStats.PastData)
	JobBackUpTime_mean, JobBackUpTime_SD, JobBackUpTime_data := UpdateValues(Stats_response.JobBackUpTimeStats.SD, Stats_response.JobBackUpTimeStats.Mean, float64(Response_for_elastic.AverageJobbackupTimeForAllJobsUsecs), Stats_response.JobBackUpTimeStats.PastData)
	JobArchivalTime_mean, JobArchivalTime_SD, JobArchivalTime_data := UpdateValues(Stats_response.JobArchivalTimeStats.SD, Stats_response.JobArchivalTimeStats.Mean, float64(Response_for_elastic.AverageJobreplicationTimeForAllJobsUsecs), Stats_response.JobArchivalTimeStats.PastData)
	JobReplicationTime_mean, JobReplicationTime_SD, JobReplicationTime_data := UpdateValues(Stats_response.JobReplicationTimeStats.SD, Stats_response.JobReplicationTimeStats.Mean, float64(Response_for_elastic.AverageJobarchivalTimeForAllJobsUsecs), Stats_response.JobReplicationTimeStats.PastData)

	newstats := datastats{

		Id:                               int64(previous_sample_size),
		FileCreateRateStats:              AssignValues(FileCreateRate_mean, FileCreateRate_SD, FileCreateRate_data),
		FileCreateSumStats:               AssignValues(FileCreateSum_mean, FileCreateSum_SD, FileCreateSum_data),
		FileLatencyStats:                 AssignValues(FileLatency_mean, FileLatency_SD, FileLatency_data),
		SystemUtilizationChangeRateStats: AssignValues(SystemUtilizationChangeRate_mean, SystemUtilizationChangeRate_SD, SystemUtilizationChangeRate_data),
		ClusterSpaceUtilizationStats:     AssignValues(ClusterSpaceUtilization_mean, ClusterSpaceUtilization_SD, ClusterSpaceUtilization_data),
		GarbageCollectionStats:           AssignValues(GarbageCollection_mean, GarbageCollection_SD, GarbageCollection_data),
		MetaDataUtilizationStats:         AssignValues(MetaDataUtilization_mean, MetaDataUtilization_SD, MetaDataUtilization_data),
		DeduplicationStats:               AssignValues(Deduplication_mean, Deduplication_SD, Deduplication_data),
		JobBackUpTimeStats:               AssignValues(JobBackUpTime_mean, JobBackUpTime_SD, JobBackUpTime_data),
		JobArchivalTimeStats:             AssignValues(JobArchivalTime_mean, JobArchivalTime_SD, JobArchivalTime_data),
		JobReplicationTimeStats:          AssignValues(JobReplicationTime_mean, JobReplicationTime_SD, JobReplicationTime_data),
		Protectiongroupstats:             myProtectionGroups,
	}

	stats_data, _ := json.Marshal(newstats)
	StoreInElasticDb(stats_data, "stats_test")

}

func UserRequestForData(starttime, endtime int64) {
	colorWhite := "\033[37m"
	fmt.Print(string(colorWhite))
	timeT := time.Unix(endtime, 0)
	loc, _ := time.LoadLocation("Asia/Kolkata")
	now := timeT.In(loc)
	start_timeT := time.Unix(starttime, 0)
	start_now := start_timeT.In(loc)
	fmt.Println("During the period from ", start_now, " to ", now)
	hrs := (endtime - starttime)/3600
	fmt.Println("For the duration of ",hrs," Hours")
	fmt.Println()
	var filecreatesumdata, filecreateratedata, SystemUtilizationChangeRatedata, GarbageCollectiondata, ClusterSpaceUtilizationdata, filelatencydata []float64
	var MetaDataUtilizationdata, Deduplicationdata, jobbackupdata, jobrepdata, jobarchdata []float64
	var filecreaterate_mean, filecreaterate_sd float64
	var filecreatesum_mean, filecreatesum_sd float64
	var SystemUtilizationChangeRate_mean, SystemUtilizationChangeRate_sd float64
	var GarbageCollection_mean, GarbageCollection_sd float64
	var ClusterSpaceUtilization_mean, ClusterSpaceUtilization_sd float64
	var MetaDataUtilization_mean, MetaDataUtilization_sd float64
	var Deduplication_mean, Deduplication_sd float64
	var jobbackup_mean, jobbackup_sd float64
	var jobarchival_mean, jobarchival_sd float64
	var jobreplication_mean, jobreplication_sd float64
	var filelatency_mean, filelatency_sd float64
	var ProtectionMap = make(map[string]protectiongroupstats)

	for timestamp := starttime; timestamp <= endtime; timestamp += 3600 {

		RetriveFronElasticDb(int64(timestamp))

		filecreaterate_mean, filecreaterate_sd, filecreateratedata = UpdateValues(filecreaterate_sd, filecreaterate_mean, float64(Response_for_elastic.FileCreateRate), filecreateratedata)
		filecreatesum_mean, filecreatesum_sd, filecreatesumdata = UpdateValues(filecreatesum_sd, filecreatesum_mean, float64(Response_for_elastic.FileCreateSum), filecreatesumdata)
		SystemUtilizationChangeRate_mean, SystemUtilizationChangeRate_sd, SystemUtilizationChangeRatedata = UpdateValues(SystemUtilizationChangeRate_sd, SystemUtilizationChangeRate_mean, float64(Response_for_elastic.SystemUtilizationChangeRate), SystemUtilizationChangeRatedata)
		GarbageCollection_mean, GarbageCollection_sd, GarbageCollectiondata = UpdateValues(GarbageCollection_sd, GarbageCollection_mean, float64(Response_for_elastic.GarbageCollection), GarbageCollectiondata)
		ClusterSpaceUtilization_mean, ClusterSpaceUtilization_sd, ClusterSpaceUtilizationdata = UpdateValues(ClusterSpaceUtilization_sd, ClusterSpaceUtilization_mean, float64(Response_for_elastic.ClusterSpaceUtilization), ClusterSpaceUtilizationdata)
		MetaDataUtilization_mean, MetaDataUtilization_sd, MetaDataUtilizationdata = UpdateValues(MetaDataUtilization_sd, MetaDataUtilization_mean, float64(Response_for_elastic.MetaDataUtilization), MetaDataUtilizationdata)
		Deduplication_mean, Deduplication_sd, Deduplicationdata = UpdateValues(Deduplication_sd, Deduplication_mean, float64(Response_for_elastic.Deduplication), Deduplicationdata)
		filelatency_mean, filelatency_sd, filelatencydata = UpdateValues(filelatency_mean, filelatency_sd, float64(Response_for_elastic.FileLatency), filelatencydata)
		jobbackup_mean, jobbackup_sd, jobbackupdata = UpdateValues(jobbackup_sd, jobbackup_mean, float64(Response_for_elastic.AverageJobbackupTimeForAllJobsUsecs), jobbackupdata)
		jobreplication_mean, jobreplication_sd, jobrepdata = UpdateValues(jobreplication_sd, jobreplication_mean, float64(Response_for_elastic.AverageJobreplicationTimeForAllJobsUsecs), jobrepdata)
		jobarchival_mean, jobarchival_sd, jobarchdata = UpdateValues(jobarchival_sd, jobarchival_mean, float64(Response_for_elastic.AverageJobarchivalTimeForAllJobsUsecs), jobarchdata)

		num_of_jobs := len(Response_for_elastic.ProtectionJobsInfo)

		for i := 0; i < num_of_jobs; i++ {
			for j := 0; j < Response_for_elastic.ProtectionJobsInfo[i].Runs; j++ {

				protectiongroup_jobname := Response_for_elastic.ProtectionJobsInfo[i].ProtectionGroupName
				temp := ProtectionMap[protectiongroup_jobname]
				temp.Name = protectiongroup_jobname
				protectiongroup_jb_new_data := float64(Response_for_elastic.ProtectionJobsInfo[i].JobBackuptimeUsecs)
				if protectiongroup_jb_new_data > 0 {
					protectiongroup_jb_prev_sample_size := float64(len(temp.JobBackUpTimeStats.PastData))
					temp.JobBackUpTimeStats.PastData = append(temp.JobBackUpTimeStats.PastData, protectiongroup_jb_new_data)
					protectiongroup_jb_new_mean := mean(temp.JobBackUpTimeStats.Mean, protectiongroup_jb_prev_sample_size, protectiongroup_jb_new_data)
					protectiongroup_jb_new_sd := standardDeviation(temp.JobBackUpTimeStats.SD, temp.JobBackUpTimeStats.Mean, protectiongroup_jb_new_mean, protectiongroup_jb_prev_sample_size, protectiongroup_jb_new_data)

					protectiongroup_jb_prev_sample_size++

					temp.JobBackUpTimeStats.Mean = protectiongroup_jb_new_mean

					temp.JobBackUpTimeStats.SD = protectiongroup_jb_new_sd

				}

				protectiongroup_jr_new_data := float64(Response_for_elastic.ProtectionJobsInfo[i].JobReplicationTimeUsecs)
				if protectiongroup_jr_new_data > 0 {
					protectiongroup_jr_prev_sample_size := float64(len(temp.JobReplicationTimeStats.PastData))
					temp.JobReplicationTimeStats.PastData = append(temp.JobReplicationTimeStats.PastData, protectiongroup_jr_new_data)
					protectiongroup_jr_new_mean := mean(temp.JobReplicationTimeStats.Mean, protectiongroup_jr_prev_sample_size, protectiongroup_jr_new_data)
					protectiongroup_jr_new_sd := standardDeviation(temp.JobReplicationTimeStats.SD, temp.JobReplicationTimeStats.Mean, protectiongroup_jr_new_mean, protectiongroup_jr_prev_sample_size, protectiongroup_jr_new_data)
					protectiongroup_jr_prev_sample_size++
					temp.JobReplicationTimeStats.Mean = protectiongroup_jr_new_mean
					temp.JobReplicationTimeStats.SD = protectiongroup_jr_new_sd
				}

				protectiongroup_ja_new_data := float64(Response_for_elastic.ProtectionJobsInfo[i].JobArchivalTimeUsecs)
				if protectiongroup_ja_new_data > 0 {
					protectiongroup_ja_prev_sample_size := float64(len(temp.JobArchivalTimeStats.PastData))
					temp.JobArchivalTimeStats.PastData = append(temp.JobArchivalTimeStats.PastData, protectiongroup_ja_new_data)
					protectiongroup_ja_new_mean := mean(temp.JobArchivalTimeStats.Mean, protectiongroup_ja_prev_sample_size, protectiongroup_ja_new_data)
					protectiongroup_ja_new_sd := standardDeviation(temp.JobArchivalTimeStats.SD, temp.JobArchivalTimeStats.Mean, protectiongroup_ja_new_mean, protectiongroup_ja_prev_sample_size, protectiongroup_ja_new_data)
					protectiongroup_ja_prev_sample_size++
					temp.JobArchivalTimeStats.Mean = protectiongroup_ja_new_mean
					temp.JobArchivalTimeStats.SD = protectiongroup_ja_new_sd
				}

				ProtectionMap[protectiongroup_jobname] = temp
			}
		}

	}

	// var ProtectionGroupsinfo []protectiongroupstats
	// for _, v := range ProtectionMap {

	// 	myProtectionGroups = append(myProtectionGroups, v)
	// }

	// mystats := datastats{

	// 	ClusterSoftwareVersion:           Response_for_elastic.ClusterSoftwareVersion,
	// 	FileCreateRateStats:              AssignValues(filecreaterate_mean, filecreaterate_sd, filecreateratedata),
	// 	FileCreateSumStats:               AssignValues(filecreatesum_mean, filecreatesum_sd, filecreatesumdata),
	// 	FileLatencyStats:                 AssignValues(filelatency_mean, filelatency_sd, filelatencydata),
	// 	SystemUtilizationChangeRateStats: AssignValues(SystemUtilizationChangeRate_mean, SystemUtilizationChangeRate_sd, SystemUtilizationChangeRatedata),
	// 	ClusterSpaceUtilizationStats:     AssignValues(ClusterSpaceUtilization_mean, ClusterSpaceUtilization_sd, ClusterSpaceUtilizationdata),
	// 	GarbageCollectionStats:           AssignValues(GarbageCollection_mean, GarbageCollection_sd, GarbageCollectiondata),
	// 	MetaDataUtilizationStats:         AssignValues(MetaDataUtilization_mean, MetaDataUtilization_sd, MetaDataUtilizationdata),
	// 	DeduplicationStats:               AssignValues(Deduplication_mean, Deduplication_sd, Deduplicationdata),
	// 	JobBackUpTimeStats:               AssignValues(jobbackup_mean, jobbackup_sd, jobbackupdata),
	// 	JobArchivalTimeStats:             AssignValues(jobarchival_mean, jobarchival_sd, jobarchdata),
	// 	JobReplicationTimeStats:          AssignValues(jobreplication_mean, jobreplication_sd, jobrepdata),
	// 	Protectiongroupstats:             ProtectionGroupsinfo,
	// }

	//fmt.Println(mystats)
	
	fmt.Println("ClusterSoftwareVersion : ", Response_for_elastic.ClusterSoftwareVersion)
	if(jobbackup_mean >0 ){ fmt.Println("Average BackUp Time : ", ConvertUnixTime(int64(jobbackup_mean))) }
	if(jobarchival_mean >0){fmt.Println("Average Archival Time : ", ConvertUnixTime(int64(jobarchival_mean)))}
	if(jobreplication_mean > 0){fmt.Println("Average Replication Time : ", ConvertUnixTime(int64(jobreplication_mean)))}
	fmt.Println("Average File Create Rate : ", roundFloat(filecreaterate_mean,1))
	fmt.Println("Average File latency : ", roundFloat(filelatency_mean,1), "ms")
	fmt.Println("Average System Utilization Change : ", roundFloat(SystemUtilizationChangeRate_mean,1),"GiB/sec")
	fmt.Println("Average Cluster Space Utilization : ", roundFloat(ClusterSpaceUtilization_mean,1),"TiB")
	fmt.Println("Average Garbage Collection : ", roundFloat(GarbageCollection_mean,1),"GiB")
	fmt.Println("Average Meta Data Utilization : ", roundFloat(MetaDataUtilization_mean,1), "TiB")
	fmt.Println("Average Deduplication : ", roundFloat(Deduplication_mean,1))
	fmt.Println()
	fmt.Println("The jobs which have been completed duirng this period are:")
	fmt.Println()
	for k, _ := range ProtectionMap {
		fmt.Println("	Job Name : ", k)
		if(ProtectionMap[k].JobBackUpTimeStats.Mean > 0){fmt.Println("	Average BackUp Time : ", ConvertUnixTime(int64(ProtectionMap[k].JobBackUpTimeStats.Mean)))  }
		if(ProtectionMap[k].JobArchivalTimeStats.Mean > 0 ){fmt.Println("	Average Archival Time : ", ConvertUnixTime(int64(ProtectionMap[k].JobArchivalTimeStats.Mean))) }
		if(ProtectionMap[k].JobReplicationTimeStats.Mean > 0){fmt.Println("	Average Replication Time : ", ConvertUnixTime(int64(ProtectionMap[k].JobReplicationTimeStats.Mean))) }
		if(len(ProtectionMap[k].JobBackUpTimeStats.PastData) > 0){fmt.Println("	Total runs : ", len(ProtectionMap[k].JobBackUpTimeStats.PastData))}
		fmt.Println()

	}
}

func roundFloat(val float64, precision uint) float64 {
    ratio := math.Pow(10, float64(precision))
    return math.Round(val*ratio) / ratio
}