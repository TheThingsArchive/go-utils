package roots

import (
	"testing"

	. "github.com/smartystreets/assertions"
)

func TestMozillaRootCAs(t *testing.T) {
	a := New(t)
	a.So(MozillaRootCAs.Subjects(), ShouldNotBeEmpty)
}
