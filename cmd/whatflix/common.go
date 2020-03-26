package main

func removeDuplicate(items []*creditsData) []string {
	var key = make(map[string]bool)
	lists := []string{}
	for _, item := range items {
		entry := *item
		if _, value := key[entry.Title]; !value {
			key[entry.Title] = true
			lists = append(lists, entry.Title)
		}
	}
	return lists
}
