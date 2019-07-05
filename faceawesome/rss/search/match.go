// 匹配器接口模型，凡事要实现匹配逻辑的必须实现matcher接口
// 这是一个抽象接口，有点像模板方法设计模式
// 匹配器返回结果集模型
// 通过种子检索
package search

import "log"

// matcher defines the behavior required by types that want
// to implement a new search type
type Matcher interface {
	Search(feed *Feed, searchTerm string) ([]*Result, error)
}

// result contains the result of a search
type Result struct {
	Field   string
	Content string
}

// match is a launched（启动装置）as a goroutine for each individual feed to run
// searches concurrently
func Match(matcher Matcher, feed *Feed, searchTerm string, results chan<- *Result) {
	// perform the search against the specified matcher
	searchResults, err := matcher.Search(feed, searchTerm)
	if err != nil {
		log.Println(err)
		return
	}

	// write the results to the channel
	for _, result := range searchResults {
		results <- result
	}
}

// display writes results to the console window as they
// are received by the individual goroutines
func Display(results chan *Result) {

	// the channel blocks until a result is written to the channel.
	// once the channel is closed the for loop terminates
	for result := range results {
		log.Printf("%s:\n%s\n\n,", result.Field, result.Content)
	}
}
