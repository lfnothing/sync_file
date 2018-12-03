package sync_file

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
)

const (
	register_sync_file_path1 = "syncFileService1.json"
	register_sync_file_path2 = "syncFileService2.json"
)

func init() {
	RegisterSyncFile(register_sync_file_path1, false)
	RegisterSyncFile(register_sync_file_path2, true)
}

func TestSyncFileService_Operation_Write(t *testing.T) {
	go SyncFileServiceStart()

	register1FileData := []TestDataStruct{testData1, testData2, testData3, testData4}
	register1FileBytes, _ := json.MarshalIndent(&register1FileData, "", "  ")
	register1FileOperation := SyncFileServiceData{
		Data:      register1FileBytes,
		Filepath:  register_sync_file_path1,
		Operation: SyncFileWrite,
	}

	register2FileBytes1, _ := json.MarshalIndent(&testData1, "", "  ")
	register2FileBytes2, _ := json.MarshalIndent(&testData2, "", "  ")
	register2FileBytes3, _ := json.MarshalIndent(&testData3, "", "  ")
	register2FileBytes4, _ := json.MarshalIndent(&testData4, "", "  ")
	register2FileOperation1 := SyncFileServiceData{
		Data:      register2FileBytes1,
		Filepath:  register_sync_file_path2,
		Operation: SyncFileWrite,
	}
	register2FileOperation2 := SyncFileServiceData{
		Data:      register2FileBytes2,
		Filepath:  register_sync_file_path2,
		Operation: SyncFileWrite,
	}
	register2FileOperation3 := SyncFileServiceData{
		Data:      register2FileBytes3,
		Filepath:  register_sync_file_path2,
		Operation: SyncFileWrite,
	}
	register2FileOperation4 := SyncFileServiceData{
		Data:      register2FileBytes4,
		Filepath:  register_sync_file_path2,
		Operation: SyncFileWrite,
	}

	var wg sync.WaitGroup
	wg.Add(5)
	go func() {
		defer wg.Done()
		Operation(register1FileOperation)
		register1Result := Result(register_sync_file_path1)
		if register1Result.Err != nil {
			fmt.Printf("Write %s error\n", register_sync_file_path1)
		}
	}()

	go func() {
		defer wg.Done()
		Operation(register2FileOperation1)
		register2Result1 := Result(register_sync_file_path2)
		if register2Result1.Err != nil {
			fmt.Printf("Write %s first error\n", register_sync_file_path2)
		}
	}()

	go func() {
		defer wg.Done()
		Operation(register2FileOperation2)
		register2Result2 := Result(register_sync_file_path2)
		if register2Result2.Err != nil {
			fmt.Printf("Write %s  second\n", register_sync_file_path2)
		}
	}()

	go func() {
		defer wg.Done()
		Operation(register2FileOperation3)
		register2Result3 := Result(register_sync_file_path2)
		if register2Result3.Err != nil {
			fmt.Printf("Write %s third error\n", register_sync_file_path2)
		}
	}()

	go func() {
		defer wg.Done()
		Operation(register2FileOperation4)
		register2Result4 := Result(register_sync_file_path2)
		if register2Result4.Err != nil {
			fmt.Printf("Write %s forth error\n", register_sync_file_path2)
		}
	}()
	wg.Wait()
	fmt.Println("Write finsh...")
}

func TestSyncFileService_Operation_Read(t *testing.T) {
	TestSyncFileService_Operation_Write(t)

	var register1FileBytes []byte
	var register2FileBytes1 []byte
	var register2FileBytes2 []byte
	var register2FileBytes3 []byte
	var register2FileBytes4 []byte
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		register1FileOperation := SyncFileServiceData{
			Filepath:  register_sync_file_path1,
			Operation: SyncFileRead,
		}
		Operation(register1FileOperation)
		register1Result := Result(register_sync_file_path1)
		if register1Result.Err != nil {
			fmt.Printf("Read %s error\n", register_sync_file_path1)
		}
		register1FileBytes = register1Result.Data
	}()

	var offset int64
	go func() {
		defer wg.Done()

		register2FileOperation1 := SyncFileServiceData{
			Offset:    offset,
			Filepath:  register_sync_file_path2,
			Operation: SyncFileRead,
		}
		Operation(register2FileOperation1)
		register2Result1 := Result(register_sync_file_path2)
		if register2Result1.Err != nil {
			fmt.Printf("Read %s error\n", register_sync_file_path1)
		}
		offset += ChainFileEntryOffset + int64(len(register2Result1.Data))
		register2FileBytes1 = register2Result1.Data

		register2FileOperation2 := SyncFileServiceData{
			Offset:    offset,
			Filepath:  register_sync_file_path2,
			Operation: SyncFileRead,
		}
		Operation(register2FileOperation2)
		register2Result2 := Result(register_sync_file_path2)
		if register2Result2.Err != nil {
			fmt.Printf("Read %s error\n", register_sync_file_path1)
		}
		offset += ChainFileEntryOffset + int64(len(register2Result1.Data))
		register2FileBytes2 = register2Result2.Data

		register2FileOperation3 := SyncFileServiceData{
			Offset:    offset,
			Filepath:  register_sync_file_path2,
			Operation: SyncFileRead,
		}
		Operation(register2FileOperation3)
		register2Result3 := Result(register_sync_file_path2)
		if register2Result3.Err != nil {
			fmt.Printf("Read %s error\n", register_sync_file_path1)
		}
		offset += ChainFileEntryOffset + int64(len(register2Result1.Data))
		register2FileBytes3 = register2Result3.Data

		register2FileOperation4 := SyncFileServiceData{
			Offset:    offset,
			Filepath:  register_sync_file_path2,
			Operation: SyncFileRead,
		}
		Operation(register2FileOperation4)
		register2Result4 := Result(register_sync_file_path2)
		if register2Result4.Err != nil {
			fmt.Printf("Read %s error\n", register_sync_file_path1)
		}
		offset += ChainFileEntryOffset + int64(len(register2Result1.Data))
		register2FileBytes4 = register2Result4.Data
	}()
	wg.Wait()

	var register1FileData []TestDataStruct
	json.Unmarshal(register1FileBytes, &register1FileData)

	var register2FileData1 TestDataStruct
	var register2FileData2 TestDataStruct
	var register2FileData3 TestDataStruct
	var register2FileData4 TestDataStruct
	json.Unmarshal(register2FileBytes1, &register2FileData1)
	json.Unmarshal(register2FileBytes2, &register2FileData2)
	json.Unmarshal(register2FileBytes3, &register2FileData3)
	json.Unmarshal(register2FileBytes4, &register2FileData4)

	for _, v := range register1FileData {
		if v == register2FileData1 {
			fmt.Println("Test data1 match")
		}
		if v == register2FileData2 {
			fmt.Println("Test data2 match")
		}
		if v == register2FileData3 {
			fmt.Println("Test data3 match")
		}
		if v == register2FileData4 {
			fmt.Println("Test data4 match")
		}
	}
}

func TestSyncFileService_Operation_Cut(t *testing.T) {
	TestSyncFileService_Operation_Write(t)

	var register1FileBytes []byte
	var register2FileBytes1 []byte
	var register2FileBytes2 []byte
	var register2FileBytes3 []byte
	var register2FileBytes4 []byte
	var wg sync.WaitGroup
	wg.Add(5)
	go func() {
		defer wg.Done()
		register1FileOperation := SyncFileServiceData{
			Filepath:  register_sync_file_path1,
			Operation: SyncFileCut,
		}
		Operation(register1FileOperation)
		register1Result := Result(register_sync_file_path1)
		if register1Result.Err != nil {
			fmt.Printf("Read %s error\n", register_sync_file_path1)
		}
		register1FileBytes = register1Result.Data
	}()

	go func() {
		defer wg.Done()
		register2FileOperation1 := SyncFileServiceData{
			Filepath:  register_sync_file_path2,
			Operation: SyncFileCut,
		}
		Operation(register2FileOperation1)
		register2Result1 := Result(register_sync_file_path2)
		if register2Result1.Err != nil {
			fmt.Printf("Read %s error\n", register_sync_file_path2)
		}
		register2FileBytes1 = register2Result1.Data
	}()

	go func() {
		defer wg.Done()
		register2FileOperation2 := SyncFileServiceData{
			Filepath:  register_sync_file_path2,
			Operation: SyncFileCut,
		}
		Operation(register2FileOperation2)
		register2Result2 := Result(register_sync_file_path2)
		if register2Result2.Err != nil {
			fmt.Printf("Read %s error\n", register_sync_file_path2)
		}
		register2FileBytes2 = register2Result2.Data
	}()

	go func() {
		defer wg.Done()
		register2FileOperation3 := SyncFileServiceData{
			Filepath:  register_sync_file_path2,
			Operation: SyncFileCut,
		}
		Operation(register2FileOperation3)
		register2Result3 := Result(register_sync_file_path2)
		if register2Result3.Err != nil {
			fmt.Printf("Read %s error\n", register_sync_file_path2)
		}
		register2FileBytes3 = register2Result3.Data
	}()

	go func() {
		defer wg.Done()
		register2FileOperation4 := SyncFileServiceData{
			Filepath:  register_sync_file_path2,
			Operation: SyncFileCut,
		}
		Operation(register2FileOperation4)
		register2Result4 := Result(register_sync_file_path2)
		if register2Result4.Err != nil {
			fmt.Printf("Read %s error\n", register_sync_file_path2)
		}
		register2FileBytes4 = register2Result4.Data
	}()
	wg.Wait()

	var register1FileData []TestDataStruct
	json.Unmarshal(register1FileBytes, &register1FileData)

	var register2FileData1 TestDataStruct
	var register2FileData2 TestDataStruct
	var register2FileData3 TestDataStruct
	var register2FileData4 TestDataStruct
	json.Unmarshal(register2FileBytes1, &register2FileData1)
	json.Unmarshal(register2FileBytes2, &register2FileData2)
	json.Unmarshal(register2FileBytes3, &register2FileData3)
	json.Unmarshal(register2FileBytes4, &register2FileData4)

	for _, v := range register1FileData {
		if v == register2FileData1 {
			fmt.Println("Test data1 match")
		}
		if v == register2FileData2 {
			fmt.Println("Test data2 match")
		}
		if v == register2FileData3 {
			fmt.Println("Test data3 match")
		}
		if v == register2FileData4 {
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
