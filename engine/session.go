package engine

import "time"

type Phase int

const (
	PhaseWork Phase = iota
	PhaseShortBreak
	PhaseLongBreak
	PhaseDone
)

func (p Phase) String() string {
	switch p {
	case PhaseWork:
		return "Work"
	case PhaseShortBreak:
		return "Short Break"
	case PhaseLongBreak:
		return "Long Break"
	case PhaseDone:
		return "Done"
	default:
		return "Unknown"
	}
}

type Config struct {
	WorkDuration       time.Duration
	ShortBreakDuration time.Duration
	LongBreakDuration  time.Duration
	LongBreakEvery     int
	TotalCycles        int
}

// Handles state transitions
type Session struct {
	config         Config
	currentPhase   Phase
	cyclesComplete int
	totalPhases    int
	phasesComplete int
}

func NewSession(cfg Config) *Session {
	s := &Session{
		config:       cfg,
		currentPhase: PhaseWork,
	}
	s.totalPhases = s.calculateTotalPhases()
	return s
}

func (s *Session) calculateTotalPhases() int {
	if s.config.TotalCycles == 0 {
		return 0
	}

	cycles := s.config.TotalCycles
	phases := cycles

	if s.config.LongBreakEvery > 0 {
		longBreaks := (cycles - 1) / s.config.LongBreakEvery
		shortBreaks := cycles - 1 - longBreaks
		phases += longBreaks + shortBreaks
	} else {
		phases += cycles - 1
	}

	return phases
}

func (s *Session) CurrentPhase() Phase { return s.currentPhase }
func (s *Session) CyclesComplete() int { return s.cyclesComplete }
func (s *Session) TotalCycles() int    { return s.config.TotalCycles }
func (s *Session) TotalPhases() int    { return s.totalPhases }
func (s *Session) PhasesComplete() int { return s.phasesComplete }

func (s *Session) PhaseDuration() time.Duration {
	switch s.currentPhase {
	case PhaseWork:
		return s.config.WorkDuration
	case PhaseShortBreak:
		return s.config.ShortBreakDuration
	case PhaseLongBreak:
		return s.config.LongBreakDuration
	default:
		return 0
	}
}

func (s *Session) NextPhase() Phase {
	if s.currentPhase == PhaseDone {
		return PhaseDone
	}

	s.phasesComplete++

	switch s.currentPhase {
	case PhaseWork:
		s.cyclesComplete++

		if s.config.TotalCycles > 0 && s.cyclesComplete >= s.config.TotalCycles {
			s.currentPhase = PhaseDone
			return s.currentPhase
		}

		if s.config.LongBreakEvery > 0 && s.cyclesComplete%s.config.LongBreakEvery == 0 {
			s.currentPhase = PhaseLongBreak
		} else {
			s.currentPhase = PhaseShortBreak
		}

	case PhaseShortBreak, PhaseLongBreak:
		s.currentPhase = PhaseWork
	}

	return s.currentPhase
}
