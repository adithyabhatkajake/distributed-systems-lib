package crypto

var mapType = make(map[string]PKIAlgo)

// AddPKIAlgo A function that implements the interface must call this function to register their Algorithm
func AddPKIAlgo(name string, ctx PKIAlgo) {
	mapType[name] = ctx
}

// CheckExists checks if the algorithm is registered in the interface
func CheckExists(name string) bool {
	_, exists := mapType[name]
	return exists
}

// GetAlgo fetches the algorithm
// If the algorithm does not exist, it returns nil
func GetAlgo(name string) (alg PKIAlgo, exists bool) {
	alg, exists = mapType[name]
	return alg, exists
}
