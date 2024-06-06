package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "time"
	"net/url"

    _ "github.com/go-sql-driver/mysql"
    "github.com/gorilla/mux"
)


func connectDB() (*sql.DB, error) {
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    dbHost := os.Getenv("DB_HOST")

    dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":3306)/" + dbName
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }
    return db, nil
}

type Log struct {
	Table_ID  int		`json:"table_id"`
    USERID    int       `json:"user_id"`
    Timestamp time.Time `json:"timestamp"`
    LogLevel  string    `json:"log_level"`
    Message   string    `json:"message"`
    Event_id  string    `json:"source"`
}


//文字列から番号を求める関数
//DBの接続を切らさないために　引数からもらう
func getEvent_ID(db *sql.DB, eventid string) (int,error){
	var event_id int
	rows, err :=db.Query("Select Event_id from Events where Event_message = ?", eventid)
	if err != nil {
		return  -1 , err;
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&event_id)
	if err != nil {
		return  -1 , err
	}
	return event_id, nil
}

//
func addLog(w http.ResponseWriter, r* http.Request){
	var logEntry Log
	err := json.NewDecoder(r.Body).Decode(&logEntry)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := connectDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	event_id ,err := getEvent_ID(db, logEntry.Event_id);
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	

	stmt, err := db.Prepare("Insert Into logs(User_id, timestamp, severity, message, Event_ID, table_id) Values (?,?,?,?,?,?)")
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(logEntry.USERID,logEntry.Timestamp, logEntry.LogLevel, url.QueryEscape(logEntry.Message), event_id,logEntry.Table_ID)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	w.WriteHeader(http.StatusCreated)
}

func getLogs(w http.ResponseWriter, r *http.Request){
	db, err := connectDB()
	if  err != nil{
		http.Error(w, err.Error(),http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("select logs.User_id, logs.timestamp, Events.Event_message, logs.table_id, logs.severity,logs.message from logs join Events on logs.Event_id = Events.Event_id")
	if err != nil {
		http.Error(w, err.Error(),http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var logs []Log
	for rows.Next(){
		var logEntry Log
		var timestampStr string
		var messageencode string
		err := rows.Scan(&logEntry.USERID,&timestampStr, &logEntry.Event_id,&logEntry.Table_ID,&logEntry.LogLevel,&messageencode)
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logEntry.Message , _ = url.QueryUnescape(messageencode);
        logEntry.Timestamp, err = time.Parse("2006-01-02 15:04:05", timestampStr)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
		logs = append(logs, logEntry)
	} 
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}


func main() {
    r := mux.NewRouter()
    r.HandleFunc("/logs", addLog).Methods("POST")
    r.HandleFunc("/logs", getLogs).Methods("GET")

    fmt.Print(http.ListenAndServe(":8080", r))
}