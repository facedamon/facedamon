package arraylist

import "encoding/json"

func (list *List) ToJSON() ([]byte, error) {
	return json.Marshal(list.elements[:])
}

func (list *List) FromJSON(data []byte) error {
	err := json.Unmarshal(data, &list.elements)
	if err == nil {
		list.size = len(list.elements)
	}
	return err
}
