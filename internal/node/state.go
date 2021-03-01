package node

import "time"

type state struct {
	curDifficulty     int
	lastMineTimestamp time.Time
}

func newState() state {
	return state{STARTING_DIFFICULTY, time.Now()}
}

func (s *state) adjustDifficulty(time time.Time) {
	if time.Sub(s.lastMineTimestamp) < MINING_RATE {
		if s.curDifficulty != MAX_DIFFICULTY {
			s.curDifficulty++
		}
	} else if s.curDifficulty > MIN_DIFFICULTY {
		s.curDifficulty--
	}
}

func (s *state) setLastMineTime() {
	s.lastMineTimestamp = time.Now()
}
