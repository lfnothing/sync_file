package sync_file

import (
	"os"
	"sync"
	"unsafe"
)

//--------------------------------------
// sync file
//--------------------------------------

const (
	ChainFileEntryOffset = int64(unsafe.Sizeof(int64(0)))
)

type SyncFile struct {
	lock     *sync.RWMutex
	chain    bool
	filepath string
	filesize int64
}

func NewSyncFile(filepath string, chain bool) (this *SyncFile) {
	this = &SyncFile{
		lock:     &sync.RWMutex{},
		chain:    chain,
		filepath: filepath,
		filesize: 0,
	}

	if f, err := os.Stat(filepath); err == nil {
		this.filesize = f.Size()
	}
	return
}

func (this *SyncFile) SetChain(chain bool) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.setChain(chain)
}

func (this *SyncFile) GetChain() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.getChain()
}

func (this *SyncFile) SetFileSize(size int64) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.setFileSize(size)
}

func (this *SyncFile) GetFileSize() int64 {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.getFileSize()
}

func (this *SyncFile) setChain(chain bool) {
	this.chain = chain
}

func (this *SyncFile) getChain() bool {
	return this.chain
}

func (this *SyncFile) setFileSize(size int64) {
	this.filesize = size
}

func (this *SyncFile) getFileSize() int64 {
	return this.filesize
}

func (this *SyncFile) Write(head bool, data ...[]byte) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	if !this.getChain() {
		return this.write(data...)
	}
	if head == true {
		return this.insert(data...)
	}
	return this.append(data...)
}

func (this *SyncFile) write(data ...[]byte) (err error) {
	var file *os.File
	if file, err = os.OpenFile(this.filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600); err != nil {
		return
	}
	var list []byte
	for i := 0; i < len(data); i++ {
		list = append(list, data[i]...)
	}
	file.Write(list)
	file.Sync()
	file.Close()
	this.setFileSize(int64(len(list)))
	return
}

func (this *SyncFile) append(data ...[]byte) (err error) {
	var file *os.File
	fz := this.getFileSize()
	if file, err = os.OpenFile(this.filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600); err != nil {
		return
	}
	var list []byte
	for i := 0; i < len(data); i++ {
		list = append(list, Int64ToBytes(int64(len(data[i])))...)
		list = append(list, data[i]...)
	}
	file.Write(list)
	file.Sync()
	file.Close()
	this.setFileSize(fz + int64(len(list)))
	return
}

func (this *SyncFile) insert(data ...[]byte) (err error) {
	var rest []byte
	var file *os.File
	this.setChain(false)
	defer this.setChain(true)
	if this.getFileSize() != 0 {
		if rest, err = this.read(0); err != nil {
			return
		}
	}
	if file, err = os.OpenFile(this.filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600); err != nil {
		return
	}
	var list []byte
	for i := 0; i < len(data); i++ {
		list = append(list, Int64ToBytes(int64(len(data[i])))...)
		list = append(list, data[i]...)
	}
	list = append(list, rest...)
	file.Write(list)
	file.Sync()
	file.Close()
	this.setFileSize(int64(len(list)))
	return
}

func (this *SyncFile) read(offset int64) (data []byte, err error) {
	var file *os.File
	if file, err = os.Open(this.filepath); err != nil {
		return
	}

	var size int64
	if this.getChain() {
		offset += ChainFileEntryOffset
		buffer := make([]byte, ChainFileEntryOffset)
		if _, err = file.Read(buffer); err != nil {
			return
		}
		size = BytesToInt64(buffer)
	} else {
		f, _ := os.Stat(this.filepath)
		size = f.Size() - offset
	}

	data = make([]byte, size)
	_, err = file.ReadAt(data, offset)
	file.Close()
	return
}

func (this *SyncFile) Read(offset int64) (data []byte, err error) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.read(offset)
}

func (this *SyncFile) Cut() (data []byte, err error) {
	if data, err = this.Read(0); err != nil {
		return
	}

	if !this.GetChain() {
		this.Write(true, []byte{})
		return
	}

	var rest []byte
	this.SetChain(false)
	if rest, err = this.Read(ChainFileEntryOffset + int64(len(data))); err != nil {
		this.SetChain(true)
		return
	}
	err = this.Write(true, rest)
	this.SetChain(true)
	return
}
