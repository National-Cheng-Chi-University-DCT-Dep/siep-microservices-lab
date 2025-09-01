package mqtt

import (
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"

	pkglogger "github.com/lipeichen/Ultimate-Security-Intelligence-Platform/backend/pkg/logger"
)

// MQTTClient MQTT客戶端介面
type MQTTClientInterface interface {
	Connect() error
	Disconnect()
	Publish(topic string, message interface{}) error
	Subscribe(topic string, handler MessageHandler) error
	Unsubscribe(topic string) error
	IsConnected() bool
}

// MessageHandler 訊息處理器
type MessageHandler func(topic string, payload []byte) error

// ThreatNotification 威脅通知結構
type ThreatNotification struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`        // created, updated, deleted, alert
	ThreatID    string                 `json:"threat_id"`
	IPAddress   string                 `json:"ip_address,omitempty"`
	Domain      string                 `json:"domain,omitempty"`
	ThreatType  string                 `json:"threat_type"`
	Severity    string                 `json:"severity"`
	RiskScore   string                 `json:"risk_score"`
	Source      string                 `json:"source"`
	Description string                 `json:"description"`
	Timestamp   time.Time              `json:"timestamp"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// SubscriptionFilter 訂閱篩選器
type SubscriptionFilter struct {
	ThreatTypes        []string `json:"threat_types,omitempty"`
	Severities         []string `json:"severities,omitempty"`
	MinConfidenceScore int      `json:"min_confidence_score,omitempty"`
	Sources            []string `json:"sources,omitempty"`
}

// MQTTClient MQTT客戶端實作
type MQTTClient struct {
	client   mqtt.Client
	broker   string
	clientID string
	username string
	password string
	qos      byte
	options  *mqtt.ClientOptions
}

// NewMQTTClient 建立MQTT客戶端
func NewMQTTClient(broker, username, password string) MQTTClientInterface {
	clientID := fmt.Sprintf("threat-intel-%s", uuid.New().String()[:8])
	
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetConnectTimeout(10 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(10 * time.Second)
	opts.SetCleanSession(true)

	// 設定連線處理器
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		pkglogger.Info("MQTT client connected", pkglogger.Fields{
			"broker":    broker,
			"client_id": clientID,
		})
	})

	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		pkglogger.Error("MQTT connection lost", pkglogger.Fields{
			"error":     err.Error(),
			"client_id": clientID,
		})
	})

	client := mqtt.NewClient(opts)

	return &MQTTClient{
		client:   client,
		broker:   broker,
		clientID: clientID,
		username: username,
		password: password,
		qos:      1, // At least once
		options:  opts,
	}
}

// Connect 連接到MQTT broker
func (m *MQTTClient) Connect() error {
	pkglogger.Info("Connecting to MQTT broker", pkglogger.Fields{
		"broker":    m.broker,
		"client_id": m.clientID,
	})

	token := m.client.Connect()
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}

	pkglogger.Info("Successfully connected to MQTT broker", pkglogger.Fields{
		"broker":    m.broker,
		"client_id": m.clientID,
	})

	return nil
}

// Disconnect 斷開MQTT連接
func (m *MQTTClient) Disconnect() {
	if m.client.IsConnected() {
		m.client.Disconnect(1000) // 1秒超時
		pkglogger.Info("MQTT client disconnected", pkglogger.Fields{
			"client_id": m.clientID,
		})
	}
}

// Publish 發布訊息
func (m *MQTTClient) Publish(topic string, message interface{}) error {
	if !m.client.IsConnected() {
		return fmt.Errorf("MQTT client not connected")
	}

	// 序列化訊息
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// 發布訊息
	token := m.client.Publish(topic, m.qos, false, payload)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to publish message: %w", token.Error())
	}

	pkglogger.Debug("Message published", pkglogger.Fields{
		"topic":      topic,
		"payload_len": len(payload),
	})

	return nil
}

// Subscribe 訂閱主題
func (m *MQTTClient) Subscribe(topic string, handler MessageHandler) error {
	if !m.client.IsConnected() {
		return fmt.Errorf("MQTT client not connected")
	}

	// 包裝處理器
	messageHandler := func(client mqtt.Client, msg mqtt.Message) {
		pkglogger.Debug("Message received", pkglogger.Fields{
			"topic":   msg.Topic(),
			"payload": string(msg.Payload()),
		})

		if err := handler(msg.Topic(), msg.Payload()); err != nil {
			pkglogger.Error("Failed to handle MQTT message", pkglogger.Fields{
				"error": err.Error(),
				"topic": msg.Topic(),
			})
		}
	}

	// 訂閱主題
	token := m.client.Subscribe(topic, m.qos, messageHandler)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, token.Error())
	}

	pkglogger.Info("Subscribed to MQTT topic", pkglogger.Fields{
		"topic": topic,
	})

	return nil
}

// Unsubscribe 取消訂閱
func (m *MQTTClient) Unsubscribe(topic string) error {
	if !m.client.IsConnected() {
		return fmt.Errorf("MQTT client not connected")
	}

	token := m.client.Unsubscribe(topic)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to unsubscribe from topic %s: %w", topic, token.Error())
	}

	pkglogger.Info("Unsubscribed from MQTT topic", pkglogger.Fields{
		"topic": topic,
	})

	return nil
}

// IsConnected 檢查連接狀態
func (m *MQTTClient) IsConnected() bool {
	return m.client.IsConnected()
}

// ThreatNotificationPublisher 威脅通知發布器
type ThreatNotificationPublisher struct {
	client MQTTClientInterface
}

// NewThreatNotificationPublisher 建立威脅通知發布器
func NewThreatNotificationPublisher(client MQTTClientInterface) *ThreatNotificationPublisher {
	return &ThreatNotificationPublisher{
		client: client,
	}
}

// PublishThreatCreated 發布威脅建立通知
func (p *ThreatNotificationPublisher) PublishThreatCreated(threat *ThreatNotification) error {
	threat.Type = "created"
	threat.Timestamp = time.Now()
	
	topic := fmt.Sprintf("threats/created/%s", threat.Severity)
	return p.client.Publish(topic, threat)
}

// PublishThreatUpdated 發布威脅更新通知
func (p *ThreatNotificationPublisher) PublishThreatUpdated(threat *ThreatNotification) error {
	threat.Type = "updated"
	threat.Timestamp = time.Now()
	
	topic := fmt.Sprintf("threats/updated/%s", threat.Severity)
	return p.client.Publish(topic, threat)
}

// PublishThreatDeleted 發布威脅刪除通知
func (p *ThreatNotificationPublisher) PublishThreatDeleted(threatID, threatType, severity string) error {
	notification := &ThreatNotification{
		ID:         uuid.New().String(),
		Type:       "deleted",
		ThreatID:   threatID,
		ThreatType: threatType,
		Severity:   severity,
		Timestamp:  time.Now(),
	}
	
	topic := fmt.Sprintf("threats/deleted/%s", severity)
	return p.client.Publish(topic, notification)
}

// PublishHighRiskAlert 發布高風險警報
func (p *ThreatNotificationPublisher) PublishHighRiskAlert(threat *ThreatNotification) error {
	threat.Type = "alert"
	threat.Timestamp = time.Now()
	
	// 高風險威脅發送到特殊主題
	topic := "threats/alerts/high-risk"
	return p.client.Publish(topic, threat)
}

// ThreatSubscriber 威脅訂閱器
type ThreatSubscriber struct {
	client  MQTTClientInterface
	filters map[string]SubscriptionFilter
}

// NewThreatSubscriber 建立威脅訂閱器
func NewThreatSubscriber(client MQTTClientInterface) *ThreatSubscriber {
	return &ThreatSubscriber{
		client:  client,
		filters: make(map[string]SubscriptionFilter),
	}
}

// SubscribeToThreats 訂閱威脅通知
func (s *ThreatSubscriber) SubscribeToThreats(subscriberID string, filter SubscriptionFilter, handler MessageHandler) error {
	// 儲存篩選器
	s.filters[subscriberID] = filter

	// 包裝處理器以應用篩選器
	filteredHandler := func(topic string, payload []byte) error {
		// 解析通知
		var notification ThreatNotification
		if err := json.Unmarshal(payload, &notification); err != nil {
			pkglogger.Error("Failed to unmarshal threat notification", pkglogger.Fields{
				"error": err.Error(),
			})
			return err
		}

		// 應用篩選器
		if s.shouldForwardNotification(subscriberID, &notification) {
			return handler(topic, payload)
		}

		return nil
	}

	// 訂閱相關主題
	topics := s.generateTopicsFromFilter(filter)
	for _, topic := range topics {
		if err := s.client.Subscribe(topic, filteredHandler); err != nil {
			return fmt.Errorf("failed to subscribe to topic %s: %w", topic, err)
		}
	}

	pkglogger.Info("Subscribed to threat notifications", pkglogger.Fields{
		"subscriber_id": subscriberID,
		"topics":        topics,
	})

	return nil
}

// UnsubscribeFromThreats 取消訂閱威脅通知
func (s *ThreatSubscriber) UnsubscribeFromThreats(subscriberID string) error {
	filter, exists := s.filters[subscriberID]
	if !exists {
		return fmt.Errorf("subscriber %s not found", subscriberID)
	}

	// 取消訂閱相關主題
	topics := s.generateTopicsFromFilter(filter)
	for _, topic := range topics {
		if err := s.client.Unsubscribe(topic); err != nil {
			pkglogger.Error("Failed to unsubscribe from topic", pkglogger.Fields{
				"error": err.Error(),
				"topic": topic,
			})
		}
	}

	// 移除篩選器
	delete(s.filters, subscriberID)

	pkglogger.Info("Unsubscribed from threat notifications", pkglogger.Fields{
		"subscriber_id": subscriberID,
		"topics":        topics,
	})

	return nil
}

// shouldForwardNotification 檢查是否應該轉發通知
func (s *ThreatSubscriber) shouldForwardNotification(subscriberID string, notification *ThreatNotification) bool {
	filter, exists := s.filters[subscriberID]
	if !exists {
		return true // 如果沒有篩選器，轉發所有通知
	}

	// 檢查威脅類型
	if len(filter.ThreatTypes) > 0 {
		found := false
		for _, threatType := range filter.ThreatTypes {
			if threatType == notification.ThreatType {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// 檢查嚴重程度
	if len(filter.Severities) > 0 {
		found := false
		for _, severity := range filter.Severities {
			if severity == notification.Severity {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// 檢查來源
	if len(filter.Sources) > 0 {
		found := false
		for _, source := range filter.Sources {
			if source == notification.Source {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// generateTopicsFromFilter 根據篩選器生成主題列表
func (s *ThreatSubscriber) generateTopicsFromFilter(filter SubscriptionFilter) []string {
	var topics []string

	// 如果沒有指定嚴重程度，訂閱所有
	severities := filter.Severities
	if len(severities) == 0 {
		severities = []string{"low", "medium", "high", "critical"}
	}

	// 生成主題
	for _, severity := range severities {
		topics = append(topics, fmt.Sprintf("threats/created/%s", severity))
		topics = append(topics, fmt.Sprintf("threats/updated/%s", severity))
		topics = append(topics, fmt.Sprintf("threats/deleted/%s", severity))
	}

	// 總是訂閱高風險警報
	topics = append(topics, "threats/alerts/high-risk")

	return topics
} 