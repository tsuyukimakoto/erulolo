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



func notice_to_aws(data int) {
	cli := cloudwatch.New(aws.IAMCreds(), "ap-northeast-1", nil)

	dimensionParam := &cloudwatch.Dimension{
		Name:  aws.String("InstanceId"),
		Value: aws.String(getInstanceId()),
	}

	metricDataParam := &cloudwatch.MetricDatum{
		Dimensions: []cloudwatch.Dimension{*dimensionParam},
		MetricName: aws.String("Modify File Count"),
		Unit:	   aws.String("Count"),
		Value:	  aws.Double(float64(data)),
	}

	putMetricDataInput := &cloudwatch.PutMetricDataInput{
		MetricData: []cloudwatch.MetricDatum{*metricDataParam},
		Namespace:  aws.String("Erulolo"),
	}

	fmt.Println("put metrict:", cli.PutMetricData(putMetricDataInput))
}

func main() {
	go Watch()
	select {}
	// notice_to_aws(1)
}
