package chat

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	messages  []Message
	createdAt time.Time
}

func (c *Conversation) GetCreatedAt() time.Time {
	return c.createdAt
}

func (c *Conversation) GetMessages() []Message {
	return c.messages
}

func (c *Conversation) AddMessage(message Message) {
	c.messages = append(c.messages, message)
}

func (c *Conversation) RemoveMessage(id uuid.UUID) {
	for i, message := range c.messages {
		if message.ID == id {
			c.messages = append(c.messages[:i], c.messages[i+1:]...)
			break
		}
	}
}

func (c *Conversation) Reset() (*Conversation, error) {
	conversation := NewStackedConversation()

	// Copy system messages
	for _, message := range c.messages {
		if message.Role != RoleSystem {
			break
		}
		conversation.AddMessage(
			message,
		)
	}

	c.createdAt = conversation.GetCreatedAt()
	c.messages = conversation.GetMessages()

	return c, nil
}

func (c *Conversation) Clean() (*Conversation, error) {
	conversation := NewStackedConversation()

	c.createdAt = conversation.GetCreatedAt()
	c.messages = conversation.GetMessages()

	return c, nil
}

func NewStackedConversation() *Conversation {
	return &Conversation{
		messages:  make([]Message, 0),
		createdAt: time.Now(),
	}
}
