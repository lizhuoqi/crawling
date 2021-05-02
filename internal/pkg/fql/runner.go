package fql

var ferret *Ferret

func init() {
	// initial Ferret instance
	ferret = newFerret()
}

// get the ferret instance
func GetFerret() *Ferret {
	return ferret
}
