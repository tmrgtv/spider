package spider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type RespOpenFile struct {
	ErrorResp string `json:"error"`
	Errcode   int    `json:"errcode"`
	DocHandle int    `json:"docHandle"`
}

type RespGetTableHandle struct {
	ErrorResp   string `json:"error"`
	Errcode     int    `json:"errcode"`
	TableHandle int    `json:"tableHandle"`
}

type RespErr struct {
	ErrorResp string `json:"error"`
	Errcode   int    `json:"errcode"`
}

type TableResource struct {
	Code string `json:"Code"`
	Guid string `json:"Guid"`
	/*
		Name     string  `json:"Name"`
		Calen    string  `json:"Calen"`
		Cost_8x5 float64 `json:"c_tim_8x5_FTE"`
		Overtime float64 `json:"c_tim_Out_of_Hours"`
		Grade    string  `json:"GRD"`
		DUName   string  `json:"DUN"`
		Line_man string  `json:"ЛИН_РУК"`
		Location string  `json:"ЛКЦ"`
		Number   int     `json:"Number"`
	*/
}

type FieldsTable struct {
	Code string `json:"Code"`
	Name string `json:"Name"`
}

type RespJSONResource struct {
	Fields    []FieldsTable   `json:"fields"`
	Array     []TableResource `json:"array"`
	Total     int             `json:"total"`
	ErrCode   int             `json:"errcode"`
	ErrorResp string          `json:"error"`
}

type TableGanttAct struct {
	Level        string  `json:"Level"`
	Code         string  `json:"Code"`
	Name         string  `json:"Name"`
	Guid         string  `json:"Guid"`
	Start        int64   `json:"Start"`
	Final        int64   `json:"Fin"`
	FactStart    int64   `json:"FactStart"`
	FactFinal    int64   `json:"FactFin"`
	WorkLoadPlan float64 `json:"WorkLoadSum"`
	WorkLoadFact float64 `json:"WorkLoadFact"`
}

type RespJSONGanttAct struct {
	Fields    []FieldsTable   `json:"fields"`
	Array     []TableGanttAct `json:"array"`
	Total     int             `json:"total"`
	ErrCode   int             `json:"errcode"`
	ErrorResp string          `json:"error"`
}

func OpenFile(url, PathSpiderDB string) (string, error) {
	var respOF RespOpenFile
	body := []byte(`{"command":"openFile", "fileName":"` + strings.ReplaceAll(PathSpiderDB, `\`, `\\`) + `", "sessId":""}`) //открываем файл с базой Спайдера
	reqOpenFile, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil || reqOpenFile.StatusCode != 200 {
		return "", fmt.Errorf("ошибка метода http Post на команду openFile: %v", err)
	}
	defer reqOpenFile.Body.Close()
	jsondec := json.NewDecoder(reqOpenFile.Body)
	err = jsondec.Decode(&respOF)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения ответа json команды openFile: %v", err)
	}
	docHandle := fmt.Sprint(respOF.DocHandle)
	if len(respOF.ErrorResp) > 0 {
		return docHandle, fmt.Errorf("ошибка на openFile таблицы: %v", respOF.ErrorResp)
	}
	if respOF.Errcode > 0 {
		return docHandle, fmt.Errorf("ошибка без описания на openFile таблицы")
	}
	return docHandle, nil
}

func GetTableHandle(url, docHandle, tableName string) (string, error) {
	var respGTH RespGetTableHandle
	body := []byte(`{"command":"getTableHandle", "docHandle":` + docHandle + `,"table":"` + tableName + `", "sessId":""}`) //получем код таблицы ресурсов
	reqGetTableH, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil || reqGetTableH.StatusCode != 200 {
		return "", fmt.Errorf("ошибка метода http Post на команду getTableHandle: %v", err)
	}
	defer reqGetTableH.Body.Close()
	jsondec := json.NewDecoder(reqGetTableH.Body)
	err = jsondec.Decode(&respGTH)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения ответа json команды getTableHandle: %v", err)
	}
	tableHandle := fmt.Sprint(respGTH.TableHandle)
	if len(respGTH.ErrorResp) > 0 {
		return tableHandle, fmt.Errorf("ошибка на getTableHandle таблицы %v: %v", tableName, respGTH.ErrorResp)
	}
	if respGTH.Errcode > 0 {
		return tableHandle, fmt.Errorf("ошибка без описания на getTableHandle таблицы %v: %v", tableName, respGTH.ErrorResp)
	}
	return tableHandle, nil
}

//map[Guid] = Code (email сотрудника)
func GetTableResource(url, tableHandle string) (map[string]string, error) {
	var respGT RespJSONResource
	body := []byte(`{"command":"getTable","tableHandle":` + tableHandle + `,"sessId":""}`) //получем значения в таблице
	reqGetTable, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil || reqGetTable.StatusCode != 200 {
		return nil, fmt.Errorf("ошибка метода http Post команды getTable таблицы Resource: %v", err)
	}
	defer reqGetTable.Body.Close()
	jsondec := json.NewDecoder(reqGetTable.Body)
	err = jsondec.Decode(&respGT)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа json команды getTable таблицы Resource: %v", err)
	}
	mapGuid := make(map[string]string, len(respGT.Array))
	for _, el := range respGT.Array {
		mapGuid[el.Guid] = el.Code
	}
	if len(respGT.ErrorResp) > 0 {
		return mapGuid, fmt.Errorf("ошибка на getTable таблицы Resource: %v", respGT.ErrorResp)
	}
	if respGT.ErrCode > 0 {
		return mapGuid, fmt.Errorf("ошибка без описания на getTable таблицы Resource: %v", respGT.ErrorResp)
	}

	return mapGuid, nil
}

func GetTableGanttAct(url, tableHandle string) ([]TableGanttAct, error) {
	var respGT RespJSONGanttAct
	body := []byte(`{"command":"getTable","tableHandle":` + tableHandle + `,"sessId":""}`) //получем значения в таблице
	reqGetTable, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil || reqGetTable.StatusCode != 200 {
		return nil, fmt.Errorf("ошибка метода http Post команды getTable таблицы GanttAct: %v", err)
	}
	defer reqGetTable.Body.Close()
	jsondec := json.NewDecoder(reqGetTable.Body)
	err = jsondec.Decode(&respGT)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа json команды getTable таблицы GanttAct: %v", err)
	}
	if len(respGT.ErrorResp) > 0 {
		return respGT.Array, fmt.Errorf("ошибка на getTable таблицы GanttAct: %v", respGT.ErrorResp)
	}
	if respGT.ErrCode > 0 {
		return respGT.Array, fmt.Errorf("ошбика без описания на getTable таблицы GanttAct")
	}
	return respGT.Array, nil
}

func ClearFilterTable(url, tableHandle, nameTable string) error {
	var respClFil RespErr
	body := []byte(`{"command":"clearFilter","tableHandle":` + tableHandle + `,"sessId":""}`) //очищаем фильтра в таблице
	reqClearFIl, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil || reqClearFIl.StatusCode != 200 {
		return fmt.Errorf("ошибка метода http Post команды clearFilter: %v", err)
	}
	defer reqClearFIl.Body.Close()
	jsondec := json.NewDecoder(reqClearFIl.Body)
	err = jsondec.Decode(&respClFil)
	if err != nil {
		return fmt.Errorf("ошибка чтения ответа json команды clearFilter: %v", err)
	}
	if len(respClFil.ErrorResp) > 0 {
		return fmt.Errorf("ошибка на clearFilter таблицы %v: %v", nameTable, respClFil.ErrorResp)
	}
	if respClFil.Errcode > 0 {
		return fmt.Errorf("ошибка без описания на clearFilter таблицы %v", nameTable)
	}
	return nil
}

func LoadLayout(url, docHandle, tableCode, layoutCode string) error {
	var respLL RespErr
	body := []byte(`{"command":"loadLayout","docHandle":` + docHandle + `,"table":"` + tableCode + `","layoutCode":"` + layoutCode + `","sessId":""}`) //вызываем созданную конфигурацию
	reqLL, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil || reqLL.StatusCode != 200 {
		return fmt.Errorf("ошибка метода http Post команды loadLayout: %v", err)
	}
	defer reqLL.Body.Close()
	jsondec := json.NewDecoder(reqLL.Body)
	err = jsondec.Decode(&respLL)
	if err != nil {
		return fmt.Errorf("ошибка чтения ответа json команды loadLayout: %v", err)
	}
	if len(respLL.ErrorResp) > 0 {
		return fmt.Errorf("ошибка на LoadLayout %v таблицы %v: %v", layoutCode, tableCode, respLL.ErrorResp)
	}
	if respLL.Errcode > 0 {
		return fmt.Errorf("ошибка без описания на LoadLayout %v таблицы %v", layoutCode, tableCode)
	}
	return nil
}

func CloseFile(url, docHandle string) error {
	var respErrorClose RespErr
	body := []byte(`{"command":"closeFile", "docHandle":` + docHandle + `, "sessId":""}`) //закрываем файл Спайдер
	reqCloseTable, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil || reqCloseTable.StatusCode != 200 {
		return fmt.Errorf("ошибка метода http Post на команду closeFile: %v", err)
	}
	defer reqCloseTable.Body.Close()
	jsondec := json.NewDecoder(reqCloseTable.Body)
	err = jsondec.Decode(&respErrorClose)
	if err != nil {
		return fmt.Errorf("ошибка чтения ответа json на команду closeFile в API Spider: %v", err)
	}
	if len(respErrorClose.ErrorResp) > 0 {
		return fmt.Errorf("ошибка на команду closeFile: %v", respErrorClose.ErrorResp)
	}
	if respErrorClose.Errcode > 0 {
		return fmt.Errorf("ошибка без описания на команду closeFile")
	}
	return nil
}

func SaveFile(url, docHandle string) error {
	var respErrorSave RespErr
	body := []byte(`{"command":"saveFile", "docHandle":` + docHandle + `, "sessId":""}`) //сохраняем файл Спайдер
	reqSaveTable, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil || reqSaveTable.StatusCode != 200 {
		return fmt.Errorf("ошибка метода http Post в API Spider:%v", err)
	}
	defer reqSaveTable.Body.Close()
	jsondec := json.NewDecoder(reqSaveTable.Body)
	err = jsondec.Decode(&respErrorSave)
	if err != nil {
		return fmt.Errorf("ошибка чтения ответа json на saveFile в API Spider: %v", err)
	}
	if len(respErrorSave.ErrorResp) > 0 {
		return fmt.Errorf("ошибка на saveFile: %v", respErrorSave.ErrorResp)
	}
	if respErrorSave.Errcode > 0 {
		return fmt.Errorf("ошибка без описания на saveFile")
	}
	return nil
}

func SetTable(url, tableHandle, fields, data string) error {
	var respErrorRes RespErr
	body := []byte(fmt.Sprint(`{"command":"setTable","tableHandle":`, tableHandle, `,"sessId":"","fields":`, fields, `,"data":`, data, `}`))
	reqSetTablePerRes, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil || reqSetTablePerRes.StatusCode != 200 {
		return fmt.Errorf("ошибка метода http Post на SetTable в API Spider: %v", err)
	}
	defer reqSetTablePerRes.Body.Close()
	jsondec := json.NewDecoder(reqSetTablePerRes.Body)
	err = jsondec.Decode(&respErrorRes)
	if err != nil {
		return fmt.Errorf("ошибка чтения ответа json на setTable в API Spider: %v", err)
	}
	if len(respErrorRes.ErrorResp) > 0 {
		return fmt.Errorf("ошибка API Spider на SetTable: %v", respErrorRes.ErrorResp)
	}
	if respErrorRes.Errcode > 0 {
		return fmt.Errorf("ошибка без описания API Spider на SetTable")
	}
	return nil
}
