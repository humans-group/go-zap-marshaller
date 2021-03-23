package zapmarshaller

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/mkorolyov/astparser"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

func TestGenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("zap_generator.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "ZapGenerator", []Reporter{junitReporter})
}

var _ = Describe("backoff", func() {
	It("should generate golden files", func() {
		cfg := astparser.Config{
			InputDir:      "fixtures_test",
			IncludeRegexp: "custom.go",
		}
		sources, err := astparser.Load(cfg)
		Ω(err).ShouldNot(HaveOccurred())

		generator := Generator{Cfg: Config{}}
		files := generator.Generate(sources)
		for name, got := range files {
			want, err := ioutil.ReadFile(
				fmt.Sprintf("fixtures_test/%s.zap.go", strings.Split(name, ".")[0]))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(string(want)).Should(BeEquivalentTo(string(got)))
		}
	})
})
