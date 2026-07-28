package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Job struct {
	ID          string
	Data        string
	CallbackURL string
}

// 전역 큐 (worker pool 시뮬레이션)
var jobCh = make(chan Job, 100)

func main() {
	go startClientServer() // webhook 수신 클라이언트
	go startMainServer()   // 메인 async 서버
	go workerPool(3)       // 워커 3개 가동

	select {} // main 종료 방지
}

// ① 클라이언트 서버 (Port 9090) : webhook 수신 & 요청 발생
func startClientServer() {
	mux := http.NewServeMux()

	// Webhook 수신 핸들러
	mux.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		json.NewDecoder(r.Body).Decode(&payload)
		log.Printf("[Client] Webhook 수신! Payload: %+v", payload)
		w.WriteHeader(http.StatusOK)
	})

	// 테스트 요청 (메인 서버에 async 요청 보내기)
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := json.Marshal(map[string]string{
			"data":        "example-data",
			"callbackUrl": "http://localhost:9090/webhook",
		})

		resp, err := http.Post("http://localhost:8080/task", "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var result map[string]any
		json.NewDecoder(resp.Body).Decode(&result)
		log.Printf("[Client] 메인 서버 응답: %+v", result)
		json.NewEncoder(w).Encode(result)
	})

	log.Println("[Client] Webhook 서버 시작 (http://localhost:9090)")
	http.ListenAndServe(":9090", mux)
}

// ② 메인 서버 (Port 8080) : 비동기 요청 수신 + webhook 호출
func startMainServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/task", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Data        string `json:"data"`
			CallbackURL string `json:"callbackUrl"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		job := Job{
			ID:          fmt.Sprintf("job-%d", rand.Intn(100000)),
			Data:        req.Data,
			CallbackURL: req.CallbackURL,
		}

		select {
		case jobCh <- job:
			log.Printf("[Server] Job 등록됨: %+v", job)
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(map[string]string{
				"jobId":   job.ID,
				"message": "Job accepted. Processing asynchronously.",
			})
		default:
			http.Error(w, "Job queue full", http.StatusServiceUnavailable)
		}
	})

	log.Println("[Server] Async 메인 서버 시작 (http://localhost:8080)")
	http.ListenAndServe(":8080", mux)
}

// ③ Worker Pool (비동기 처리 + webhook 전송)
func workerPool(n int) {
	for i := 0; i < n; i++ {
		go func(id int) {
			log.Printf("[Worker-%d] 시작됨", id)
			for job := range jobCh {
				log.Printf("[Worker-%d] Job 처리중: %s", id, job.ID)
				result := processJob(job)
				sendWebhook(job.CallbackURL, result)
			}
		}(i)
	}
}

func processJob(job Job) map[string]any {
	// 실제로는 DB, 외부 API, AI inference 등의 긴 작업 가능
	time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second)
	return map[string]any{
		"jobId":  job.ID,
		"status": "completed",
		"result": fmt.Sprintf("Processed data: %s", job.Data),
	}
}

func sendWebhook(callbackURL string, payload map[string]any) {
	b, _ := json.Marshal(payload)
	resp, err := http.Post(callbackURL, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Printf("[Server] Webhook 전송 실패: %v", err)
		return
	}
	defer resp.Body.Close()
	log.Printf("[Server] Webhook 전송 완료 (%s): status=%d", callbackURL, resp.StatusCode)
}
