package utils



type ObjectsBox struct {
	objects map[string] interface{}
}

func NewObjectsBox() *ObjectsBox{
	return &ObjectsBox{ objects: make(map[string] interface{}), }
}

func (o *ObjectsBox) AddObject(key string, obj interface{}) {
	o.objects[key] = obj
}

func (o *ObjectsBox) GetObject(key string) interface{} {
	var exist bool
	var object interface{}

	if object, exist = o.objects[key]; !exist {
		return nil
	}
	return object
}