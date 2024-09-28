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


// We are going to need to attach Evolve to state
type Event interface {
	evolve(State) State
}

// We are going to attach Command to state
type Command interface {
	decide(State) Event
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

/* Event logic */
func (e EventFitted) evolve(s State) State {
	new_data := StateData{Fitted: true, IsOn: false, RemainingUses: e.MaxUses}
	return StateWorking{new_data}
}

func (e EventSwitchedOn) evolve(s State) State {
	current_data := s.GetData()
	new_data := StateData{Fitted: true, IsOn: true, RemainingUses: current_data.RemainingUses - 1}
	return StateWorking{new_data}
}

func (e EventSwitchedOff) evolve(s State) State {
	current_data := s.GetData()
	new_data := StateData{Fitted: true, IsOn: false, RemainingUses: current_data.RemainingUses}
	return StateWorking{new_data}
}

func (e EventBlew) evolve(s State) State {
	new_data := StateData{Fitted: true, IsBlown: true, IsOn: false, RemainingUses: 0}

	return StateBlown{new_data}
}

/* Commands logic */
func (c CommandFit) decide(s State) Event {
	return EventFitted{MaxUses: c.MaxUses}
}

func (c CommandSwitchOn) decide(s State) Event {
	if s.GetData().RemainingUses <= 0 {
		return EventBlew{}
	}
	return EventSwitchedOn{}
}

func (c CommandSwitchOff) decide(s State) Event {
	return EventSwitchedOff{}
}

func initial_state() State {
	state_data := StateData{Fitted: false, IsOn: false, RemainingUses: 0}
	return StateNotFitted{state_data}
}

func main() {
	state := initial_state()
	// Setup commands
	command_fit := CommandFit{MaxUses: 1}
	command_switch_on := CommandSwitchOn{}
	command_switch_off := CommandSwitchOff{}
	
	// Start
	event := command_fit.decide(state)
	state = event.evolve(state)

	fmt.Println("---------------")
	fmt.Println("Is on         :", state.GetData().IsOn)
	fmt.Println("Uses remaining:", state.GetData().RemainingUses)
	fmt.Println("Is bulb blown :", state.GetData().IsBlown)

	// Switch on
	event = command_switch_on.decide(state)
	state = event.evolve(state)

	fmt.Println("---------------")
	fmt.Println("Is on         :", state.GetData().IsOn)
	fmt.Println("Uses remaining:", state.GetData().RemainingUses)
	fmt.Println("Is bulb blown :", state.GetData().IsBlown)

	// Switch off
	event = command_switch_off.decide(state)
	state = event.evolve(state)

	fmt.Println("---------------")
	fmt.Println("Is on         :", state.GetData().IsOn)
	fmt.Println("Uses remaining:", state.GetData().RemainingUses)
	fmt.Println("Is bulb blown :", state.GetData().IsBlown)

	// Switch on
	event = command_switch_on.decide(state)
	state = event.evolve(state)

	fmt.Println("---------------")
	fmt.Println("Is on         :", state.GetData().IsOn)
	fmt.Println("Uses remaining:", state.GetData().RemainingUses)
	fmt.Println("Is bulb blown :", state.GetData().IsBlown)

}
