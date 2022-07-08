package mypackages

import (
	"encoding/json"
	"log"
	"math"
	"net/url"
	"strconv"
)

type metric struct {
	Name    string `json:"metricName"`
	Datavec []rate `json:"dataPointVec"`
}

type rate struct {
	Time  int64     `json:"timestampMsecs"`
	Value ratevalue `json:"data"`
}

type ratevalue struct {
	Data int64 `json:"int64Value"`
}

type smartfileviews struct{
	Views []view `json:"views"`
}
 
type view struct {
	ViewId int64 `json:"viewId"`
	CreateId int64 `json:"createTimeMsecs"`
}
const (

	filecreate_metricName = "kCreateFileOps"
	filecreate_schemaName = "kBridgeViewPerfStats"
	filecreate_metricUnitType = "5"
	filecreate_rate_rollupFunction = "rate"
	filecreate_sum_rollupFunction = "sum"

	filelatency_metricName = "kCreateFileLatUsecs"
	filelatency_schemaName = "kBridgeViewPerfStats"
	filelatency_metricUnitType = "1"
	filelatency_rollupFunction = "average"

	GarbageCollection_metricName = "EstimatedGarbageBytes"
	GarbageCollection_schemaName = "ApolloV2ClusterStats"
	GarbageCollection_metricUnitType = "0"
	GarbageCollection_rollupFunction = "average"
	GarbageCollection_entityId = "st-longevity+(ID+2790138600742128)"

	Cluster_entityId = "2790138600742128"

	SystemUtilizationChangeRate_metricName = "kSystemUsageBytes"
	SystemUtilizationChangeRate_schemaName = "kBridgeClusterStats"
	SystemUtilizationChangeRate_metricUnitType = "0"
	SystemUtilizationChangeRate_rollupFunction = "rate"

	IndexingBacklog_metricName = "yoda_backlog_size"
	IndexingBacklog_schemaName = "YodaBacklogStats"
	IndexingBacklog_metricUnitType = "5"
	IndexingBacklog_rollupFunction = "count"

	MorphedUsage_metricName = "kMorphedUsageBytes"
	UnMorphedUsage_metricName = "kUnmorphedUsageBytes"
	Morphed_and_Unmorphed_schemaName = "kBridgeClusterStats"
	Morphed_and_Unmorphed_metricUnitType = "0"
	Morphed_and_Unmorphed_rollupFunction = "average"
	
	timePeriod = "day"
	rollupIntervalSecs = "180"
	
)

var All_Entities_Id []string
var StartTime_Msecs , Endtime_Msecs string

func GetAllEntitiesId() {
	All_Entities_Id = nil
	url := "https://10.14.19.226/v2/file-services/views?useCachedData=false&maxCount=2000&includeTenants=false&includeStats=false&includeProtectionGroups=false&includeInactive=false"
	
	response := PostRequestForAccessToken()
	data := GetRequestForJsonData(response, url)

	var responsedata smartfileviews

	json.Unmarshal(data, &responsedata)
	instances_of_data := len(responsedata.Views)
	for i := 0; i < instances_of_data; i++ {
		entity_id := strconv.FormatInt(responsedata.Views[i].ViewId, 10)
		All_Entities_Id = append(All_Entities_Id, entity_id)
	}
}

func Filecreaterate() (average_file_create_rate int64) {

	GetAllEntitiesId()
	 StartTime_Msecs, Endtime_Msecs = TimeStampGneerator(timePeriod, "M")

	for i := 0; i < len(All_Entities_Id); i++ {
		entityId := All_Entities_Id[i]
		average_file_create_rate_of_entity := AverageValueOfTimeSeries(Endtime_Msecs, entityId, filecreate_metricName,
			filecreate_metricUnitType,timePeriod, filecreate_rate_rollupFunction, rollupIntervalSecs, filecreate_schemaName,StartTime_Msecs)
		
			average_file_create_rate += average_file_create_rate_of_entity
	}
	average_file_create_rate = average_file_create_rate / int64(len(All_Entities_Id))

	return
}

func FileCreateSum()(average_file_create_sum int64){
	
	for i := 0; i < len(All_Entities_Id); i++ {
		entityId := All_Entities_Id[i]
		average_file_create_sum_of_entity := AverageValueOfTimeSeries(Endtime_Msecs, entityId, filecreate_metricName,
			filecreate_metricUnitType,timePeriod, filecreate_sum_rollupFunction, rollupIntervalSecs, filecreate_schemaName,StartTime_Msecs)
		average_file_create_sum += average_file_create_sum_of_entity
	}
	average_file_create_sum = average_file_create_sum 

	return
	
}

func FileLatency() (average_file_latency float64) {

	for i := 0; i < len(All_Entities_Id); i++ {
		entityId := All_Entities_Id[i]
		average_file_latency_of_entity := float64(AverageValueOfTimeSeries(Endtime_Msecs, entityId, filelatency_metricName, 
			filelatency_metricUnitType, timePeriod, filelatency_rollupFunction, rollupIntervalSecs, filelatency_schemaName, StartTime_Msecs))
		average_file_latency += average_file_latency_of_entity 
	}
	average_file_latency /= 1000
	average_file_latency = math.Round(average_file_latency*100) / 100

	return
}

func SystemUtilizationChangeRate() (average_utilization_change_rate_converted float64) {
	
	average_utilization_change_rate := AverageValueOfTimeSeries(Endtime_Msecs, Cluster_entityId, SystemUtilizationChangeRate_metricName, 
		SystemUtilizationChangeRate_metricUnitType, timePeriod, SystemUtilizationChangeRate_rollupFunction, rollupIntervalSecs, SystemUtilizationChangeRate_schemaName, StartTime_Msecs)
	average_utilization_change_rate_converted = ConvertToMegaBytes(average_utilization_change_rate)
	return
}

func GarbageCollection() (average_garbage_collection_converted float64) {
	
	average_garbage_collection := AverageValueOfTimeSeries(Endtime_Msecs, GarbageCollection_entityId, GarbageCollection_metricName, 
		GarbageCollection_metricUnitType, timePeriod, GarbageCollection_rollupFunction, rollupIntervalSecs, GarbageCollection_schemaName, StartTime_Msecs)
	average_garbage_collection_converted = ConvertToMegaBytes(average_garbage_collection)
	return
}

func IndexingBacklog() (average_indexing_backlog int64) {
	
	average_indexing_backlog = AverageValueOfTimeSeries(Endtime_Msecs, Cluster_entityId, IndexingBacklog_metricName, 
		IndexingBacklog_metricUnitType, timePeriod, IndexingBacklog_rollupFunction, rollupIntervalSecs, IndexingBacklog_schemaName, StartTime_Msecs)
	return
}

func MorphedUsage()(average_morphed_usage int64){
	
	average_morphed_usage = AverageValueOfTimeSeries(Endtime_Msecs, Cluster_entityId , MorphedUsage_metricName, 
		Morphed_and_Unmorphed_metricUnitType, timePeriod, Morphed_and_Unmorphed_rollupFunction, rollupIntervalSecs, Morphed_and_Unmorphed_schemaName, StartTime_Msecs)
	return
}

func UnMorphedUsage()(average_unmorphed_usage int64){
	
	average_unmorphed_usage = AverageValueOfTimeSeries(Endtime_Msecs, Cluster_entityId, UnMorphedUsage_metricName, 
		Morphed_and_Unmorphed_metricUnitType, timePeriod, Morphed_and_Unmorphed_rollupFunction, rollupIntervalSecs, Morphed_and_Unmorphed_schemaName, StartTime_Msecs)
	return
}

func Deduplication()(dedublication float64){

	average_morphed_usage := MorphedUsage()
	average_unmorphed_usage := UnMorphedUsage()
	dedublication = float64(average_unmorphed_usage)/float64(average_morphed_usage)
	dedublication = math.Round(dedublication*100) / 100
	return 
}

func AverageValueOfTimeSeries(endtime_Msecs, entityId, metricName, metricUnitType, timePeriod, 
	rollupFunction, rollupIntervalSecs, schemaName, startTime_Msecs string) (average int64) {

	newURL := GenerateNewURLforTimeSeries(endtime_Msecs, entityId, metricName, metricUnitType, 
		timePeriod, rollupFunction, rollupIntervalSecs, schemaName, startTime_Msecs)
	
	path, err := url.PathUnescape(newURL)
	if err != nil {
		log.Fatal(err)
	}
	response := PostRequestForAccessToken()
	data := GetRequestForJsonData(response, path)

	var responsedata metric

	json.Unmarshal(data, &responsedata)
	count := 0
	for i := 0; i < len(responsedata.Datavec); i++ {
		average += responsedata.Datavec[i].Value.Data
		if(responsedata.Datavec[i].Value.Data > 0){ count++ }

	}
	if rollupFunction == "sum"{
		return average
	}else{
	if count > 0 {
		average /= int64(count)
	} else {
		average = 0
	}
	return
}
}

func ConvertToMegaBytes(average int64) (converted_avg float64) {
	mega_byte := 1024 * 1024 * 1024
	converted_avg = float64(average) / float64(mega_byte)
	converted_avg = math.Round(converted_avg*100) / 100
	return
}
