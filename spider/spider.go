package spider

func (sp *Spider) Run() {
	sp.ApiFinder()
	if sp.SearchEnable {
		sp.Search()
	}
}
