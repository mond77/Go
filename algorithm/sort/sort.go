package sort

//冒泡排序
func BubbleSort(nums []int){
	for i := len(nums)-1;i>=0;i--{
		for j :=0 ;j<i;j++{
			if nums[j]>nums[j+1]{
				nums[j],nums[j+1] = nums[j+1],nums[j]
			}
		}
	}
}

//快速排序
func QuickSort1(nums []int)[]int{
	if len(nums) == 0{
		return []int{}
	}
	if len(nums) == 1{
		return []int{nums[0]}
	}
	pivot := nums[0]
	ls,rs := []int{},[]int{}
	for i:= 1;i<len(nums);i++{
		if nums[i]<=pivot{
			ls = append(ls, nums[i])
		}else{
			rs = append(rs, nums[i])
		}
	}
	return append(QuickSort1(ls),append([]int{pivot},QuickSort1(rs)...)...)
}

//快速排序（不需要额外空间）
func QuickSort2(nums []int,l,r int){
	if l>=r{
		return
	}
	idx := partition(nums,l,r)
	QuickSort2(nums,l,idx-1)
	QuickSort2(nums,idx+1,r)
}

func partition(nums []int,l,r int)int{

	pivot := nums[r]
	for l<r{
		for l<r && nums[l]<=pivot{
			l++
		}
		if nums[l]>pivot{
			nums[r] = nums[l]
		}
		for l<r && nums[r] >pivot{
			r--
		}
		if nums[r]<=pivot{
			nums[l] = nums[r]
		}
	}
	nums[l] = pivot
	return l
}

//希尔排序
func ShellSort(arr []int) []int{
	n := len(arr)
	gap := n/2
	for gap>0{
		for i:= gap;i<n;i++{
			tmp := arr[i]
			pre := i-gap
			for pre>= 0 && arr[pre]>tmp{
				arr[pre+gap] = arr[pre]
				pre-= gap
			}
			arr[pre+gap] = tmp
		}
		gap /=2
	} 
	return arr
}

