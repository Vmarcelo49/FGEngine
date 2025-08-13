package animation

type AnimationRegistry struct {
	components map[int]*AnimationManager
	nextID     int
}

func NewAnimationRegistry() *AnimationRegistry {
	return &AnimationRegistry{
		components: make(map[int]*AnimationManager),
		nextID:     1,
	}
}

func (ar *AnimationRegistry) CreateComponent(char *Character) *AnimationManager {
	component := CreateAnimationManager(ar.nextID, char)
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
func (ar *AnimationRegistry) GetComponent(id int) *AnimationManager {
	return ar.components[id]
}

// GetComponentCount returns the number of active components
func (ar *AnimationRegistry) GetComponentCount() int {
	return len(ar.components)
}
