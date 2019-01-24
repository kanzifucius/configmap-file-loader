package service

import (
	"github.com/kanzifucius/configmap_file_loader/pkg/log"
	gocache "github.com/patrickmn/go-cache"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Echo is simple echo service.
type ConfigMapFileWriter interface {
	// EchoObj echoes the received object.
	WriteConfigMap(obj runtime.Object) error
	DeleteConfigMap(objKey string) error
}

// ConfigMapFileWriter echoes the received object name.
type ConfigMapFileWriterSrv struct {
	logger         log.Logger
	directory      string
	keyNameFilter  string
	contollerCache *gocache.Cache
	webhookMethod string
	webhookStatusCode int
	webhook string

}


func NewConfigMapFileWriter(logger log.Logger, directory ,keyNameFilter ,webhook,webhookMethod string,webhookStatusCode int) *ConfigMapFileWriterSrv {
	return &ConfigMapFileWriterSrv{
		logger:    logger,
		directory: directory,
		keyNameFilter: keyNameFilter,
		contollerCache: gocache.New(5*time.Minute, 10*time.Minute),
		webhook: webhook,
		webhookMethod: webhookMethod,
		webhookStatusCode: webhookStatusCode,


	}
}



func (s *ConfigMapFileWriterSrv) DeleteConfigMap(objKey string) error {


	sanitizedName := strings.Replace(objKey, "/", "-", -1)
	pattern := join(s.directory, "/", sanitizedName, "*")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	for _, f := range files {
		s.logger.Infof("Deleting Files form ConfigMap %s entry -> %s", objKey, f)
		if err := os.Remove(f); err != nil {
			return err
		}
	}
	s.Notify()
	return nil
}

func join(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}

func (s *ConfigMapFileWriterSrv) WriteConfigMap(obj runtime.Object)  error {


	cmap, ok := obj.(*corev1.ConfigMap)
	if !ok {
		s.logger.Errorf("Could not get config map")
	}


	key,err :=  cache.MetaNamespaceKeyFunc(obj) ; if err !=nil{

		return err
	}

	s.logger.Infof("*************** [%s] ConfigMap Detected *********************** ", key)
	processed ,err := s.isVersionProcessed(cmap)
	if 	processed {
		s.logger.Infof("Skipping ConfigMap [%s] , version already processed ", key)
		return nil
	}

	err = s.DeleteConfigMap(key);if err !=nil {
		s.logger.Errorf("%+v", err)
		return err
	}
	err = s.writeDataKeysAsFiles(cmap);if err !=nil {
		s.logger.Errorf("%+v", err)
		return err
	}
	s.contollerCache.Set(key, cmap.ResourceVersion, gocache.NoExpiration)
	s.Notify()
	return nil
}


func (s *ConfigMapFileWriterSrv) writeDataKeysAsFiles(cmap *corev1.ConfigMap  ) error {

	key,err :=  cache.MetaNamespaceKeyFunc(cmap) ; if err !=nil{
		return err
	}

	for dataKey, v := range cmap.Data {
		if !s.isFilterMatch(dataKey) {
			s.logger.Infof("Skipping Key [%s] in ConfigMap [%s] , filter match exclusion-> %s", dataKey,key, s.keyNameFilter)
			break
		}

		fileName := join(key, "-", dataKey)
		sanitizedName := strings.Replace(fileName, "/", "-", -1)
		path := filepath.Join(s.directory, sanitizedName)
  		file, err := os.Create(path) ;if err !=nil {
			s.logger.Errorf("%+v", err)
			return err
		}

		_, err = file.WriteString(v) ; if err !=nil {
			s.logger.Errorf("%+v", err)
			return err
		}

		s.logger.Infof("Written Files form ConfigMap %s entry -> %s", cmap.Name, path)
	}

	return nil
}

func (s *ConfigMapFileWriterSrv) isFilterMatch(filename string ) bool{

	match,_ := regexp.MatchString(s.keyNameFilter,filename)
	return match

}

func (s *ConfigMapFileWriterSrv) isVersionProcessed(cmap *corev1.ConfigMap ) (bool , error){

	key,err :=  cache.MetaNamespaceKeyFunc(cmap) ; if err !=nil{
		return false,err
	}

	version, found := s.contollerCache.Get(key)
	if found {
		versionString :=version.(string)
		if cmap.ResourceVersion == versionString{
			return true,nil
		}
	}

	return false,nil
}

func (s *ConfigMapFileWriterSrv) Notify(){

	if s.webhook !="" {

		req, err := http.NewRequest(s.webhookMethod, s.webhook, nil)
		if err != nil {
			s.logger.Errorf("%+v", err)

		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			s.logger.Errorf("%+v", err)

		}

		resp.Body.Close()
		if resp.StatusCode != s.webhookStatusCode {
			s.logger.Errorf("Received response code %s, expected %s",resp.StatusCode,  s.webhookStatusCode)
			return
		}
		s.logger.Infof("successfully triggered webhook notifcation")
	}
}