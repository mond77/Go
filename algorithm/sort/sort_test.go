package sort

import (
	"fmt"
	"testing"
)

func TestBubbleSort(t *testing.T) {

	nums := []int{1,9,2,8,3,7,6,4,5,112,102,1,9,2,8,3,7,6,4,5,112,102,1,9,2,8,3,7,6,4,5,112,102,1,9,2,8,3,7,6,4,5,112,102,1,9,2,8,3,7,6,4,5,112,102,1,9,2,8,3,7,6,4,5,112}
	BubbleSort(nums)
	fmt.Println(nums)

}

func TestQuickSort1(t *testing.T) {

	nums := []int{1,9,2,8,3,7,6,4,5,112,102,1,9,2,8,3,7,6,4,5,112,102,1,9,2,8,3,7,6,4,5,112,102,1,9,2,8,3,7,6,4,5,112,102,1,9,2,8,3,7,6,4,5,112,102,1,9,2,8,3,7,6,4,5,112}
	nums = QuickSort1(nums)
	fmt.Println(nums)
}

func TestQuickSort2(t *testing.T) {

	nums := []int{1,9,2,8,3,7,6,4,5,112,102,1,9,2,8,3,7,6,4,5,112,102,1,9,2,8,3,7,6,4,5,112,102,1,9,2,8,3,7,6,4,5,112,102,1,9,2,8,3,7,6,4,5,112,102,1,9,2,8,3,7,6,4,5,112}
	QuickSort2(nums,0,len(nums)-1)
	fmt.Println(nums)
}