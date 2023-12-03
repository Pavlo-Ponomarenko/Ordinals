package ord_math

func findRelation(index uint32, outputValues []int64, inputValues []int64) uint32 {
	outputIndex := 0
	inputIndex := uint32(0)
	differ := int64(0)
	for uint32(outputIndex) < index {
		if differ == 0 {
			differ = inputValues[inputIndex] - outputValues[outputIndex]
			if inputValues[inputIndex] > outputValues[outputIndex] {
				outputIndex++
			} else if inputValues[inputIndex] == outputValues[outputIndex] {
				outputIndex++
				inputIndex++
			} else {
				inputIndex++
			}
		} else if differ > 0 {
			oldOutputIndex := outputIndex
			if differ > outputValues[outputIndex] {
				outputIndex++
			} else if differ == outputValues[outputIndex] {
				outputIndex++
				inputIndex++
			} else {
				inputIndex++
			}
			differ = differ - outputValues[oldOutputIndex]
		} else {
			oldInputIndex := inputIndex
			if inputValues[inputIndex] > differ*-1 {
				outputIndex++
			} else if inputValues[inputIndex] == differ*-1 {
				outputIndex++
				inputIndex++
			} else {
				inputIndex++
			}
			differ = inputValues[oldInputIndex] + differ
		}
	}
	return inputIndex
}

func calcCurrentOutputOrder(outputIndex uint32, inputIndex uint32, outputValues []int64, inputValues []int64) int64 {
	beforeInputCoins := int64(0)
	for i := uint32(0); i < inputIndex; i++ {
		beforeInputCoins += inputValues[i]
	}
	beforeOutputCoins := int64(0)
	for i := uint32(0); i < outputIndex; i++ {
		beforeOutputCoins += outputValues[i]
	}
	return beforeOutputCoins - beforeInputCoins
}

func firstOrdinal(height int64) int64 {
	start := int64(0)
	for h := int64(0); h < height; h++ {
		start += subsidy(h)
	}
	return start
}

func subsidy(height int64) int64 {
	return int64(50*100000000) >> (uint64(height) / 210000)
}
