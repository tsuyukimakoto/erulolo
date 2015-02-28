package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/gen/cloudwatch"
)

func getInstanceId() (instanceId string) {
	resp, _ := http.Get("http://169.254.169.254/latest/meta-data/instance-id")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	instanceId = string(body)
	return
}

func noticeToAwsFileModified(data int) {
	noticeToAws("Modify File Count", "Count", float64(data))
}

func noticeToAwsDiskUse(data float64) {
	noticeToAws("Disk Use", "Percent", data)
}

func noticeToAws(metricName string, unit string, value float64) {
	//TODO region
	cli := cloudwatch.New(aws.IAMCreds(), "ap-northeast-1", nil)


	dimensionParam := &cloudwatch.Dimension{
		Name:  aws.String("InstanceId"),
		Value: aws.String(getInstanceId()),
	}

	metricDataParam := &cloudwatch.MetricDatum{
		Dimensions: []cloudwatch.Dimension{*dimensionParam},
		MetricName: aws.String(metricName),
		Unit:	   aws.String(unit),
		Value:	  aws.Double(float64(value)),
	}

	putMetricDataInput := &cloudwatch.PutMetricDataInput{
		MetricData: []cloudwatch.MetricDatum{*metricDataParam},
		Namespace:  aws.String("Erulolo"),
	}
	err := cli.PutMetricData(putMetricDataInput)
	if err != nil {
		fmt.Println("put metrics:", err)
	}
}

func main() {
	go Buoy()
	go SwellWatch()
	select {}
}
