package controllers

type Topic struct {
	*revel.Controller
}

func (t *Topic) Index() {
	topics, err := models.TopicGetAll()

	return t.Render(topics)
}

func (t *Topic) Get(code string) {
	topic, ok, err := models.TopicGet(code)

	return t.Render(topic)
}
