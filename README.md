# Jira Tempo Timesheet Automation

Automate your daily worklog submissions to **Tempo Timesheets (Jira)** using a simple Go script and CSV files.  
Easily manage and batch-upload your work activities for compliance and productivity reporting.

---

## 📖 Overview

This repository provides a Go script that reads your work activities from a CSV file and sends them to your organization's **Tempo Timesheets** API (Jira).  
Each worklog is sent as a separate POST request, and all results (including errors and worklog IDs) are written to an output CSV for easy review.

---

## 📦 File Structure

.
├── main.go                # The main Go script (edit credential variables at the top)
├── worklog.csv            # Example input CSV file (copy and edit for your own activities)
├── worklog_result.csv     # Output CSV file (generated automatically; do not edit manually)
└── README.md              # This documentation

---

## 🛠️ Prerequisites

- [Go (Golang)](https://golang.org/dl/) installed (1.18+ recommended)
- Access to your organization's **Jira Tempo Timesheets API**
- API token (or Jira password) and your email/user for authentication
- Worklog data in a CSV file with the following columns:
  worker,started,timeSpentSeconds,originTaskId,WorkType,comment

---

## 📝 Preparing Your Worklog CSV

Example worklog_1707.csv:

worker,started,timeSpentSeconds,originTaskId,WorkType,comment
00000001,2025-07-17T07:30:00.000,5400,4445761,Doa+BriefingPagi,Briefing dan doa pagi bersama tim.
00000001,2025-07-17T09:00:00.000,7200,4445721,Development,Mengembangkan fitur baru dan perbaikan bug.
...

- worker: Your Jira/Tempo user or employee ID
- started: Start time, ISO format (YYYY-MM-DDTHH:MM:SS.000)
- timeSpentSeconds: Duration of activity in seconds (e.g., 7200 = 2 hours)
- originTaskId: The Jira issue/task ID (numeric)
- WorkType: Type of work (Development, Dokumentasi, etc.)
- comment: Short description of the activity

---

## ⚙️ Configuration

Open main.go and configure these variables near the top:

csvFile := "worklog.csv"                // Input CSV file
jiraUrl := "https://your-jira-url/rest/tempo-timesheets/4/worklogs/" // Change to your Tempo API endpoint
email := "your-email"                        // Your Jira/Tempo account email/user
passEmail := "your-api-token-or-password"    // API token or password for authentication

Note:  
- For Atlassian Cloud, use API token ([generate here](https://id.atlassian.com/manage-profile/security/api-tokens)) instead of your password.

---

## 🚀 How To Run

1. Place your worklog CSV in the same directory as main.go.
2. Open terminal in that directory.
3. Run:
   go run main.go
4. The script will:
    - Read each row in the CSV file
    - POST each worklog to the Tempo API
    - Print progress log for each request in your terminal
    - Write all results (success or fail) into worklog_result.csv

---

## ✅ Output: worklog_result.csv

Each row contains:

rowNum, status, error, tempoWorklogId, timeSpentSeconds, worker, started, originTaskId, WorkType, comment, issueKey, issueSummary

- status: SUCCESS if uploaded, FAILED otherwise
- error: error message if failed, empty if success
- tempoWorklogId: Tempo's unique ID for the uploaded worklog (if successful)

---

## 🐞 Troubleshooting

- Error "Parse JSON gagal":  
  Check that the Tempo API response is a JSON array (e.g., [ { ... } ]).  
  Adjust parsing if API response format changes.
- POST failed:  
  Check your credentials, Jira URL, and CSV file content.
- CSV errors:  
  Ensure all required columns are filled for each row.

---

## 💡 Tips

- Batch Upload: You can upload as many worklogs as needed in one CSV file.
- Custom Fields: You can extend the Go script to add more fields or change mapping as needed.

---

## 📜 License

This repository is licensed under the MIT License.

---

Need help or want to contribute? Open an issue or pull request!

---

Happy automating 🚀
