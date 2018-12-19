package main

type Vector2 [2]int

func (v Vector2) Add(v2 Vector2) Vector2 {
	return Vector2{v[0] + v2[0], v[1] + v2[1]}
}

func (v Vector2) Sub(v2 Vector2) Vector2 {
	return Vector2{v[0] - v2[0], v[1] - v2[1]}
}

type VectorList []Vector2

func (vl VectorList) Less(i, j int) bool {
	if vl[i][1] < vl[j][1] {
		return true
	}
	if vl[i][1] == vl[j][1] && vl[i][0] < vl[j][0] {
		return true
	}
	return false
}

func (vl VectorList) Swap(i, j int) {
	vl[i], vl[j] = vl[j], vl[i]
}

func (vl VectorList) Len() int {
	return len(vl)
}
