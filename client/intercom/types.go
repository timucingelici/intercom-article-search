package intercom

type Canvas struct {
	Canvas CanvasStatic `json:"canvas"`
}

type CanvasStatic struct {
	Content Content `json:"content"`
	StoredData interface{} `json:"stored_data"`
}

type Content struct {
	Components []interface{} `json:"components"`
}

type InputComponent struct {
	Type 		string	`json:"type"`
	Id			string	`json:"id"`
	Label		string	`json:"label"`
	Placeholder	string	`json:"placeholder"`
	Value		string	`json:"value"`
	SaveState	string	`json:"save_state"`
	Disabled	bool	`json:"disabled"`
	Action 		SubmitAction `json:"action"`
}

type SubmitAction struct {
	Type 	string	`json:"type"`
}

type SheetAction struct {
	Type	string `json:"type"`
	Url		string `json:"url"`
}

type Request struct {
	AppId	string `json:"app_id"`
	CurrentCanvas	struct{
		Content	Content
	} `json:"current_canvas"`
	ComponentId	string `json:"component_id"`
	InputValues map[string]string `json:"input_values"`
	User struct {
		CustomAttributes struct {
			Region	string `json:"region"`
			Role	string `json:"role"`
			Roles	string `json:"roles"`
			Source	string `json:"source"`
			Tenant	string `json:"tenant"`
		} `json:"custom_attributes"`
	} `json:"user"`

}

type SearchRequest struct {
	ComponentId string `json:"component_id"`
	InputValues struct {
		Query string `json:"query"`
	} `json:"input_values"`
}

type ListResponse struct {
	Type 		string			`json:"type"`
	Disabled 	bool			`json:"disabled"`
	Items 		[]ListResponseItem	`json:"items"`
}

type ListResponseItem struct {
	Type 		string `json:"type"`
	Id			string `json:"id"`
	Title		string `json:"title"`
	Subtitle	string `json:"subtitle"`
	Image		string `json:"image"`
	ImageWidth	int `json:"image_width"`
	ImageHeight	int `json:"image_height"`
	RoundedImage	bool `json:"rounded_image"`
	Disabled	bool `json:"disabled"`
	Action		interface{} `json:"action"`
}

type ConfigResults struct {
	Results map[string]string `json:"results"`
}
