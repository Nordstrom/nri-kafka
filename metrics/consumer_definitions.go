package metrics

import (
	"strings"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
)

const consumerHolder = "%CONSUMER%"

// Consumer Metrics
var consumerMetricDefs = []*JMXMetricSet{
	{
		MBean:        "kafka.consumer:type=consumer-fetch-manager-metrics,client-id=" + consumerHolder,
		MetricPrefix: "kafka.consumer:type=consumer-fetch-manager-metrics,client-id=" + consumerHolder,
		MetricDefs: []*MetricDefinition{
			{
				Name:       "consumer.bytesInPerSecond",
				SourceType: metric.RATE,
				JMXAttr:    "attr=consumed-rate",
			},
			{
				Name:       "consumer.fetchPerSecond",
				SourceType: metric.RATE,
				JMXAttr:    "attr=fetch-rate",
			},
			{
				Name:       "consumer.maxLag",
				SourceType: metric.GAUGE,
				JMXAttr:    "attr=records-lag-max",
			},
			{
				Name:       "consumer.MessageConsumptionPerSecond",
				SourceType: metric.RATE,
				JMXAttr:    "attr=records-consumed-rate",
			},
		},
	},
	{
		MBean:        "kafka.consumer:type=ZookeeperConsumerConnector,name=*,clientId=" + consumerHolder,
		MetricPrefix: "kafka.consumer:type=ZookeeperConsumerConnector,",
		MetricDefs: []*MetricDefinition{
			{
				Name:       "consumer.offsetKafkaCommitsPerSecond",
				SourceType: metric.RATE,
				JMXAttr:    "name=KafkaCommitsPerSec,clientId=" + consumerHolder + "attr=Count",
			},
			{
				Name:       "consumer.offsetZooKeeperCommitsPerSecond",
				SourceType: metric.RATE,
				JMXAttr:    "name=ZooKeeperCommitsPerSec,clientId=" + consumerHolder + "attr=Count",
			},
		},
	},
}

// ConsumerTopicMetricDefs metric definitions for topic metrics that are specific to a Consumer
var ConsumerTopicMetricDefs = []*JMXMetricSet{
	{
		MBean:        "kafka.consumer:type=consumer-fetch-manager-metrics,client-id=" + consumerHolder + ",topic=*",
		MetricPrefix: "kafka.consumer:type=consumer-fetch-manager-metrics,client-id=" + consumerHolder + ",topic=" + topicHolder + ",",
		MetricDefs: []*MetricDefinition{
			{
				Name:       "consumer.avgFetchSizeInBytes",
				SourceType: metric.GAUGE,
				JMXAttr:    "attr=fetch-size-avg",
			},
			{
				Name:       "consumer.maxFetchSizeInBytes",
				SourceType: metric.GAUGE,
				JMXAttr:    "attr=fetch-size-max",
			},
			{
				Name:       "consumer.avgRecordConsumedPerTopicPerSecond",
				SourceType: metric.RATE,
				JMXAttr:    "attr=records-consumed-rate",
			},
			{
				Name:       "consumer.avgRecordConsumedPerTopic",
				SourceType: metric.GAUGE,
				JMXAttr:    "attr=records-per-request-avg",
			},
		},
	},
}

// applyConsumerName to be used when passed to CollectMetricDefinitions to modified bean name for Consumer
func applyConsumerName(consumerName string) func(string) string {
	return func(beanName string) string {
		return strings.Replace(beanName, consumerHolder, consumerName, -1)
	}
}

// ApplyconsumerTopicName to be used when passed to CollectMetricDefinitions to modified bean name
// for Consumer and Topic
func ApplyconsumerTopicName(consumerName, topicName string) func(string) string {
	return func(beanName string) string {
		modifiedBeanName := strings.Replace(beanName, consumerHolder, consumerName, -1)
		return strings.Replace(modifiedBeanName, topicHolder, topicName, -1)
	}
}