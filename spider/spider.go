package spider

//todo: add comment
func (sp *Spider) Run() {
	sp.ApiFinder()
	if sp.SearchEnable {
		sp.Search()
	}
}
