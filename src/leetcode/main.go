package main

import (
    "fmt"
)

// https://leetcode-cn.com/problems/regular-expression-matching/
func isMatch(s string, p string) bool {
    s_len := len(s)
    p_len := len(p)
    if p_len == 0 {
        return s_len == 0
    }

    if p_len >= 2 {
        if p[1] == '*' {
            for i := 0; i < s_len; i++ {
                if isMatch(s[i:], p[2:]) {
                    return true
                }

                if p[0] != '.' && s[i] != p[0] {
                    return false
                }
            }
            return isMatch(s[s_len:], p[2:])
        }
    }

    if s_len == 0 || (p[0] != '.' && p[0] != s[0]){
        return false
    }

    return isMatch(s[1:], p[1:])
}

// m个不同整数中找出n个整数的所有组合
func gen_m_num (m int, num *[]int) {
    for i:=1; i <= m; i++ {
        *num = append(*num, i)
    }
}
func gen_n_combin (n int, pre []int, input []int) {
    switch {
    case n == 1:
        for _, v := range input {
            fmt.Println(append(pre, v))
        }
    case len(input) > n:
        for i := 1; i <= len(input) - (n-1); i++ {
            new_pre := append(pre, input[i-1])
            gen_n_combin(n-1, new_pre, input[i:])
        }
    case len(input) == n:
        fmt.Println(append(pre, input...))
    }
}

// https://leetcode-cn.com/problems/add-two-numbers/
type ListNode struct {
    Val int
    Next *ListNode
}
func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
    switch {
    case l1 == nil && l2 == nil:
        return nil
    case l1 == nil:
        l1 = &ListNode{Val:0}
    case l2 == nil:
        l2 = &ListNode{Val:0}
    }

    l1.Val += l2.Val
    if l1.Val >= 10 {
        a := l1.Val / 10
        l1.Val %= 10
        if l1.Next == nil {
            l1.Next = &ListNode{Val: a}
        } else {
            l1.Next.Val += a
        }
    }

    l1.Next = addTwoNumbers(l1.Next, l2.Next)
    return l1
}

// https://leetcode-cn.com/problems/two-sum/
func twoSum(nums []int, target int) []int {
    for i, v := range nums {
        for j, v2 := range nums[i+1:] {
            if v + v2 == target {
                return []int{i, j+i+1}
            }
        }
    }
    return nil
}

// https://leetcode-cn.com/problems/remove-nth-node-from-end-of-list/
func removeNthFromEnd(head *ListNode, n int) *ListNode {
    first := head
    for i := 0; i < n; i++ {
        first = first.Next
    }
    pre_head := &ListNode{Val:0, Next:head}
    second := pre_head
    for first != nil {
        second = second.Next
        first = first.Next
    }
    second.Next = second.Next.Next
    return pre_head.Next
}

// https://leetcode-cn.com/problems/longest-valid-parentheses/
type Range struct {
    Begin int
    Len int
}
func longestValidParentheses(s string) int {
    var stack []int
    var r []Range
    var max int
    for i, v := range s {
        if v == '(' {
            stack = append(stack, i)
        } else {
            if len(stack) > 0 {
                begin := stack[len(stack)-1]
                size := i + 1 - begin

                j := len(r)
                for j > 0 && begin < r[j-1].Begin {
                   j--
                }

                if j == 0 {
                    r = append(r, Range{Begin: begin, Len: size})
                } else {
                    if r[j-1].Begin + r[j-1].Len == begin {
                        r = r[:j]
                        r[j-1].Len += size
                    } else {
                        r = r[:j]
                        r = append(r, Range{Begin: begin, Len: size})
                    }
                }

                if r[len(r)-1].Len > max {
                    max = r[len(r)-1].Len
                }
                stack = stack[:len(stack)-1]
            }
        }
    }
    return max
}

// https://leetcode-cn.com/problems/substring-with-concatenation-of-all-words/
func findSubstring(s string, words []string) []int {
    var result []int
    words_count := len(words)
    if words_count == 0 {
        return result
    }
    words_len := len(words[0])
    if words_len == 0 {
        return result
    }

    words_map := make(map[string]int)
    for _, v := range words {
        words_map[v] += 1
    }
    for i := 0; i < words_len; i++ {
        var temp_map = make(map[string]int)
        var match_count int
        first_word_index := i
        for j := i; j + words_len <= len(s); j += words_len {
            cur_s := s[j:j+words_len]
            if v, key_exist := words_map[cur_s]; key_exist {
                if temp_map[cur_s] < v {
                    temp_map[cur_s] += 1
                    match_count++
                    if match_count == words_count {
                        result = append(result, first_word_index)

                        temp_map[s[first_word_index:first_word_index+words_len]] -= 1
                        match_count--
                        first_word_index += words_len
                    }
                } else {
                    for k := first_word_index; k < j; k += words_len {
                        if cur_s == s[k:k+words_len] {
                            first_word_index = k + words_len
                            break
                        }
                        temp_map[s[k:k+words_len]] -= 1
                        match_count--
                    }
                }
            } else {
                temp_map = make(map[string]int)
                first_word_index = j + words_len
                match_count = 0
            }
        }
    }
    return result
}

func main() {
    var m int
    var n int
    fmt.Println("find n combin from m")
    fmt.Println("input m and n(m >= n):")
    fmt.Scanln(&m, &n)
    list := make([]int, 0)
    gen_m_num(m, &list)
    fmt.Println("gen num:", list)
    fmt.Println("result:")
    gen_n_combin(n, make([]int, 0), list)

    fmt.Println("##############################################################")
    fmt.Println("start test.....................")
    if isMatch("aab", "c*a*b") {
        fmt.Println("match")
    } else {
        fmt.Println("not match")
    }
    fmt.Println(twoSum([]int{2, 7, 11, 15}, 9))
    fmt.Println(longestValidParentheses(")()())"))
    fmt.Println(longestValidParentheses("()(())"))
    fmt.Println(findSubstring("barfoothefoobarman", []string{"foo", "bar"}))
    fmt.Println(findSubstring("wordgoodgoodgoodbestword", []string{"word","good","best","word"}))
    fmt.Println(findSubstring("wordgoodgoodgoodbestword", []string{"word","good","best","good"}))
    fmt.Println(findSubstring("barfoofoobarthefoobarman", []string{"bar","foo","the"}))
    fmt.Println("finish test.....................")
}
