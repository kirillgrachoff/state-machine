package mask_acceptor

import "unsafe"

type MaskAcceptor struct {
	data map[uint]bool
}

func (acc *MaskAcceptor) In(mask uint) bool {
	for _, value := range FromMask(mask) {
		if acc.data[value] {
			return true
		}
	}
	return false
}

func New(positive []uint) *MaskAcceptor {
	ans := &MaskAcceptor{make(map[uint]bool)}
	for _, value := range positive {
		ans.data[value] = true
	}
	return ans
}

func FromMask(mask uint) []uint {
	ans := make([]uint, 0)
	for i := uint(0); i < 8*uint(unsafe.Sizeof(mask)); i++ {
		if (mask & (1 << i)) != 0 {
			ans = append(ans, i)
		}
	}
	return ans
}

func ToMask(source []uint) uint {
	ans := uint(0)
	for _, value := range source {
		ans |= uint(1) << value
	}
	return ans
}
