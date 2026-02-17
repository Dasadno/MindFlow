// Package world provides the simulation orchestrator.
//
// Orchestrator — центральный координатор симуляции.
// Каждые 15 секунд берёт двух случайных агентов и запускает диалог.

package world

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"milk/server/internal/agent"
	"milk/server/internal/api"
	"milk/server/internal/storage"
	"milk/server/pkg/llm"
)

// Orchestrator — тикер диалогов между агентами.
type Orchestrator struct {
	repo         *storage.Repository
	llm          *llm.Client
	hub          *api.Hub
	tickInterval time.Duration
	turns        int // реплик за тик
	currentTick  int64
	mu           sync.Mutex
	cancel       context.CancelFunc
}

// NewOrchestrator создаёт Orchestrator.
func NewOrchestrator(repo *storage.Repository, llmClient *llm.Client, hub *api.Hub) *Orchestrator {
	return &Orchestrator{
		repo:         repo,
		llm:          llmClient,
		hub:          hub,
		tickInterval: 15 * time.Second,
		turns:        4,
	}
}

// Start запускает фоновый тикер. Блокирует до отмены ctx.
func (o *Orchestrator) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	o.cancel = cancel

	log.Println("orchestrator: started, tick interval", o.tickInterval)

	ticker := time.NewTicker(o.tickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			o.mu.Lock()
			o.currentTick++
			tick := o.currentTick
			o.mu.Unlock()

			go o.runTick(ctx, tick)

		case <-ctx.Done():
			log.Println("orchestrator: stopped")
			return
		}
	}
}

// Stop останавливает тикер.
func (o *Orchestrator) Stop() {
	if o.cancel != nil {
		o.cancel()
	}
}

func (o *Orchestrator) runTick(ctx context.Context, tick int64) {
	agents, err := o.repo.GetRandomActiveAgents(2)
	if err != nil {
		log.Printf("orchestrator tick %d: getAgents error: %v", tick, err)
		return
	}
	if len(agents) < 2 {
		log.Printf("orchestrator tick %d: not enough agents (%d)", tick, len(agents))
		return
	}

	log.Printf("orchestrator tick %d: %s <-> %s", tick, agents[0].Name, agents[1].Name)
	o.runConversation(ctx, agents[0], agents[1], tick)
}

func (o *Orchestrator) runConversation(
	ctx context.Context,
	a1, a2 storage.AgentRecord,
	tick int64,
) {
	p1 := parsePersonality(a1.Personality)
	p2 := parsePersonality(a2.Personality)

	brain1 := agent.NewBrain(&p1)
	brain2 := agent.NewBrain(&p2)

	mood1 := agent.MoodNeutral
	mood2 := agent.MoodNeutral

	goals1 := parseGoals(a1.Goals)
	goals2 := parseGoals(a2.Goals)

	var history1, history2 []llm.Message

	// Начальное сообщение: агент 1 видит агента 2
	opener := fmt.Sprintf(
		"You notice %s nearby. Start a conversation naturally based on your personality and current mood.",
		a2.Name,
	)
	history1 = append(history1, llm.Message{Role: "user", Content: opener})

	for i := 0; i < o.turns; i++ {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if i%2 == 0 {
			// Ход агента 1
			o.injectHumanMessages(&history1, a1.ID)

			reply, err := brain1.Think(ctx, o.llm, a1.Name, mood1, goals1, history1)
			if err != nil {
				log.Printf("orchestrator: agent1 error: %v", err)
				return
			}

			o.saveAndBroadcast(a1, a2, reply, tick)

			history1 = append(history1, llm.Message{Role: "assistant", Content: reply})
			history2 = append(history2, llm.Message{
				Role:    "user",
				Content: fmt.Sprintf("%s says: %s", a1.Name, reply),
			})
		} else {
			// Ход агента 2
			o.injectHumanMessages(&history2, a2.ID)

			if len(history2) == 0 {
				history2 = append(history2, llm.Message{
					Role:    "user",
					Content: fmt.Sprintf("%s is talking to you. Respond naturally.", a1.Name),
				})
			}

			reply, err := brain2.Think(ctx, o.llm, a2.Name, mood2, goals2, history2)
			if err != nil {
				log.Printf("orchestrator: agent2 error: %v", err)
				return
			}

			o.saveAndBroadcast(a2, a1, reply, tick)

			history2 = append(history2, llm.Message{Role: "assistant", Content: reply})
			history1 = append(history1, llm.Message{
				Role:    "user",
				Content: fmt.Sprintf("%s says: %s", a2.Name, reply),
			})
		}
	}
}

func (o *Orchestrator) injectHumanMessages(history *[]llm.Message, agentID string) {
	injections := o.hub.DrainInjections(agentID)
	for _, inj := range injections {
		*history = append(*history, llm.Message{
			Role:    "user",
			Content: fmt.Sprintf("[Human says to you]: %s", inj),
		})
	}
}

func (o *Orchestrator) saveAndBroadcast(speaker, target storage.AgentRecord, reply string, tick int64) {
	_ = o.repo.SaveConversationEvent(speaker.ID, target.ID, reply, tick)
	o.hub.Broadcast(api.SSEEvent{
		Type:    "conversation",
		Speaker: speaker.Name,
		Target:  target.Name,
		Content: reply,
		AgentID: speaker.ID,
		Tick:    tick,
	})
}

func parsePersonality(raw string) agent.Personality {
	var p agent.Personality
	if raw != "" {
		json.Unmarshal([]byte(raw), &p)
	}
	return p
}

func parseGoals(goals sql.NullString) []agent.Goal {
	if !goals.Valid || goals.String == "" {
		return nil
	}
	var g []agent.Goal
	json.Unmarshal([]byte(goals.String), &g)
	return g
}
