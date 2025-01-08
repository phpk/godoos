package vector

import (
	"cmp"
	"container/heap"
	"context"
	"fmt"
	"runtime"
	"slices"
	"strings"
	"sync"
)

var supportedFilters = []string{"$contains", "$not_contains"}

type docSim struct {
	docID      string
	similarity float32
}

// docMaxHeap 是基于相似度的最大堆。
type docMaxHeap []docSim

func (h docMaxHeap) Len() int           { return len(h) }
func (h docMaxHeap) Less(i, j int) bool { return h[i].similarity < h[j].similarity }
func (h docMaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *docMaxHeap) Push(x any) {
	*h = append(*h, x.(docSim))
}

func (h *docMaxHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// maxDocSims 管理一个固定大小的最大堆，保存最高的 n 个相似度。并发安全，但 values() 返回的结果不是。
type maxDocSims struct {
	h    docMaxHeap
	lock sync.RWMutex
	size int
}

// newMaxDocSims 创建一个新的固定大小的最大堆。
func newMaxDocSims(size int) *maxDocSims {
	return &maxDocSims{
		h:    make(docMaxHeap, 0, size),
		size: size,
	}
}

// add 插入一个新的 docSim 到堆中，保持最高的 n 个相似度。
func (mds *maxDocSims) add(doc docSim) {
	mds.lock.Lock()
	defer mds.lock.Unlock()
	if mds.h.Len() < mds.size {
		heap.Push(&mds.h, doc)
	} else if mds.h.Len() > 0 && mds.h[0].similarity < doc.similarity {
		heap.Pop(&mds.h)
		heap.Push(&mds.h, doc)
	}
}

// values 返回堆中的 docSim，按相似度降序排列。调用是并发安全的，但结果不是。
func (d *maxDocSims) values() []docSim {
	d.lock.RLock()
	defer d.lock.RUnlock()
	slices.SortFunc(d.h, func(i, j docSim) int {
		return cmp.Compare(j.similarity, i.similarity)
	})
	return d.h
}

// filterDocs 并发过滤文档，根据元数据和内容进行筛选。
func filterDocs(docs map[string]*Document, where, whereDocument map[string]string) []*Document {
	filteredDocs := make([]*Document, 0, len(docs))
	var filteredDocsLock sync.Mutex

	numCPUs := runtime.NumCPU()
	numDocs := len(docs)
	concurrency := min(numCPUs, numDocs)

	docChan := make(chan *Document, concurrency*2)

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for doc := range docChan {
				if documentMatchesFilters(doc, where, whereDocument) {
					filteredDocsLock.Lock()
					filteredDocs = append(filteredDocs, doc)
					filteredDocsLock.Unlock()
				}
			}
		}()
	}

	for _, doc := range docs {
		docChan <- doc
	}
	close(docChan)

	wg.Wait()

	if len(filteredDocs) == 0 {
		return nil
	}
	return filteredDocs
}

// documentMatchesFilters 检查文档是否匹配给定的过滤条件。
func documentMatchesFilters(document *Document, where, whereDocument map[string]string) bool {
	for k, v := range where {
		if document.Metadata[k] != v {
			return false
		}
	}

	for k, v := range whereDocument {
		switch k {
		case "$contains":
			if !strings.Contains(document.Content, v) {
				return false
			}
		case "$not_contains":
			if strings.Contains(document.Content, v) {
				return false
			}
		}
	}

	return true
}

// getMostSimilarDocs 获取与查询向量最相似的前 n 个文档。
func getMostSimilarDocs(ctx context.Context, queryVectors []float32, docs []*Document, n int) ([]docSim, error) {
	nMaxDocs := newMaxDocSims(n)

	numCPUs := runtime.NumCPU()
	numDocs := len(docs)
	concurrency := min(numCPUs, numDocs)

	var sharedErr error
	var sharedErrLock sync.Mutex
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	setSharedErr := func(err error) {
		sharedErrLock.Lock()
		defer sharedErrLock.Unlock()
		if sharedErr == nil {
			sharedErr = err
			cancel(sharedErr)
		}
	}

	var wg sync.WaitGroup
	subSliceSize := len(docs) / concurrency
	rem := len(docs) % concurrency

	for i := 0; i < concurrency; i++ {
		start := i * subSliceSize
		end := start + subSliceSize
		if i == concurrency-1 {
			end += rem
		}

		wg.Add(1)
		go func(subSlice []*Document) {
			defer wg.Done()
			for _, doc := range subSlice {
				if ctx.Err() != nil {
					return
				}

				sim, err := dotProduct(queryVectors, doc.Embedding)
				if err != nil {
					setSharedErr(fmt.Errorf("无法计算文档 '%s' 的相似度: %w", doc.ID, err))
					return
				}

				nMaxDocs.add(docSim{docID: doc.ID, similarity: sim})
			}
		}(docs[start:end])
	}

	wg.Wait()

	if sharedErr != nil {
		return nil, sharedErr
	}

	return nMaxDocs.values(), nil
}

// 辅助函数：返回两个数中的最小值。
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
