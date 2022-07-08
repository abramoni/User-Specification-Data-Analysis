package mypackages

import (
	"log"
	"net/url"
)

func GenerateNewURLforTimeSeries(endTimeMsecs, entityId, metricName, metricUnitType, timePeriod, rollupFunction, rollupIntervalSecs, schemaName, startTimeMsecs string) (Result string) {

	myurl := "https://10.14.19.226/irisservices/api/v1/public/statistics/timeSeriesStats?endTimeMsecs=1655103600000&entityId=12522945&metricName=kCreateFileOps&metricUnitType=5&range=day&rollupFunction=rate&rollupIntervalSecs=180&schemaName=kBridgeViewPerfStats&startTimeMsecs=1655017200000"
	result, _ := url.Parse(myurl)

	q := result.Query()
	q.Set("metricName", metricName)
	q.Set("entityId", entityId)
	q.Set("endTimeMsecs", endTimeMsecs)
	q.Set("metricName", metricName)
	q.Set("metricUnitType", metricUnitType)
	q.Set("range", timePeriod)
	q.Set("rollupFunction", rollupFunction)
	q.Set("rollupIntervalSecs", rollupIntervalSecs)
	q.Set("schemaName", schemaName)
	q.Set("startTimeMsecs", startTimeMsecs)
	result.RawQuery = q.Encode()

	query, err := url.QueryUnescape(result.String())

	if err != nil {
		log.Fatal(err)
	}

	return query
}

func GenerateNewURLforProtectionJobs(startTimeUsecs, endTimeUsecs, jobId string) (Result string) {

	baseurl := "https://10.14.19.226/v2/data-protect/protection-groups/2790138600742128%3A1647109707001%3A"
	resturl := "/runs?useCachedData=false&startTimeUsecs=1654671600000000&numRuns=360&includeTenants=true&includeObjectDetails=false&endTimeUsecs=1655276399999000"
	newURL := baseurl + jobId + resturl
	result, _ := url.Parse(newURL)

	q := result.Query()
	q.Set("useCacheData", "false")
	q.Set("startTimeUsecs", startTimeUsecs)
	q.Set("numRuns", "360")
	q.Set("includeTenants", "true")
	q.Set("includeObjectDetails", "false")
	q.Set("endTimeUsecs", endTimeUsecs)
	result.RawQuery = q.Encode()

	return result.String()
}

func GenerateNewURLforAlerts(startTimeUsecs, endTimeUsecs string) (Result string) {
	URL := "https://10.14.19.226/irisservices/api/v1/public/alerts?maxAlerts=1000&startDateUsecs=1655276400000000&endDateUsecs=1655881199999000&alertStateList=kOpen&alertSeverityList=kCritical"

	result, _ := url.Parse(URL)

	q := result.Query()
	q.Set("startDateUsecs", startTimeUsecs)
	q.Set("endDateUsecs", endTimeUsecs)
	result.RawQuery = q.Encode()

	return result.String()

}
