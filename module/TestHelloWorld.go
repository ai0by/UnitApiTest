package module

import "io"
import . "gopkg.in/check.v1"
func Helloworld(c *C){
	c.Assert(42, Equals, 42)
	c.Assert(io.ErrClosedPipe, ErrorMatches, "io: .*on closed pipe")
	c.Check(42, Equals, 42)
}