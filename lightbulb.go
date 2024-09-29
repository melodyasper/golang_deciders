package main

import "fmt"

/************************ EVENTS  ****************/
// Event interface with isEvent marker method.
type Event interface {
	isEvent()
}

// Event defintions
type EventFitted struct {
	MaxUses uint
}
type EventSwitchedOn struct{}
type EventSwitchedOff struct{}
type EventBlew struct{}

// Event markers
func (e EventFitted) isEvent() {}
func (e EventSwitchedOn) isEvent() {}
func (e EventSwitchedOff) isEvent() {}
func (e EventBlew) isEvent() {}

// Implement the Stringer (%s formatter) interfaces for Events
func (e EventFitted) String() string {
	return fmt.Sprintf("EventFitted {MaxUses: %d}", e.MaxUses)
}
func (e EventSwitchedOn) String() string {
	return "EventSwitchedOn {}"
}
func (e EventSwitchedOff) String() string {
	return "EventSwitchedOff {}"
}
func (e EventBlew) String() string {
	return "EventBlew {}"
}


/************************ STATES  ****************/
// State interface with marker function
type State interface {
	isState()
}

// State definitions
type StateNotFitted struct {}
type StateWorking struct {
	IsOn          bool
	RemainingUses uint
}

type StateBlown struct {}

// State marker implementation
func (s StateNotFitted) isState() {}
func (s StateWorking) isState() {}
func (s StateBlown) isState() {}

// Implement the Stringer (%s formatter) interfaces for States
func (s StateNotFitted) String() string {
	return "StateNotFitted {}"
}
func (s StateWorking) String() string {
	return fmt.Sprintf("StateWorking {IsOn: %t, RemainingUses: %d}", s.IsOn, s.RemainingUses)
}
func (s StateBlown) String() string {
	return "StateBlown {}"
}

/************************ COMMANDS ****************/
// Command interface isCommand marker method.
type Command interface {
	isCommand() 
}
// Command Defintions
type CommandFit struct {
	MaxUses uint
}
type CommandSwitchOn struct{}
type CommandSwitchOff struct{}

// Command marker method implementations
func (c CommandFit) isCommand() {}
func (c CommandSwitchOn) isCommand() {}
func (c CommandSwitchOff) isCommand() {}

// Implement the Stringer (%s formatter) interfaces for Commands
func (c CommandFit) String() string {
	return fmt.Sprintf("CommandFit {MaxUses: %d}", c.MaxUses)
}
func (s CommandSwitchOn) String() string {
	return "CommandSwitchOn {}"
}
func (s CommandSwitchOff) String() string {
	return "CommandSwitchOff {}"
}


/************************ AGGREGATE AND BULB **************/
type Aggregate interface {
	decide(Command) []Event
	evolve(Event)
}
type Bulb struct {
	state State
}

func (b Bulb) decide(c Command) []Event {
	switch cv := c.(type) {
	case CommandFit:
		if _, sok := b.state.(StateNotFitted); sok {
			return []Event{EventFitted{MaxUses: cv.MaxUses}}
		}
	case CommandSwitchOn:
		if sv, sok := b.state.(StateWorking); sok {
			if !sv.IsOn {
				if sv.RemainingUses > 0 {
					return []Event{EventSwitchedOn{}}
				}
				return []Event{EventBlew{}}
			}
		}
	case CommandSwitchOff:
		if sv, sok := b.state.(StateWorking); sok {
			if sv.IsOn {
				return []Event{EventSwitchedOff{}}
			}
		}
	}
	return []Event{}
}
func (b *Bulb) evolve(e Event) {
	/* Clone data from pointer */
	switch ev := e.(type) {
		case EventFitted:
			new_state := StateWorking { RemainingUses: ev.MaxUses, IsOn: false}
			b.state = new_state
			return

		case EventSwitchedOn:
			if sv, sok := b.state.(StateWorking); sok {
				new_state := sv
				new_state.IsOn = true
				new_state.RemainingUses -= 1
				b.state = new_state
			}
			return

		case EventSwitchedOff:
			if sv, sok := b.state.(StateWorking); sok {
				new_state := sv
				new_state.IsOn = false
				b.state = new_state
			}
			return

		case EventBlew:
			b.state = StateBlown{}
			return
	}
}


// Helper function
func initial_state() State {
	return StateNotFitted{}
}


func main() {
	// Setup commands
	command_fit := CommandFit{MaxUses: 1}
	command_switch_on := CommandSwitchOn{}
	command_switch_off := CommandSwitchOff{}
	commands := []Command{command_fit, command_switch_on, command_switch_off, command_switch_on}
	bulb := Bulb{state: initial_state()}
	for idx, command := range commands {
		fmt.Printf("Evaluating command: %s\n", command)
		events := bulb.decide(command)
		if len(events) == 0 && idx < (len(commands) - 1) {
			fmt.Println("Unexpected end state")
		}
		for _, event := range events {
			fmt.Printf("Event             : %s\n", event)
			bulb.evolve(event)
			fmt.Printf("Results in state  : %s\n", bulb.state)
		}
	}
	
}
