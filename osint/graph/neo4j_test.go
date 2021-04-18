package graph

import "testing"

func Test_Graph(t *testing.T) {
	neo, _ := New("bolt://localhost:7687", "neo4j", "changeit")

	if _, err := neo.AdjustItem(CleanAll, map[string]interface{}{}); err != nil {
		t.Fatal(err)
	}
	d := map[string]interface{}{
		"id":     "0",
		"name":   "user1",
		"social": "twitter",
	}
	if _, err := neo.AdjustItem(AddNode, d); err != nil {
		t.Fatal(err)
	}

	d = map[string]interface{}{
		"id":     "1",
		"name":   "user2",
		"social": "twitter",
	}
	if _, err := neo.AdjustItem(AddNode, d); err != nil {
		t.Fatal(err)
	}
	d = map[string]interface{}{
		"id1":       "0",
		"name1":     "user1",
		"id2":       "1",
		"name2":     "user2",
		"condition": "n1.name = $name1 AND n1.id_str = $id1 AND n2.name = $name2 AND n2.id_str = $id2",
	}
	if _, err := neo.AdjustItem(AddRelation, d); err != nil {
		t.Fatal(err)
	}
}
