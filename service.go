package sync_file

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

//--------------------------------------
// sync file service
//--------------------------------------

type SyncFileOperation int

const (
	SyncFileRead SyncFileOperation = iota
	SyncFileWrite
	SyncFileCut
	SyncFileSize
)

var (
	syncFileService = NewSyncFileService()
)

type SyncFileServiceData struct {
	Err       error
	Data      []byte
	Offset    int64
	Filepath  string
	Filesize  int64
	Operation SyncFileOperation
}

type SyncFileService struct {
	files   map[string]*SyncFile
	closing chan bool
	groups  *sync.WaitGroup
	inputs  map[string]chan SyncFileServiceData
	outputs map[string]chan SyncFileServiceData
}

func NewSyncFileService() *SyncFileService {
	return &SyncFileService{
		files:   make(map[string]*SyncFile),
		closing: make(chan bool),
		groups:  &sync.WaitGroup{},
		inputs:  make(map[string]chan SyncFileServiceData),
		outputs: make(map[string]chan SyncFileServiceData),
	}
}

func (this *SyncFileService) Register(filepath string, chain bool) {
	if _, ok := this.files[filepath]; ok {
		return
	}
	this.files[filepath] = NewSyncFile(filepath, chain)
	this.inputs[filepath] = make(chan SyncFileServiceData)
	this.outputs[filepath] = make(chan SyncFileServiceData)
}

func (this *SyncFileService) Serve() {
	defer this.groups.Done()
	for k, v := range this.files {
		this.groups.Add(1)

		// pass param
		go func(key string, syncfile *SyncFile) {
			defer this.groups.Done()
			for {
				select {
				case <-this.closing:
					return

				case input := <-this.inputs[key]:
					switch input.Operation {
					case SyncFileRead:
						data, err := syncfile.Read(input.Offset)
						serviceData := SyncFileServiceData{
							Data: data,
							Err:  err,
						}
						this.outputs[key] <- serviceData
					case SyncFileWrite:
						err := syncfile.Write(false, input.Data)
						serviceData := SyncFileServiceData{
							Err: err,
						}
						this.outputs[key] <- serviceData
					case SyncFileCut:
						data, err := syncfile.Cut()
						serviceData := SyncFileServiceData{
							Data: data,
							Err:  err,
						}
						this.outputs[key] <- serviceData
					case SyncFileSize:
						fz := syncfile.getFileSize()
						serviceData := SyncFileServiceData{
							Filesize: fz,
						}
						this.outputs[key] <- serviceData
					}
				}
			}
		}(k, v)
	}
}

func (this *SyncFileService) Stop() {
	syncFileService.closing <- true
	syncFileService.groups.Wait()
}

func (this *SyncFileService) Groups() *sync.WaitGroup {
	return this.groups
}

func Operation(input SyncFileServiceData) {
	if _, ok := syncFileService.files[input.Filepath]; !ok {
		return
	}
	syncFileService.inputs[input.Filepath] <- input
}

func Result(filepath string) SyncFileServiceData {
	return <-syncFileService.outputs[filepath]
}

func SyncFileServiceStart() {
	syncFileService.groups.Add(1)
	go syncFileService.Serve()

	// recevie program quit
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	syncFileService.Stop()
}

//--------------------------------------
// sync file register
//--------------------------------------

func RegisterSyncFile(filepath string, chain bool) {
	syncFileService.Register(filepath, chain)
}
