package txlog

import (
	"testing"

	. "github.com/adjust/gocheck"
)

func TestContextSuite(t *testing.T) {
	TestingSuiteT(&ContextSuite{}, t)
}

type ContextSuite struct{}

func (suite *ContextSuite) TestContextUpdate(c *C) {
	ctx1 := Context{
		"a": 1,
		"b": "string",
	}
	ctx2 := Context{
		"c": false,
	}

	ctx1.Update(ctx2)

	c.Check(ctx2, HasLen, 1)
	c.Check(ctx2["c"], Equals, false)

	c.Check(ctx1, HasLen, 3)
	c.Check(ctx1["a"], Equals, 1)
	c.Check(ctx1["b"], Equals, "string")
	c.Check(ctx1["c"], Equals, false)
}

func (suite *ContextSuite) TestContextUpdateWithNil(c *C) {
	ctx := Context{
		"a": 1,
		"b": "string",
	}

	ctx.Update(nil)

	c.Check(ctx, HasLen, 2)
	c.Check(ctx["a"], Equals, 1)
	c.Check(ctx["b"], Equals, "string")
}

func (suite *ContextSuite) TestContextUpdateWithConflict(c *C) {
	ctx1 := Context{
		"a": 1,
		"b": "string",
	}
	ctx2 := Context{
		"a": false,
	}

	ctx1.Update(ctx2)

	c.Check(ctx1, HasLen, 2)
	c.Check(ctx1["a"], Equals, false)
	c.Check(ctx1["b"], Equals, "string")
}
