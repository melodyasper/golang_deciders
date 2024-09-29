package main

import "fmt"

/* State */
type StateData struct {
	Fitted        bool
	IsOn          bool
	RemainingUses uint
	IsBlown       bool
}

type State interface {
	GetData() *StateData
	IsValid() bool
}

type StateNotFitted struct {
	data StateData
}

type StateWorking struct {
	data StateData
}

type StateBlown struct {
	data StateData
}

type Aggregate interface {
	decide(Command) []Event
	evolve(Event)
}
type Bulb struct {
	state State
}



func (s StateNotFitted) GetData() *StateData {
	return &s.data
}
func (s StateWorking) GetData() *StateData {
	return &s.data
}
func (s StateBlown) GetData() *StateData {
	return &s.data
}

func (s StateNotFitted) IsValid() bool {
	data := s.GetData()
	return !data.Fitted && !data.IsBlown && data.RemainingUses == 0
}

func (s StateWorking) IsValid() bool {
	data := s.GetData()
	return data.Fitted && !data.IsBlown
}

func (s StateBlown) IsValid() bool {
	data := s.GetData()
	return data.Fitted && data.IsBlown  && !data.IsOn && data.RemainingUses == 0
}


// Event with isEvent marker method.
type Event interface {
	isEvent()
}

// Command with isCommand marker method.
type Command interface {
	isCommand() 
}

/* Commands */
type CommandFit struct {
	MaxUses uint
}
type CommandSwitchOn struct{}
type CommandSwitchOff struct{}

/* Events */
type EventFitted struct {
	MaxUses uint
}
type EventSwitchedOn struct{}
type EventSwitchedOff struct{}
type EventBlew struct{}

/* Event markers */
func (e EventFitted) isEvent() {}
func (e EventSwitchedOn) isEvent() {}
func (e EventSwitchedOff) isEvent() {}
func (e EventBlew) isEvent() {}

/* Command marker method implementations */
func (c CommandFit) isCommand() {}
func (c CommandSwitchOn) isCommand() {}
func (c CommandSwitchOff) isCommand() {}

func initial_state() State {
	state_data := StateData{Fitted: false, IsOn: false, RemainingUses: 0}
	return StateNotFitted{state_data}
}

func (b Bulb) decide(c Command) []Event {
	switch cv := c.(type) {
	case CommandFit:
		if _, sok := b.state.(StateNotFitted); sok {
			return []Event{EventFitted{MaxUses: cv.MaxUses}}
		}
	case CommandSwitchOn:
		if sv, sok := b.state.(StateWorking); sok {
			data := sv.GetData()
			if !data.IsOn {
				if data.RemainingUses > 0 {
					return []Event{EventSwitchedOn{}}
				}
				return []Event{EventBlew{}}
			}
		}
	case CommandSwitchOff:
		if sv, sok := b.state.(StateWorking); sok {
			if sv.GetData().IsOn {
				return []Event{EventSwitchedOff{}}
			}
		}
	}
	return []Event{}
}
func (b *Bulb) evolve(e Event) {
	/* Clone data from pointer */
	new_data := *(b.state.GetData())

	switch ev := e.(type) {
		case EventFitted:
			new_data.Fitted = true
			new_data.RemainingUses = ev.MaxUses
			b.state = StateWorking{new_data}
			return

		case EventSwitchedOn:
			new_data.RemainingUses -= 1
			new_data.IsOn = true
			b.state = StateWorking{new_data}
			return

		case EventSwitchedOff:
			new_data.IsOn = false
			b.state = StateWorking{new_data}
			return

		case EventBlew:
			new_data.IsBlown = true
			new_data.IsOn = false
			b.state = StateBlown{new_data}
			return
	}
}

func main() {
	// Setup commands
	command_fit := CommandFit{MaxUses: 1}
	command_switch_on := CommandSwitchOn{}
	command_switch_off := CommandSwitchOff{}
	commands := []Command{command_fit, command_switch_on, command_switch_off, command_switch_on}
	bulb := Bulb{state: initial_state()}
	for idx, command := range commands {
		fmt.Println("---------------")
		fmt.Println("Evaluating command ", idx)
		events := bulb.decide(command)
		if len(events) == 0 && idx < (len(commands) - 1) {
			fmt.Println("Unexpected end state")
		}
		for _, event := range events {
			bulb.evolve(event)
			fmt.Println("Is bulb state valid:", bulb.state.IsValid())
			fmt.Println("Is on              :", bulb.state.GetData().IsOn)
			fmt.Println("Uses remaining     :", bulb.state.GetData().RemainingUses)
			fmt.Println("Is bulb blown      :", bulb.state.GetData().IsBlown)
		}
	}
	
}
