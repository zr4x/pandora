package netsample

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var _ = Describe("Phout", func() {
	const fileName = "out.txt"
	var (
		fs     afero.Fs
		conf   PhoutConfig
		testee Aggregator
		ctx    context.Context
		cancel context.CancelFunc
		runErr chan error
	)
	getOutput := func() string {
		data, err := afero.ReadFile(fs, fileName)
		Expect(err).NotTo(HaveOccurred())
		return string(data)
	}

	BeforeEach(func() {
		fs = afero.NewMemMapFs()
		conf = PhoutConfig{Destination: fileName}
		ctx, cancel = context.WithCancel(context.Background())
	})
	JustBeforeEach(func() {
		var err error
		testee, err = NewPhout(fs, conf)
		Expect(err).NotTo(HaveOccurred())
		runErr = make(chan error)
		go func() {
			runErr <- testee.Run(ctx)
		}()
	})
	It("no id by default", func() {
		testee.Report(newTestSample())
		testee.Report(newTestSample())
		cancel()
		Expect(<-runErr).NotTo(HaveOccurred())
		Expect(getOutput()).To(Equal(strings.Repeat(testSampleNoIdPhout+"\n", 2)))
	}, 1)
	Context("id option set", func() {
		BeforeEach(func() {
			conf.Id = true
		})
		It("id printed", func() {
			testee.Report(newTestSample())
			cancel()
			Expect(<-runErr).NotTo(HaveOccurred())
			Expect(getOutput()).To(Equal(testSamplePhout + "\n"))
		}, 1)

	})

})

const (
	testSamplePhout     = "1484660999.002	tag1|tag2#42	333333	0	0	0	0	0	0	0	13	999"
	testSampleNoIdPhout = "1484660999.002	tag1|tag2	333333	0	0	0	0	0	0	0	13	999"
)

func TestOUt(t *testing.T) {
	fmt.Println(newTestSample().String())
}

func newTestSample() *Sample {
	s := &Sample{}
	s.timeStamp = time.Unix(1484660999, 002*1000000)
	s.SetId(42)
	s.AddTag("tag1|tag2")
	s.setDuration(keyRTTMicro, time.Second/3)
	s.set(keyErrno, 13)
	s.set(keyProtoCode, ProtoCodeError)
	return s
}
