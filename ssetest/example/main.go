package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ssetest/sse"
)

var (
	hub            *sse.Hub
	sessionManager *sse.SessionManager
)

func main() {
	// Initialize Hub and SessionManager
	hub = sse.NewHub()
	sessionManager = sse.NewSessionManager(hub)

	// Start the hub
	go hub.Run()

	// Start cleanup routine (check every 1 minute, remove sessions inactive for 5 minutes)
	sessionManager.StartCleanupRoutine(1*time.Minute, 5*time.Minute)

	r := gin.Default()

	// SSE endpoints
	r.GET("/sse/:sessionID", handleSSE)
	r.GET("/sse/user/:userID", handleUserSSE)

	// Message sending endpoints
	r.POST("/send/:sessionID", handleSendToSession)
	r.POST("/send/user/:userID", handleSendToUser)
	r.POST("/broadcast", handleBroadcast)

	// Session management endpoints
	r.GET("/sessions", handleGetSessions)
	r.GET("/sessions/:sessionID", handleGetSession)
	r.DELETE("/sessions/:sessionID", handleDeleteSession)
	r.PUT("/sessions/:sessionID/ping", handlePing)

	// Test page
	r.GET("/", handleTestPage)

	log.Println("Server starting on :8080")
	if err := r.Run(":8083"); err != nil {
		log.Fatal(err)
	}
}

// handleSSE establishes SSE connection with given session ID
func handleSSE(c *gin.Context) {
	sessionID := c.Param("sessionID")
	userID := c.Query("userID")
	if userID == "" {
		userID = "anonymous"
	}

	// Create session
	sessionManager.CreateSession(sessionID, userID)
	defer sessionManager.DeleteSession(sessionID)

	// Serve SSE
	sse.ServeSSERaw(hub, sessionID)(c)
}

// handleUserSSE establishes SSE connection with auto-generated session ID for a user
func handleUserSSE(c *gin.Context) {
	userID := c.Param("userID")
	sessionID := uuid.New().String()

	sessionManager.CreateSession(sessionID, userID)
	defer sessionManager.DeleteSession(sessionID)

	c.Header("X-Session-ID", sessionID)
	sse.ServeSSERaw(hub, sessionID)(c)
}

// handleSendToSession sends a message to a specific session
func handleSendToSession(c *gin.Context) {
	sessionID := c.Param("sessionID")

	var event sse.EventData
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if sessionManager.SendToSession(sessionID, event) {
		c.JSON(http.StatusOK, gin.H{"status": "sent"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
	}
}

// handleSendToUser sends a message to all sessions of a user
func handleSendToUser(c *gin.Context) {
	userID := c.Param("userID")

	var event sse.EventData
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sent := sessionManager.SendToUser(userID, event)
	c.JSON(http.StatusOK, gin.H{
		"status":       "sent",
		"sent_to":      sent,
		"total_sessions": len(sessionManager.GetSessionsByUserID(userID)),
	})
}

// handleBroadcast sends a message to all connected clients
func handleBroadcast(c *gin.Context) {
	var event sse.EventData
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sessionManager.Broadcast(event)
	c.JSON(http.StatusOK, gin.H{
		"status":        "broadcast",
		"total_clients": hub.GetClientCount(),
	})
}

// handleGetSessions returns all active sessions
func handleGetSessions(c *gin.Context) {
	sessions := sessionManager.GetAllSessions()
	result := make([]gin.H, 0, len(sessions))
	for _, s := range sessions {
		result = append(result, gin.H{
			"id":         s.ID,
			"user_id":    s.UserID,
			"created_at": s.CreatedAt,
			"last_ping":  s.LastPing,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"sessions": result,
		"count":    len(sessions),
	})
}

// handleGetSession returns a specific session
func handleGetSession(c *gin.Context) {
	sessionID := c.Param("sessionID")
	session, ok := sessionManager.GetSession(sessionID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         session.ID,
		"user_id":    session.UserID,
		"created_at": session.CreatedAt,
		"last_ping":  session.LastPing,
		"metadata":   session.Metadata,
	})
}

// handleDeleteSession removes a session
func handleDeleteSession(c *gin.Context) {
	sessionID := c.Param("sessionID")
	sessionManager.DeleteSession(sessionID)
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// handlePing updates session's last ping time
func handlePing(c *gin.Context) {
	sessionID := c.Param("sessionID")
	sessionManager.UpdateLastPing(sessionID)
	c.JSON(http.StatusOK, gin.H{"status": "pinged"})
}

// handleTestPage serves a simple HTML test page
func handleTestPage(c *gin.Context) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>SSE Test</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .container { max-width: 800px; margin: 0 auto; }
        .section { margin: 20px 0; padding: 15px; border: 1px solid #ddd; border-radius: 5px; }
        h2 { margin-top: 0; }
        input, button { padding: 8px 12px; margin: 5px 0; }
        input { width: 200px; }
        button { cursor: pointer; background: #007bff; color: white; border: none; border-radius: 3px; }
        button:hover { background: #0056b3; }
        #events { height: 300px; overflow-y: auto; border: 1px solid #ccc; padding: 10px; background: #f9f9f9; }
        .event { padding: 5px; margin: 5px 0; border-left: 3px solid #007bff; padding-left: 10px; }
        .status { padding: 5px 10px; border-radius: 3px; display: inline-block; }
        .connected { background: #d4edda; color: #155724; }
        .disconnected { background: #f8d7da; color: #721c24; }
    </style>
</head>
<body>
    <div class="container">
        <h1>SSE Test Page</h1>

        <div class="section">
            <h2>Connection</h2>
            <input type="text" id="sessionId" placeholder="Session ID">
            <input type="text" id="userId" placeholder="User ID (optional)">
            <button onclick="connect()">Connect</button>
            <button onclick="disconnect()">Disconnect</button>
            <p>Status: <span id="status" class="status disconnected">Disconnected</span></p>
        </div>

        <div class="section">
            <h2>Send Message</h2>
            <input type="text" id="targetSession" placeholder="Target Session ID">
            <input type="text" id="eventType" placeholder="Event Type" value="message">
            <input type="text" id="eventData" placeholder="Data">
            <input type="text" id="eventId" placeholder="Event ID (optional)">
            <br>
            <button onclick="sendToSession()">Send to Session</button>
            <button onclick="broadcast()">Broadcast</button>
        </div>

        <div class="section">
            <h2>Events</h2>
            <button onclick="clearEvents()">Clear</button>
            <div id="events"></div>
        </div>
    </div>

    <script>
        let eventSource = null;
        let currentSessionId = null;

        function connect() {
            const sessionId = document.getElementById('sessionId').value || 'session-' + Date.now();
            const userId = document.getElementById('userId').value;

            if (eventSource) {
                eventSource.close();
            }

            const url = '/sse/' + sessionId + (userId ? '?userID=' + userId : '');
            eventSource = new EventSource(url);
            currentSessionId = sessionId;

            eventSource.onopen = function() {
                document.getElementById('status').textContent = 'Connected (' + sessionId + ')';
                document.getElementById('status').className = 'status connected';
                document.getElementById('sessionId').value = sessionId;
                document.getElementById('targetSession').value = sessionId;
                addEvent('system', 'Connected to SSE');
            };

            eventSource.onmessage = function(e) {
                addEvent('default', e.data);
            };

            eventSource.onerror = function(e) {
                document.getElementById('status').textContent = 'Disconnected';
                document.getElementById('status').className = 'status disconnected';
                addEvent('error', 'Connection error');
            };

            // Listen for common event types
            ['message', 'notification', 'alert', 'update', 'announcement'].forEach(function(type) {
                eventSource.addEventListener(type, function(e) {
                    addEvent(type, e.data);
                });
            });
        }

        function disconnect() {
            if (eventSource) {
                eventSource.close();
                eventSource = null;
                currentSessionId = null;
            }
            document.getElementById('status').textContent = 'Disconnected';
            document.getElementById('status').className = 'status disconnected';
            addEvent('system', 'Disconnected from SSE');
        }

        function sendToSession() {
            const targetSession = document.getElementById('targetSession').value || currentSessionId;
            const eventType = document.getElementById('eventType').value || 'message';
            const eventData = document.getElementById('eventData').value;
            const eventId = document.getElementById('eventId').value;

            if (!targetSession) {
                addEvent('error', 'No target session specified');
                return;
            }

            fetch('/send/' + targetSession, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    event: eventType,
                    data: eventData,
                    id: eventId
                })
            })
            .then(r => r.json())
            .then(res => {
                addEvent('api', 'Send response: ' + JSON.stringify(res));
            })
            .catch(err => {
                addEvent('error', 'Send failed: ' + err);
            });
        }

        function broadcast() {
            const eventType = document.getElementById('eventType').value || 'message';
            const eventData = document.getElementById('eventData').value;
            const eventId = document.getElementById('eventId').value;

            fetch('/broadcast', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    event: eventType,
                    data: eventData,
                    id: eventId
                })
            })
            .then(r => r.json())
            .then(res => {
                addEvent('api', 'Broadcast response: ' + JSON.stringify(res));
            })
            .catch(err => {
                addEvent('error', 'Broadcast failed: ' + err);
            });
        }

        function addEvent(type, data) {
            const events = document.getElementById('events');
            const div = document.createElement('div');
            div.className = 'event';
            div.innerHTML = '<strong>[' + type + ']</strong> ' + data + ' <small>(' + new Date().toLocaleTimeString() + ')</small>';
            events.appendChild(div);
            events.scrollTop = events.scrollHeight;
        }

        function clearEvents() {
            document.getElementById('events').innerHTML = '';
        }
    </script>
</body>
</html>`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// Example of sending events programmatically
func sendNotification(sessionID string, message string) {
	event, _ := sse.NewEventData("notification", message, fmt.Sprintf("%d", time.Now().UnixNano()))
	sessionManager.SendToSession(sessionID, event)
}

// Example of broadcasting to all users
func broadcastAnnouncement(message string) {
	event := sse.NewEventDataString("announcement", message, "")
	sessionManager.Broadcast(event)
}
