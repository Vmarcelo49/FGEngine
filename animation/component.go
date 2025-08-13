package animation

type ComponentManager struct {
	components map[int]*AnimationManager
	nextID     int
}

func NewComponentManager() *ComponentManager {
	return &ComponentManager{
		components: make(map[int]*AnimationManager),
		nextID:     1,
	}
}

func (cm *ComponentManager) CreateComponent(char *Character) *AnimationManager {
	component := CreateAnimationManager(cm.nextID, char)
	cm.components[cm.nextID] = component
	cm.nextID++
	return component
}

// UpdateAll updates all registered animation components
func (cm *ComponentManager) UpdateAll() {
	for _, component := range cm.components {
		component.Update()
	}
}

// RemoveComponent removes a component from the manager
func (cm *ComponentManager) RemoveComponent(id int) {
	delete(cm.components, id)
}

// GetComponent retrieves a component by ID
func (cm *ComponentManager) GetComponent(id int) *AnimationManager {
	return cm.components[id]
}

// GetComponentCount returns the number of active components
func (cm *ComponentManager) GetComponentCount() int {
	return len(cm.components)
}
