package txlog

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	. "github.com/adjust/gocheck"
)

func TestJsonLoggerSuite(t *testing.T) {
	TestingSuiteT(&JsonLoggerSuite{}, t)
}

type JsonLoggerSuite struct {
	buffer *bytes.Buffer
	logger *jsonWriter
}

func (suite *JsonLoggerSuite) SetUpSuite(c *C) {
	suite.buffer = new(bytes.Buffer)
	suite.logger = NewJsonLogger(suite.buffer, "testing_service")
}

func (suite *JsonLoggerSuite) SetUpTest(c *C) {
	suite.buffer.Reset()
}

func (suite *JsonLoggerSuite) TestLog(c *C) {
	context := map[string]interface{}{
		"string":   "test",
		"int":      7,
		"datetime": time.Now(),
	}
	suite.logger.Log("testing", "Sample Message", context)

	logRecord := suite.unmarshalBuffer(c)

	c.Check(logRecord["message"], Equals, "Sample Message")
	c.Check(logRecord["tag"], Equals, "testing")
	c.Check(logRecord["timestamp"], Not(Equals), "")
	c.Check(logRecord["service"], Equals, "testing_service")
	c.Check(logRecord["host"], Matches, ".+")

	_, ok := logRecord["context"]
	c.Assert(ok, Equals, true)

	obtainedContext := Context(logRecord["context"].(map[string]interface{}))

	c.Check(obtainedContext["int"], Equals, float64(7))
	c.Check(obtainedContext["string"], Equals, "test")

	timestamp := logRecord["timestamp"].(string)
	datetime, err := time.Parse(time.RFC3339Nano, timestamp)
	c.Assert(err, IsNil, Commentf("Time string: %s", timestamp))
	c.Check(datetime.Truncate(time.Millisecond), Equals, context["datetime"].(time.Time).Truncate(time.Millisecond))
}

func (suite *JsonLoggerSuite) TestLogNoContext(c *C) {
	suite.logger.Log("testing", "Sample Message", nil)

	logRecord := suite.unmarshalBuffer(c)

	c.Check(logRecord["message"], Equals, "Sample Message")
	c.Check(logRecord["tag"], Equals, "testing")
	c.Check(logRecord["timestamp"], Not(Equals), "")
	c.Check(logRecord["service"], Equals, "testing_service")
	c.Check(logRecord["host"], Matches, ".+")
	c.Check(logRecord["context"], IsNil)
}

func (suite *JsonLoggerSuite) unmarshalBuffer(c *C) map[string]interface{} {
	c.Assert(suite.buffer.Len() > 0, Equals, true)

	output := suite.buffer.Bytes()
	logRecord := make(map[string]interface{})
	err := json.Unmarshal(output, &logRecord)
	c.Assert(err, IsNil, Commentf("JSON: %s", output))

	return logRecord
}

func BenchmarkJsonLog(b *testing.B) {
	l := NewJsonLogger(new(bytes.Buffer), "benchmark")

	for i := 0; i < b.N; i++ {
		l.Log("test", "test", Context{})
	}
}
