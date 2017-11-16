package main

import "fmt"

type HumanI interface {
	Profile() string
}
type PhilosopherI interface {
	Thinking() string
}
type MusicianI interface {
	Create() string
}
type AthleteI interface {
	Compete() string
}

type Human struct{}

func (h *Human) Major() string { return "Unknown" }

type Philosopher struct{}

func (p *Philosopher) Thinking() string { return "Thinking the nature of the world." }
func (p *Philosopher) Major() string    { return "Philosopher" }

type Musician struct{}

func (m *Musician) Create() string { return "Creating good music." }
func (m *Musician) Major() string  { return "Musician" }

type Athlete struct{}

func (a *Athlete) Compete() string { return "For the champion." }
func (a *Athlete) Major() string   { return "Athlete" }

// Emulate Mixins.
type Person struct {
	Name string
	Human
	Philosopher
	Musician
	Athlete
}

func (p *Person) Profile() string {
	return fmt.Sprintf(`  name: %s
  major: %s,%s,%s
  hobbies:
    %s
    %s
    %s`, p.Name, p.Philosopher.Major(), p.Musician.Major(), p.Athlete.Major(),
		p.Thinking(), p.Create(), p.Compete())
}

// Implementing interfaces concerned.
type Student struct {
	Name string
	Human
}

func (s *Student) Thinking() string { return "Asking 'what','why', and 'how'." }
func (s *Student) Create() string   { return "Learning the music." }
func (s *Student) Compete() string  { return "Playing game." }
func (s *Student) Profile() string {
	h := &Hobby{}
	return fmt.Sprintf(`  name: %s
  major: %s
  hobbies:
    %s
    %s
    %s`, s.Name, s.Major(), h.Philosopher(PhilosopherI(s)), h.Musician(MusicianI(s)), h.Athlete(AthleteI(s)))
}

// interface methods caller.
type Hobby struct{}

func (h *Hobby) Philosopher(p PhilosopherI) string { return p.Thinking() }
func (h *Hobby) Musician(m MusicianI) string       { return m.Create() }
func (h *Hobby) Athlete(a AthleteI) string         { return a.Compete() }

type Profile struct{}

func (p Profile) Profile(h HumanI) {
	fmt.Println(h.Profile())
}

func main() {
	person := &Person{Name: "Bob"}
	Profile{}.Profile(HumanI(person))
	stu := &Student{Name: "Tom"}
	Profile{}.Profile(HumanI(stu))
}
