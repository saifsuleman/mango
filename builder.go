package main

type LoadingCacheBuilder[K comparable, V any] LoadingCache[K, V]

func NewLoadingCacheBuilder[K comparable, V any]() *LoadingCacheBuilder[K, V] {
	return &LoadingCacheBuilder[K, V]{}
}

func (b *LoadingCacheBuilder[K, V]) MaximumSize(maximumSize int64) *LoadingCacheBuilder[K, V] {
	b.maximumSize = maximumSize
	return b
}

func (b *LoadingCacheBuilder[K, V]) MaximumWeight(maximumWeight int64) *LoadingCacheBuilder[K, V] {
	b.maximumWeight = maximumWeight
	return b
}

func (b *LoadingCacheBuilder[K, V]) ExpirationSeconds(expirationSeconds int64) *LoadingCacheBuilder[K, V] {
	b.expirationSeconds = expirationSeconds
	return b
}

func (b *LoadingCacheBuilder[K, V]) WithLoad(loadingFunction LoadingFunction[K, V]) *LoadingCacheBuilder[K, V] {
	b.loadingFunction = loadingFunction
	return b
}

func (b *LoadingCacheBuilder[K, V]) WithWeight(weighingFunction WeighingFunction[K, V]) *LoadingCacheBuilder[K, V] {
	b.weighingFunction = weighingFunction
	return b
}

func (b *LoadingCacheBuilder[K, V]) Build() LoadingCache[K, V] {
	return LoadingCache[K, V]{
		maximumSize:       b.maximumSize,
		maximumWeight:     b.maximumWeight,
		expirationSeconds: b.expirationSeconds,
		loadingFunction:   b.loadingFunction,
		weighingFunction:  b.weighingFunction,
		cache:             make(map[K]*cacheEntry[V]),
	}
}
