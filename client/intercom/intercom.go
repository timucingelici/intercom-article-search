package intercom

func (c *Canvas) AddComponent(component interface{}) *Canvas {
	c.Canvas.Content.Components = append(c.Canvas.Content.Components, component)
	return c
}


func NewCanvas() Canvas {
	return Canvas{
		CanvasStatic {
			Content{},
			nil,
		},
	}
}

func NewInput(label string, placeholder string) InputComponent {
	return InputComponent{
		"input",
		"query",
		label,
		placeholder,
		"",
		"unsaved",
		false,
		SubmitAction {
			"submit",
		},
	}
}