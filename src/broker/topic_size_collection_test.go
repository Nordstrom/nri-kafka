package broker

import (
	"errors"
	"testing"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-kafka/src/connection"
	"github.com/newrelic/nri-kafka/src/connection/mocks"
	"github.com/newrelic/nri-kafka/src/jmxwrapper"
	"github.com/newrelic/nri-kafka/src/testutils"
	"github.com/stretchr/testify/assert"
)

func TestGatherTopicSize_Single(t *testing.T) {
	testutils.SetupJmxTesting()
	testutils.SetupTestArgs()

	i, _ := integration.New("test", "1.0.0")

	jmxwrapper.JMXQuery = func(query string, timeout int) (map[string]interface{}, error) {
		return map[string]interface{}{
			"one":   float64(1),
			"two":   float64(2),
			"three": float64(3),
			"four":  float64(4),
		}, nil
	}

	mockBroker := &mocks.SaramaBroker{}
	mockBroker.On("Addr").Return("kafkabroker:9090")

	broker := &connection.Broker{
		Host:         "localhost",
		JMXPort:      9999,
		SaramaBroker: mockBroker,
	}

	e, _ := broker.Entity(i)
	collectedTopics := map[string]*metric.Set{
		"topic": e.NewMetricSet("KafkaBrokerSample",
			metric.Attribute{Key: "displayName", Value: "testEntity"},
			metric.Attribute{Key: "entityName", Value: "broker:testEntity"},
			metric.Attribute{Key: "topic", Value: "topic"},
		),
	}

	gatherTopicSizes(broker, collectedTopics, i)

	expected := map[string]interface{}{
		"topic.diskSize": float64(10),
		"event_type":     "KafkaBrokerSample",
		"entityName":     "broker:testEntity",
		"displayName":    "testEntity",
		"topic":          "topic",
	}

	entity, err := broker.Entity(i)
	assert.NoError(t, err)
	assert.Len(t, entity.Metrics, 1)
	assert.Equal(t, expected, entity.Metrics[0].Metrics)
}

func TestGatherTopicSize_QueryError(t *testing.T) {
	testutils.SetupJmxTesting()
	testutils.SetupTestArgs()

	i, _ := integration.New("test", "1.0.0")

	jmxwrapper.JMXQuery = func(query string, timeout int) (map[string]interface{}, error) { return nil, errors.New("error") }

	mockBroker := &mocks.SaramaBroker{}
	mockBroker.On("Addr").Return("kafkabroker:9090")

	broker := &connection.Broker{
		Host:         "localhost",
		JMXPort:      9999,
		SaramaBroker: mockBroker,
	}

	e, _ := broker.Entity(i)

	collectedTopics := map[string]*metric.Set{
		"topic": e.NewMetricSet("KafkaBrokerSample",
			metric.Attribute{Key: "displayName", Value: "testEntity"},
			metric.Attribute{Key: "entityName", Value: "broker:testEntity"},
			metric.Attribute{Key: "topic", Value: "topic"},
		),
	}

	gatherTopicSizes(broker, collectedTopics, i)

	assert.Len(t, e.Metrics, 1)
	assert.NotContains(t, e.Metrics[0].Metrics, "topic.diskSize", "Metric was unexpectedly set after query error")
}

func TestGatherTopicSize_QueryBlank(t *testing.T) {
	testutils.SetupJmxTesting()
	testutils.SetupTestArgs()

	i, _ := integration.New("test", "1.0.0")

	jmxwrapper.JMXQuery = func(query string, timeout int) (map[string]interface{}, error) {
		return make(map[string]interface{}), nil
	}

	mockBroker := &mocks.SaramaBroker{}
	mockBroker.On("Addr").Return("kafkabroker:9090")

	broker := &connection.Broker{
		Host:         "localhost",
		JMXPort:      9999,
		SaramaBroker: mockBroker,
	}

	e, _ := broker.Entity(i)

	collectedTopics := map[string]*metric.Set{
		"topic": e.NewMetricSet("KafkaBrokerSample",
			metric.Attribute{Key: "displayName", Value: "testEntity"},
			metric.Attribute{Key: "entityName", Value: "broker:testEntity"},
			metric.Attribute{Key: "topic", Value: "topic"},
		),
	}

	gatherTopicSizes(broker, collectedTopics, i)

	assert.Len(t, e.Metrics, 1)
	assert.NotContains(t, e.Metrics[0].Metrics, "topic.diskSize", "Metric was unexpectedly set after empty query result")
}

func TestGatherTopicSize_AggregateError(t *testing.T) {
	testutils.SetupJmxTesting()
	testutils.SetupTestArgs()

	i, _ := integration.New("test", "1.0.0")

	jmxwrapper.JMXQuery = func(query string, timeout int) (map[string]interface{}, error) {
		return map[string]interface{}{
			"one":  "nope",
			"four": float64(4),
		}, nil
	}

	mockBroker := &mocks.SaramaBroker{}
	mockBroker.On("Addr").Return("kafkabroker:9090")

	broker := &connection.Broker{
		Host:         "localhost",
		JMXPort:      9999,
		SaramaBroker: mockBroker,
	}

	e, _ := broker.Entity(i)

	collectedTopics := map[string]*metric.Set{
		"topic": e.NewMetricSet("KafkaBrokerSample",
			metric.Attribute{Key: "displayName", Value: "testEntity"},
			metric.Attribute{Key: "entityName", Value: "broker:testEntity"},
			metric.Attribute{Key: "topic", Value: "topic"},
		),
	}

	gatherTopicSizes(broker, collectedTopics, i)

	assert.Len(t, e.Metrics, 1)
	assert.NotContains(t, e.Metrics[0].Metrics, "topic.diskSize", "Metric was unexpectedly set after aggregate error")
}
