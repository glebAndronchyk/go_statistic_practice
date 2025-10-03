package interfaces

type IDistribution interface {
	GetVariants() []int
	GetOccurences() []int
}
