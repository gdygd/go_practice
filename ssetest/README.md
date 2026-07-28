# SSE Library for Go

Go 언어로 작성된 Server-Sent Events(SSE) 라이브러리입니다. Gin 프레임워크와 통합되어 있으며, 브로드캐스트, 단일 유저 전송, 세션 관리 기능을 제공합니다.

## 설치

```bash
go get github.com/gin-gonic/gin
go get github.com/google/uuid
```

## 프로젝트 구조

```
ssetest/
├── go.mod
├── sse/
│   ├── sse.go        # 핵심 SSE 로직
│   └── session.go    # 세션 관리
└── example/
    └── main.go       # Gin 예제
```

---

## 패키지: sse

SSE 연결 관리 및 메시지 전송을 위한 핵심 패키지입니다.

---

## Types (타입)

### EventData

SSE 이벤트 데이터를 표현하는 구조체입니다.

```go
type EventData struct {
    Msgtype string `json:"event"` // 이벤트 타입 (예: "message", "notification")
    Data    string `json:"data"`  // 전송할 데이터
    Id      string `json:"id"`    // 이벤트 ID (선택사항)
}
```

**SSE 와이어 포맷 출력 예시:**
```
id: 123
event: notification
data: Hello World

```

---

### Client

개별 SSE 클라이언트 연결을 나타냅니다.

```go
type Client struct {
    ID      string           // 클라이언트 고유 식별자
    Channel chan EventData   // 이벤트 수신 채널 (버퍼 크기: 256)
}
```

---

### Hub

모든 SSE 연결을 관리하고 메시지 브로드캐스트를 처리하는 중앙 허브입니다.

```go
type Hub struct {
    clients    map[string]*Client  // 연결된 클라이언트 맵
    register   chan *Client        // 클라이언트 등록 채널
    unregister chan *Client        // 클라이언트 해제 채널
    broadcast  chan EventData      // 브로드캐스트 채널
    mu         sync.RWMutex        // 동시성 제어용 뮤텍스
}
```

---

### Session

사용자 세션 정보를 저장하는 구조체입니다.

```go
type Session struct {
    ID        string         // 세션 고유 ID
    UserID    string         // 사용자 ID
    CreatedAt time.Time      // 생성 시간
    LastPing  time.Time      // 마지막 활동 시간
    Metadata  map[string]any // 추가 메타데이터
}
```

---

### SessionManager

세션 생명주기를 관리하는 매니저입니다.

```go
type SessionManager struct {
    sessions map[string]*Session  // 세션 맵
    mu       sync.RWMutex         // 동시성 제어용 뮤텍스
    hub      *Hub                 // 연결된 Hub 참조
}
```

---

## Functions (함수)

### Hub 관련

#### NewHub

새로운 Hub 인스턴스를 생성합니다.

```go
func NewHub() *Hub
```

**예시:**
```go
hub := sse.NewHub()
go hub.Run() // 반드시 고루틴으로 실행
```

---

#### (*Hub) Run

Hub의 메인 이벤트 루프를 시작합니다. 클라이언트 등록/해제 및 브로드캐스트를 처리합니다.

```go
func (h *Hub) Run()
```

**주의:** 반드시 고루틴으로 실행해야 합니다.

```go
go hub.Run()
```

---

#### (*Hub) Register

새 클라이언트를 Hub에 등록합니다.

```go
func (h *Hub) Register(client *Client)
```

---

#### (*Hub) Unregister

클라이언트를 Hub에서 해제합니다. 클라이언트의 채널도 함께 닫힙니다.

```go
func (h *Hub) Unregister(client *Client)
```

---

#### (*Hub) Broadcast

모든 연결된 클라이언트에게 이벤트를 전송합니다.

```go
func (h *Hub) Broadcast(event EventData)
```

**예시:**
```go
event := sse.NewEventDataString("announcement", "서버 점검 예정", "1")
hub.Broadcast(event)
```

---

#### (*Hub) SendToClient

특정 클라이언트에게 이벤트를 전송합니다.

```go
func (h *Hub) SendToClient(clientID string, event EventData) bool
```

**반환값:**
- `true`: 전송 성공
- `false`: 클라이언트가 없거나 채널이 가득 참

**예시:**
```go
event := sse.NewEventDataString("message", "개인 메시지", "2")
success := hub.SendToClient("client-123", event)
```

---

#### (*Hub) GetClientCount

현재 연결된 클라이언트 수를 반환합니다.

```go
func (h *Hub) GetClientCount() int
```

---

#### (*Hub) HasClient

특정 클라이언트가 연결되어 있는지 확인합니다.

```go
func (h *Hub) HasClient(clientID string) bool
```

---

### Gin Handler 함수

#### ServeSSE

Gin 프레임워크용 SSE 핸들러입니다. Gin의 `c.SSEvent()`를 사용합니다.

```go
func ServeSSE(hub *Hub, clientID string) gin.HandlerFunc
```

**예시:**
```go
r.GET("/sse/:id", func(c *gin.Context) {
    clientID := c.Param("id")
    sse.ServeSSE(hub, clientID)(c)
})
```

---

#### ServeSSERaw

SSE 와이어 포맷(id, event, data)을 직접 출력하는 핸들러입니다.

```go
func ServeSSERaw(hub *Hub, clientID string) gin.HandlerFunc
```

**차이점:**
- `ServeSSE`: Gin의 내장 SSE 메서드 사용
- `ServeSSERaw`: EventData의 모든 필드(id, event, data) 직접 출력

---

### EventData 생성 함수

#### NewEventData

데이터를 JSON으로 마샬링하여 EventData를 생성합니다.

```go
func NewEventData(msgType string, data any, id string) (EventData, error)
```

**예시:**
```go
payload := map[string]string{"message": "Hello", "user": "John"}
event, err := sse.NewEventData("chat", payload, "msg-1")
// 결과: Data = `{"message":"Hello","user":"John"}`
```

---

#### NewEventDataString

문자열 데이터로 EventData를 생성합니다.

```go
func NewEventDataString(msgType string, data string, id string) EventData
```

**예시:**
```go
event := sse.NewEventDataString("notification", "새 메시지가 도착했습니다", "3")
```

---

### SessionManager 관련

#### NewSessionManager

새로운 SessionManager를 생성합니다.

```go
func NewSessionManager(hub *Hub) *SessionManager
```

**예시:**
```go
hub := sse.NewHub()
sessionManager := sse.NewSessionManager(hub)
```

---

#### (*SessionManager) CreateSession

새 세션을 생성합니다.

```go
func (sm *SessionManager) CreateSession(sessionID, userID string) *Session
```

**예시:**
```go
session := sessionManager.CreateSession("sess-123", "user-456")
```

---

#### (*SessionManager) GetSession

세션 ID로 세션을 조회합니다.

```go
func (sm *SessionManager) GetSession(sessionID string) (*Session, bool)
```

---

#### (*SessionManager) DeleteSession

세션을 삭제합니다.

```go
func (sm *SessionManager) DeleteSession(sessionID string)
```

---

#### (*SessionManager) UpdateLastPing

세션의 마지막 활동 시간을 업데이트합니다.

```go
func (sm *SessionManager) UpdateLastPing(sessionID string)
```

---

#### (*SessionManager) SetMetadata / GetMetadata

세션에 메타데이터를 저장하고 조회합니다.

```go
func (sm *SessionManager) SetMetadata(sessionID, key string, value any)
func (sm *SessionManager) GetMetadata(sessionID, key string) (any, bool)
```

**예시:**
```go
sessionManager.SetMetadata("sess-123", "role", "admin")
role, ok := sessionManager.GetMetadata("sess-123", "role")
```

---

#### (*SessionManager) GetSessionsByUserID

특정 사용자의 모든 세션을 조회합니다. (다중 디바이스 지원)

```go
func (sm *SessionManager) GetSessionsByUserID(userID string) []*Session
```

---

#### (*SessionManager) SendToUser

특정 사용자의 모든 세션에 이벤트를 전송합니다.

```go
func (sm *SessionManager) SendToUser(userID string, event EventData) int
```

**반환값:** 성공적으로 전송된 세션 수

**예시:**
```go
event := sse.NewEventDataString("alert", "새 알림", "")
sent := sessionManager.SendToUser("user-456", event)
fmt.Printf("%d개 세션에 전송됨\n", sent)
```

---

#### (*SessionManager) SendToSession

특정 세션에 이벤트를 전송합니다.

```go
func (sm *SessionManager) SendToSession(sessionID string, event EventData) bool
```

---

#### (*SessionManager) Broadcast

모든 세션에 이벤트를 브로드캐스트합니다.

```go
func (sm *SessionManager) Broadcast(event EventData)
```

---

#### (*SessionManager) GetAllSessions

모든 활성 세션 목록을 반환합니다.

```go
func (sm *SessionManager) GetAllSessions() []*Session
```

---

#### (*SessionManager) GetSessionCount

활성 세션 수를 반환합니다.

```go
func (sm *SessionManager) GetSessionCount() int
```

---

#### (*SessionManager) CleanupStaleSessions

지정된 시간 동안 활동이 없는 세션을 정리합니다.

```go
func (sm *SessionManager) CleanupStaleSessions(timeout time.Duration) int
```

**반환값:** 삭제된 세션 수

---

#### (*SessionManager) StartCleanupRoutine

주기적으로 비활성 세션을 정리하는 고루틴을 시작합니다.

```go
func (sm *SessionManager) StartCleanupRoutine(interval, timeout time.Duration)
```

**예시:**
```go
// 1분마다 체크, 5분 이상 비활성 세션 삭제
sessionManager.StartCleanupRoutine(1*time.Minute, 5*time.Minute)
```

---

## 전체 사용 예시

```go
package main

import (
    "ssetest/sse"
    "github.com/gin-gonic/gin"
)

func main() {
    // 1. Hub 및 SessionManager 초기화
    hub := sse.NewHub()
    sessionManager := sse.NewSessionManager(hub)

    // 2. Hub 실행 (필수)
    go hub.Run()

    // 3. 세션 정리 루틴 시작 (선택)
    sessionManager.StartCleanupRoutine(1*time.Minute, 5*time.Minute)

    // 4. Gin 라우터 설정
    r := gin.Default()

    // SSE 연결 엔드포인트
    r.GET("/sse/:sessionID", func(c *gin.Context) {
        sessionID := c.Param("sessionID")
        userID := c.Query("userID")

        sessionManager.CreateSession(sessionID, userID)
        defer sessionManager.DeleteSession(sessionID)

        sse.ServeSSERaw(hub, sessionID)(c)
    })

    // 메시지 전송 엔드포인트
    r.POST("/broadcast", func(c *gin.Context) {
        var event sse.EventData
        c.ShouldBindJSON(&event)
        sessionManager.Broadcast(event)
        c.JSON(200, gin.H{"status": "ok"})
    })

    r.Run(":8080")
}
```

---

## API 엔드포인트 (예제 서버)

| Method | Path | 설명 |
|--------|------|------|
| GET | `/` | 테스트 HTML 페이지 |
| GET | `/sse/:sessionID` | SSE 연결 |
| GET | `/sse/user/:userID` | 유저별 SSE 연결 (세션 자동 생성) |
| POST | `/send/:sessionID` | 특정 세션에 메시지 전송 |
| POST | `/send/user/:userID` | 특정 유저에게 메시지 전송 |
| POST | `/broadcast` | 전체 브로드캐스트 |
| GET | `/sessions` | 전체 세션 목록 |
| GET | `/sessions/:sessionID` | 세션 상세 정보 |
| DELETE | `/sessions/:sessionID` | 세션 삭제 |
| PUT | `/sessions/:sessionID/ping` | 세션 활동 시간 갱신 |

---

## 실행 방법

```bash
cd example
go run main.go
```

브라우저에서 `http://localhost:8080` 접속

---

## 메시지 전송 예시 (curl)

```bash
# 브로드캐스트
curl -X POST http://localhost:8080/broadcast \
  -H "Content-Type: application/json" \
  -d '{"event":"notification","data":"Hello Everyone!","id":"1"}'

# 특정 세션에 전송
curl -X POST http://localhost:8080/send/session-123 \
  -H "Content-Type: application/json" \
  -d '{"event":"message","data":"Private message","id":"2"}'

# 특정 유저에게 전송
curl -X POST http://localhost:8080/send/user/user-456 \
  -H "Content-Type: application/json" \
  -d '{"event":"alert","data":"You have new notification","id":"3"}'
```

---

## JavaScript 클라이언트 예시

```javascript
const eventSource = new EventSource('/sse/my-session-id?userID=user-123');

eventSource.onopen = () => {
    console.log('Connected');
};

eventSource.addEventListener('notification', (e) => {
    console.log('Notification:', e.data);
});

eventSource.addEventListener('message', (e) => {
    console.log('Message:', e.data);
});

eventSource.onerror = (e) => {
    console.error('Error:', e);
};
```
