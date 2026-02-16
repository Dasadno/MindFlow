// Package world provides the TimeManager for simulation clock control.
//
// =============================================================================
// PURPOSE:
// =============================================================================
// The TimeManager controls the simulation's notion of time. It manages tick
// intervals, simulation speed, and provides the world clock that agents use
// to understand "when" things happen. Simulation time is decoupled from
// wall-clock time to allow speed-up, slow-down, and pause.
//
// =============================================================================
// TIME MANAGER STRUCTURE:
// =============================================================================
//
// type TimeManager struct {
//     currentTick    int64          // Monotonically increasing tick counter
//     tickDuration   time.Duration  // Base duration of one tick (e.g., 1 second)
//     speedMultiplier float64       // 0.5 = half speed, 2.0 = double speed
//     isPaused       bool
//     startTime      time.Time      // When simulation started
//     simTime        time.Time      // Current simulation time (may differ from wall clock)
//     ticker         *time.Ticker   // Controls tick frequency
//     mu             sync.Mutex
// }
//
// =============================================================================
// TIME OPERATIONS:
// =============================================================================
//
// func NewTimeManager(tickDuration time.Duration) *TimeManager
//   - Create time manager with specified base tick duration
//   - Default speed multiplier: 1.0
//   - Initialize simulation time to current wall clock
//
// func (tm *TimeManager) Advance() TickInfo
//   - Increment tick counter
//   - Advance simulation time by (tickDuration / speedMultiplier)
//   - Return TickInfo with current tick number, simulation time, delta
//
// func (tm *TimeManager) SetSpeed(multiplier float64)
//   - Adjust simulation speed
//   - Recalculate ticker interval: tickDuration / multiplier
//   - Clamp to reasonable range (0.1x to 10.0x)
//
// func (tm *TimeManager) Pause() / Resume()
//   - Pause/resume tick advancement
//   - When paused, Advance() blocks until resumed or Step() is called
//
// func (tm *TimeManager) Step()
//   - Advance exactly one tick while paused
//   - Useful for debugging and observation
//
// =============================================================================
// TIME QUERIES:
// =============================================================================
//
// type TickInfo struct {
//     Tick      int64         // Current tick number
//     SimTime   time.Time     // Current simulation time
//     Delta     time.Duration // Time elapsed since last tick
//     WallTime  time.Time     // Actual wall-clock time
// }
//
// func (tm *TimeManager) CurrentTick() int64
// func (tm *TimeManager) SimulationTime() time.Time
// func (tm *TimeManager) Uptime() time.Duration
// func (tm *TimeManager) TickChannel() <-chan TickInfo
//   - Returns channel that emits TickInfo at each tick
//   - Used by Orchestrator to drive the main loop
//
// =============================================================================
// SCHEDULED EVENTS:
// =============================================================================
//
// func (tm *TimeManager) ScheduleAt(tick int64, callback func())
//   - Register a callback to fire at a specific future tick
//   - Used for: memory consolidation, mood decay, random events
//
// func (tm *TimeManager) ScheduleEvery(interval int64, callback func())
//   - Register a recurring callback every N ticks
//   - Used for: periodic reflection, statistics collection

package world
