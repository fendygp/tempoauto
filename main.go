package main

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type WorkTypeAttr struct {
	Name            string `json:"name"`
	WorkAttributeId int    `json:"workAttributeId"`
	Value           string `json:"value"`
}
type Attributes struct {
	WorkType WorkTypeAttr `json:"_WorkType_"`
}
type Payload struct {
	Attributes            Attributes `json:"attributes"`
	BillableSeconds       string     `json:"billableSeconds"`
	OriginId              int        `json:"originId"`
	Worker                string     `json:"worker"`
	Comment               *string    `json:"comment"`
	Started               string     `json:"started"`
	TimeSpentSeconds      int        `json:"timeSpentSeconds"`
	OriginTaskId          string     `json:"originTaskId"`
	RemainingEstimate     int        `json:"remainingEstimate"`
	EndDate               *string    `json:"endDate"`
	IncludeNonWorkingDays bool       `json:"includeNonWorkingDays"`
}

// Struct sesuai sebagian response Tempo
type Issue struct {
	Key     string `json:"key"`
	ID      int    `json:"id"`
	Summary string `json:"summary"`
}
type AttrResp struct {
	WorkAttributeId int    `json:"workAttributeId"`
	Value           string `json:"value"`
}
type TempoResponse struct {
	TempoWorklogId   int                 `json:"tempoWorklogId"`
	TimeSpentSeconds int                 `json:"timeSpentSeconds"`
	Comment          string              `json:"comment"`
	Worker           string              `json:"worker"`
	Started          string              `json:"started"`
	OriginTaskId     int                 `json:"originTaskId"`
	Attributes       map[string]AttrResp `json:"attributes"`
	Issue            Issue               `json:"issue"`
	Status           string              // custom, not in response
	RowNum           int                 // custom, for CSV
	ErrorMsg         string              // custom, for CSV
}

func main() {
	csvFile := "workload.csv"
	jiraUrl := "https://xx.xx.xx.xx/rest/tempo-timesheets/4/worklogs/"
	email := "xxxx" // Ganti dengan email Atlassian-mu
	passEmail := "xxxx"

	auth := email + ":" + passEmail
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	// Open CSV input
	file, err := os.Open(csvFile)
	if err != nil {
		fmt.Println("Gagal buka file CSV:", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	header, err := reader.Read()
	if err != nil {
		fmt.Println("Gagal baca header CSV:", err)
		return
	}

	// Prepare output CSV
	outFile, err := os.Create("worklog_result.csv")
	if err != nil {
		fmt.Println("Gagal buat file CSV output:", err)
		return
	}
	defer outFile.Close()
	outWriter := csv.NewWriter(outFile)
	defer outWriter.Flush()

	// CSV header
	outWriter.Write([]string{
		"rowNum", "status", "error", "tempoWorklogId", "timeSpentSeconds", "worker", "started", "originTaskId",
		"WorkType", "comment", "issueKey", "issueSummary",
	})

	rowNum := 1
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		rowNum++
		if err != nil {
			fmt.Printf("Baris %d dilewati (error): %v\n", rowNum, err)
			continue
		}
		if len(row) < 5 {
			fmt.Printf("Baris %d dilewati: data tidak lengkap\n", rowNum)
			continue
		}
		worker := getCol(header, row, "worker")
		started := getCol(header, row, "started")
		timeSpentSeconds, _ := strconv.Atoi(getCol(header, row, "timeSpentSeconds"))
		originTaskId := getCol(header, row, "originTaskId")
		workType := getCol(header, row, "WorkType")
		commentStr := getCol(header, row, "comment")
		var comment *string
		if commentStr != "" {
			comment = &commentStr
		}

		// LOG setiap request
		fmt.Printf("Baris %d: POST [%s] %s, %s, %ss, %s\n",
			rowNum, started, workType, originTaskId, strconv.Itoa(timeSpentSeconds), worker)

		payload := Payload{
			Attributes: Attributes{
				WorkType: WorkTypeAttr{
					Name:            "Work Type",
					WorkAttributeId: 7,
					Value:           workType,
				},
			},
			BillableSeconds:       "",
			OriginId:              -1,
			Worker:                worker,
			Comment:               comment,
			Started:               started,
			TimeSpentSeconds:      timeSpentSeconds,
			OriginTaskId:          originTaskId,
			RemainingEstimate:     0,
			EndDate:               nil,
			IncludeNonWorkingDays: false,
		}
		payloadBytes, _ := json.Marshal(payload)

		req, err := http.NewRequest("POST", jiraUrl, bytes.NewBuffer(payloadBytes))
		if err != nil {
			fmt.Printf("Baris %d: gagal create request: %v\n", rowNum, err)
			continue
		}
		req.Header.Set("Authorization", authHeader)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Baris %d: gagal POST: %v\n", rowNum, err)
			continue
		}
		body, _ := ioutil.ReadAll(resp.Body)

		result := TempoResponse{
			Status:   "FAILED",
			RowNum:   rowNum,
			ErrorMsg: "",
		}

		if resp.StatusCode == 200 || resp.StatusCode == 201 {
			var resultArr []TempoResponse
			err = json.Unmarshal(body, &resultArr)
			if err != nil {
				result.ErrorMsg = fmt.Sprintf("Parse JSON gagal: %v", err)
			} else if len(resultArr) > 0 {
				result = resultArr[0]
				result.Status = "SUCCESS"
			} else {
				result.ErrorMsg = "Response array kosong"
			}
		} else {
			result.ErrorMsg = string(body)
		}
		resp.Body.Close()

		// Write to output CSV
		outWriter.Write([]string{
			strconv.Itoa(result.RowNum),
			result.Status,
			result.ErrorMsg,
			strconv.Itoa(result.TempoWorklogId),
			strconv.Itoa(result.TimeSpentSeconds),
			result.Worker,
			result.Started,
			strconv.Itoa(result.OriginTaskId),
			workType,
			commentStr,
			result.Issue.Key,
			result.Issue.Summary,
		})
	}
}

// Helper: get value kolom by name
func getCol(header, row []string, name string) string {
	for i, h := range header {
		if h == name && i < len(row) {
			return row[i]
		}
	}
	return ""
}
