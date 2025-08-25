package animation

type AnimationRegistry struct {
	components map[int]*AnimationSystem
	nextID     int
}

func NewAnimationRegistry() *AnimationRegistry {
	return &AnimationRegistry{
		components: make(map[int]*AnimationSystem),
		nextID:     1,
	}
}

func (ar *AnimationRegistry) CreateComponent(char *Character) *AnimationSystem {
	component := NewAnimationSystem(ar.nextID, char)
	ar.components[ar.nextID] = component
	ar.nextID++
	return component
}

// UpdateAll updates all registered animation components
func (ar *AnimationRegistry) UpdateAll() {
	for _, component := range ar.components {
		component.Update()
	}
}

// RemoveComponent removes a component from the manager
func (ar *AnimationRegistry) RemoveComponent(id int) {
	delete(ar.components, id)
}

// GetComponent retrieves a component by ID
func (ar *AnimationRegistry) GetComponent(id int) *AnimationSystem {
	return ar.components[id]
}

// GetComponentCount returns the number of active components
func (ar *AnimationRegistry) GetComponentCount() int {
	return len(ar.components)
}
