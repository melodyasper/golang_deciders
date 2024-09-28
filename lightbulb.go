package main

import "fmt"

/* State */
type State struct {
	Fitted        bool
	IsOn          bool
	RemainingUses int
	IsBlown       bool
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
	MaxUses int
}
type CommandSwitchOn struct{}
type CommandSwitchOff struct{}

/* Events */
type EventFitted struct {
	MaxUses int
}
type EventSwitchedOn struct{}
type EventSwitchedOff struct{}
type EventBlew struct{}

/* Event logic */
func (e EventFitted) evolve(s State) State {
	return State{Fitted: true, IsOn: false, RemainingUses: e.MaxUses}
}

func (e EventSwitchedOn) evolve(s State) State {
	return State{Fitted: true, IsOn: true, RemainingUses: s.RemainingUses - 1}
}

func (e EventSwitchedOff) evolve(s State) State {
	return State{Fitted: true, IsOn: false, RemainingUses: s.RemainingUses}
}

func (e EventBlew) evolve(s State) State {
	return State{Fitted: true, IsOn: false, RemainingUses: s.RemainingUses, IsBlown: true}
}

/* Commands logic */
func (c CommandFit) decide(s State) Event {
	return EventFitted{MaxUses: c.MaxUses}
}

func (c CommandSwitchOn) decide(s State) Event {
	if s.RemainingUses <= 0 {
		return EventBlew{}
	}
	return EventSwitchedOn{}
}

func (c CommandSwitchOff) decide(s State) Event {
	return EventSwitchedOff{}
}

func main() {
	state := State{Fitted: false, IsOn: false, RemainingUses: 0}
	// Setup commands
	command_fit := CommandFit{MaxUses: 1}
	command_switch_on := CommandSwitchOn{}
	command_switch_off := CommandSwitchOff{}

	// Start
	event := command_fit.decide(state)
	state = event.evolve(state)

	fmt.Println("---------------")
	fmt.Println("Is on         :", state.IsOn)
	fmt.Println("Uses remaining:", state.RemainingUses)
	fmt.Println("Is bulb blown :", state.IsBlown)

	// Switch on
	event = command_switch_on.decide(state)
	state = event.evolve(state)

	fmt.Println("---------------")
	fmt.Println("Is on         :", state.IsOn)
	fmt.Println("Uses remaining:", state.RemainingUses)
	fmt.Println("Is bulb blown :", state.IsBlown)

	// Switch off
	event = command_switch_off.decide(state)
	state = event.evolve(state)

	fmt.Println("---------------")
	fmt.Println("Is on         :", state.IsOn)
	fmt.Println("Uses remaining:", state.RemainingUses)
	fmt.Println("Is bulb blown :", state.IsBlown)

	// Switch on
	event = command_switch_on.decide(state)
	state = event.evolve(state)

	fmt.Println("---------------")
	fmt.Println("Is on         :", state.IsOn)
	fmt.Println("Uses remaining:", state.RemainingUses)
	fmt.Println("Is bulb blown :", state.IsBlown)

}
