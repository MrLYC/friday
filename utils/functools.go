package utils

// CartesianProductHandler :
type CartesianProductHandler func(interface{}, interface{})

// CartesianProduct :
func CartesianProduct(v1 []interface{}, v2 []interface{}, handler CartesianProductHandler) {
	for _, i := range v1 {
		for _, j := range v2 {
			handler(i, j)
		}
	}
}
