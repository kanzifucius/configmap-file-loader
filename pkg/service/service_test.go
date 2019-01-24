package service_test

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kanzifucius/configmap_file_loader/pkg/service"
)

type logKind int

const (
	infoKind logKind = iota
	warnignKind
	errorKind
)

type logEvent struct {
	kind logKind
	line string
}

type testLogger struct {
	events []logEvent
	sync.Mutex
}

func (t *testLogger) logLine(kind logKind, format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	t.events = append(t.events, logEvent{kind: kind, line: str})
}

func (t *testLogger) Infof(format string, args ...interface{}) {
	t.logLine(infoKind, format, args...)
}
func (t *testLogger) Warningf(format string, args ...interface{}) {
	t.logLine(warnignKind, format, args...)
}
func (t *testLogger) Errorf(format string, args ...interface{}) {
	t.logLine(errorKind, format, args...)
}



func TestWriteConfigMapTest(t *testing.T) {
	tests := []struct {
		name       string
		prefix     string
		obj        runtime.Object
		expResults []logEvent
	}{
		{
			name:   "Logging a pod should print pod name.",
			prefix: "test",
			obj: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name: "Configmap",

				},Data: map[string]string{"test.rules": "testcontent", },
			},
			expResults: []logEvent{
				logEvent{kind: infoKind, line: "Creating Files form ConfigMap Configmap entry -> /home/kanzi/temp/test/Configmap-test.rules"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// Mocks.
			ml := &testLogger{events: []logEvent{}}

			// Create aservice and run.
			srv := service.NewConfigMapFileWriter(ml,"/home/kanzi/temp/test/","","","",200)
			srv.WriteConfigMap(test.obj)

			// Check.
			assert.Equal(test.expResults, ml.events)
		})
	}
}

func TestDeleteConfigMapTest(t *testing.T) {
	tests := []struct {
		name       string
		prefix     string
		obj        runtime.Object
		expResults []logEvent
	}{
		{
			name:   "Deleting ConfigMap",
			prefix: "test",
			obj: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name: "Configmap",

				},Data: map[string]string{"test.rules": "testcontent", },
			},
			expResults: []logEvent{

				logEvent{kind: infoKind, line: "Config Map [configmaptest] ->  Deletion Detected "},
				logEvent{kind: infoKind, line: "Deleting Files form ConfigMap configmaptest entry -> /home/kanzi/temp/test/configmaptest-datatest"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// Mocks.
			ml := &testLogger{events: []logEvent{}}


			// Create aservice and run.
			srv := service.NewConfigMapFileWriter(ml,"/home/kanzi/temp/test/","*","","",200)
			_, _ = os.Create("/home/kanzi/temp/test/configmaptest-datatest")
			_ = srv.DeleteConfigMap("configmaptest")

			// Check.
			assert.Equal(test.expResults, ml.events)
		})
	}
}
