package leafnodes

import (
	"fmt"
	"reflect"
	"time"

	"github.com/onsi/ginkgo/internal/failer"
	"github.com/onsi/ginkgo/types"
)

type MeasureNode struct {
	runner *runner

	text        string
	flag        types.FlagType
	samples     int
	benchmarker *benchmarker
}

func NewMeasureNode(text string, body interface{}, flag types.FlagType, codeLocation types.CodeLocation, timeout time.Duration, samples int, failer *failer.Failer, componentIndex int) *MeasureNode {
	benchmarker := newBenchmarker()

	var wrappedBody interface{}

	bodyType := reflect.TypeOf(body)
	if bodyType.Kind() != reflect.Func {
		panic(fmt.Sprintf("Expected a function but got something else at %v", codeLocation))
	}

	switch bodyType.NumIn() {
	case 1:
		wrappedBody = func() {
			fmt.Println(reflect.ValueOf(body))
			reflect.ValueOf(body).Call([]reflect.Value{reflect.ValueOf(benchmarker)})
		}
	case 2:
		wrappedBody = func(done chan<- interface{}) {
			reflect.ValueOf(body).Call([]reflect.Value{reflect.ValueOf(benchmarker), reflect.ValueOf(done)})
		}
	default:
		panic(fmt.Sprintf("Wrong number of arguments to function at %v", codeLocation))
	}

	return &MeasureNode{
		runner: newRunner(wrappedBody, codeLocation, timeout, failer, types.SpecComponentTypeMeasure, componentIndex),

		text:        text,
		flag:        flag,
		samples:     samples,
		benchmarker: benchmarker,
	}
}

func (node *MeasureNode) Run() (outcome types.SpecState, failure types.SpecFailure) {
	return node.runner.run()
}

func (node *MeasureNode) MeasurementsReport() map[string]*types.SpecMeasurement {
	return node.benchmarker.measurementsReport()
}

func (node *MeasureNode) Type() types.SpecComponentType {
	return types.SpecComponentTypeMeasure
}

func (node *MeasureNode) Text() string {
	return node.text
}

func (node *MeasureNode) Flag() types.FlagType {
	return node.flag
}

func (node *MeasureNode) CodeLocation() types.CodeLocation {
	return node.runner.codeLocation
}

func (node *MeasureNode) Samples() int {
	return node.samples
}
