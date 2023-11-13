package slice

// PageOfSliceStr - pagination for slice
func PageOfSliceStr(s []string, page, perPage uint32) []string {
	if page == 0 || perPage == 0 {
		return []string{}
	}
	offset := (page - 1) * perPage
	if int(offset) > len(s) {
		return []string{}
	}
	end := offset + perPage
	if int(end) > len(s) {
		end = uint32(len(s))
	}
	return s[offset:end]
}

// SplitSliceStr []string -> [][]string
func SplitSliceStr(s []string, batchSize int) [][]string {
	if batchSize < 1 {
		return [][]string{}
	}
	var batchesCount int
	if len(s)%batchSize == 0 {
		batchesCount = len(s) / batchSize
	} else {
		batchesCount = len(s)/batchSize + 1
	}
	batches := make([][]string, 0, batchesCount)
	for i := 0; i < batchesCount; i++ {
		batches = append(batches, make([]string, 0, batchSize))
		for j := 0; j < batchSize; j++ {
			index := i*batchSize + j
			if index == len(s) {
				break
			}
			batches[i] = append(batches[i], s[i*batchSize+j])
		}
	}
	return batches
}
