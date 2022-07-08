package mypackages

import (
	"encoding/json"
	//"fmt"
	"math"
)

type clustersummary struct {
	ClusterSoftwareVersion		string  `json:"clusterSoftwareVersion"`
	Stats                       stats   `json:"stats"`
	AvailableMetadataSpace      int64   `json:"availableMetadataSpace"`
	UsedMetadataSpacePercentage float64 `json:"usedMetadataSpacePct"`
}

type stats struct {
	LocalUsageStats localusagestats `json:"localUsagePerfStats"`
}
type localusagestats struct {
	TotalPhysicalUsageBytes int64 `json:"totalPhysicalUsageBytes"`
}

type alert struct {
	Id string `json:"id"`
	AlertCode string `json:"alertCode"`
	LatestTimeStamp int64 `json:"latestTimestampUsecs"`
	AlertDocument alertdocumet `json:"alertDocument"`

}
type alertdocumet struct {
	AlertName string `json:"alertName"`
}

var AlertsList []alert

func UrlData() (responsedata clustersummary) {

	url := "https://10.14.19.226/irisservices/api/v1/public/cluster?fetchStats=true"
	response := PostRequestForAccessToken()
	data := GetRequestForJsonData(response, url)

	json.Unmarshal(data, &responsedata)
	return
}

func ClusterSpaceUtilization()(cluster_usage float64) {
	responsedata := UrlData()
	total_physical_usage_bytes := responsedata.Stats.LocalUsageStats.TotalPhysicalUsageBytes
	cluster_usage = ContertToTeraBytes(total_physical_usage_bytes)
	return
}

func ClusterVersion()(cluster_version string){

	responsedata := UrlData()
	cluster_version = responsedata.ClusterSoftwareVersion
	return
}

func MetaDataUtilization() (metadata float64, used_meta_space_percentage float64) {
	responsedata := UrlData()
	available_meta_space := responsedata.AvailableMetadataSpace

	used_meta_space_percentage = responsedata.UsedMetadataSpacePercentage
	used_meta_space_percentage = math.Round(used_meta_space_percentage*100) / 100

	metadata = ContertToTeraBytes(available_meta_space)

	return
}

func Alerts(){
	timePeriod := "day"
	start_time_Usecs, end_time_Usecs := TimeStampGneerator(timePeriod, "U")
	newUrl := GenerateNewURLforAlerts(start_time_Usecs,end_time_Usecs)
	response := PostRequestForAccessToken()
	data := GetRequestForJsonData(response, newUrl)

	json.Unmarshal(data, &AlertsList)

	for i :=0;i<len(AlertsList);i++{
		new_id := GenerateAlertId(AlertsList[i].Id)
		AlertsList[i].Id = new_id
	}

}

func GenerateAlertId(id string)(key string){
	
	first_digit_in_id := 0

	for id[first_digit_in_id] != 58 {
		key = key + string(id[first_digit_in_id])
		first_digit_in_id += 1
	}

	return
}

func ContertToTeraBytes(space int64) (space_converted float64) {

	terabyte := 1024 * 1024 * 1024 * 1024

	space_converted = float64(space) / float64(terabyte)

	space_converted = math.Round(space_converted*100) / 100

	return
}
