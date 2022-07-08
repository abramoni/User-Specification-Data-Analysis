package mypackages

import (
	"context"
	"encoding/json"
	"fmt"

	elastic "github.com/olivere/elastic/v7"
)

var Cluster_Space_Utilization , Meta_Data_Utilization, Meta_Data_Space_Percentage float64

func GetESClient() (*elastic.Client, error) {

	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))

	// fmt.Println("ES initialized...")

	return client, err

}

func StoreInElasticDb(body []byte, IndexName string) {
	ctx := context.Background()
	esclient, err := GetESClient()
	if err != nil {
		fmt.Println("Error initializing elastcidb client : ", err)
		panic("Client fail ")
	}

	jsondata_string := string(body)
	ind, err := esclient.Index().
		Index(IndexName).
		BodyJson(jsondata_string).
		Do(ctx)
	_ = ind

	if err != nil {
		panic(err)
	}
	//fmt.Println("Stored in Db....")
}

func RetriveFronElasticDb(time_stamp int64) {
	
	ctx := context.Background()
	esclient, err := GetESClient()
	if err != nil {
		fmt.Println("Error initializing : ", err)
		panic("Client fail ")
	}

	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMatchQuery("timeStamp", time_stamp))
	queryStr, err1 := searchSource.Source()
	_, err2 := json.Marshal(queryStr)

	if err1 != nil || err2 != nil {
		fmt.Println("[esclient][GetResponse]err during query marshal=", err1, err2)
	}
	

	searchService := esclient.Search().Index("demodata").SearchSource(searchSource)

	searchResult, err := searchService.Do(ctx)
	if err != nil {
		fmt.Println("Error=", err)
		return
	}

	for _, hit := range searchResult.Hits.Hits {
		
		err := json.Unmarshal(hit.Source, &Response_for_elastic)
		if err != nil {
			fmt.Println("Searching Result Err = ", err)
		}

	}

	if err != nil {
		fmt.Println("Fetching from elastic DB failed: ", err)
	} else {
		
		Cluster_Space_Utilization = Response_for_elastic.ClusterSpaceUtilization
		Meta_Data_Utilization = Response_for_elastic.MetaDataUtilization
		Meta_Data_Space_Percentage = Response_for_elastic.MetaDataSpacePercentage 
		}
	}

	func RetriveStatsFronElasticDb(id int64) {

		ctx := context.Background()
		esclient, err := GetESClient()
		if err != nil {
			fmt.Println("Error initializing : ", err)
			panic("Client fail ")
		}
	
		searchSource := elastic.NewSearchSource()
		searchSource.Query(elastic.NewMatchQuery("id", id))
		queryStr, err1 := searchSource.Source()
		_, err2 := json.Marshal(queryStr)
	
		if err1 != nil || err2 != nil {
			fmt.Println("[esclient][GetResponse]err during query marshal=", err1, err2)
		}
		//fmt.Println("[esclient]Final ESQuery=\n", string(queryJs))
	
		searchService := esclient.Search().Index("stats_test").SearchSource(searchSource)
	
		searchResult, err := searchService.Do(ctx)
		if err != nil {
			fmt.Println("Error=", err)
			return
		}
	
		for _, hit := range searchResult.Hits.Hits {
	
			err := json.Unmarshal(hit.Source, &Stats_response)
			// fmt.Println()
			 //fmt.Println("fcr ; ",Stats_response.ClusterSoftwareVersion)
			// fmt.Println()
			// fmt.Println("fcs; ",Stats_response.FileCreateSumStats.PastData)
			// fmt.Println()
			// fmt.Println("fl ; ",Stats_response.FileLatencyStats.PastData)
			// fmt.Println()
			// fmt.Println("gc ; ",Stats_response.GarbageCollectionStats.PastData)
			// fmt.Println()
			// fmt.Println("mdu ; ",Stats_response.MetaDataUtilizationStats.PastData)
			// fmt.Println()
			// fmt.Println("suc ; ",Stats_response.SystemUtilizationChangeRateStats.PastData)
			// fmt.Println()
			// fmt.Println("jbt ; ",Stats_response.JobBackUpTimeStats.PastData)
			// fmt.Println()
			// fmt.Println("jat ; ",Stats_response.JobArchivalTimeStats.PastData)
			// fmt.Println()
			// fmt.Println("jrt ; ",Stats_response.JobReplicationTimeStats.PastData)
			
			
			if err != nil {
				fmt.Println("Searching Result Err = ", err)
			}
	
		}
	
		if err != nil {
			fmt.Println("Fetching from elastic DB failed: ", err)
		} 
	}