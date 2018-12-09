package sync_file

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
)

type TestDataStruct struct {
	Id   int    `json:"id"`
	Data string `json:"data"`
}

var (
	testData1 = TestDataStruct{1, "1"}
	testData2 = TestDataStruct{2, "2"}
	testData3 = TestDataStruct{3, "3"}
	testData4 = TestDataStruct{4, "4"}

	syncFile      = NewSyncFile("syncfile.json", false)
	syncChainFile = NewSyncFile("syncChainFile.json", true)
)

func TestSyncFile_Write(t *testing.T) {
	syncFileData := []TestDataStruct{testData1, testData2, testData3, testData4}
	syncFileBytes, _ := json.MarshalIndent(&syncFileData, "", "  ")

	syncChainFileBytes1, _ := json.MarshalIndent(&testData1, "", "  ")
	syncChainFileBytes2, _ := json.MarshalIndent(&testData2, "", "  ")
	syncChainFileBytes3, _ := json.MarshalIndent(&testData3, "", "  ")
	syncChainFileBytes4, _ := json.MarshalIndent(&testData4, "", "  ")

	var wg sync.WaitGroup
	wg.Add(5)
	go func() {
		defer wg.Done()
		syncFile.Write(true, syncFileBytes)
	}()

	// chain file write
	go func() {
		defer wg.Done()
		syncChainFile.Write(false, syncChainFileBytes1)
	}()
	go func() {
		defer wg.Done()
		syncChainFile.Write(false, syncChainFileBytes2)
	}()
	go func() {
		defer wg.Done()
		syncChainFile.Write(false, syncChainFileBytes3)
	}()
	go func() {
		defer wg.Done()
		syncChainFile.Write(false, syncChainFileBytes4)
	}()

	wg.Wait()
	syncChainFile.Write(true, syncChainFileBytes4)
	syncChainFile.Write(true, syncChainFileBytes3)
	syncChainFile.Write(true, syncChainFileBytes2)
	syncChainFile.Write(true, syncChainFileBytes1)
}

func TestSyncFile_Read(t *testing.T) {
	TestSyncFile_Write(t)

	var syncFileBytes []byte
	var syncChainFileBytes1 []byte
	var syncChainFileBytes2 []byte
	var syncChainFileBytes3 []byte
	var syncChainFileBytes4 []byte
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		syncFileBytes, _ = syncFile.Read(0)
	}()

	go func() {
		defer wg.Done()

		var offset int64
		syncChainFileBytes1, _ = syncChainFile.Read(offset)
		offset += ChainFileEntryOffset + int64(len(syncChainFileBytes1))

		syncChainFileBytes2, _ = syncChainFile.Read(offset)
		offset += ChainFileEntryOffset + int64(len(syncChainFileBytes2))

		syncChainFileBytes3, _ = syncChainFile.Read(offset)
		offset += ChainFileEntryOffset + int64(len(syncChainFileBytes3))

		syncChainFileBytes4, _ = syncChainFile.Read(offset)
		offset += ChainFileEntryOffset + int64(len(syncChainFileBytes4))
	}()
	wg.Wait()

	var syncFileData []TestDataStruct
	json.Unmarshal(syncFileBytes, &syncFileData)

	var syncChainFileData1 TestDataStruct
	var syncChainFileData2 TestDataStruct
	var syncChainFileData3 TestDataStruct
	var syncChainFileData4 TestDataStruct
	json.Unmarshal(syncChainFileBytes1, &syncChainFileData1)
	json.Unmarshal(syncChainFileBytes2, &syncChainFileData2)
	json.Unmarshal(syncChainFileBytes3, &syncChainFileData3)
	json.Unmarshal(syncChainFileBytes4, &syncChainFileData4)

	for _, v := range syncFileData {
		if v == syncChainFileData1 {
			fmt.Println("Test data1 match")
		}
		if v == syncChainFileData2 {
			fmt.Println("Test data2 match")
		}
		if v == syncChainFileData3 {
			fmt.Println("Test data3 match")
		}
		if v == syncChainFileData4 {
			fmt.Println("Test data4 match")
		}
	}
}

func TestSyncFile_Cut(t *testing.T) {
	TestSyncFile_Write(t)

	var syncFileBytes []byte
	var syncChainFileBytes1 []byte
	var syncChainFileBytes2 []byte
	var syncChainFileBytes3 []byte
	var syncChainFileBytes4 []byte
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		syncFileBytes, _ = syncFile.Cut()
	}()

	go func() {
		defer wg.Done()
		syncChainFileBytes1, _ = syncChainFile.Cut()
		syncChainFileBytes2, _ = syncChainFile.Cut()
		syncChainFileBytes3, _ = syncChainFile.Cut()
		syncChainFileBytes4, _ = syncChainFile.Cut()
		syncChainFileBytes1, _ = syncChainFile.Cut()
		syncChainFileBytes2, _ = syncChainFile.Cut()
		syncChainFileBytes3, _ = syncChainFile.Cut()
		syncChainFileBytes4, _ = syncChainFile.Cut()
	}()
	wg.Wait()

	var syncFileData []TestDataStruct
	json.Unmarshal(syncFileBytes, &syncFileData)

	var syncChainFileData1 TestDataStruct
	var syncChainFileData2 TestDataStruct
	var syncChainFileData3 TestDataStruct
	var syncChainFileData4 TestDataStruct
	json.Unmarshal(syncChainFileBytes1, &syncChainFileData1)
	json.Unmarshal(syncChainFileBytes2, &syncChainFileData2)
	json.Unmarshal(syncChainFileBytes3, &syncChainFileData3)
	json.Unmarshal(syncChainFileBytes4, &syncChainFileData4)

	for _, v := range syncFileData {
		if v == syncChainFileData1 {
			fmt.Println("Test data1 match")
		}
		if v == syncChainFileData2 {
			fmt.Println("Test data2 match")
		}
		if v == syncChainFileData3 {
			fmt.Println("Test data3 match")
		}
		if v == syncChainFileData4 {
			fmt.Println("Test data4 match")
		}
	}

	if syncFile.getFileSize() != 0 {
		fmt.Println("Sync file cut data error")
	}

	if syncChainFile.getFileSize() != 0 {
		fmt.Println("Sync chain file cut data error")
	}
}
