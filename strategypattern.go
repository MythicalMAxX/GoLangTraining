// Question: Implement a sorting system where different sorting algorithms (e.g., bubble sort, quick sort, merge sort) can be applied to a list of numbers.
// Scenario: Create a SortStrategy interface and concrete strategy classes for each sorting algorithm, then use a Context class to apply the chosen strategy.

package main

import "fmt"

// SortStrategy interface
type SortStrategy interface {
    sort([]int) []int  // Changed return type to []int
}

// BubbleSort struct
type BubbleSort struct{}

func (bs *BubbleSort) sort(arr []int) []int {
    n := len(arr)
    for i := 0; i < n; i++ {
        for j := 0; j < n-i-1; j++ {
            if arr[j] > arr[j+1] {
                arr[j], arr[j+1] = arr[j+1], arr[j]
            }
        }
    }
    return arr
}

// QuickSort struct
type QuickSort struct{}

func (qs *QuickSort) sort(arr []int) []int {
    QuickSortUtil(arr, 0, len(arr)-1)
    return arr
}

func QuickSortUtil(arr []int, low int, high int) {
    if low < high {
        pivot := partition(arr, low, high)
        QuickSortUtil(arr, low, pivot-1)
        QuickSortUtil(arr, pivot+1, high)
    }
}

func partition(arr []int, low int, high int) int {
    pivot := arr[high]
    i := low - 1
    for j := low; j < high; j++ {
        if arr[j] < pivot {
            i++
            arr[i], arr[j] = arr[j], arr[i]
        }
    }
    arr[i+1], arr[high] = arr[high], arr[i+1]
    return i + 1
}

// MergeSort struct
type MergeSort struct{}

func (ms *MergeSort) sort(arr []int) []int {
    MergeSortUtil(arr, 0, len(arr)-1)
    return arr
}

func MergeSortUtil(arr []int, low int, high int) {
    if low < high {
        mid := low + (high-low)/2
        MergeSortUtil(arr, low, mid)
        MergeSortUtil(arr, mid+1, high)
        merge(arr, low, mid, high)
    }
}

func merge(arr []int, low int, mid int, high int) {
    n1 := mid - low + 1
    n2 := high - mid

    L := make([]int, n1)
    R := make([]int, n2)

    for i := 0; i < n1; i++ {
        L[i] = arr[low+i]
    }

    for j := 0; j < n2; j++ {
        R[j] = arr[mid+1+j]
    }

    i, j := 0, 0
    k := low

    for i < n1 && j < n2 {
        if L[i] <= R[j] {
            arr[k] = L[i]
            i++
        } else {
            arr[k] = R[j]
            j++
        }
        k++
    }

    for i < n1 {
        arr[k] = L[i]
        i++
        k++
    }

    for j < n2 {
        arr[k] = R[j]
        j++
        k++
    }
}

func main() {
    arr := []int{12, 11, 13, 5, 6}
    fmt.Println("Original array:", arr)

    bubbleSort := &BubbleSort{}
    bubbleSort.sort(arr)

    fmt.Println("Sorted array using Bubble Sort:", arr)

    quickSort := &QuickSort{}
    quickSort.sort(arr)

    fmt.Println("Sorted array using Quick Sort:", arr)

    mergeSort := &MergeSort{}
    mergeSort.sort(arr)

    fmt.Println("Sorted array using Merge Sort:", arr)
}